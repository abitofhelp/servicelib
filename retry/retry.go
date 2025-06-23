// Copyright (c) 2025 A Bit of Help, Inc.

// Package retry provides functionality for retrying operations with configurable backoff and jitter.
package retry

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/abitofhelp/servicelib/errors"
)

// RetryableFunc is a function that can be retried
type RetryableFunc func(ctx context.Context) error

// IsRetryableError is a function that determines if an error is retryable
type IsRetryableError func(err error) bool

// Config contains retry configuration parameters
type Config struct {
	MaxRetries      int           // Maximum number of retry attempts
	InitialBackoff  time.Duration // Initial backoff duration
	MaxBackoff      time.Duration // Maximum backoff duration
	BackoffFactor   float64       // Factor by which the backoff increases
	JitterFactor    float64       // Factor for random jitter (0-1)
	RetryableErrors []error       // Errors that are considered retryable
}

// DefaultConfig returns a default retry configuration
func DefaultConfig() Config {
	return Config{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     2 * time.Second,
		BackoffFactor:  2.0,
		JitterFactor:   0.2,
	}
}

// WithMaxRetries sets the maximum number of retry attempts
func (c Config) WithMaxRetries(maxRetries int) Config {
	c.MaxRetries = maxRetries
	return c
}

// WithInitialBackoff sets the initial backoff duration
func (c Config) WithInitialBackoff(initialBackoff time.Duration) Config {
	c.InitialBackoff = initialBackoff
	return c
}

// WithMaxBackoff sets the maximum backoff duration
func (c Config) WithMaxBackoff(maxBackoff time.Duration) Config {
	c.MaxBackoff = maxBackoff
	return c
}

// WithBackoffFactor sets the factor by which the backoff increases
func (c Config) WithBackoffFactor(backoffFactor float64) Config {
	c.BackoffFactor = backoffFactor
	return c
}

// WithJitterFactor sets the factor for random jitter
func (c Config) WithJitterFactor(jitterFactor float64) Config {
	c.JitterFactor = jitterFactor
	return c
}

// WithRetryableErrors sets the errors that are considered retryable
func (c Config) WithRetryableErrors(retryableErrors []error) Config {
	c.RetryableErrors = retryableErrors
	return c
}

// Do executes the given function with retry logic
func Do(ctx context.Context, fn RetryableFunc, config Config, isRetryable IsRetryableError) error {
	var err error
	backoff := config.InitialBackoff

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Check if context is done before each attempt
		if ctx.Err() != nil {
			return errors.NewContextError("context cancelled or timed out during retry", ctx.Err())
		}

		// Execute the function
		err = fn(ctx)
		if err == nil {
			return nil // Success, no need to retry
		}

		// Check if we've reached the maximum number of retries
		if attempt == config.MaxRetries {
			return errors.NewRetryError("maximum retry attempts reached", err, attempt, config.MaxRetries)
		}

		// Check if the error is retryable
		if isRetryable != nil && !isRetryable(err) {
			return err // Non-retryable error, return immediately
		}

		// Calculate jitter
		jitter := time.Duration(rand.Float64() * config.JitterFactor * float64(backoff))

		// Apply jitter (randomly add or subtract)
		if rand.Intn(2) == 0 {
			backoff = backoff + jitter
		} else {
			backoff = backoff - jitter
		}

		// Wait for backoff duration
		select {
		case <-ctx.Done():
			return errors.NewContextError("context cancelled or timed out during retry backoff", ctx.Err())
		case <-time.After(backoff):
			// Continue with next attempt
		}

		// Increase backoff for next attempt
		backoff = time.Duration(float64(backoff) * config.BackoffFactor)
		if backoff > config.MaxBackoff {
			backoff = config.MaxBackoff
		}
	}

	// This should never happen due to the return in the loop, but just in case
	return errors.NewRetryError("retry loop exited unexpectedly", err, config.MaxRetries, config.MaxRetries)
}

// IsNetworkError returns true if the error is a network-related error
func IsNetworkError(err error) bool {
	// This is a simplified implementation
	// In a real-world scenario, you would check for specific network error types
	// such as net.Error, net.OpError, etc.
	if err == nil {
		return false
	}

	// Check if the error is a timeout error
	if IsTimeoutError(err) {
		return true
	}

	// Check if the error message contains common network error strings
	errMsg := err.Error()
	networkErrorStrings := []string{
		"connection refused",
		"connection reset",
		"connection closed",
		"no route to host",
		"network is unreachable",
		"broken pipe",
	}

	for _, s := range networkErrorStrings {
		if strings.Contains(strings.ToLower(errMsg), s) {
			return true
		}
	}

	return false
}

// IsTimeoutError returns true if the error is a timeout error
func IsTimeoutError(err error) bool {
	// This is a simplified implementation
	// In a real-world scenario, you would check for specific timeout error types
	if err == nil {
		return false
	}

	// Check if the error is a context timeout
	if errors.IsTimeout(err) {
		return true
	}

	// Check if the error message contains common timeout error strings
	errMsg := err.Error()
	timeoutErrorStrings := []string{
		"timeout",
		"timed out",
		"deadline exceeded",
	}

	for _, s := range timeoutErrorStrings {
		if strings.Contains(strings.ToLower(errMsg), s) {
			return true
		}
	}

	return false
}

// IsTransientError returns true if the error is a transient error
func IsTransientError(err error) bool {
	// This is a simplified implementation
	// In a real-world scenario, you would check for specific transient error types
	if err == nil {
		return false
	}

	// Check if the error is a network error or timeout error
	if IsNetworkError(err) || IsTimeoutError(err) {
		return true
	}

	// Check if the error message contains common transient error strings
	errMsg := err.Error()
	transientErrorStrings := []string{
		"server is busy",
		"too many requests",
		"rate limit exceeded",
		"try again",
		"temporary",
		"transient",
	}

	for _, s := range transientErrorStrings {
		if strings.Contains(strings.ToLower(errMsg), s) {
			return true
		}
	}

	return false
}
