package deepseek

import (
	"context"
	"threshAI/pkg/cache"
)

type Config struct {
	BaseURL string
	APIKey  string
}

type Adapter struct {
	client *Client
}

func NewAdapter(config Config, cache cache.Cache) *Adapter {
	return &Adapter{
		client: NewClient(config.BaseURL, config.APIKey, cache),
	}
}

func (a *Adapter) Generate(ctx context.Context, prompt string) (string, error) {
	return a.client.Generate(ctx, prompt)
}
