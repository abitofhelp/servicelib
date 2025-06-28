//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating basic usage of the rate package
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/rate"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a rate limiter with default configuration
	cfg := rate.DefaultConfig()
	options := rate.DefaultOptions().WithName("example")
	rl := rate.NewRateLimiter(cfg, options)

	// Define a function to execute with rate limiting
	fn := func(ctx context.Context) (string, error) {
		// Your operation here
		// For demonstration, we'll simulate a successful operation
		fmt.Println("Executing operation...")
		return "success", nil
	}

	// Execute the function with rate limiting
	fmt.Println("Executing function with rate limiting...")
	result, err := rate.Execute(ctx, rl, "example-operation", fn)
	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
	} else {
		fmt.Printf("Operation succeeded with result: %s\n", result)
	}

	// Demonstrate rate limiting by executing multiple requests in quick succession
	fmt.Println("\nDemonstrating rate limiting...")
	for i := 0; i < 10; i++ {
		result, err := rate.Execute(ctx, rl, "example-operation", fn)
		if err != nil {
			fmt.Printf("Request %d: Rate limit exceeded\n", i+1)
		} else {
			fmt.Printf("Request %d: Succeeded with result: %s\n", i+1, result)
		}
	}

	// Demonstrate waiting for tokens
	fmt.Println("\nDemonstrating waiting for tokens...")
	for i := 0; i < 3; i++ {
		start := time.Now()
		result, err := rate.ExecuteWithWait(ctx, rl, "example-operation", fn)
		duration := time.Since(start)
		if err != nil {
			fmt.Printf("Request %d: Failed after %v: %v\n", i+1, duration, err)
		} else {
			fmt.Printf("Request %d: Succeeded after %v with result: %s\n", i+1, duration, result)
		}
	}
}