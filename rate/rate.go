// Copyright (c) 2025 A Bit of Help, Inc.

// Package rate provides functionality for rate limiting to protect resources.
//
// This package implements a token bucket rate limiter to protect resources
// from being overwhelmed by too many requests.
package rate

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Config contains rate limiter configuration parameters.
// It defines the behavior of the rate limiter, including whether it's enabled,
// how many requests are allowed per second, and the maximum burst size.
type Config struct {
	// Enabled determines if the rate limiter is enabled.
	// If set to false, all requests will be allowed without rate limiting.
	Enabled bool

	// RequestsPerSecond is the number of requests allowed per second.
	// This determines the rate at which tokens are added to the bucket.
	RequestsPerSecond int

	// BurstSize is the maximum number of requests allowed in a burst.
	// This determines the maximum capacity of the token bucket.
	BurstSize int
}

// DefaultConfig returns a default rate limiter configuration.
// The default configuration includes:
//   - Enabled: true (rate limiting is enabled)
//   - RequestsPerSecond: 100 (100 requests allowed per second)
//   - BurstSize: 50 (maximum of 50 requests allowed in a burst)
//
// Returns:
//   - A Config instance with default values.
func DefaultConfig() Config {
	return Config{
		Enabled:           true,
		RequestsPerSecond: 100,
		BurstSize:         50,
	}
}

// WithEnabled sets whether the rate limiter is enabled.
// If enabled is set to false, all requests will be allowed without rate limiting.
//
// Parameters:
//   - enabled: A boolean indicating whether the rate limiter should be enabled.
//
// Returns:
//   - A new Config instance with the updated Enabled value.
func (c Config) WithEnabled(enabled bool) Config {
	c.Enabled = enabled
	return c
}

// WithRequestsPerSecond sets the number of requests allowed per second.
// This determines the rate at which tokens are added to the bucket.
// If a non-positive value is provided, it will be set to 1.
//
// Parameters:
//   - requestsPerSecond: The number of requests allowed per second.
//
// Returns:
//   - A new Config instance with the updated RequestsPerSecond value.
func (c Config) WithRequestsPerSecond(requestsPerSecond int) Config {
	if requestsPerSecond <= 0 {
		requestsPerSecond = 1
	}
	c.RequestsPerSecond = requestsPerSecond
	return c
}

// WithBurstSize sets the maximum number of requests allowed in a burst.
// This determines the maximum capacity of the token bucket.
// If a non-positive value is provided, it will be set to 1.
//
// Parameters:
//   - burstSize: The maximum number of requests allowed in a burst.
//
// Returns:
//   - A new Config instance with the updated BurstSize value.
func (c Config) WithBurstSize(burstSize int) Config {
	if burstSize <= 0 {
		burstSize = 1
	}
	c.BurstSize = burstSize
	return c
}

// Options contains additional options for the rate limiter.
// These options are not directly related to the rate limiting behavior itself,
// but provide additional functionality like logging, tracing, and identification.
type Options struct {
	// Logger is used for logging rate limiter operations.
	// If nil, a no-op logger will be used.
	Logger *logging.ContextLogger

	// Tracer is used for tracing rate limiter operations.
	// It provides integration with OpenTelemetry for distributed tracing.
	Tracer telemetry.Tracer

	// Name is the name of the rate limiter.
	// This is useful for identifying the rate limiter in logs and traces.
	Name string
}

// DefaultOptions returns default options for rate limiter operations.
// The default options include:
//   - No logger (a no-op logger will be used)
//   - A no-op tracer (no OpenTelemetry integration)
//   - Name: "default"
//
// Returns:
//   - An Options instance with default values.
func DefaultOptions() Options {
	return Options{
		Logger: nil,
		Tracer: telemetry.NewNoopTracer(),
		Name:   "default",
	}
}

// WithLogger sets the logger for the rate limiter.
// The logger is used to log rate limiter operations, such as when requests are
// allowed or rejected due to rate limiting.
//
// Parameters:
//   - logger: A ContextLogger instance for logging rate limiter operations.
//
// Returns:
//   - A new Options instance with the updated Logger value.
func (o Options) WithLogger(logger *logging.ContextLogger) Options {
	o.Logger = logger
	return o
}

// WithOtelTracer returns Options with an OpenTelemetry tracer.
// This allows users to opt-in to OpenTelemetry tracing if they need it.
// The tracer is used to create spans for rate limiter operations, which can be
// viewed in a distributed tracing system.
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

