package api

import (
	"fmt"
	"github.com/gynshu-one/in-memory-storage/internal/domain"
	"net/http"
	"time"
)

// RateLimiterMiddleware returns a middleware function that limits the number of requests per second for a given IP address.
func RateLimiterMiddleware(rl domain.RateLimiter) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				ip = forwardedFor
			}

			if !rl.Check(ip) {
				w.WriteHeader(http.StatusTooManyRequests)
				_, err := w.Write([]byte("Too many requests"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
			rl.Limit(ip)
			next.ServeHTTP(w, r)
		}
	}
}

// LoggingMiddleware returns a middleware function that logs the HTTP requests and responses.
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		fmt.Printf("%s %s %s %d %d %v\n", r.Method, r.URL.Path, r.Proto, rw.status, rw.length, duration)
	}
}

// responseWriter is a custom http.ResponseWriter that keeps track of the status code and response length.
type responseWriter struct {
	http.ResponseWriter
	status int
	length int
}

// WriteHeader writes the HTTP status code to the response.
func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Write writes the response body to the response and keeps track of the response length.
func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.length += n
	return n, err
}
