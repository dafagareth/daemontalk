package templates

import "sync/atomic"

var radarActive atomic.Bool

func init() {
	radarActive.Store(true)
}

func SetRadarEnabled(enabled bool) {
	radarActive.Store(enabled)
}

func IsRadarEnabled() bool {
	return radarActive.Load()
}
