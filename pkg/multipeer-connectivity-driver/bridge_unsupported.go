//go:build !darwin || (darwin && !cgo) || catalyst || noproximitytransport
// +build !darwin darwin,!cgo catalyst noproximitytransport

package mc

import (
	"go.uber.org/zap"

	proximity "berty.tech/weshnet/pkg/proximitytransport"
)

const Supported = false

// Noop implementation for platform that are not Darwin

func NewDriver(logger *zap.Logger) proximity.ProximityDriver {
	logger = logger.Named("MC")
	logger.Info("NewDriver(): incompatible system")

	return proximity.NewNoopProximityDriver(ProtocolCode, ProtocolName, DefaultAddr)
}
