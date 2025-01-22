package memory

import "strings"

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

func (m *Memory) RetrieveLastInteraction() (string, string) {
	if len(m.Interactions) == 0 {
		return "", "No previous interactions found."
	}
	last := m.Interactions[len(m.Interactions)-1]
	return last.UserInput, last.EidosResp
}
