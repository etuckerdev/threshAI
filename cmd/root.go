package cmd

import (
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "thresh",
	Short: "threshAI Quantum Assimilation Protocol",
	Long: `Unified AI crucible merging Eidos and NOUSx remnants.
Implements quantum cognition with strict performance dogma.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		// Initialize flags through centralized package
		rootCmd.PersistentFlags().Bool("brutal", false, "Global brutal mode")

		// Add commands
		rootCmd.AddCommand(generateCmd)
		rootCmd.AddCommand(executeCmd)
		rootCmd.AddCommand(nukeCmd)
	})
}
