package tinder

import (
	"context"

	p2p_routing "github.com/libp2p/go-libp2p-core/routing"
	discovery "github.com/libp2p/go-libp2p-discovery"
	"go.uber.org/zap"
)

type DriverRouting interface {
	p2p_routing.Routing

	Driver
}

type routing struct {
	p2p_routing.ContentRouting
	p2p_routing.PeerRouting
	p2p_routing.ValueStore

	Driver

	bootstrap func(context.Context) error
}

func (r *routing) Bootstrap(ctx context.Context) error {
	return r.bootstrap(ctx)
}

func NewDriverRouting(logger *zap.Logger, name string, r p2p_routing.Routing) DriverRouting {
	rdisc := discovery.NewRoutingDiscovery(r)
	md := ComposeDriver(name, rdisc, rdisc, NoopUnregisterer)
	return &routing{
		ContentRouting: r,
		ValueStore:     r,
		PeerRouting:    r,
		Driver:         md,

		bootstrap: r.Bootstrap,
	}
}
