package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Ollama struct {
		URL string `yaml:"url"`
	} `yaml:"ollama"`

	DeepSeek struct {
		APIKey         string  `yaml:"-"` // From environment variable
		BaseURL        string  `yaml:"base_url"`
		MaxTokens      int     `yaml:"max_tokens"`
		CacheTTL       string  `yaml:"cache_ttl"`
		RequestTimeout string  `yaml:"request_timeout"`
		MaxRetries     int     `yaml:"max_retries"`
		Temperature    float64 `yaml:"temperature"`
	} `yaml:"deepseek"`
}

type SecurityConfig struct {
	BrutalPresets map[int]struct {
		Quant           string  `yaml:"quant"`
		VramLimit       string  `yaml:"vram_limit"`
		Description     string  `yaml:"description"`
		ChaosMultiplier float64 `yaml:"chaos_multiplier"`
		Unstable        bool    `yaml:"unstable,omitempty"`
	} `yaml:"brutal_presets"`
	Validation []struct {
		Name string `yaml:"name"`
		Cmd  string `yaml:"cmd"`
		Max  string `yaml:"max,omitempty"`
	} `yaml:"validation"`
	Monitoring struct {
		ChaosMetrics []string `yaml:"chaos_metrics"`
	} `yaml:"monitoring"`
}

func LoadConfigFile(filename string) ([]byte, error) {
	filepath := fmt.Sprintf("config/%s", filename)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", filepath, err)
	}

	var config SecurityConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s: %v", filepath, err)
	}

	return data, nil
}
