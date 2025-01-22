package telemetry

import (
	"threshAI/internal/analytics"
)

// Submit sends metrics to the telemetry system
func Submit(metrics analytics.MetricPayload) error {
	// TODO: Implement actual telemetry submission
	// For now just log the metrics
	return nil
}
