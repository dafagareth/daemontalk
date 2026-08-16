package templates

import "sync/atomic"

var radarActive atomic.Bool

// SetRadarEnabled sets the global visibility state for the systems radar feature.
func SetRadarEnabled(enabled bool) {
	radarActive.Store(enabled)
}

// IsRadarEnabled reports whether the systems radar / graph feature is active.
func IsRadarEnabled() bool {
	return radarActive.Load()
}
