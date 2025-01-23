package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Memory represents the chat memory system
type Memory struct {
	Interactions []Interaction
	filepath     string
}

// Interaction represents a single chat interaction
type Interaction struct {
	UserInput string `json:"user_input"`
	EidosResp string `json:"eidos_response"`
}

var defaultMemoryPath = filepath.Join(os.Getenv("HOME"), ".thresh/memory/chat_history.json")

// LoadMemory initializes or loads existing memory
func LoadMemory() *Memory {
	mem := &Memory{
		filepath: defaultMemoryPath,
	}

	// Ensure directory exists
	os.MkdirAll(filepath.Dir(defaultMemoryPath), 0755)

	// Try to load existing memory
	data, err := os.ReadFile(defaultMemoryPath)
	if err == nil {
		json.Unmarshal(data, &mem.Interactions)
	}

	return mem
}

// Save persists the memory to disk
func (m *Memory) Save() error {
	data, err := json.MarshalIndent(m.Interactions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filepath, data, 0644)
}

// AddInteraction stores a new interaction
func (m *Memory) AddInteraction(input, response string) {
	m.Interactions = append(m.Interactions, Interaction{
		UserInput: input,
		EidosResp: response,
	})
}

// RetrieveRelevantContext finds relevant past interactions
func (m *Memory) RetrieveRelevantContext(input string) []Interaction {
	var relevant []Interaction
	input = strings.ToLower(input)

	// Simple relevance matching - can be enhanced with more sophisticated matching
	for i := len(m.Interactions) - 1; i >= 0; i-- {
		interaction := m.Interactions[i]
		if strings.Contains(strings.ToLower(interaction.UserInput), input) ||
			strings.Contains(input, strings.ToLower(interaction.UserInput)) {
			relevant = append(relevant, interaction)
			if len(relevant) >= 3 { // Limit context to last 3 relevant interactions
				break
			}
		}
	}

	return relevant
}

// RetrieveLastInteraction gets the most recent interaction
func (m *Memory) RetrieveLastInteraction() (*Interaction, error) {
	if len(m.Interactions) == 0 {
		return nil, ErrNoHistory
	}
	return &m.Interactions[len(m.Interactions)-1], nil
}

// Clear removes all stored interactions
func (m *Memory) Clear() error {
	m.Interactions = nil
	return m.Save()
}
