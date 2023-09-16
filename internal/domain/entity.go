package domain

import "time"

// Entity represents a key-value pair in the in-memory storage.
type Entity struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	// Expiration is the time in nanoseconds when the key-value pair will expire.
	Expiration int64 `json:"expiration"`
}

func (e *Entity) IsExpired() bool {
	return e.Expiration > 0 && time.Now().UnixNano() > e.Expiration
}
