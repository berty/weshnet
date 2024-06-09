//go:build js

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
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
		ret[i] = &connFromJS{conn: c, hint: "network conns"}
	}
	return ret
}

type connFromJS struct {
	conn js.Value
	hint string
}

var _ network.Conn = (*connFromJS)(nil)

func (jc *connFromJS) Close() error {
	fmt.Println("called conn close")
	panic("not implemented") // TODO: Implement
}

// LocalPeer returns our peer ID
func (jc *connFromJS) LocalPeer() peer.ID {
	fmt.Println("called conn local peer")
	panic("not implemented") // TODO: Implement
}

// RemotePeer returns the peer ID of the remote peer.
func (jc *connFromJS) RemotePeer() peer.ID {
	fmt.Println("called conn remote peer")
	p, err := peer.Decode(jc.conn.Get("remotePeer").Call("toString").String())
	if err != nil {
		panic(err)
	}
	return p
}

// RemotePublicKey returns the public key of the remote peer.
func (jc *connFromJS) RemotePublicKey() ic.PubKey {
	fmt.Println("called conn remote public key")
	panic("not implemented") // TODO: Implement
}

// ConnState returns information about the connection state.
func (jc *connFromJS) ConnState() network.ConnectionState {
	fmt.Println("called conn state")
	panic("not implemented") // TODO: Implement
}

// LocalMultiaddr returns the local Multiaddr associated
// with this connection
func (jc *connFromJS) LocalMultiaddr() ma.Multiaddr {
	fmt.Println("called conn local multiaddr")
	panic("not implemented") // TODO: Implement
}

// RemoteMultiaddr returns the remote Multiaddr associated
// with this connection
func (jc *connFromJS) RemoteMultiaddr() ma.Multiaddr {
	fmt.Println("called conn remote multiaddr")
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
	fmt.Println("called conn scope")
	panic("not implemented") // TODO: Implement
}

// ID returns an identifier that uniquely identifies this Conn within this
// host, during this run. Connection IDs may repeat across restarts.
func (jc *connFromJS) ID() string {
	fmt.Println("called conn id")
	panic("not implemented") // TODO: Implement
}

// NewStream constructs a new Stream over this conn.
func (jc *connFromJS) NewStream(ctx context.Context) (network.Stream, error) {
	fmt.Println("FIXME: ignored input context in conn.NewStream")
	s, err := await(jc.conn.Call("newStream"))
	if err != nil {
		return nil, err
	}
	return newStreamFromJS(s, jc.conn, jc.hint), nil
}

type streamFromJS struct {
	s          js.Value
	conn       js.Value
	sourceChan chan js.Value
	readBuf    []byte
	readClosed bool
	hint       string
}

func newStreamFromJS(s js.Value, conn js.Value, hint string) *streamFromJS {
	ch := make(chan js.Value)
	jstrm := &streamFromJS{s: s, conn: conn, sourceChan: ch, hint: hint}
	recv := js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("wrapper cb")
		ch <- args[0]
		return nil
	})
	end := js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("stream en cb")
		if !jstrm.readClosed {
			close(jstrm.sourceChan)
			jstrm.readClosed = true
		}
		return nil
	})
	js.Global().Get("wrapAsyncGenerator").Invoke(s.Get("source"), recv, end)
	return jstrm
}

var _ network.Stream = (*streamFromJS)(nil)

