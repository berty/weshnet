package weshnet

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"go.uber.org/zap"

	"berty.tech/go-orbit-db/stores"
	"berty.tech/weshnet/pkg/cryptoutil"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/ipfsutil"
	"berty.tech/weshnet/pkg/logutil"
	"berty.tech/weshnet/pkg/protocoltypes"
)

type GroupContext struct {
	ctx             context.Context
	cancel          context.CancelFunc
	group           *protocoltypes.Group
	metadataStore   *MetadataStore
	messageStore    *MessageStore
	messageKeystore *cryptoutil.MessageKeystore
	memberDevice    *cryptoutil.OwnMemberDevice
	logger          *zap.Logger
	closed          uint32
	tasks           sync.WaitGroup

	devicesAdded   map[string]chan struct{}
	muDevicesAdded sync.RWMutex
}

func (gc *GroupContext) MessageKeystore() *cryptoutil.MessageKeystore {
	return gc.messageKeystore
}

func (gc *GroupContext) getMemberPrivKey() crypto.PrivKey {
	return gc.memberDevice.PrivateMember()
}

func (gc *GroupContext) MessageStore() *MessageStore {
	return gc.messageStore
}

func (gc *GroupContext) MetadataStore() *MetadataStore {
	return gc.metadataStore
}

func (gc *GroupContext) Group() *protocoltypes.Group {
	return gc.group
}

func (gc *GroupContext) MemberPubKey() crypto.PubKey {
	return gc.memberDevice.PrivateMember().GetPublic()
}

func (gc *GroupContext) DevicePubKey() crypto.PubKey {
	return gc.memberDevice.PrivateDevice().GetPublic()
}

func (gc *GroupContext) Close() error {
	gc.cancel()

	// @NOTE(gfanton): wait for active tasks to end, doing this we avoid to do
	// some operations on a closed store
	gc.tasks.Wait()

	// mark group context has closed
	atomic.StoreUint32(&gc.closed, 1)

	// @FIXME(gfanton): should we really handle store closing here ?
	gc.metadataStore.Close()
	gc.messageStore.Close()

	gc.logger.Debug("group context closed", zap.String("groupID", gc.group.GroupIDAsString()))
	return nil
}

func (gc *GroupContext) IsClosed() bool {
	return atomic.LoadUint32(&gc.closed) != 0
}

func NewContextGroup(group *protocoltypes.Group, metadataStore *MetadataStore, messageStore *MessageStore, messageKeystore *cryptoutil.MessageKeystore, memberDevice *cryptoutil.OwnMemberDevice, logger *zap.Logger) *GroupContext {
	ctx, cancel := context.WithCancel(context.Background())

	if logger == nil {
		logger = zap.NewNop()
	}

	return &GroupContext{
		ctx:             ctx,
		cancel:          cancel,
		group:           group,
		metadataStore:   metadataStore,
		messageStore:    messageStore,
		messageKeystore: messageKeystore,
		memberDevice:    memberDevice,
		logger:          logger.With(logutil.PrivateString("group-id", fmt.Sprintf("%.6s", base64.StdEncoding.EncodeToString(group.PublicKey)))),
		closed:          0,
		devicesAdded:    make(map[string]chan struct{}),
	}
}

