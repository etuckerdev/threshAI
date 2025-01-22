// File: internal/analytics/metrics.go
package analytics

type MetricPayload struct {
	Coherence    float32
	Entanglement float32
	Slang        float32
	BrutalLevel  int
	PromptHash   string
}

// Deprecated: Use MetricPayload instead
type Metrics struct {
	Coherence    float32
	Entanglement float32
	SlangRatio   float32
}
