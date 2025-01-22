package cmd

import (
	"fmt"
	"time"

	"github.com/etuckerdev/threshAI/internal/core/config"
	"github.com/etuckerdev/threshAI/internal/render"
	"github.com/etuckerdev/threshAI/internal/security"
	"github.com/etuckerdev/threshAI/pkg/execution"
	pkgflags "github.com/etuckerdev/threshAI/pkg/flags"
	"github.com/etuckerdev/threshAI/pkg/models"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [text]",
	Short: "Generate content using NOUSx heritage",
	Long: `Quantum content generation with optional brutal mode.
Implements quantum execution constraints and model validation.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		text := args[0]
		quantum := pkgflags.Quantum()
		model, _ := cmd.Flags().GetString("model")
		quantize, _ := cmd.Flags().GetString("quantize")

		// Validate quantization level
		switch quantize {
		case "Q4_K_M", "Q5_K_S":
			// Valid quantization level
		default:
			fmt.Printf("Invalid quantization level: %s\n", quantize)
			return
		}

		// Validate and select model
		validatedModel := models.ValidateModel(model)

		// Apply quantum constraints if enabled
		if quantum {
			execution.QuantumCheck()
		}

		// Security scan
		securityModel, _ := cmd.Flags().GetString("security-model")
		if !models.IsValidSecurityModel(securityModel) {
			fmt.Printf("Invalid security model: %s\n", securityModel)
			return
		}

		// Perform quantum-resistant security scan
		threatScore, err := security.DetectInjections(text)
		if err != nil {
			fmt.Printf("Security scan failed: %v\n", err)
			return
		}
		if threatScore > 0.89 {
			security.NuclearIsolation(fmt.Sprintf(
				"Threat detected: score=%.2f using model %s",
				threatScore,
				securityModel,
			))
			return
		}

		// Apply temporal smearing for quantum-resistant timing protection
		smearedDuration := security.TemporalSmearing()
		time.Sleep(smearedDuration)

		// Process content based on brutal mode level
		var content string
		if config.BrutalLevel == int(config.ModeBrutal) {
			content = render.Brutalize(text)
		} else {
			content = text
		}

		// Generate output
		fmt.Printf("# Quantum Generation Report\n\n")
		fmt.Printf("## Execution Constraints\n")
		fmt.Printf("Model: %s\n", validatedModel)
		fmt.Printf("Brutal Mode: %v\n", config.BrutalLevel)
		fmt.Printf("Quantum Lock: %v\n\n", quantum)
		fmt.Printf("## Generated Content\n%s\n", content)
	},
}

func init() {
	// Initialize generate command flags
	// In cmd/generate.go
	generateCmd.Flags().IntVar(&config.BrutalLevel, "brutal", 0, "Brutal generation level")
	generateCmd.Flags().String("model", "", "Specify quantum model override")
	generateCmd.Flags().String("security-model", "withsecure/llama3-8b", "Specify quantum security model override")
	generateCmd.Flags().String("quantize", "Q4_K_M", "Quantization level (Q4_K_M, Q5_K_S)")

	// Register generate command
	rootCmd.AddCommand(generateCmd)
}
