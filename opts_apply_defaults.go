//go:build !js

package weshnet

import (
	"context"
	"fmt"
	mrand "math/rand"
	"path/filepath"
	"time"
	"unsafe"

	"berty.tech/go-orbit-db/baseorbitdb"
	"berty.tech/go-orbit-db/pubsub/directchannel"
	"berty.tech/go-orbit-db/pubsub/pubsubraw"
	"berty.tech/weshnet/internal/datastoreutil"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/ipfsutil"
	ipfs_mobile "berty.tech/weshnet/pkg/ipfsutil/mobile"
	"berty.tech/weshnet/pkg/secretstore"
	tinder "berty.tech/weshnet/pkg/tinder"
	pubsub_fix "github.com/berty/go-libp2p-pubsub"
	ds "github.com/ipfs/go-datastore"
	ds_sync "github.com/ipfs/go-datastore/sync"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	backoff "github.com/libp2p/go-libp2p/p2p/discovery/backoff"
	"go.uber.org/zap"
	"moul.io/srand"
)

func (opts *Opts) applyDefaults(ctx context.Context) error {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	rng := mrand.New(mrand.NewSource(srand.MustSecure())) // nolint:gosec // we need to use math/rand here, but it is seeded from crypto/rand

	if err := opts.applyDefaultsGetDatastore(); err != nil {
		return err
	}

	opts.applyPushDefaults()

	if opts.SecretStore == nil {
		secretStore, err := secretstore.NewSecretStore(opts.RootDatastore, &secretstore.NewSecretStoreOptions{
			Logger: opts.Logger,
		})
		if err != nil {
			return errcode.ErrInternal.Wrap(err)
		}

		opts.SecretStore = secretStore
	}

	var mnode *ipfs_mobile.IpfsMobile
	if opts.IpfsCoreAPI == nil {
		dsync := opts.RootDatastore
		if dsync == nil {
			dsync = ds_sync.MutexWrap(ds.NewMapDatastore())
		}

		repo, err := ipfsutil.CreateMockedRepo(dsync)
		if err != nil {
			return err
		}

		mrepo := ipfs_mobile.NewRepoMobile("", repo)
		mnode, err = ipfsutil.NewIPFSMobile(ctx, mrepo, &ipfsutil.MobileOptions{})
		if err != nil {
			return err
		}

		opts.IpfsCoreAPI, err = ipfsutil.NewExtendedCoreAPIFromNode(mnode.IpfsNode)
		if err != nil {
			return err
		}
		opts.Host = mnode.PeerHost()

		oldClose := opts.close
		opts.close = func() error {
			if oldClose != nil {
				_ = oldClose()
			}

			return mnode.Close()
		}
	}

	if opts.Host == nil {
		opts.Host = opts.IpfsCoreAPI
	}

	// setup default tinder service
	if opts.TinderService == nil {
		drivers := []tinder.IDriver{}

		// setup loac disc
		localdisc, err := tinder.NewLocalDiscovery(opts.Logger, opts.Host, rng)
		if err != nil {
			return fmt.Errorf("unable to setup tinder localdiscovery: %w", err)
		}
		drivers = append(drivers, localdisc)

		if mnode != nil {
			dhtdisc := tinder.NewRoutingDiscoveryDriver("dht", mnode.DHT)
			drivers = append(drivers, dhtdisc)
		}

		opts.TinderService, err = tinder.NewService(opts.Host, opts.Logger, drivers...)
		if err != nil {
			return fmt.Errorf("unable to setup tinder service: %w", err)
		}
	}

	if opts.OrbitDBPubSub == nil {
		var err error

		popts := []pubsub_fix.Option{
			pubsub_fix.WithMessageSigning(true),
			pubsub_fix.WithPeerExchange(true),
		}

		backoffstrat := backoff.NewExponentialBackoff(
			time.Second*10, time.Hour,
			backoff.FullJitter,
			time.Second, 10.0, 0, rng)

		cacheSize := 100
		dialTimeout := time.Second * 20
		backoffconnector := func(host host.Host) (*backoff.BackoffConnector, error) {
			return backoff.NewBackoffConnector(host, cacheSize, dialTimeout, backoffstrat)
		}

		adaptater := tinder.NewDiscoveryAdaptater(opts.Logger.Named("disc"), opts.TinderService)
		popts = append(popts, pubsub_fix.WithDiscovery(adaptater, pubsub_fix.WithDiscoverConnector(backoffconnector)))

		// pubsub.DiscoveryPollInterval = m.Node.Protocol.PollInterval
		ps, err := pubsub_fix.NewGossipSub(ctx, opts.Host, popts...)
		if err != nil {
			return fmt.Errorf("unable to init gossipsub: %w", err)
		}

		// @NOTE(gfanton): we need to force cast here until our fix is push
		// upstream on the original go-libp2p-pubsub
		// see: https://github.com/gfanton/go-libp2p-pubsub/commit/8f4fd394f8dfcb3a5eb724a03f9e4e1e33194cbd
		opts.OrbitDBPubSub = (*pubsub.PubSub)(unsafe.Pointer(ps))
	}

	if opts.OrbitDB == nil {
		orbitDirectory := InMemoryDirectory
		if opts.DatastoreDir != InMemoryDirectory {
			orbitDirectory = filepath.Join(opts.DatastoreDir, NamespaceOrbitDBDirectory)
		}

		pubsub := pubsubraw.NewPubSub(opts.OrbitDBPubSub, opts.Host.ID(), opts.Logger, nil)
		odbOpts := &NewOrbitDBOptions{
			NewOrbitDBOptions: baseorbitdb.NewOrbitDBOptions{
				Directory: &orbitDirectory,
				PubSub:    pubsub,
				Logger:    opts.Logger,
			},
			PrometheusRegister:     opts.PrometheusRegister,
			Datastore:              datastoreutil.NewNamespacedDatastore(opts.RootDatastore, ds.NewKey(NamespaceOrbitDBDatastore)),
			SecretStore:            opts.SecretStore,
			GroupMetadataStoreType: opts.GroupMetadataStoreType,
			GroupMessageStoreType:  opts.GroupMessageStoreType,
		}

		if opts.Host != nil {
			odbOpts.DirectChannelFactory = directchannel.InitDirectChannelFactory(opts.Logger, opts.Host)
		}

		odb, err := NewWeshOrbitDB(ctx, opts.IpfsCoreAPI, odbOpts)
		if err != nil {
			return err
		}

		oldClose := opts.close
		opts.close = func() error {
			if oldClose != nil {
				_ = oldClose()
			}

			return odb.Close()
		}

		opts.OrbitDB = odb
	}

	return nil
}
