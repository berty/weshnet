package weshnet

import (
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"slices"
	"strings"

	coreiface "github.com/ipfs/kubo/core/coreiface"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/p2p/host/eventbus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	ipfslog "berty.tech/go-ipfs-log"
	"berty.tech/go-ipfs-log/identityprovider"
	ipliface "berty.tech/go-ipfs-log/iface"
	"berty.tech/go-orbit-db/address"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/stores"
	"berty.tech/go-orbit-db/stores/basestore"
	"berty.tech/go-orbit-db/stores/operation"
	"berty.tech/weshnet/v2/pkg/errcode"
	"berty.tech/weshnet/v2/pkg/protocoltypes"
	"berty.tech/weshnet/v2/pkg/secretstore"
	"berty.tech/weshnet/v2/pkg/tyber"
)

type MetadataStore struct {
	basestore.BaseStore
	eventBus event.Bus
	emitters struct {
		groupMetadata    event.Emitter
		metadataReceived event.Emitter
	}

	group              *protocoltypes.Group
	memberDevice       secretstore.OwnMemberDevice
	devicePublicKeyRaw []byte
	secretStore        secretstore.SecretStore
	logger             *zap.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

func isMultiMemberGroup(m *MetadataStore) bool {
	return m.group.GroupType == protocoltypes.GroupType_GroupTypeMultiMember
}

func isAccountGroup(m *MetadataStore) bool {
	return m.group.GroupType == protocoltypes.GroupType_GroupTypeAccount
}

func isContactGroup(m *MetadataStore) bool {
	return m.group.GroupType == protocoltypes.GroupType_GroupTypeContact
}

func (m *MetadataStore) typeChecker(types ...func(m *MetadataStore) bool) bool {
	for _, t := range types {
		if t(m) {
			return true
		}
	}

	return false
}

func (m *MetadataStore) setLogger(l *zap.Logger) {
	if l == nil {
		return
	}

	// m.logger = l.Named("store").With(logutil.PrivateString("group-id", fmt.Sprintf("%.6s", base64.StdEncoding.EncodeToString(m.group.PublicKey))))
	m.logger = l.Named("metastore")

	if index, ok := m.Index().(loggable); ok {
		index.setLogger(m.logger)
	}
}

func openMetadataEntry(log ipfslog.Log, e ipfslog.Entry, g *protocoltypes.Group) (*protocoltypes.GroupMetadataEvent, proto.Message, error) {
	op, err := operation.ParseOperation(e)
	if err != nil {
		return nil, nil, err
	}

	meta, event, err := openGroupEnvelope(g, op.GetValue())
	if err != nil {
		return nil, nil, err
	}

	metaEvent, err := newGroupMetadataEventFromEntry(log, e, meta, event, g)
	if err != nil {
		return nil, nil, err
	}

	return metaEvent, event, err
}

// not used
// func (m *MetadataStore) openMetadataEntry(e ipfslog.Entry) (*protocoltypes.GroupMetadataEvent, proto.Message, error) {
// 	return openMetadataEntry(m.OpLog(), e, m.group, m.devKS)
// }

// FIXME: use iterator instead to reduce resource usage (require go-ipfs-log improvements)
func (m *MetadataStore) ListEvents(_ context.Context, since, until []byte, reverse bool) (<-chan *protocoltypes.GroupMetadataEvent, error) {
	entries, err := getEntriesInRange(m.OpLog().GetEntries().Reverse().Slice(), since, until)
	if err != nil {
		return nil, err
	}

	out := make(chan *protocoltypes.GroupMetadataEvent)

	go func() {
		iterateOverEntries(
			entries,
			reverse,
			func(entry ipliface.IPFSLogEntry) {
				event, _, err := openMetadataEntry(m.OpLog(), entry, m.group)
				if err != nil {
					m.logger.Error("unable to open metadata event", zap.Error(err))
				} else {
					out <- event
					m.logger.Info("metadata store - sent 1 event from log history")
				}
			},
		)

		close(out)
	}()

	return out, nil
}

func (m *MetadataStore) AddDeviceToGroup(ctx context.Context) (operation.Operation, error) {
	md, err := m.secretStore.GetOwnMemberDeviceForGroup(m.group)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return MetadataStoreAddDeviceToGroup(ctx, m, m.group, md)
}

func MetadataStoreAddDeviceToGroup(ctx context.Context, m *MetadataStore, g *protocoltypes.Group, md secretstore.OwnMemberDevice) (operation.Operation, error) {
	device, err := md.Device().Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	member, err := md.Member().Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	k, err := m.GetMemberByDevice(md.Device())
	if err == nil && k != nil {
		return nil, nil
	}

	memberSig, err := md.MemberSign(device)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoSignature.Wrap(err)
	}

	event := &protocoltypes.GroupMemberDeviceAdded{
		MemberPk:  member,
		DevicePk:  device,
		MemberSig: memberSig,
	}

	sig, err := signProtoWithDevice(event, md)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoSignature.Wrap(err)
	}

	m.logger.Info("announcing device on store")

	return metadataStoreAddEvent(ctx, m, g, protocoltypes.EventType_EventTypeGroupMemberDeviceAdded, event, sig)
}

