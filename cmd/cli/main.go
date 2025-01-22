package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"threshAI/pkg/cache"
	"threshAI/pkg/core/generation"
	"threshAI/pkg/llm/deepseek"
	"threshAI/pkg/llm/ollama"
)

func main() {
	prompt := flag.String("prompt", "", "Prompt for code generation")
	provider := flag.String("provider", "ollama", "LLM provider (ollama or deepseek)")
	flag.Parse()

	cache := cache.NewInMemoryCache()

	var generator generation.Generator
	var err error

	switch *provider {
	case "ollama":
		generator, err = generation.NewGenerator(
			generation.ProviderOllama,
			ollama.Config{BaseURL: "http://localhost:11434"},
			cache,
		)
	case "deepseek":
		deepseekAPIKey := os.Getenv("DEEPSEEK_API_KEY")
		if deepseekAPIKey == "" {
			fmt.Println("DEEPSEEK_API_KEY environment variable is required")
			return
		}
		generator, err = generation.NewGenerator(
			generation.ProviderDeepSeek,
			deepseek.Config{BaseURL: "https://api.deepseek.com", APIKey: deepseekAPIKey},
			cache,
		)
	default:
		fmt.Println("Invalid provider")
		return
	}

	if err != nil {
		fmt.Println("Error creating generator:", err)
		return
	}

	output, err := generation.Generate(context.Background(), generator, *prompt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(output)
}
