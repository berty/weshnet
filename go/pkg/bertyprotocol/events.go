package bertyprotocol

import (
	"crypto/rand"

	ipfslog "berty.tech/go-ipfs-log"
	"github.com/golang/protobuf/proto"
	"github.com/ipfs/go-cid"
	"golang.org/x/crypto/nacl/secretbox"

	"berty.tech/berty/go/pkg/errcode"
)

var eventTypesMapper = map[EventType]struct {
	Message    proto.Message
	SigChecker SigChecker
}{
	EventTypeGroupMemberDeviceAdded:                 {Message: &GroupAddMemberDevice{}, SigChecker: SigCheckerMemberDeviceAdded},
	EventTypeGroupDeviceSecretAdded:                 {Message: &GroupAddDeviceSecret{}, SigChecker: SigCheckerDeviceSigned},
	EventTypeAccountGroupJoined:                     {Message: &AccountGroupJoined{}, SigChecker: SigCheckerMissing},
	EventTypeAccountGroupLeft:                       {Message: &AccountGroupLeft{}, SigChecker: SigCheckerMissing},
	EventTypeAccountContactRequestDisabled:          {Message: &AccountContactRequestDisabled{}, SigChecker: SigCheckerMissing},
	EventTypeAccountContactRequestEnabled:           {Message: &AccountContactRequestEnabled{}, SigChecker: SigCheckerMissing},
	EventTypeAccountContactRequestReferenceReset:    {Message: &AccountContactRequestReferenceReset{}, SigChecker: SigCheckerMissing},
	EventTypeAccountContactRequestOutgoingEnqueued:  {Message: &AccountContactRequestEnqueued{}, SigChecker: SigCheckerMissing},
	EventTypeAccountContactRequestOutgoingSent:      {Message: &AccountContactRequestSent{}, SigChecker: SigCheckerMissing},
	EventTypeAccountContactRequestIncomingReceived:  {Message: &AccountContactRequestReceived{}, SigChecker: SigCheckerMissing},
	EventTypeAccountContactRequestIncomingDiscarded: {Message: &AccountContactRequestDiscarded{}, SigChecker: SigCheckerMissing},
	EventTypeAccountContactRequestIncomingAccepted:  {Message: &AccountContactRequestAccepted{}, SigChecker: SigCheckerMissing},
	EventTypeAccountContactBlocked:                  {Message: &AccountContactBlocked{}, SigChecker: SigCheckerMissing},
	EventTypeAccountContactUnblocked:                {Message: &AccountContactUnblocked{}, SigChecker: SigCheckerMissing},
	EventTypeContactAliasKeyAdded:                   {Message: &ContactAddAliasKey{}, SigChecker: SigCheckerMissing},
	EventTypeMultiMemberGroupAliasResolverAdded:     {Message: &MultiMemberGroupAddAliasResolver{}, SigChecker: SigCheckerMissing},
	EventTypeMultiMemberGroupInitialMemberAnnounced: {Message: &MultiMemberInitialMember{}, SigChecker: SigCheckerGroupSigned},
	EventTypeMultiMemberGroupAdminRoleGranted:       {Message: &MultiMemberGrantAdminRole{}, SigChecker: SigCheckerMissing},
	EventTypeGroupMetadataPayloadSent:               {Message: &AppMetadata{}, SigChecker: SigCheckerMissing},
}

func NewEventContext(eventID cid.Cid, parentIDs []cid.Cid, g *Group) (*EventContext, error) {
	parentIDsBytes := make([][]byte, len(parentIDs))
	for i, parentID := range parentIDs {
		parentIDsBytes[i] = parentID.Bytes()
	}

	return &EventContext{
		ID:        eventID.Bytes(),
		ParentIDs: parentIDsBytes,
		GroupPK:   g.PublicKey,
	}, nil
}

func getParentsForCID(log ipfslog.Log, c cid.Cid) []cid.Cid {
	if log == nil {
		// TODO: this should not happen
		return []cid.Cid{}
	}

	parent, ok := log.GetEntries().Get(c.String())

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

func NewGroupMetadataEventFromEntry(log ipfslog.Log, e ipfslog.Entry, metadata *GroupMetadata, event proto.Message, g *Group) (*GroupMetadataEvent, error) {
	// TODO: if parent is a merge node we should return the next nodes of it
	evtCtx, err := NewEventContext(e.GetHash(), getParentsForCID(log, e.GetHash()), g)
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return nil, errcode.ErrSerialization
	}

	return &GroupMetadataEvent{
		EventContext: evtCtx,
		Metadata:     metadata,
		Event:        eventBytes,
	}, nil
}

func OpenGroupEnvelope(g *Group, envelopeBytes []byte) (*GroupMetadata, proto.Message, error) {
	sharedSecret, err := g.GetSharedSecret()
	if err != nil {
		return nil, nil, errcode.TODO.Wrap(err)
	}

	env := &GroupEnvelope{}
	if err := env.Unmarshal(envelopeBytes); err != nil {
		return nil, nil, errcode.ErrInvalidInput.Wrap(err)
	}

	if len(env.Nonce) != 24 {
		return nil, nil, errcode.ErrInvalidInput
	}

	var nonce [24]byte
	copy(nonce[:], env.Nonce[:24])

	data, ok := secretbox.Open(nil, env.Event, &nonce, sharedSecret)
	if !ok {
		return nil, nil, errcode.ErrGroupMemberLogEventOpen
	}

	metadataEvent := &GroupMetadata{}

	err = metadataEvent.Unmarshal(data)
	if err != nil {
		return nil, nil, errcode.TODO.Wrap(err)
	}

	et, ok := eventTypesMapper[metadataEvent.EventType]
	if !ok {
		return nil, nil, errcode.ErrInvalidInput
	}

	payload := proto.Clone(et.Message)
	if err := proto.Unmarshal(metadataEvent.Payload, payload); err != nil {
		return nil, nil, errcode.ErrDeserialization.Wrap(err)
	}

	if err := et.SigChecker(g, metadataEvent, payload); err != nil {
		return nil, nil, errcode.ErrSignatureVerificationFailed.Wrap(err)
	}

	return metadataEvent, payload, nil
}

func SealGroupEnvelope(g *Group, eventType EventType, payload proto.Marshaler, payloadSig []byte) ([]byte, error) {
	payloadBytes, err := payload.Marshal()
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	var nonce [24]byte
	nonceArr := make([]byte, 24)
	if n, err := rand.Read(nonceArr); err != nil || n != 24 {
		return nil, errcode.ErrRandomGenerationFailed.Wrap(err)
	}

	for i := range nonceArr {
		nonce[i] = nonceArr[i]
	}

	sharedSecret, err := g.GetSharedSecret()
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	event := &GroupMetadata{
		EventType: eventType,
		Payload:   payloadBytes,
		Sig:       payloadSig,
	}

	eventClearBytes, err := event.Marshal()
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	eventBytes := secretbox.Seal(nil, eventClearBytes, &nonce, sharedSecret)

	env := &GroupEnvelope{
		Event: eventBytes,
		Nonce: nonceArr,
	}

	return env.Marshal()
}
