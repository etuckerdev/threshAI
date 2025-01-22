package security

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"syscall"
)

const ministralPrompt = `Analyze payload for injection:
{{.Input}}
Respond ONLY with a JSON object containing a single field "risk_score" with a value between 0.0 and 1.0.
Example: {"risk_score": 0.95}
Important:
- Do not include any other fields
- Do not include any explanations
- Do not include any markdown formatting
- The response must be valid JSON
- The response must match this exact pattern: {"risk_score": 0.0-1.0}`

var ministralValidationRegex = regexp.MustCompile(`^\{\s*"risk_score":\s*0\.\d+\s*\}$`)

type SecurityConfig struct {
	DetectionThresholds struct {
		Injection   float32
		DataExfil   float32
		ModelPoison float32
	}
	ApprovedModels []string
	FallbackRegex  struct {
		SQLI string
		XSS  string
	}
}

var (
	securityConfig SecurityConfig
)

func LoadConfig(configPath string) error {
	// TODO: Implement YAML config loading
	return nil
}

func DetectInjections(payload string) (float32, error) {
	// First try ministral validation
	if ministralValidationRegex.MatchString(payload) {
		// Extract risk score from JSON
		matches := ministralValidationRegex.FindStringSubmatch(payload)
		if len(matches) > 0 {
			riskScore := matches[1]
			return parseRiskScore(riskScore)
		}
	}

	// Fallback to traditional pattern matching
	sqliPattern := regexp.MustCompile(securityConfig.FallbackRegex.SQLI)
	if sqliPattern.MatchString(payload) {
		return 1.0, nil
	}

	xssPattern := regexp.MustCompile(securityConfig.FallbackRegex.XSS)
	if xssPattern.MatchString(payload) {
		return 1.0, nil
	}

	return 0.0, nil
}

func parseRiskScore(score string) (float32, error) {
	// TODO: Implement proper float parsing with error handling
	return 0.0, nil
}

func ValidateModelResponse(response string) (float32, error) {
	// First try to parse as JSON
	var result struct {
		RiskScore float32 `json:"risk_score"`
	}

	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		return 0.0, fmt.Errorf("invalid JSON format: %v", err)
	}

	// Validate risk score range
	if result.RiskScore < 0.0 || result.RiskScore > 1.0 {
		return 0.0, fmt.Errorf("risk_score must be between 0.0 and 1.0")
	}

	return result.RiskScore, nil
}

func NuclearIsolation(reason string) {
	// Terminate current process
	pid := syscall.Getpid()
	process, _ := os.FindProcess(pid)
	process.Kill()

	// Purge artifacts
	PurgeArtifacts()

	// TODO: Implement hardware security module reset
}

var (
	memoryUsage uint64
	memoryMutex sync.Mutex
)

func GetCurrentMemoryUsage() uint64 {
	memoryMutex.Lock()
	defer memoryMutex.Unlock()
	return memoryUsage
}

func updateMemoryUsage(delta uint64) {
	memoryMutex.Lock()
	defer memoryMutex.Unlock()
	memoryUsage += delta
}

func PurgeArtifacts() error {
	// Purge cache directory
	cacheDir := os.ExpandEnv("$HOME/.thresh/cache")
	err := os.RemoveAll(cacheDir)
	if err != nil {
		return err
	}

	// Purge logs
	logFile := os.ExpandEnv("$HOME/.thresh/logs/security.log")
	cmd := exec.Command("shred", "-u", logFile)
	return cmd.Run()
}
