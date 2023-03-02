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

func (t EventType) Has(other EventType) bool {
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

// WaitForStateChange waits until the currentState changes from sourceState or ctx expires.
// The eventType argument allow you to filter out the event you want to wait for.
// A true value is returned in former case and false in latter.
// The EventType is also returned to know which events has been triggered.
func (m *NetManager) WaitForStateChange(ctx context.Context, sourceState *ConnectivityInfo, eventType EventType) (bool, EventType) {
	m.locker.Lock()

	var currentEventType EventType
	ok := true

	for ok {
		currentEventType = 0

		if (sourceState.State != m.currentState.State) {
			currentEventType |= ConnectivityStateChanged
		}
		if (sourceState.Metering != m.currentState.Metering) {
			currentEventType |= ConnectivityMeteringChanged
		}
		if (sourceState.Bluetooth != m.currentState.Bluetooth) {
			currentEventType |= ConnectivityBluetoothChanged
		}
		if (sourceState.NetType != m.currentState.NetType) {
			currentEventType |= ConnectivityNetTypeChanged
		}
		if (sourceState.CellularType != m.currentState.CellularType) {
			currentEventType |= ConnectivityCellularTypeChanged
		}

		if ((eventType & currentEventType) != 0) {
			break
		}
		// wait until state has been changed or context has been cancel
		ok = m.notify.Wait(ctx)
	}

	m.locker.Unlock()
	return ok, currentEventType
}

// GetCurrentState return the current state of the Manager
func (m *NetManager) GetCurrentState() (state ConnectivityInfo) {
	m.locker.RLock()
	state = m.currentState
	m.locker.RUnlock()
	return
}
