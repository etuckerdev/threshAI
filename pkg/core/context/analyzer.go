package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Message represents a chat message with metadata
type Message struct {
	Content   string    `json:"content"`
	Role      string    `json:"role"`
	Language  string    `json:"language,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// Context maintains the session state and language preferences
type Context struct {
	UserID     string    `json:"user_id"`
	Messages   []Message `json:"messages"`
	LastActive time.Time `json:"last_active"`
	mu         sync.RWMutex
	MLModel    string `json:"ml_model"`
}

// contextStore stores active user contexts
var (
	store   = make(map[string]*Context)
	storeMu sync.RWMutex
)

// OllamaResponse represents the response from Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
}

// DetectLanguage analyzes the prompt and conversation history to determine
// the programming language context. Uses Nous-Hermes2 model for inference.
func DetectLanguage(prompt string, history []Message) string {
	// Build context from history
	var contextBuilder strings.Builder
	for _, msg := range history {
		if msg.Language != "" {
			contextBuilder.WriteString(msg.Language + " ")
		}
	}

	// Prepare prompt for language detection
	request := map[string]interface{}{
		"model": "nous-hermes2:10.7b-ctx",
		"prompt": fmt.Sprintf(
			`Given the coding conversation context and new prompt, detect the main programming language.
Previous context: %s
New prompt: %s
Respond with ONLY the programming language name in lowercase.`,
			contextBuilder.String(),
			prompt,
		),
		"stream": false,
	}

	// Call Ollama API
	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(jsonEncode(request)),
	)
	if err != nil {
		return "unknown" // Fallback
	}
	defer resp.Body.Close()

	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "unknown"
	}

	// If response is empty, default to a sensible language
	if result.Response == "" {
		return "javascript" // Default to JavaScript if no language detected
	}

	// Clean and normalize the detected language
	detectedLang := strings.TrimSpace(result.Response)
	return strings.ToLower(detectedLang)
}

// MaintainSessionState retrieves or creates a new context for the given user
func MaintainSessionState(userID string) *Context {
	storeMu.Lock()
	defer storeMu.Unlock()

	if ctx, exists := store[userID]; exists {
		ctx.LastActive = time.Now()
		return ctx
	}

	// Create new context
	newCtx := &Context{
		UserID:     userID,
		Messages:   make([]Message, 0),
		LastActive: time.Now(),
		MLModel:    "nous-hermes2:10.7b-ctx",
	}
	store[userID] = newCtx
	return newCtx
}

// AddMessage adds a new message to the context
func (c *Context) AddMessage(content, role string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	msg := Message{
		Content:   content,
		Role:      role,
		Language:  DetectLanguage(content, c.Messages),
		Timestamp: time.Now(),
	}
	c.Messages = append(c.Messages, msg)
	c.LastActive = time.Now()
}

// GetRecentMessages returns the n most recent messages
func (c *Context) GetRecentMessages(n int) []Message {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Messages) <= n {
		return c.Messages
	}
	return c.Messages[len(c.Messages)-n:]
}

func jsonEncode(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return []byte("{}")
	}
	return data
}
