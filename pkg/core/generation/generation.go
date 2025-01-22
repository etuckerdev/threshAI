package generation

import (
	"context"
	"fmt"

	"threshAI/pkg/cache"
	"threshAI/pkg/llm/deepseek"
	"threshAI/pkg/llm/ollama"
)

type Generator interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type ProviderType string

const (
	ProviderOllama   ProviderType = "ollama"
	ProviderDeepSeek ProviderType = "deepseek"
)

func NewGenerator(provider ProviderType, config interface{}, cache cache.Cache) (Generator, error) {
	switch provider {
	case ProviderOllama:
		cfg := config.(ollama.Config)
		return ollama.NewAdapter(cfg), nil
	case ProviderDeepSeek:
		cfg := config.(deepseek.Config)
		return deepseek.NewAdapter(cfg, cache), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}

func Generate(ctx context.Context, generator Generator, prompt string) (string, error) {
	return generator.Generate(ctx, prompt)
}
