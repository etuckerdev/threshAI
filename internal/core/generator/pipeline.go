package generator

import "fmt"

// generateFunc holds the function to be used for generating content.
var generateFunc func(string) (string, error)

// Generate generates content based on the given prompt.
func Generate(prompt string) (string, error) {
	if generateFunc == nil {
		return "", fmt.Errorf("generateFunc is not set")
	}
	return generateFunc(prompt)
}

// SetGenerateMock sets a mock function for Generate.
func SetGenerateMock(mock func(string) (string, error)) {
	generateFunc = mock
}

// init initializes generateFunc with the default generation logic.
func init() {
	generateFunc = func(prompt string) (string, error) {
		// Replace this with your actual generation logic
		return fmt.Sprintf("Generated: %s", prompt), nil
	}
}
