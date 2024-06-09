package ipfsutil

import (
	ipfs_interface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p/core/connmgr"
	ipfs_host "github.com/libp2p/go-libp2p/core/host"
)

type ConnMgr interface {
	connmgr.ConnManager
}

type ExtendedCoreAPI interface {
	ipfs_interface.CoreAPI
	ipfs_host.Host

	ConnMgr() ConnMgr
}
