package weshnet

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	ds_sync "github.com/ipfs/go-datastore/sync"
	coreapi "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p/core/crypto"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/host/eventbus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	ipfslog "berty.tech/go-ipfs-log"
	"berty.tech/go-ipfs-log/enc"
	"berty.tech/go-ipfs-log/entry"
	"berty.tech/go-ipfs-log/identityprovider"
	"berty.tech/go-ipfs-log/io"
	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/baseorbitdb"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/pubsub/pubsubcoreapi"
	"berty.tech/go-orbit-db/stores"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/ipfsutil"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/rendezvous"
	"berty.tech/weshnet/pkg/secretstore"
	"berty.tech/weshnet/pkg/tyber"
)

type GroupOpenMode uint64

const (
	GroupOpenModeUndefined GroupOpenMode = iota
	GroupOpenModeReplicate
	GroupOpenModeWrite
)

var _ = GroupOpenModeUndefined

type loggable interface {
	setLogger(*zap.Logger)
}

type NewOrbitDBOptions struct {
	baseorbitdb.NewOrbitDBOptions
	Datastore          datastore.Batching
	SecretStore        secretstore.SecretStore
	RotationInterval   *rendezvous.RotationInterval
	PrometheusRegister prometheus.Registerer

	GroupMetadataStoreType string
	GroupMessageStoreType  string
	ReplicationMode        bool
}

func (n *NewOrbitDBOptions) applyDefaults() {
	if n.Datastore == nil {
		n.Datastore = ds_sync.MutexWrap(datastore.NewMapDatastore())
	}

	if n.PrometheusRegister == nil {
		n.PrometheusRegister = prometheus.DefaultRegisterer
	}

	if n.Cache == nil {
		n.Cache = NewOrbitDatastoreCache(n.Datastore)
	}

	if n.Logger == nil {
		n.Logger = zap.NewNop()
	}

	if n.RotationInterval == nil {
		n.RotationInterval = rendezvous.NewStaticRotationInterval()
	}

	if n.Logger == nil {
		n.Logger = zap.NewNop()
	}
	n.Logger = n.Logger.Named("odb")

	if n.GroupMetadataStoreType == "" {
		n.GroupMetadataStoreType = "wesh_group_metadata"
	}

	if n.GroupMessageStoreType == "" {
		n.GroupMessageStoreType = "wesh_group_messages"
	}
}

type (
	GroupMap           = sync.Map
	GroupContextMap    = sync.Map
	GroupsSigPubKeyMap = sync.Map
)

type WeshOrbitDB struct {
	baseorbitdb.BaseOrbitDB
	keyStore           *BertySignedKeyStore
	secretStore        secretstore.SecretStore
	pubSub             iface.PubSubInterface
	rotationInterval   *rendezvous.RotationInterval
	messageMarshaler   *OrbitDBMessageMarshaler
	replicationMode    bool
	prometheusRegister prometheus.Registerer

	groupMetadataStoreType string
	groupMessageStoreType  string

	ctx context.Context
	// FIXME(gfanton): use real map instead of sync.Map
	groups          *GroupMap           // map[string]*protocoltypes.Group
	groupContexts   *GroupContextMap    // map[string]*GroupContext
	groupsSigPubKey *GroupsSigPubKeyMap // map[string]crypto.PubKey
}

func (s *WeshOrbitDB) registerGroupPrivateKey(g *protocoltypes.Group) error {
	groupID := g.GroupIDAsString()

	gSigSK, err := g.GetSigningPrivKey()
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	if err := s.SetGroupSigPubKey(groupID, gSigSK.GetPublic()); err != nil {
		return errcode.TODO.Wrap(err)
	}

	if err := s.keyStore.SetKey(gSigSK); err != nil {
		return errcode.TODO.Wrap(err)
	}

	return nil
}

func (s *WeshOrbitDB) registerGroupSigningPubKey(g *protocoltypes.Group) error {
	groupID := g.GroupIDAsString()

	var gSigPK crypto.PubKey

	gSigSK, err := g.GetSigningPrivKey()
	if err == nil && gSigSK != nil {
		gSigPK = gSigSK.GetPublic()
	} else {
		gSigPK, err = g.GetSigningPubKey()
		if err != nil {
			return errcode.TODO.Wrap(err)
		}
	}

	if err := s.SetGroupSigPubKey(groupID, gSigPK); err != nil {
		return errcode.TODO.Wrap(err)
	}

	return nil
}

