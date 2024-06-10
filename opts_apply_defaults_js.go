//go:build js

package weshnet

import (
	"context"
	"errors"
	"fmt"
	mrand "math/rand"
	"path/filepath"

	"berty.tech/go-orbit-db/baseorbitdb"
	"berty.tech/go-orbit-db/pubsub/directchannel"
	"berty.tech/weshnet/internal/datastoreutil"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/secretstore"
	tinder "berty.tech/weshnet/pkg/tinder"
	ds "github.com/ipfs/go-datastore"
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

	if opts.IpfsCoreAPI == nil {
		return errors.New("must provide an ipfs core api")
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

		/*
			// FIXME: find a way to inject this from js
			if opts.IpfsCoreAPI != nil {
				dhtdisc := tinder.NewRoutingDiscoveryDriver("dht", mnode.DHT)
				drivers = append(drivers, dhtdisc)
			}
		*/

		opts.TinderService, err = tinder.NewService(opts.Host, opts.Logger, drivers...)
		if err != nil {
			return fmt.Errorf("unable to setup tinder service: %w", err)
		}
	}

	if opts.OrbitDBPubSub != nil {
		return errors.New("can't inject orbit-db PubSub in js")
	}

	if opts.OrbitDB == nil {
		orbitDirectory := InMemoryDirectory
		if opts.DatastoreDir != InMemoryDirectory {
			orbitDirectory = filepath.Join(opts.DatastoreDir, NamespaceOrbitDBDirectory)
		}

		odbOpts := &NewOrbitDBOptions{
			NewOrbitDBOptions: baseorbitdb.NewOrbitDBOptions{
				Directory: &orbitDirectory,
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
