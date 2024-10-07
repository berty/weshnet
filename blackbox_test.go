package weshnet_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"berty.tech/weshnet"
	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/secretstore"
	"berty.tech/weshnet/pkg/testutil"
)

func TestTestingClient_impl(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	secretStore, err := secretstore.NewInMemSecretStore(nil)
	require.NoError(t, err)

	client, cleanup := weshnet.TestingService(ctx, t, weshnet.Opts{
		Logger:      logger,
		SecretStore: secretStore,
	})
	defer cleanup()

	// test service
	_, _ = client.ServiceGetConfiguration(ctx, &protocoltypes.ServiceGetConfiguration_Request{})
	status := client.Status()
	expected := weshnet.Status{}
	assert.Equal(t, expected, status)
}

func ExampleNewInMemoryServiceClient_basic() {
	// disable ressources manager for test
	os.Setenv("LIBP2P_RCMGR", "false")

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
	// disable ressources manager for test
	os.Setenv("LIBP2P_RCMGR", "false")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create a temporary path to host data of our persistent service
	path, err := os.MkdirTemp("", "weshnet-test-persistent")
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

		peerid = ret.PeerId

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

		if peerid != ret.PeerId {
			panic("peerid should be identical")
		}
	}

	// Output:
}

func ExampleNewServiceClient_basic() {
	// disable ressources manager for test
	os.Setenv("LIBP2P_RCMGR", "false")

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
	// disable ressources manager for test
	os.Setenv("LIBP2P_RCMGR", "false")

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
