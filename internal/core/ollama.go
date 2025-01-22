package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type OllamaRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Options struct {
		Temperature float32 `json:"temperature"`
		MaxTokens   int     `json:"num_predict"`
	} `json:"options"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func GenerateWithOllama(request ContentRequest) (string, error) {
	ollamaReq := OllamaRequest{
		Model:  "deepseek-r1:7b",
		Prompt: fmt.Sprintf("Write a %d-word %s article about '%s'. Do not include any internal monologue, thinking process, or self-referential text.", request.Length, request.Tone, request.Topic),
		Options: struct {
			Temperature float32 `json:"temperature"`
			MaxTokens   int     `json:"num_predict"`
		}{
			Temperature: 0.7,
			MaxTokens:   request.Length * 5,
		},
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Ollama request: %v", err)
	}

	fmt.Println("🚀 Sending request to Ollama...") // Progress feedback

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ollama connection failed: %v", err)
	}
	defer resp.Body.Close()

	var fullResponse strings.Builder
	decoder := json.NewDecoder(resp.Body)
	for {
		var result OllamaResponse
		if err := decoder.Decode(&result); err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("failed to decode Ollama response: %v", err)
		}
		fullResponse.WriteString(result.Response)
		fmt.Print(".") // Progress dots
		if result.Done {
			break
		}
	}
	fmt.Println() // Newline after progress dots

	// Extract and log thoughts
	cleanedResponse, thoughts := extractThoughts(fullResponse.String())

	// Log thoughts to a file
	if err := logThoughts(thoughts, request.Topic, request.Tone); err != nil {
		return "", fmt.Errorf("failed to log thoughts: %v", err)
	}

	if len(cleanedResponse) == 0 {
		return "", fmt.Errorf("ollama returned an empty response")
	}

	return cleanedResponse, nil
}

func extractThoughts(content string) (string, string) {
	// Extract thoughts between <think> tags
	startTag := "<think>"
	endTag := "</think>"
	startIdx := strings.Index(content, startTag)
	endIdx := strings.Index(content, endTag)

	if startIdx == -1 || endIdx == -1 {
		return content, "" // No thoughts found
	}

	thoughts := content[startIdx+len(startTag) : endIdx]
	cleanedResponse := content[:startIdx] + content[endIdx+len(endTag):]

	return cleanedResponse, thoughts
}

func logThoughts(thoughts string, topic string, tone string) error {
	if thoughts == "" {
		return nil // No thoughts to log
	}

	file, err := os.OpenFile("thoughts.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add metadata
	logEntry := fmt.Sprintf(`
=== Thought Log ===
Topic: %s
Tone: %s
Timestamp: %s
Thoughts:
%s
`, topic, tone, time.Now().Format(time.RFC3339), thoughts)

	_, err = file.WriteString(logEntry)
	return err
}
