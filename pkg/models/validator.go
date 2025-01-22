package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	fallbackSequence   = []string{"nous-hermes2:10.7b", "mistral-r3:7b"}
	cachePath          = filepath.Join(os.Getenv("HOME"), ".thresh/cache/last_working_model")
	cacheTTL           = time.Hour
	validationEndpoint = "http://localhost:11434/api/tags"
)

type OllamaModel struct {
	Name       string `json:"name"`
	ModifiedAt string `json:"modified_at"`
}

func ValidateModel(target string) string {
	if IsModelCached() && !modelAvailable(target) {
		return ReadCacheModel()
	}

	models := FetchOllamaModels()
	if !contains(models, target) {
		target = FirstAvailable(fallbackSequence)
	}

	WriteCacheModel(target)
	return target
}

func IsModelCached() bool {
	info, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		return false
	}
	return time.Since(info.ModTime()) < cacheTTL
}

func ReadCacheModel() string {
	data, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return fallbackSequence[0]
	}
	return string(data)
}

func WriteCacheModel(model string) error {
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	return ioutil.WriteFile(cachePath, []byte(model), 0644)
}

func FetchOllamaModels() []string {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	var result struct {
		Models []OllamaModel `json:"models"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	models := make([]string, len(result.Models))
	for i, m := range result.Models {
		models[i] = m.Name
	}
	return models
}

func contains(models []string, target string) bool {
	for _, m := range models {
		if m == target {
			return true
		}
	}
	return false
}

func FirstAvailable(models []string) string {
	available := FetchOllamaModels()
	for _, m := range models {
		if contains(available, m) {
			return m
		}
	}
	return models[0]
}

func modelAvailable(target string) bool {
	return contains(FetchOllamaModels(), target)
}

func IsValidSecurityModel(model string) bool {
	approvedModels := []string{
		"withsecure/llama3-8b",
		"nous-hermes2:10.7b-secure",
	}
	return contains(approvedModels, model)
}
