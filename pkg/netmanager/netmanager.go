package netmanager

import (
	"context"
	"sync"

	"berty.tech/weshnet/internal/notify"
)

type NetManager struct {
	currentState ConnectivityInfo

	locker *sync.RWMutex
	notify *notify.Notify
}

type EventType uint

const (
	ConnectivityStateChanged EventType = 1 << iota
	ConnectivityMeteringChanged
	ConnectivityBluetoothChanged
	ConnectivityNetTypeChanged
	ConnectivityCellularTypeChanged

	ConnectivityChanged = 0 |
		ConnectivityStateChanged |
		ConnectivityMeteringChanged |
		ConnectivityBluetoothChanged |
		ConnectivityNetTypeChanged |
		ConnectivityCellularTypeChanged
)

func (t EventType) has(other EventType) bool {
	return (t & other) == other
}

func NewNetManager(initialState ConnectivityInfo) *NetManager {
	var locker sync.RWMutex
	return &NetManager{
		currentState: initialState,
		locker:       &locker,
		notify:       notify.New(&locker),
	}
}

// UpdateState update the current state of the Manager
func (m *NetManager) UpdateState(state ConnectivityInfo) {
	m.locker.Lock()
	if m.currentState != state {
		m.currentState = state
		m.notify.Broadcast()
	}
	m.locker.Unlock()
}

// WaitForStateChange waits until the currentState changes from sourceState or ctx expires. A true value is returned in former case and false in latter.
func (m *NetManager) WaitForStateChange(ctx context.Context, sourceState *ConnectivityInfo, eventType EventType) bool {
	m.locker.Lock()

	ok := true
	for ok {
		if (eventType.has(ConnectivityStateChanged) && sourceState.State != m.currentState.State) ||
			(eventType.has(ConnectivityMeteringChanged) && sourceState.Metering != m.currentState.Metering) ||
			(eventType.has(ConnectivityBluetoothChanged) && sourceState.Bluetooth != m.currentState.Bluetooth) ||
			(eventType.has(ConnectivityNetTypeChanged) && sourceState.NetType != m.currentState.NetType) ||
			(eventType.has(ConnectivityCellularTypeChanged) && sourceState.CellularType != m.currentState.CellularType) {
			break
		}
		// wait until state has been changed or context has been cancel
		ok = m.notify.Wait(ctx)
	}

	m.locker.Unlock()
	return ok
}

// GetCurrentState return the current state of the Manager
func (m *NetManager) GetCurrentState() (state ConnectivityInfo) {
	m.locker.RLock()
	state = m.currentState
	m.locker.RUnlock()
	return
}