func NewWeshOrbitDB(ctx context.Context, ipfs coreapi.CoreAPI, options *NewOrbitDBOptions) (*WeshOrbitDB, error) {
	var err error

	if options == nil {
		options = &NewOrbitDBOptions{}
	}

	options.applyDefaults()

	ks := &BertySignedKeyStore{}
	options.Keystore = ks
	options.Identity = &identityprovider.Identity{}

	self, err := ipfs.Key().Self(ctx)
	if err != nil {
		return nil, err
	}

	if options.PubSub == nil {
		options.PubSub = pubsubcoreapi.NewPubSub(ipfs, self.ID(), time.Second, options.Logger, options.Tracer)
	}

	mm := NewOrbitDBMessageMarshaler(self.ID(), options.SecretStore, options.RotationInterval, options.ReplicationMode)
	options.MessageMarshaler = mm

	orbitDB, err := baseorbitdb.NewOrbitDB(ctx, ipfs, &options.NewOrbitDBOptions)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	bertyDB := &WeshOrbitDB{
		ctx:                    ctx,
		messageMarshaler:       mm,
		BaseOrbitDB:            orbitDB,
		keyStore:               ks,
		secretStore:            options.SecretStore,
		rotationInterval:       options.RotationInterval,
		pubSub:                 options.PubSub,
		groups:                 &GroupMap{},
		groupContexts:          &GroupContextMap{},    // map[string]*GroupContext
		groupsSigPubKey:        &GroupsSigPubKeyMap{}, // map[string]crypto.PubKey
		groupMetadataStoreType: options.GroupMetadataStoreType,
		groupMessageStoreType:  options.GroupMessageStoreType,
		replicationMode:        options.ReplicationMode,
		prometheusRegister:     options.PrometheusRegister,
	}

	if err := bertyDB.RegisterAccessControllerType(NewSimpleAccessController); err != nil {
		return nil, errcode.TODO.Wrap(err)
	}
	bertyDB.RegisterStoreType(bertyDB.groupMetadataStoreType, constructorFactoryGroupMetadata(bertyDB, options.Logger))
	bertyDB.RegisterStoreType(bertyDB.groupMessageStoreType, constructorFactoryGroupMessage(bertyDB, options.Logger))

	return bertyDB, nil
}

func (s *WeshOrbitDB) openAccountGroup(ctx context.Context, options *orbitdb.CreateDBOptions, ipfsCoreAPI ipfsutil.ExtendedCoreAPI) (*GroupContext, error) {
	l := s.Logger()

	if options == nil {
		options = &orbitdb.CreateDBOptions{}
	}

	if options.EventBus == nil {
		options.EventBus = s.EventBus()
	}

	group, _, err := s.secretStore.GetGroupForAccount()
	if err != nil {
		return nil, errcode.ErrOrbitDBOpen.Wrap(err)
	}

	l.Debug("Got account group", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "Group", Description: group.String()}})...)

	gc, err := s.OpenGroup(ctx, group, options)
	if err != nil {
		return nil, errcode.ErrGroupOpen.Wrap(err)
	}
	l.Debug("Opened account group", tyber.FormatStepLogFields(ctx, []tyber.Detail{})...)

	gc.TagGroupContextPeers(ipfsCoreAPI, 84)

	if err := gc.ActivateGroupContext(nil); err != nil {
		return nil, errcode.TODO.Wrap(err)
	}
	l.Debug("TagGroupContextPeers done", tyber.FormatStepLogFields(ctx, []tyber.Detail{})...)
	return gc, nil
}

