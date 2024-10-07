package weshnet

import (
	"context"
	"encoding/hex"
	"fmt"
	mrand "math/rand"
	"path/filepath"
	"sync"
	"time"
	"unsafe"

	"github.com/dgraph-io/badger/v2/options"
	ds "github.com/ipfs/go-datastore"
	ds_sync "github.com/ipfs/go-datastore/sync"
	badger "github.com/ipfs/go-ds-badger2"
	coreiface "github.com/ipfs/kubo/core/coreiface"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	backoff "github.com/libp2p/go-libp2p/p2p/discovery/backoff"
	"github.com/libp2p/go-libp2p/p2p/host/eventbus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"moul.io/srand"

	"berty.tech/go-orbit-db/baseorbitdb"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/pubsub/directchannel"
	"berty.tech/go-orbit-db/pubsub/pubsubraw"
	"berty.tech/weshnet/internal/bertyversion"
	"berty.tech/weshnet/internal/datastoreutil"
	"berty.tech/weshnet/pkg/bertyvcissuer"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/ipfsutil"
	ipfs_mobile "berty.tech/weshnet/pkg/ipfsutil/mobile"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/secretstore"
	tinder "berty.tech/weshnet/pkg/tinder"
	"berty.tech/weshnet/pkg/tyber"
)

var _ Service = (*service)(nil)

// Service is the main Berty Protocol interface
type Service interface {
	protocoltypes.ProtocolServiceServer

	Close() error
	Status() Status
	IpfsCoreAPI() coreiface.CoreAPI
}

type service struct {
	// variables
	ctx                    context.Context
	ctxCancel              context.CancelFunc
	logger                 *zap.Logger
	ipfsCoreAPI            ipfsutil.ExtendedCoreAPI
	odb                    *WeshOrbitDB
	accountGroupCtx        *GroupContext
	openedGroups           map[string]*GroupContext
	lock                   sync.RWMutex
	close                  func() error
	startedAt              time.Time
	host                   host.Host
	grpcInsecure           bool
	refreshprocess         map[string]context.CancelFunc
	muRefreshprocess       sync.RWMutex
	swiper                 *Swiper
	peerStatusManager      *ConnectednessManager
	accountEventBus        event.Bus
	contactRequestsManager *contactRequestsManager
	vcClient               *bertyvcissuer.Client
	secretStore            secretstore.SecretStore

	protocoltypes.UnimplementedProtocolServiceServer
}

// Opts contains optional configuration flags for building a new Client
type Opts struct {
	Logger             *zap.Logger
	IpfsCoreAPI        ipfsutil.ExtendedCoreAPI
	DatastoreDir       string
	RootDatastore      ds.Batching
	OrbitDB            *WeshOrbitDB
	TinderService      *tinder.Service
	Host               host.Host
	PubSub             *pubsub.PubSub
	GRPCInsecureMode   bool
	LocalOnly          bool
	close              func() error
	SecretStore        secretstore.SecretStore
	PrometheusRegister prometheus.Registerer

	// These are used if OrbitDB is nil.
	GroupMetadataStoreType string
	GroupMessageStoreType  string
}

func (opts *Opts) applyPushDefaults() {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	if opts.PrometheusRegister == nil {
		opts.PrometheusRegister = prometheus.DefaultRegisterer
	}
}

func (opts *Opts) applyDefaultsGetDatastore() error {
	if opts.RootDatastore == nil {
		if opts.DatastoreDir == "" || opts.DatastoreDir == InMemoryDirectory {
			opts.RootDatastore = ds_sync.MutexWrap(ds.NewMapDatastore())
		} else {
			bopts := badger.DefaultOptions
			bopts.ValueLogLoadingMode = options.FileIO

			ds, err := badger.NewDatastore(opts.DatastoreDir, &bopts)
			if err != nil {
				return fmt.Errorf("unable to init badger datastore: %w", err)
			}
			opts.RootDatastore = ds

			oldClose := opts.close
			opts.close = func() error {
				var err error
				if oldClose != nil {
					err = oldClose()
				}

				if dserr := ds.Close(); dserr != nil {
					err = multierr.Append(err, fmt.Errorf("unable to close datastore: %w", dserr))
				}

				return err
			}
		}
	}

	return nil
}