func (m *MetadataStore) SendSecret(ctx context.Context, memberPK crypto.PubKey) (operation.Operation, error) {
	ok, err := m.Index().(*metadataStoreIndex).areSecretsAlreadySent(memberPK)
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	if ok {
		return nil, errcode.ErrCode_ErrGroupSecretAlreadySentToMember
	}

	if devs, err := m.GetDevicesForMember(memberPK); len(devs) == 0 || err != nil {
		m.logger.Warn("sending secret to an unknown group member")
	}

	encryptedSecret, err := m.secretStore.GetShareableChainKey(ctx, m.group, memberPK)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoEncrypt.Wrap(err)
	}

	return MetadataStoreSendSecret(ctx, m, m.group, m.memberDevice, memberPK, encryptedSecret)
}

func MetadataStoreSendSecret(ctx context.Context, m *MetadataStore, g *protocoltypes.Group, md secretstore.OwnMemberDevice, memberPK crypto.PubKey, encryptedSecret []byte) (operation.Operation, error) {
	devicePKRaw, err := md.Device().Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	memberPKRaw, err := memberPK.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	event := &protocoltypes.GroupDeviceChainKeyAdded{
		DevicePk:     devicePKRaw,
		DestMemberPk: memberPKRaw,
		Payload:      encryptedSecret,
	}

	sig, err := signProtoWithDevice(event, md)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoSignature.Wrap(err)
	}

	return metadataStoreAddEvent(ctx, m, g, protocoltypes.EventType_EventTypeGroupDeviceChainKeyAdded, event, sig)
}

func (m *MetadataStore) ClaimGroupOwnership(ctx context.Context, groupSK crypto.PrivKey) (operation.Operation, error) {
	if !m.typeChecker(isMultiMemberGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	event := &protocoltypes.MultiMemberGroupInitialMemberAnnounced{
		MemberPk: m.devicePublicKeyRaw,
	}

	sig, err := signProtoWithPrivateKey(event, groupSK)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoSignature.Wrap(err)
	}

	return metadataStoreAddEvent(ctx, m, m.group, protocoltypes.EventType_EventTypeMultiMemberGroupInitialMemberAnnounced, event, sig)
}

func signProtoWithDevice(message proto.Message, memberDevice secretstore.OwnMemberDevice) ([]byte, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	sig, err := memberDevice.DeviceSign(data)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoSignature.Wrap(err)
	}

	return sig, nil
}

func signProtoWithPrivateKey(message proto.Message, sk crypto.PrivKey) ([]byte, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	sig, err := sk.Sign(data)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoSignature.Wrap(err)
	}

	return sig, nil
}

