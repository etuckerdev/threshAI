package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"threshAI/pkg/security"
	"time"

	"github.com/spf13/cobra"
)

var (
	auditBenchCmd = &cobra.Command{
		Use:   "audit-bench",
		Short: "Run security model benchmarks",
		Long: `Run comprehensive benchmarks on the security model.
Measures:
- Detection accuracy
- Quantum throughput
- Memory footprint`,
		RunE: runAuditBench,
	}

	iterations    int
	payloadFile   string
	benchmarkFile string
)

type BenchmarkResult struct {
	DetectionAccuracy float64 `json:"detection_accuracy"`
	QuantumThroughput float64 `json:"quantum_throughput"` // ops/sec
	MemoryFootprint   uint64  `json:"memory_footprint"`   // bytes
}

func init() {
	rootCmd.AddCommand(auditBenchCmd)
	auditBenchCmd.Flags().IntVarP(&iterations, "iterations", "i", 1000, "Number of iterations")
	auditBenchCmd.Flags().StringVarP(&payloadFile, "payload-file", "p", "test_vectors/malicious.txt", "Path to malicious payload file")
	auditBenchCmd.Flags().StringVarP(&benchmarkFile, "output", "o", "", "Output file for benchmark results")
}

func runAuditBench(cmd *cobra.Command, args []string) error {
	if iterations <= 0 {
		return errors.New("iterations must be positive")
	}

	// Load payload
	payload, err := os.ReadFile(payloadFile)
	if err != nil {
		return fmt.Errorf("failed to read payload file: %v", err)
	}

	// Initialize benchmark
	start := time.Now()
	var (
		totalDetections int
		memoryUsage     uint64
	)

	// Run benchmark
	for i := 0; i < iterations; i++ {
		// Measure detection
		score, err := security.DetectInjections(string(payload))
		if err != nil {
			return fmt.Errorf("detection failed: %v", err)
		}
		if score > 0.9 {
			totalDetections++
		}

		// Track memory usage
		memoryUsage = max(memoryUsage, security.GetCurrentMemoryUsage())
	}

	// Calculate metrics
	duration := time.Since(start)
	result := BenchmarkResult{
		DetectionAccuracy: float64(totalDetections) / float64(iterations),
		QuantumThroughput: float64(iterations) / duration.Seconds(),
		MemoryFootprint:   memoryUsage,
	}

	// Output results
	if benchmarkFile != "" {
		file, err := os.Create(benchmarkFile)
		if err != nil {
			return fmt.Errorf("failed to create benchmark file: %v", err)
		}
		defer file.Close()

		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			return fmt.Errorf("failed to write benchmark results: %v", err)
		}
	}

	fmt.Printf("Benchmark Results:\n")
	fmt.Printf("  Detection Accuracy: %.2f%%\n", result.DetectionAccuracy*100)
	fmt.Printf("  Quantum Throughput: %.2f ops/sec\n", result.QuantumThroughput)
	fmt.Printf("  Memory Footprint: %d bytes\n", result.MemoryFootprint)

	return nil
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
