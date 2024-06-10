//go:build js

package wasm

import (
	"context"
	"fmt"
	"sync"
	"syscall/js"

	"berty.tech/weshnet/pkg/ipfsutil"
	ipld "github.com/ipfs/go-ipld-format"
	ipfs_interface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
)

type coreAPIFromJS struct {
	helia         js.Value
	nextHandlerID int
	handlerMutex  sync.Mutex
}

var _ ipfsutil.ExtendedCoreAPI = (*coreAPIFromJS)(nil)

func NewCoreAPIFromJS(helia js.Value) (ipfsutil.ExtendedCoreAPI, error) {
	if helia.Type() != js.TypeObject || helia.Get("libp2p").Type() != js.TypeObject {
		return nil, errors.New("argument does not look like an helia object")
	}
	return &coreAPIFromJS{helia: helia}, nil
}

// IPFS API

// Unixfs returns an implementation of Unixfs API
func (jca *coreAPIFromJS) Unixfs() ipfs_interface.UnixfsAPI {
	panic("not implemented")
}

// Block returns an implementation of Block API
func (jca *coreAPIFromJS) Block() ipfs_interface.BlockAPI {
	panic("not implemented")
}

// Dag returns an implementation of Dag API
func (jca *coreAPIFromJS) Dag() ipfs_interface.APIDagService {
	return newDagFromJS(jca.helia)
}

// Name returns an implementation of Name API
func (jca *coreAPIFromJS) Name() ipfs_interface.NameAPI {
	panic("not implemented")
}

// Key returns an implementation of Key API
func (jca *coreAPIFromJS) Key() ipfs_interface.KeyAPI {
	return &keyAPIFromJS{helia: jca.helia}
}

// Pin returns an implementation of Pin API
func (jca *coreAPIFromJS) Pin() ipfs_interface.PinAPI {
	panic("not implemented")
}

// Object returns an implementation of Object API
func (jca *coreAPIFromJS) Object() ipfs_interface.ObjectAPI {
	panic("not implemented")
}

// Dht returns an implementation of Dht API
func (jca *coreAPIFromJS) Dht() ipfs_interface.DhtAPI {
	panic("not implemented")
}

// Swarm returns an implementation of Swarm API
func (jca *coreAPIFromJS) Swarm() ipfs_interface.SwarmAPI {
	return &swarmAPIFromJS{helia: jca.helia}
}

// PubSub returns an implementation of PubSub API
func (jca *coreAPIFromJS) PubSub() ipfs_interface.PubSubAPI {
	return &pubSubFromJS{helia: jca.helia}
}

// Routing returns an implementation of Routing API
func (jca *coreAPIFromJS) Routing() ipfs_interface.RoutingAPI {
	panic("not implemented")
}

// ResolvePath resolves the path using Unixfs resolver
func (jca *coreAPIFromJS) ResolvePath(context.Context, path.Path) (path.Resolved, error) {
	panic("not implemented")
}

// ResolveNode resolves the path (if not resolved already) using Unixfs
// resolver, gets and returns the resolved Node
func (jca *coreAPIFromJS) ResolveNode(context.Context, path.Path) (ipld.Node, error) {
	panic("not implemented")
}

// WithOptions creates new instance of CoreAPI based on this instance with
// a set of options applied
func (jca *coreAPIFromJS) WithOptions(...options.ApiOption) (ipfs_interface.CoreAPI, error) {
	panic("not implemented")
}

// HOST API

// ID returns the (local) peer.ID associated with this Host
func (jca *coreAPIFromJS) ID() peer.ID {
	peerId := jca.helia.Get("libp2p").Get("peerId").Call("toString").String()
	obj, err := peer.Decode(peerId)
	if err != nil {
		panic(err)
	}
	return obj
}

// Peerstore returns the Host's repository of Peer Addresses and Keys.
func (jca *coreAPIFromJS) Peerstore() peerstore.Peerstore {
	return &peerStoreFromJS{helia: jca.helia}
}

