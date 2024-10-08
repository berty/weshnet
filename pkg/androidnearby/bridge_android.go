//go:build android && !noproximitytransport
// +build android,!noproximitytransport

package androidnearby

import (
	"go.uber.org/zap"

	proximity "berty.tech/weshnet/v2/pkg/proximitytransport"
)

// Supported is used by main package as default value for enable this driver.
// While UI actually enable or not the Java Android Nearby driver.
// TODO: remove this when UI will be able to handle this for the first App launching.
const Supported = true

// Noop implementation for Android
// Real driver is given from Java directly here: berty/js/android/app/src/main/java/tech/berty/gobridge/nearby
func NewDriver(logger *zap.Logger) proximity.ProximityDriver {
	logger = logger.Named("Nearby")
	logger.Info("NewDriver(): Java driver not found")

	return proximity.NewNoopProximityDriver(ProtocolCode, ProtocolName, DefaultAddr)
}
