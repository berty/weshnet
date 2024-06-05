//go:build js

package main

import (
	"context"
	"syscall/js"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

type networkFromJS struct {
	helia js.Value
}

var _ network.Network = (*networkFromJS)(nil)

// Peerstore returns the internal peerstore
// This is useful to tell the dialer about a new address for a peer.
// Or use one of the public keys found out over the network.
func (jn *networkFromJS) Peerstore() peerstore.Peerstore {
	panic("not implemented") // TODO: Implement
}

// LocalPeer returns the local peer associated with this network
func (jn *networkFromJS) LocalPeer() peer.ID {
	panic("not implemented") // TODO: Implement
}

// DialPeer establishes a connection to a given peer
func (jn *networkFromJS) DialPeer(_ context.Context, _ peer.ID) (network.Conn, error) {
	panic("not implemented") // TODO: Implement
}

// ClosePeer closes the connection to a given peer
func (jn *networkFromJS) ClosePeer(_ peer.ID) error {
	panic("not implemented") // TODO: Implement
}

// Connectedness returns a state signaling connection capabilities
func (jn *networkFromJS) Connectedness(_ peer.ID) network.Connectedness {
	panic("not implemented") // TODO: Implement
}

// Peers returns the peers connected
func (jn *networkFromJS) Peers() []peer.ID {
	panic("not implemented") // TODO: Implement
}

// Conns returns the connections in this Network
func (jn *networkFromJS) Conns() []network.Conn {
	panic("not implemented") // TODO: Implement
}

// ConnsToPeer returns the connections in this Network for given peer.
func (jn *networkFromJS) ConnsToPeer(p peer.ID) []network.Conn {
	panic("not implemented") // TODO: Implement
}

// Notify/StopNotify register and unregister a notifiee for signals
func (jn *networkFromJS) Notify(_ network.Notifiee) {
	panic("not implemented") // TODO: Implement
}

func (jn *networkFromJS) StopNotify(_ network.Notifiee) {
	panic("not implemented") // TODO: Implement
}

func (jn *networkFromJS) Close() error {
	panic("not implemented") // TODO: Implement
}

// SetStreamHandler sets the handler for new streams opened by the
// remote side. This operation is threadsafe.
func (jn *networkFromJS) SetStreamHandler(_ network.StreamHandler) {
	panic("not implemented") // TODO: Implement
}

// NewStream returns a new stream to given peer p.
// If there is no connection to p, attempts to create one.
func (jn *networkFromJS) NewStream(_ context.Context, _ peer.ID) (network.Stream, error) {
	panic("not implemented") // TODO: Implement
}

// Listen tells the network to start listening on given multiaddrs.
func (jn *networkFromJS) Listen(_ ...ma.Multiaddr) error {
	panic("not implemented") // TODO: Implement
}

// ListenAddresses returns a list of addresses at which this network listens.
func (jn *networkFromJS) ListenAddresses() []ma.Multiaddr {
	panic("not implemented") // TODO: Implement
}

// InterfaceListenAddresses returns a list of addresses at which this network
// listens. It expands "any interface" addresses (/ip4/0.0.0.0, /ip6/::) to
// use the known local interfaces.
func (jn *networkFromJS) InterfaceListenAddresses() ([]ma.Multiaddr, error) {
	panic("not implemented") // TODO: Implement
}

// ResourceManager returns the ResourceManager associated with this network
func (jn *networkFromJS) ResourceManager() network.ResourceManager {
	panic("not implemented") // TODO: Implement
}
