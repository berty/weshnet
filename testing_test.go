//go:build !js

package weshnet

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"berty.tech/weshnet/pkg/protocoltypes"
)

func TestClient_impl(t *testing.T) {
	var _ Service = (*service)(nil)
	var _ protocoltypes.ProtocolServiceServer = (*service)(nil)
}

func TestEmptyArgs(t *testing.T) {
	// disable ressources manager for test
	os.Setenv("LIBP2P_RCMGR", "false")

	// initialize new client
	client, err := NewService(Opts{})
	require.NoError(t, err)
	err = client.Close()
	require.NoError(t, err)
	client.Close()
}

func TestTestingProtocol(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	opts := TestingOpts{}
	tp, cleanup := NewTestingProtocol(ctx, t, &opts, nil)
	assert.NotNil(t, tp)
	cleanup()
	cancel()
}

func TestTestingProtocolWithMockedPeers(t *testing.T) {
	for amount := 0; amount < 5; amount++ {
		t.Run(fmt.Sprintf("%d-peers", amount), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			opts := TestingOpts{}
			tp, cleanup := NewTestingProtocolWithMockedPeers(ctx, t, &opts, nil, amount)
			assert.NotNil(t, tp)
			cleanup()
			cancel()
		})
	}
}
