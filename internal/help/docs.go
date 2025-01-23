package help

import (
	"fmt"
	"strings"
)

// CommandHelp stores help information for a command
type CommandHelp struct {
	Command     string
	Usage       string
	Description string
	Examples    []string
	SubCommands []CommandHelp
	Flags       []FlagHelp
	Category    string
}

// FlagHelp stores help information for a flag
type FlagHelp struct {
	Name       string
	Shorthand  string
	Usage      string
	Category   string
	DefaultVal string
}

var commandHelp = map[string]CommandHelp{
	"prompt": {
		Command:     "prompt",
		Usage:       "thresh prompt [flags]",
		Description: "Execute prompt chains for task automation",
		Examples: []string{
			"thresh prompt run codecraft.xml",
			"thresh prompt run oracle.xml --verbose",
		},
		Category: "core",
	},
	"config": {
		Command:     "config",
		Usage:       "thresh config [command]",
		Description: "Manage ThreshAI configuration settings",
		Examples: []string{
			"thresh config show",
			"thresh config set model.endpoint http://localhost:11434",
		},
		SubCommands: []CommandHelp{
			{
				Command:     "show",
				Usage:       "thresh config show",
				Description: "Display current configuration",
			},
			{
				Command:     "set",
				Usage:       "thresh config set [key] [value]",
				Description: "Set a configuration value",
			},
		},
		Category: "config",
	},
	"system": {
		Command:     "system",
		Usage:       "thresh system [command]",
		Description: "System management and monitoring",
		Examples: []string{
			"thresh system status",
			"thresh system diag --verbose",
			"thresh system metrics",
		},
		SubCommands: []CommandHelp{
			{
				Command:     "status",
				Usage:       "thresh system status",
				Description: "Check system status",
			},
			{
				Command:     "diag",
				Usage:       "thresh system diag",
				Description: "Run system diagnostics",
			},
			{
				Command:     "metrics",
				Usage:       "thresh system metrics",
				Description: "Display system metrics",
			},
		},
		Category: "system",
	},
}

// GetCommandHelp retrieves help information for a command
func GetCommandHelp(cmdPath string) (CommandHelp, error) {
	parts := strings.Split(cmdPath, " ")
	cmd := parts[0]

	help, ok := commandHelp[cmd]
	if !ok {
		return CommandHelp{}, fmt.Errorf("no help found for command: %s", cmd)
	}

	// For subcommands
	if len(parts) > 1 {
		subCmd := parts[1]
		for _, sub := range help.SubCommands {
			if sub.Command == subCmd {
				return sub, nil
			}
		}
		return CommandHelp{}, fmt.Errorf("no help found for subcommand: %s", subCmd)
	}

	return help, nil
}

// GetCommandExamples retrieves examples for a command
func GetCommandExamples(cmdPath string) []string {
	help, err := GetCommandHelp(cmdPath)
	if err != nil {
		return nil
	}
	return help.Examples
}

// SearchHelp searches help documentation for keywords
func SearchHelp(query string) []CommandHelp {
	var results []CommandHelp
	query = strings.ToLower(query)

	for _, help := range commandHelp {
		// Search in command name and description
		if strings.Contains(strings.ToLower(help.Command), query) ||
			strings.Contains(strings.ToLower(help.Description), query) {
			results = append(results, help)
			continue
		}

		// Search in examples
		for _, example := range help.Examples {
			if strings.Contains(strings.ToLower(example), query) {
				results = append(results, help)
				break
			}
		}

		// Search in subcommands
		for _, sub := range help.SubCommands {
			if strings.Contains(strings.ToLower(sub.Command), query) ||
				strings.Contains(strings.ToLower(sub.Description), query) {
				results = append(results, help)
				break
			}
		}
	}

	return results
}

// GetHelpByCategory retrieves all commands in a category
func GetHelpByCategory(category string) []CommandHelp {
	var results []CommandHelp
	for _, help := range commandHelp {
		if help.Category == category {
			results = append(results, help)
		}
	}
	return results
}
