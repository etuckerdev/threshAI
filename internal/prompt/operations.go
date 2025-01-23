package prompt

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

// LoadPrompt loads and parses a prompt template from a file
func LoadPrompt(filename string) (*TaskPrompt, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var prompt TaskPrompt
	err = xml.Unmarshal(file, &prompt)
	if err != nil {
		return nil, fmt.Errorf("parsing XML: %w", err)
	}

	return &prompt, nil
}

// ExecuteChain executes a chain of prompts and returns the output
func ExecuteChain(prompts ...*TaskPrompt) string {
	output := "# Prompt Chain Output\n\n"

	for i, p := range prompts {
		output += formatPromptOutput(fmt.Sprintf("Prompt %d", i+1), p)
	}

	return output
}

// FormatPromptOutput formats a single prompt's output
func formatPromptOutput(name string, prompt *TaskPrompt) string {
	output := fmt.Sprintf("## %s\n", name)
	output += fmt.Sprintf("**System:** %s\n\n", prompt.System)

	output += fmt.Sprintf("**Inputs:**\n- RepoContent: %s\n- UserRequest: %s\n- FocusAreas: %s\n\n",
		prompt.Inputs.RepoContent, prompt.Inputs.UserRequest, prompt.Inputs.FocusAreas)

	output += fmt.Sprintf("**Process Flow:** %s\n\n", prompt.ProcessFlow)

	output += "**Output Template:**\n"
	output += fmt.Sprintf("- ProblemStatement: %s\n", prompt.OutputTemplate.Task.ProblemStatement)

	// Format CodeTargets
	output += "- CodeTargets:\n"
	for i, file := range prompt.OutputTemplate.Task.CodeTargets.Files {
		output += fmt.Sprintf("  - File: %s\n", file)
		if i < len(prompt.OutputTemplate.Task.CodeTargets.Functions) {
			output += fmt.Sprintf("    Function: %s\n", prompt.OutputTemplate.Task.CodeTargets.Functions[i])
		}
	}

	// Format SuccessMetrics
	output += "- SuccessMetrics:\n"
	for _, perf := range prompt.OutputTemplate.Task.SuccessMetrics.Performance {
		output += fmt.Sprintf("  - Performance: %s\n", perf)
	}
	for _, read := range prompt.OutputTemplate.Task.SuccessMetrics.Readability {
		output += fmt.Sprintf("  - Readability: %s\n", read)
	}
	output += "\n"

	return output
}

// ListPromptFiles lists all available prompt template files
func ListPromptFiles(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".xml" {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// ValidatePrompt checks if a prompt file is valid
func ValidatePrompt(filename string) error {
	_, err := LoadPrompt(filename)
	return err
}
