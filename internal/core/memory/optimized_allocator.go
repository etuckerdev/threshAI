package memory

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	// Memory pool tiers
	smallPoolSize  = 1 << 20 // 1MB for small allocations
	mediumPoolSize = 1 << 24 // 16MB for medium allocations
	largePoolSize  = 1 << 28 // 256MB for large allocations

	// Size thresholds
	smallAllocationLimit  = 1 << 16 // 64KB
	mediumAllocationLimit = 1 << 20 // 1MB

	// Dynamic scaling
	growthFactor = 1.3 // Reduced growth factor

	// Graduated pressure thresholds
	lowWatermark    = 0.50
	mediumWatermark = 0.75
	highWatermark   = 0.85
	criticalMark    = 0.95

	// Optimization parameters
	cooldownCycles    = 5 // Increased cooldown
	minCheckpointFreq = 100
	maxCacheEntries   = 1000 // Prevent unbounded cache growth
)

type MemoryPool struct {
	// Tiered memory pools
	small struct {
		blocks    []unsafe.Pointer
		sizes     []uint64
		freeList  []int
		totalSize uint64
		usedSize  atomic.Uint64
	}
	medium struct {
		blocks    []unsafe.Pointer
		sizes     []uint64
		freeList  []int
		totalSize uint64
		usedSize  atomic.Uint64
	}
	large struct {
		blocks    []unsafe.Pointer
		sizes     []uint64
		freeList  []int
		totalSize uint64
		usedSize  atomic.Uint64
	}
	mu         sync.RWMutex
	checkpoint struct {
		counter int32
		cooling int32
	}
}

func (p *MemoryPool) getTotalSize() uint64 {
	return p.small.totalSize + p.medium.totalSize + p.large.totalSize
}

func (p *MemoryPool) getTotalUsedSize() uint64 {
	return p.small.usedSize.Load() + p.medium.usedSize.Load() + p.large.usedSize.Load()
}

type poolTier struct {
	blocks    []unsafe.Pointer
	sizes     []uint64
	freeList  []int
	totalSize uint64
	usedSize  atomic.Uint64
}

// OptimizedAllocator implements adaptive memory management
type OptimizedAllocator struct {
	pool           *MemoryPool
	tensorCache    sync.Map // Cache for frequently used tensor sizes
	lastAllocation uint64   // Track allocation patterns
	configRegistry *ConfigRegistry
}

func NewOptimizedAllocator() *OptimizedAllocator {
	registry := NewConfigRegistry()
	config := registry.GetCurrentConfig()

	pool := &MemoryPool{
		small: struct {
			blocks    []unsafe.Pointer
			sizes     []uint64
			freeList  []int
			totalSize uint64
			usedSize  atomic.Uint64
		}{
			blocks:    make([]unsafe.Pointer, 0),
			sizes:     make([]uint64, 0),
			freeList:  make([]int, 0),
			totalSize: config.SmallPoolSize,
		},
		medium: struct {
			blocks    []unsafe.Pointer
			sizes     []uint64
			freeList  []int
			totalSize uint64
			usedSize  atomic.Uint64
		}{
			blocks:    make([]unsafe.Pointer, 0),
			sizes:     make([]uint64, 0),
			freeList:  make([]int, 0),
			totalSize: config.MediumPoolSize,
		},
		large: struct {
			blocks    []unsafe.Pointer
			sizes     []uint64
			freeList  []int
			totalSize uint64
			usedSize  atomic.Uint64
		}{
			blocks:    make([]unsafe.Pointer, 0),
			sizes:     make([]uint64, 0),
			freeList:  make([]int, 0),
			totalSize: config.LargePoolSize,
		},
	}

	return &OptimizedAllocator{
		pool:           pool,
		configRegistry: registry,
	}
}

func (a *OptimizedAllocator) getAppropriatePoolTier(size uint64) *poolTier {
	if size <= smallAllocationLimit {
		return (*poolTier)(unsafe.Pointer(&a.pool.small))
	} else if size <= mediumAllocationLimit {
		return (*poolTier)(unsafe.Pointer(&a.pool.medium))
	}
	return (*poolTier)(unsafe.Pointer(&a.pool.large))
}