// Returns the listen addresses of the Host
func (jca *coreAPIFromJS) Addrs() []multiaddr.Multiaddr {
	panic("not implemented")
}

// Networks returns the Network interface of the Host
func (jca *coreAPIFromJS) Network() network.Network {
	return &networkFromJS{helia: jca.helia}
}

// Mux returns the Mux multiplexing incoming streams to protocol handlers
func (jca *coreAPIFromJS) Mux() protocol.Switch {
	panic("not implemented")
}

// Connect ensures there is a connection between this host and the peer with
// given peer.ID. Connect will absorb the addresses in pi into its internal
// peerstore. If there is not an active connection, Connect will issue a
// h.Network.Dial, and block until a connection is open, or an error is
// returned. // TODO: Relay + NAT.
func (jca *coreAPIFromJS) Connect(ctx context.Context, pi peer.AddrInfo) error {
	panic("not implemented")
}

// SetStreamHandler sets the protocol handler on the Host's Mux.
// This is equivalent to:
//
//	host.Mux().SetHandler(proto, handler)
//
// (Threadsafe)
func (jca *coreAPIFromJS) SetStreamHandler(pid protocol.ID, handler network.StreamHandler) {
	_, err := await(jca.helia.Get("libp2p").Call("handle", string(pid), js.FuncOf(func(this js.Value, args []js.Value) any {
		return Promisify(func() ([]any, error) {
			jca.handlerMutex.Lock()
			defer jca.handlerMutex.Unlock()
			id := jca.nextHandlerID
			jca.nextHandlerID += 1
			if len(args) == 0 {
				return nil, errors.New("empty stream arg")
			}
			handler(newStreamFromJS(
				args[0].Get("stream"),
				args[0].Get("connection"),
				fmt.Sprintf("handler %s (%d)", string(pid), id),
			))
			return nil, nil
		})
	})))
	if err != nil {
		panic(err)
	}
}

// SetStreamHandlerMatch sets the protocol handler on the Host's Mux
// using a matching function for protocol selection.
func (jca *coreAPIFromJS) SetStreamHandlerMatch(protocol.ID, func(protocol.ID) bool, network.StreamHandler) {
	panic("not implemented")
}

// RemoveStreamHandler removes a handler on the mux that was set by
// SetStreamHandler
func (jca *coreAPIFromJS) RemoveStreamHandler(pid protocol.ID) {
	panic("not implemented")
}

func sliceMap[I any, O any](slice []I, transform func(I) O) []O {
	ret := make([]O, len(slice))
	for i, in := range slice {
		ret[i] = transform(in)
	}
	return ret
}

// NewStream opens a new stream to given peer p, and writes a p2p/protocol
// header with given ProtocolID. If there is no connection to p, attempts
// to create one. If ProtocolID is "", writes no header.
// (Threadsafe)
func (jca *coreAPIFromJS) NewStream(ctx context.Context, p peer.ID, pids ...protocol.ID) (network.Stream, error) {
	fmt.Println("dialing", pids)
	conn, err := await(jca.helia.Get("libp2p").Call("dial", jsPeerID(p)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial peer")
	}
	strm, err := await(conn.Call("newStream", JSArrayTransform(pids, func(pid protocol.ID) string { return string(pid) })))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open stream")
	}
	return newStreamFromJS(strm, conn, "new stream "+p.String()+" "+fmt.Sprint(pids)), nil
}

// Close shuts down the host, its Network, and services.
func (jca *coreAPIFromJS) Close() error {
	panic("not implemented")
}

// ConnManager returns this hosts connection manager
func (jca *coreAPIFromJS) ConnManager() connmgr.ConnManager {
	fmt.Println("FIXME: providing nil ConnManager")
	return nil
}

// EventBus returns the hosts eventbus
func (jca *coreAPIFromJS) EventBus() event.Bus {
	return &eventBusFromJS{helia: jca.helia}
}

// CONNMGR

func (jca *coreAPIFromJS) ConnMgr() ipfsutil.ConnMgr {
	return &weshConnMgrFromJS{}
}
