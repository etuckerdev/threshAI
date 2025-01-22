package nous

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/etuckerdev/threshAI/internal/core"
	"github.com/etuckerdev/threshAI/internal/core/config"
)

var (
	// DefaultTimeout is the timeout duration for brutal mode
	DefaultTimeout = 30 * time.Second
)

// Client represents the Ollama client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// IsModelLoaded checks if a specific model is loaded and available
func IsModelLoaded(modelName string) bool {
	// In brutal mode, assume model is loaded if it exists in Ollama
	if config.BrutalLevel > 0 {
		listCmd := exec.Command("ollama", "list")
		listCmd.Env = append(os.Environ(), "OLLAMA_HOST=localhost:11434")
		output, err := listCmd.CombinedOutput()
		return err == nil && strings.Contains(string(output), modelName)
	}

	// Standard mode requires full validation
	listCmd := exec.Command("ollama", "list")
	listCmd.Env = append(os.Environ(), "OLLAMA_HOST=localhost:11434")
	output, err := listCmd.CombinedOutput()
	if err != nil || !strings.Contains(string(output), modelName) {
		// If model not in list, try to pull it
		pullCmd := exec.Command("ollama", "pull", modelName)
		pullCmd.Env = append(os.Environ(), "OLLAMA_HOST=localhost:11434")
		if _, err := pullCmd.CombinedOutput(); err != nil {
			return false
		}
	}
	return true
}

func LoadModel() error {
	// Simple check for available memory
	if core.GetFreeMem() < core.MAX_GPU_MEM {
		return errors.New("insufficient GPU memory")
	}
	return nil
}

func NuclearSanitization() {
	// Emergency memory purge protocol
	core.FlushAllBuffers()
	core.ResetCache()
	core.FreeOrphanedMemory()
}

// GenerateText handles direct Ollama interaction
func GenerateText(ctx context.Context, prompt string, brutalMode int, securityModel string, quantize string) (string, error) {
	modelTag := securityModel
	if quantize != "" {
		modelTag = fmt.Sprintf("%s:%s", modelTag, quantize)
	}

	cmd := exec.CommandContext(ctx, "ollama", "run", "cas/ministral-8b-instruct-2410_q4km", prompt)

	cmd.Env = append(os.Environ(),
		"OLLAMA_HOST=127.0.0.1:11434",
		"OLLAMA_REGISTRY_AUTH_TOKEN=your_cas_token_here")
	out, err := cmd.CombinedOutput()
	if err != nil && len(out) == 0 {
		return "", fmt.Errorf("ollama execution failed: %w", err)
	}

	// Return any output we got, even if there was an error
	if len(out) > 0 {
		return strings.TrimSpace(string(out)), nil
	}

	return "", errors.New("no output received")
}
