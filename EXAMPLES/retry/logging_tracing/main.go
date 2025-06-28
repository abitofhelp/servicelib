//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating logging and tracing with the retry package
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/retry"
	"github.com/abitofhelp/servicelib/telemetry"
	"go.uber.org/zap"
)

// simulateFailingOperation simulates an operation that fails on the first two attempts
func simulateFailingOperation(ctx context.Context, attempt int) error {
	fmt.Printf("Attempt %d: Executing operation...\n", attempt+1)

	// Simulate failure on first two attempts
	if attempt < 2 {
		return errors.New("temporary failure")
	}

	// Succeed on third attempt
	return nil
}

func main() {
	// Initialize a development logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create a context logger
	contextLogger := logging.NewContextLogger(logger)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Track the current attempt
	attempt := 0

	// Define a function to retry that uses our simulateFailingOperation function
	fn := func(ctx context.Context) error {
		err := simulateFailingOperation(ctx, attempt)
		attempt++
		return err
	}

	// Create retry configuration
	config := retry.DefaultConfig().
		WithMaxRetries(5).
		WithInitialBackoff(100 * time.Millisecond)

	// Create options with logger
	options := retry.DefaultOptions()
	options.Logger = contextLogger

	fmt.Println("Example 1: Using logger only")
	fmt.Println("----------------------------")

	// Execute with retry and logging
	err := retry.DoWithOptions(ctx, fn, config, nil, options)
	if err != nil {
		fmt.Printf("Operation failed after retries: %v\n", err)
	} else {
		fmt.Println("Operation succeeded!")
	}

	// Reset attempt counter for the next example
	attempt = 0

	fmt.Println("\nExample 2: Using logger and tracer")
	fmt.Println("----------------------------------")

	// Create a no-op tracer (this is the default)
	noopTracer := telemetry.NewNoopTracer()

	// Create options with logger and tracer
	optionsWithTracer := retry.Options{
		Logger: contextLogger,
		Tracer: noopTracer,
	}

	// Execute with retry, logging, and tracing
	err = retry.DoWithOptions(ctx, fn, config, nil, optionsWithTracer)
	if err != nil {
		fmt.Printf("Operation failed after retries: %v\n", err)
	} else {
		fmt.Println("Operation succeeded!")
	}

	fmt.Println("\nNote: In a real application, you would use a real tracer instead of a no-op tracer.")
	fmt.Println("For example, to use OpenTelemetry tracing:")
	fmt.Println("  options := retry.DefaultOptions()")
	fmt.Println("  options.Logger = contextLogger")
	fmt.Println("  options = options.WithOtelTracer(otel.Tracer(\"my-service\"))")
}