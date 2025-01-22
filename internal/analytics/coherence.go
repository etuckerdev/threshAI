// File: internal/analytics/coherence.go
package analytics

import (
	"fmt"

	"threshAI/pkg/nlpvalidator"
)

func CalculateCoherence(text string) float32 {
	sentences := nlpvalidator.SplitSentences(text)
	fmt.Printf("DEBUG: Coherence input: %d sentences\n", len(sentences))

	if len(sentences) < 2 {
		return 1.0 // Perfect coherence for single sentences
	}

	if len(sentences) > 5 { // Safe guard for very long inputs
		sentences = sentences[:5]
		fmt.Println("WARN: Input text has many sentences, considering only first 5 for coherence calculation.")
	}

	similarityScore := nlpvalidator.CalculateTFIDFConsistency(sentences)
	penalty := nlpvalidator.DetectRepetition(text) * 0.2
	return similarityScore - penalty
}
