package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type InMemoryCache struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]string),
	}
}

func (c *InMemoryCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, ok := c.data[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("key not found")
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	return nil
}