func (opts *Opts) applyDefaults(ctx context.Context) error {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	rng := mrand.New(mrand.NewSource(srand.MustSecure())) // nolint:gosec // we need to use math/rand here, but it is seeded from crypto/rand

	if err := opts.applyDefaultsGetDatastore(); err != nil {
		return err
	}

	opts.applyPushDefaults()

	if opts.SecretStore == nil {
		secretStore, err := secretstore.NewSecretStore(opts.RootDatastore, &secretstore.NewSecretStoreOptions{
			Logger: opts.Logger,
		})
		if err != nil {
			return errcode.ErrCode_ErrInternal.Wrap(err)
		}

		opts.SecretStore = secretStore
	}

	var mnode *ipfs_mobile.IpfsMobile
	if opts.IpfsCoreAPI == nil {
		dsync := opts.RootDatastore
		if dsync == nil {
			dsync = ds_sync.MutexWrap(ds.NewMapDatastore())
		}

		repo, err := ipfsutil.CreateMockedRepo(dsync)
		if err != nil {
			return err
		}

		mrepo := ipfs_mobile.NewRepoMobile("", repo)
		mnode, err = ipfsutil.NewIPFSMobile(ctx, mrepo, &ipfsutil.MobileOptions{})
		if err != nil {
			return err
		}

		opts.IpfsCoreAPI, err = ipfsutil.NewExtendedCoreAPIFromNode(mnode.IpfsNode)
		if err != nil {
			return err
		}
		opts.Host = mnode.PeerHost()

		oldClose := opts.close
		opts.close = func() error {
			if oldClose != nil {
				_ = oldClose()
			}

			return mnode.Close()
		}
	}

	if opts.Host == nil {
		opts.Host = opts.IpfsCoreAPI
	}

	// setup default tinder service
	if opts.TinderService == nil {
		drivers := []tinder.IDriver{}

		// setup loac disc
		localdisc, err := tinder.NewLocalDiscovery(opts.Logger, opts.Host, rng)
		if err != nil {
			return fmt.Errorf("unable to setup tinder localdiscovery: %w", err)
		}
		drivers = append(drivers, localdisc)

		if mnode != nil {
			dhtdisc := tinder.NewRoutingDiscoveryDriver("dht", mnode.DHT)
			drivers = append(drivers, dhtdisc)
		}

		opts.TinderService, err = tinder.NewService(opts.Host, opts.Logger, drivers...)
		if err != nil {
			return fmt.Errorf("unable to setup tinder service: %w", err)
		}
	}

	if opts.PubSub == nil {
		var err error

		popts := []pubsub.Option{
			pubsub.WithMessageSigning(true),
			pubsub.WithPeerExchange(true),
		}

		backoffstrat := backoff.NewExponentialBackoff(
			time.Second*10, time.Hour,
			backoff.FullJitter,
			time.Second, 10.0, 0, rng)

		cacheSize := 100
		dialTimeout := time.Second * 20
		backoffconnector := func(host host.Host) (*backoff.BackoffConnector, error) {
			return backoff.NewBackoffConnector(host, cacheSize, dialTimeout, backoffstrat)
		}

		adaptater := tinder.NewDiscoveryAdaptater(opts.Logger.Named("disc"), opts.TinderService)
		popts = append(popts, pubsub.WithDiscovery(adaptater, pubsub.WithDiscoverConnector(backoffconnector)))

		// pubsub.DiscoveryPollInterval = m.Node.Protocol.PollInterval
		ps, err := pubsub.NewGossipSub(ctx, opts.Host, popts...)
		if err != nil {
			return fmt.Errorf("unable to init gossipsub: %w", err)
		}

		// @NOTE(gfanton): we need to force cast here until our fix is push
		// upstream on the original go-libp2p-pubsub
		// see: https://github.com/gfanton/go-libp2p-pubsub/commit/8f4fd394f8dfcb3a5eb724a03f9e4e1e33194cbd
		opts.PubSub = (*pubsub.PubSub)(unsafe.Pointer(ps))
	}

	if opts.OrbitDB == nil {
		orbitDirectory := InMemoryDirectory
		if opts.DatastoreDir != InMemoryDirectory {
			orbitDirectory = filepath.Join(opts.DatastoreDir, NamespaceOrbitDBDirectory)
		}

		pubsub := pubsubraw.NewPubSub(opts.PubSub, opts.Host.ID(), opts.Logger, nil)
		odbOpts := &NewOrbitDBOptions{
			NewOrbitDBOptions: baseorbitdb.NewOrbitDBOptions{
				Directory: &orbitDirectory,
				PubSub:    pubsub,
				Logger:    opts.Logger,
			},
			PrometheusRegister:     opts.PrometheusRegister,
			Datastore:              datastoreutil.NewNamespacedDatastore(opts.RootDatastore, ds.NewKey(NamespaceOrbitDBDatastore)),
			SecretStore:            opts.SecretStore,
			GroupMetadataStoreType: opts.GroupMetadataStoreType,
			GroupMessageStoreType:  opts.GroupMessageStoreType,
		}

		if opts.Host != nil {
			odbOpts.DirectChannelFactory = directchannel.InitDirectChannelFactory(opts.Logger, opts.Host)
		}

		odb, err := NewWeshOrbitDB(ctx, opts.IpfsCoreAPI, odbOpts)
		if err != nil {
			return err
		}

		oldClose := opts.close
		opts.close = func() error {
			if oldClose != nil {
				_ = oldClose()
			}

			return odb.Close()
		}

		opts.OrbitDB = odb
	}

	return nil
}

