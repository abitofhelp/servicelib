// Copyright (c) 2025 A Bit of Help, Inc.

// Package context provides utilities for working with Go's standard context package.
//
// This package extends the functionality of the standard context package by providing:
//   - Context creation with various timeout configurations (default, short, long, etc.)
//   - Context enrichment with common metadata (request ID, trace ID, user ID, etc.)
//   - Helper functions for retrieving context values
//   - Context validation and error handling
//
// The package is designed to standardize context usage across the application and
// ensure that all contexts contain the necessary metadata for tracing, logging,
// and error handling.
//
// Key features:
//   - Predefined timeout constants for different types of operations
//   - Automatic generation of unique IDs (request ID, trace ID, correlation ID)
//   - Context validation to handle timeouts and cancellations gracefully
//   - Structured context information for logging and debugging
//
// Example usage:
//
//	// Create a new context with default timeout and operation name
//	ctx, cancel := context.NewContext(context.ContextOptions{
//	    Timeout:   context.DefaultTimeout,
//	    Operation: "fetch_user_data",
//	})
//	defer cancel()
//
//	// Add user ID to the context
//	ctx = context.WithUserID(ctx, "user123")
//
//	// Check if the context is still valid before proceeding
//	if err := context.CheckContext(ctx); err != nil {
//	    return err
//	}
//
//	// Get context information for logging
//	logger.Info("Processing request", zap.String("context", context.ContextInfo(ctx)))
package context