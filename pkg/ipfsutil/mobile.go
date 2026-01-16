package ipfsutil

import (
	"context"
	"fmt"

	ipfs_config "github.com/ipfs/kubo/config"
	ipfs_p2p "github.com/ipfs/kubo/core/node/libp2p"
	p2p "github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	host "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	p2p_routing "github.com/libp2p/go-libp2p/core/routing"
	quict "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/quicreuse"
	tcpt "github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"go.uber.org/zap"

	ipfs_mobile "berty.tech/weshnet/v2/pkg/ipfsutil/mobile"
)

type DHTNetworkMode int

const (
	DHTNetworkLan DHTNetworkMode = iota
	DHTNetworkWan
	DHTNetworkDual
)

type Config func(cfg *ipfs_config.Config) ([]p2p.Option, error)

type MobileOptions struct {
	Logger          *zap.Logger
	IpfsConfigPatch Config
	// P2PStaticRelays and PeerStorePeers are only used if IpfsConfigPatch is nil
	P2PStaticRelays []string
	PeerStorePeers  []string

	HostOption    ipfs_p2p.HostOption
	RoutingOption ipfs_p2p.RoutingOption

	HostConfigFunc    ipfs_mobile.HostConfigFunc
	RoutingConfigFunc ipfs_mobile.RoutingConfigFunc

	ExtraOpts map[string]bool
}

func (o *MobileOptions) fillDefault() {
	if o.Logger == nil {
		o.Logger = zap.NewNop()
	}

	if o.HostOption == nil {
		o.HostOption = ipfs_p2p.DefaultHostOption
	}

	if o.RoutingOption == nil {
		o.RoutingOption = CustomRoutingOption(dht.ModeClient, DHTNetworkDual, dht.Concurrency(2))
	}

	if o.IpfsConfigPatch == nil {
		o.IpfsConfigPatch = o.defaultIpfsConfigPatch

		// P2PStaticRelays and PeerStorePeers are only used by defaultIpfsConfigPatch
		if o.P2PStaticRelays == nil {
			o.P2PStaticRelays = []string{DefaultP2PStaticRelay}
		}
		if o.PeerStorePeers == nil {
			o.PeerStorePeers = []string{DefaultP2PRdvpMaddr}
		}
	}

	// apply default extras
	if o.ExtraOpts == nil {
		o.ExtraOpts = make(map[string]bool)
	}

	//  if not set, disable pubsub by default to avoid collision
	if _, ok := o.ExtraOpts["pubsub"]; !ok {
		o.ExtraOpts["pubsub"] = false
	}
}

func NewIPFSMobile(ctx context.Context, repo *ipfs_mobile.RepoMobile, opts *MobileOptions) (*ipfs_mobile.IpfsMobile, error) {
	opts.fillDefault()

	var p2popts []p2p.Option

	err := repo.ApplyPatchs(func(cfg *ipfs_config.Config) error {
		var err error
		if p2popts, err = opts.IpfsConfigPatch(cfg); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// check that p2p opt is set
	if p2popts == nil {
		return nil, fmt.Errorf("unable p2p option: cannot be nil")
	}

	// configure host
	hostconfig := &ipfs_mobile.HostConfig{
		// called after host init
		ConfigFunc: opts.HostConfigFunc,

		// p2p options
		Options: p2popts,
	}

	// configure routing
	routingconfig := &ipfs_mobile.RoutingConfig{
		// called after host init
		ConfigFunc: opts.RoutingConfigFunc,
	}

	// configure ipfs mobile
	ipfsconfig := ipfs_mobile.IpfsConfig{
		HostConfig:    hostconfig,
		RoutingConfig: routingconfig,
		RepoMobile:    repo,
		ExtraOpts:     opts.ExtraOpts,
		HostOption:    opts.HostOption,
		RoutingOption: opts.RoutingOption,
	}

	return ipfs_mobile.NewNode(ctx, &ipfsconfig)
}

func CustomRoutingOption(mode dht.ModeOpt, net DHTNetworkMode, opts ...dht.Option) func(args ipfs_p2p.RoutingOptionArgs) (p2p_routing.Routing, error) {
	return func(args ipfs_p2p.RoutingOptionArgs) (p2p_routing.Routing, error) {
		opts = append(opts,
			dht.Mode(mode),
			dht.Datastore(args.Datastore),
			dht.Validator(args.Validator),
			dht.BootstrapPeers(args.BootstrapPeers...),
		)

		return newDualDHT(args.Ctx, args.Host, net, opts...)
	}
}

func (o *MobileOptions) defaultIpfsConfigPatch(cfg *ipfs_config.Config) ([]p2p.Option, error) {
	// Imitate berty setupIPFSConfig
	// https://github.com/berty/berty/blob/5a8b9cb8524c1287ab2533a9e186ac8bde7f2b57/go/internal/initutil/ipfs.go#L474C19-L474C34
	p2popts := []p2p.Option{}

	// make sure relay is enabled
	cfg.Swarm.RelayClient.Enabled = ipfs_config.True
	cfg.Swarm.Transports.Network.Relay = ipfs_config.True

	// add static relay
	pis, err := ParseAndResolveMaddrs(context.TODO(), o.Logger, o.P2PStaticRelays)
	if err != nil {
		return nil, err
	}
	if len(pis) > 0 {
		peers := make([]peer.AddrInfo, len(pis))
		for i, p := range pis {
			peers[i] = *p
		}

		p2popts = append(p2popts, p2p.EnableAutoRelayWithStaticRelays(peers))
	}

	// prefill peerstore with known rdvp servers
	peers, err := ParseAndResolveMaddrs(context.TODO(), o.Logger, o.PeerStorePeers)
	if err != nil {
		return nil, err
	}
	for _, p := range peers {
		cfg.Peering.Peers = append(cfg.Peering.Peers, *p)
	}

	// @NOTE(gfanton): disable quic transport so we can init a custom transport
	// with reusport disabled
	cfg.Swarm.Transports.Network.QUIC = ipfs_config.False
	p2popts = append(p2popts, p2p.Transport(quict.NewTransport), p2p.QUICReuse(quicreuse.NewConnManager, quicreuse.DisableReuseport()))

	// @NOTE(gfanton): disable tcp transport so we can init a custom transport
	// with reusport disabled
	cfg.Swarm.Transports.Network.TCP = ipfs_config.False
	p2popts = append(p2popts, p2p.Transport(tcpt.NewTCPTransport,
		tcpt.DisableReuseport(),
	))

	return p2popts, nil
}

const (
	// from dual package dht
	maxPrefixCountPerCpl = 2
	maxPrefixCount       = 3
)

func newDualDHT(ctx context.Context, h host.Host, net DHTNetworkMode, options ...dht.Option) (p2p_routing.Routing, error) {
	switch net {
	case DHTNetworkWan:
		options = append(options,
			dht.QueryFilter(dht.PublicQueryFilter),
			dht.RoutingTableFilter(dht.PublicRoutingTableFilter),
			dht.RoutingTablePeerDiversityFilter(dht.NewRTPeerDiversityFilter(h, maxPrefixCountPerCpl, maxPrefixCount)),
		)

		return dht.New(ctx, h, options...)
	case DHTNetworkLan:
		options = append(options,
			dht.ProtocolExtension(dual.LanExtension),
			dht.QueryFilter(dht.PrivateQueryFilter),
			dht.RoutingTableFilter(dht.PrivateRoutingTableFilter),
		)

		return dht.New(ctx, h, options...)
	default: // dual
		return dual.New(ctx, h, dual.DHTOption(options...))
	}
}
