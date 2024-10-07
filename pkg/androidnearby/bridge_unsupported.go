//go:build !android || noproximitytransport
// +build !android noproximitytransport

package androidnearby

import (
	"go.uber.org/zap"

	proximity "berty.tech/weshnet/v2/pkg/proximitytransport"
)

const Supported = false

// Noop implementation for platform that are not Darwin

func NewDriver(logger *zap.Logger) proximity.ProximityDriver {
	logger = logger.Named("Nearby")
	logger.Info("NewDriver(): incompatible system")

	return proximity.NewNoopProximityDriver(ProtocolCode, ProtocolName, DefaultAddr)
}
