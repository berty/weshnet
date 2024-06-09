package weshnet

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	ds "github.com/ipfs/go-datastore"
	ds_sync "github.com/ipfs/go-datastore/sync"
	ipfs_interface "github.com/ipfs/interface-go-ipfs-core"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/host/eventbus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"berty.tech/go-orbit-db/baseorbitdb"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/weshnet/internal/bertyversion"
	"berty.tech/weshnet/pkg/bertyvcissuer"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/ipfsutil"
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
	IpfsCoreAPI() ipfs_interface.CoreAPI
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
			return errors.New("needs a datastore")
		}
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
		return nil, errcode.TODO.Wrap(err)
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
		return nil, errcode.TODO.Wrap(err)
	}

	opts.Logger.Debug("Opened account group", tyber.FormatStepLogFields(ctx, []tyber.Detail{{Name: "AccountGroup", Description: accountGroupCtx.group.String()}})...)

	var contactRequestsManager *contactRequestsManager
	var swiper *Swiper
	if opts.TinderService != nil {
		swiper = NewSwiper(opts.Logger, opts.TinderService, opts.OrbitDB.rotationInterval)
		opts.Logger.Debug("Tinder swiper is enabled", tyber.FormatStepLogFields(ctx, []tyber.Detail{})...)

		if contactRequestsManager, err = newContactRequestsManager(swiper, accountGroupCtx.metadataStore, opts.IpfsCoreAPI, opts.Logger); err != nil {
			cancel()
			return nil, errcode.TODO.Wrap(err)
		}
	} else {
		opts.Logger.Warn("No tinder driver provided, incoming and outgoing contact requests won't be enabled", tyber.FormatStepLogFields(ctx, []tyber.Detail{})...)
	}

	if err := opts.SecretStore.PutGroup(ctx, accountGroupCtx.Group()); err != nil {
		cancel()
		return nil, errcode.ErrInternal.Wrap(fmt.Errorf("unable to add account group to group datastore, err: %w", err))
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

func (s *service) IpfsCoreAPI() ipfs_interface.CoreAPI {
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
