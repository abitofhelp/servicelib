// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic usage of the middleware package
package example_middleware

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

	// Create a simple HTTP handler
	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, World!"))
	})

	// Apply middleware to the handler
	handler := middleware.WithRequestContext(helloHandler)
	handler = middleware.WithRecovery(contextLogger, handler)
	handler = middleware.WithLogging(contextLogger, handler)
	handler = middleware.WithErrorHandling(handler)
	handler = middleware.WithCORS(handler)

	// Register the handler
	http.Handle("/", handler)

	// Start the server
	fmt.Println("Server starting on :8080")
	fmt.Println("Try accessing: http://localhost:8080/")

	// In a real application, you would start the server:
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	//     log.Fatalf("Failed to start server: %v", err)
	// }
}
