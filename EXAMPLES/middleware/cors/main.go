//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using the CORS middleware
package main

import (
	"fmt"
	"net/http"

	"github.com/abitofhelp/servicelib/middleware"
)

func main() {
	// Create a simple API handler
	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set content type
		w.Header().Set("Content-Type", "application/json")

		// Return a simple JSON response
		w.Write([]byte(`{"message":"This is a CORS-enabled API endpoint"}`))
	})

	// Apply CORS middleware to the handler
	corsHandler := middleware.WithCORS(apiHandler)

	// Register the handler with request context
	http.Handle("/api", middleware.WithRequestContext(corsHandler))

	// Start the server
	fmt.Println("Server starting on :8080")
	fmt.Println("API endpoint: http://localhost:8080/api")
	fmt.Println("")
	fmt.Println("To test CORS, you can use curl:")
	fmt.Println("  curl -H \"Origin: http://example.com\" -v http://localhost:8080/api")
	fmt.Println("")
	fmt.Println("Or test a preflight request:")
	fmt.Println("  curl -H \"Origin: http://example.com\" -H \"Access-Control-Request-Method: POST\" -H \"Access-Control-Request-Headers: Content-Type\" -X OPTIONS -v http://localhost:8080/api")

	// In a real application, you would start the server:
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	//     log.Fatalf("Failed to start server: %v", err)
	// }
}