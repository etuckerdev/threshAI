package render

import (
	"testing"
)

func TestBrutalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Brutalization Integrity Check",
			input:    "# Header\n**Bold** text\n```code block```",
			expected: "text code block│\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Brutalize(tt.input)
			if result != tt.expected {
				t.Errorf("Brutalize() = %v, want %v", result, tt.expected)
			}
		})
	}
}