func metadataStoreAddEvent(ctx context.Context, m *MetadataStore, g *protocoltypes.Group, eventType protocoltypes.EventType, event proto.Message, sig []byte) (operation.Operation, error) {
	ctx, newTrace := tyber.ContextWithTraceID(ctx)
	tyberLogError := tyber.LogError
	if newTrace {
		m.logger.Debug(fmt.Sprintf("Sending %s to %s group %s", strings.TrimPrefix(eventType.String(), "EventType"), strings.TrimPrefix(g.GroupType.String(), "GroupType"), base64.RawURLEncoding.EncodeToString(g.PublicKey)), tyber.FormatTraceLogFields(ctx)...)
		tyberLogError = tyber.LogFatalError
	}

	env, err := sealGroupEnvelope(g, eventType, event, sig)
	if err != nil {
		return nil, tyberLogError(ctx, m.logger, "Failed to seal group envelope", errcode.ErrCode_ErrCryptoSignature.Wrap(err))
	}
	m.logger.Debug(fmt.Sprintf("Sealed group envelope (%d bytes)", len(env)), tyber.FormatStepLogFields(ctx, []tyber.Detail{})...)

	op := operation.NewOperation(nil, "ADD", env)
	e, err := m.AddOperation(ctx, op, nil)
	if err != nil {
		return nil, tyberLogError(ctx, m.logger, "Failed to add operation on log", errcode.ErrCode_ErrOrbitDBAppend.Wrap(err))
	}
	m.logger.Debug("Added operation on log", tyber.FormatStepLogFields(ctx, []tyber.Detail{
		{Name: "CID", Description: e.GetHash().String()},
	})...)

	op, err = operation.ParseOperation(e)
	if err != nil {
		return nil, tyberLogError(ctx, m.logger, "Failed to parse operation returned by log", errcode.ErrCode_ErrOrbitDBDeserialization.Wrap(err))
	}

	if newTrace {
		m.logger.Debug("Added metadata on log successfully", tyber.FormatStepLogFields(ctx, []tyber.Detail{}, tyber.EndTrace)...)
	}
	return op, nil
}

func (m *MetadataStore) ListContacts() map[string]*AccountContact {
	return m.Index().(*metadataStoreIndex).listContacts()
}

func (m *MetadataStore) ListVerifiedCredentials() []*protocoltypes.AccountVerifiedCredentialRegistered {
	return m.Index().(*metadataStoreIndex).listVerifiedCredentials()
}

func (m *MetadataStore) GetMemberByDevice(pk crypto.PubKey) (crypto.PubKey, error) {
	return m.Index().(*metadataStoreIndex).getMemberByDevice(pk)
}

func (m *MetadataStore) GetDevicesForMember(pk crypto.PubKey) ([]crypto.PubKey, error) {
	return m.Index().(*metadataStoreIndex).getDevicesForMember(pk)
}

func (m *MetadataStore) ListAdmins() []crypto.PubKey {
	if m.typeChecker(isContactGroup, isAccountGroup) {
		return m.ListMembers()
	}

	return m.Index().(*metadataStoreIndex).listAdmins()
}

func (m *MetadataStore) GetIncomingContactRequestsStatus() (bool, *protocoltypes.ShareableContact) {
	if !m.typeChecker(isAccountGroup) {
		return false, nil
	}

	enabled := m.Index().(*metadataStoreIndex).contactRequestsEnabled()
	seed := m.Index().(*metadataStoreIndex).contactRequestsSeed()

	rawMemberDevice, err := m.memberDevice.Member().Raw()
	if err != nil {
		m.logger.Error("unable to serialize member public key", zap.Error(err))
		return enabled, nil
	}

	contactRef := &protocoltypes.ShareableContact{
		Pk:                   rawMemberDevice,
		PublicRendezvousSeed: seed,
	}

	return enabled, contactRef
}

func (m *MetadataStore) ListMembers() []crypto.PubKey {
	if m.typeChecker(isAccountGroup, isContactGroup, isMultiMemberGroup) {
		return m.Index().(*metadataStoreIndex).listMembers()
	}

	return nil
}

func (m *MetadataStore) ListDevices() []crypto.PubKey {
	return m.Index().(*metadataStoreIndex).listDevices()
}

func (m *MetadataStore) ListMultiMemberGroups() []*protocoltypes.Group {
	if !m.typeChecker(isAccountGroup) {
		return nil
	}

	idx, ok := m.Index().(*metadataStoreIndex)
	if !ok {
		return nil
	}
	idx.lock.Lock()
	defer idx.lock.Unlock()

	groups := []*protocoltypes.Group(nil)

	for _, g := range idx.groups {
		if g.state != accountGroupJoinedStateJoined {
			continue
		}

		groups = append(groups, g.group)
	}

	return groups
}

func (m *MetadataStore) ListOtherMembersDevices() []crypto.PubKey {
	return m.Index().(*metadataStoreIndex).listOtherMembersDevices()
}

func (m *MetadataStore) GetRequestOwnMetadataForContact(pk []byte) ([]byte, error) {
	idx, ok := m.Index().(*metadataStoreIndex)
	if !ok {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("invalid index type"))
	}

	idx.lock.Lock()
	defer idx.lock.Unlock()

	meta, ok := idx.contactRequestMetadata[string(pk)]
	if !ok {
		return nil, errcode.ErrCode_ErrMissingMapKey.Wrap(fmt.Errorf("no metadata found for specified contact"))
	}

	return meta, nil
}