func (s *WeshOrbitDB) setHeadsForGroup(ctx context.Context, g *protocoltypes.Group, metaHeads, messageHeads []cid.Cid) error {
	groupID := g.GroupIDAsString()

	var (
		err                    error
		metaImpl, messagesImpl orbitdb.Store
	)

	existingGC, err := s.getGroupContext(groupID)
	if err != nil && !errcode.Is(err, errcode.ErrMissingMapKey) {
		return errcode.ErrInternal.Wrap(err)
	}
	if err == nil {
		metaImpl = existingGC.metadataStore
		messagesImpl = existingGC.messageStore
	}
	if metaImpl == nil || messagesImpl == nil {
		s.groups.Store(groupID, g)

		if err := s.registerGroupSigningPubKey(g); err != nil {
			return errcode.ErrInternal.Wrap(err)
		}

		s.Logger().Debug("OpenGroup", zap.Any("public key", g.PublicKey), zap.Any("secret", g.Secret), zap.Stringer("type", g.GroupType))

		if metaImpl == nil {
			metaImpl, err = s.storeForGroup(ctx, s, g, nil, s.groupMetadataStoreType, GroupOpenModeReplicate)
			if err != nil {
				return errcode.ErrOrbitDBOpen.Wrap(err)
			}

			defer func() { _ = metaImpl.Close() }()
		}

		if messagesImpl == nil {
			messagesImpl, err = s.storeForGroup(ctx, s, g, nil, s.groupMessageStoreType, GroupOpenModeReplicate)
			if err != nil {
				return errcode.ErrOrbitDBOpen.Wrap(err)
			}

			defer func() { _ = messagesImpl.Close() }()
		}
	}

	if messagesImpl == nil {
		return errcode.ErrInternal.Wrap(fmt.Errorf("message store is nil"))
	}

	if metaImpl == nil {
		return errcode.ErrInternal.Wrap(fmt.Errorf("metadata store is nil"))
	}

	var wg sync.WaitGroup

	// load and wait heads for metadata and message stores
	wg.Add(2)

	go func() {
		// load meta heads
		if err := s.loadHeads(ctx, metaImpl, metaHeads); err != nil {
			s.Logger().Error("unable to load metadata heads", zap.Error(err))
		}

		wg.Done()
	}()

	go func() {
		// load message heads
		if err := s.loadHeads(ctx, messagesImpl, messageHeads); err != nil {
			s.Logger().Error("unable to load message heads", zap.Error(err))
		}

		wg.Done()
	}()

	wg.Wait()

	return nil
}

func (s *WeshOrbitDB) loadHeads(ctx context.Context, store iface.Store, heads []cid.Cid) (err error) {
	sub, err := store.EventBus().Subscribe(new(stores.EventReplicated),
		eventbus.Name("weshnet/load-heads"))
	if err != nil {
		return fmt.Errorf("unable to subscribe to EventReplicated")
	}
	defer sub.Close()

	// check and generate missing entries if needed
	headsEntries := make([]ipfslog.Entry, len(heads))
	for i, h := range heads {
		if _, ok := store.OpLog().Get(h); !ok {
			headsEntries[i] = &entry.Entry{Hash: h}
		}
	}

	if len(headsEntries) == 0 {
		return nil
	}

	store.Replicator().Load(ctx, headsEntries)

	for found := 0; found < len(heads); {
		// wait for load to finish
		select {
		case e := <-sub.Out():
			evt := e.(stores.EventReplicated)

			// iterate over entries from replicated event to search for our heads
			for _, headEntry := range headsEntries {
				for _, evtEntry := range evt.Entries {
					if evtEntry.Equals(headEntry) {
						found++
						break
					}
				}
			}

		case <-s.ctx.Done():
			return s.ctx.Err()
		}
	}

	return nil
}

