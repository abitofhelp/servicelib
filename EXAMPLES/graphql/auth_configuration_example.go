// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating how to configure the auth service for GraphQL
package example_graphql

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/auth"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create a context
	ctx := context.Background()

	// This example shows how to configure the auth service for use with GraphQL
	fmt.Println("Example: Configuring auth service for GraphQL")

	// Create a configuration for the auth service
	fmt.Println("\nStep 1: Configure the auth service")
	fmt.Println(`
	authConfig := auth.DefaultConfig()
	authConfig.JWT.SecretKey = cfg.Auth.JWT.SecretKey
	authConfig.JWT.Issuer = cfg.Auth.JWT.Issuer
	authConfig.JWT.TokenDuration = cfg.Auth.JWT.TokenDuration
	authConfig.Middleware.SkipPaths = []string{"/health", "/metrics", "/playground"}
	`)

	// For demonstration purposes, we'll create a simple config
	authConfig := auth.DefaultConfig()
	authConfig.JWT.SecretKey = "your-secret-key"
	authConfig.JWT.Issuer = "your-service"
	authConfig.JWT.TokenDuration = 24 * time.Hour
	authConfig.Middleware.SkipPaths = []string{"/health", "/metrics", "/playground"}

	// Initialize the auth service
	fmt.Println("\nStep 2: Initialize the auth service")
	fmt.Println(`
	authService, err := auth.New(ctx, authConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth service: %w", err)
	}
	`)

	// Actually initialize the service for demonstration
	authService, err := auth.New(ctx, authConfig, logger)
	if err != nil {
		fmt.Printf("Error initializing auth service: %v\n", err)
		return
	}

	fmt.Println("\nAuth service configured successfully!")
	fmt.Printf("JWT Issuer: %s\n", authConfig.JWT.Issuer)
	fmt.Printf("Token Duration: %v\n", authConfig.JWT.TokenDuration)
	fmt.Printf("Skip Paths: %v\n", authConfig.Middleware.SkipPaths)

	// Note about the auth service
	fmt.Println("\nThe auth service provides:")
	fmt.Println("- JWT token generation and validation")
	fmt.Println("- Middleware for authenticating HTTP requests")
	fmt.Println("- Functions for checking authorization")
	fmt.Println("- Context utilities for working with user information")

	// Prevent unused variable warning
	_ = authService
}
