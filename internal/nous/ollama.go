package nous

import "fmt"

// generateBrutalFunc holds the function to be used for brutal generation.
var generateBrutalFunc func(string, int) (string, error)

// GenerateBrutal generates content based on the given prompt using the brutal method.
func GenerateBrutal(prompt string, tier int) (string, error) {
	if generateBrutalFunc == nil {
		return "", fmt.Errorf("generateBrutalFunc is not set")
	}
	return generateBrutalFunc(prompt, tier)
}

// SetGenerateBrutalMock sets a mock function for GenerateBrutal.
func SetGenerateBrutalMock(mock func(string, int) (string, error)) {
	generateBrutalFunc = mock
}

// init initializes generateBrutalFunc with the default brutal generation logic.
func init() {
	generateBrutalFunc = func(prompt string, tier int) (string, error) {
		// Replace this with your actual brutal generation logic
		return fmt.Sprintf("Brutal Generated: %s", prompt), nil
	}
}
