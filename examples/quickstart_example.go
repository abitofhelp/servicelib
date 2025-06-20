// Copyright (c) 2025 A Bit of Help, Inc.

// Example of a basic ServiceLib application
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/middleware"
)

func main() {
	// Create a logger
	logger, err := logging.NewLogger(logging.Config{
		Level:  "info",
		Format: "json",
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Create a simple HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, ServiceLib!"))
	})

	// Add middleware for logging, metrics, and recovery
	handler := middleware.Chain(
		mux,
		middleware.RequestID(),
		middleware.Logging(logger),
		middleware.Recovery(logger),
	)

	// Start the server
	logger.Info(context.Background(), "Starting server", "address", ":8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		logger.Fatal(context.Background(), "Server failed", "error", err)
	}
}
