package ollama

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func GenerateBrutal(prompt string, tier int, maxTokens int) (string, error) {
	// Your custom brutal logic here
	model := fmt.Sprintf("ministral:8b-q%d_k_m", tier*2) // e.g., tier=1 → q2_k_m
	cmd := exec.Command("ollama", "run", model, "--prompt", prompt, "--max-tokens", strconv.Itoa(maxTokens))
	// ... execute and return output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute ollama command: %v, output: %s", err, output)
	}

	return string(output), nil
}

func Generate(prompt string, brutalMode int) (string, error) {
	// Use the correct model name and parameters
	model := "mistral:7b-instruct"
	cmd := exec.Command("ollama", "run", model, fmt.Sprintf("[INST] %s [/INST]", prompt))

	// Capture output properly
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ollama error: %v | Stderr: %s", err, stderr.String())
	}

	// Extract only the generated text (after the prompt)
	rawOutput := stdout.String()
	generatedText := strings.Split(rawOutput, "[/INST]")
	if len(generatedText) < 2 {
		return rawOutput, nil
	}
	return strings.TrimSpace(generatedText[1]), nil
}

// GenerateQuantumShards generates multiple outputs using different models
// to create diverse shards for quantum entanglement calculation
func GenerateQuantumShards(prompt string, count int, maxTokens int) ([]string, error) {
	models := []string{
		"mistral:7b-instruct-v0.2-q4_K_M",
		"llama2:13b-chat-q4_K_M",
		"codellama:7b-instruct-q4_K_M",
	}

	if count > len(models) {
		return nil, fmt.Errorf("maximum %d shards supported", len(models))
	}

	shards := make([]string, 0, count)
	for i := 0; i < count; i++ {
		cmd := exec.Command("ollama", "run", models[i], "--prompt", prompt, "--max-tokens", strconv.Itoa(maxTokens))
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to generate shard %d: %v", i+1, err)
		}
		shards = append(shards, string(output))
	}

	return shards, nil
}
