package memory

import "errors"

// Common memory-related errors
var (
	ErrNoHistory      = errors.New("no interaction history found")
	ErrInvalidFormat  = errors.New("invalid memory format")
	ErrStorageFailure = errors.New("failed to store memory")
)
