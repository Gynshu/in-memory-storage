package limit

import (
	"github.com/gynshu-one/in-memory-storage/internal/config"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	tests := []struct {
		name string
		want *rateLimiter
	}{
		{
			name: "NewRateLimiter returns a new instance of rateLimiter",
			want: &rateLimiter{
				Requests: make(map[string]int64),
				MaxRPN:   time.Second.Nanoseconds() / config.GetConf().RateLimit,
				Mutex:    &sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRateLimiter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRateLimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rateLimiter_Check(t *testing.T) {
	newLimiter := NewRateLimiter()

	limitCFG := time.Second.Nanoseconds() / config.GetConf().RateLimit
	for i := 0; i < 10; i++ {
		newLimiter.Limit("i")
		time.Sleep(time.Duration(limitCFG/2) * time.Nanosecond)
		if newLimiter.Check("i") {
			t.Errorf("Check() = %v, want %v", false, true)
		}
	}

	time.Sleep(time.Duration(limitCFG+2) * time.Nanosecond)
	if !newLimiter.Check("i") {
		t.Errorf("Check() = %v, want %v", true, false)
	}
}
