package cmd

import (
	"os"

	webcmd "threshAI/cmd/web"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "thresh",
	Short: "Thresh AI command line interface",
	Long:  `A command line interface for interacting with the Thresh AI system`,
}

func init() {
	// Add all commands
	AddCommands()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// AddCommands registers all the command modules
func AddCommands() {
	// Add web interface command
	rootCmd.AddCommand(webcmd.WebCmd)
	// Add generate command
	rootCmd.AddCommand(generateCmd)
}
