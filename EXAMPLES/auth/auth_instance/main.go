//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Package main demonstrates how to create and configure an Auth instance.
//
// This example shows the basic steps to initialize the authentication system:
// - Creating a logger for the auth system to use
// - Setting up a configuration with appropriate security settings
// - Initializing the auth instance with the configuration
//
// The auth package provides a comprehensive authentication and authorization
// system that can be used to secure your application's endpoints and resources.
// This example focuses on the initial setup, which is the foundation for all
// other authentication features.
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
