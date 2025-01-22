package core

import (
	"errors"

	"github.com/etuckerdev/threshAI/internal/core/config"
)

// internal/core/validator.go
func ValidateMode(mode config.GenerationMode) error {
	if mode == config.ModeBrutal && !config.AllowBrutal {
		return errors.New("brutal mode requires license activation")
	}
	return nil
}
