//go:build js

package wasm

import (
	"context"
	"fmt"
	"syscall/js"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	iface "github.com/ipfs/interface-go-ipfs-core"
	_ "github.com/ipld/go-ipld-prime/codec/cbor"
	"github.com/pkg/errors"
)

type dagAPIFromJS struct {
	ipld.DAGService
}

var _ iface.APIDagService = (*dagAPIFromJS)(nil)

func newDagFromJS(helia js.Value) iface.APIDagService {
	svc := merkledag.NewDAGService(&blocksvcFromJS{helia: helia})
	return &dagAPIFromJS{
		DAGService: svc,
	}
}

/*

// Add adds a node to this DAG.
func (jdag *dagAPIFromJS) Add(ctx context.Context, data ipld.Node) error {
	fmt.Println("FIXME: ignored input context in dag.Add")

	addType := reflect.TypeOf(data)
	fmt.Println("dag add", addType)

	var (
		bs  []byte
		cid cid.Cid
		err error
	)
	switch node := data.(type) {
	case *cbornode.Node:
		bs, err = node.MarshalJSON()
		if err != nil {
			return fmt.Errorf("failed to marshal cbor: %w", err)
		}
		cid = node.Cid()
	case *dag.ProtoNode:
		bs, err = node.Marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal proto: %w", err)
		}
	default:
		return fmt.Errorf("unknown ipld type %s", addType)
	}

	dst := js.Global().Get("Uint8Array").New(len(bs))
	n := js.CopyBytesToJS(dst, bs)
	if n != len(bs) {
		return errors.New("failed to copy bytes")
	}
	jsCID := js.Global().Get("Multiformats").Get("CID").Call("parse", cid.String())
	_, err = await(jdag.helia.Get("blockstore").Call("put", jsCID, dst))
	if err != nil {
		return fmt.Errorf("failed to put in blockstore: %w", err)
	}
	fmt.Println("dag added", addType, cid)
	return nil
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
	fmt.Println("dag get", cid)
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
	node, err := ipld.Decode(blocks)
	if err != nil {
		return nil, err
	}
	fmt.Println("dag got", cid, len(bs), "bytes")
	return node, nil
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

*/

// Pinning returns special NodeAdder which recursively pins added nodes
func (jdag *dagAPIFromJS) Pinning() ipld.NodeAdder {
	panic("dag Pinning not implemented") // TODO: Implement
}

type blocksvcFromJS struct {
	helia js.Value
}

var _ blockservice.BlockService = (*blocksvcFromJS)(nil)

func (jbsvc *blocksvcFromJS) Close() error {
	panic("blocksvc Close not implemented") // TODO: Implement
}

// GetBlock gets the requested block.
func (jbsvc *blocksvcFromJS) GetBlock(ctx context.Context, cid cid.Cid) (blocks.Block, error) {
	fmt.Println("FIXME: ignored input context in blockstore.GetBlock")
	fmt.Println("block get", cid)
	jsCID := js.Global().Get("Multiformats").Get("CID").Call("parse", cid.String())
	jsNode, err := await(jbsvc.helia.Get("blockstore").Call("get", jsCID))
	if err != nil {
		return nil, fmt.Errorf("failed to get block from js: %w", err)
	}
	bs := make([]byte, jsNode.Get("length").Int())
	js.CopyBytesToGo(bs, jsNode)
	block, err := blocks.NewBlockWithCid(bs, cid)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate block: %w", err)
	}
	fmt.Println("block got", cid)
	return block, nil
}

// GetBlocks does a batch request for the given cids, returning blocks as
// they are found, in no particular order.
//
// It may not be able to find all requested blocks (or the context may
// be canceled). In that case, it will close the channel early. It is up
// to the consumer to detect this situation and keep track which blocks
// it has received and which it hasn't.
func (jbsvc *blocksvcFromJS) GetBlocks(ctx context.Context, ks []cid.Cid) <-chan blocks.Block {
	panic("blocksvc GetBlocks not implemented") // TODO: Implement
}

// Blockstore returns a reference to the underlying blockstore
func (jbsvc *blocksvcFromJS) Blockstore() blockstore.Blockstore {
	panic("blocksvc Blockstore not implemented") // TODO: Implement
}

// Exchange returns a reference to the underlying exchange (usually bitswap)
func (jbsvc *blocksvcFromJS) Exchange() exchange.Interface {
	panic("blocksvc Exchange not implemented") // TODO: Implement
}

// AddBlock puts a given block to the underlying datastore
func (jbsvc *blocksvcFromJS) AddBlock(ctx context.Context, o blocks.Block) error {
	fmt.Println("FIXME: ignored input context in blockstore.AddBlock")
	id := o.Cid()
	fmt.Println("block add", id)
	bs := o.RawData()
	dst := js.Global().Get("Uint8Array").New(len(bs))
	n := js.CopyBytesToJS(dst, bs)
	if n != len(bs) {
		return errors.New("failed to copy bytes")
	}
	jsCID := js.Global().Get("Multiformats").Get("CID").Call("parse", id.String())
	if _, err := await(jbsvc.helia.Get("blockstore").Call("put", jsCID, dst)); err != nil {
		return fmt.Errorf("failed to put in blockstore: %w", err)
	}
	fmt.Println("block added", id)
	return nil
}

// AddBlocks adds a slice of blocks at the same time using batching
// capabilities of the underlying datastore whenever possible.
func (jbsvc *blocksvcFromJS) AddBlocks(ctx context.Context, bs []blocks.Block) error {
	panic("blocksvc AddBlocks not implemented") // TODO: Implement
}

// DeleteBlock deletes the given block from the blockservice.
func (jbsvc *blocksvcFromJS) DeleteBlock(ctx context.Context, o cid.Cid) error {
	panic("blocksvc DeleteBlock not implemented") // TODO: Implement
}
