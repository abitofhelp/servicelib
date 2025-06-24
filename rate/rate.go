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

// Config contains rate limiter configuration parameters
type Config struct {
	// Enabled determines if the rate limiter is enabled
	Enabled bool
	// RequestsPerSecond is the number of requests allowed per second
	RequestsPerSecond int
	// BurstSize is the maximum number of requests allowed in a burst
	BurstSize int
}

// DefaultConfig returns a default rate limiter configuration
func DefaultConfig() Config {
	return Config{
		Enabled:           true,
		RequestsPerSecond: 100,
		BurstSize:         50,
	}
}

// WithEnabled sets whether the rate limiter is enabled
func (c Config) WithEnabled(enabled bool) Config {
	c.Enabled = enabled
	return c
}

// WithRequestsPerSecond sets the number of requests allowed per second
func (c Config) WithRequestsPerSecond(requestsPerSecond int) Config {
	if requestsPerSecond <= 0 {
		requestsPerSecond = 1
	}
	c.RequestsPerSecond = requestsPerSecond
	return c
}

// WithBurstSize sets the maximum number of requests allowed in a burst
func (c Config) WithBurstSize(burstSize int) Config {
	if burstSize <= 0 {
		burstSize = 1
	}
	c.BurstSize = burstSize
	return c
}

// Options contains additional options for the rate limiter
type Options struct {
	// Logger is used for logging rate limiter operations
	Logger *logging.ContextLogger
	// Tracer is used for tracing rate limiter operations
	Tracer telemetry.Tracer
	// Name is the name of the rate limiter
	Name string
}

// DefaultOptions returns default options for rate limiter operations
func DefaultOptions() Options {
	return Options{
		Logger: nil,
		Tracer: telemetry.NewNoopTracer(),
		Name:   "default",
	}
}

// WithLogger sets the logger for the rate limiter
func (o Options) WithLogger(logger *logging.ContextLogger) Options {
	o.Logger = logger
	return o
}

// WithOtelTracer returns Options with an OpenTelemetry tracer
func (o Options) WithOtelTracer(tracer trace.Tracer) Options {
	o.Tracer = telemetry.NewOtelTracer(tracer)
	return o
}

// WithName sets the name of the rate limiter
func (o Options) WithName(name string) Options {
	o.Name = name
	return o
}

// RateLimiter implements a token bucket rate limiter to protect resources
// from being overwhelmed by too many requests.
type RateLimiter struct {
	name           string
	config         Config
	logger         *logging.ContextLogger
	tracer         telemetry.Tracer
	tokens         int
	lastRefillTime time.Time
	mutex          sync.Mutex
}

// NewRateLimiter creates a new rate limiter
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

// Allow checks if a request should be allowed based on the rate limit
// It returns true if the request is allowed, false otherwise
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

// Execute executes the given function with rate limiting
// If the rate limit is exceeded, it will return an error immediately
// Otherwise, it will execute the function
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

// ExecuteWithWait executes the given function with rate limiting
// If the rate limit is exceeded, it will wait until a token is available
// and then execute the function
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

// Reset resets the rate limiter to its initial state
func (rl *RateLimiter) Reset() {
	if rl == nil {
		return
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.tokens = rl.config.BurstSize
	rl.lastRefillTime = time.Now()
}
