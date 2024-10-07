package weshnet

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/crypto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	ipfslog "berty.tech/go-ipfs-log"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/weshnet/pkg/cryptoutil"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/secretstore"
)

// FIXME: replace members, devices, sentSecrets, contacts and groups by a circular buffer to avoid an attack by RAM saturation
type metadataStoreIndex struct {
	members                  map[string][]secretstore.MemberDevice
	devices                  map[string]secretstore.MemberDevice
	handledEvents            map[string]struct{}
	sentSecrets              map[string]struct{}
	admins                   map[crypto.PubKey]struct{}
	contacts                 map[string]*AccountContact
	contactsFromGroupPK      map[string]*AccountContact
	groups                   map[string]*accountGroup
	contactRequestMetadata   map[string][]byte
	verifiedCredentials      []*protocoltypes.AccountVerifiedCredentialRegistered
	contactRequestSeed       []byte
	contactRequestEnabled    *bool
	eventHandlers            map[protocoltypes.EventType][]func(event proto.Message) error
	postIndexActions         []func() error
	eventsContactAddAliasKey []*protocoltypes.ContactAliasKeyAdded
	ownAliasKeySent          bool
	otherAliasKey            []byte
	group                    *protocoltypes.Group
	ownMemberDevice          secretstore.MemberDevice
	secretStore              secretstore.SecretStore
	ctx                      context.Context
	lock                     sync.RWMutex
	logger                   *zap.Logger
}

//nolint:revive
func (m *metadataStoreIndex) Get(key string) interface{} {
	return nil
}

func (m *metadataStoreIndex) setLogger(logger *zap.Logger) {
	if logger == nil {
		return
	}

	m.logger = logger
}

func (m *metadataStoreIndex) UpdateIndex(log ipfslog.Log, _ []ipfslog.Entry) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	entries := log.GetEntries().Slice()

	// Resetting state
	m.contacts = map[string]*AccountContact{}
	m.contactsFromGroupPK = map[string]*AccountContact{}
	m.groups = map[string]*accountGroup{}
	m.contactRequestMetadata = map[string][]byte{}
	m.contactRequestEnabled = nil
	m.contactRequestSeed = []byte(nil)
	m.verifiedCredentials = nil
	m.handledEvents = map[string]struct{}{}

	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]

		_, alreadyHandledEvent := m.handledEvents[e.GetHash().String()]

		// TODO: improve account events handling
		if m.group.GroupType != protocoltypes.GroupType_GroupTypeAccount && alreadyHandledEvent {
			continue
		}

		metaEvent, event, err := openMetadataEntry(log, e, m.group)
		if err != nil {
			m.logger.Error("unable to open metadata entry", zap.Error(err))
			continue
		}

		handlers, ok := m.eventHandlers[metaEvent.Metadata.EventType]
		if !ok {
			m.handledEvents[e.GetHash().String()] = struct{}{}
			m.logger.Error("handler for event type not found", zap.String("event-type", metaEvent.Metadata.EventType.String()))
			continue
		}

		var lastErr error

		for _, h := range handlers {
			err = h(event)
			if err != nil {
				m.logger.Error("unable to handle event", zap.Error(err))
				lastErr = err
			}
		}

		if lastErr != nil {
			m.handledEvents[e.GetHash().String()] = struct{}{}
			continue
		}

		m.handledEvents[e.GetHash().String()] = struct{}{}
	}

	for _, h := range m.postIndexActions {
		if err := h(); err != nil {
			return errcode.ErrCode_ErrInternal.Wrap(err)
		}
	}

	return nil
}

func (m *metadataStoreIndex) handleGroupMemberDeviceAdded(event proto.Message) error {
	e, ok := event.(*protocoltypes.GroupMemberDeviceAdded)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	member, err := crypto.UnmarshalEd25519PublicKey(e.MemberPk)
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	device, err := crypto.UnmarshalEd25519PublicKey(e.DevicePk)
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if _, ok := m.devices[string(e.DevicePk)]; ok {
		return nil
	}

	memberDevice := secretstore.NewMemberDevice(member, device)

	m.devices[string(e.DevicePk)] = memberDevice
	m.members[string(e.MemberPk)] = append(m.members[string(e.MemberPk)], memberDevice)

	return nil
}

