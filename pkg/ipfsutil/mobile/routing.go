package node

import (
	"fmt"

	ipfs_p2p "github.com/ipfs/kubo/core/node/libp2p"
	p2p_host "github.com/libp2p/go-libp2p/core/host"
	p2p_routing "github.com/libp2p/go-libp2p/core/routing"
)

type RoutingConfigFunc func(p2p_host.Host, p2p_routing.Routing) error

type RoutingConfig struct {
	ConfigFunc RoutingConfigFunc
}

func NewRoutingConfigOption(ro ipfs_p2p.RoutingOption, rc *RoutingConfig) ipfs_p2p.RoutingOption {
	return func(args ipfs_p2p.RoutingOptionArgs) (p2p_routing.Routing, error) {
		routing, err := ro(args)
		if err != nil {
			return nil, err
		}

		if rc.ConfigFunc != nil {
			if err := rc.ConfigFunc(args.Host, routing); err != nil {
				return nil, fmt.Errorf("failed to config routing: %w", err)
			}
		}

		return routing, nil
	}
}
