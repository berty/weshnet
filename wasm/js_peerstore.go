//go:build js

package main

import (
	"context"
	"syscall/js"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	ic "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/crypto/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
)

type peerStoreFromJS struct {
	helia js.Value
}

var _ peerstore.Peerstore = (*peerStoreFromJS)(nil)

func (jps *peerStoreFromJS) Close() error {
	panic("not implemented") // TODO: Implement
}

// AddAddr calls AddAddrs(p, []ma.Multiaddr{addr}, ttl)
func (jps *peerStoreFromJS) AddAddr(p peer.ID, addr ma.Multiaddr, ttl time.Duration) {
	panic("not implemented") // TODO: Implement
}

// AddAddrs gives this AddrBook addresses to use, with a given ttl
// (time-to-live), after which the address is no longer valid.
// If the manager has a longer TTL, the operation is a no-op for that address
func (jps *peerStoreFromJS) AddAddrs(p peer.ID, addrs []ma.Multiaddr, ttl time.Duration) {
	panic("not implemented") // TODO: Implement
}

// SetAddr calls mgr.SetAddrs(p, addr, ttl)
func (jps *peerStoreFromJS) SetAddr(p peer.ID, addr ma.Multiaddr, ttl time.Duration) {
	panic("not implemented") // TODO: Implement
}

// SetAddrs sets the ttl on addresses. This clears any TTL there previously.
// This is used when we receive the best estimate of the validity of an address.
func (jps *peerStoreFromJS) SetAddrs(p peer.ID, addrs []ma.Multiaddr, ttl time.Duration) {
	panic("not implemented") // TODO: Implement
}

// UpdateAddrs updates the addresses associated with the given peer that have
// the given oldTTL to have the given newTTL.
func (jps *peerStoreFromJS) UpdateAddrs(p peer.ID, oldTTL time.Duration, newTTL time.Duration) {
	panic("not implemented") // TODO: Implement
}

// Addrs returns all known (and valid) addresses for a given peer.
func (jps *peerStoreFromJS) Addrs(p peer.ID) []ma.Multiaddr {
	panic("not implemented") // TODO: Implement
}

// AddrStream returns a channel that gets all addresses for a given
// peer sent on it. If new addresses are added after the call is made
// they will be sent along through the channel as well.
func (jps *peerStoreFromJS) AddrStream(_ context.Context, _ peer.ID) <-chan ma.Multiaddr {
	panic("not implemented") // TODO: Implement
}

// ClearAddresses removes all previously stored addresses.
func (jps *peerStoreFromJS) ClearAddrs(p peer.ID) {
	panic("not implemented") // TODO: Implement
}

// PeersWithAddrs returns all of the peer IDs stored in the AddrBook.
func (jps *peerStoreFromJS) PeersWithAddrs() peer.IDSlice {
	panic("not implemented") // TODO: Implement
}

// PubKey stores the public key of a peer.
func (jps *peerStoreFromJS) PubKey(_ peer.ID) ic.PubKey {
	panic("not implemented") // TODO: Implement
}

// AddPubKey stores the public key of a peer.
func (jps *peerStoreFromJS) AddPubKey(_ peer.ID, _ ic.PubKey) error {
	panic("not implemented") // TODO: Implement
}

// PrivKey returns the private key of a peer, if known. Generally this might only be our own
// private key, see
// https://discuss.libp2p.io/t/what-is-the-purpose-of-having-map-peer-id-privatekey-in-peerstore/74.
func (jps *peerStoreFromJS) PrivKey(id peer.ID) ic.PrivKey {
	keychain := jps.helia.Get("libp2p").Get("services").Get("keychain")
	key, err := await(keychain.Call("findKeyById", id.String()))
	if err != nil {
		panic(err)
	}
	return &privKeyFromJS{key: key}
}

type privKeyFromJS struct {
	key js.Value
}

// Equals checks whether two PubKeys are the same
func (jpk *privKeyFromJS) Equals(_ crypto.Key) bool {
	panic("not implemented") // TODO: Implement
}

