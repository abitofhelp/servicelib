// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using the error handling middleware
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/middleware"
)

func main() {
	// Create a handler that may return errors
	errorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the error type from the query parameter
		errorType := r.URL.Query().Get("error")

		// Simulate different error scenarios
		switch errorType {
		case "validation":
			// Simulate a validation error
			// This will be mapped to HTTP 400 Bad Request
			panic(errors.Validation("Invalid input: field 'name' is required"))
		case "notfound":
			// Simulate a not found error
			// This will be mapped to HTTP 404 Not Found
			panic(errors.NotFound("Resource with ID %s not found", "12345"))
		case "unauthorized":
			// Simulate an unauthorized error
			// This will be mapped to HTTP 401 Unauthorized
			panic(errors.Unauthorized("Authentication required"))
		case "forbidden":
			// Simulate a forbidden error
			// This will be mapped to HTTP 403 Forbidden
			panic(errors.Forbidden("Insufficient permissions"))
		case "internal":
			// Simulate an internal error
			// This will be mapped to HTTP 500 Internal Server Error
			err := fmt.Errorf("database connection failed")
			panic(errors.Internal(err, "An unexpected error occurred"))
		case "timeout":
			// Simulate a timeout error
			// This will be mapped to HTTP 504 Gateway Timeout
			panic(errors.Timeout("Operation timed out after %d seconds", 30))
		default:
			// No error, normal response
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(fmt.Sprintf("No error occurred. Current time: %s", time.Now().Format(time.RFC3339))))
		}
	})

	// Apply error handling middleware to catch and handle errors
	handler := middleware.WithErrorHandling(errorHandler)

	// Apply recovery middleware to catch panics and convert them to errors
	handler = middleware.WithRequestContext(handler)

	// Register the handler
	http.Handle("/", handler)

	// Start the server
	fmt.Println("Server starting on :8080")
	fmt.Println("Try accessing:")
	fmt.Println("  http://localhost:8080/ (success response)")
	fmt.Println("  http://localhost:8080/?error=validation (validation error)")
	fmt.Println("  http://localhost:8080/?error=notfound (not found error)")
	fmt.Println("  http://localhost:8080/?error=unauthorized (unauthorized error)")
	fmt.Println("  http://localhost:8080/?error=forbidden (forbidden error)")
	fmt.Println("  http://localhost:8080/?error=internal (internal error)")
	fmt.Println("  http://localhost:8080/?error=timeout (timeout error)")

	// In a real application, you would start the server:
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	//     log.Fatalf("Failed to start server: %v", err)
	// }
}