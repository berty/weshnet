//go:build js

package main

import (
	"context"
	"errors"
	"fmt"
	"syscall/js"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-libipfs/blocks"
	iface "github.com/ipfs/interface-go-ipfs-core"
)

type dagAPIFromJS struct {
	helia js.Value
}

var _ iface.APIDagService = (*dagAPIFromJS)(nil)

// Add adds a node to this DAG.
func (jdag *dagAPIFromJS) Add(ctx context.Context, data ipld.Node) error {
	fmt.Println("FIXME: ignored input context in dag.Add")
	/*
			const buf = codec.encode(obj)
		    const hash = await (options.hasher ?? sha256).digest(buf)
		    const cid = CID.createV1(codec.code, hash)

		    await this.components.blockstore.put(cid, buf, options)

		    return cid
	*/
	cid := data.Cid()
	bs := data.RawData()
	dst := js.Global().Get("Uint8Array").New(len(bs))
	n := js.CopyBytesToJS(dst, bs)
	if n != len(bs) {
		return errors.New("failed to copy bytes")
	}
	jsCID := js.Global().Get("Multiformats").Get("CID").Call("parse", cid.String())
	_, err := await(jdag.helia.Get("blockstore").Call("put", jsCID, dst))
	return err
}

// AddMany adds many nodes to this DAG.
//
// Consider using the Batch NodeAdder (`NewBatch`) if you make
// extensive use of this function.
func (jdag *dagAPIFromJS) AddMany(_ context.Context, _ []ipld.Node) error {
	panic("not implemented") // TODO: Implement
}

// Get retrieves nodes by CID. Depending on the NodeGetter
// implementation, this may involve fetching the Node from a remote
// machine; consider setting a deadline in the context.
func (jdag *dagAPIFromJS) Get(ctx context.Context, cid cid.Cid) (ipld.Node, error) {
	fmt.Println("FIXME: ignored input context in dag.Get")
	jsCID := js.Global().Get("Multiformats").Get("CID").Call("parse", cid.String())
	jsNode, err := await(jdag.helia.Get("blockstore").Call("get", jsCID))
	if err != nil {
		return nil, err
	}
	bs := make([]byte, jsNode.Get("length").Int())
	js.CopyBytesToGo(bs, jsNode)
	blocks, err := blocks.NewBlockWithCid(bs, cid)
	if err != nil {
		return nil, err
	}
	return ipld.Decode(blocks)
}

// GetMany returns a channel of NodeOptions given a set of CIDs.
func (jdag *dagAPIFromJS) GetMany(_ context.Context, _ []cid.Cid) <-chan *ipld.NodeOption {
	panic("not implemented") // TODO: Implement
}

// Remove removes a node from this DAG.
//
// Remove returns no error if the requested node is not present in this DAG.
func (jdag *dagAPIFromJS) Remove(context.Context, cid.Cid) error {
	panic("not implemented") // TODO: Implement
}

// RemoveMany removes many nodes from this DAG.
//
// It returns success even if the nodes were not present in the DAG.
func (jdag *dagAPIFromJS) RemoveMany(context.Context, []cid.Cid) error {
	panic("not implemented") // TODO: Implement
}

// Pinning returns special NodeAdder which recursively pins added nodes
func (jdag *dagAPIFromJS) Pinning() ipld.NodeAdder {
	panic("not implemented") // TODO: Implement
}
