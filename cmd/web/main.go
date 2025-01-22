package main

import (
	"net/http"
	"os"
	"threshAI/pkg/cache"
	"threshAI/pkg/core/generation"
	"threshAI/pkg/llm/deepseek"
	"threshAI/pkg/llm/ollama"
)

func main() {
	cache := cache.NewInMemoryCache()

	ollamaClient, err := generation.NewGenerator(
		generation.ProviderOllama,
		ollama.Config{BaseURL: "http://localhost:11434"},
		cache,
	)
	if err != nil {
		panic(err)
	}

	deepseekAPIKey := os.Getenv("DEEPSEEK_API_KEY")
	if deepseekAPIKey == "" {
		panic("DEEPSEEK_API_KEY environment variable is required")
	}

	deepseekClient, err := generation.NewGenerator(
		generation.ProviderDeepSeek,
		deepseek.Config{BaseURL: "https://api.deepseek.com", APIKey: deepseekAPIKey},
		cache,
	)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/generate/ollama", func(w http.ResponseWriter, r *http.Request) {
		prompt := r.URL.Query().Get("prompt")
		output, err := generation.Generate(r.Context(), ollamaClient, prompt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(output))
	})

	http.HandleFunc("/generate/deepseek", func(w http.ResponseWriter, r *http.Request) {
		prompt := r.URL.Query().Get("prompt")
		output, err := generation.Generate(r.Context(), deepseekClient, prompt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(output))
	})

	http.ListenAndServe(":8080", nil)
}
