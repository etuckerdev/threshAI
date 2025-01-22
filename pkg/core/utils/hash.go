package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// hashPrompt generates a SHA-256 hash of the input prompt
func HashPrompt(prompt string) string {
	hash := sha256.Sum256([]byte(prompt))
	return hex.EncodeToString(hash[:])
}
