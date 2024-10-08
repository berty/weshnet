package weshnet

import (
	"context"

	datastore "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"

	"berty.tech/go-orbit-db/address"
	"berty.tech/go-orbit-db/cache"
	"berty.tech/weshnet/v2/internal/datastoreutil"
)

type datastoreCache struct {
	ds datastore.Batching
}

//nolint:revive
func (d *datastoreCache) Load(directory string, dbAddress address.Address) (datastore.Datastore, error) {
	return datastoreutil.NewNamespacedDatastore(d.ds, datastore.NewKey(dbAddress.String())), nil
}

func (d *datastoreCache) Close() error {
	return nil
}

//nolint:revive
func (d *datastoreCache) Destroy(directory string, dbAddress address.Address) error {
	keys, err := datastoreutil.NewNamespacedDatastore(d.ds, datastore.NewKey(dbAddress.String())).Query(context.TODO(), query.Query{KeysOnly: true})
	if err != nil {
		return nil
	}

	for {
		val, hasValue := keys.NextSync()
		if !hasValue {
			return nil
		}

		if err := d.ds.Delete(context.TODO(), datastore.NewKey(val.Key)); err != nil {
			return err
		}
	}
}

func NewOrbitDatastoreCache(ds datastore.Batching) cache.Interface {
	return &datastoreCache{
		ds: datastoreutil.NewNamespacedDatastore(ds, datastore.NewKey(NamespaceOrbitDBDatastore)),
	}
}

var _ cache.Interface = (*datastoreCache)(nil)