func (m *metadataStoreIndex) handleGroupDeviceChainKeyAdded(event proto.Message) error {
	e, ok := event.(*protocoltypes.GroupDeviceChainKeyAdded)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	_, err := crypto.UnmarshalEd25519PublicKey(e.DestMemberPk)
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	senderPK, err := crypto.UnmarshalEd25519PublicKey(e.DevicePk)
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if m.ownMemberDevice.Device().Equals(senderPK) {
		m.sentSecrets[string(e.DestMemberPk)] = struct{}{}
	}

	return nil
}

func (m *metadataStoreIndex) getMemberByDevice(devicePublicKey crypto.PubKey) (crypto.PubKey, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	publicKeyBytes, err := devicePublicKey.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	return m.unsafeGetMemberByDevice(publicKeyBytes)
}

func (m *metadataStoreIndex) unsafeGetMemberByDevice(publicKeyBytes []byte) (crypto.PubKey, error) {
	if l := len(publicKeyBytes); l != cryptoutil.KeySize {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("invalid private key size, expected %d got %d", cryptoutil.KeySize, l))
	}

	device, ok := m.devices[string(publicKeyBytes)]
	if !ok {
		return nil, errcode.ErrCode_ErrMissingInput
	}

	return device.Member(), nil
}

func (m *metadataStoreIndex) getDevicesForMember(pk crypto.PubKey) ([]crypto.PubKey, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	id, err := pk.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	mds, ok := m.members[string(id)]
	if !ok {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	ret := make([]crypto.PubKey, len(mds))
	for i, md := range mds {
		ret[i] = md.Device()
	}

	return ret, nil
}

func (m *metadataStoreIndex) MemberCount() int {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return len(m.members)
}

func (m *metadataStoreIndex) DeviceCount() int {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return len(m.devices)
}

func (m *metadataStoreIndex) listContacts() map[string]*AccountContact {
	m.lock.RLock()
	defer m.lock.RUnlock()

	contacts := make(map[string]*AccountContact)

	for k, contact := range m.contacts {
		contacts[k] = &AccountContact{
			state: contact.state,
			contact: &protocoltypes.ShareableContact{
				Pk:                   contact.contact.Pk,
				PublicRendezvousSeed: contact.contact.PublicRendezvousSeed,
				Metadata:             contact.contact.Metadata,
			},
		}
	}

	return contacts
}

func (m *metadataStoreIndex) listVerifiedCredentials() []*protocoltypes.AccountVerifiedCredentialRegistered {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.verifiedCredentials
}

func (m *metadataStoreIndex) listMembers() []crypto.PubKey {
	m.lock.RLock()
	defer m.lock.RUnlock()

	members := make([]crypto.PubKey, len(m.members))
	i := 0

	for _, md := range m.members {
		members[i] = md[0].Member()
		i++
	}

	return members
}

func (m *metadataStoreIndex) listDevices() []crypto.PubKey {
	m.lock.RLock()
	defer m.lock.RUnlock()

	devices := make([]crypto.PubKey, len(m.devices))
	i := 0

	for _, md := range m.devices {
		devices[i] = md.Device()
		i++
	}

	return devices
}

func (m *metadataStoreIndex) areSecretsAlreadySent(pk crypto.PubKey) (bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	key, err := pk.Raw()
	if err != nil {
		return false, errcode.ErrCode_ErrInvalidInput.Wrap(err)
	}

	_, ok := m.sentSecrets[string(key)]
	return ok, nil
}

type accountGroupJoinedState uint32

const (
	accountGroupJoinedStateJoined accountGroupJoinedState = iota + 1
	accountGroupJoinedStateLeft
)

type accountGroup struct {
	state accountGroupJoinedState
	group *protocoltypes.Group
}

type AccountContact struct {
	state   protocoltypes.ContactState
	contact *protocoltypes.ShareableContact
}

func (m *metadataStoreIndex) handleGroupJoined(event proto.Message) error {
	evt, ok := event.(*protocoltypes.AccountGroupJoined)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	_, ok = m.groups[string(evt.Group.PublicKey)]
	if ok {
		return nil
	}

	m.groups[string(evt.Group.PublicKey)] = &accountGroup{
		group: evt.Group,
		state: accountGroupJoinedStateJoined,
	}

	return nil
}

func (m *metadataStoreIndex) handleGroupLeft(event proto.Message) error {
	evt, ok := event.(*protocoltypes.AccountGroupLeft)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	_, ok = m.groups[string(evt.GroupPk)]
	if ok {
		return nil
	}

	m.groups[string(evt.GroupPk)] = &accountGroup{
		state: accountGroupJoinedStateLeft,
	}

	return nil
}

func (m *metadataStoreIndex) handleContactRequestDisabled(event proto.Message) error {
	if m.contactRequestEnabled != nil {
		return nil
	}

	_, ok := event.(*protocoltypes.AccountContactRequestDisabled)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	f := false
	m.contactRequestEnabled = &f

	return nil
}

func (m *metadataStoreIndex) handleContactRequestEnabled(event proto.Message) error {
	if m.contactRequestEnabled != nil {
		return nil
	}

	_, ok := event.(*protocoltypes.AccountContactRequestEnabled)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	t := true
	m.contactRequestEnabled = &t

	return nil
}

func (m *metadataStoreIndex) handleContactRequestReferenceReset(event proto.Message) error {
	evt, ok := event.(*protocoltypes.AccountContactRequestReferenceReset)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	if m.contactRequestSeed != nil {
		return nil
	}

	m.contactRequestSeed = evt.PublicRendezvousSeed

	return nil
}

func (m *metadataStoreIndex) registerContactFromGroupPK(ac *AccountContact) error {
	if m.group.GroupType != protocoltypes.GroupType_GroupTypeAccount {
		return errcode.ErrCode_ErrGroupInvalidType
	}

	contactPK, err := crypto.UnmarshalEd25519PublicKey(ac.contact.Pk)
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	group, err := m.secretStore.GetGroupForContact(contactPK)
	if err != nil {
		return errcode.ErrCode_ErrOrbitDBOpen.Wrap(err)
	}

	m.contactsFromGroupPK[string(group.PublicKey)] = ac

	return nil
}

func (m *metadataStoreIndex) handleContactRequestOutgoingEnqueued(event proto.Message) error {
	evt, ok := event.(*protocoltypes.AccountContactRequestOutgoingEnqueued)
	if ko := !ok || evt.Contact == nil; ko {
		return errcode.ErrCode_ErrInvalidInput
	}

	if _, ok := m.contacts[string(evt.Contact.Pk)]; ok {
		if m.contacts[string(evt.Contact.Pk)].contact.Metadata == nil {
			m.contacts[string(evt.Contact.Pk)].contact.Metadata = evt.Contact.Metadata
		}

		if m.contacts[string(evt.Contact.Pk)].contact.PublicRendezvousSeed == nil {
			m.contacts[string(evt.Contact.Pk)].contact.PublicRendezvousSeed = evt.Contact.PublicRendezvousSeed
		}

		return nil
	}

	if data, ok := m.contactRequestMetadata[string(evt.Contact.Pk)]; !ok || len(data) == 0 {
		m.contactRequestMetadata[string(evt.Contact.Pk)] = evt.OwnMetadata
	}

	ac := &AccountContact{
		state: protocoltypes.ContactState_ContactStateToRequest,
		contact: &protocoltypes.ShareableContact{
			Pk:                   evt.Contact.Pk,
			Metadata:             evt.Contact.Metadata,
			PublicRendezvousSeed: evt.Contact.PublicRendezvousSeed,
		},
	}

	m.contacts[string(evt.Contact.Pk)] = ac
	err := m.registerContactFromGroupPK(ac)

	return err
}

func (m *metadataStoreIndex) handleContactRequestOutgoingSent(event proto.Message) error {
	evt, ok := event.(*protocoltypes.AccountContactRequestOutgoingSent)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	if _, ok := m.contacts[string(evt.ContactPk)]; ok {
		return nil
	}

	ac := &AccountContact{
		state: protocoltypes.ContactState_ContactStateAdded,
		contact: &protocoltypes.ShareableContact{
			Pk: evt.ContactPk,
		},
	}

	m.contacts[string(evt.ContactPk)] = ac
	err := m.registerContactFromGroupPK(ac)

	return err
}

func (m *metadataStoreIndex) handleContactRequestIncomingReceived(event proto.Message) error {
	evt, ok := event.(*protocoltypes.AccountContactRequestIncomingReceived)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	if _, ok := m.contacts[string(evt.ContactPk)]; ok {
		if m.contacts[string(evt.ContactPk)].contact.Metadata == nil {
			m.contacts[string(evt.ContactPk)].contact.Metadata = evt.ContactMetadata
		}

		if m.contacts[string(evt.ContactPk)].contact.PublicRendezvousSeed == nil {
			m.contacts[string(evt.ContactPk)].contact.PublicRendezvousSeed = evt.ContactRendezvousSeed
		}

		return nil
	}

	ac := &AccountContact{
		state: protocoltypes.ContactState_ContactStateReceived,
		contact: &protocoltypes.ShareableContact{
			Pk:                   evt.ContactPk,
			Metadata:             evt.ContactMetadata,
			PublicRendezvousSeed: evt.ContactRendezvousSeed,
		},
	}

	m.contacts[string(evt.ContactPk)] = ac
	err := m.registerContactFromGroupPK(ac)

	return err
}

func (m *metadataStoreIndex) handleContactRequestIncomingDiscarded(event proto.Message) error {
	evt, ok := event.(*protocoltypes.AccountContactRequestIncomingDiscarded)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	if _, ok := m.contacts[string(evt.ContactPk)]; ok {
		return nil
	}

	ac := &AccountContact{
		state: protocoltypes.ContactState_ContactStateDiscarded,
		contact: &protocoltypes.ShareableContact{
			Pk: evt.ContactPk,
		},
	}

	m.contacts[string(evt.ContactPk)] = ac
	err := m.registerContactFromGroupPK(ac)

	return err
}

func (m *metadataStoreIndex) handleContactRequestIncomingAccepted(event proto.Message) error {
	evt, ok := event.(*protocoltypes.AccountContactRequestIncomingAccepted)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	if _, ok := m.contacts[string(evt.ContactPk)]; ok {
		return nil
	}

	ac := &AccountContact{
		state: protocoltypes.ContactState_ContactStateAdded,
		contact: &protocoltypes.ShareableContact{
			Pk: evt.ContactPk,
		},
	}

	m.contacts[string(evt.ContactPk)] = ac
	err := m.registerContactFromGroupPK(ac)

	return err
}

func (m *metadataStoreIndex) handleContactBlocked(event proto.Message) error {
	evt, ok := event.(*protocoltypes.AccountContactBlocked)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	if _, ok := m.contacts[string(evt.ContactPk)]; ok {
		return nil
	}

	ac := &AccountContact{
		state: protocoltypes.ContactState_ContactStateBlocked,
		contact: &protocoltypes.ShareableContact{
			Pk: evt.ContactPk,
		},
	}

	m.contacts[string(evt.ContactPk)] = ac
	err := m.registerContactFromGroupPK(ac)

	return err
}

func (m *metadataStoreIndex) handleContactUnblocked(event proto.Message) error {
	evt, ok := event.(*protocoltypes.AccountContactUnblocked)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	if _, ok := m.contacts[string(evt.ContactPk)]; ok {
		return nil
	}

	ac := &AccountContact{
		state: protocoltypes.ContactState_ContactStateRemoved,
		contact: &protocoltypes.ShareableContact{
			Pk: evt.ContactPk,
		},
	}

	m.contacts[string(evt.ContactPk)] = ac
	err := m.registerContactFromGroupPK(ac)

	return err
}

func (m *metadataStoreIndex) handleContactAliasKeyAdded(event proto.Message) error {
	evt, ok := event.(*protocoltypes.ContactAliasKeyAdded)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	m.eventsContactAddAliasKey = append(m.eventsContactAddAliasKey, evt)

	return nil
}

func (m *metadataStoreIndex) handleMultiMemberInitialMember(event proto.Message) error {
	e, ok := event.(*protocoltypes.MultiMemberGroupInitialMemberAnnounced)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	pk, err := crypto.UnmarshalEd25519PublicKey(e.MemberPk)
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if _, ok := m.admins[pk]; ok {
		return errcode.ErrCode_ErrInternal
	}

	m.admins[pk] = struct{}{}

	return nil
}

//nolint:revive
func (m *metadataStoreIndex) handleMultiMemberGrantAdminRole(event proto.Message) error {
	// TODO:

	return nil
}

func (m *metadataStoreIndex) handleGroupMetadataPayloadSent(_ proto.Message) error {
	return nil
}

func (m *metadataStoreIndex) handleAccountVerifiedCredentialRegistered(event proto.Message) error {
	e, ok := event.(*protocoltypes.AccountVerifiedCredentialRegistered)
	if !ok {
		return errcode.ErrCode_ErrInvalidInput
	}

	m.verifiedCredentials = append(m.verifiedCredentials, e)

	return nil
}

func (m *metadataStoreIndex) listAdmins() []crypto.PubKey {
	m.lock.RLock()
	defer m.lock.RUnlock()

	admins := make([]crypto.PubKey, len(m.admins))
	i := 0

	for admin := range m.admins {
		admins[i] = admin
		i++
	}

	return admins
}

func (m *metadataStoreIndex) listOtherMembersDevices() []crypto.PubKey {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if m.ownMemberDevice == nil || m.ownMemberDevice.Member() == nil {
		return nil
	}

	ownMemberPK, err := m.ownMemberDevice.Member().Raw()
	if err != nil {
		m.logger.Warn("unable to serialize member pubkey", zap.Error(err))
		return nil
	}

	devices := []crypto.PubKey(nil)
	for pk, devicesForMember := range m.members {
		if string(ownMemberPK) == pk {
			continue
		}

		for _, md := range devicesForMember {
			devices = append(devices, md.Device())
		}
	}

	return devices
}

func (m *metadataStoreIndex) contactRequestsEnabled() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.contactRequestEnabled != nil && *m.contactRequestEnabled
}

