//go:build (!darwin && !android) || noproximitytransport
// +build !darwin,!android noproximitytransport

package ble

import (
	"go.uber.org/zap"

	proximity "berty.tech/weshnet/v2/pkg/proximitytransport"
)

const Supported = false

// Noop implementation for platform that are not Darwin

func NewDriver(logger *zap.Logger) proximity.ProximityDriver {
	logger = logger.Named("BLE")
	logger.Info("NewDriver(): incompatible system")

	return proximity.NewNoopProximityDriver(ProtocolCode, ProtocolName, DefaultAddr)
}