// UpdateConfig atomically updates the allocator configuration
func (a *OptimizedAllocator) UpdateConfig(config *AllocatorConfig) error {
	if err := a.configRegistry.UpdateConfig(config); err != nil {
		return err
	}

	// Apply immediate changes that don't require restart
	current := a.configRegistry.GetCurrentConfig()
	if config.MixedPrecision != current.MixedPrecision {
		a.tensorCache.Range(func(key, value interface{}) bool {
			a.tensorCache.Delete(key)
			return true
		})
	}

	return nil
}

// GetConfig returns the current configuration
func (a *OptimizedAllocator) GetConfig() *AllocatorConfig {
	return a.configRegistry.GetCurrentConfig()
}

// AllocTensorWithCache implements smart caching for tensor allocations
func (a *OptimizedAllocator) AllocTensorWithCache(size int) (unsafe.Pointer, error) {
	if size <= 0 {
		return nil, errors.New("invalid size: must be positive")
	}

	config := a.configRegistry.GetCurrentConfig()

	// Check cache first if auto-tuning enabled
	if config.EnableAutoTuning {
		if cached, ok := a.tensorCache.Load(size); ok {
			return cached.(unsafe.Pointer), nil
		}
	}

	alignedSize := ((size + 15) / 16) * 16 // 16-byte alignment
	usize := uint64(alignedSize)

	poolTier := a.getAppropriatePoolTier(usize)

	// Check memory pressure for the selected tier
	memoryPressure := float64(poolTier.usedSize.Load()) / float64(poolTier.totalSize)

	// Graduated memory pressure handling
	if memoryPressure > config.CriticalWatermark {
		atomic.StoreInt32(&a.pool.checkpoint.cooling, config.CooldownCycles*2)
		a.reclaimUnusedMemory()
	} else if memoryPressure > config.HighWatermark && atomic.LoadInt32(&a.pool.checkpoint.cooling) == 0 {
		atomic.StoreInt32(&a.pool.checkpoint.cooling, config.CooldownCycles)
		a.reclaimUnusedMemory()
	}

	// Try to allocate
	a.pool.mu.Lock()
	defer a.pool.mu.Unlock()

	// Look for existing block
	for i := 0; i < len(poolTier.freeList); i++ {
		idx := poolTier.freeList[i]
		if poolTier.sizes[idx] >= usize {
			poolTier.freeList = append(poolTier.freeList[:i], poolTier.freeList[i+1:]...)
			poolTier.usedSize.Add(usize)
			return poolTier.blocks[idx], nil
		}
	}

	// Allocate new block
	newSize := uint64(float64(usize) * config.GrowthFactor)
	if newSize < usize {
		newSize = usize // Protect against overflow
	}

	ptr, err := cuda.Malloc(int(newSize))
	if err != nil {
		ptr, err = cuda.Malloc(int(usize)) // Try without growth
		if err != nil {
			return nil, err
		}
		newSize = usize
	}

	poolTier.blocks = append(poolTier.blocks, ptr)
	poolTier.sizes = append(poolTier.sizes, newSize)
	poolTier.totalSize += newSize
	poolTier.usedSize.Add(usize)

	// Update cache if enabled
	if config.EnableAutoTuning && a.lastAllocation == uint64(size) {
		cacheSize := 0
		a.tensorCache.Range(func(_, _ interface{}) bool {
			cacheSize++
			return cacheSize < maxCacheEntries
		})

		if cacheSize < maxCacheEntries {
			a.tensorCache.Store(size, ptr)
		}
	}
	a.lastAllocation = uint64(size)

	return ptr, nil
}

// FreeTensor deallocates a previously allocated tensor
func (a *OptimizedAllocator) FreeTensor(ptr unsafe.Pointer) {
	a.pool.mu.Lock()
	defer a.pool.mu.Unlock()

	// Try to find the pointer in each pool tier
	for _, tier := range []*poolTier{
		(*poolTier)(unsafe.Pointer(&a.pool.small)),
		(*poolTier)(unsafe.Pointer(&a.pool.medium)),
		(*poolTier)(unsafe.Pointer(&a.pool.large)),
	} {
		for i, block := range tier.blocks {
			if block == ptr {
				// Add to free list
				tier.freeList = append(tier.freeList, i)
				// Update used size (subtract from current value)
				if i < len(tier.sizes) {
					currentUsed := tier.usedSize.Load()
					if currentUsed >= tier.sizes[i] {
						tier.usedSize.Store(currentUsed - tier.sizes[i])
					}
				}
				return
			}
		}
	}
}

