package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"threshAI/internal/core/config"
)

type Stream struct {
	client *http.Client
	config *config.Config
}

func NewStream(client *http.Client, config *config.Config) *Stream {
	return &Stream{
		client: client,
		config: config,
	}
}

func (s *Stream) StreamResponse(ctx context.Context, prompt string) error {
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.Ollama.URL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		var v interface{}
		if err := decoder.Decode(&v); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error decoding response: %v", err)
		}
		fmt.Println(v)
	}

	return nil
}
