//go:build !js

package ipfsutil

import (
	ipfs_interface "github.com/ipfs/interface-go-ipfs-core"
	ipfs_core "github.com/ipfs/kubo/core"
	ipfs_coreapi "github.com/ipfs/kubo/core/coreapi"
	ipfs_host "github.com/libp2p/go-libp2p/core/host"
)

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

func NewExtendedCoreAPIFromNode(node *ipfs_core.IpfsNode) (ExtendedCoreAPI, error) {
	api, err := ipfs_coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	return NewExtendedCoreAPI(node.PeerHost, api), nil
}
