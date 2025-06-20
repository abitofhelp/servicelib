// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating how to apply the auth middleware to a GraphQL handler
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

	// This example shows how to apply the auth middleware to a GraphQL handler
	fmt.Println("Example: Applying auth middleware to a GraphQL handler")

	// Step 1: Create an auth service (simplified for this example)
	fmt.Println("\nStep 1: Create an auth service")
	authConfig := auth.DefaultConfig()
	authConfig.JWT.SecretKey = "your-secret-key"
	authService, err := auth.New(ctx, authConfig, logger)
	if err != nil {
		fmt.Printf("Error initializing auth service: %v\n", err)
		return
	}

	// Step 2: Create a dummy handler (this would be your GraphQL handler)
	fmt.Println("\nStep 2: Create a GraphQL handler")
	graphqlHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GraphQL response"))
	})

	// Step 3: Apply the auth middleware to the handler
	fmt.Println("\nStep 3: Apply the auth middleware to the handler")
	fmt.Println(`
	// Apply the auth middleware to the GraphQL handler
	handler = authService.Middleware()(handler)
	`)

	// Actually apply the middleware
	protectedHandler := authService.Middleware()(graphqlHandler)

	// Explain what the middleware does
	fmt.Println("\nThe auth middleware:")
	fmt.Println("1. Extracts the JWT token from the Authorization header")
	fmt.Println("2. Validates the token")
	fmt.Println("3. Extracts user information (ID, roles, scopes, resources)")
	fmt.Println("4. Adds the user information to the request context")
	fmt.Println("5. Passes the request to the next handler if the token is valid")
	fmt.Println("6. Returns an error response if the token is invalid")

	// Prevent unused variable warning
	_ = protectedHandler
}