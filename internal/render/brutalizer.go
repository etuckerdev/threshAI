package render

import (
	"regexp"
	"strings"
)

var markdownRegex = regexp.MustCompile(`(?s)(\x60{3}(?:.|\n)*?\x60{3}|\*\*.*?\*\*|^#+\s+)`)

// Brutalize strips markdown formatting and enforces line constraints
func Brutalize(content string) string {
	// Regex pattern execution
	stripped := markdownRegex.ReplaceAllString(content, "")

	// UTF-8 pipe wrapping
	lines := strings.Split(stripped, "\n")
	var builder strings.Builder
	for _, line := range lines {
		for len(line) > 0 {
			chunk := line
			if len(chunk) > 80 {
				chunk = chunk[:80]
			}
			builder.WriteString(chunk + "│\n")
			line = line[len(chunk):]
		}
	}
	return builder.String()
}
