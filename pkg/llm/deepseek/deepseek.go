package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"threshAI/pkg/cache"
	"threshAI/pkg/logging"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	Cache      cache.Cache
}

func NewClient(baseURL, apiKey string, cache cache.Cache) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		Cache:      cache,
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	// Check cache first
	if cached, err := c.Cache.Get(ctx, prompt); err == nil {
		logging.Logger.Printf("Cache hit for prompt: %s", prompt)
		return cached, nil
	}

	reqBody := Request{
		Model: "deepseek-chat",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/v1/chat/completions", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	logging.Logger.Printf("Making request to DeepSeek API with prompt: %s", prompt)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	output := response.Choices[0].Message.Content

	// Cache the result
	if err := c.Cache.Set(ctx, prompt, output, 24*time.Hour); err != nil {
		logging.Logger.Printf("Warning: failed to cache result: %v", err)
	}

	return output, nil
}
