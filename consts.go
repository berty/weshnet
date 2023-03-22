package weshnet

import (
	"berty.tech/go-orbit-db/cache/cacheleveldown"
)

const (
	NamespaceOrbitDBDatastore = "orbitdb_datastore"
	NamespaceOrbitDBDirectory = "orbitdb"
	NamespaceIPFSDatastore    = "ipfs_datastore"
)

var InMemoryDirectory = cacheleveldown.InMemoryDirectory
