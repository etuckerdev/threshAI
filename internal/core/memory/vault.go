package memory

import (
	"sync"

	"github.com/cornelk/hashmap"
)

type Vault struct {
	mu    sync.RWMutex
	store *hashmap.Map[string, interface{}]
}

func NewVault() *Vault {
	return &Vault{
		store: hashmap.New[string, interface{}](),
	}
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
