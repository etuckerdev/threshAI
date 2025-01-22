package deepseek

import (
	"threshAI/pkg/cache"
	"time"
)

type Builder struct {
	config Config
	cache  cache.Cache
}

func NewBuilder(baseURL, apiKey string, cache cache.Cache) *Builder {
	return &Builder{
		config: Config{
			BaseURL: baseURL,
			APIKey:  apiKey,
		},
		cache: cache,
	}
}

func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	return b
}

func (b *Builder) Build() *Client {
	return NewClient(b.config.BaseURL, b.config.APIKey, b.cache)
}
