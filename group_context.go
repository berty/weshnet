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
	"google.golang.org/protobuf/proto"

	"berty.tech/go-orbit-db/stores"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/ipfsutil"
	"berty.tech/weshnet/pkg/logutil"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/secretstore"
)

type GroupContext struct {
	ctx             context.Context
	cancel          context.CancelFunc
	group           *protocoltypes.Group
	metadataStore   *MetadataStore
	messageStore    *MessageStore
	secretStore     secretstore.SecretStore
	ownMemberDevice secretstore.OwnMemberDevice

	logger            *zap.Logger
	closed            uint32
	tasks             sync.WaitGroup
	devicesAdded      map[string]chan struct{}
	muDevicesAdded    sync.RWMutex
	selfAnnounced     chan struct{}
	selfAnnouncedOnce sync.Once
}

func (gc *GroupContext) SecretStore() secretstore.SecretStore {
	return gc.secretStore
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
	return gc.ownMemberDevice.Member()
}

func (gc *GroupContext) DevicePubKey() crypto.PubKey {
	return gc.ownMemberDevice.Device()
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

func NewContextGroup(group *protocoltypes.Group, metadataStore *MetadataStore, messageStore *MessageStore, secretStore secretstore.SecretStore, memberDevice secretstore.OwnMemberDevice, logger *zap.Logger) *GroupContext {
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
		secretStore:     secretStore,
		ownMemberDevice: memberDevice,
		logger:          logger.With(logutil.PrivateString("group-id", fmt.Sprintf("%.6s", base64.StdEncoding.EncodeToString(group.PublicKey)))),
		closed:          0,
		devicesAdded:    make(map[string]chan struct{}),
		selfAnnounced:   make(chan struct{}),
	}
}

