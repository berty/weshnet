package ipfsutil

import (
	ipfs_core "github.com/ipfs/kubo/core"
	ipfs_coreapi "github.com/ipfs/kubo/core/coreapi"
	coreiface "github.com/ipfs/kubo/core/coreiface"
	"github.com/libp2p/go-libp2p/core/connmgr"
	ipfs_host "github.com/libp2p/go-libp2p/core/host"
)

type ConnMgr interface {
	connmgr.ConnManager
}

type ExtendedCoreAPI interface {
	coreiface.CoreAPI
	ipfs_host.Host

	ConnMgr() ConnMgr
}

type extendedCoreAPI struct {
	coreiface.CoreAPI
	ipfs_host.Host
}

func (e *extendedCoreAPI) ConnMgr() ConnMgr {
	return e.Host.ConnManager()
}

func NewExtendedCoreAPI(host ipfs_host.Host, api coreiface.CoreAPI) ExtendedCoreAPI {
	return &extendedCoreAPI{
		CoreAPI: api,
		Host:    host,
	}
}

func NewExtendedCoreAPIFromNode(node *ipfs_core.IpfsNode) (ExtendedCoreAPI, error) {
	api, err := ipfs_coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	return NewExtendedCoreAPI(node.PeerHost, api), nil
}
