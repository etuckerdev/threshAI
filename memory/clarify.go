package memory

import (
	"fmt"
	"strings"
)

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