// NewService initializes a new Service using the opts.
// If opts.RootDatastore is nil and opts.DatastoreDir is "" or InMemoryDirectory, then set
// opts.RootDatastore to an in-memory data store. Otherwise, if opts.RootDatastore is nil then set
// opts.RootDatastore to a persistent data store at opts.DatastoreDir .
func NewService(opts Opts) (_ Service, err error) {
	ctx, cancel := context.WithCancel(context.Background())

	if err := opts.applyDefaults(ctx); err != nil {
		cancel()
		return nil, errcode.ErrCode_TODO.Wrap(err)
	}

	opts.Logger = opts.Logger.Named("pt")

	ctx, _, endSection := tyber.Section(tyber.ContextWithoutTraceID(ctx), opts.Logger, fmt.Sprintf("Initializing ProtocolService version %s", bertyversion.Version))
	defer func() { endSection(err, "") }()

	accountEventBus := eventbus.NewBus(
		eventbus.WithMetricsTracer(eventbus.NewMetricsTracer(eventbus.WithRegisterer(opts.PrometheusRegister))))

	dbOpts := &iface.CreateDBOptions{
		EventBus:  accountEventBus,
		LocalOnly: &opts.LocalOnly,
	}

	accountGroupCtx, err := opts.OrbitDB.openAccountGroup(ctx, dbOpts, opts.IpfsCoreAPI)
	if err != nil {
		cancel()
		return nil, errcode.ErrCode_TODO.Wrap(err)
	}

	opts.Logger.Debug("Opened account group", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "AccountGroup", Description: accountGroupCtx.group.String()}})...)

	var contactRequestsManager *contactRequestsManager
	var swiper *Swiper
	if opts.TinderService != nil {
		swiper = NewSwiper(opts.Logger, opts.TinderService, opts.OrbitDB.rotationInterval)
		opts.Logger.Debug("Tinder swiper is enabled", tyber.FormatStepLogFields(ctx, []tyber.Detail{})...)

		if contactRequestsManager, err = newContactRequestsManager(swiper, accountGroupCtx.metadataStore, opts.IpfsCoreAPI, opts.Logger); err != nil {
			cancel()
			return nil, errcode.ErrCode_TODO.Wrap(err)
		}
	} else {
		opts.Logger.Warn("No tinder driver provided, incoming and outgoing contact requests won't be enabled", tyber.FormatStepLogFields(ctx, []tyber.Detail{})...)
	}

	if err := opts.SecretStore.PutGroup(ctx, accountGroupCtx.Group()); err != nil {
		cancel()
		return nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unable to add account group to group datastore, err: %w", err))
	}

	s := &service{
		ctx:             ctx,
		ctxCancel:       cancel,
		host:            opts.Host,
		ipfsCoreAPI:     opts.IpfsCoreAPI,
		logger:          opts.Logger,
		odb:             opts.OrbitDB,
		close:           opts.close,
		accountGroupCtx: accountGroupCtx,
		swiper:          swiper,
		startedAt:       time.Now(),
		openedGroups: map[string]*GroupContext{
			string(accountGroupCtx.Group().PublicKey): accountGroupCtx,
		},
		secretStore:            opts.SecretStore,
		grpcInsecure:           opts.GRPCInsecureMode,
		refreshprocess:         make(map[string]context.CancelFunc),
		peerStatusManager:      NewConnectednessManager(),
		accountEventBus:        accountEventBus,
		contactRequestsManager: contactRequestsManager,
	}

	s.startGroupDeviceMonitor()

	return s, nil
}