// Raw returns the raw bytes of the key (not wrapped in the
// libp2p-crypto protobuf).
//
// This function is the inverse of {Priv,Pub}KeyUnmarshaler.
func (jpk *privKeyFromJS) Raw() ([]byte, error) {
	panic("not implemented") // TODO: Implement
}

// Type returns the protobuf key type.
func (jpk *privKeyFromJS) Type() pb.KeyType {
	panic("not implemented") // TODO: Implement
}

// Cryptographically sign the given bytes
func (jpk *privKeyFromJS) Sign(_ []byte) ([]byte, error) {
	panic("not implemented") // TODO: Implement
}

// Return a public key paired with this private key
func (jpk *privKeyFromJS) GetPublic() crypto.PubKey {
	panic("not implemented") // TODO: Implement
}

// AddPrivKey stores the private key of a peer.
func (jps *peerStoreFromJS) AddPrivKey(_ peer.ID, _ ic.PrivKey) error {
	panic("not implemented") // TODO: Implement
}

// PeersWithKeys returns all the peer IDs stored in the KeyBook.
func (jps *peerStoreFromJS) PeersWithKeys() peer.IDSlice {
	panic("not implemented") // TODO: Implement
}

// RemovePeer removes all keys associated with a peer.
func (jps *peerStoreFromJS) RemovePeer(_ peer.ID) {
	panic("not implemented") // TODO: Implement
}

// Get / Put is a simple registry for other peer-related key/value pairs.
// If we find something we use often, it should become its own set of
// methods. This is a last resort.
func (jps *peerStoreFromJS) Get(p peer.ID, key string) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}

func (jps *peerStoreFromJS) Put(p peer.ID, key string, val interface{}) error {
	panic("not implemented") // TODO: Implement
}

// RecordLatency records a new latency measurement
func (jps *peerStoreFromJS) RecordLatency(_ peer.ID, _ time.Duration) {
	panic("not implemented") // TODO: Implement
}

// LatencyEWMA returns an exponentially-weighted moving avg.
// of all measurements of a peer's latency.
func (jps *peerStoreFromJS) LatencyEWMA(_ peer.ID) time.Duration {
	panic("not implemented") // TODO: Implement
}

func (jps *peerStoreFromJS) GetProtocols(_ peer.ID) ([]protocol.ID, error) {
	panic("not implemented") // TODO: Implement
}

func (jps *peerStoreFromJS) AddProtocols(_ peer.ID, _ ...protocol.ID) error {
	panic("not implemented") // TODO: Implement
}

func (jps *peerStoreFromJS) SetProtocols(_ peer.ID, _ ...protocol.ID) error {
	panic("not implemented") // TODO: Implement
}

func (jps *peerStoreFromJS) RemoveProtocols(_ peer.ID, _ ...protocol.ID) error {
	panic("not implemented") // TODO: Implement
}

// SupportsProtocols returns the set of protocols the peer supports from among the given protocols.
// If the returned error is not nil, the result is indeterminate.
func (jps *peerStoreFromJS) SupportsProtocols(_ peer.ID, _ ...protocol.ID) ([]protocol.ID, error) {
	panic("not implemented") // TODO: Implement
}

// FirstSupportedProtocol returns the first protocol that the peer supports among the given protocols.
// If the peer does not support any of the given protocols, this function will return an empty protocol.ID and a nil error.
// If the returned error is not nil, the result is indeterminate.
func (jps *peerStoreFromJS) FirstSupportedProtocol(_ peer.ID, _ ...protocol.ID) (protocol.ID, error) {
	panic("not implemented") // TODO: Implement
}

// PeerInfo returns a peer.PeerInfo struct for given peer.ID.
// This is a small slice of the information Peerstore has on
// that peer, useful to other services.
func (jps *peerStoreFromJS) PeerInfo(_ peer.ID) peer.AddrInfo {
	panic("not implemented") // TODO: Implement
}

// Peers returns all of the peer IDs stored across all inner stores.
func (jps *peerStoreFromJS) Peers() peer.IDSlice {
	panic("not implemented") // TODO: Implement
}
