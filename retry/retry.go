// Copyright (c) 2025 A Bit of Help, Inc.

// Package retry provides functionality for retrying operations with configurable backoff and jitter.
package retry

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
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

// Options contains additional options for the retry operation
type Options struct {
	// Logger is used for logging retry operations
	Logger *logging.ContextLogger
	// Tracer is used for tracing retry operations
	Tracer Tracer
}

// DefaultOptions returns default options for retry operations
func DefaultOptions() Options {
	// Use OpenTelemetry tracer if available, otherwise use no-op tracer
	return Options{
		Logger: nil,
		Tracer: NewOtelTracer(otel.Tracer("github.com/abitofhelp/servicelib/retry")),
	}
}

// Do executes the given function with retry logic using default options
func Do(ctx context.Context, fn RetryableFunc, config Config, isRetryable IsRetryableError) error {
	return DoWithOptions(ctx, fn, config, isRetryable, DefaultOptions())
}

// DoWithOptions executes the given function with retry logic and custom options
func DoWithOptions(ctx context.Context, fn RetryableFunc, config Config, isRetryable IsRetryableError, options Options) error {
	// Create a span for the entire retry operation
	var span Span
	ctx, span = options.Tracer.Start(ctx, "retry.Do")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int("retry.max_attempts", config.MaxRetries+1),
		attribute.Int64("retry.initial_backoff_ms", config.InitialBackoff.Milliseconds()),
		attribute.Int64("retry.max_backoff_ms", config.MaxBackoff.Milliseconds()),
		attribute.Float64("retry.backoff_factor", config.BackoffFactor),
		attribute.Float64("retry.jitter_factor", config.JitterFactor),
	)

	// Use the provided logger or create a no-op logger
	logger := options.Logger
	if logger == nil {
		logger = logging.NewContextLogger(zap.NewNop())
	}

	var err error
	backoff := config.InitialBackoff

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Check if context is done before each attempt
		if ctx.Err() != nil {
			ctxErr := errors.NewContextError("context cancelled or timed out during retry", ctx.Err())
			logger.Error(ctx, "Retry operation cancelled", 
				zap.Error(ctxErr),
				zap.Int("attempt", attempt),
				zap.Int("max_attempts", config.MaxRetries+1))

			span.SetAttributes(
				attribute.String("retry.result", "cancelled"),
				attribute.Int("retry.attempts", attempt),
			)

			return ctxErr
		}

		// Log the attempt
		logger.Debug(ctx, "Executing retry attempt", 
			zap.Int("attempt", attempt+1), 
			zap.Int("max_attempts", config.MaxRetries+1))

		// Create a span for this attempt
		var attemptCtx context.Context
		var attemptSpan Span
		attemptCtx, attemptSpan = options.Tracer.Start(ctx, "retry.Attempt")
		attemptSpan.SetAttributes(attribute.Int("retry.attempt", attempt+1))

		// Use the context with the attempt span
		ctx = attemptCtx

		// Execute the function
		startTime := time.Now()
		err = fn(ctx)
		duration := time.Since(startTime)

		// End the attempt span
		attemptSpan.SetAttributes(
			attribute.Bool("retry.success", err == nil),
			attribute.Int64("retry.duration_ms", duration.Milliseconds()),
		)
		if err != nil {
			attemptSpan.RecordError(err)
		}
		attemptSpan.End()

		if err == nil {
			logger.Debug(ctx, "Retry operation succeeded", 
				zap.Int("attempt", attempt+1),
				zap.Duration("duration", duration))

			span.SetAttributes(
				attribute.String("retry.result", "success"),
				attribute.Int("retry.attempts", attempt+1),
			)

			return nil // Success, no need to retry
		}

		// Log the error
		logger.Debug(ctx, "Retry attempt failed", 
			zap.Error(err),
			zap.Int("attempt", attempt+1), 
			zap.Int("max_attempts", config.MaxRetries+1),
			zap.Duration("duration", duration))

		// Check if we've reached the maximum number of retries
		if attempt == config.MaxRetries {
			retryErr := errors.NewRetryError("maximum retry attempts reached", err, attempt, config.MaxRetries)
			logger.Error(ctx, "Maximum retry attempts reached", 
				zap.Error(retryErr),
				zap.Int("attempts", attempt+1),
				zap.Int("max_attempts", config.MaxRetries+1))

			span.SetAttributes(
				attribute.String("retry.result", "max_attempts_reached"),
				attribute.Int("retry.attempts", attempt+1),
			)

			return retryErr
		}

		// Check if the error is retryable
		if isRetryable != nil && !isRetryable(err) {
			logger.Error(ctx, "Non-retryable error encountered", 
				zap.Error(err),
				zap.Int("attempt", attempt+1))

			span.SetAttributes(
				attribute.String("retry.result", "non_retryable_error"),
				attribute.Int("retry.attempts", attempt+1),
			)

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

		logger.Debug(ctx, "Waiting before next retry attempt", 
			zap.Duration("backoff", backoff),
			zap.Int("attempt", attempt+1),
			zap.Int("next_attempt", attempt+2))

		// Wait for backoff duration
		select {
		case <-ctx.Done():
			ctxErr := errors.NewContextError("context cancelled or timed out during retry backoff", ctx.Err())
			logger.Error(ctx, "Retry operation cancelled during backoff", 
				zap.Error(ctxErr),
				zap.Int("attempt", attempt+1),
				zap.Duration("backoff", backoff))

			span.SetAttributes(
				attribute.String("retry.result", "cancelled_during_backoff"),
				attribute.Int("retry.attempts", attempt+1),
			)

			return ctxErr
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
	unexpectedErr := errors.NewRetryError("retry loop exited unexpectedly", err, config.MaxRetries, config.MaxRetries)
	logger.Error(ctx, "Retry loop exited unexpectedly", zap.Error(unexpectedErr))

	span.SetAttributes(attribute.String("retry.result", "unexpected_exit"))

	return unexpectedErr
}

// IsNetworkError returns true if the error is a network-related error
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Check if the error is a known network error type from our errors package
	if errors.IsNetworkError(err) {
		return true
	}

	// Check if the error implements net.Error interface with Timeout() or Temporary() methods
	if netErr, ok := err.(interface {
		Timeout() bool
		Temporary() bool
	}); ok {
		return netErr.Temporary() || netErr.Timeout()
	}

	// Check if the error is a timeout error
	if IsTimeoutError(err) {
		return true
	}

	// Fall back to string matching for errors that don't implement specific interfaces
	// This is less reliable but provides backward compatibility
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
	if err == nil {
		return false
	}

	// Check if the error is a known timeout error from our errors package
	if errors.IsTimeout(err) {
		return true
	}

	// Check if the error is a context deadline exceeded error
	if err == context.DeadlineExceeded {
		return true
	}

	// Check if the error implements net.Error interface with Timeout() method
	if netErr, ok := err.(interface {
		Timeout() bool
	}); ok {
		return netErr.Timeout()
	}

	// Fall back to string matching for errors that don't implement specific interfaces
	// This is less reliable but provides backward compatibility
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
	if err == nil {
		return false
	}

	// Check if the error is a known transient error type from our errors package
	// such as RetryError, NetworkError, or ExternalServiceError
	if errors.IsRetryError(err) || errors.IsNetworkError(err) || errors.IsExternalServiceError(err) {
		return true
	}

	// Check if the error is a network error or timeout error
	if IsNetworkError(err) || IsTimeoutError(err) {
		return true
	}

	// Check if the error implements a Temporary() method (common in Go standard library)
	if tempErr, ok := err.(interface {
		Temporary() bool
	}); ok {
		return tempErr.Temporary()
	}

	// Fall back to string matching for errors that don't implement specific interfaces
	// This is less reliable but provides backward compatibility
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
