package domain

import "time"

// Repository defines the methods for interacting with the in-memory storage.
type Repository interface {
	// Set adds a new key-value pair to the storage or replaces it if it already exists.
	// If the key already exists, it returns an error.
	// If the ttl is 0, the key-value pair will not expire.
	Set(key string, value string, ttl time.Duration) error
	// Delete deletes a key from the storage.
	Delete(key string) error
	// Get gets the value of a key from the storage.
	Get(key string) (string, error)
	// GetAll gets all the key-value pairs from the storage. Returns copy
	GetAll() ([]Entity, error)
}
