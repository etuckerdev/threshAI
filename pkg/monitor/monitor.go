package monitor

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	"threshAI/pkg/analytics"
)

type ModelMetrics struct {
	// Memory metrics
	TotalVRAM   uint64
	UsedVRAM    uint64
	PeakVRAM    uint64
	CurrentVRAM uint64

	// Layer timing metrics
	LayerTiming map[string]time.Duration
	TotalTime   time.Duration

	// Throughput metrics
	TokensPerSec float64
	BatchSize    int

	// Additional metrics
	CoherenceScore     float64
	EntanglementFactor float64
	SlangRatio         float64
	PromptSHA256       string
	BrutalLevel        int
}

// NewModelMetrics initializes a new ModelMetrics instance
func NewModelMetrics() *ModelMetrics {
	return &ModelMetrics{
		LayerTiming: make(map[string]time.Duration),
	}
}

// AddLayerTime records timing for a specific layer
func (m *ModelMetrics) AddLayerTime(layerName string, duration time.Duration) {
	m.LayerTiming[layerName] = duration
	m.TotalTime += duration
}

// UpdateMemoryStats updates memory statistics
func (m *ModelMetrics) UpdateMemoryStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.CurrentVRAM = memStats.Alloc
	if m.CurrentVRAM > m.PeakVRAM {
		m.PeakVRAM = m.CurrentVRAM
	}
	m.UsedVRAM = memStats.Sys
	m.TotalVRAM = memStats.TotalAlloc
}

// CalculateTokensPerSec computes token throughput
func (m *ModelMetrics) CalculateTokensPerSec(tokensProcessed int) {
	if m.TotalTime > 0 {
		seconds := float64(m.TotalTime) / float64(time.Second)
		m.TokensPerSec = float64(tokensProcessed) / seconds
	}
}

// GetLayerBreakdown returns timing breakdown per layer
func (m *ModelMetrics) GetLayerBreakdown() map[string]float64 {
	breakdown := make(map[string]float64)
	for layer, duration := range m.LayerTiming {
		breakdown[layer] = float64(duration) / float64(m.TotalTime)
	}
	return breakdown
}

// GetTopBottlenecks identifies slowest layers
func (m *ModelMetrics) GetTopBottlenecks(n int) []string {
	type layerTime struct {
		name     string
		duration time.Duration
	}

	layers := make([]layerTime, 0, len(m.LayerTiming))
	for name, duration := range m.LayerTiming {
		layers = append(layers, layerTime{name, duration})
	}

	sort.Slice(layers, func(i, j int) bool {
		return layers[i].duration > layers[j].duration
	})

	result := make([]string, 0, n)
	for i := 0; i < n && i < len(layers); i++ {
		result = append(result, layers[i].name)
	}
	return result
}

func LogMetrics(data interface{}) {
	switch metrics := data.(type) {
	case *ModelMetrics:
		// Memory metrics
		fmt.Printf("[ThreshAI-Metrics] total_vram=%d MB\n", metrics.TotalVRAM/1024/1024)
		fmt.Printf("[ThreshAI-Metrics] used_vram=%d MB\n", metrics.UsedVRAM/1024/1024)
		fmt.Printf("[ThreshAI-Metrics] peak_vram=%d MB\n", metrics.PeakVRAM/1024/1024)

		// Performance metrics
		fmt.Printf("[ThreshAI-Metrics] tokens_per_sec=%.2f\n", metrics.TokensPerSec)
		fmt.Printf("[ThreshAI-Metrics] total_time=%s\n", metrics.TotalTime)

		// Layer timing breakdown
		breakdown := metrics.GetLayerBreakdown()
		for layer, percentage := range breakdown {
			fmt.Printf("[ThreshAI-Metrics] layer_%s_time_pct=%.2f%%\n",
				strings.ToLower(layer), percentage*100)
		}

		// Bottleneck analysis
		bottlenecks := metrics.GetTopBottlenecks(3)
		if len(bottlenecks) > 0 {
			fmt.Printf("[ThreshAI-Metrics] top_bottlenecks=%s\n",
				strings.Join(bottlenecks, ","))
		}

		// Additional metrics
		if metrics.CoherenceScore > 0 {
			fmt.Printf("[ThreshAI-Metrics] coherence_score=%.2f\n", metrics.CoherenceScore)
		}
		if metrics.EntanglementFactor > 0 {
			fmt.Printf("[ThreshAI-Metrics] entanglement_factor=%.2f\n", metrics.EntanglementFactor)
		}
		if metrics.SlangRatio > 0 {
			fmt.Printf("[ThreshAI-Metrics] slang_ratio=%.2f\n", metrics.SlangRatio)
		}
		if metrics.PromptSHA256 != "" {
			fmt.Printf("[ThreshAI-Metrics] prompt_sha256=%s\n", metrics.PromptSHA256)
		}
		if metrics.BrutalLevel > 0 {
			fmt.Printf("[ThreshAI-Metrics] brutal_level=%d\n", metrics.BrutalLevel)
		}

	case map[string]interface{}:
		// Legacy format support
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