func (a *OptimizedAllocator) reclaimUnusedMemory() {
	a.pool.mu.Lock()
	defer a.pool.mu.Unlock()

	// Clear tensor cache
	a.tensorCache.Range(func(key, value interface{}) bool {
		a.tensorCache.Delete(key)
		return true
	})

	// Helper function to compact a pool tier
	compactTier := func(tier *poolTier) {
		newBlocks := make([]unsafe.Pointer, 0)
		newSizes := make([]uint64, 0)
		for i := range tier.blocks {
			if i < len(tier.freeList) {
				cuda.Free(tier.blocks[i])
				tier.totalSize -= tier.sizes[i]
			} else {
				newBlocks = append(newBlocks, tier.blocks[i])
				newSizes = append(newSizes, tier.sizes[i])
			}
		}
		tier.blocks = newBlocks
		tier.sizes = newSizes
		tier.freeList = make([]int, 0)
	}

	// Compact each pool tier
	compactTier((*poolTier)(unsafe.Pointer(&a.pool.small)))
	compactTier((*poolTier)(unsafe.Pointer(&a.pool.medium)))
	compactTier((*poolTier)(unsafe.Pointer(&a.pool.large)))
}

// EnableAdaptiveCheckpointing implements smart checkpoint frequency with config validation
// GetMetrics returns the current memory pool metrics
func (a *OptimizedAllocator) GetMetrics() map[string]interface{} {
	getPoolPressure := func(tier *poolTier) float64 {
		return float64(tier.usedSize.Load()) / float64(tier.totalSize)
	}

	config := a.configRegistry.GetCurrentConfig()

	return map[string]interface{}{
		"small_pool_pressure":  getPoolPressure((*poolTier)(unsafe.Pointer(&a.pool.small))),
		"medium_pool_pressure": getPoolPressure((*poolTier)(unsafe.Pointer(&a.pool.medium))),
		"large_pool_pressure":  getPoolPressure((*poolTier)(unsafe.Pointer(&a.pool.large))),
		"total_used_size":      a.pool.getTotalUsedSize(),
		"total_size":           a.pool.getTotalSize(),
		"cache_enabled":        config.EnableAutoTuning,
		"high_watermark":       config.HighWatermark,
		"critical_watermark":   config.CriticalWatermark,
	}
}

func (a *OptimizedAllocator) EnableAdaptiveCheckpointing() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			config := a.configRegistry.GetCurrentConfig()
			totalUsed := float64(a.pool.getTotalUsedSize())
			totalSize := float64(a.pool.getTotalSize())
			pressure := totalUsed / totalSize
			cooling := atomic.LoadInt32(&a.pool.checkpoint.cooling)

			if pressure > config.CriticalWatermark {
				if cooling == 0 {
					atomic.StoreInt32(&a.pool.checkpoint.counter, config.CheckpointFreq*2)
					atomic.StoreInt32(&a.pool.checkpoint.cooling, config.CooldownCycles*2)
				}
			} else if pressure > config.HighWatermark {
				if cooling == 0 {
					atomic.StoreInt32(&a.pool.checkpoint.counter, config.CheckpointFreq)
					atomic.StoreInt32(&a.pool.checkpoint.cooling, config.CooldownCycles)
				}
			} else if pressure < config.LowWatermark && cooling > 0 {
				atomic.AddInt32(&a.pool.checkpoint.cooling, -1)
			}

			// Tier-specific tuning
			for _, tier := range []*poolTier{
				(*poolTier)(unsafe.Pointer(&a.pool.small)),
				(*poolTier)(unsafe.Pointer(&a.pool.medium)),
				(*poolTier)(unsafe.Pointer(&a.pool.large)),
			} {
				tierPressure := float64(tier.usedSize.Load()) / float64(tier.totalSize)
				if tierPressure > config.MediumWatermark {
					// Consider proactive reclamation for this tier
					if cooling == 0 {
						atomic.StoreInt32(&a.pool.checkpoint.counter, config.CheckpointFreq)
					}
				}
			}
		}
	}()
}
