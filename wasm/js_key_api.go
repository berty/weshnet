//go:build js

package main

import (
	"context"
	"syscall/js"

	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p/core/peer"
)

type keyAPIFromJS struct {
	helia js.Value
}

var _ iface.KeyAPI = (*keyAPIFromJS)(nil)

// Generate generates new key, stores it in the keystore under the specified
// name and returns a base58 encoded multihash of it's public key
func (kapi *keyAPIFromJS) Generate(ctx context.Context, name string, opts ...options.KeyGenerateOption) (iface.Key, error) {
	panic("not implemented") // TODO: Implement
}

// Rename renames oldName key to newName. Returns the key and whether another
// key was overwritten, or an error
func (kapi *keyAPIFromJS) Rename(ctx context.Context, oldName string, newName string, opts ...options.KeyRenameOption) (iface.Key, bool, error) {
	panic("not implemented") // TODO: Implement
}

// List lists keys stored in keystore
func (kapi *keyAPIFromJS) List(ctx context.Context) ([]iface.Key, error) {
	panic("not implemented") // TODO: Implement
}

// Self returns the 'main' node key
func (kapi *keyAPIFromJS) Self(ctx context.Context) (iface.Key, error) {
	keychain := kapi.helia.Get("libp2p").Get("services").Get("keychain")
	key, err := await(keychain.Call("findKeyByName", "self"))
	if err != nil {
		return nil, err
	}
	return &keyFromJS{key: key}, nil
}

// Remove removes keys from keystore. Returns ipns path of the removed key
func (kapi *keyAPIFromJS) Remove(ctx context.Context, name string) (iface.Key, error) {
	panic("not implemented") // TODO: Implement
}

type keyFromJS struct {
	key js.Value
}

var _ iface.Key = (*keyFromJS)(nil)

// Key returns key name
func (jk *keyFromJS) Name() string {
	return jk.key.Get("name").String()
}

// Path returns key path
func (jk *keyFromJS) Path() path.Path {
	consoleLog("FIXME: used hacked-in iface.Key.Path")
	id := jk.ID()
	return path.New(id.String())
}

// ID returns key PeerID
func (jk *keyFromJS) ID() peer.ID {
	raw := jk.key.Get("id").String()
	id, err := peer.Decode(raw)
	if err != nil {
		panic(err)
	}
	return id
}
