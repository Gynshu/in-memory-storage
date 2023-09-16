package main

import (
	"context"
	"errors"
	"github.com/gynshu-one/in-memory-storage/internal/api"
	"github.com/gynshu-one/in-memory-storage/internal/config"
	ratelimiter "github.com/gynshu-one/in-memory-storage/internal/infra/limit"
	"github.com/gynshu-one/in-memory-storage/internal/infra/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Init repo and rate limiter
	NewLimiter := ratelimiter.NewRateLimiter()
	repo := storage.NewInMemory()

	Rlm := api.RateLimiterMiddleware(NewLimiter)

	// Create a new router
	router := api.NewRouter()

	// Create a new handlers
	hands := api.NewHandlers(repo)

	// Add the middlewares to the router
	router.Use(api.LoggingMiddleware)
	router.Use(Rlm)

	// Add the routes to the router
	router.Post("/set", hands.Set)
	router.Delete("/delete", hands.Delete)
	router.Get("/get", hands.Get)
	router.Get("/all", hands.GetAll)

	// Init the server
	srv := &http.Server{
		Addr:        ":" + config.GetConf().ServerPort,
		Handler:     router,
		ReadTimeout: 10 * time.Second,
	}

	// Start the server in a separate goroutine
	go func() {
		log.Printf("Server listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 11*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Println("server shutdown successfully")
}
