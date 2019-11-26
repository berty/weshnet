package ipfsutil

import (
	"context"
	"time"

	"berty.tech/go/pkg/errcode"
	"github.com/libp2p/go-libp2p/p2p/discovery"

	ipfs_core "github.com/ipfs/go-ipfs/core"
	ipfs_interface "github.com/ipfs/interface-go-ipfs-core"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
)

type DiscoveryNotifee struct {
	api ipfs_interface.CoreAPI
}

func (n *DiscoveryNotifee) HandlePeerFound(pi peerstore.PeerInfo) {
	if err := n.api.Swarm().Connect(context.Background(), pi); err != nil {
		println("HandlePeerFound: Unable to connect to peer", err.Error())
	}
}

func OptionMDNSDiscovery(ctx context.Context, node *ipfs_core.IpfsNode, api ipfs_interface.CoreAPI) error {
	s, err := discovery.NewMdnsService(ctx, node.PeerHost, time.Second, "")
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	n := &DiscoveryNotifee{api: api}

	s.RegisterNotifee(n)

	return nil
}