func (m *MetadataStore) ListContactsByStatus(states ...protocoltypes.ContactState) []*protocoltypes.ShareableContact {
	if !m.typeChecker(isAccountGroup) {
		return nil
	}

	idx, ok := m.Index().(*metadataStoreIndex)
	if !ok {
		return nil
	}
	idx.lock.Lock()
	defer idx.lock.Unlock()

	contacts := []*protocoltypes.ShareableContact(nil)

	for _, c := range idx.contacts {
		hasState := slices.Contains(states, c.state)

		if hasState {
			contacts = append(contacts, c.contact)
		}
	}

	return contacts
}

func (m *MetadataStore) GetContactFromGroupPK(groupPK []byte) *protocoltypes.ShareableContact {
	if !m.typeChecker(isAccountGroup) {
		return nil
	}

	idx, ok := m.Index().(*metadataStoreIndex)
	if !ok {
		return nil
	}
	idx.lock.Lock()
	defer idx.lock.Unlock()

	contact, ok := idx.contactsFromGroupPK[string(groupPK)]
	if !ok || contact == nil {
		return nil
	}

	return contact.contact
}

func (m *MetadataStore) checkIfInGroup(pk []byte) bool {
	idx, ok := m.Index().(*metadataStoreIndex)
	if !ok {
		return false
	}

	idx.lock.Lock()
	defer idx.lock.Unlock()

	if existingGroup, ok := idx.groups[string(pk)]; ok && existingGroup.state == accountGroupJoinedStateJoined {
		return true
	}

	return false
}

// GroupJoin indicates the payload includes that the deviceKeystore has joined a group
func (m *MetadataStore) GroupJoin(ctx context.Context, g *protocoltypes.Group) (operation.Operation, error) {
	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	if err := g.IsValid(); err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if m.checkIfInGroup(g.PublicKey) {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("already present in group"))
	}

	return m.attributeSignAndAddEvent(ctx, &protocoltypes.AccountGroupJoined{
		Group: g,
	}, protocoltypes.EventType_EventTypeAccountGroupJoined)
}

// GroupLeave indicates the payload includes that the deviceKeystore has left a group
func (m *MetadataStore) GroupLeave(ctx context.Context, pk crypto.PubKey) (operation.Operation, error) {
	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	if pk == nil {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	bytes, err := pk.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	if !m.checkIfInGroup(bytes) {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	return m.groupAction(ctx, pk, &protocoltypes.AccountGroupLeft{}, protocoltypes.EventType_EventTypeAccountGroupLeft)
}

// ContactRequestDisable indicates the payload includes that the deviceKeystore has disabled incoming contact requests
func (m *MetadataStore) ContactRequestDisable(ctx context.Context) (operation.Operation, error) {
	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	return m.attributeSignAndAddEvent(ctx, &protocoltypes.AccountContactRequestDisabled{}, protocoltypes.EventType_EventTypeAccountContactRequestDisabled)
}

// ContactRequestEnable indicates the payload includes that the deviceKeystore has enabled incoming contact requests
func (m *MetadataStore) ContactRequestEnable(ctx context.Context) (operation.Operation, error) {
	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	return m.attributeSignAndAddEvent(ctx, &protocoltypes.AccountContactRequestEnabled{}, protocoltypes.EventType_EventTypeAccountContactRequestEnabled)
}

// ContactRequestReferenceReset indicates the payload includes that the deviceKeystore has a new contact request reference
func (m *MetadataStore) ContactRequestReferenceReset(ctx context.Context) (operation.Operation, error) {
	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	seed, err := genNewSeed()
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoKeyGeneration.Wrap(err)
	}

	return m.attributeSignAndAddEvent(ctx, &protocoltypes.AccountContactRequestReferenceReset{
		PublicRendezvousSeed: seed,
	}, protocoltypes.EventType_EventTypeAccountContactRequestReferenceReset)
}

// ContactRequestOutgoingEnqueue indicates the payload includes that the deviceKeystore will attempt to send a new contact request
func (m *MetadataStore) ContactRequestOutgoingEnqueue(ctx context.Context, contact *protocoltypes.ShareableContact, ownMetadata []byte) (operation.Operation, error) {
	ctx, _ = tyber.ContextWithTraceID(ctx)

	b64GroupPK := base64.RawURLEncoding.EncodeToString(m.group.PublicKey)
	m.logger.Debug("Enqueuing contact request", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "GroupPK", Description: fmt.Sprint(b64GroupPK)}})...)

	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	if err := contact.CheckFormat(); err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	accountPublicKey := m.memberDevice.Member()
	if contact.IsSamePK(accountPublicKey) {
		return nil, errcode.ErrCode_ErrContactRequestSameAccount
	}

	pk, err := contact.GetPubKey()
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if m.checkContactStatus(pk, protocoltypes.ContactState_ContactStateAdded) {
		return nil, errcode.ErrCode_ErrContactRequestContactAlreadyAdded
	}

	if m.checkContactStatus(pk, protocoltypes.ContactState_ContactStateRemoved, protocoltypes.ContactState_ContactStateDiscarded, protocoltypes.ContactState_ContactStateReceived) {
		return m.ContactRequestOutgoingSent(ctx, pk)
	}

	op, err := m.attributeSignAndAddEvent(ctx, &protocoltypes.AccountContactRequestOutgoingEnqueued{
		Contact: &protocoltypes.ShareableContact{
			Pk:                   contact.Pk,
			PublicRendezvousSeed: contact.PublicRendezvousSeed,
			Metadata:             contact.Metadata,
		},
		OwnMetadata: ownMetadata,
	}, protocoltypes.EventType_EventTypeAccountContactRequestOutgoingEnqueued)

	m.logger.Debug("Enqueued contact request", tyber.FormatStepLogFields(ctx, []tyber.Detail{})...)

	return op, err
}

