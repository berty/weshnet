package ipfsutil

import (
	ipfs_interface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p-core/connmgr"
	ipfs_host "github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

type ConnMgr interface {
	TagPeer(peer.ID, string, int)
	UntagPeer(p peer.ID, tag string)
	Protect(id peer.ID, tag string)
	Unprotect(id peer.ID, tag string) (protected bool)
	GetTagInfo(p peer.ID) *connmgr.TagInfo
}

type ExtendedCoreAPI interface {
	ipfs_interface.CoreAPI
	ipfs_host.Host

	ConnMgr() ConnMgr
}

type extendedCoreAPI struct {
	ipfs_interface.CoreAPI
	ipfs_host.Host
}

func (e *extendedCoreAPI) ConnMgr() ConnMgr {
	return e.Host.ConnManager()
}

func NewExtendedCoreAPI(host ipfs_host.Host, api ipfs_interface.CoreAPI) ExtendedCoreAPI {
	return &extendedCoreAPI{
		CoreAPI: api,
		Host:    host,
	}
}