func (s *WeshOrbitDB) OpenGroup(ctx context.Context, g *protocoltypes.Group, options *orbitdb.CreateDBOptions) (*GroupContext, error) {
	if s.secretStore == nil {
		return nil, errcode.ErrInvalidInput.Wrap(fmt.Errorf("db open in naive mode"))
	}

	groupID := g.GroupIDAsString()

	existingGC, err := s.getGroupContext(groupID)
	if err != nil && !errcode.Is(err, errcode.ErrMissingMapKey) {
		return nil, errcode.ErrInternal.Wrap(err)
	}
	if err == nil {
		return existingGC, nil
	}

	s.groups.Store(groupID, g)

	if err := s.registerGroupPrivateKey(g); err != nil {
		return nil, err
	}

	s.Logger().Debug("OpenGroup", tyber.FormatStepLogFields(s.ctx, tyber.ZapFieldsToDetails(zap.Any("public key", g.PublicKey), zap.Any("secret", g.Secret), zap.Stringer("type", g.GroupType)))...)

	memberDevice, err := s.secretStore.GetOwnMemberDeviceForGroup(g)
	if err != nil {
		return nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	mpkb, err := crypto.MarshalPublicKey(memberDevice.Member())
	if err != nil {
		mpkb = []byte{}
	}
	s.Logger().Debug("Got member device", tyber.FormatStepLogFields(s.ctx, []tyber.Detail{{Name: "DevicePublicKey", Description: base64.RawURLEncoding.EncodeToString(mpkb)}})...)

	// Force secret generation if missing
	if _, err := s.secretStore.GetShareableChainKey(s.ctx, g, memberDevice.Member()); err != nil {
		return nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	s.Logger().Debug("Got device chain key", tyber.FormatStepLogFields(s.ctx, []tyber.Detail{})...)

	metaImpl, err := s.groupMetadataStore(ctx, g, options)
	if err != nil {
		return nil, errcode.ErrOrbitDBOpen.Wrap(err)
	}
	s.messageMarshaler.RegisterGroup(metaImpl.Address().String(), g)

	s.Logger().Debug("Got metadata store", tyber.FormatStepLogFields(s.ctx, []tyber.Detail{})...)

	// force to unshare the same EventBus between groupMetadataStore and groupMessageStore
	// to avoid having a bunch of events which are not for the correct group
	if options != nil && options.EventBus != nil {
		options.EventBus = eventbus.NewBus(
			eventbus.WithMetricsTracer(eventbus.NewMetricsTracer(eventbus.WithRegisterer(s.prometheusRegister))))
	}

	messagesImpl, err := s.groupMessageStore(ctx, g, options)
	if err != nil {
		metaImpl.Close()
		return nil, errcode.ErrOrbitDBOpen.Wrap(err)
	}
	s.messageMarshaler.RegisterGroup(messagesImpl.Address().String(), g)

	s.Logger().Debug("Got message store", tyber.FormatStepLogFields(s.ctx, []tyber.Detail{})...)

	gc := NewContextGroup(g, metaImpl, messagesImpl, s.secretStore, memberDevice, s.Logger())

	s.Logger().Debug("Created group context", tyber.FormatStepLogFields(s.ctx, []tyber.Detail{})...)

	s.groupContexts.Store(groupID, gc)

	s.Logger().Debug("Stored group context", tyber.FormatStepLogFields(s.ctx, []tyber.Detail{})...)

	return gc, nil
}

func (s *WeshOrbitDB) OpenGroupReplication(ctx context.Context, g *protocoltypes.Group, options *orbitdb.CreateDBOptions) (iface.Store, iface.Store, error) {
	if g == nil || len(g.PublicKey) == 0 {
		return nil, nil, errcode.ErrInvalidInput.Wrap(fmt.Errorf("missing group or group pubkey"))
	}

	groupID := g.GroupIDAsString()

	gc, err := s.getGroupContext(groupID)
	if err != nil && !errcode.Is(err, errcode.ErrMissingMapKey) {
		return nil, nil, errcode.ErrInternal.Wrap(err)
	}
	if err == nil {
		return gc.metadataStore, gc.messageStore, nil
	}

	s.groups.Store(groupID, g)

	if err := s.registerGroupSigningPubKey(g); err != nil {
		return nil, nil, err
	}

	metadataStore, err := s.storeForGroup(ctx, s, g, options, s.groupMetadataStoreType, GroupOpenModeReplicate)
	if err != nil {
		_ = metadataStore.Close()
		return nil, nil, errors.Wrap(err, "unable to open database")
	}

	messageStore, err := s.storeForGroup(ctx, s, g, options, s.groupMessageStoreType, GroupOpenModeReplicate)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to open database")
	}

	return metadataStore, messageStore, nil
}

func (s *WeshOrbitDB) getGroupContext(id string) (*GroupContext, error) {
	g, ok := s.groupContexts.Load(id)
	if !ok {
		return nil, errcode.ErrMissingMapKey
	}

	gc, ok := g.(*GroupContext)
	if !ok {
		s.groupContexts.Delete(id)
		return nil, errors.New("cannot cast object to GroupContext")
	}

	if gc.IsClosed() {
		s.groupContexts.Delete(id)
		return nil, errcode.ErrMissingMapKey
	}

	return g.(*GroupContext), nil
}

// SetGroupSigPubKey registers a new group signature pubkey, mainly used to
// replicate a store data without needing to access to its content
func (s *WeshOrbitDB) SetGroupSigPubKey(groupID string, pubKey crypto.PubKey) error {
	if pubKey == nil {
		return errcode.ErrInvalidInput
	}

	s.groupsSigPubKey.Store(groupID, pubKey)

	return nil
}

func (s *WeshOrbitDB) storeForGroup(ctx context.Context, o iface.BaseOrbitDB, g *protocoltypes.Group, options *orbitdb.CreateDBOptions, storeType string, groupOpenMode GroupOpenMode) (iface.Store, error) {
	l := s.Logger()

	if options == nil {
		options = &orbitdb.CreateDBOptions{}
	}

	// setup eventbus metrics
	if options.EventBus == nil {
		options.EventBus = eventbus.NewBus(eventbus.WithMetricsTracer(eventbus.NewMetricsTracer(eventbus.WithRegisterer(s.prometheusRegister))))
	}

	options, err := DefaultOrbitDBOptions(g, options, s.keyStore, storeType, groupOpenMode)
	if err != nil {
		return nil, err
	}

	l.Debug("Opening store", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "Group", Description: g.String()}, {Name: "Options", Description: fmt.Sprint(options)}}, tyber.Status(tyber.Running))...)

	options.StoreType = &storeType
	name := fmt.Sprintf("%s_%s", g.GroupIDAsString(), storeType)

	addr, err := o.DetermineAddress(ctx, name, storeType, &orbitdb.DetermineAddressOptions{AccessController: options.AccessController})
	if err != nil {
		return nil, err
	}

	s.messageMarshaler.RegisterGroup(addr.String(), g)

	linkKey, err := g.GetLinkKeyArray()
	if err != nil {
		return nil, err
	}

	if key := linkKey[:]; len(key) > 0 {
		sk, err := enc.NewSecretbox(key)
		if err != nil {
			return nil, err
		}

		cborIO := io.CBOR()
		cborIO.ApplyOptions(&io.CBOROptions{LinkKey: sk})
		options.IO = cborIO

		l.Debug("opening store: register rotation", zap.String("topic", addr.String()))

		s.messageMarshaler.RegisterSharedKeyForTopic(addr.String(), sk)
		s.rotationInterval.RegisterRotation(time.Now(), addr.String(), key)
	}

	store, err := o.Open(ctx, name, options)
	if err != nil {
		return nil, errcode.ErrOrbitDBOpen.Wrap(err)
	}

	l.Debug("Loading store", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "Group", Description: g.String()}, {Name: "StoreType", Description: store.Type()}, {Name: "Store", Description: store.Address().String()}}, tyber.Status(tyber.Running))...)

	_ = store.Load(ctx, -1)

	l.Debug("Loaded store", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "Group", Description: g.String()}})...)

	return store, nil
}

