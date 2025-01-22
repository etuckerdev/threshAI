package core

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"

	"threshAI/pkg/analytics"
	"threshAI/pkg/quantum"
	"threshAI/pkg/security"
)

var securityModel = security.SecurityConfig{}.ApprovedModels[0]

func GetCoherenceScore() float64 {
	return float64(analytics.CalculateCoherence(""))
}

func GetEntanglementFactor() float64 {
	return float64(quantum.CalculateEntanglement([]string{}))
}

func GetSlangRatio() float64 {
	return float64(analytics.CalculateSlangRatio(""))
}

type BrutalConfig struct {
	Quantization    string
	VRAMLimit       string
	Description     string
	ChaosMultiplier float64
	Unstable        bool
}

var brutalPresets = map[int]BrutalConfig{
	1: {
		Quantization:    "Q2_K",
		VRAMLimit:       "4G",
		Description:     "Gremlin Mode - 2-bit chaos",
		ChaosMultiplier: 0.3,
	},
	2: {
		Quantization:    "Q4_K_M",
		VRAMLimit:       "6G",
		Description:     "Standard Brutalization",
		ChaosMultiplier: 1.0,
	},
	3: {
		Quantization:    "Q6_K",
		VRAMLimit:       "8G",
		Description:     "Maximum Chaos",
		ChaosMultiplier: 2.5,
		Unstable:        true,
	},
}

func GenerateBrutal(brutalMode int, prompt string, maxTokens int) {
	preset, exists := brutalPresets[brutalMode]
	if !exists {
		return
	}

	quantize := preset.Quantization
	modelName := fmt.Sprintf("%s:%s-%s", securityModel, "8b", quantize)
	cmd := exec.Command("ollama", "run", modelName, "--prompt", prompt, "--max-tokens", strconv.Itoa(maxTokens))
	cmd.Stdout = io.Discard // Suppress output
	cmd.Stderr = io.Discard
	go cmd.Run() // Fire and forget
}

func ValidateBrutalMode(brutalMode int) bool {
	_, exists := brutalPresets[brutalMode]
	return exists
}

func GetChaosMetrics() map[string]float64 {
	return map[string]float64{
		"output_coherence_score":      GetCoherenceScore(),
		"quantum_entanglement_factor": GetEntanglementFactor(),
		"tiktok_slang_ratio":          GetSlangRatio(),
	}
}
