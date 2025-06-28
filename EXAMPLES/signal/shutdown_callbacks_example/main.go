//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using shutdown callbacks
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/signal"
	"go.uber.org/zap"
)

// MockDatabase simulates a database connection
type MockDatabase struct {
	logger *logging.ContextLogger
}

// Close simulates closing a database connection
func (db *MockDatabase) Close(ctx context.Context) error {
	db.logger.Info(ctx, "Closing database connection")
	// Simulate some work being done during close
	time.Sleep(500 * time.Millisecond)
	return nil
}

// MockCache simulates a cache service
type MockCache struct {
	logger *logging.ContextLogger
}

// Shutdown simulates shutting down a cache service
func (c *MockCache) Shutdown(ctx context.Context) error {
	c.logger.Info(ctx, "Shutting down cache service")
	// Simulate some work being done during shutdown
	time.Sleep(300 * time.Millisecond)
	return nil
}

func main() {
	// Create a logger
	baseLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	logger := logging.NewContextLogger(baseLogger)
	defer baseLogger.Sync()

	// Create resources
	db := &MockDatabase{logger: logger}
	cache := &MockCache{logger: logger}

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	// Register a simple handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Start the server in a goroutine
	go func() {
		logger.Info(context.Background(), "Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "Server error", zap.Error(err))
		}
	}()

	// Setup graceful shutdown with callbacks
	shutdownTimeout := 30 * time.Second
	ctx, gs := signal.SetupSignalHandler(shutdownTimeout, logger)

	// Register shutdown callbacks
	// Order matters: callbacks are executed concurrently, but it's good practice
	// to register them in reverse order of creation

	// Register database shutdown callback
	gs.RegisterCallback(func(ctx context.Context) error {
		return db.Close(ctx)
	})

	// Register cache shutdown callback
	gs.RegisterCallback(func(ctx context.Context) error {
		return cache.Shutdown(ctx)
	})

	// Register HTTP server shutdown callback
	gs.RegisterCallback(func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		logger.Info(ctx, "Shutting down HTTP server")
		return server.Shutdown(shutdownCtx)
	})

	// Wait for context to be canceled (when a signal is received)
	<-ctx.Done()

	// The registered callbacks will be executed automatically
	// Just wait here until the process exits
	fmt.Println("Shutdown signal received. Waiting for graceful shutdown to complete...")

	// In a real application, we would block here until the process exits
	// For this example, we'll sleep briefly to allow logs to be displayed
	time.Sleep(100 * time.Millisecond)
}
