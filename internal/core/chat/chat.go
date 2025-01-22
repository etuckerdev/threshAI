package chat

import (
	"fmt"
	"strings"

	"threshAI/memory"
)

func LoadMemory() *memory.Memory {
	return memory.LoadMemory()
}

func NeedsClarification(userInput string) (bool, string) {
	// Detect ambiguous terms using NLP or regex
	ambiguousTerms := []string{"this", "that", "it", "they"}
	for _, term := range ambiguousTerms {
		if strings.Contains(strings.ToLower(userInput), term) {
			return true, fmt.Sprintf("You mentioned '%s'. Could you specify what you're referring to?", term)
		}
	}
	return false, ""
}

func GetUserInput() string {
	fmt.Print("User: ")
	var userInput string
	fmt.Scanln(&userInput)
	return userInput
}

func GenerateResponse(userInput string, context []memory.Interaction) string {
	// Dummy response generator for now
	if len(context) > 0 {
		return "Acknowledged. Context: " + context[0].UserInput
	}
	return "Response to: " + userInput
}

func shouldExit(userInput string) bool {
	return strings.ToLower(userInput) == "exit" || strings.ToLower(userInput) == "quit"
}

func ChatLoop() {
	mem := LoadMemory()
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
				fmt.Printf("Eidos: In our last conversation, you asked: %s\n", lastInteraction.UserInput)
				fmt.Printf("Eidos: My response was: %s\n", lastInteraction.EidosResp)
			}
			continue
		}

		// Generate response
		response := GenerateResponse(userInput, mem.RetrieveRelevantContext(userInput))
		fmt.Printf("Eidos: %s\n", response)
		mem.AddInteraction(userInput, response)
	}
}
