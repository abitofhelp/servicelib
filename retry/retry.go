// Copyright (c) 2025 A Bit of Help, Inc.

// Package retry provides functionality for retrying operations with configurable backoff and jitter.
//
// This package uses RetryError from the errors/infra package to represent errors that occur during
// retry operations. This is different from RetryableError in the errors/wrappers package, which is
// used to wrap errors that should be retried by external systems.
//
// RetryError: Used internally by this package to indicate that all retry attempts have been exhausted.
// RetryableError: Used by external systems to indicate that an error should be retried.
package retry

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Package-level random number generator for jitter calculations
// This is more efficient than creating a new one for each call to DoWithOptions
// Protected by rngMutex for concurrent access
var (
	rngMutex sync.Mutex
	rng      = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RetryableFunc is a function that can be retried.
// It takes a context.Context parameter and returns an error.
// If the function returns nil, it is considered successful.
// If it returns an error, the retry mechanism will determine whether to retry based on the IsRetryableError function.
type RetryableFunc func(ctx context.Context) error

// IsRetryableError is a function that determines if an error is retryable.
// It takes an error parameter and returns a boolean indicating whether the error should be retried.
// Return true to retry the operation, false to stop retrying and return the error.
type IsRetryableError func(err error) bool

// Config contains retry configuration parameters.
// It defines how retry operations should behave, including the number of retries,
// backoff durations, and jitter factors.
type Config struct {
	// MaxRetries is the maximum number of retry attempts that will be made.
	// The total number of attempts will be MaxRetries + 1 (including the initial attempt).
	MaxRetries int

	// InitialBackoff is the duration to wait before the first retry attempt.
	// This value will be multiplied by BackoffFactor for subsequent retries.
	InitialBackoff time.Duration

	// MaxBackoff is the maximum duration to wait between retry attempts.
	// The backoff duration will not exceed this value, regardless of the BackoffFactor.
	MaxBackoff time.Duration

	// BackoffFactor is the factor by which the backoff duration increases after each retry.
	// For example, a BackoffFactor of 2.0 will double the backoff duration after each retry.
	BackoffFactor float64

	// JitterFactor is a factor for adding random jitter to the backoff duration.
	// It should be a value between 0 and 1, where 0 means no jitter and 1 means maximum jitter.
	// Jitter helps prevent multiple retries from occurring simultaneously (thundering herd).
	JitterFactor float64

	// RetryableErrors is a list of specific errors that should be considered retryable.
	// If an error matches one of these errors (using errors.Is), it will be retried.
	RetryableErrors []error
}

// DefaultConfig returns a default retry configuration with reasonable values.
// The default configuration includes:
//   - 3 maximum retry attempts (4 total attempts including the initial one)
//   - 100ms initial backoff
//   - 2s maximum backoff
//   - 2.0 backoff factor (doubles the backoff after each retry)
//   - 0.2 jitter factor (adds up to 20% random jitter to the backoff)
//
// Returns:
//   - A Config instance with default values.
func DefaultConfig() Config {
	return Config{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     2 * time.Second,
		BackoffFactor:  2.0,
		JitterFactor:   0.2,
	}
}

// WithMaxRetries sets the maximum number of retry attempts.
// If a negative value is provided, it will be set to 0 (no retries).
//
// Parameters:
//   - maxRetries: The maximum number of retry attempts.
//
// Returns:
//   - A new Config instance with the updated MaxRetries value.
func (c Config) WithMaxRetries(maxRetries int) Config {
	// Ensure maxRetries is non-negative
	if maxRetries < 0 {
		maxRetries = 0
	}
	c.MaxRetries = maxRetries
	return c
}

// WithInitialBackoff sets the initial backoff duration.
// If a non-positive value is provided, it will be set to 1ms.
//
// Parameters:
//   - initialBackoff: The duration to wait before the first retry attempt.
//
// Returns:
//   - A new Config instance with the updated InitialBackoff value.
func (c Config) WithInitialBackoff(initialBackoff time.Duration) Config {
	// Ensure initialBackoff is positive
	if initialBackoff <= 0 {
		initialBackoff = 1 * time.Millisecond
	}
	c.InitialBackoff = initialBackoff
	return c
}

// WithMaxBackoff sets the maximum backoff duration.
// If a non-positive value is provided, it will be set to the initial backoff duration.
//
// Parameters:
//   - maxBackoff: The maximum duration to wait between retry attempts.
//
// Returns:
//   - A new Config instance with the updated MaxBackoff value.
func (c Config) WithMaxBackoff(maxBackoff time.Duration) Config {
	// Ensure maxBackoff is positive
	if maxBackoff <= 0 {
		maxBackoff = c.InitialBackoff
	}
	c.MaxBackoff = maxBackoff
	return c
}

// WithBackoffFactor sets the factor by which the backoff duration increases after each retry.
// If a non-positive value is provided, it will be set to 1.0 (no increase).
//
// Parameters:
//   - backoffFactor: The factor by which the backoff increases.
//
// Returns:
//   - A new Config instance with the updated BackoffFactor value.
func (c Config) WithBackoffFactor(backoffFactor float64) Config {
	// Ensure backoffFactor is positive
	if backoffFactor <= 0 {
		backoffFactor = 1.0
	}
	c.BackoffFactor = backoffFactor
	return c
}

// WithJitterFactor sets the factor for adding random jitter to the backoff duration.
// The value should be between 0 and 1, where 0 means no jitter and 1 means maximum jitter.
// If a value outside this range is provided, it will be clamped to the valid range.
//
// Parameters:
//   - jitterFactor: The factor for random jitter (0-1).
//
// Returns:
//   - A new Config instance with the updated JitterFactor value.
func (c Config) WithJitterFactor(jitterFactor float64) Config {
	// Ensure jitterFactor is between 0 and 1
	if jitterFactor < 0 {
		jitterFactor = 0
	} else if jitterFactor > 1 {
		jitterFactor = 1
	}
	c.JitterFactor = jitterFactor
	return c
}

// WithRetryableErrors sets the list of specific errors that should be considered retryable.
// If an error matches one of these errors (using errors.Is), it will be retried.
//
// Parameters:
//   - retryableErrors: A slice of errors that should be considered retryable.
//
// Returns:
//   - A new Config instance with the updated RetryableErrors value.
func (c Config) WithRetryableErrors(retryableErrors []error) Config {
	c.RetryableErrors = retryableErrors
	return c
}

// Options contains additional options for the retry operation.
// These options are not directly related to the retry behavior itself,
// but provide additional functionality like logging and tracing.
type Options struct {
	// Logger is used for logging retry operations.
	// If nil, a no-op logger will be used.
	Logger *logging.ContextLogger

	// Tracer is used for tracing retry operations.
	// It provides integration with OpenTelemetry for distributed tracing.
	Tracer telemetry.Tracer
}

// DefaultOptions returns default options for retry operations.
// The default options include:
//   - No logger (a no-op logger will be used)
//   - A no-op tracer (no OpenTelemetry integration)
//
// Returns:
//   - An Options instance with default values.
func DefaultOptions() Options {
	// Use no-op tracer by default to make OpenTelemetry dependency optional
	return Options{
		Logger: nil,
		Tracer: telemetry.NewNoopTracer(),
	}
}

// WithOtelTracer returns Options with an OpenTelemetry tracer.
// This allows users to opt-in to OpenTelemetry tracing if they need it.
//
// Parameters:
//   - tracer: An OpenTelemetry trace.Tracer instance.
//
// Returns:
//   - A new Options instance with the provided OpenTelemetry tracer.
func (o Options) WithOtelTracer(tracer trace.Tracer) Options {
	o.Tracer = telemetry.NewOtelTracer(tracer)
	return o
}

// Do executes the given function with retry logic using default options.
// This is a convenience wrapper around DoWithOptions that uses DefaultOptions().
//
// Parameters:
//   - ctx: The context for the operation. Can be used to cancel the retry operation.
//   - fn: The function to execute with retry logic.
//   - config: The retry configuration parameters.
//   - isRetryable: A function that determines if an error is retryable.
//
// Returns:
//   - The error from the last execution of the function, or nil if successful.
func Do(ctx context.Context, fn RetryableFunc, config Config, isRetryable IsRetryableError) error {
	return DoWithOptions(ctx, fn, config, isRetryable, DefaultOptions())
}

// DoWithOptions executes the given function with retry logic and custom options.
// It will retry the function according to the provided configuration and retry condition.
// The function will be retried if it returns an error and the isRetryable function returns true.
// The retry operation can be canceled by canceling the context.
//
// Parameters:
//   - ctx: The context for the operation. Can be used to cancel the retry operation.
//   - fn: The function to execute with retry logic.
//   - config: The retry configuration parameters.
//   - isRetryable: A function that determines if an error is retryable.
//   - options: Additional options for the retry operation, such as logging and tracing.
//
// Returns:
//   - The error from the last execution of the function, or nil if successful.
func DoWithOptions(ctx context.Context, fn RetryableFunc, config Config, isRetryable IsRetryableError, options Options) error {
	// Create a span for the entire retry operation
	var span telemetry.Span
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
		var attemptSpan telemetry.Span
		attemptCtx, attemptSpan = options.Tracer.Start(ctx, "retry.Attempt")
		attemptSpan.SetAttributes(attribute.Int("retry.attempt", attempt+1))

		// Use a separate variable for the attempt context to avoid modifying the input parameter
		currentCtx := attemptCtx

		// Execute the function
		startTime := time.Now()
		err = fn(currentCtx)
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

		// Calculate jitter factor between (1-jitterFactor) and (1+jitterFactor)
		// This ensures backoff is always positive and properly distributed
		// Use mutex to protect access to the random number generator
		rngMutex.Lock()
		jitterValue := rng.Float64()
		rngMutex.Unlock()
		jitterMultiplier := 1.0 + (jitterValue*2.0-1.0)*config.JitterFactor

		// Apply jitter by multiplying the backoff by the jitter factor
		backoff = time.Duration(float64(backoff) * jitterMultiplier)

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

// IsNetworkError checks if an error is a network-related error.
//
// Deprecated: Use errors.IsNetworkError instead.
// This function is maintained for backward compatibility and will be removed in a future version.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is a network-related error, false otherwise.
func IsNetworkError(err error) bool {
	return errors.IsNetworkError(err)
}

// IsTimeoutError checks if an error is a timeout error.
//
// Deprecated: Use errors.IsTimeout instead.
// This function is maintained for backward compatibility and will be removed in a future version.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error indicates a timeout, false otherwise.
func IsTimeoutError(err error) bool {
	return errors.IsTimeout(err)
}

// IsTransientError checks if an error is a transient error that may be resolved by retrying.
//
// Deprecated: Use errors.IsTransientError instead.
// This function is maintained for backward compatibility and will be removed in a future version.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is likely transient and may be resolved by retrying, false otherwise.
func IsTransientError(err error) bool {
	return errors.IsTransientError(err)
}
