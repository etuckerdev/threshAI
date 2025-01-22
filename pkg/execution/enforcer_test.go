package execution

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestQuantumCheck(t *testing.T) {
	// Test GPU detection
	t.Run("NVIDIA GPU Detection", func(t *testing.T) {
		// Mock command execution
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo", "50")
		}
		defer func() { execCommand = exec.Command }()

		if !HasNvidiaGPU() {
			t.Error("HasNvidiaGPU() should return true when GPU is present")
		}
	})

	t.Run("Quantum Timeout Enforcement", func(t *testing.T) {
		// Setup signal handler
		sig := make(chan os.Signal, 1)
		go func() {
			time.Sleep(600 * time.Millisecond)
			sig <- syscall.SIGTERM
		}()

		// Run quantum check and mark process as done immediately
		MarkProcessDone()
		go QuantumCheck()

		select {
		case <-sig:
			t.Error("QuantumCheck() should not trigger SIGTERM when process completes")
		case <-time.After(600 * time.Millisecond):
			// Expected no SIGTERM received
		}
	})
}
