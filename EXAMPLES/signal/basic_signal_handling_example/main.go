// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic signal handling
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

func main() {
	// Create a logger
	baseLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	logger := logging.NewContextLogger(baseLogger)
	defer baseLogger.Sync()

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

	// Wait for shutdown signal
	shutdownTimeout := 30 * time.Second
	ctx := signal.WaitForShutdown(shutdownTimeout, logger)

	// Perform shutdown when context is canceled
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info(ctx, "Shutting down HTTP server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "Error shutting down server", zap.Error(err))
	}

	logger.Info(ctx, "Server stopped")
	fmt.Println("Server has been gracefully shut down.")
}
