package netmanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNetManager(t *testing.T) {
	initial := ConnectivityInfo{
		State: ConnectivityStateOn,
		NetType:  ConnectivityNetWifi,
		Bluetooth: ConnectivityStateOn,
	}

	netmanager := NewNetManager(initial)

	require.Equal(t, initial, netmanager.GetCurrentState())
	initial.State = ConnectivityStateOff
	require.NotEqual(t, initial, netmanager.GetCurrentState())
}

func TestNetManagerSingleUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := ConnectivityInfo{
		State: ConnectivityStateOn,
	}
	state := ConnectivityInfo{}

	netmanager := NewNetManager(state)

	netmanager.UpdateState(a)
	require.Equal(t, a, netmanager.GetCurrentState())

	netmanager.WaitForStateChange(ctx, &state, NetManagerConnectivityChanged)

	require.Equal(t, a, netmanager.GetCurrentState())
}

func TestNetManagerDoubleUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := ConnectivityInfo{
		State: ConnectivityStateOn,
	}
	b := ConnectivityInfo{
		State: ConnectivityStateOff,
	}
	state := ConnectivityInfo{}

	netmanager := NewNetManager(state)

	netmanager.UpdateState(a)
	require.Equal(t, a, netmanager.GetCurrentState())
	netmanager.UpdateState(b)
	require.Equal(t, b, netmanager.GetCurrentState())

	netmanager.WaitForStateChange(ctx, &state, NetManagerConnectivityChanged)
	require.Equal(t, b, netmanager.GetCurrentState())
}

func TestNetManagerFilterUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := ConnectivityInfo{
		State: ConnectivityStateOff,
	}
	b := ConnectivityInfo{
		State: ConnectivityStateOn,
		NetType: ConnectivityNetCellular,
		CellularType: ConnectivityCellular3G,
	}
	state := ConnectivityInfo{}

	netmanager := NewNetManager(state)

	netmanager.UpdateState(a)
	require.Equal(t, a, netmanager.GetCurrentState())
	netmanager.UpdateState(b)
	require.Equal(t, b, netmanager.GetCurrentState())

	netmanager.WaitForStateChange(ctx, &state, NetManagerConnectivityCellularTypeChanged)
	require.Equal(t, b, netmanager.GetCurrentState())
}
