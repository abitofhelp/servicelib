//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of shutting down multiple resources
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/shutdown"
	"go.uber.org/zap"
)

// MockDatabase simulates a database connection
type MockDatabase struct {
	logger *logging.ContextLogger
}

// NewMockDatabase creates a new mock database
func NewMockDatabase(logger *logging.ContextLogger) *MockDatabase {
	logger.Info(context.Background(), "Opening database connection")
	return &MockDatabase{
		logger: logger,
	}
}

// Close simulates closing a database connection
func (db *MockDatabase) Close() error {
	db.logger.Info(context.Background(), "Closing database connection")
	// Simulate some work being done during close
	time.Sleep(500 * time.Millisecond)
	return nil
}

// MockCache simulates a cache service
type MockCache struct {
	logger *logging.ContextLogger
}

// NewMockCache creates a new mock cache
func NewMockCache(logger *logging.ContextLogger) *MockCache {
	logger.Info(context.Background(), "Initializing cache service")
	return &MockCache{
		logger: logger,
	}
}

// Shutdown simulates shutting down a cache service
func (c *MockCache) Shutdown() error {
	c.logger.Info(context.Background(), "Shutting down cache service")
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
	db := NewMockDatabase(logger)
	cache := NewMockCache(logger)

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

	// Define shutdown function for multiple resources
	shutdownFunc := func() error {
		// Create a context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Shutdown HTTP server first
		logger.Info(ctx, "Shutting down HTTP server")
		if err := server.Shutdown(ctx); err != nil {
			logger.Error(ctx, "Error shutting down server", zap.Error(err))
			return err
		}

		// Then shutdown cache
		if err := cache.Shutdown(); err != nil {
			logger.Error(ctx, "Error shutting down cache", zap.Error(err))
			return err
		}

		// Finally close database connection
		if err := db.Close(); err != nil {
			logger.Error(ctx, "Error closing database", zap.Error(err))
			return err
		}

		return nil
	}

	// Wait for shutdown signal
	ctx := context.Background()
	fmt.Println("Server is running with multiple resources.")
	fmt.Println("Press Ctrl+C to initiate graceful shutdown.")
	err = shutdown.GracefulShutdown(ctx, logger, shutdownFunc)
	if err != nil {
		logger.Error(ctx, "Error during shutdown", zap.Error(err))
	}

	logger.Info(ctx, "Application stopped")
	fmt.Println("All resources have been gracefully shut down.")
}
