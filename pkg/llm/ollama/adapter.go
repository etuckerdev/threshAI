package ollama

import (
	"context"
)

type Adapter struct {
	client *Client
}

func NewAdapter(config Config) *Adapter {
	return &Adapter{
		client: NewClient(config.BaseURL),
	}
}

func (a *Adapter) Generate(ctx context.Context, prompt string) (string, error) {
	return a.client.Generate(ctx, prompt)
}
