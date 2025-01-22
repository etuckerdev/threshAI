// File: internal/monitor/monitor.go
package monitor

import (
	"fmt"

	"threshAI/internal/analytics"
)

func LogMetrics(data interface{}) {
	switch metrics := data.(type) {
	case map[string]interface{}:
		// New format with complete metrics payload
		fmt.Printf("[ThreshAI-Metrics] timestamp=%s\n", metrics["timestamp"])
		fmt.Printf("[ThreshAI-Metrics] coherence_score=%.2f\n", metrics["coherence_score"])
		fmt.Printf("[ThreshAI-Metrics] entanglement_factor=%.2f\n", metrics["entanglement_factor"])
		fmt.Printf("[ThreshAI-Metrics] slang_ratio=%.2f\n", metrics["slang_ratio"])
		fmt.Printf("[ThreshAI-Metrics] prompt_sha256=%s\n", metrics["prompt_sha256"])
		fmt.Printf("[ThreshAI-Metrics] brutal_level=%d\n", metrics["brutal_level"])
	case analytics.Metrics:
		// Backward compatibility
		fmt.Printf("[ThreshAI-Metrics] coherence_score=%.2f\n", metrics.Coherence)
		fmt.Printf("[ThreshAI-Metrics] entanglement=%.2f\n", metrics.Entanglement)
		fmt.Printf("[ThreshAI-Metrics] slang_ratio=%.2f\n", metrics.SlangRatio)
	default:
		fmt.Println("[ThreshAI-Metrics] Invalid metrics format")
	}
}
