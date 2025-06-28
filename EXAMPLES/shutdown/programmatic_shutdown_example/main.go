//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of programmatic shutdown initiation
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

	// Register handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Add a shutdown endpoint
	var shutdownTrigger context.CancelFunc
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Shutting down server..."))

		// Trigger shutdown after response is sent
		go func() {
			time.Sleep(100 * time.Millisecond)
			if shutdownTrigger != nil {
				logger.Info(context.Background(), "Programmatic shutdown triggered via /shutdown endpoint")
				shutdownTrigger()
			}
		}()
	})

	// Start the server in a goroutine
	go func() {
		logger.Info(context.Background(), "Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "Server error", zap.Error(err))
		}
	}()

	// Define shutdown function
	shutdownFunc := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.Info(ctx, "Shutting down HTTP server")
		return server.Shutdown(ctx)
	}

	// Setup graceful shutdown with programmatic trigger
	ctx := context.Background()
	fmt.Println("Server is running. Visit http://localhost:8080/shutdown to trigger shutdown.")
	fmt.Println("You can also press Ctrl+C to initiate graceful shutdown.")
	shutdownTrigger, errCh := shutdown.SetupGracefulShutdown(ctx, logger, shutdownFunc)

	// Wait for shutdown to complete
	err = <-errCh
	if err != nil {
		logger.Error(ctx, "Error during shutdown", zap.Error(err))
	}

	logger.Info(ctx, "Server stopped")
	fmt.Println("Server has been gracefully shut down.")
}
