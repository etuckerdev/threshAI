package render

import (
	"math/rand"
	"strings"
	"time"
)

// Quantumize applies quantum superposition effects to text
func Quantumize(content string) string {
	rand.Seed(time.Now().UnixNano())
	lines := strings.Split(content, "\n")
	var builder strings.Builder

	for _, line := range lines {
		// Apply quantum superposition by randomly flipping characters
		runes := []rune(line)
		for i := 0; i < len(runes); i++ {
			if rand.Intn(100) < 30 { // 30% chance to flip
				runes[i] = flipRune(runes[i])
			}
		}
		builder.WriteString(string(runes) + "\n")
	}

	return builder.String()
}

func flipRune(r rune) rune {
	// Basic ASCII flipping for demonstration
	if r >= 'a' && r <= 'z' {
		return 'z' - (r - 'a')
	}
	if r >= 'A' && r <= 'Z' {
		return 'Z' - (r - 'A')
	}
	return r
}
