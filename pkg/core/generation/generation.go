package generation

import (
"context"
"fmt"

"threshAI/pkg/cache"
"threshAI/pkg/llm/deepseek"
"threshAI/pkg/llm/ollama"
"threshAI/pkg/llm/transformer"
)

type Generator interface {
Generate(ctx context.Context, prompt string) (string, error)
}

type ProviderType string

const (
ProviderOllama      ProviderType = "ollama"
ProviderDeepSeek    ProviderType = "deepseek"
ProviderTransformer ProviderType = "transformer"
)

func NewGenerator(provider ProviderType, config interface{}, cache cache.Cache) (Generator, error) {
switch provider {
case ProviderOllama:
cfg := config.(ollama.Config)
return ollama.NewAdapter(cfg), nil
case ProviderDeepSeek:
cfg := config.(deepseek.Config)
return deepseek.NewAdapter(cfg, cache), nil
case ProviderTransformer:
cfg := config.(transformer.Config)
adapter, err := transformer.NewAdapter(cfg)
if err != nil {
return nil, fmt.Errorf("failed to create transformer adapter: %v", err)
}
return adapter, nil
default:
return nil, fmt.Errorf("unknown provider: %s", provider)
}
}

func Generate(ctx context.Context, generator Generator, prompt string) (string, error) {
return generator.Generate(ctx, prompt)
}
