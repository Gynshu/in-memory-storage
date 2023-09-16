// Package config provides configuration for the application.
// init function reads the configuration from environment variables.
// GetConf returns a new config instance with default values.
package config

import (
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func init() {
	port := os.Getenv("SERVER_PORT")
	rateLimit := os.Getenv("RATE_LIMIT")
	if port != "" {
		cfg.ServerPort = port
	}
	if rateLimit != "" {
		limit, err := strconv.Atoi(rateLimit)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to convert rate limit to int")
		}
		cfg.RateLimit = int64(limit)
	}

}

var cfg = &config{
	ServerPort: "8080",
	RateLimit:  10,
}

// config represents the configuration for the application.
type config struct {
	ServerPort string `json:"server_port"`
	// RateLimitWindow is the time window for rate limiting in seconds.
	RateLimit int64 `json:"rate_limit"`
}

// GetConf returns a new config instance with default values.
func GetConf() *config {
	return cfg
}
