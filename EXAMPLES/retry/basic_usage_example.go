// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating basic usage of the retry package
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/retry"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define a function to retry
	fn := func(ctx context.Context) error {
		// Your operation here
		// For demonstration, we'll simulate a successful operation
		fmt.Println("Executing operation...")
		return nil
	}

	// Use default retry configuration
	config := retry.DefaultConfig()

	// Execute with retry
	err := retry.Do(ctx, fn, config, nil)
	if err != nil {
		fmt.Printf("Operation failed after retries: %v\n", err)
	} else {
		fmt.Println("Operation succeeded!")
	}
}