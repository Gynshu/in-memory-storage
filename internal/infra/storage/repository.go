// Package storage provides simple redis-like in-memory storage.
// It is used to store key-value pairs.
// Implements the fallowing methods:
// - Set(key string, value string, ttl time.Duration) error
// - Delete(key string) error
// - Get(key string) (string, error)
// - GetAll() (map[string]string, error)

package storage

import (
	"sync"
	"time"

	"github.com/gynshu-one/in-memory-storage/internal/domain"
)

// storage represents the in-memory storage.
type storage struct {
	mu      *sync.RWMutex
	storage map[string]domain.Entity
}

// NewInMemory creates a new instance of storage.
// It returns a pointer to the newly created instance.
func NewInMemory() *storage {
	return &storage{
		mu:      &sync.RWMutex{},
		storage: make(map[string]domain.Entity),
	}
}

// Set adds a new key-value pair to the storage or replaces it if it already exists.
// If the key already exists, it returns an error.
// If the ttl is 0, the key-value pair will not expire.
func (i *storage) Set(key string, value string, ttl time.Duration) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	exp := time.Now().Add(ttl).UnixNano()

	if ttl == 0 {
		exp = 0
	}

	i.storage[key] = domain.Entity{
		Key:        key,
		Value:      value,
		Expiration: exp,
	}

	return nil
}

// Delete deletes a key from the storage.
func (i *storage) Delete(key string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if entity, ok := i.storage[key]; !ok {
		if entity.IsExpired() {
			delete(i.storage, key)
			return domain.ErrKeyExpired
		}
		return domain.ErrKeyNotFound
	}

	delete(i.storage, key)

	return nil
}

// Get gets the value of a key from the storage.
func (i *storage) Get(key string) (string, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if entity, ok := i.storage[key]; ok {
		if entity.IsExpired() {
			delete(i.storage, key)
			return "", domain.ErrKeyExpired
		}
		return entity.Value, nil
	}

	return "", domain.ErrKeyNotFound
}

// GetAll gets all the key-value pairs from the storage. Returns copy
func (i *storage) GetAll() ([]domain.Entity, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var result []domain.Entity
	for key, entity := range i.storage {
		if entity.IsExpired() {
			delete(i.storage, key)
			continue
		}

		result = append(result, entity)
	}

	if len(result) == 0 {
		return nil, domain.ErrStorageEmpty
	}
	return result, nil
}
