package cmd

import (
	"fmt"
	"strings"

	"threshAI/memory"
	"threshAI/pkg/llm/ollama"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

var (
	ollamaConfig = ollama.OllamaConfig{
		BaseURL:    "http://localhost:11434",
		Model:      "llama2",
		MaxRetries: 3,
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
)

func GetUserInput() string {
	fmt.Print("User: ")
	var userInput string
	fmt.Scanln(&userInput)
	return userInput
}

func GenerateResponse(userInput string, context []memory.Interaction) string {
	// Build prompt with context
	prompt := userInput
	if len(context) > 0 {
		contextStr := "Previous relevant context:\n"
		for _, interaction := range context {
			contextStr += fmt.Sprintf("User: %s\nEidos: %s\n", interaction.UserInput, interaction.EidosResp)
		}
		prompt = contextStr + "\nCurrent query: " + userInput
	}

	// Generate response using ollama adapter
	adapter := ollama.NewOllamaAdapter(ollamaConfig, redisClient)
	output, err := adapter.Complete(prompt)
	if err != nil {
		return fmt.Sprintf("Oops, something went wrong generating a response. Error: %s", err)
	}

	if output == "" {
		return "Apologies, but I couldn't generate a proper response at the moment."
	}

	// Clean up any output formatting
	output = strings.TrimSpace(output)
	if strings.HasPrefix(output, "Generated:") {
		output = strings.TrimSpace(strings.TrimPrefix(output, "Generated:"))
	}

	return output
}

func shouldExit(userInput string) bool {
	return strings.ToLower(userInput) == "exit" || strings.ToLower(userInput) == "quit"
}

func ChatLoop() {
	mem := memory.LoadMemory()
	defer mem.Save()

	for {
		userInput := GetUserInput()
		if shouldExit(userInput) {
			break
		}

		// Handle "go back" requests
		if strings.Contains(strings.ToLower(userInput), "go back") {
			lastInteraction, err := mem.RetrieveLastInteraction()
			if err != nil {
				fmt.Println("Eidos: No previous conversation found.")
			} else {
				lastInput := lastInteraction.UserInput
				lastResp := lastInteraction.EidosResp
				fmt.Printf("Eidos: In our last conversation, you asked: %s\n", lastInput)
				fmt.Printf("Eidos: My response was: %s\n", lastResp)
			}
			continue
		}

		// Check for ambiguity
		if needsClarify, msg := memory.NeedsClarification(userInput); needsClarify {
			fmt.Printf("Eidos: %s\n", msg)
			userInput += " " + GetUserInput() // Append clarification
		}

		// Retrieve context
		context := mem.RetrieveRelevantContext(userInput)
		response := GenerateResponse(userInput, context)

		fmt.Printf("Eidos: %s\n", response)
		mem.AddInteraction(userInput, response)
	}
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat session with Eidos",
	Long:  `Initiates an interactive chat session with the Eidos personality engine, allowing for conversational interaction and persistent memory.`,
	Run: func(cmd *cobra.Command, args []string) {
		ChatLoop()
	},
}

func init() {
	// Add chat command to root command
	rootCmd.AddCommand(chatCmd)
}
