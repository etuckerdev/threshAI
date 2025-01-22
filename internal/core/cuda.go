package core

/*
#cgo CFLAGS: -I/usr/local/cuda-12.3/include
#cgo LDFLAGS: -L/usr/local/cuda-12.3/lib64 -lcudart
#include <cuda_runtime.h>
*/
import "C"
import (
	"errors"
)

var (
	cudaInitialized = false
)

func initCUDA() error {
	if cudaInitialized {
		return nil
	}
	if err := C.cudaSetDevice(0); err != C.cudaSuccess {
		return errors.New("failed to initialize CUDA device")
	}
	cudaInitialized = true
	return nil
}

func GetFreeMem() uint64 {
	if err := initCUDA(); err != nil {
		panic(err)
	}
	var free C.size_t
	var total C.size_t
	if err := C.cudaMemGetInfo(&free, &total); err != C.cudaSuccess {
		panic("failed to get CUDA memory info")
	}
	return uint64(free)
}

func FlushAllBuffers() {
	C.cudaDeviceSynchronize()
}

func ResetCache() {
	C.cudaDeviceReset()
}

func FreeOrphanedMemory() {
	C.cudaDeviceReset()
}

func CUDACheck() {
	if GetFreeMem() < 500*1024*1024 {
		TriggerQuantumRollback()
		panic("CUDA MEM CRISIS: 500MB remaining")
	}
}

func TriggerQuantumRollback() {
	FlushAllBuffers()
	ResetCache()
	FreeOrphanedMemory()
}