func (gc *GroupContext) ActivateGroupContext(contactPK crypto.PubKey) (err error) {
	ctx := gc.ctx

	// start watching for GroupMetadataEvent to send secret and register
	// chainkey of new members.
	// notify `wgSelfAnnouncement` when we found ourself
	var wgSelfAnnouncement sync.WaitGroup
	{
		m := gc.MetadataStore()
		sub, err := m.EventBus().Subscribe(new(protocoltypes.GroupMetadataEvent))
		if err != nil {
			return fmt.Errorf("unable to subscribe to group metadata event: %w", err)
		}

		gc.tasks.Add(1)
		wgSelfAnnouncement.Add(1)
		go func() {
			var once sync.Once

			defer gc.tasks.Done()                  // ultimately, mark bg task has done
			defer once.Do(wgSelfAnnouncement.Done) // avoid deadlock
			defer sub.Close()

			for {
				var evt interface{}
				select {
				case <-ctx.Done():
					return
				case evt = <-sub.Out():
				}

				// @TODO(gfanton): should we handle this in a sub gorouting ?
				e := evt.(protocoltypes.GroupMetadataEvent)
				self, err := gc.handleGroupMetadataEvent(&e)
				if err != nil {
					gc.logger.Error("unable to handle EventTypeGroupDeviceSecretAdded", zap.Error(err))
				}

				if self {
					once.Do(wgSelfAnnouncement.Done)
				}
			}
		}()
	}

	// send secret and register key from existing memebers.
	// we should wait until all the events have been retreived.
	{
		var wgExistingMembers sync.WaitGroup

		wgExistingMembers.Add(2)
		go func() {
			start := time.Now()
			gc.fillMessageKeysHolderUsingPreviousData()
			wgExistingMembers.Done()
			gc.logger.Info(fmt.Sprintf("FillMessageKeysHolderUsingPreviousData took %s", time.Since(start)))
		}()

		go func() {
			start := time.Now()
			gc.sendSecretsToExistingMembers(contactPK)
			wgExistingMembers.Done()
			gc.logger.Info(fmt.Sprintf("SendSecretsToExistingMembers took %s", time.Since(start)))
		}()
	}

	// if selfAnnouncement is enable wait until we get our own secret
	start := time.Now()
	op, err := gc.MetadataStore().AddDeviceToGroup(gc.ctx)
	if err != nil {
		return fmt.Errorf("unable to add device to groupo: %w", err)
	}

	gc.logger.Info(fmt.Sprintf("AddDeviceToGroup took %s", time.Since(start)))
	if op != nil {
		// Waiting for async events to be handled
		wgSelfAnnouncement.Wait()
	}

	return nil
}

func (gc *GroupContext) handleGroupMetadataEvent(e *protocoltypes.GroupMetadataEvent) (isSelf bool, err error) {
	switch e.Metadata.EventType {
	case protocoltypes.EventTypeGroupMemberDeviceAdded:
		event := &protocoltypes.GroupAddMemberDevice{}
		if err := event.Unmarshal(e.Event); err != nil {
			gc.logger.Error("unable to unmarshal payload", zap.Error(err))
		}

		memberPK, err := crypto.UnmarshalEd25519PublicKey(event.MemberPK)
		if err != nil {
			return isSelf, fmt.Errorf("unable to unmarshal sender member pk: %w", err)
		}

		if memberPK.Equals(gc.memberDevice.PrivateMember().GetPublic()) {
			isSelf = true
		}

		if _, err := gc.MetadataStore().SendSecret(gc.ctx, memberPK); err != nil {
			if !errcode.Is(err, errcode.ErrGroupSecretAlreadySentToMember) {
				return isSelf, fmt.Errorf("unable to send secret to member: %w", err)
			}
		}

	case protocoltypes.EventTypeGroupDeviceSecretAdded:
		gc.muDevicesAdded.RLock()
		defer gc.muDevicesAdded.RUnlock()

		pk, ds, err := openDeviceSecret(e.Metadata, gc.getMemberPrivKey(), gc.Group())
		switch err {
		case nil: // ok
		case errcode.ErrInvalidInput, errcode.ErrGroupSecretOtherDestMember:
			// @FIXME(gfanton): should we log this ?
			return isSelf, nil
		default:
			return isSelf, fmt.Errorf("an error occurred while opening device secrets: %w", err)
		}

		if err = gc.MessageKeystore().RegisterChainKey(gc.ctx, gc.Group(), pk, ds, gc.DevicePubKey().Equals(pk)); err != nil {
			return isSelf, fmt.Errorf("unable to register chain key: %w", err)
		}

		if rawPK, err := pk.Raw(); err == nil {
			// A new chainKey has been registered
			// notify watcher, start in subroutine to avoid deadlock
			go gc.notifyDeviceAdded(rawPK)
			// process queued message and check if cached messages can be opened with it
			gc.MessageStore().ProcessMessageQueueForDevicePK(gc.ctx, rawPK)
		}
	}

	return isSelf, nil
}

func (gc *GroupContext) fillMessageKeysHolderUsingPreviousData() {
	publishedSecrets := gc.metadataStoreListSecrets()

	for pk, sec := range publishedSecrets {
		if err := gc.MessageKeystore().RegisterChainKey(gc.ctx, gc.Group(), pk, sec, gc.DevicePubKey().Equals(pk)); err != nil {
			gc.logger.Error("unable to register chain key", zap.Error(err))
			continue
		}
		// A new chainKey is registered, check if cached messages can be opened with it
		if rawPK, err := pk.Raw(); err == nil {
			gc.MessageStore().ProcessMessageQueueForDevicePK(gc.ctx, rawPK)
		}
	}
}

