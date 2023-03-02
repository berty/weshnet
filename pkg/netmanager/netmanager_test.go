package netmanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNetManager(t *testing.T) {
	initial := ConnectivityInfo{
		State:     ConnectivityStateOn,
		NetType:   ConnectivityNetWifi,
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

	ok, eventType := netmanager.WaitForStateChange(ctx, &state, ConnectivityChanged)

	require.Equal(t, a, netmanager.GetCurrentState())
	require.True(t, ok)
	require.Equal(t, ConnectivityStateChanged, eventType)
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

	ok, eventType := netmanager.WaitForStateChange(ctx, &state, ConnectivityChanged)

	require.Equal(t, b, netmanager.GetCurrentState())
	require.True(t, ok)
	require.Equal(t, ConnectivityStateChanged, eventType)
}

func TestNetManagerFilterUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := ConnectivityInfo{
		State: ConnectivityStateOff,
	}
	b := ConnectivityInfo{
		State:        ConnectivityStateOn,
		NetType:      ConnectivityNetCellular,
		CellularType: ConnectivityCellular3G,
	}
	state := ConnectivityInfo{}

	netmanager := NewNetManager(state)

	netmanager.UpdateState(a)
	require.Equal(t, a, netmanager.GetCurrentState())
	netmanager.UpdateState(b)
	require.Equal(t, b, netmanager.GetCurrentState())

	ok, eventType := netmanager.WaitForStateChange(ctx, &state, ConnectivityCellularTypeChanged)

	require.Equal(t, b, netmanager.GetCurrentState())
	require.True(t, ok)
	require.Equal(t, ConnectivityStateChanged | ConnectivityNetTypeChanged | ConnectivityCellularTypeChanged, eventType)
}