func (s *WeshOrbitDB) groupMetadataStore(ctx context.Context, g *protocoltypes.Group, options *orbitdb.CreateDBOptions) (*MetadataStore, error) {
	if options == nil {
		options = &orbitdb.CreateDBOptions{}
	}

	l := s.Logger().Named("metadataStore")
	options.Logger = l

	l.Debug("Opening group metadata store", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "Group", Description: g.String()}, {Name: "Options", Description: fmt.Sprint(options)}}, tyber.Status(tyber.Running))...)

	store, err := s.storeForGroup(ctx, s, g, options, s.groupMetadataStoreType, GroupOpenModeWrite)
	if err != nil {
		return nil, tyber.LogFatalError(ctx, l, "Failed to get group store", errors.Wrap(err, "unable to open database"))
	}

	l.Debug("Got group store", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "DBName", Description: store.DBName()}})...)

	sStore, ok := store.(*MetadataStore)
	if !ok {
		return nil, tyber.LogFatalError(ctx, l, "Failed to cast group store", errors.New("unable to cast store to metadata store"))
	}

	l.Debug("Opened group metadata store", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "Group", Description: g.String()}})...)

	return sStore, nil
}

func (s *WeshOrbitDB) groupMessageStore(ctx context.Context, g *protocoltypes.Group, options *orbitdb.CreateDBOptions) (*MessageStore, error) {
	if options == nil {
		options = &orbitdb.CreateDBOptions{}
	}

	l := s.Logger().Named("messageStore")
	options.Logger = l

	l.Debug("Opening group message store", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "Group", Description: g.String()}, {Name: "Options", Description: fmt.Sprint(options)}}, tyber.Status(tyber.Running))...)

	store, err := s.storeForGroup(ctx, s, g, options, s.groupMessageStoreType, GroupOpenModeWrite)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open database")
	}

	mStore, ok := store.(*MessageStore)
	if !ok {
		return nil, errors.New("unable to cast store to message store")
	}

	return mStore, nil
}

func (s *WeshOrbitDB) getGroupFromOptions(options *iface.NewStoreOptions) (*protocoltypes.Group, error) {
	groupIDs, err := options.AccessController.GetAuthorizedByRole(identityGroupIDKey)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	if len(groupIDs) != 1 {
		return nil, errcode.ErrInvalidInput
	}

	g, ok := s.groups.Load(groupIDs[0])
	if !ok {
		return nil, errcode.ErrInvalidInput
	}

	typed, ok := g.(*protocoltypes.Group)
	if !ok {
		return nil, errcode.ErrInvalidInput
	}

	return typed, nil
}

func (s *WeshOrbitDB) IsGroupLoaded(groupID string) bool {
	gc, ok := s.groups.Load(groupID)

	return ok && gc != nil
}

func (s *WeshOrbitDB) GetDevicePKForPeerID(id peer.ID) (pdg *PeerDeviceGroup, ok bool) {
	return s.messageMarshaler.GetDevicePKForPeerID(id)
}
