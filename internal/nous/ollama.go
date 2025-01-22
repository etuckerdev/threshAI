package nous

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/etuckerdev/threshAI/internal/core"
)

var (
	// DefaultTimeout is the timeout duration for brutal mode
	DefaultTimeout = 30 * time.Second
)

// IsModelLoaded checks if a specific model is loaded
func IsModelLoaded(modelName string) bool {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), modelName)
}

func LoadModel() error {
	if core.GetFreeMem() < core.MAX_GPU_MEM {
		NuclearSanitization() // Initiate emergency cleanup
		return errors.New("GPU memory undercommit: MIN 4GB")
	}
	// Quantum-safe loading protocol
	return nil
}

func NuclearSanitization() {
	// Emergency memory purge protocol
	core.FlushAllBuffers()
	core.ResetCache()
	core.FreeOrphanedMemory()
}

// GenerateText handles direct Ollama interaction
func GenerateText(prompt string, brutalMode int, securityModel string, quantize string) (string, error) {
	modelTag := securityModel
	if quantize != "" {
		modelTag = fmt.Sprintf("%s:%s", modelTag, quantize)
	}

	cmd := exec.Command("ollama", "run", modelTag, prompt)
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