// ContactRequestOutgoingSent indicates the payload includes that the deviceKeystore has sent a contact request
func (m *MetadataStore) ContactRequestOutgoingSent(ctx context.Context, pk crypto.PubKey) (operation.Operation, error) {
	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	switch m.getContactStatus(pk) {
	case protocoltypes.ContactState_ContactStateToRequest:
	case protocoltypes.ContactState_ContactStateReceived:
	case protocoltypes.ContactState_ContactStateRemoved:
	case protocoltypes.ContactState_ContactStateDiscarded:

	case protocoltypes.ContactState_ContactStateUndefined:
		return nil, errcode.ErrCode_ErrContactRequestContactUndefined
	case protocoltypes.ContactState_ContactStateAdded:
		return nil, errcode.ErrCode_ErrContactRequestContactAlreadyAdded
	case protocoltypes.ContactState_ContactStateBlocked:
		return nil, errcode.ErrCode_ErrContactRequestContactBlocked
	default:
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	return m.contactAction(ctx, pk, &protocoltypes.AccountContactRequestOutgoingSent{}, protocoltypes.EventType_EventTypeAccountContactRequestOutgoingSent)
}

// ContactRequestIncomingReceived indicates the payload includes that the deviceKeystore has received a contact request
func (m *MetadataStore) ContactRequestIncomingReceived(ctx context.Context, contact *protocoltypes.ShareableContact) (operation.Operation, error) {
	m.logger.Debug("Sending ContactRequestIncomingReceived on Account group", tyber.FormatStepLogFields(ctx, []tyber.Detail{})...)

	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	if err := contact.CheckFormat(protocoltypes.ShareableContactOptionsAllowMissingRDVSeed); err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	accountPublicKey := m.memberDevice.Member()
	if contact.IsSamePK(accountPublicKey) {
		return nil, errcode.ErrCode_ErrContactRequestSameAccount
	}

	pk, err := contact.GetPubKey()
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	switch m.getContactStatus(pk) {
	case protocoltypes.ContactState_ContactStateUndefined:
	case protocoltypes.ContactState_ContactStateRemoved:
	case protocoltypes.ContactState_ContactStateDiscarded:

	// If incoming request comes from an account for which an outgoing request
	// is in "sending" state, mark the outgoing request as "sent"
	case protocoltypes.ContactState_ContactStateToRequest:
		return m.ContactRequestOutgoingSent(ctx, pk)

	// Errors
	case protocoltypes.ContactState_ContactStateReceived:
		return nil, errcode.ErrCode_ErrContactRequestIncomingAlreadyReceived
	case protocoltypes.ContactState_ContactStateAdded:
		return nil, errcode.ErrCode_ErrContactRequestContactAlreadyAdded
	case protocoltypes.ContactState_ContactStateBlocked:
		return nil, errcode.ErrCode_ErrContactRequestContactBlocked
	default:
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	return m.attributeSignAndAddEvent(ctx, &protocoltypes.AccountContactRequestIncomingReceived{
		ContactPk:             contact.Pk,
		ContactRendezvousSeed: contact.PublicRendezvousSeed,
		ContactMetadata:       contact.Metadata,
	}, protocoltypes.EventType_EventTypeAccountContactRequestIncomingReceived)
}

// ContactRequestIncomingDiscard indicates the payload includes that the deviceKeystore has ignored a contact request
func (m *MetadataStore) ContactRequestIncomingDiscard(ctx context.Context, pk crypto.PubKey) (operation.Operation, error) {
	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	if !m.checkContactStatus(pk, protocoltypes.ContactState_ContactStateReceived) {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	return m.contactAction(ctx, pk, &protocoltypes.AccountContactRequestIncomingDiscarded{}, protocoltypes.EventType_EventTypeAccountContactRequestIncomingDiscarded)
}

// ContactRequestIncomingAccept indicates the payload includes that the deviceKeystore has accepted a contact request
func (m *MetadataStore) ContactRequestIncomingAccept(ctx context.Context, pk crypto.PubKey) (operation.Operation, error) {
	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	if !m.checkContactStatus(pk, protocoltypes.ContactState_ContactStateReceived) {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	return m.contactAction(ctx, pk, &protocoltypes.AccountContactRequestIncomingAccepted{}, protocoltypes.EventType_EventTypeAccountContactRequestIncomingAccepted)
}

// ContactBlock indicates the payload includes that the deviceKeystore has blocked a contact
func (m *MetadataStore) ContactBlock(ctx context.Context, pk crypto.PubKey) (operation.Operation, error) {
	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	accountPublicKey := m.memberDevice.Member()
	if accountPublicKey.Equals(pk) {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	if m.checkContactStatus(pk, protocoltypes.ContactState_ContactStateBlocked) {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	return m.contactAction(ctx, pk, &protocoltypes.AccountContactBlocked{}, protocoltypes.EventType_EventTypeAccountContactBlocked)
}

// ContactUnblock indicates the payload includes that the deviceKeystore has unblocked a contact
func (m *MetadataStore) ContactUnblock(ctx context.Context, pk crypto.PubKey) (operation.Operation, error) {
	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	if !m.checkContactStatus(pk, protocoltypes.ContactState_ContactStateBlocked) {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	return m.contactAction(ctx, pk, &protocoltypes.AccountContactUnblocked{}, protocoltypes.EventType_EventTypeAccountContactUnblocked)
}

func (m *MetadataStore) ContactSendAliasKey(ctx context.Context) (operation.Operation, error) {
	if !m.typeChecker(isContactGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	accountProofPublicKey, err := m.secretStore.GetAccountProofPublicKey()
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	alias, err := accountProofPublicKey.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return m.attributeSignAndAddEvent(ctx, &protocoltypes.ContactAliasKeyAdded{
		AliasPk: alias,
	}, protocoltypes.EventType_EventTypeContactAliasKeyAdded)
}

func (m *MetadataStore) SendAliasProof(ctx context.Context) (operation.Operation, error) {
	if !m.typeChecker(isMultiMemberGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	resolver := []byte(nil) // TODO: should be a hmac value of something for quicker searches
	proof := []byte(nil)    // TODO: should be a signed value of something

	return m.attributeSignAndAddEvent(ctx, &protocoltypes.MultiMemberGroupAliasResolverAdded{
		AliasResolver: resolver,
		AliasProof:    proof,
	}, protocoltypes.EventType_EventTypeMultiMemberGroupAliasResolverAdded)
}

func (m *MetadataStore) SendAppMetadata(ctx context.Context, message []byte) (operation.Operation, error) {
	return m.attributeSignAndAddEvent(ctx, &protocoltypes.GroupMetadataPayloadSent{
		Message: message,
	}, protocoltypes.EventType_EventTypeGroupMetadataPayloadSent)
}

func (m *MetadataStore) SendAccountVerifiedCredentialAdded(ctx context.Context, token *protocoltypes.AccountVerifiedCredentialRegistered) (operation.Operation, error) {
	if !m.typeChecker(isAccountGroup) {
		return nil, errcode.ErrCode_ErrGroupInvalidType
	}

	return m.attributeSignAndAddEvent(ctx, token, protocoltypes.EventType_EventTypeAccountVerifiedCredentialRegistered)
}

func (m *MetadataStore) SendGroupReplicating(ctx context.Context, authenticationURL, replicationServer string) (operation.Operation, error) {
	return m.attributeSignAndAddEvent(ctx, &protocoltypes.GroupReplicating{
		AuthenticationUrl: authenticationURL,
		ReplicationServer: replicationServer,
	}, protocoltypes.EventType_EventTypeGroupReplicating)
}

type accountSignableEvent interface {
	proto.Message
	SetDevicePK([]byte)
}

type accountContactEvent interface {
	accountSignableEvent
	SetContactPK([]byte)
}

type accountGroupEvent interface {
	accountSignableEvent
	SetGroupPK([]byte)
}

func (m *MetadataStore) attributeSignAndAddEvent(ctx context.Context, evt accountSignableEvent, eventType protocoltypes.EventType) (operation.Operation, error) {
	evt.SetDevicePK(m.devicePublicKeyRaw)

	sig, err := signProtoWithDevice(evt, m.memberDevice)
	if err != nil {
		return nil, errcode.ErrCode_ErrCryptoSignature.Wrap(err)
	}

	m.logger.Debug("Signed event", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "Signature", Description: base64.RawURLEncoding.EncodeToString(sig)}})...)

	return metadataStoreAddEvent(ctx, m, m.group, eventType, evt, sig)
}

func (m *MetadataStore) contactAction(ctx context.Context, pk crypto.PubKey, event accountContactEvent, evtType protocoltypes.EventType) (operation.Operation, error) {
	ctx, newTrace := tyber.ContextWithTraceID(ctx)
	var tyberFields []zap.Field
	if newTrace {
		tyberFields = tyber.FormatTraceLogFields(ctx)
	} else {
		tyberFields = tyber.FormatStepLogFields(ctx, []tyber.Detail{})
	}
	m.logger.Debug("Sending "+strings.TrimPrefix(evtType.String(), "EventType")+" on Account group", tyberFields...)

	if pk == nil || event == nil {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	pkBytes, err := pk.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	event.SetContactPK(pkBytes)

	op, err := m.attributeSignAndAddEvent(ctx, event, evtType)
	if err != nil {
		return nil, err
	}

	if newTrace {
		m.logger.Debug("Event added successfully", tyber.FormatStepLogFields(ctx, []tyber.Detail{}, tyber.EndTrace)...)
	}
	return op, nil
}

func (m *MetadataStore) groupAction(ctx context.Context, pk crypto.PubKey, event accountGroupEvent, evtType protocoltypes.EventType) (operation.Operation, error) {
	pkBytes, err := pk.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	event.SetGroupPK(pkBytes)

	return m.attributeSignAndAddEvent(ctx, event, evtType)
}

func (m *MetadataStore) getContactStatus(pk crypto.PubKey) protocoltypes.ContactState {
	if pk == nil {
		return protocoltypes.ContactState_ContactStateUndefined
	}

	contact, err := m.Index().(*metadataStoreIndex).getContact(pk)
	if err != nil {
		m.logger.Warn("unable to get contact for public key", zap.Error(err))
		return protocoltypes.ContactState_ContactStateUndefined
	}

	return contact.state
}

func (m *MetadataStore) checkContactStatus(pk crypto.PubKey, states ...protocoltypes.ContactState) bool {
	contactStatus := m.getContactStatus(pk)

	return slices.Contains(states, contactStatus)
}

type EventMetadataReceived struct {
	MetaEvent *protocoltypes.GroupMetadataEvent
	Event     proto.Message
}

func constructorFactoryGroupMetadata(s *WeshOrbitDB, logger *zap.Logger) iface.StoreConstructor {
	return func(ipfs coreiface.CoreAPI, identity *identityprovider.Identity, addr address.Address, options *iface.NewStoreOptions) (iface.Store, error) {
		g, err := s.getGroupFromOptions(options)
		if err != nil {
			return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
		}
		shortGroupType := strings.TrimPrefix(g.GetGroupType().String(), "GroupType")
		b64GroupPK := base64.RawURLEncoding.EncodeToString(g.PublicKey)

		replication := false

		if options.EventBus == nil {
			options.EventBus = eventbus.NewBus()
		}

		store := &MetadataStore{
			eventBus:    options.EventBus,
			group:       g,
			logger:      logger,
			secretStore: s.secretStore,
		}

		if s.replicationMode {
			replication = true
		} else {
			var err error

			store.memberDevice, err = s.secretStore.GetOwnMemberDeviceForGroup(g)
			if err != nil {
				if errcode.Is(err, errcode.ErrCode_ErrInvalidInput) {
					replication = true
				} else {
					return nil, errcode.ErrCode_ErrOrbitDBInit.Wrap(err)
				}
			} else {
				store.devicePublicKeyRaw, err = store.memberDevice.Device().Raw()
				if err != nil {
					return nil, errcode.ErrCode_ErrOrbitDBInit.Wrap(err)
				}
			}
		}

		store.ctx, store.cancel = context.WithCancel(context.Background())

		if err := store.initEmitter(); err != nil {
			return nil, fmt.Errorf("unable to init emitters: %w", err)
		}

		if replication {
			options.Index = basestore.NewNoopIndex
			if err := store.InitBaseStore(ipfs, identity, addr, options); err != nil {
				store.cancel()
				return nil, errcode.ErrCode_ErrOrbitDBInit.Wrap(err)
			}

			return store, nil
		}

		chSub, err := store.eventBus.Subscribe([]interface{}{
			new(stores.EventWrite),
			new(stores.EventReplicated),
		}, eventbus.BufSize(128))
		if err != nil {
			store.cancel()
			return nil, fmt.Errorf("unable to subscribe to store events")
		}

		// Enable logs in the metadata index
		store.setLogger(logger)

		go func(ctx context.Context) {
			defer chSub.Close()

			for {
				var e interface{}
				select {
				case e = <-chSub.Out():
				case <-ctx.Done():
					return
				}

				var entries []ipfslog.Entry

				switch evt := e.(type) {
				case stores.EventWrite:
					entries = []ipfslog.Entry{evt.Entry}

				case stores.EventReplicated:
					entries = evt.Entries
				}

				for _, entry := range entries {
					ctx = tyber.ContextWithConstantTraceID(ctx, "msgrcvd-"+entry.GetHash().String())
					tyber.LogTraceStart(ctx, store.logger, fmt.Sprintf("Received metadata from %s group %s", shortGroupType, b64GroupPK))

					metaEvent, event, err := openMetadataEntry(store.OpLog(), entry, g)
					if err != nil {
						_ = tyber.LogFatalError(ctx, store.logger, "Unable to open metadata event", err, tyber.WithDetail("RawEvent", fmt.Sprint(e)), tyber.ForceReopen)
						continue
					}

					tyber.LogStep(ctx, store.logger, "Opened metadata store event",
						tyber.ForceReopen,
						tyber.EndTrace,
						tyber.WithJSONDetail("MetaEvent", metaEvent),
						tyber.WithJSONDetail("Event", event),
						tyber.UpdateTraceName(fmt.Sprintf("Received %s from %s group %s", strings.TrimPrefix(metaEvent.GetMetadata().GetEventType().String(), "EventType"), shortGroupType, b64GroupPK)),
					)

					recvEvent := EventMetadataReceived{
						MetaEvent: metaEvent,
						Event:     event,
					}

					if err := store.emitters.metadataReceived.Emit(recvEvent); err != nil {
						store.logger.Warn("unable to emit recv event", zap.Error(err))
					}

					if err := store.emitters.groupMetadata.Emit(metaEvent); err != nil {
						store.logger.Warn("unable to emit group metadata event", zap.Error(err))
					}
				}
			}
		}(store.ctx)

		options.Index = newMetadataIndex(store.ctx, g, store.memberDevice, s.secretStore)
		if err := store.InitBaseStore(ipfs, identity, addr, options); err != nil {
			store.cancel()
			return nil, errcode.ErrCode_ErrOrbitDBInit.Wrap(err)
		}

		return store, nil
	}
}

func (m *MetadataStore) initEmitter() (err error) {
	if m.emitters.metadataReceived, err = m.eventBus.Emitter(new(EventMetadataReceived)); err != nil {
		return err
	}

	if m.emitters.groupMetadata, err = m.eventBus.Emitter(new(*protocoltypes.GroupMetadataEvent)); err != nil {
		return err
	}

	return err
}

func genNewSeed() (seed []byte, err error) {
	seed, err = io.ReadAll(io.LimitReader(crand.Reader, protocoltypes.RendezvousSeedLength))
	return seed, err
}

func (m *MetadataStore) Close() error {
	m.cancel()
	return m.BaseStore.Close()
}