func (gc *GroupContext) ActivateGroupContext(contactPK crypto.PubKey) (err error) {
	ctx := gc.ctx

	// start watching for GroupMetadataEvent to send secret and register
	// chainkey of new members.
	{
		m := gc.MetadataStore()
		sub, err := m.EventBus().Subscribe(new(*protocoltypes.GroupMetadataEvent))
		if err != nil {
			return fmt.Errorf("unable to subscribe to group metadata event: %w", err)
		}

		gc.tasks.Add(1)
		go func() {
			defer gc.tasks.Done() // ultimately, mark bg task has done
			defer sub.Close()

			for {
				var evt interface{}
				select {
				case <-ctx.Done():
					return
				case evt = <-sub.Out():
				}

				// @TODO(gfanton): should we handle this in a sub gorouting ?
				e := evt.(*protocoltypes.GroupMetadataEvent)
				// start := time.Now()
				if err := gc.handleGroupMetadataEvent(e); err != nil {
					gc.logger.Error("unable to handle EventTypeGroupDeviceSecretAdded", zap.Error(err))
				}

				// if t := time.Since(start).Milliseconds(); t > 0 {
				// 	fmt.Printf("elapsed: %dms\n", t)
				// }
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

		wgExistingMembers.Wait()
	}

	start := time.Now()
	op, err := gc.MetadataStore().AddDeviceToGroup(gc.ctx)
	if err != nil {
		return fmt.Errorf("unable to add device to groupo: %w", err)
	}

	gc.logger.Info(fmt.Sprintf("AddDeviceToGroup took %s", time.Since(start)))
	if op != nil {
		// Waiting for async events to be handled
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-gc.selfAnnounced: // device has been selfAnnounced
		}
	}

	return nil
}

func (gc *GroupContext) handleGroupMetadataEvent(e *protocoltypes.GroupMetadataEvent) (err error) {
	switch e.Metadata.EventType {
	case protocoltypes.EventType_EventTypeGroupMemberDeviceAdded:
		event := &protocoltypes.GroupMemberDeviceAdded{}
		if err := proto.Unmarshal(e.Event, event); err != nil {
			gc.logger.Error("unable to unmarshal payload", zap.Error(err))
		}

		memberPK, err := crypto.UnmarshalEd25519PublicKey(event.MemberPk)
		if err != nil {
			return fmt.Errorf("unable to unmarshal sender member pk: %w", err)
		}

		if memberPK.Equals(gc.ownMemberDevice.Member()) {
			gc.selfAnnouncedOnce.Do(func() { close(gc.selfAnnounced) }) // mark has self announced
		}

		if _, err := gc.MetadataStore().SendSecret(gc.ctx, memberPK); err != nil {
			if !errcode.Is(err, errcode.ErrCode_ErrGroupSecretAlreadySentToMember) {
				return fmt.Errorf("unable to send secret to member: %w", err)
			}
		}

	case protocoltypes.EventType_EventTypeGroupDeviceChainKeyAdded:
		senderPublicKey, encryptedDeviceChainKey, err := getAndFilterGroupDeviceChainKeyAddedPayload(e.Metadata, gc.ownMemberDevice.Member())
		switch err {
		case nil: // ok
		case errcode.ErrCode_ErrInvalidInput, errcode.ErrCode_ErrGroupSecretOtherDestMember:
			// @FIXME(gfanton): should we log this ?
			return nil
		default:
			return fmt.Errorf("an error occurred while opening device secrets: %w", err)
		}

		if err = gc.SecretStore().RegisterChainKey(gc.ctx, gc.Group(), senderPublicKey, encryptedDeviceChainKey); err != nil {
			return fmt.Errorf("unable to register chain key: %w", err)
		}

		if rawPK, err := senderPublicKey.Raw(); err == nil {
			// A new chainKey has been registered, notify watcher
			go gc.notifyDeviceAdded(rawPK)
			// process queued message and check if cached messages can be opened with it
			gc.MessageStore().ProcessMessageQueueForDevicePK(gc.ctx, rawPK)
		}
	}

	return nil
}

func (gc *GroupContext) fillMessageKeysHolderUsingPreviousData() {
	publishedSecrets := gc.metadataStoreListSecrets()

	for senderPublicKey, encryptedSecret := range publishedSecrets {
		if err := gc.SecretStore().RegisterChainKey(gc.ctx, gc.Group(), senderPublicKey, encryptedSecret); err != nil {
			gc.logger.Error("unable to register chain key", zap.Error(err))
			continue
		}
		// A new chainKey is registered, check if cached messages can be opened with it
		if rawPK, err := senderPublicKey.Raw(); err == nil {
			gc.MessageStore().ProcessMessageQueueForDevicePK(gc.ctx, rawPK)
		}
	}
}

func (gc *GroupContext) metadataStoreListSecrets() map[crypto.PubKey][]byte {
	publishedSecrets := map[crypto.PubKey][]byte{}

	m := gc.MetadataStore()

	metadatas, err := m.ListEvents(gc.ctx, nil, nil, false)
	if err != nil {
		return nil
	}
	for metadata := range metadatas {
		if metadata == nil {
			continue
		}

		pk, encryptedDeviceChainKey, err := getAndFilterGroupDeviceChainKeyAddedPayload(metadata.Metadata, gc.MemberPubKey())
		if errcode.Is(err, errcode.ErrCode_ErrInvalidInput) || errcode.Is(err, errcode.ErrCode_ErrGroupSecretOtherDestMember) {
			continue
		}

		if err != nil {
			gc.logger.Error("unable to open chain key", zap.Error(err))
			continue
		}

		publishedSecrets[pk] = encryptedDeviceChainKey
	}

	return publishedSecrets
}

func (gc *GroupContext) sendSecretsToExistingMembers(contact crypto.PubKey) {
	members := gc.MetadataStore().ListMembers()

	// Force sending secret to contact member in contact group
	if gc.group.GroupType == protocoltypes.GroupType_GroupTypeContact && len(members) < 2 && contact != nil {
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
			if !errcode.Is(err, errcode.ErrCode_ErrGroupSecretAlreadySentToMember) {
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

func (gc *GroupContext) WaitForDeviceAdded(ctx context.Context, devicePK crypto.PubKey) (found chan struct{}) {
	gc.muDevicesAdded.Lock()
	defer gc.muDevicesAdded.Unlock()

	rawpk, err := devicePK.Raw()
	if err != nil {
		gc.logger.Error("unable to get raw public key", zap.Error(err))
		return
	}

	k := string(rawpk)
	var ok bool
	if found, ok = gc.devicesAdded[k]; ok {
		return
	}

	groupPublicKey, err := gc.group.GetPubKey()
	if err != nil {
		gc.logger.Error("unable to get group public key", zap.Error(err))
		return
	}

	found = make(chan struct{})
	if gc.secretStore.IsChainKeyKnownForDevice(ctx, groupPublicKey, devicePK) {
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
