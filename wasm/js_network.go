//go:build js

package main

import (
	"context"
	"fmt"
	"syscall/js"
	"time"

	ic "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
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
func (jn *networkFromJS) Connectedness(p peer.ID) network.Connectedness {
	conns := jn.ConnsToPeer(p)
	if len(conns) == 0 {
		return network.NotConnected
	}
	return network.Connected
}

// Peers returns the peers connected
func (jn *networkFromJS) Peers() []peer.ID {
	ret, err := heliaConnectedPeers(jn.helia)
	if err != nil {
		panic(err)
	}
	return ret
}

// Conns returns the connections in this Network
func (jn *networkFromJS) Conns() []network.Conn {
	conns := jn.helia.Get("libp2p").Call("getConnections")
	l := conns.Length()
	ret := make([]network.Conn, l)
	for i := 0; i < l; i++ {
		c := conns.Index(i)
		ret[i] = &connFromJS{conn: c}
	}
	return ret
}

type connFromJS struct {
	conn js.Value
}

var _ network.Conn = (*connFromJS)(nil)

func (jc *connFromJS) Close() error {
	panic("not implemented") // TODO: Implement
}

// LocalPeer returns our peer ID
func (jc *connFromJS) LocalPeer() peer.ID {
	panic("not implemented") // TODO: Implement
}

// RemotePeer returns the peer ID of the remote peer.
func (jc *connFromJS) RemotePeer() peer.ID {
	panic("not implemented") // TODO: Implement
}

// RemotePublicKey returns the public key of the remote peer.
func (jc *connFromJS) RemotePublicKey() ic.PubKey {
	panic("not implemented") // TODO: Implement
}

// ConnState returns information about the connection state.
func (jc *connFromJS) ConnState() network.ConnectionState {
	panic("not implemented") // TODO: Implement
}

// LocalMultiaddr returns the local Multiaddr associated
// with this connection
func (jc *connFromJS) LocalMultiaddr() ma.Multiaddr {
	panic("not implemented") // TODO: Implement
}

// RemoteMultiaddr returns the remote Multiaddr associated
// with this connection
func (jc *connFromJS) RemoteMultiaddr() ma.Multiaddr {
	maddrStr := jc.conn.Get("remoteAddr").Call("toString").String()
	maddr := ma.StringCast(maddrStr)
	return maddr
}

// Stat stores metadata pertaining to this conn.
func (jc *connFromJS) Stat() network.ConnStats {
	fmt.Println("FIFXME: partial conn Stat")
	consoleLog("stat", jc.conn)
	return network.ConnStats{
		Stats: network.Stats{
			Direction: directionFromJS(jc.conn.Get("direction")),
			Transient: jc.conn.Get("transient").Bool(),
			// Opened:    TODO,
			// Extra:     TODO,
		},
		NumStreams: jc.conn.Get("streams").Length(),
	}
}

func directionFromJS(val js.Value) network.Direction {
	switch val.String() {
	case "outbound":
		return network.DirOutbound
	case "inbound":
		return network.DirInbound
	default:
		return network.DirUnknown
	}
}

// Scope returns the user view of this connection's resource scope
func (jc *connFromJS) Scope() network.ConnScope {
	panic("not implemented") // TODO: Implement
}

// ID returns an identifier that uniquely identifies this Conn within this
// host, during this run. Connection IDs may repeat across restarts.
func (jc *connFromJS) ID() string {
	panic("not implemented") // TODO: Implement
}

// NewStream constructs a new Stream over this conn.
func (jc *connFromJS) NewStream(ctx context.Context) (network.Stream, error) {
	fmt.Println("FIXME: ignored input context in conn.NewStream")
	s, err := await(jc.conn.Call("newStream"))
	if err != nil {
		return nil, err
	}
	return &jsStream{s: s}, nil
}

type jsStream struct {
	s js.Value
}

var _ network.Stream = (*jsStream)(nil)

func (jstrm *jsStream) Read(p []byte) (n int, err error) {
	panic("not implemented") // TODO: Implement
}

func (jstrm *jsStream) Write(p []byte) (n int, err error) {
	panic("not implemented") // TODO: Implement
}

func (jstrm *jsStream) Close() error {
	panic("not implemented") // TODO: Implement
}

// CloseWrite closes the stream for writing but leaves it open for
// reading.
//
// CloseWrite does not free the stream, users must still call Close or
// Reset.
func (jstrm *jsStream) CloseWrite() error {
	panic("not implemented") // TODO: Implement
}

// CloseRead closes the stream for reading but leaves it open for
// writing.
//
// When CloseRead is called, all in-progress Read calls are interrupted with a non-EOF error and
// no further calls to Read will succeed.
//
// The handling of new incoming data on the stream after calling this function is implementation defined.
//
// CloseRead does not free the stream, users must still call Close or
// Reset.
func (jstrm *jsStream) CloseRead() error {
	panic("not implemented") // TODO: Implement
}

// Reset closes both ends of the stream. Use this to tell the remote
// side to hang up and go away.
func (jstrm *jsStream) Reset() error {
	panic("not implemented") // TODO: Implement
}

func (jstrm *jsStream) SetDeadline(_ time.Time) error {
	panic("not implemented") // TODO: Implement
}

func (jstrm *jsStream) SetReadDeadline(_ time.Time) error {
	panic("not implemented") // TODO: Implement
}

func (jstrm *jsStream) SetWriteDeadline(_ time.Time) error {
	panic("not implemented") // TODO: Implement
}

// ID returns an identifier that uniquely identifies this Stream within this
// host, during this run. Stream IDs may repeat across restarts.
func (jstrm *jsStream) ID() string {
	panic("not implemented") // TODO: Implement
}

func (jstrm *jsStream) Protocol() protocol.ID {
	panic("not implemented") // TODO: Implement
}

func (jstrm *jsStream) SetProtocol(id protocol.ID) error {
	panic("not implemented") // TODO: Implement
}

// Stat returns metadata pertaining to this stream.
func (jstrm *jsStream) Stat() network.Stats {
	panic("not implemented") // TODO: Implement
}

// Conn returns the connection this stream is part of.
func (jstrm *jsStream) Conn() network.Conn {
	panic("not implemented") // TODO: Implement
}

// Scope returns the user's view of this stream's resource scope
func (jstrm *jsStream) Scope() network.StreamScope {
	panic("not implemented") // TODO: Implement
}

// GetStreams returns all open streams over this conn.
func (jc *connFromJS) GetStreams() []network.Stream {
	panic("not implemented") // TODO: Implement
}

// IsClosed returns whether a connection is fully closed, so it can
// be garbage collected.
func (jc *connFromJS) IsClosed() bool {
	panic("not implemented") // TODO: Implement
}

// ConnsToPeer returns the connections in this Network for given peer.
func (jn *networkFromJS) ConnsToPeer(p peer.ID) []network.Conn {
	conns := jn.helia.Get("libp2p").Call("getConnections", p.String())
	l := conns.Length()
	ret := make([]network.Conn, l)
	for i := 0; i < l; i++ {
		c := conns.Index(i)
		ret[i] = &connFromJS{conn: c}
	}
	return ret
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
	ret, err := heliaListenAddresses(jn.helia)
	if err != nil {
		panic(err)
	}
	return ret
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
