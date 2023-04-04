package netmanager

import (
	"fmt"
	"strings"
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

func ParseConnectivityState(s string) (ConnectivityState, error) {
	switch strings.ToLower(s) {
	case "unknown":
		return ConnectivityStateUnknown, nil
	case "off":
		return ConnectivityStateOff, nil
	case "on":
		return ConnectivityStateOn, nil
	default:
		return ConnectivityStateUnknown, fmt.Errorf("invalid connectivity state (unknown/off/on): %q", s)
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

func ParseConnectivityNetType(s string) (ConnectivityNetType, error) {
	switch strings.ToLower(s) {
	case "unknown":
		return ConnectivityNetUnknown, nil
	case "none":
		return ConnectivityNetNone, nil
	case "wifi":
		return ConnectivityNetWifi, nil
	case "ethernet":
		return ConnectivityNetEthernet, nil
	case "cellular":
		return ConnectivityNetCellular, nil
	default:
		return ConnectivityNetUnknown, fmt.Errorf("invalid connectivity net type (unknown/none/wifi/ethernet/cellular): %q", s)
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

func ParseConnectivityCellularType(s string) (ConnectivityCellularType, error) {
	switch strings.ToLower(s) {
	case "unknown":
		return ConnectivityCellularUnknown, nil
	case "none":
		return ConnectivityCellularNone, nil
	case "2g":
		return ConnectivityCellular2G, nil
	case "3g":
		return ConnectivityCellular3G, nil
	case "4g":
		return ConnectivityCellular4G, nil
	case "5g":
		return ConnectivityCellular5G, nil
	default:
		return ConnectivityCellularUnknown, fmt.Errorf("invalid connectivity cellular type (unknown/none/2g/3g/4g/5g): %q", s)
	}
}

type ConnectivityInfo struct {
	// False when the device is not connected to a network.
	State ConnectivityState

	// True when the device is connected to a metered network.
	Metering ConnectivityState

	// True when the device is connected to a bluetooth network.
	Bluetooth ConnectivityState

	// The type of the network the device is connected to: wifi/ethernet/cellular.
	NetType ConnectivityNetType

	// If the device is connected to a cellular network:
	// The type of the cellular network the device is connected to: 2G/3G/4G/5G.
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
