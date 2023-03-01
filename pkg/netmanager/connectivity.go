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

func (cs ConnectivityState) String() string {
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

func (cnt ConnectivityNetType) String() string {
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

func (cct ConnectivityCellularType) String() string {
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

func (ci ConnectivityInfo) String() string {
	return fmt.Sprint("ConnectivityInfo{ ",
		"State: ", ci.State.String(), ", ",
		"Metering: ", ci.Metering.String(), ", ",
		"Bluetooth: ", ci.Bluetooth.String(), ", ",
		"NetType: ", ci.NetType.String(), ", ",
		"CellularType: ", ci.CellularType.String(),
		" }")
}