// WithName sets the name of the rate limiter.
// The name is used to identify the rate limiter in logs and traces,
// which is especially useful when multiple rate limiters are used in the same application.
//
// Parameters:
//   - name: A string identifier for the rate limiter.
//
// Returns:
//   - A new Options instance with the updated Name value.
func (o Options) WithName(name string) Options {
	o.Name = name
	return o
}

// RateLimiter implements a token bucket rate limiter to protect resources
// from being overwhelmed by too many requests.
//
// The token bucket algorithm works by maintaining a bucket of tokens that are
// added at a fixed rate (RequestsPerSecond). Each request consumes one token.
// If tokens are available, the request is allowed; otherwise, it is rejected
// or delayed until tokens become available.
//
// This implementation is thread-safe and can be used concurrently from multiple
// goroutines.
type RateLimiter struct {
	name           string
	config         Config
	logger         *logging.ContextLogger
	tracer         telemetry.Tracer
	tokens         int
	lastRefillTime time.Time
	mutex          sync.Mutex
}

// NewRateLimiter creates a new rate limiter with the specified configuration and options.
// If the rate limiter is disabled (config.Enabled is false), a special no-op rate limiter
// is returned that allows all requests without rate limiting.
//
// Parameters:
//   - config: The configuration parameters for the rate limiter.
//   - options: Additional options for the rate limiter, such as logging and tracing.
//
// Returns:
//   - A new RateLimiter instance configured according to the provided parameters.
func NewRateLimiter(config Config, options Options) *RateLimiter {
	if !config.Enabled {
		if options.Logger != nil {
			options.Logger.Info(context.Background(), "Rate limiter is disabled", zap.String("name", options.Name))
		}
		return nil
	}

	// Use the provided logger or create a no-op logger
	logger := options.Logger
	if logger == nil {
		logger = logging.NewContextLogger(zap.NewNop())
	}

	// Use the provided tracer or create a no-op tracer
	tracer := options.Tracer
	if tracer == nil {
		tracer = telemetry.NewNoopTracer()
	}

	logger.Info(context.Background(), "Initializing rate limiter",
		zap.String("name", options.Name),
		zap.Int("requests_per_second", config.RequestsPerSecond),
		zap.Int("burst_size", config.BurstSize))

	return &RateLimiter{
		name:           options.Name,
		config:         config,
		logger:         logger,
		tracer:         tracer,
		tokens:         config.BurstSize,
		lastRefillTime: time.Now(),
		mutex:          sync.Mutex{},
	}
}

