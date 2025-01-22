package memory

import (
	"encoding/json"
	"os"
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
