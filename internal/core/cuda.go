package core

/*
#cgo CFLAGS: -I/usr/local/cuda-12.7/include
#cgo LDFLAGS: -L/usr/local/cuda-12.7/lib64 -lcudart
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
	// Placeholder for flushing buffers
}

func ResetCache() {
	// Placeholder for resetting cache
}

func FreeOrphanedMemory() {
	// Placeholder for freeing orphaned memory
}