// Allow checks if a request should be allowed based on the rate limit.
// This is a non-blocking method that immediately returns whether the request
// is allowed or not. If the rate limiter is disabled or nil, all requests are allowed.
//
// The method is thread-safe and can be called concurrently from multiple goroutines.
//
// Returns:
//   - true if the request is allowed (a token was available or rate limiting is disabled).
//   - false if the request is not allowed (no tokens were available).
func (rl *RateLimiter) Allow() bool {
	if rl == nil {
		// If rate limiter is disabled, allow all requests
		return true
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Refill tokens based on time elapsed since last refill
	rl.refillTokens()

	// Check if we have tokens available
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// Execute executes a function with rate limiting.
// This is a generic function that works with any return type. If the rate limit
// is exceeded, it returns an error without executing the function. If the rate
// limit is not exceeded, it executes the function and returns its result.
//
// This is a non-blocking operation - it will not wait for tokens to become available.
//
// Type Parameters:
//   - T: The return type of the function to execute.
//
// Parameters:
//   - ctx: The context for the operation. Can be used to cancel the operation.
//   - rl: The rate limiter to use. If nil, the function is executed without rate limiting.
//   - operation: A name for the operation being performed, used in logs and traces.
//   - fn: The function to execute if the rate limit is not exceeded.
//
// Returns:
//   - The result of the function execution, or the zero value of T if the rate limit is exceeded.
//   - An error if the rate limit is exceeded or if the function returns an error.
func Execute[T any](ctx context.Context, rl *RateLimiter, operation string, fn func(ctx context.Context) (T, error)) (T, error) {
	var zero T

	if rl == nil {
		// If rate limiter is disabled, just execute the function
		return fn(ctx)
	}

	// Create a span for the rate limiter operation
	var span telemetry.Span
	ctx, span = rl.tracer.Start(ctx, "rate.Execute")
	defer span.End()

	span.SetAttributes(
		attribute.String("rate_limiter.name", rl.name),
		attribute.String("rate_limiter.operation", operation),
		attribute.Int("rate_limiter.requests_per_second", rl.config.RequestsPerSecond),
		attribute.Int("rate_limiter.burst_size", rl.config.BurstSize),
	)

	// Check if the request is allowed
	if !rl.Allow() {
		err := errors.New(errors.ResourceExhaustedCode, fmt.Sprintf("rate limit exceeded for %s", rl.name))
		rl.logger.Warn(ctx, "Rate limit exceeded, rejecting request",
			zap.String("rate_limiter", rl.name),
			zap.String("operation", operation))

		span.SetAttributes(attribute.String("rate_limiter.result", "rejected"))
		span.RecordError(err)

		return zero, err
	}

	span.SetAttributes(attribute.String("rate_limiter.result", "allowed"))

	// Execute the function
	startTime := time.Now()
	result, err := fn(ctx)
	duration := time.Since(startTime)

	span.SetAttributes(
		attribute.Bool("rate_limiter.success", err == nil),
		attribute.Int64("rate_limiter.duration_ms", duration.Milliseconds()),
	)

	if err != nil {
		span.RecordError(err)
	}

	return result, err
}

// ExecuteWithWait executes a function with rate limiting, waiting if necessary.
// This is a generic function that works with any return type. If the rate limit
// is exceeded, it will wait until a token becomes available before executing the function.
// If the context is canceled while waiting, it returns an error without executing the function.
//
// This is a blocking operation - it will wait for tokens to become available.
//
// Type Parameters:
//   - T: The return type of the function to execute.
//
// Parameters:
//   - ctx: The context for the operation. Can be used to cancel the operation.
//   - rl: The rate limiter to use. If nil, the function is executed without rate limiting.
//   - operation: A name for the operation being performed, used in logs and traces.
//   - fn: The function to execute when a token becomes available.
//
// Returns:
//   - The result of the function execution, or the zero value of T if the context is canceled.
//   - An error if the context is canceled or if the function returns an error.
func ExecuteWithWait[T any](ctx context.Context, rl *RateLimiter, operation string, fn func(ctx context.Context) (T, error)) (T, error) {
	var zero T

	if rl == nil {
		// If rate limiter is disabled, just execute the function
		return fn(ctx)
	}

	// Create a span for the rate limiter operation
	var span telemetry.Span
	ctx, span = rl.tracer.Start(ctx, "rate.ExecuteWithWait")
	defer span.End()

	span.SetAttributes(
		attribute.String("rate_limiter.name", rl.name),
		attribute.String("rate_limiter.operation", operation),
		attribute.Int("rate_limiter.requests_per_second", rl.config.RequestsPerSecond),
		attribute.Int("rate_limiter.burst_size", rl.config.BurstSize),
	)

	// Wait until a token is available or context is canceled
	waitStartTime := time.Now()
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			span.SetAttributes(
				attribute.String("rate_limiter.result", "context_canceled"),
				attribute.Int64("rate_limiter.wait_ms", time.Since(waitStartTime).Milliseconds()),
			)
			span.RecordError(err)
			return zero, err
		default:
			if rl.Allow() {
				waitDuration := time.Since(waitStartTime)
				span.SetAttributes(
					attribute.String("rate_limiter.result", "allowed_after_wait"),
					attribute.Int64("rate_limiter.wait_ms", waitDuration.Milliseconds()),
				)

				// Execute the function
				startTime := time.Now()
				result, err := fn(ctx)
				duration := time.Since(startTime)

				span.SetAttributes(
					attribute.Bool("rate_limiter.success", err == nil),
					attribute.Int64("rate_limiter.duration_ms", duration.Milliseconds()),
				)

				if err != nil {
					span.RecordError(err)
				}

				return result, err
			}
			// Wait a bit before trying again
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// refillTokens refills tokens based on time elapsed since last refill
func (rl *RateLimiter) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefillTime)
	tokensToAdd := int(elapsed.Seconds() * float64(rl.config.RequestsPerSecond))

	if tokensToAdd > 0 {
		rl.tokens = min(rl.tokens+tokensToAdd, rl.config.BurstSize)
		rl.lastRefillTime = now
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Reset resets the rate limiter to its initial state.
// This method refills the token bucket to its maximum capacity (BurstSize)
// and resets the last refill time to the current time. This is useful for
// testing or when you want to clear any rate limiting history.
//
// The method is thread-safe and can be called concurrently from multiple goroutines.
// If the rate limiter is nil, this method does nothing.
func (rl *RateLimiter) Reset() {
	if rl == nil {
		return
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.tokens = rl.config.BurstSize
	rl.lastRefillTime = time.Now()
}
