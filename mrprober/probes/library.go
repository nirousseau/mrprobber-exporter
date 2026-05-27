package probes

import "github.com/VictoriaMetrics/metrics"

// UnregisterAllMetrics flush the registry when on configuration changes
func UnregisterAllMetrics() {
	metrics.UnregisterAllMetrics()
}
