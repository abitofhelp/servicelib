//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating custom configuration of the rate package
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/rate"
	"go.uber.org/zap"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a logger
	logger, _ := zap.NewDevelopment()
	contextLogger := logging.NewContextLogger(logger)

	// Create a rate limiter with custom configuration
	cfg := rate.DefaultConfig().
		WithEnabled(true).                // Explicitly enable the rate limiter
		WithRequestsPerSecond(5).         // Allow 5 requests per second
		WithBurstSize(3)                  // Allow a burst of 3 requests

	options := rate.DefaultOptions().
		WithName("custom-example").       // Set a custom name
		WithLogger(contextLogger)         // Use a custom logger

	rl := rate.NewRateLimiter(cfg, options)

	// Define a function to execute with rate limiting
	fn := func(ctx context.Context) (string, error) {
		// Your operation here
		// For demonstration, we'll simulate a successful operation
		fmt.Println("Executing operation with custom configuration...")
		return "success", nil
	}

	// Execute the function with rate limiting
	fmt.Println("Executing function with custom rate limiting...")
	result, err := rate.Execute(ctx, rl, "custom-operation", fn)
	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
	} else {
		fmt.Printf("Operation succeeded with result: %s\n", result)
	}

	// Demonstrate rate limiting with custom configuration
	fmt.Println("\nDemonstrating rate limiting with custom configuration...")
	fmt.Println("With burst size of 3 and 5 requests per second, the first 3 requests should succeed immediately,")
	fmt.Println("then subsequent requests should be rate limited.")

	for i := 0; i < 6; i++ {
		start := time.Now()
		result, err := rate.Execute(ctx, rl, "custom-operation", fn)
		duration := time.Since(start)
		if err != nil {
			fmt.Printf("Request %d: Rate limit exceeded after %v\n", i+1, duration)
		} else {
			fmt.Printf("Request %d: Succeeded after %v with result: %s\n", i+1, duration, result)
		}
	}

	// Demonstrate individual parameter configuration
	fmt.Println("\nDemonstrating individual parameter configuration:")
	
	// Start with default config
	cfg = rate.DefaultConfig()
	
	// Configure requests per second only
	cfg = cfg.WithRequestsPerSecond(20)
	fmt.Printf("Requests per second: %d\n", cfg.RequestsPerSecond)
	
	// Configure burst size only
	cfg = cfg.WithBurstSize(10)
	fmt.Printf("Burst size: %d\n", cfg.BurstSize)
	
	// Configure enabled only
	cfg = cfg.WithEnabled(false)
	fmt.Printf("Enabled: %t\n", cfg.Enabled)
}