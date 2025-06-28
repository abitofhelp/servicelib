// Copyright (c) 2025 A Bit of Help, Inc.

// Package retry provides functionality for retrying operations with configurable backoff and jitter.
//
// This package implements a flexible retry mechanism that can be used to retry operations
// that may fail due to transient errors. It supports:
//   - Configurable maximum retry attempts
//   - Exponential backoff with configurable initial and maximum durations
//   - Jitter to prevent thundering herd problems
//   - Customizable retry conditions
//   - Integration with OpenTelemetry for tracing
//   - Comprehensive logging of retry attempts
//
// The package distinguishes between two types of errors related to retries:
//   - RetryError: Used internally by this package to indicate that all retry attempts have been exhausted.
//   - RetryableError: Used by external systems to indicate that an error should be retried.
//
// Example usage:
//
//	// Define a function to retry
//	retryableFunc := func(ctx context.Context) error {
//	    // Some operation that might fail transiently
//	    return callExternalService(ctx)
//	}
//
//	// Define a function to determine if an error is retryable
//	isRetryable := func(err error) bool {
//	    return errors.IsTransientError(err)
//	}
//
//	// Configure retry parameters
//	config := retry.DefaultConfig().
//	    WithMaxRetries(5).
//	    WithInitialBackoff(100 * time.Millisecond).
//	    WithMaxBackoff(2 * time.Second)
//
//	// Execute with retry
//	err := retry.Do(ctx, retryableFunc, config, isRetryable)
//	if err != nil {
//	    // Handle error after all retries have failed
//	}
package retry
