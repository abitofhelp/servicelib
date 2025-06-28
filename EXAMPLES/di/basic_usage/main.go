//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic usage of the dependency injection package
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/abitofhelp/servicelib/di"
	"go.uber.org/zap"
)

func main() {
	// Create a context
	ctx := context.Background()

	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create a simple configuration
	config := map[string]interface{}{
		"app_name": "example",
		"version":  "1.0.0",
	}

	// Create a container
	container, err := di.NewContainer(ctx, logger, config)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	// Use the container
	fmt.Println("Container created successfully")

	// Get the context from the container
	// containerCtx := container.GetContext()
	fmt.Println("Context retrieved from container")

	// Get the logger from the container
	containerLogger := container.GetLogger()
	containerLogger.Info("Logger retrieved from container")

	// Get the configuration from the container
	containerConfig := container.GetConfig()
	fmt.Printf("Config retrieved from container: %v\n", containerConfig)
}
