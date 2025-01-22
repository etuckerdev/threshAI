package chat

import (
	"fmt"
	"os"
	"strings"
	"time"

	"encoding/json"
)

type Interaction struct {
	UserInput string
	EidosResp string
	Timestamp string
}

type Memory struct {
	Interactions []Interaction
	Context      map[string]interface{}
}

func LoadMemory() *Memory {
	// Load from JSON file or database
	data, _ := os.ReadFile("memory.json")
	var mem Memory
	json.Unmarshal(data, &mem)
	return &mem
}

func (m *Memory) Save() {
	data, _ := json.Marshal(m)
	os.WriteFile("memory.json", data, 0644)
}

func (m *Memory) AddInteraction(userInput, eidosResp string) {
	m.Interactions = append(m.Interactions, Interaction{
		UserInput: userInput,
		EidosResp: eidosResp,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func (m *Memory) RetrieveRelevantContext(userInput string) []Interaction {
	// Use vector similarity search (e.g., OpenAI embeddings)
	// or keyword matching to find related interactions
	var relevant []Interaction
	for _, interaction := range m.Interactions {
		if strings.Contains(strings.ToLower(interaction.UserInput), strings.ToLower(userInput)) {
			relevant = append(relevant, interaction)
		}
	}
	return relevant
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

func GenerateResponse(userInput string, context []Interaction) string {
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
			lastInput, lastResp := mem.RetrieveLastInteraction()
			if lastInput == "" {
				fmt.Println("Eidos: No previous conversation found.")
			} else {
				fmt.Printf("Eidos: In our last conversation, you asked: %s\n", lastInput)
				fmt.Printf("Eidos: My response was: %s\n", lastResp)
			}
			continue
		}

		// Generate response
		response := GenerateResponse(userInput, mem.RetrieveRelevantContext(userInput))
		fmt.Printf("Eidos: %s\n", response)
		mem.AddInteraction(userInput, response)
	}
}
