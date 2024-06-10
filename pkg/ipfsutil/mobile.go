//go:build !js

package ipfsutil

import (
	"context"
	"fmt"

	ds "github.com/ipfs/go-datastore"
	ipfs_config "github.com/ipfs/kubo/config"
	ipfs_p2p "github.com/ipfs/kubo/core/node/libp2p"
	p2p "github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	p2p_dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	p2p_record "github.com/libp2p/go-libp2p-record"
	host "github.com/libp2p/go-libp2p/core/host"
	p2p_host "github.com/libp2p/go-libp2p/core/host"
	p2p_peer "github.com/libp2p/go-libp2p/core/peer"
	p2p_routing "github.com/libp2p/go-libp2p/core/routing"
	"github.com/pkg/errors"

	ipfs_mobile "berty.tech/weshnet/pkg/ipfsutil/mobile"
)

type DHTNetworkMode int

const (
	DHTNetworkLan DHTNetworkMode = iota
	DHTNetworkWan
	DHTNetworkDual
)

type Config func(cfg *ipfs_config.Config) ([]p2p.Option, error)

type MobileOptions struct {
	IpfsConfigPatch Config

	HostOption    ipfs_p2p.HostOption
	RoutingOption ipfs_p2p.RoutingOption

	HostConfigFunc    ipfs_mobile.HostConfigFunc
	RoutingConfigFunc ipfs_mobile.RoutingConfigFunc

	ExtraOpts map[string]bool
}

func (o *MobileOptions) fillDefault() {
	if o.HostOption == nil {
		o.HostOption = ipfs_p2p.DefaultHostOption
	}

	if o.RoutingOption == nil {
		o.RoutingOption = CustomRoutingOption(p2p_dht.ModeClient, DHTNetworkDual, p2p_dht.Concurrency(2))
	}

	if o.IpfsConfigPatch == nil {
		o.IpfsConfigPatch = defaultIpfsConfigPatch
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

	relayAddrs := []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
	}

	relayPeerInfo := make([]p2p_peer.AddrInfo, 0, len(relayAddrs))
	for _, addr := range relayAddrs {
		info, err := p2p_peer.AddrInfoFromString(addr)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse relay addr info")
		}
		relayPeerInfo = append(relayPeerInfo, *info)
	}
	// configure host
	hostconfig := &ipfs_mobile.HostConfig{
		// called after host init
		ConfigFunc: opts.HostConfigFunc,

		// p2p options
		Options: append(p2popts, p2p.EnableAutoRelayWithStaticRelays(relayPeerInfo)),
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

func CustomRoutingOption(mode p2p_dht.ModeOpt, net DHTNetworkMode, opts ...p2p_dht.Option) func(
	ctx context.Context,
	host p2p_host.Host,
	dstore ds.Batching,
	validator p2p_record.Validator,
	bootstrapPeers ...p2p_peer.AddrInfo,
) (p2p_routing.Routing, error) {
	return func(
		ctx context.Context,
		host p2p_host.Host,
		dstore ds.Batching,
		validator p2p_record.Validator,
		bootstrapPeers ...p2p_peer.AddrInfo,
	) (p2p_routing.Routing, error) {
		opts = append(opts,
			p2p_dht.Mode(mode),
			p2p_dht.Datastore(dstore),
			p2p_dht.Validator(validator),
			p2p_dht.BootstrapPeers(bootstrapPeers...),
		)

		return newDualDHT(ctx, host, net, opts...)
	}
}

func defaultIpfsConfigPatch(_ *ipfs_config.Config) ([]p2p.Option, error) {
	return []p2p.Option{}, nil
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
