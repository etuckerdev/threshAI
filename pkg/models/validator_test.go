package models

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateModel(t *testing.T) {
	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"models":[{"name":"nous-hermes2:10.7b"}]}`))
	}))
	defer ts.Close()

	// Override API endpoint
	originalEndpoint := validationEndpoint
	validationEndpoint = ts.URL
	defer func() { validationEndpoint = originalEndpoint }()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Fallback Model Activation",
			input:    "invalid-model",
			expected: "nous-hermes2:10.7b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateModel(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateModel() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCacheOperations(t *testing.T) {
	// Setup test cache
	testCachePath := "test_cache"
	cachePath = testCachePath
	defer func() {
		os.Remove(testCachePath)
		cachePath = filepath.Join(os.Getenv("HOME"), ".thresh/cache/last_working_model")
	}()

	// Test cache write and read
	model := "test-model"
	err := WriteCacheModel(model)
	if err != nil {
		t.Fatalf("WriteCacheModel failed: %v", err)
	}

	cachedModel := ReadCacheModel()
	if cachedModel != model {
		t.Errorf("ReadCacheModel() = %v, want %v", cachedModel, model)
	}

	// Test cache expiration
	info, _ := os.Stat(testCachePath)
	oldTime := info.ModTime().Add(-2 * time.Hour)
	os.Chtimes(testCachePath, oldTime, oldTime)

	if IsModelCached() {
		t.Error("IsModelCached() should return false for expired cache")
	}
}
