package transformer

import (
	"context"
	"fmt"
)

type Adapter struct {
	model  *TransformerModel
	config Config
}

func NewAdapter(config Config) (*Adapter, error) {
	model, err := NewTransformerModel(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create transformer model: %v", err)
	}

	return &Adapter{
		model:  model,
		config: config,
	}, nil
}

func (a *Adapter) Generate(ctx context.Context, prompt string) (string, error) {
	// Tokenize the input
	tokens, err := a.model.tokenizer.Encode(prompt, a.config.MaxContext)
	if err != nil {
		return "", fmt.Errorf("tokenization failed: %v", err)
	}

	// Convert tokens to float64
	input := make([]float64, len(tokens))
	for i, t := range tokens {
		input[i] = float64(t)
	}

	// Calculate target length
	maxLen := a.config.MaxContext
	if len(input) >= maxLen {
		maxLen = len(input) + 100 // Generate 100 more tokens
	}

	// Generate tokens
	generated, err := a.model.Generate(input, maxLen)
	if err != nil {
		return "", fmt.Errorf("generation failed: %v", err)
	}

	// Convert generated tokens back to integers
	outputTokens := make([]int, len(generated)-len(input))
	for i, token := range generated[len(input):] {
		outputTokens[i] = int(token)
	}

	// Decode tokens to text
	result := a.model.tokenizer.Decode(outputTokens)

	return result, nil
}

// Save saves the model weights to a file
func (a *Adapter) Save(path string) error {
	return a.model.SaveCheckpoint(path)
}

// Load loads the model weights from a file
func (a *Adapter) Load(path string) error {
	model, err := LoadCheckpoint(path)
	if err != nil {
		return fmt.Errorf("failed to load checkpoint: %v", err)
	}
	a.model = model
	return nil
}

// RegisterProvider registers the transformer provider type
func RegisterProvider() {
	// TODO: Register the transformer provider in the generation package
	// This would require updating the generation package to support registration
}
