package cmd

import (
	"fmt"

	"github.com/etuckerdev/threshAI/internal/core/config"
	"github.com/etuckerdev/threshAI/internal/core/generator"
	"github.com/etuckerdev/threshAI/internal/nous"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [text]",
	Short: "Generate content using NOUSx heritage",
	Long: `Quantum content generation with optional brutal mode.
    Implements quantum execution constraints and model validation.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGenerate(cmd, args)
	},
}

func runGenerate(cmd *cobra.Command, args []string) error {
	brutalMode, _ := cmd.Flags().GetInt("brutal")
	securityModel, _ := cmd.Flags().GetString("security-model")
	quantize, _ := cmd.Flags().GetString("quantize")

	if brutalMode > 0 {
		fmt.Println("⚠️  WARNING: Mode checks bypassed (--brutal accepts any int now)")
		config.BrutalLevel = brutalMode
		config.AllowBrutal = true
		config.CurrentGenerationMode = config.ModeBrutal
	}

	config.SecurityModel = securityModel
	config.Quantize = quantize

	modelName := "cas/ministral-8b-instruct-2410_q4km"
	fmt.Printf("\nDEBUG: Model loaded? %v\n", nous.IsModelLoaded(modelName))

	output, err := generator.Generate(args[0])
	if err != nil {
		fmt.Printf("Error generating content: %v\n", err)
		return err
	}
	fmt.Println(output)
	return nil
}

func init() {
	// Initialize generate command flags
	generateCmd.Flags().IntVar(&config.BrutalLevel, "brutal", 0, "Brutality tier (0-3)")
	generateCmd.Flags().StringVar(&config.SecurityModel, "security-model", "cas/ministral-8b-instruct-2410_q4km", "Specify model override")
	generateCmd.Flags().StringVar(&config.Quantize, "quantize", "", "Quantization level (Q4_K_M, Q5_K_S)")
	generateCmd.Flags().StringVar(&config.Quantize, "quantize", "Q4_K_M", "Quantization level (Q4_K_M, Q5_K_S)")

	// Register generate command
	rootCmd.AddCommand(generateCmd)
}
