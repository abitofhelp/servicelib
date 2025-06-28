//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of manual cancellation with HandleShutdown
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

	// Create a graceful shutdown handler
	shutdownTimeout := 30 * time.Second
	gs := signal.NewGracefulShutdown(shutdownTimeout, logger)

	// Get a context and cancel function
	ctx, cancel := gs.HandleShutdown()

	// Register HTTP server shutdown callback
	gs.RegisterCallback(func(ctx context.Context) error {
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownCancel()

		logger.Info(ctx, "Shutting down HTTP server")
		return server.Shutdown(shutdownCtx)
	})

	// Register a simple handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Add a shutdown endpoint that triggers manual cancellation
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Shutting down server..."))

		// Trigger shutdown after response is sent
		go func() {
			time.Sleep(100 * time.Millisecond)
			logger.Info(context.Background(), "Manual shutdown triggered via /shutdown endpoint")
			cancel() // This will trigger the same shutdown process as a signal would
		}()
	})

	// Start the server in a goroutine
	go func() {
		logger.Info(context.Background(), "Starting server on :8080")
		logger.Info(context.Background(), "Visit http://localhost:8080/shutdown to trigger manual shutdown")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "Server error", zap.Error(err))
		}
	}()

	// Wait for context to be canceled (either by signal or manual cancellation)
	<-ctx.Done()

	// The registered callbacks will be executed automatically
	// Just wait here until the process exits
	fmt.Println("Shutdown triggered. Waiting for graceful shutdown to complete...")

	// In a real application, we would block here until the process exits
	// For this example, we'll sleep briefly to allow logs to be displayed
	time.Sleep(100 * time.Millisecond)
}
