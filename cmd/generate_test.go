package cmd

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strings"
	"testing"

	"threshAI/internal/core/generator"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

// Helper function to execute the command and capture output
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

// Helper function to execute the command and capture output with context
func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

// Helper function to reset command state before each test
func resetState() {
	brutalMode = 0
	quantumMode = false
	metrics = false
}

func TestGenerateCmdValidSubcommands(t *testing.T) {
	resetState()

	tests := []struct {
		name        string
		args        []string
		outputMatch string
		wantErr     bool
	}{
		{
			name:        "valid_meme_generation",
			args:        []string{"generate", "meme", "--input", "test.txt"},
			outputMatch: `Meme Generation Complete:.*SHA3-256:[\da-f]{64}`,
			wantErr:     false,
		},
		{
			name:        "vector_with_quantum_smearing",
			args:        []string{"generate", "vector", "--quantum"},
			outputMatch: `11D Vector:.*smear_factor=[\d\.]+`,
			wantErr:     false,
		},
		{
			name:        "crisis_mode_validation",
			args:        []string{"generate", "crisis", "--severity", "5"},
			outputMatch: `CRISIS PROTOCOL ACTIVATED: LEVEL 5`,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set mock functions
			// Setup mock generators
			generator.SetGenerateMock(func(prompt string) (string, error) {
				switch {
				case strings.Contains(prompt, "meme"):
					return fmt.Sprintf("Meme Generation Complete: %s\n// CodeHash: %x", prompt, sha256.Sum256([]byte(prompt))), nil
				case strings.Contains(prompt, "vector"):
					return fmt.Sprintf("11D Vector: %s smear_factor=%.2f", prompt, rand.Float64()), nil
				case strings.Contains(prompt, "crisis"):
					return fmt.Sprintf("CRISIS PROTOCOL ACTIVATED: LEVEL %s", strings.TrimPrefix(prompt, "crisis")), nil
				default:
					return "", fmt.Errorf("invalid generator prompt")
				}
			})

			// Reset mock after test
			defer generator.SetGenerateMock(nil)

			output, err := executeCommand(rootCmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			matched, _ := regexp.MatchString(tt.outputMatch, output)
			if !matched {
				t.Errorf("executeCommand() output = %v, want match %v", output, tt.outputMatch)
			}
		})
	}
}

func TestGenerateCmdInvalid(t *testing.T) {
	resetState()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "generate no subcommand",
			args:    []string{"generate"},
			wantErr: true,
		},
		{
			name:    "generate invalid subcommand",
			args:    []string{"generate", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeCommand(rootCmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateCmdBrutalFlag(t *testing.T) {
	resetState()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "generate meme --brutal invalid",
			args:    []string{"generate", "meme", "--brutal", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeCommand(rootCmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateCmdQuantumFlag(t *testing.T) {
	resetState()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "generate vector --quantum invalid",
			args:    []string{"generate", "vector", "--quantum", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeCommand(rootCmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