func (m *metadataStoreIndex) contactRequestsSeed() []byte {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.contactRequestSeed
}

func (m *metadataStoreIndex) getContact(pk crypto.PubKey) (*AccountContact, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	bytes, err := pk.Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	contact, ok := m.contacts[string(bytes)]
	if !ok {
		return nil, errcode.ErrCode_ErrMissingMapKey.Wrap(err)
	}

	return contact, nil
}

func (m *metadataStoreIndex) postHandlerSentAliases() error {
	for _, evt := range m.eventsContactAddAliasKey {
		memberPublicKey, err := m.unsafeGetMemberByDevice(evt.DevicePk)
		if err != nil {
			return fmt.Errorf("couldn't get member for device")
		}

		if memberPublicKey.Equals(m.ownMemberDevice.Member()) {
			m.ownAliasKeySent = true
			continue
		}

		if l := len(evt.AliasPk); l != cryptoutil.KeySize {
			return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("invalid alias key size, expected %d, got %d", cryptoutil.KeySize, l))
		}

		m.otherAliasKey = evt.AliasPk
	}

	m.eventsContactAddAliasKey = nil

	return nil
}

// nolint:staticcheck,revive
// newMetadataIndex returns a new index to manage the list of the group members
func newMetadataIndex(ctx context.Context, g *protocoltypes.Group, md secretstore.MemberDevice, secretStore secretstore.SecretStore) iface.IndexConstructor {
	return func(publicKey []byte) iface.StoreIndex {
		m := &metadataStoreIndex{
			members:                map[string][]secretstore.MemberDevice{},
			devices:                map[string]secretstore.MemberDevice{},
			admins:                 map[crypto.PubKey]struct{}{},
			sentSecrets:            map[string]struct{}{},
			handledEvents:          map[string]struct{}{},
			contacts:               map[string]*AccountContact{},
			contactsFromGroupPK:    map[string]*AccountContact{},
			groups:                 map[string]*accountGroup{},
			contactRequestMetadata: map[string][]byte{},
			group:                  g,
			ownMemberDevice:        md,
			secretStore:            secretStore,
			ctx:                    ctx,
			logger:                 zap.NewNop(),
		}

		m.eventHandlers = map[protocoltypes.EventType][]func(event proto.Message) error{
			protocoltypes.EventType_EventTypeAccountContactBlocked:                  {m.handleContactBlocked},
			protocoltypes.EventType_EventTypeAccountContactRequestDisabled:          {m.handleContactRequestDisabled},
			protocoltypes.EventType_EventTypeAccountContactRequestEnabled:           {m.handleContactRequestEnabled},
			protocoltypes.EventType_EventTypeAccountContactRequestIncomingAccepted:  {m.handleContactRequestIncomingAccepted},
			protocoltypes.EventType_EventTypeAccountContactRequestIncomingDiscarded: {m.handleContactRequestIncomingDiscarded},
			protocoltypes.EventType_EventTypeAccountContactRequestIncomingReceived:  {m.handleContactRequestIncomingReceived},
			protocoltypes.EventType_EventTypeAccountContactRequestOutgoingEnqueued:  {m.handleContactRequestOutgoingEnqueued},
			protocoltypes.EventType_EventTypeAccountContactRequestOutgoingSent:      {m.handleContactRequestOutgoingSent},
			protocoltypes.EventType_EventTypeAccountContactRequestReferenceReset:    {m.handleContactRequestReferenceReset},
			protocoltypes.EventType_EventTypeAccountContactUnblocked:                {m.handleContactUnblocked},
			protocoltypes.EventType_EventTypeAccountGroupJoined:                     {m.handleGroupJoined},
			protocoltypes.EventType_EventTypeAccountGroupLeft:                       {m.handleGroupLeft},
			protocoltypes.EventType_EventTypeContactAliasKeyAdded:                   {m.handleContactAliasKeyAdded},
			protocoltypes.EventType_EventTypeGroupDeviceChainKeyAdded:               {m.handleGroupDeviceChainKeyAdded},
			protocoltypes.EventType_EventTypeGroupMemberDeviceAdded:                 {m.handleGroupMemberDeviceAdded},
			protocoltypes.EventType_EventTypeMultiMemberGroupAdminRoleGranted:       {m.handleMultiMemberGrantAdminRole},
			protocoltypes.EventType_EventTypeMultiMemberGroupInitialMemberAnnounced: {m.handleMultiMemberInitialMember},
			protocoltypes.EventType_EventTypeGroupMetadataPayloadSent:               {m.handleGroupMetadataPayloadSent},
			protocoltypes.EventType_EventTypeAccountVerifiedCredentialRegistered:    {m.handleAccountVerifiedCredentialRegistered},
		}

		m.postIndexActions = []func() error{
			m.postHandlerSentAliases,
		}

		return m
	}
}

var _ iface.StoreIndex = &metadataStoreIndex{}
