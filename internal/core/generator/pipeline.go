package generator

import (
	"context"
	"errors"
	"fmt"

	"github.com/etuckerdev/threshAI/internal/core/config"
	"github.com/etuckerdev/threshAI/internal/nous"
	"github.com/etuckerdev/threshAI/internal/render"
)

type Brutalizer interface {
	Brutalize(input string) string
}

type brutalizer struct{}
type quantumizer struct{}

func (b *brutalizer) brutalizePayload(input string) string {
	return render.Brutalize(input)
}

func (q *quantumizer) quantumizePayload(input string) string {
	return render.Quantumize(input)
}

// GenerateWithOllama handles direct model interaction
func GenerateWithOllama(input string, stream bool) (string, error) {
	if err := nous.LoadModel(); err != nil {
		return "", fmt.Errorf("model load failed: %w", err)
	}

	// In brutal mode we skip streaming and validation
	if config.BrutalLevel > 0 {
		// Direct generation through Ollama
		ctx := context.Background()
		result, err := nous.GenerateText(ctx, input, config.BrutalLevel, config.SecurityModel, config.Quantize)
		if err != nil {
			return "", fmt.Errorf("brutal generation failed: %w", err)
		}
		return result, nil
	}

	return "", errors.New("standard generation not implemented")
}

func ProcessInput(mode config.GenerationMode, input string) (string, error) {
	// Nuclear option enforcement
	if config.BrutalLevel > 0 {
		// Nuclear bypass - force brutal mode
		b := &brutalizer{}
		brutalized := b.brutalizePayload(input)

		// Direct Ollama call in brutal mode - no streaming
		result, err := GenerateWithOllama(brutalized, false)
		if err != nil {
			return "", fmt.Errorf("ollama generation failed: %w", err)
		}
		return result, nil
	}

	// Original validation
	if mode == config.ModeBrutal && !config.AllowBrutal {
		return "", errors.New("nuclear option required: brutal mode needs --unsafe")
	}

	switch mode {
	case config.ModeBrutal:
		b := &brutalizer{}
		return b.brutalizePayload(input), nil
	case config.ModeQuantum:
		q := &quantumizer{}
		return q.quantumizePayload(input), nil
	default:
		return input, nil
	}
}

func Generate(input string) (string, error) {
	return ProcessInput(config.CurrentGenerationMode, input)
}
