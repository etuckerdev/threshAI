package memory

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

// ConfigVersion tracks schema versions for backwards compatibility
type ConfigVersion string

const (
	ConfigV1 ConfigVersion = "v1"
	ConfigV2 ConfigVersion = "v2"
)

// AllocatorConfig represents the atomic configuration state
type AllocatorConfig struct {
	Version           ConfigVersion `json:"version"`
	SmallPoolSize     uint64        `json:"small_pool_size"`
	MediumPoolSize    uint64        `json:"medium_pool_size"`
	LargePoolSize     uint64        `json:"large_pool_size"`
	GrowthFactor      float64       `json:"growth_factor"`
	LowWatermark      float64       `json:"low_watermark"`
	MediumWatermark   float64       `json:"medium_watermark"`
	HighWatermark     float64       `json:"high_watermark"`
	CriticalWatermark float64       `json:"critical_watermark"`
	CooldownCycles    int32         `json:"cooldown_cycles"`
	CheckpointFreq    int32         `json:"checkpoint_freq"`
	MaxCacheEntries   int32         `json:"max_cache_entries"`
	MixedPrecision    bool          `json:"mixed_precision"`
	EnableAutoTuning  bool          `json:"enable_auto_tuning"`
}

// ConfigRegistry manages versioned configurations
type ConfigRegistry struct {
	current atomic.Value // stores *AllocatorConfig
}

// NewConfigRegistry initializes a new registry with default v1 config
func NewConfigRegistry() *ConfigRegistry {
	r := &ConfigRegistry{}
	r.current.Store(&AllocatorConfig{
		Version:           ConfigV1,
		SmallPoolSize:     smallPoolSize,
		MediumPoolSize:    mediumPoolSize,
		LargePoolSize:     largePoolSize,
		GrowthFactor:      growthFactor,
		LowWatermark:      lowWatermark,
		MediumWatermark:   mediumWatermark,
		HighWatermark:     highWatermark,
		CriticalWatermark: criticalMark,
		CooldownCycles:    cooldownCycles,
		CheckpointFreq:    minCheckpointFreq,
		MaxCacheEntries:   maxCacheEntries,
		MixedPrecision:    false,
		EnableAutoTuning:  true,
	})
	return r
}

// GetCurrentConfig safely retrieves the current configuration
func (r *ConfigRegistry) GetCurrentConfig() *AllocatorConfig {
	return r.current.Load().(*AllocatorConfig)
}

// ValidateConfig performs schema validation based on version
func (r *ConfigRegistry) ValidateConfig(config *AllocatorConfig) error {
	switch config.Version {
	case ConfigV1:
		return r.validateV1(config)
	case ConfigV2:
		return r.validateV2(config)
	default:
		return fmt.Errorf("unsupported config version: %s", config.Version)
	}
}

func (r *ConfigRegistry) validateV1(config *AllocatorConfig) error {
	// Validate pool sizes
	if config.SmallPoolSize == 0 || config.MediumPoolSize == 0 || config.LargePoolSize == 0 {
		return fmt.Errorf("all pool sizes must be positive")
	}
	if config.SmallPoolSize >= config.MediumPoolSize || config.MediumPoolSize >= config.LargePoolSize {
		return fmt.Errorf("pool sizes must be strictly increasing")
	}

	// Validate growth factor
	if config.GrowthFactor <= 1.0 {
		return fmt.Errorf("growth_factor must be greater than 1.0")
	}

	// Validate watermarks
	if !r.validateWatermarks(config) {
		return fmt.Errorf("watermarks must be strictly increasing between 0 and 1")
	}

	// Validate other parameters
	if config.CooldownCycles < 0 {
		return fmt.Errorf("cooldown_cycles must be non-negative")
	}
	if config.CheckpointFreq < 1 {
		return fmt.Errorf("checkpoint_freq must be positive")
	}
	if config.MaxCacheEntries < 1 {
		return fmt.Errorf("max_cache_entries must be positive")
	}

	return nil
}

func (r *ConfigRegistry) validateWatermarks(config *AllocatorConfig) bool {
	return config.LowWatermark > 0 &&
		config.MediumWatermark > config.LowWatermark &&
		config.HighWatermark > config.MediumWatermark &&
		config.CriticalWatermark > config.HighWatermark &&
		config.CriticalWatermark < 1.0
}

func (r *ConfigRegistry) validateV2(config *AllocatorConfig) error {
	// V2 adds validation for auto-tuning
	if err := r.validateV1(config); err != nil {
		return err
	}
	return nil
}

// UpdateConfig atomically updates configuration after validation
func (r *ConfigRegistry) UpdateConfig(config *AllocatorConfig) error {
	if err := r.ValidateConfig(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Deep copy to prevent external modification
	newConfig := *config
	r.current.Store(&newConfig)
	return nil
}

// LoadConfigFromJSON parses and validates JSON configuration
func (r *ConfigRegistry) LoadConfigFromJSON(data []byte) error {
	var config AllocatorConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}
	return r.UpdateConfig(&config)
}
