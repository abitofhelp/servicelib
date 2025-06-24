// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating how to check authorization in GraphQL operations
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abitofhelp/servicelib/graphql"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// ExampleAuthorizationCheck demonstrates how to properly check authorization
// for GraphQL operations using the servicelib authorization middleware.
func ExampleAuthorizationCheck(ctx context.Context, logger *logging.ContextLogger) error {
	// Create a context with timeout for safety
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Define authorization parameters
	allowedRoles := []string{"ADMIN", "EDITOR"}
	requiredScopes := []string{"READ"}
	resource := "ITEM"
	operationName := "ExampleOperation"

	// Check authorization with proper error handling
	if err := graphql.CheckAuthorization(
		ctx,
		allowedRoles,
		requiredScopes,
		resource,
		operationName,
		logger,
	); err != nil {
		// Log the error with context
		logger.Error(ctx, "Authorization check failed",
			zap.Error(err),
			zap.Strings("roles", allowedRoles),
			zap.Strings("scopes", requiredScopes),
			zap.String("resource", resource),
			zap.String("operation", operationName),
		)
		return fmt.Errorf("authorization failed: %w", err)
	}

	return nil
}

func main() {
	// Initialize logger
	baseLogger, err := logging.NewLogger("info", true)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer baseLogger.Sync()

	// Create context logger
	contextLogger := logging.NewContextLogger(baseLogger)

	// Create a base context
	ctx := context.Background()

	// Run the example
	if err := ExampleAuthorizationCheck(ctx, contextLogger); err != nil {
		log.Fatalf("Failed to run authorization example: %v", err)
	}

	fmt.Println("Authorization check completed successfully")
}
