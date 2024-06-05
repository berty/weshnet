//go:build js

package main

import (
	"context"
	"fmt"
	"syscall/js"

	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

type swarmAPIFromJS struct {
	helia js.Value
}

var _ iface.SwarmAPI = (*swarmAPIFromJS)(nil)

// Connect to a given peer
func (jswarm *swarmAPIFromJS) Connect(_ context.Context, _ peer.AddrInfo) error {
	panic("not implemented") // TODO: Implement
}

// Disconnect from a given address
func (jswarm *swarmAPIFromJS) Disconnect(_ context.Context, _ ma.Multiaddr) error {
	panic("not implemented") // TODO: Implement
}

// Peers returns the list of peers we are connected to
func (jswarm *swarmAPIFromJS) Peers(_ context.Context) ([]iface.ConnectionInfo, error) {
	fmt.Println("FIXME: mocked swarm.Peers")
	return nil, nil
}

// KnownAddrs returns the list of all addresses this node is aware of
func (jswarm *swarmAPIFromJS) KnownAddrs(_ context.Context) (map[peer.ID][]ma.Multiaddr, error) {
	panic("not implemented") // TODO: Implement
}

// LocalAddrs returns the list of announced listening addresses
func (jswarm *swarmAPIFromJS) LocalAddrs(_ context.Context) ([]ma.Multiaddr, error) {
	panic("not implemented") // TODO: Implement
}

// ListenAddrs returns the list of all listening addresses
func (jswarm *swarmAPIFromJS) ListenAddrs(_ context.Context) ([]ma.Multiaddr, error) {
	return heliaListenAddresses(jswarm.helia)
}