func (s *service) IpfsCoreAPI() coreiface.CoreAPI {
	return s.ipfsCoreAPI
}

func (s *service) Close() error {
	endSection := tyber.SimpleSection(tyber.ContextWithoutTraceID(s.ctx), s.logger, "Closing ProtocolService")

	var err error
	pks := []crypto.PubKey{}

	// gather public keys
	s.lock.Lock()

	if s.contactRequestsManager != nil {
		s.contactRequestsManager.close()
		s.contactRequestsManager = nil
	}

	for _, gc := range s.openedGroups {
		pk, subErr := gc.group.GetPubKey()
		if subErr != nil {
			err = multierr.Append(err, subErr)
			continue
		}

		pks = append(pks, pk)
	}
	s.lock.Unlock()

	// deactivate all groups
	for _, pk := range pks {
		derr := s.deactivateGroup(pk)
		if derr != nil {
			err = multierr.Append(derr, derr)
		}
	}

	err = multierr.Append(err, s.odb.Close())

	if s.close != nil {
		err = multierr.Append(err, s.close())
	}

	endSection(err)

	s.ctxCancel()

	return err
}

func (s *service) startGroupDeviceMonitor() {
	if s.host == nil {
		return
	}

	// monitor exchange heads events
	subHead, err := s.odb.EventBus().Subscribe(new(baseorbitdb.EventExchangeHeads),
		eventbus.Name("weshnet/service/monitor-exchange-heads"))
	if err != nil {
		s.logger.Error("startGroupDeviceMonitor", zap.Error(errors.Wrap(err, "unable to subscribe odb event")))
		return
	}

	// monitor peer connectednesschanged
	subPeer, err := s.host.EventBus().Subscribe(new(event.EvtPeerConnectednessChanged),
		eventbus.Name("weshnet/service/monitor-peer-connectedness"))
	if err != nil {
		s.logger.Error("startGroupDeviceMonitor", zap.Error(errors.Wrap(err, "unable to subscribe odb event")))
		subHead.Close()
		return
	}

	go func() {
		defer subHead.Close()
		defer subPeer.Close()

		for {
			var evt interface{}

			select {
			case evt = <-subHead.Out():
			case evt = <-subPeer.Out():
			case <-s.ctx.Done():
				return
			}

			switch e := evt.(type) {
			case event.EvtPeerConnectednessChanged:
				switch e.Connectedness {
				case network.Connected:
					s.peerStatusManager.UpdateState(e.Peer, ConnectednessTypeConnected)
				case network.NotConnected:
					s.peerStatusManager.UpdateState(e.Peer, ConnectednessTypeDisconnected)
				}
			case baseorbitdb.EventExchangeHeads:
				if dpk, ok := s.odb.GetDevicePKForPeerID(e.Peer); ok {
					gkey := hex.EncodeToString(dpk.Group.PublicKey)
					s.peerStatusManager.AssociatePeer(gkey, e.Peer)
				}
			}
		}
	}()

	// get status of peers in the peerstore
	peers := s.host.Peerstore().Peers()
	for _, peer := range peers {
		// if we got some connected peer check their status
		if s.host.Network().Connectedness(peer) == network.Connected {
			s.peerStatusManager.UpdateState(peer, ConnectednessTypeConnected)
		}

		// if we already have some head exchange with this peer, associate it
		if dpk, ok := s.odb.GetDevicePKForPeerID(peer); ok {
			gkey := hex.EncodeToString(dpk.Group.PublicKey)
			s.peerStatusManager.AssociatePeer(gkey, peer)
		}
	}
}

// Status contains results of status checks
type Status struct {
	DB       error
	Protocol error
}

func (s *service) Status() Status {
	return Status{
		Protocol: nil,
	}
}
