// File: internal/analytics/slang_detector.go
package analytics

import (
	"fmt"
	"strings"
)

var tiktokSlang = map[string]bool{
	"rizz":    true,
	"sus":     true,
	"bussin'": true,
	"no cap":  true,
	"cool":    true,
	"awesome": true,
	"yeet":    true,
	"lit":     true,
	"bro":     true,
	"vibe":    true,
	"flex":    true,
	"ghosted": true,
	"clout":   true,
	"drip":    true,
	"simp":    true,
	// Add more slang terms here
}

func CalculateSlangRatio(text string) float32 {
	words := strings.Fields(text)
	slangCount := 0

	for _, word := range words {
		if tiktokSlang[strings.ToLower(word)] {
			slangCount++
		}
	}

	fmt.Printf("DEBUG: Slang words detected: %d/%d\n", slangCount, len(words))
	return float32(slangCount) / float32(len(words))
}
