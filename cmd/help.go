package cmd

import (
	"fmt"
	"strings"

	"threshAI/internal/help"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "Get detailed help and examples for commands",
	Long: `Enhanced help system providing detailed documentation, examples, and search capabilities.
Use 'thresh help [command]' for detailed information about a command.
Use 'thresh help search [query]' to search the help documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showGeneralHelp()
			return
		}

		if args[0] == "search" && len(args) > 1 {
			searchHelp(strings.Join(args[1:], " "))
			return
		}

		showCommandHelp(strings.Join(args, " "))
	},
}

var helpSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search help documentation",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchHelp(strings.Join(args, " "))
	},
}

func showGeneralHelp() {
	fmt.Println("ThreshAI CLI Help")
	fmt.Println("\nCommand Categories:")

	// Show core commands
	fmt.Println("\nCore Commands:")
	for _, cmd := range help.GetHelpByCategory("core") {
		fmt.Printf("  %-15s %s\n", cmd.Command, cmd.Description)
	}

	// Show config commands
	fmt.Println("\nConfiguration Commands:")
	for _, cmd := range help.GetHelpByCategory("config") {
		fmt.Printf("  %-15s %s\n", cmd.Command, cmd.Description)
	}

	// Show system commands
	fmt.Println("\nSystem Commands:")
	for _, cmd := range help.GetHelpByCategory("system") {
		fmt.Printf("  %-15s %s\n", cmd.Command, cmd.Description)
	}

	fmt.Println("\nUse 'thresh help [command]' for detailed information about a command")
	fmt.Println("Use 'thresh help search [query]' to search the help documentation")
}

func showCommandHelp(cmdPath string) {
	cmdHelp, err := help.GetCommandHelp(cmdPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Display command information
	fmt.Printf("\n%s - %s\n\n", cmdHelp.Command, cmdHelp.Description)
	fmt.Printf("Usage:\n  %s\n\n", cmdHelp.Usage)

	// Display subcommands if any
	if len(cmdHelp.SubCommands) > 0 {
		fmt.Println("Available Commands:")
		for _, sub := range cmdHelp.SubCommands {
			fmt.Printf("  %-15s %s\n", sub.Command, sub.Description)
		}
		fmt.Println()
	}

	// Display examples
	if len(cmdHelp.Examples) > 0 {
		fmt.Println("Examples:")
		for _, example := range cmdHelp.Examples {
			fmt.Printf("  %s\n", example)
		}
	}
}

func searchHelp(query string) {
	results := help.SearchHelp(query)
	if len(results) == 0 {
		fmt.Printf("No results found for '%s'\n", query)
		return
	}

	fmt.Printf("Search results for '%s':\n\n", query)
	for _, result := range results {
		fmt.Printf("Command: %s\n", result.Command)
		fmt.Printf("Description: %s\n", result.Description)
		if len(result.Examples) > 0 {
			fmt.Println("Example:", result.Examples[0])
		}
		fmt.Println()
	}
}

func init() {
	helpCmd.AddCommand(helpSearchCmd)
	rootCmd.AddCommand(helpCmd)
}
