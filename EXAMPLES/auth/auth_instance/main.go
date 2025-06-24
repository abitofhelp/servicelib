// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of creating an Auth instance
package main

import (
	"context"
	"fmt"

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

	// Create an auth instance
	_, err := auth.New(ctx, config, logger)
	if err != nil {
		logger.Fatal("Failed to create auth instance", zap.Error(err))
	}

	fmt.Println("Auth instance created successfully")
}