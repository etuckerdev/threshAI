// Package deepseek provides a client implementation for interacting with the DeepSeek API.
// It handles API requests, response parsing, and caching of results.
//
// Example:
//
//	client := deepseek.NewClient("https://api.deepseek.com", "api-key", cache)
//	response, err := client.Generate(context.Background(), "What is threshAI?")
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

// Client manages connections and requests to the DeepSeek API.
// Handles authentication, request building, and response processing.
//
// Fields:
//
//	BaseURL: The base URL for DeepSeek API endpoints
//	APIKey: Authentication key for API access
//	HTTPClient: Configured HTTP client with timeout settings
//	Cache: Cache implementation for storing API responses
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	Cache      cache.Cache
}

// NewClient creates a new DeepSeek API client instance.
// Validates configuration parameters and initializes HTTP client with timeout.
//
// Parameters:
//
//	baseURL: Base URL for DeepSeek API endpoints
//	apiKey: Authentication key for API access
//	cache: Cache implementation for storing API responses
//
// Returns:
//
//	*Client: Initialized DeepSeek client instance
func NewClient(baseURL, apiKey string, cache cache.Cache) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		Cache:      cache,
	}
}

// Message represents a single message in the chat conversation.
// Used for both input messages and API responses.
//
// Fields:
//
//	Role: The role of the message sender (e.g., "user", "assistant")
//	Content: The text content of the message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request represents the payload sent to the DeepSeek API.
// Contains the model identifier and conversation history.
//
// Fields:
//
//	Model: The model identifier to use for generation
//	Messages: Conversation history as a sequence of messages
type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Response represents the structure of API responses from DeepSeek.
// Contains generated message choices from the model.
//
// Fields:
//
//	Choices: Array of generated message options
type Response struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Generate sends a prompt to the DeepSeek API and returns the generated response.
// Implements caching to reduce API calls for repeated prompts.
//
// Parameters:
//
//	ctx: Context for request cancellation and timeout
//	prompt: Input text to send to the API
//
// Returns:
//
//	string: Generated response from DeepSeek
//	error: API request or processing errors
//
// Error Handling:
//   - Returns error for network failures
//   - Returns error for invalid API responses
//   - Returns error for empty responses
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
