// Copyright (c) 2025 A Bit of Help, Inc.

// Example of a basic ServiceLib application
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/middleware"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, err := logging.NewLogger("info", true)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Create a context logger
	contextLogger := logging.NewContextLogger(logger)

	// Create a simple HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, ServiceLib!"))
	})

	// Add middleware for logging, metrics, and recovery
	handler := middleware.Chain(
		mux,
		middleware.WithRequestID(context.Background()),
		middleware.Logging(contextLogger),
		middleware.Recovery(contextLogger),
	)

	// Start the server
	contextLogger.Info(context.Background(), "Starting server", zap.String("address", ":8080"))
	if err := http.ListenAndServe(":8080", handler); err != nil {
		contextLogger.Fatal(context.Background(), "Server failed", zap.Error(err))
	}
}