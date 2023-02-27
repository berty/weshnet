package netmanager

import (
	"fmt"
)

type ConnectivityState int
const (
	ConnectivityStateUnknown ConnectivityState = iota
	ConnectivityStateOff
	ConnectivityStateOn
)
func (cs ConnectivityState) ToString() string {
	switch cs {
	case ConnectivityStateUnknown:
		return "unknown"
	case ConnectivityStateOff:
		return "off"
	case ConnectivityStateOn:
		return "on"
	default:
		return "error"
	}
}

type ConnectivityNetType int
const (
	ConnectivityNetUnknown ConnectivityNetType = iota
	ConnectivityNetNone
	ConnectivityNetWifi
	ConnectivityNetEthernet
	ConnectivityNetCellular
)
func (cnt ConnectivityNetType) ToString() string {
	switch cnt {
	case ConnectivityNetUnknown:
		return "unknown"
	case ConnectivityNetNone:
		return "none"
	case ConnectivityNetWifi:
		return "wifi"
	case ConnectivityNetEthernet:
		return "ethernet"
	case ConnectivityNetCellular:
		return "cellular"
	default:
		return "error"
	}
}

type ConnectivityCellularType int
const (
	ConnectivityCellularUnknown ConnectivityCellularType = iota
	ConnectivityCellularNone
	ConnectivityCellular2G
	ConnectivityCellular3G
	ConnectivityCellular4G
	ConnectivityCellular5G
)
func (cct ConnectivityCellularType) ToString() string {
	switch cct {
	case ConnectivityCellularUnknown:
		return "unknown"
	case ConnectivityCellularNone:
		return "none"
	case ConnectivityCellular2G:
		return "2G"
	case ConnectivityCellular3G:
		return "3G"
	case ConnectivityCellular4G:
		return "4G"
	case ConnectivityCellular5G:
		return "5G"
	default:
		return "error"
	}
}

type ConnectivityInfo struct {
	State        ConnectivityState
	Metering     ConnectivityState
	Bluetooth    ConnectivityState
	NetType      ConnectivityNetType
	CellularType ConnectivityCellularType
}

func NewConnectivityInfo() *ConnectivityInfo {
	return &ConnectivityInfo{
		State:        ConnectivityState(ConnectivityStateUnknown),
		Metering:     ConnectivityState(ConnectivityStateUnknown),
		Bluetooth:    ConnectivityState(ConnectivityStateUnknown),
		NetType:      ConnectivityNetType(ConnectivityNetUnknown),
		CellularType: ConnectivityCellularType(ConnectivityCellularUnknown),
	}
}

func (ci *ConnectivityInfo) GetState() ConnectivityState                 { return ci.State }
func (ci *ConnectivityInfo) GetMetering() ConnectivityState              { return ci.Metering }
func (ci *ConnectivityInfo) GetBluetooth() ConnectivityState             { return ci.Bluetooth }
func (ci *ConnectivityInfo) GetNetType() ConnectivityNetType             { return ci.NetType }
func (ci *ConnectivityInfo) GetCellularType() ConnectivityCellularType   { return ci.CellularType }

func (ci *ConnectivityInfo) SetState(state ConnectivityState)                      { ci.State = state }
func (ci *ConnectivityInfo) SetMetering(metering ConnectivityState)                { ci.Metering = metering }
func (ci *ConnectivityInfo) SetBluetooth(bluetooth ConnectivityState)              { ci.Bluetooth = bluetooth }
func (ci *ConnectivityInfo) SetNetType(netType ConnectivityNetType)                { ci.NetType = netType }
func (ci *ConnectivityInfo) SetCellularType(cellularType ConnectivityCellularType) { ci.CellularType = cellularType }

func (ci *ConnectivityInfo) ToString() string {
	return fmt.Sprint("ConnectivityInfo{ ",
		"State: ", ci.State.ToString(), ", ",
		"Metering: ", ci.Metering.ToString(), ", ",
		"Bluetooth: ", ci.Bluetooth.ToString(), ", ",
		"NetType: ", ci.NetType.ToString(), ", ",
		"CellularType: ", ci.CellularType.ToString(),
	" }")
}