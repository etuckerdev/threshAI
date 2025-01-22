package nous

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"threshAI/internal/core"
)

var (
	ErrMemoryOvercommit = errors.New("GPU memory overcommit: MAX 4GB")
)

func StreamModelChunks(modelPath string) error {
	// Calculate chunk size as 80% of available free memory
	chunkSize := int(float64(core.GetFreeMem()) * 0.8)

	// Streaming logic with memory ceiling enforcement
	modelSize, err := totalModelSize(modelPath)
	if err != nil {
		return err
	}
	for currentOffset := 0; currentOffset < modelSize; currentOffset += chunkSize {
		if err := loadModelChunk(modelPath, currentOffset, chunkSize); err != nil {
			return err
		}


// NuclearSanitization performs extreme input sanitization for sensitive operations
func NuclearSanitization(input string) string {
	// Remove all non-ASCII characters
	clean := strings.Map(func(r rune) rune {
		if r > 127 {
			return -1
		}
		return r
	}, input)

	// Remove potentially dangerous patterns
	clean = regexp.MustCompile(`(?i)(<script|&#|\\x|%[0-9a-f]{2})`).ReplaceAllString(clean, "")
	
	// Truncate to safe length
	if len(clean) > 4096 {
		clean = clean[:4096]
	}
	return clean
}

func StreamModelChunks(modelPath string) error {
	// Calculate chunk size as 80% of available free memory
	chunkSize := int(float64(core.GetFreeMem()) * 0.8)

	// Streaming logic with memory ceiling enforcement
	modelSize, err := totalModelSize(modelPath)
	if err != nil {
		return err
	}
	
	for currentOffset := 0; currentOffset < modelSize; currentOffset += chunkSize {
		if err := loadModelChunk(modelPath, currentOffset, chunkSize); err != nil {
			return err
		}

		// Enforce memory ceiling after each chunk
		if core.GetFreeMem() < core.MAX_GPU_MEM {
			NuclearSanitization("")
			return ErrMemoryOvercommit
		}
	}
	
	return nil
}
func totalModelSize(modelPath string) (int, error) {
	info, err := os.Stat(modelPath)
	if err != nil {
		return 0, err
	}
	return int(info.Size()), nil
}

func loadModelChunk(modelPath string, currentOffset int, chunkSize int) error {
	// Implementation using parameters
	_, err := os.Stat(modelPath)
	if err != nil {
		return fmt.Errorf("model path error: %w", err)
	}

	// Add actual chunk loading logic here
	return nil
}
