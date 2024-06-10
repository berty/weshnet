//go:build js

package wasm

import (
	"context"
	"errors"
	"fmt"
	"syscall/js"

	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/libp2p/go-libp2p/core/peer"
)

type pubSubFromJS struct {
	helia js.Value
}

var _ iface.PubSubAPI = (*pubSubFromJS)(nil)

// Ls lists subscribed topics by name
func (jps *pubSubFromJS) Ls(_ context.Context) ([]string, error) {
	panic("Ls not implemented") // TODO: Implement
}

// Peers list peers we are currently pubsubbing with
func (jps *pubSubFromJS) Peers(ctx context.Context, opts ...options.PubSubPeersOption) ([]peer.ID, error) {
	settings := options.PubSubPeersSettings{}
	for _, opt := range opts {
		err := opt(&settings)
		if err != nil {
			return nil, fmt.Errorf("failed to apply peers options: %w", err)
		}
	}

	pubsub := jps.helia.Get("libp2p").Get("services").Get("pubsub")
	if pubsub.Type() != js.TypeObject {
		return nil, errors.New("pubsub is not an object")
	}
	var jsPeers js.Value
	if settings.Topic == "" {
		jsPeers = pubsub.Call("getPeers")
	} else {
		jsPeers = pubsub.Call("getSubscribers", settings.Topic)
	}
	return peersFromJS(jsPeers)
}

// Publish a message to a given pubsub topic
func (jps *pubSubFromJS) Publish(_ context.Context, _ string, _ []byte) error {
	panic("Publish not implemented") // TODO: Implement
}

// Subscribe to messages on a given topic
func (jps *pubSubFromJS) Subscribe(_ context.Context, topic string, _ ...options.PubSubSubscribeOption) (iface.PubSubSubscription, error) {
	pubsub := jps.helia.Get("libp2p").Get("services").Get("pubsub")
	if pubsub.Type() != js.TypeObject {
		return nil, errors.New("pubsub is not an object")
	}

	fmt.Println("FIXME: partial sub on topic", topic)

	// TODO: check if this should be buffered
	ch := make(chan js.Value)

	pubsub.Call("addEventListener", "gossipsub:message", js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("pubsub msg")
		if len(args) < 1 {
			panic("not enough args")
		}
		ch <- args[0]
		return nil
	}))

	pubsub.Call("subscribe", topic)

	return &subFromJS{ch: ch}, nil
}

type subFromJS struct {
	ch chan js.Value
}

func (jsub *subFromJS) Close() error {
	close(jsub.ch)
	return nil
}

// Next return the next incoming message
func (jsub *subFromJS) Next(_ context.Context) (iface.PubSubMessage, error) {
	msg := <-jsub.ch
	return &msgFromJS{msg: msg}, nil
}

type msgFromJS struct {
	msg js.Value
}

// From returns id of a peer from which the message has arrived
func (jmsg *msgFromJS) From() peer.ID {
	jspid := jmsg.msg.Get("detail").Get("propagationSource")
	p, err := peerFromJS(jspid)
	if err != nil {
		panic(err)
	}
	return p
}

// Data returns the message body
func (jmsg *msgFromJS) Data() []byte {
	panic("not implemented") // TODO: Implement
}

// Seq returns message identifier
func (jmsg *msgFromJS) Seq() []byte {
	panic("not implemented") // TODO: Implement
}

// Topics returns list of topics this message was set to
func (jmsg *msgFromJS) Topics() []string {
	panic("not implemented") // TODO: Implement
}
