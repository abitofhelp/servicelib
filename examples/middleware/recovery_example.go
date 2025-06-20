// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using the recovery middleware
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/middleware"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create a context logger
	contextLogger := logging.NewContextLogger(logger)

	// Create a handler that will panic
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request wants to trigger a panic
		if r.URL.Query().Get("panic") == "true" {
			// Simulate a panic
			panic("This is a simulated panic!")
		}

		// Normal operation
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("No panic occurred"))
	})

	// Apply recovery middleware to catch panics
	safeHandler := middleware.WithRecovery(contextLogger, panicHandler)

	// Register the handler
	http.Handle("/", middleware.WithRequestContext(safeHandler))

	// Start the server
	fmt.Println("Server starting on :8080")
	fmt.Println("Try accessing: http://localhost:8080/ (normal operation)")
	fmt.Println("Try accessing: http://localhost:8080/?panic=true (will recover from panic)")

	// In a real application, you would start the server:
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	//     log.Fatalf("Failed to start server: %v", err)
	// }
}
