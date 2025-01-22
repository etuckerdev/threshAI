package memory

import (
	"sync"
	"time"
	"unsafe"

	"github.com/cornelk/hashmap"
)

type ShardState int

const (
	ShardUnloaded ShardState = iota
	ShardLoading
	ShardLoaded
	ShardError
)

type ModelShard struct {
	ID       string
	State    ShardState
	Size     uint64
	Location unsafe.Pointer
	LastUsed time.Time
}

type Vault struct {
	mu         sync.RWMutex
	store      *hashmap.Map[string, interface{}]
	shards     map[string]*ModelShard
	allocator  *GPUMemoryAllocator
	brutalMode int
}

func NewVault(allocator *GPUMemoryAllocator) *Vault {
	return &Vault{
		store:     hashmap.New[string, interface{}](),
		shards:    make(map[string]*ModelShard),
		allocator: allocator,
	}
}

func (v *Vault) SetBrutalMode(mode int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.brutalMode = mode
}

func (v *Vault) Store(key string, value interface{}) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.store.Set(key, value)
}

func (v *Vault) Retrieve(key string) (interface{}, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.store.Get(key)
}

func (v *Vault) Purge() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.store = hashmap.New[string, interface{}]()
}

func (v *Vault) Size() int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.store.Len()
}

func (v *Vault) Export() map[string]interface{} {
	v.mu.RLock()
	defer v.mu.RUnlock()

	result := make(map[string]interface{})
	v.store.Range(func(key string, value interface{}) bool {
		result[key] = value
		return true
	})
	return result
}

func (v *Vault) LoadShard(shardID string, size uint64) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if shard, exists := v.shards[shardID]; exists && shard.State == ShardLoaded {
		return nil
	}

	ptr, err := v.allocator.Alloc4BitTensorCache(int(size))
	if err != nil {
		return err
	}

	v.shards[shardID] = &ModelShard{
		ID:       shardID,
		State:    ShardLoaded,
		Size:     size,
		Location: ptr,
		LastUsed: time.Now(),
	}
	return nil
}

func (v *Vault) UnloadShard(shardID string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	shard, exists := v.shards[shardID]
	if !exists || shard.State != ShardLoaded {
		return nil
	}

	err := v.allocator.Free(shard.Location)
	if err != nil {
		return err
	}

	shard.State = ShardUnloaded
	shard.Location = nil
	return nil
}

func (v *Vault) GetShardState(shardID string) ShardState {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if shard, exists := v.shards[shardID]; exists {
		return shard.State
	}
	return ShardUnloaded
}

func (v *Vault) ActiveShardMemory() uint64 {
	v.mu.RLock()
	defer v.mu.RUnlock()

	var total uint64
	for _, shard := range v.shards {
		if shard.State == ShardLoaded {
			total += shard.Size
		}
	}
	return total
}

const (
	brutalMode1VRAM = 8 * 1024 * 1024 * 1024  // 8GB
	brutalMode2VRAM = 12 * 1024 * 1024 * 1024 // 12GB
	brutalMode3VRAM = 16 * 1024 * 1024 * 1024 // 16GB
)

func (v *Vault) EnsureBrutalModeConstraints() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	var maxVRAM uint64
	switch v.brutalMode {
	case 1:
		maxVRAM = brutalMode1VRAM
	case 2:
		maxVRAM = brutalMode2VRAM
	case 3:
		maxVRAM = brutalMode3VRAM
	default:
		return nil
	}

	// Unload shards until we're under the VRAM limit
	for v.ActiveShardMemory() > maxVRAM {
		// Find the least recently used shard
		var oldestShard *ModelShard
		for _, shard := range v.shards {
			if shard.State == ShardLoaded && (oldestShard == nil || shard.LastUsed.Before(oldestShard.LastUsed)) {
				oldestShard = shard
			}
		}

		if oldestShard == nil {
			break
		}

		if err := v.UnloadShard(oldestShard.ID); err != nil {
			return err
		}
	}
	return nil
}

func (v *Vault) LoadQuantumCache() error {
	if v.brutalMode < 3 {
		return nil
	}

	// Allocate quantum cache memory
	quantumSize := uint64(2 * 1024 * 1024 * 1024) // 2GB
	ptr, err := v.allocator.Alloc4BitTensorCache(int(quantumSize))
	if err != nil {
		return err
	}

	v.shards["quantum_cache"] = &ModelShard{
		ID:       "quantum_cache",
		State:    ShardLoaded,
		Size:     quantumSize,
		Location: ptr,
	}
	return nil
}
