//go:build js

package weshnet

import (
	"fmt"
	"syscall/js"

	"berty.tech/weshnet/pkg/ipfsutil/wasm"
)

// NewPersistentServiceClientJS creates a Wesh protocol service using persistent storage files in the
// directory given by the directory path. If the directory doesn't exist, this creates it with files
// of a new Wesh account and peer identity. (If the directory doesn't exist, this will create it only
// if the parent directory exists. Otherwise you must first create the parent directories.) However,
// if the persistent storage files already exist, then this opens them to use the existing Wesh
// account and peer identity. This returns a gRPC ServiceClient which uses a direct in-memory
// connection. When finished, you must call Close().
func NewPersistentServiceClientJS(path string, helia js.Value) (ServiceClient, error) {
	var opts Opts

	opts.DatastoreDir = path

	var err error
	if opts.IpfsCoreAPI, err = wasm.NewCoreAPIFromJS(helia); err != nil {
		return nil, fmt.Errorf("unable to wrap helia: %w", err)
	}

	var cleanupLogger func()
	if opts.Logger, cleanupLogger, err = setupDefaultLogger(); err != nil {
		return nil, fmt.Errorf("uanble to setup logger: %w", err)
	}

	cl, err := NewServiceClient(opts)
	if err != nil {
		return nil, err
	}

	return &persistentServiceClient{
		ServiceClient: cl,
		cleanup:       cleanupLogger,
	}, nil
}