func (gc *GroupContext) metadataStoreListSecrets() map[crypto.PubKey]*protocoltypes.DeviceSecret {
	publishedSecrets := map[crypto.PubKey]*protocoltypes.DeviceSecret{}

	m := gc.MetadataStore()
	ownSK := gc.getMemberPrivKey()
	g := gc.Group()

	metadatas, err := m.ListEvents(gc.ctx, nil, nil, false)
	if err != nil {
		return nil
	}
	for metadata := range metadatas {
		if metadata == nil {
			continue
		}

		pk, ds, err := openDeviceSecret(metadata.Metadata, ownSK, g)
		if errcode.Is(err, errcode.ErrInvalidInput) || errcode.Is(err, errcode.ErrGroupSecretOtherDestMember) {
			continue
		}

		if err != nil {
			gc.logger.Error("unable to open device secret", zap.Error(err))
			continue
		}

		publishedSecrets[pk] = ds
	}

	return publishedSecrets
}

func (gc *GroupContext) sendSecretsToExistingMembers(contact crypto.PubKey) {
	members := gc.MetadataStore().ListMembers()

	// Force sending secret to contact member in contact group
	if gc.group.GroupType == protocoltypes.GroupTypeContact && len(members) < 2 && contact != nil {
		// Check if contact member is already listed
		found := false
		for _, member := range members {
			if member.Equals(contact) {
				found = true
			}
		}

		// If not listed, add it to the list
		if !found {
			members = append(members, contact)
		}
	}

	for _, pk := range members {
		rawPK, err := pk.Raw()
		if err != nil {
			gc.logger.Error("failed to serialize pk", zap.Error(err))
			continue
		}

		if _, err := gc.MetadataStore().SendSecret(gc.ctx, pk); err != nil {
			if !errcode.Is(err, errcode.ErrGroupSecretAlreadySentToMember) {
				gc.logger.Info("secret already sent secret to member", logutil.PrivateString("memberpk", base64.StdEncoding.EncodeToString(rawPK)))
				continue
			}
		} else {
			gc.logger.Info("sent secret to existing member", logutil.PrivateString("memberpk", base64.StdEncoding.EncodeToString(rawPK)))
		}
	}
}

func (gc *GroupContext) TagGroupContextPeers(ipfsCoreAPI ipfsutil.ExtendedCoreAPI, weight int) {
	id := gc.Group().GroupIDAsString()

	chSub1, err := gc.metadataStore.EventBus().Subscribe(new(stores.EventNewPeer))
	if err != nil {
		gc.logger.Warn("unable to subscribe to metadata event new peer")
		return
	}

	chSub2, err := gc.messageStore.EventBus().Subscribe(new(stores.EventNewPeer))
	if err != nil {
		gc.logger.Warn("unable to subscribe to message event new peer")
		return
	}

	go func() {
		defer chSub1.Close()
		defer chSub2.Close()

		for {
			var e interface{}

			select {
			case e = <-chSub1.Out():
			case e = <-chSub2.Out():
			case <-gc.ctx.Done():
				return
			}

			evt := e.(stores.EventNewPeer)

			tag := fmt.Sprintf("grp_%s", id)
			gc.logger.Debug("new peer of interest", logutil.PrivateStringer("peer", evt.Peer), zap.String("tag", tag), zap.Int("score", weight))
			ipfsCoreAPI.ConnMgr().TagPeer(evt.Peer, tag, weight)
		}
	}()
}

func (gc *GroupContext) WaitForDeviceAdded(devicePK []byte) (found chan struct{}) {
	gc.muDevicesAdded.Lock()
	defer gc.muDevicesAdded.Unlock()

	k := string(devicePK)
	var ok bool
	if found, ok = gc.devicesAdded[k]; ok {
		return
	}

	found = make(chan struct{})
	if gc.messageKeystore.HasSecretForRawDevicePK(gc.ctx, gc.group.PublicKey, devicePK) {
		close(found)
		return
	}

	gc.devicesAdded[k] = found
	return
}

func (gc *GroupContext) notifyDeviceAdded(dPK []byte) {
	gc.muDevicesAdded.Lock()
	k := string(dPK)
	if cc, ok := gc.devicesAdded[k]; ok {
		close(cc)
		delete(gc.devicesAdded, k)
	}
	gc.muDevicesAdded.Unlock()
}
