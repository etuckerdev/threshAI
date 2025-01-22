package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
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

func (m *Memory) RetrieveLastInteraction() (*Interaction, error) {
	if len(m.Interactions) == 0 {
		return nil, fmt.Errorf("no interactions found")
	}
	return &m.Interactions[len(m.Interactions)-1], nil
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
