// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the Auth middleware
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abitofhelp/servicelib/auth"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create a context
	ctx := context.Background()

	// Create a configuration
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "your-secret-key"
	config.Middleware.SkipPaths = []string{"/public", "/health"}

	// Create an auth instance
	authInstance, err := auth.New(ctx, config, logger)
	if err != nil {
		logger.Fatal("Failed to create auth instance", zap.Error(err))
	}

	// Get the middleware function
	middleware := authInstance.Middleware()

	// Define a handler function
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, authenticated user!"))
	})

	// Use the middleware with the handler
	http.Handle("/", middleware(handler))
	http.Handle("/public", handler) // This path will skip authentication

	fmt.Println("Server configured with auth middleware")
	fmt.Println("Protected route: http://localhost:8080/")
	fmt.Println("Public route: http://localhost:8080/public")

	// In a real application, you would start the server:
	// http.ListenAndServe(":8080", nil)
}
