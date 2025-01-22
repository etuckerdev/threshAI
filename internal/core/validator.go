package core

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"threshAI/internal/core/config"

	"gopkg.in/yaml.v2"
)

func ValidateMode(mode string) error {
	// Always validate brutal mode
	// Load brutal_presets from config
	securityConfig, err := loadSecurityConfig()
	if err != nil {
		return fmt.Errorf("failed to load security config: %v", err)
	}

	for tier, preset := range securityConfig.BrutalPresets {
		// Validate model for the selected tier
		if err := ValidateModel(preset.Quant); err != nil {
			return err
		}

		// Validate VRAM for the selected tier
		if err := ValidateVRAM(preset.VramLimit); err != nil {
			return err
		}
		fmt.Printf("Validated brutal mode for tier %d\n", tier)
	}
	return nil
}

func ValidateModel(quant string) error {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute ollama list: %v, output: %s", err, output)
	}

	modelPattern := fmt.Sprintf("8b-%s", quant)
	if !strings.Contains(string(output), modelPattern) {
		return fmt.Errorf("model with quantization %s not found in ollama list", quant)
	}

	return nil
}

func ValidateVRAM(vramLimit string) error {
	cmd := exec.Command("nvidia-smi", "--query-gpu=memory.used", "--format=noheader,nounits")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute nvidia-smi: %v, output: %s", err, output)
	}

	usedVRAMStr := strings.TrimSpace(string(output))
	usedVRAM, err := strconv.ParseFloat(usedVRAMStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse nvidia-smi output: %v, output: [%s]", err, usedVRAMStr)
	}

	limitBytes, err := parseVRAMLimit(vramLimit)
	if err != nil {
		return fmt.Errorf("failed to parse VRAM limit: %v", err)
	}

	if usedVRAM*1024*1024 > float64(limitBytes) {
		return fmt.Errorf("VRAM usage (%.2f MB) exceeds limit (%s)", usedVRAM, vramLimit)
	}

	return nil
}

func parseVRAMLimit(vramLimit string) (int64, error) {
	if strings.HasSuffix(vramLimit, "G") {
		limitGB, err := strconv.ParseFloat(strings.TrimSuffix(vramLimit, "G"), 64)
		if err != nil {
			return 0, err
		}
		return int64(limitGB * 1024 * 1024 * 1024), nil // Convert to bytes
	}

	return 0, fmt.Errorf("invalid VRAM limit format: %s", vramLimit)
}

// Add a helper function to load security.yaml
func loadSecurityConfig() (*config.SecurityConfig, error) {
	configFile, err := config.LoadConfigFile("security.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %v", err)
	}

	var securityConfig config.SecurityConfig
	err = yaml.Unmarshal(configFile, &securityConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal security config: %v", err)
	}

	return &securityConfig, nil
}
