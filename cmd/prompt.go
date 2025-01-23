package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"threshAI/internal/prompt"

	"github.com/spf13/cobra"
)

var (
	outputFormat string
	promptFile   string
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Execute prompt operations",
	Long: `Run and manage prompts for task automation.
Supports XML-based prompt templates and chaining multiple prompts.`,
	GroupID: "core",
}

var promptRunCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Execute a prompt chain",
	Long: `Execute a prompt chain from an XML template file.
The template defines the system context, inputs, and expected outputs.`,
	Example: `thresh prompt run codecraft.xml
thresh prompt run oracle.xml --format json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 && promptFile == "" {
			return fmt.Errorf("prompt file is required")
		}

		file := promptFile
		if len(args) > 0 {
			file = args[0]
		}

		return executePrompt(file)
	},
}

var promptListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available prompts",
	Run: func(cmd *cobra.Command, args []string) {
		listPrompts()
	},
}

var promptValidateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a prompt template",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("prompt file is required")
		}
		return validatePrompt(args[0])
	},
}

func executePrompt(filename string) error {
	p, err := prompt.LoadPrompt(filename)
	if err != nil {
		return fmt.Errorf("failed to load prompt: %w", err)
	}

	fmt.Printf("✅ Loaded %s\n", filepath.Base(filename))
	fmt.Printf("System Context: %s\n", p.System)

	output := prompt.ExecuteChain(p)
	if outputFormat == "json" {
		// TODO: Implement JSON output formatting
	} else {
		fmt.Println(output)
	}

	return nil
}

func listPrompts() {
	promptsDir := "internal/core/memory/sys_prompts"
	files, err := os.ReadDir(promptsDir)
	if err != nil {
		fmt.Printf("Error reading prompts directory: %v\n", err)
		return
	}

	fmt.Println("Available prompts:")
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".xml" {
			fmt.Printf("- %s\n", file.Name())
		}
	}
}

func validatePrompt(filename string) error {
	if err := prompt.ValidatePrompt(filename); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	fmt.Printf("✅ %s is valid\n", filepath.Base(filename))
	return nil
}

func init() {
	promptRunCmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format (text/json)")
	promptRunCmd.Flags().StringVarP(&promptFile, "file", "i", "", "Input prompt file")

	promptCmd.AddCommand(promptRunCmd)
	promptCmd.AddCommand(promptListCmd)
	promptCmd.AddCommand(promptValidateCmd)
}
