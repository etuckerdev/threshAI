//go:build cuda
// +build cuda

package core

/*
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -lcudart
#include <cuda_runtime.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
)

func CheckCUDA() error {
	// Initialize CUDA driver
	err := C.cudaSetDevice(0)
	if err != C.cudaSuccess {
		return fmt.Errorf("failed to initialize CUDA: %v", err)
	}

	// Get free and total memory
	var free, total C.size_t
	err = C.cudaMemGetInfo(&free, &total)
	if err != C.cudaSuccess {
		return fmt.Errorf("failed to get CUDA memory info: %v", err)
	}

	// Print memory info
	fmt.Printf("Total CUDA memory: %d bytes\n", total)
	fmt.Printf("Free CUDA memory: %d bytes\n", free)

	// Synchronize to ensure all CUDA operations are complete
	err = C.cudaDeviceSynchronize()
	if err != C.cudaSuccess {
		return fmt.Errorf("failed to synchronize CUDA device: %v", err)
	}

	// Reset device
	err = C.cudaDeviceReset()
	if err != C.cudaSuccess {
		return fmt.Errorf("failed to reset CUDA device: %v", err)
	}

	return nil
}

func GetFreeMem() uint64 {
	var free C.size_t
	var total C.size_t

	err := C.cudaMemGetInfo(&free, &total)
	if err != C.cudaSuccess {
		fmt.Printf("failed to get CUDA memory info: %v\n", err)
		return 0
	}
	return uint64(free)
}

func FlushAllBuffers() {
	// CUDA-specific buffer flushing
}

func ResetCache() {
	// CUDA-specific cache reset
}

func FreeOrphanedMemory() {
	// CUDA-specific memory cleanup
}
