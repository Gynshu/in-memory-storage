// Package limiter  provides a rate limiter for the API.
// Underlying structure and interface for http middleware.
// It is used to limit the number of requests per second.
// Implements the fallowing methods:
// - Limit(ip string)
// - Check(ip string) bool
package limit

import (
	"github.com/gynshu-one/in-memory-storage/internal/config"
	"sync"
	"time"
)

// rateLimiter represents the rate limiter for the API.
type rateLimiter struct {
	// map of IP addresses and the last request time
	Requests map[string]int64
	MaxRPN   int64
	Mutex    *sync.Mutex
}

// NewRateLimiter creates a new rateLimiter with the given max requests per second.
func NewRateLimiter() *rateLimiter {
	return &rateLimiter{
		Requests: make(map[string]int64),
		MaxRPN:   time.Second.Nanoseconds() / config.GetConf().RateLimit,
		Mutex:    &sync.Mutex{},
	}
}

// Limit limits the requests per second for the given IP address.
func (rl *rateLimiter) Limit(ip string) {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()
	rl.Requests[ip] = time.Now().UnixNano()
}

// Check checks if the number of requests per second for the given IP address is less than the maximum.
// It takes last request time from map and compares it with the current time.
// if the difference is less than the configured limit it returns false, otherwise true.
func (rl *rateLimiter) Check(ip string) bool {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()
	if lastRequest, ok := rl.Requests[ip]; ok {
		if time.Now().UnixNano()-lastRequest < rl.MaxRPN {
			return false
		}
	}
	return true
}
