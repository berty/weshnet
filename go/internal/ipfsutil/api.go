package ipfsutil

import (
	"context"

	"berty.tech/berty/v2/go/pkg/errcode"
	ds "github.com/ipfs/go-datastore"
	ipfs_ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	ipfs_core "github.com/ipfs/go-ipfs/core"
	ipfs_coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	ipfs_node "github.com/ipfs/go-ipfs/core/node"
	ipfs_libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	ipfs_repo "github.com/ipfs/go-ipfs/repo"
	ipfs_interface "github.com/ipfs/interface-go-ipfs-core"

	p2p "github.com/libp2p/go-libp2p" // nolint:staticcheck
	p2p_host "github.com/libp2p/go-libp2p-core/host"
	p2p_peer "github.com/libp2p/go-libp2p-core/peer" // nolint:staticcheck
	p2p_ps "github.com/libp2p/go-libp2p-core/peerstore"
	// nolint:staticcheck
)

type CoreAPIOption func(context.Context, *ipfs_core.IpfsNode, ipfs_interface.CoreAPI) error

type CoreAPIConfig struct {
	BootstrapAddrs []string
	SwarmAddrs     []string

	ExtraLibp2pOption p2p.Option
	Routing           ipfs_libp2p.RoutingOption

	Options []CoreAPIOption
}

func NewCoreAPI(ctx context.Context, cfg *CoreAPIConfig) (ExtendedCoreAPI, *ipfs_core.IpfsNode, error) {
	ds := dsync.MutexWrap(ds.NewMapDatastore())
	return NewCoreAPIFromDatastore(ctx, ds, cfg)
}

func NewCoreAPIFromDatastore(ctx context.Context, ds ipfs_ds.Batching, cfg *CoreAPIConfig) (ExtendedCoreAPI, *ipfs_core.IpfsNode, error) {
	repo, err := CreateMockedRepo(ds)
	if err != nil {
		return nil, nil, errcode.TODO.Wrap(err)
	}

	return NewCoreAPIFromRepo(ctx, repo, cfg)
}

func NewCoreAPIFromRepo(ctx context.Context, repo ipfs_repo.Repo, cfg *CoreAPIConfig) (ExtendedCoreAPI, *ipfs_core.IpfsNode, error) {
	bcfg, err := CreateBuildConfig(repo, cfg)
	if err != nil {
		return nil, nil, errcode.TODO.Wrap(err)
	}

	if err := updateRepoConfig(repo, cfg); err != nil {
		return nil, nil, errcode.TODO.Wrap(err)
	}

	if cfg.Options == nil {
		cfg.Options = []CoreAPIOption{}
	}

	return NewConfigurableCoreAPI(ctx, bcfg, cfg.Options...)
}

// NewConfigurableCoreAPI returns an IPFS CoreAPI from a provided ipfs_node.BuildCfg
func NewConfigurableCoreAPI(ctx context.Context, bcfg *ipfs_node.BuildCfg, opts ...CoreAPIOption) (ExtendedCoreAPI, *ipfs_core.IpfsNode, error) {
	node, err := ipfs_core.NewNode(ctx, bcfg)
	if err != nil {
		return nil, nil, errcode.TODO.Wrap(err)
	}

	api, err := ipfs_coreapi.NewCoreAPI(node)
	if err != nil {
		node.Close()
		return nil, nil, errcode.TODO.Wrap(err)
	}

	for _, opt := range opts {
		err := opt(ctx, node, api)
		if err != nil {
			node.Close()
			return nil, nil, err
		}
	}

	return NewExtendedCoreAPI(node.PeerHost, api), node, nil
}

func CreateBuildConfig(repo ipfs_repo.Repo, opts *CoreAPIConfig) (*ipfs_node.BuildCfg, error) {
	if opts == nil {
		opts = &CoreAPIConfig{}
	}

	routingOpt := ipfs_libp2p.DHTOption
	if opts.Routing != nil {
		routingOpt = opts.Routing
	}

	hostOpt := ipfs_libp2p.DefaultHostOption
	if opts.ExtraLibp2pOption != nil {
		hostOpt = wrapP2POptionsToHost(hostOpt, opts.ExtraLibp2pOption)
	}

	return &ipfs_node.BuildCfg{
		Online:                      true,
		Permanent:                   true,
		DisableEncryptedConnections: false,
		NilRepo:                     false,
		Routing:                     routingOpt,
		Host:                        hostOpt,
		Repo:                        repo,
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	}, nil
}

func updateRepoConfig(repo ipfs_repo.Repo, cfg *CoreAPIConfig) error {
	rcfg, err := repo.Config()
	if err != nil {
		return err
	}

	if cfg.BootstrapAddrs != nil {
		rcfg.Bootstrap = cfg.BootstrapAddrs
	}

	if len(cfg.SwarmAddrs) != 0 {
		rcfg.Addresses.Swarm = cfg.SwarmAddrs
	}

	return repo.SetConfig(rcfg)
}

func wrapP2POptionsToHost(hf ipfs_libp2p.HostOption, opt ...p2p.Option) ipfs_libp2p.HostOption {
	return func(ctx context.Context, id p2p_peer.ID, ps p2p_ps.Peerstore, options ...p2p.Option) (p2p_host.Host, error) {
		return hf(ctx, id, ps, append(options, opt...)...)
	}
}
