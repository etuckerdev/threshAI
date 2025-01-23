package cmd

import (
	"fmt"
	"time"

	"threshAI/pkg/analytics"
	"threshAI/pkg/core/utils"
	"threshAI/pkg/license"
	"threshAI/pkg/llm/ollama"
	"threshAI/pkg/monitor"
	"threshAI/pkg/quantum"
	"threshAI/pkg/telemetry"

	"github.com/spf13/cobra"
)

var (
	brutalMode  int
	quantumMode bool
	metrics     bool
)

func init() {
	generateCmd.Flags().IntVar(&brutalMode, "brutal", 0, "Brutality tier (0 = off, 1-3 = severity)")
	generateCmd.Flags().BoolVar(&quantumMode, "quantum", false, "Enable quantum consensus mode")
	generateCmd.Flags().BoolVar(&metrics, "metrics", false, "Enable metrics collection and logging")
}

var generateCmd = &cobra.Command{
	Use:   "generate [prompt]",
	Short: "Generate content using AI",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if brutalMode > 0 || quantumMode {
			if !license.HasBrutalAccess() {
				return fmt.Errorf("brutal/quantum mode requires a valid license (run 'thresh auth --upgrade')")
			}
			fmt.Printf("⚠️  Brutal tier %d engaged – validation bypassed\n", brutalMode)
			return nil
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("DEBUG: RunE - brutalMode value: %d\n", brutalMode)
		prompt := args[0]

		var output string
		var shardOutputs []string
		var err error

		// Create Ollama client with default configuration
		client := ollama.NewOllamaAdapter(ollama.OllamaConfig{
			BaseURL:         "http://localhost:11434",
			Model:           "codellama:70b",
			MaxTokens:       512,
			CacheTTL:        5 * time.Minute,
			RequestTimeout:  30 * time.Second,
			MaxRetries:      3,
			Temperature:     0.7,
			CacheEnabled:    true,
			FallbackEnabled: false,
		}, nil) // No cache for now

		if quantumMode {
			// Generate multiple shards for quantum entanglement calculation
			shardOutputs, err = client.GenerateQuantumShards(prompt, 3, 512) // Generate 3 shards with max 512 tokens
			if err != nil {
				return fmt.Errorf("quantum generation failed: %v", err)
			}
			output = shardOutputs[0] // Use first shard as primary output
			fmt.Println(output)
		} else if brutalMode > 0 {
			// Call the brutal generation logic
			output, err := client.GenerateBrutal(prompt, brutalMode, 512)
			if err != nil {
				return fmt.Errorf("brutal generation failed: %v", err)
			}
			fmt.Println(output)
		} else {
			// Standard generation using updated Ollama client
			output, err = client.Generate(prompt, brutalMode)
			if err != nil {
				return fmt.Errorf("error generating content: %v", err)
			}
			fmt.Println("Generated:", output)
			if metrics {
				fmt.Printf("[ThreshAI-Metrics] slang_ratio=%.2f\n", analytics.CalculateSlangRatio(output))
			}
		}

		if metrics {
			// After generating output:
			// Calculate metrics payload
			metrics := analytics.MetricPayload{
				Coherence:    analytics.CalculateCoherence(output),
				Entanglement: 0.0, // Default value
				Slang:        analytics.CalculateSlangRatio(output),
				BrutalLevel:  brutalMode,
				PromptHash:   utils.HashPrompt(prompt),
			}

			// Only calculate entanglement if we have multiple shards
			if len(shardOutputs) > 1 {
				metrics.Entanglement = quantum.CalculateEntanglement(shardOutputs)
			}

			// Submit metrics to telemetry
			if err := telemetry.Submit(metrics); err != nil {
				fmt.Printf("WARNING: Failed to submit metrics: %v\n", err)
			}

			// Log metrics in JSON format
			logEntry := map[string]interface{}{
				"timestamp":           time.Now().UTC().Format(time.RFC3339),
				"coherence_score":     metrics.Coherence,
				"entanglement_factor": metrics.Entanglement,
				"slang_ratio":         metrics.Slang,
				"prompt_sha256":       metrics.PromptHash,
				"brutal_level":        metrics.BrutalLevel,
			}
			monitor.LogMetrics(logEntry)
		}

		return nil
	},
}
