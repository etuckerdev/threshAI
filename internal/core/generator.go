package core

import "fmt"

func GenerateContent(request ContentRequest) string {
	content, err := GenerateWithOllama(request)
	if err != nil {
		return fmt.Sprintf("⚠️ AI Error: %v\nFallback content...", err)
	}
	return content
}
