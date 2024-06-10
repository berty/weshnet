package wasm

import (
	"context"
	"fmt"

	"berty.tech/weshnet/pkg/ipfsutil"
	"github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type weshConnMgrFromJS struct {
}

var _ ipfsutil.ConnMgr = (*weshConnMgrFromJS)(nil)

// TagPeer tags a peer with a string, associating a weight with the tag.
func (jwconnmgr *weshConnMgrFromJS) TagPeer(_ peer.ID, _ string, _ int) {
	fmt.Println("FIXME: ignored TagPeer")
}

// Untag removes the tagged value from the peer.
func (jwconnmgr *weshConnMgrFromJS) UntagPeer(p peer.ID, tag string) {
	panic("not implemented") // TODO: Implement
}

// UpsertTag updates an existing tag or inserts a new one.
//
// The connection manager calls the upsert function supplying the current
// value of the tag (or zero if inexistent). The return value is used as
// the new value of the tag.
func (jwconnmgr *weshConnMgrFromJS) UpsertTag(p peer.ID, tag string, upsert func(int) int) {
	panic("not implemented") // TODO: Implement
}

// GetTagInfo returns the metadata associated with the peer,
// or nil if no metadata has been recorded for the peer.
func (jwconnmgr *weshConnMgrFromJS) GetTagInfo(p peer.ID) *connmgr.TagInfo {
	panic("not implemented") // TODO: Implement
}

// TrimOpenConns terminates open connections based on an implementation-defined
// heuristic.
func (jwconnmgr *weshConnMgrFromJS) TrimOpenConns(ctx context.Context) {
	panic("not implemented") // TODO: Implement
}

// Notifee returns an implementation that can be called back to inform of
// opened and closed connections.
func (jwconnmgr *weshConnMgrFromJS) Notifee() network.Notifiee {
	panic("not implemented") // TODO: Implement
}

// Protect protects a peer from having its connection(s) pruned.
//
// Tagging allows different parts of the system to manage protections without interfering with one another.
//
// Calls to Protect() with the same tag are idempotent. They are not refcounted, so after multiple calls
// to Protect() with the same tag, a single Unprotect() call bearing the same tag will revoke the protection.
func (jwconnmgr *weshConnMgrFromJS) Protect(id peer.ID, tag string) {
	panic("not implemented") // TODO: Implement
}

// Unprotect removes a protection that may have been placed on a peer, under the specified tag.
//
// The return value indicates whether the peer continues to be protected after this call, by way of a different tag.
// See notes on Protect() for more info.
func (jwconnmgr *weshConnMgrFromJS) Unprotect(id peer.ID, tag string) (protected bool) {
	panic("not implemented") // TODO: Implement
}

// IsProtected returns true if the peer is protected for some tag; if the tag is the empty string
// then it will return true if the peer is protected for any tag
func (jwconnmgr *weshConnMgrFromJS) IsProtected(id peer.ID, tag string) (protected bool) {
	panic("not implemented") // TODO: Implement
}

// Close closes the connection manager and stops background processes.
func (jwconnmgr *weshConnMgrFromJS) Close() error {
	panic("not implemented") // TODO: Implement
}
