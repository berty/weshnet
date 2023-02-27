package netmanager

import (
	"context"
	"sync"

	"berty.tech/weshnet/internal/notify"
)

type NetManager struct {
	currentState *ConnectivityInfo

	locker *sync.RWMutex
	notify *notify.Notify
}

type NetManagerEventType uint
const (
	NetManagerConnectivityStateChanged NetManagerEventType = 1 << iota
	NetManagerConnectivityMeteringChanged
	NetManagerConnectivityBluetoothChanged
	NetManagerConnectivityNetTypeChanged
	NetManagerConnectivityCellularTypeChanged

	NetManagerConnectivityChanged = 0             |
		NetManagerConnectivityStateChanged        |
		NetManagerConnectivityMeteringChanged     |
		NetManagerConnectivityBluetoothChanged    |
		NetManagerConnectivityNetTypeChanged      |
		NetManagerConnectivityCellularTypeChanged
)

func NewNetManager(initialState *ConnectivityInfo) *NetManager {
	var locker sync.RWMutex
	return &NetManager{
		currentState: initialState,
		locker:       &locker,
		notify:       notify.New(&locker),
	}
}

// UpdateState update the current state of the Manager
func (m *NetManager) UpdateState(state *ConnectivityInfo) {
	m.locker.Lock()
	if m.currentState != state {
		m.currentState = state
		m.notify.Broadcast()
	}
	m.locker.Unlock()
}

// WaitForStateChange waits until the currentState changes from sourceState or ctx expires. A true value is returned in former case and false in latter.
func (m *NetManager) WaitForStateChange(ctx context.Context, sourceState *ConnectivityInfo, eventType NetManagerEventType) bool {
	m.locker.Lock()

	ok := true
	for ok {
		if (eventType & NetManagerConnectivityStateChanged != 0 && sourceState.State != m.currentState.State) ||
		   (eventType & NetManagerConnectivityMeteringChanged != 0 && sourceState.Metering != m.currentState.Metering) ||
		   (eventType & NetManagerConnectivityBluetoothChanged != 0 && sourceState.Bluetooth != m.currentState.Bluetooth) ||
		   (eventType & NetManagerConnectivityNetTypeChanged != 0 && sourceState.NetType != m.currentState.NetType) ||
		   (eventType & NetManagerConnectivityCellularTypeChanged != 0 && sourceState.CellularType != m.currentState.CellularType) {
			break
		}
		// wait until state has been changed or context has been cancel
		ok = m.notify.Wait(ctx)
	}

	m.locker.Unlock()
	return ok
}

// GetCurrentState return the current state of the Manager
func (m *NetManager) GetCurrentState() (state *ConnectivityInfo) {
	m.locker.RLock()
	state = m.currentState
	m.locker.RUnlock()
	return
}
