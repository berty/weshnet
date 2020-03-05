package ipfsutil

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	ipfs_datastore "github.com/ipfs/go-datastore"
	ipfs_cfg "github.com/ipfs/go-ipfs-config"
	ipfs_node "github.com/ipfs/go-ipfs/core/node"
	ipfs_libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	ipfs_repo "github.com/ipfs/go-ipfs/repo"
	libp2p_ci "github.com/libp2p/go-libp2p-core/crypto" // nolint:staticcheck
	libp2p_peer "github.com/libp2p/go-libp2p-core/peer" // nolint:staticcheck

	"berty.tech/berty/go/pkg/errcode"
)

type BuildOpts struct {
	SwarmAddresses []string
}

func CreateBuildConfigWithDatastore(opts *BuildOpts, ds ipfs_datastore.Batching) (*ipfs_node.BuildCfg, error) {
	if opts == nil {
		opts = &BuildOpts{}
	}

	repo, err := CreateRepo(ds, opts)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	routing := ipfs_libp2p.DHTOption
	hostopts := ipfs_libp2p.DefaultHostOption
	return &ipfs_node.BuildCfg{
		Online:                      true,
		Permanent:                   true,
		DisableEncryptedConnections: false,
		NilRepo:                     false,
		Routing:                     routing,
		Host:                        hostopts,
		Repo:                        repo,
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	}, nil
}

func CreateRepo(dstore ipfs_datastore.Batching, opts *BuildOpts) (ipfs_repo.Repo, error) {
	c := ipfs_cfg.Config{}
	priv, pub, err := libp2p_ci.GenerateKeyPairWithReader(libp2p_ci.RSA, 2048, rand.Reader) // nolint:staticcheck
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	pid, err := libp2p_peer.IDFromPublicKey(pub) // nolint:staticcheck
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	privkeyb, err := priv.Bytes()
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	c.Bootstrap = ipfs_cfg.DefaultBootstrapAddresses

	if len(opts.SwarmAddresses) != 0 {
		c.Addresses.Swarm = opts.SwarmAddresses
	} else {
		portOffsetBI, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			return nil, errcode.TODO.Wrap(err)
		}

		portOffset := portOffsetBI.Int64() % 100

		c.Addresses.Swarm = []string{
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", 4001+portOffset),
			fmt.Sprintf("/ip6/0.0.0.0/tcp/%d", 4001+portOffset),
		}
	}

	fmt.Printf("IPFS listening on %s\n", strings.Join(c.Addresses.Swarm, ", "))

	c.Identity.PeerID = pid.Pretty()
	c.Identity.PrivKey = base64.StdEncoding.EncodeToString(privkeyb)
	c.Discovery.MDNS.Enabled = true
	c.Discovery.MDNS.Interval = 1

	return &ipfs_repo.Mock{
		D: dstore,
		C: c,
	}, nil
}
