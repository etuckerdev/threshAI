package memory

import (
	"errors"
	"sync"
	"unsafe"
)

// Mock CUDA implementation
type cudaMock struct{}

const (
	mathModeMixed = 1
)

func (c *cudaMock) Malloc(size int) (unsafe.Pointer, error) {
	return unsafe.Pointer(&size), nil
}

func (c *cudaMock) SetMathMode(mode int) error {
	if mode != mathModeMixed {
		return errors.New("unsupported math mode")
	}
	return nil
}

func (c *cudaMock) EnableTensorCores(enable bool) error {
	return nil
}

func (c *cudaMock) Free(ptr unsafe.Pointer) error {
	return nil
}

var cuda = &cudaMock{}

const (
	bit4Alignment = 16 // 16-byte alignment for 4-bit tensors
)

type GPUMemoryAllocator struct {
	mu          sync.Mutex
	allocations map[uintptr]uint64 // track allocations by pointer and size
}

func NewGPUMemoryAllocator() *GPUMemoryAllocator {
	return &GPUMemoryAllocator{
		allocations: make(map[uintptr]uint64),
	}
}

// Alloc4BitTensorCache allocates aligned memory for 4-bit tensors
func (a *GPUMemoryAllocator) Alloc4BitTensorCache(size int) (unsafe.Pointer, error) {
	if size <= 0 {
		return nil, errors.New("invalid size: must be positive")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	// Calculate aligned size
	alignedSize := ((size + bit4Alignment - 1) / bit4Alignment) * bit4Alignment

	// Allocate GPU memory
	ptr, err := cuda.Malloc(alignedSize)
	if err != nil {
		return nil, err
	}

	// Track allocation
	a.allocations[uintptr(ptr)] = uint64(alignedSize)

	return ptr, nil
}

// EnableMixedPrecision enables mixed precision operations
func (a *GPUMemoryAllocator) EnableMixedPrecision() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Configure CUDA for mixed precision
	err := cuda.SetMathMode(mathModeMixed)
	if err != nil {
		return err
	}

	// Enable tensor cores
	return cuda.EnableTensorCores(true)
}

// FreeAll releases all allocated GPU memory
func (a *GPUMemoryAllocator) FreeAll() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	var lastErr error
	for ptr := range a.allocations {
		if err := cuda.Free(unsafe.Pointer(ptr)); err != nil {
			lastErr = err
		}
		delete(a.allocations, ptr)
	}

	return lastErr
}

// MemoryUsage returns total allocated GPU memory in bytes
func (a *GPUMemoryAllocator) MemoryUsage() uint64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	var total uint64
	for _, size := range a.allocations {
		total += size
	}
	return total
}
