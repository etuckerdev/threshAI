package execution

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var (
	processDone = false
	execCommand = exec.Command
)

// QuantumCheck verifies GPU requirements and enforces execution timeout
func QuantumCheck() {
	if !HasNvidiaGPU() {
		panic("Quantum lock requires NVIDIA reality anchor")
	}

	time.AfterFunc(500*time.Millisecond, func() {
		if !processDone {
			syscall.Kill(os.Getpid(), syscall.SIGTERM)
		}
	})
}

// HasNvidiaGPU checks for NVIDIA GPU availability
func HasNvidiaGPU() bool {
	cmd := execCommand("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

// MarkProcessDone signals that the process has completed successfully
func MarkProcessDone() {
	processDone = true
}
