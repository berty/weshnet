package ipfsutil

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"time"

	ipfs_ds "github.com/ipfs/go-datastore"
	ipfs_cfg "github.com/ipfs/kubo/config"
	ipfs_loader "github.com/ipfs/kubo/plugin/loader"
	ipfs_repo "github.com/ipfs/kubo/repo"
	ipfs_fsrepo "github.com/ipfs/kubo/repo/fsrepo"
	p2p_ci "github.com/libp2p/go-libp2p/core/crypto"
	p2p_peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"

	"berty.tech/weshnet/pkg/errcode"
)

// defaultConnMgrHighWater is the default value for the connection managers
// 'high water' mark
const defaultConnMgrHighWater = 200

// defaultConnMgrLowWater is the default value for the connection managers 'low
// water' mark
const defaultConnMgrLowWater = 150

// defaultConnMgrGracePeriod is the default value for the connection managers
// grace period
const defaultConnMgrGracePeriod = time.Second * 20

// @NOTE(gfanton): this will be removed with gomobile-ipfs
var plugins *ipfs_loader.PluginLoader

func CreateMockedRepo(dstore ipfs_ds.Batching) (ipfs_repo.Repo, error) {
	c, err := CreateBaseConfig()
	if err != nil {
		return nil, err
	}

	return &ipfs_repo.Mock{
		D: dstore,
		C: *c,
	}, nil
}

func CreateOrLoadMockedRepo(dstore ipfs_ds.Batching) (ipfs_repo.Repo, error) {
	c, err := CreateBaseConfig()
	if err != nil {
		return nil, err
	}

	return &ipfs_repo.Mock{
		D: dstore,
		C: *c,
	}, nil
}

func LoadRepoFromPath(path string) (ipfs_repo.Repo, error) {
	dir, _ := filepath.Split(path)
	if _, err := LoadPlugins(dir); err != nil {
		return nil, errors.Wrap(err, "failed to load plugins")
	}

	if !ipfs_fsrepo.IsInitialized(path) {
		cfg, err := CreateBaseConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create base config: %w", err)
		}

		ucfg, err := upgradeToPersistentConfig(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to upgrade repo")
		}

		if err := ipfs_fsrepo.Init(path, ucfg); err != nil {
			return nil, fmt.Errorf("failed to init ipfs repo: %w", err)
		}
	}

	return ipfs_fsrepo.Open(path)
}

var DefaultSwarmListeners = []string{
	//"/ip4/0.0.0.0/udp/0/quic",
	//"/ip6/::/udp/0/quic",
	//"/ip4/0.0.0.0/tcp/0",
	//"/ip6/::/tcp/0",
	//"/ip4/0.0.0.0/tcp/0/http/p2p-webrtc-direct",
}

func CreateBaseConfig() (*ipfs_cfg.Config, error) {
	c := ipfs_cfg.Config{}

	// set default bootstrap
	c.Bootstrap = ipfs_cfg.DefaultBootstrapAddresses
	c.Peering.Peers = []p2p_peer.AddrInfo{}

	// Identity
	if err := ResetRepoIdentity(&c); err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	// Discovery
	c.Discovery.MDNS.Enabled = true

	// swarm listeners
	c.Addresses.Swarm = DefaultSwarmListeners

	// Swarm
	c.Swarm.RelayClient.Enabled = ipfs_cfg.True
	c.Swarm.ConnMgr = ipfs_cfg.ConnMgr{
		LowWater:    ipfs_cfg.NewOptionalInteger(defaultConnMgrLowWater),
		HighWater:   ipfs_cfg.NewOptionalInteger(defaultConnMgrHighWater),
		GracePeriod: ipfs_cfg.NewOptionalDuration(defaultConnMgrGracePeriod),
		Type:        ipfs_cfg.NewOptionalString("basic"),
	}

	c.Routing = ipfs_cfg.Routing{
		Type: ipfs_cfg.NewOptionalString("dhtclient"),
	}

	return &c, nil
}

func ResetRepoIdentity(c *ipfs_cfg.Config) error {
	priv, pub, err := p2p_ci.GenerateKeyPairWithReader(p2p_ci.Ed25519, 2048, crand.Reader) // nolint:staticcheck
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	pid, err := p2p_peer.IDFromPublicKey(pub) // nolint:staticcheck
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	privkeyb, err := p2p_ci.MarshalPrivateKey(priv)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	// Identity
	c.Identity.PeerID = pid.Pretty()
	c.Identity.PrivKey = base64.StdEncoding.EncodeToString(privkeyb)

	return nil
}

func upgradeToPersistentConfig(cfg *ipfs_cfg.Config) (*ipfs_cfg.Config, error) {
	cfgCopy, err := cfg.Clone()
	if err != nil {
		return nil, err
	}

	// setup the node mount points.
	cfgCopy.Mounts = ipfs_cfg.Mounts{
		IPFS: "/ipfs",
		IPNS: "/ipns",
	}

	cfgCopy.Ipns = ipfs_cfg.Ipns{
		ResolveCacheSize: 128,
	}

	cfgCopy.Reprovider = ipfs_cfg.Reprovider{
		Interval: ipfs_cfg.NewOptionalDuration(time.Hour * 12),
		Strategy: ipfs_cfg.NewOptionalString("all"),
	}

	cfgCopy.Datastore = ipfs_cfg.Datastore{
		StorageMax:         "10GB",
		StorageGCWatermark: 90, // 90%
		GCPeriod:           "1h",
		BloomFilterSize:    0,
		Spec: map[string]interface{}{
			"type": "mount",
			"mounts": []interface{}{
				map[string]interface{}{
					"mountpoint": "/blocks",
					"type":       "measure",
					"prefix":     "flatfs.datastore",
					"child": map[string]interface{}{
						"type":      "flatfs",
						"path":      "blocks",
						"sync":      true,
						"shardFunc": "/repo/flatfs/shard/v1/next-to-last/2",
					},
				},
				map[string]interface{}{
					"mountpoint": "/",
					"type":       "measure",
					"prefix":     "leveldb.datastore",
					"child": map[string]interface{}{
						"type":        "levelds",
						"path":        "datastore",
						"compression": "none",
					},
				},
			},
		},
	}

	return cfgCopy, nil
}

func ResetExistingRepoIdentity(repo ipfs_repo.Repo) (ipfs_repo.Repo, error) {
	cfg, err := repo.Config()
	if err != nil {
		_ = repo.Close()
		return nil, errcode.ErrInternal.Wrap(err)
	}

	if err := ResetRepoIdentity(cfg); err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	updatedCfg, err := upgradeToPersistentConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upgrade repo")
	}

	err = repo.SetConfig(updatedCfg)
	if err != nil {
		return nil, errcode.ErrInternal.Wrap(err)
	}

	return repo, nil
}

func LoadPlugins(repoPath string) (*ipfs_loader.PluginLoader, error) { // nolint:unparam
	if plugins != nil {
		return plugins, nil
	}

	pluginpath := filepath.Join(repoPath, "plugins")

	lp, err := ipfs_loader.NewPluginLoader(pluginpath)
	if err != nil {
		return nil, err
	}

	if err = lp.Initialize(); err != nil {
		return nil, err
	}

	if err = lp.Inject(); err != nil {
		return nil, err
	}

	plugins = lp
	return lp, nil
}
