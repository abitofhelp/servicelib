// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using the logging middleware
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

	// Create a handler that simulates different response types
	demoHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the response type from the query parameter
		responseType := r.URL.Query().Get("type")

		// Set content type
		w.Header().Set("Content-Type", "text/plain")

		// Simulate processing time
		time.Sleep(100 * time.Millisecond)

		// Handle different response types
		switch responseType {
		case "error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		case "notfound":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, World!"))
		}
	})

	// Apply logging middleware
	loggingHandler := middleware.WithLogging(contextLogger, demoHandler)

	// Register the handler with request context
	http.Handle("/", middleware.WithRequestContext(loggingHandler))

	// Start the server
	fmt.Println("Server starting on :8080")
	fmt.Println("Try accessing:")
	fmt.Println("  http://localhost:8080/ (success response)")
	fmt.Println("  http://localhost:8080/?type=error (error response)")
	fmt.Println("  http://localhost:8080/?type=notfound (not found response)")
	fmt.Println("Check the logs to see the request and response details")

	// In a real application, you would start the server:
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	//     log.Fatalf("Failed to start server: %v", err)
	// }
}