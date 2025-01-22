//go:build !cuda
// +build !cuda

package core

import "fmt"

func CheckCUDA() error {
	fmt.Println("CUDA support disabled - running in CPU-only mode")
	return nil
}

func GetFreeMem() uint64 {
	return 0
}

func FlushAllBuffers() {
	// No-op for non-CUDA
}

func ResetCache() {
	// No-op for non-CUDA
}

func FreeOrphanedMemory() {
	// No-op for non-CUDA
}
