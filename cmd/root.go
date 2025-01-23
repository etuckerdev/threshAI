package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	verbose    bool
	configPath string
	debug      bool
)

var rootCmd = &cobra.Command{
	Use:   "thresh",
	Short: "Thresh AI command line interface",
	Long: `A comprehensive command line interface for the Thresh AI system.
Provides functionality for prompt execution, configuration management,
and system monitoring.`,
	Version: "1.0.0",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			fmt.Println("Debug mode enabled")
		}
		if verbose {
			fmt.Println("Verbose mode enabled")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	// Enable shell completion generation
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	return rootCmd.Execute()
}

func init() {
	// Persistent flags available to all commands
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to config file")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")

	// Add command groups
	rootCmd.AddGroup(&cobra.Group{ID: "core", Title: "Core Commands:"})
	rootCmd.AddGroup(&cobra.Group{ID: "config", Title: "Configuration Commands:"})
	rootCmd.AddGroup(&cobra.Group{ID: "system", Title: "System Commands:"})

	// Add commands to appropriate groups
	rootCmd.AddCommand(promptCmd)
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(systemCmd)
	rootCmd.AddCommand(helpCmd)

	// Enable shell completion
	rootCmd.InitDefaultCompletionCmd()

	// Initialize configuration
	initConfig()
}

func initConfig() {
	if configPath != "" {
		// Use config file from the flag
		fmt.Printf("Using config file: %s\n", configPath)
	}

	// Default config initialization can go here
	// For now, we'll just print a message in debug mode
	if debug {
		fmt.Println("Initializing default configuration")
	}
}