func (jstrm *streamFromJS) Read(p []byte) (n int, err error) {
	/*
		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*2000)
		defer cancel()
	*/
	ctx := context.TODO()
	//fmt.Println("FIXME: using hacked-in stream.Read for", len(p), "bytes on", jstrm.hint)
	offset := 0
	for offset < len(p) {
		//fmt.Println("offset", offset)
		if len(jstrm.readBuf) != 0 {
			tn := copy(p[offset:], jstrm.readBuf)
			jstrm.readBuf = jstrm.readBuf[tn:]
			offset += tn
			continue
		}
		select {
		case jsbs := <-jstrm.sourceChan:
			//fmt.Println("got elem from chan in select")
			if jsbs.IsUndefined() {
				//fmt.Println("stream end after chan")
				return offset, io.EOF
			}
			inLen := jsbs.Get("length").Int()
			//fmt.Println("received", inLen, "bytes")
			lbuf := make([]byte, inLen)
			_ = js.CopyBytesToGo(lbuf, jsbs.Call("subarray"))
			tn := copy(p[offset:], lbuf)
			//fmt.Println("copied", tn, "bytes")
			offset += tn
			if tn < inLen {
				//fmt.Println("saved", inLen-tn, "bytes")
				jstrm.readBuf = append(jstrm.readBuf, lbuf[tn:]...)
			}
		case <-ctx.Done():
			return 0, errors.New("read timeout")
		}
	}
	if offset != len(p) {
		return 0, errors.New("unexpected offset")
	}
	fmt.Println("red", len(p), "bytes on", jstrm.Protocol(), "hint", jstrm.hint)
	return len(p), nil
}

func (jstrm *streamFromJS) Write(p []byte) (n int, err error) {
	fmt.Println("called stream write")
	panic("not implemented") // TODO: Implement
}

func (jstrm *streamFromJS) Close() error {
	fmt.Println("called stream close")
	if !jstrm.readClosed {
		close(jstrm.sourceChan)
		jstrm.readClosed = true
	}
	_, err := await(jstrm.s.Call("close"))
	return err
}

// CloseWrite closes the stream for writing but leaves it open for
// reading.
//
// CloseWrite does not free the stream, users must still call Close or
// Reset.
func (jstrm *streamFromJS) CloseWrite() error {
	fmt.Println("called stream close write")
	_, err := await(jstrm.s.Call("closeWrite"))
	return err
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
func (jstrm *streamFromJS) CloseRead() error {
	fmt.Println("called stream close read")
	if !jstrm.readClosed {
		close(jstrm.sourceChan)
		jstrm.readClosed = true
	}
	_, err := await(jstrm.s.Call("closeRead"))
	return err
}

// Reset closes both ends of the stream. Use this to tell the remote
// side to hang up and go away.
func (jstrm *streamFromJS) Reset() error {
	fmt.Println("called stream reset")

	jstrm.s.Call("abort", js.Global().Get("Error").New("go away"))
	return nil
}

func (jstrm *streamFromJS) SetDeadline(_ time.Time) error {
	fmt.Println("called stream set deadline")
	panic("not implemented") // TODO: Implement
}

func (jstrm *streamFromJS) SetReadDeadline(_ time.Time) error {
	fmt.Println("called stream set read deadline")
	panic("not implemented") // TODO: Implement
}

func (jstrm *streamFromJS) SetWriteDeadline(_ time.Time) error {
	fmt.Println("called stream set write deadline")
	panic("not implemented") // TODO: Implement
}

// ID returns an identifier that uniquely identifies this Stream within this
// host, during this run. Stream IDs may repeat across restarts.
func (jstrm *streamFromJS) ID() string {
	fmt.Println("called stream id")
	panic("not implemented") // TODO: Implement
}

func (jstrm *streamFromJS) Protocol() protocol.ID {
	ret := jstrm.s.Get("protocol")
	if ret.IsUndefined() {
		return ""
	}
	return protocol.ID(ret.String())
}

func (jstrm *streamFromJS) SetProtocol(id protocol.ID) error {
	fmt.Println("called stream set protocol")
	panic("not implemented") // TODO: Implement
}

// Stat returns metadata pertaining to this stream.
func (jstrm *streamFromJS) Stat() network.Stats {
	fmt.Println("called stream stats")
	panic("not implemented") // TODO: Implement
}

// Conn returns the connection this stream is part of.
func (jstrm *streamFromJS) Conn() network.Conn {
	fmt.Println("called stream conn")
	return &connFromJS{conn: jstrm.conn, hint: jstrm.hint}
}

// Scope returns the user's view of this stream's resource scope
func (jstrm *streamFromJS) Scope() network.StreamScope {
	fmt.Println("called stream scope")
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
		ret[i] = &connFromJS{conn: c, hint: "peer " + p.String()}
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
