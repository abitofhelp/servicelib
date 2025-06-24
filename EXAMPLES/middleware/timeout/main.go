// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using the timeout middleware
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/abitofhelp/servicelib/middleware"
)

func main() {
	// Create a handler that simulates a slow operation
	slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request wants to simulate a timeout
		if r.URL.Query().Get("slow") == "true" {
			// Sleep for 3 seconds to simulate a slow operation
			time.Sleep(3 * time.Second)
		}

		// Write response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Operation completed successfully"))
	})

	// Apply timeout middleware with a 2-second timeout
	timeoutHandler := middleware.WithTimeout(2 * time.Second)(slowHandler)

	// Register the handler
	http.Handle("/", middleware.WithRequestContext(timeoutHandler))

	// Start the server
	fmt.Println("Server starting on :8080")
	fmt.Println("Try accessing: http://localhost:8080/ (normal operation)")
	fmt.Println("Try accessing: http://localhost:8080/?slow=true (will timeout)")

	// In a real application, you would start the server:
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	//     log.Fatalf("Failed to start server: %v", err)
	// }
}