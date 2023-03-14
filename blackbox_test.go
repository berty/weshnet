package weshnet_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	keystore "github.com/ipfs/go-ipfs-keystore"
	"github.com/stretchr/testify/assert"

	"berty.tech/weshnet"
	"berty.tech/weshnet/pkg/cryptoutil"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/testutil"
)

func TestTestingClient_impl(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	client, cleanup := weshnet.TestingService(ctx, t, weshnet.Opts{
		Logger:         logger,
		DeviceKeystore: cryptoutil.NewDeviceKeystore(keystore.NewMemKeystore(), nil),
	})
	defer cleanup()

	// test service
	_, _ = client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
	status := client.Status()
	expected := weshnet.Status{}
	assert.Equal(t, expected, status)
}

func ExampleNewInMemoryServiceClient_basic() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := weshnet.NewInMemoryServiceClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ret, err := client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
	if err != nil {
		panic(err)
	}

	for _, listener := range ret.Listeners {
		if listener == "/p2p-circuit" {
			fmt.Println(listener)
		}
	}

	// Output:
	// /p2p-circuit
}

func ExampleNewPersistentServiceClient_basic() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create a temporary path to host data of our persistant service
	path, err := ioutil.TempDir("", "weshnet-test-persistant")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(path)

	var peerid string
	// open once
	{
		client, err := weshnet.NewPersistentServiceClient(path)
		if err != nil {
			panic(err)
		}

		ret, err := client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
		if err != nil {
			panic(err)
		}

		peerid = ret.PeerID

		if err := client.Close(); err != nil {
			panic(err)
		}
	}

	// open twice
	{
		client, err := weshnet.NewPersistentServiceClient(path)
		if err != nil {
			panic(err)
		}
		defer client.Close()

		ret, err := client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
		if err != nil {
			panic(err)
		}

		if peerid != ret.PeerID {
			panic("peerid should be identical")
		}
	}

	// Output:
}

func ExampleNewServiceClient_basic() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := weshnet.NewServiceClient(weshnet.Opts{})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ret, err := client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
	if err != nil {
		panic(err)
	}

	for _, listener := range ret.Listeners {
		if listener == "/p2p-circuit" {
			fmt.Println(listener)
		}
	}

	// Output:
	// /p2p-circuit
}

func ExampleNewService_basic() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := weshnet.NewService(weshnet.Opts{})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ret, err := client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
	if err != nil {
		panic(err)
	}

	for _, listener := range ret.Listeners {
		if listener == "/p2p-circuit" {
			fmt.Println(listener)
		}
	}

	// Output:
	// /p2p-circuit
}

// FIXME: create examples that actually use groups and contacts
