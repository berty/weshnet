package weshnet

import (
	"fmt"

	cid "github.com/ipfs/go-cid"
	"golang.org/x/crypto/nacl/secretbox"
	"google.golang.org/protobuf/proto"

	ipfslog "berty.tech/go-ipfs-log"
	"berty.tech/weshnet/v2/pkg/cryptoutil"
	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
)

var eventTypesMapper = map[protocoltypes.EventType]struct {
	Message    proto.Message
	SigChecker sigChecker
}{
	protocoltypes.EventType_EventTypeGroupMemberDeviceAdded:                 {Message: &protocoltypes.GroupMemberDeviceAdded{}, SigChecker: sigCheckerGroupMemberDeviceAdded},
	protocoltypes.EventType_EventTypeGroupDeviceChainKeyAdded:               {Message: &protocoltypes.GroupDeviceChainKeyAdded{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountGroupJoined:                     {Message: &protocoltypes.AccountGroupJoined{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountGroupLeft:                       {Message: &protocoltypes.AccountGroupLeft{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountContactRequestDisabled:          {Message: &protocoltypes.AccountContactRequestDisabled{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountContactRequestEnabled:           {Message: &protocoltypes.AccountContactRequestEnabled{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountContactRequestReferenceReset:    {Message: &protocoltypes.AccountContactRequestReferenceReset{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountContactRequestOutgoingEnqueued:  {Message: &protocoltypes.AccountContactRequestOutgoingEnqueued{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountContactRequestOutgoingSent:      {Message: &protocoltypes.AccountContactRequestOutgoingSent{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountContactRequestIncomingReceived:  {Message: &protocoltypes.AccountContactRequestIncomingReceived{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountContactRequestIncomingDiscarded: {Message: &protocoltypes.AccountContactRequestIncomingDiscarded{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountContactRequestIncomingAccepted:  {Message: &protocoltypes.AccountContactRequestIncomingAccepted{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountContactBlocked:                  {Message: &protocoltypes.AccountContactBlocked{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountContactUnblocked:                {Message: &protocoltypes.AccountContactUnblocked{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeContactAliasKeyAdded:                   {Message: &protocoltypes.ContactAliasKeyAdded{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeMultiMemberGroupAliasResolverAdded:     {Message: &protocoltypes.MultiMemberGroupAliasResolverAdded{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeMultiMemberGroupInitialMemberAnnounced: {Message: &protocoltypes.MultiMemberGroupInitialMemberAnnounced{}, SigChecker: sigCheckerGroupSigned},
	protocoltypes.EventType_EventTypeMultiMemberGroupAdminRoleGranted:       {Message: &protocoltypes.MultiMemberGroupAdminRoleGranted{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeGroupMetadataPayloadSent:               {Message: &protocoltypes.GroupMetadataPayloadSent{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeGroupReplicating:                       {Message: &protocoltypes.GroupReplicating{}, SigChecker: sigCheckerDeviceSigned},
	protocoltypes.EventType_EventTypeAccountVerifiedCredentialRegistered:    {Message: &protocoltypes.AccountVerifiedCredentialRegistered{}, SigChecker: sigCheckerDeviceSigned},
}

func newEventContext(eventID cid.Cid, parentIDs []cid.Cid, g *protocoltypes.Group) *protocoltypes.EventContext {
	parentIDsBytes := make([][]byte, len(parentIDs))
	for i, parentID := range parentIDs {
		parentIDsBytes[i] = parentID.Bytes()
	}

	return &protocoltypes.EventContext{
		Id:        eventID.Bytes(),
		ParentIds: parentIDsBytes,
		GroupPk:   g.PublicKey,
	}
}

// FIXME(gfanton): getParentsCID use a lot of resources
// nolint:unused
func getParentsForCID(log ipfslog.Log, c cid.Cid) []cid.Cid {
	if log == nil {
		// TODO: this should not happen
		return []cid.Cid{}
	}

	parent, ok := log.Get(c)

	// Can't fetch parent entry
	if !ok {
		return []cid.Cid{}
	}

	nextEntries := parent.GetNext()

	// Parent has only one or no parents, returning its id
	if len(nextEntries) <= 1 {
		return []cid.Cid{parent.GetHash()}
	}

	// Parent has more than one parent, returning parent entries
	var ret []cid.Cid
	for _, n := range nextEntries {
		ret = append(ret, getParentsForCID(log, n)...)
	}

	return ret
}

func newGroupMetadataEventFromEntry(_ ipfslog.Log, e ipfslog.Entry, metadata *protocoltypes.GroupMetadata, event proto.Message, g *protocoltypes.Group) (*protocoltypes.GroupMetadataEvent, error) {
	// TODO: if parent is a merge node we should return the next nodes of it

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization
	}

	// TODO(gfanton): getParentsCID use a lot of resources, disable it until we need it
	// evtCtx := newEventContext(e.GetHash(), getParentsForCID(log, e.GetHash()), group, attachmentsCIDs)
	evtCtx := newEventContext(e.GetHash(), []cid.Cid{}, g)

	gme := protocoltypes.GroupMetadataEvent{
		EventContext: evtCtx,
		Metadata:     metadata,
		Event:        eventBytes,
	}

	return &gme, nil
}

func openGroupEnvelope(g *protocoltypes.Group, envelopeBytes []byte) (*protocoltypes.GroupMetadata, proto.Message, error) {
	env := &protocoltypes.GroupEnvelope{}
	if err := proto.Unmarshal(envelopeBytes, env); err != nil {
		return nil, nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	nonce, err := cryptoutil.NonceSliceToArray(env.Nonce)
	if err != nil {
		return nil, nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	data, ok := secretbox.Open(nil, env.Event, nonce, g.GetSharedSecret())
	if !ok {
		return nil, nil, errcode.ErrCode_ErrGroupMemberLogEventOpen
	}

	metadataEvent := &protocoltypes.GroupMetadata{}

	err = proto.Unmarshal(data, metadataEvent)
	if err != nil {
		return nil, nil, errcode.ErrCode_TODO.Wrap(err)
	}

	et, ok := eventTypesMapper[metadataEvent.EventType]
	if !ok {
		return nil, nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("event type not found"))
	}

	payload := proto.Clone(et.Message)
	if err := proto.Unmarshal(metadataEvent.Payload, payload); err != nil {
		return nil, nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if err := et.SigChecker(g, metadataEvent, payload); err != nil {
		return nil, nil, errcode.ErrCode_ErrCryptoSignatureVerification.Wrap(err)
	}

	return metadataEvent, payload, nil
}

func sealGroupEnvelope(g *protocoltypes.Group, eventType protocoltypes.EventType, payload proto.Message, payloadSig []byte) ([]byte, error) {
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, errcode.ErrCode_TODO.Wrap(err)
	}

	nonce, err := cryptoutil.GenerateNonce()
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoNonceGeneration.Wrap(err)
	}

	event := &protocoltypes.GroupMetadata{
		EventType:        eventType,
		Payload:          payloadBytes,
		Sig:              payloadSig,
		ProtocolMetadata: &protocoltypes.ProtocolMetadata{},
	}

	eventClearBytes, err := proto.Marshal(event)
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	eventBytes := secretbox.Seal(nil, eventClearBytes, nonce, g.GetSharedSecret())

	env := &protocoltypes.GroupEnvelope{
		Event: eventBytes,
		Nonce: nonce[:],
	}

	return proto.Marshal(env)
}
