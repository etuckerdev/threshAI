package cmd

import (
	"fmt"
	"os" // Add this import
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "thresh",
	Short: "ThreshAI CLI",
	Long:  `ThreshAI - Your AI-powered content generator.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1) // Now properly imported
	}
}