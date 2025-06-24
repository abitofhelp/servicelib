// Copyright (c) 2025 A Bit of Help, Inc.

// Package circuit provides functionality for circuit breaking on external dependencies.
//
// This package implements the circuit breaker pattern to protect against
// cascading failures when external dependencies are unavailable.
package circuit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/errors/recovery"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// State represents the state of the circuit breaker
type State int

const (
	// Closed means the circuit is closed and requests are allowed through
	Closed State = iota
	// Open means the circuit is open and requests are not allowed through
	Open
	// HalfOpen means the circuit is allowing a limited number of requests through to test if the dependency is healthy
	HalfOpen
)

// String returns a string representation of the circuit breaker state
func (s State) String() string {
	switch s {
	case Closed:
		return "Closed"
	case Open:
		return "Open"
	case HalfOpen:
		return "HalfOpen"
	default:
		return "Unknown"
	}
}

// Config contains circuit breaker configuration parameters
type Config struct {
	// Enabled determines if the circuit breaker is enabled
	Enabled bool
	// Timeout is the maximum time allowed for a request
	Timeout time.Duration
	// MaxConcurrent is the maximum number of concurrent requests allowed
	MaxConcurrent int
	// ErrorThreshold is the percentage of errors that will trip the circuit (0.0-1.0)
	ErrorThreshold float64
	// VolumeThreshold is the minimum number of requests before the error threshold is checked
	VolumeThreshold int
	// SleepWindow is the time to wait before allowing a single request through in half-open state
	SleepWindow time.Duration
}

// DefaultConfig returns a default circuit breaker configuration
func DefaultConfig() Config {
	return Config{
		Enabled:         true,
		Timeout:         5 * time.Second,
		MaxConcurrent:   100,
		ErrorThreshold:  0.5,
		VolumeThreshold: 10,
		SleepWindow:     1 * time.Second,
	}
}

// WithEnabled sets whether the circuit breaker is enabled
func (c Config) WithEnabled(enabled bool) Config {
	c.Enabled = enabled
	return c
}

// WithTimeout sets the maximum time allowed for a request
func (c Config) WithTimeout(timeout time.Duration) Config {
	if timeout <= 0 {
		timeout = 1 * time.Millisecond
	}
	c.Timeout = timeout
	return c
}

// WithMaxConcurrent sets the maximum number of concurrent requests allowed
func (c Config) WithMaxConcurrent(maxConcurrent int) Config {
	if maxConcurrent <= 0 {
		maxConcurrent = 1
	}
	c.MaxConcurrent = maxConcurrent
	return c
}

// WithErrorThreshold sets the percentage of errors that will trip the circuit
func (c Config) WithErrorThreshold(errorThreshold float64) Config {
	if errorThreshold < 0 {
		errorThreshold = 0
	} else if errorThreshold > 1 {
		errorThreshold = 1
	}
	c.ErrorThreshold = errorThreshold
	return c
}

// WithVolumeThreshold sets the minimum number of requests before the error threshold is checked
func (c Config) WithVolumeThreshold(volumeThreshold int) Config {
	if volumeThreshold <= 0 {
		volumeThreshold = 1
	}
	c.VolumeThreshold = volumeThreshold
	return c
}

// WithSleepWindow sets the time to wait before allowing a single request through in half-open state
func (c Config) WithSleepWindow(sleepWindow time.Duration) Config {
	if sleepWindow <= 0 {
		sleepWindow = 1 * time.Millisecond
	}
	c.SleepWindow = sleepWindow
	return c
}

// Options contains additional options for the circuit breaker
type Options struct {
	// Logger is used for logging circuit breaker operations
	Logger *logging.ContextLogger
	// Tracer is used for tracing circuit breaker operations
	Tracer telemetry.Tracer
	// Name is the name of the circuit breaker
	Name string
}

// DefaultOptions returns default options for circuit breaker operations
func DefaultOptions() Options {
	return Options{
		Logger: nil,
		Tracer: telemetry.NewNoopTracer(),
		Name:   "default",
	}
}

// WithLogger sets the logger for the circuit breaker
func (o Options) WithLogger(logger *logging.ContextLogger) Options {
	o.Logger = logger
	return o
}

// WithOtelTracer returns Options with an OpenTelemetry tracer
func (o Options) WithOtelTracer(tracer trace.Tracer) Options {
	o.Tracer = telemetry.NewOtelTracer(tracer)
	return o
}

// WithName sets the name of the circuit breaker
func (o Options) WithName(name string) Options {
	o.Name = name
	return o
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name            string
	config          Config
	logger          *logging.ContextLogger
	tracer          telemetry.Tracer
	state           State
	failureCount    int
	lastFailureTime time.Time
	mutex           sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config Config, options Options) *CircuitBreaker {
	if !config.Enabled {
		if options.Logger != nil {
			options.Logger.Info(context.Background(), "Circuit breaker is disabled", zap.String("name", options.Name))
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

	logger.Info(context.Background(), "Initializing circuit breaker",
		zap.String("name", options.Name),
		zap.Duration("timeout", config.Timeout),
		zap.Int("max_concurrent", config.MaxConcurrent),
		zap.Float64("error_threshold", config.ErrorThreshold),
		zap.Int("volume_threshold", config.VolumeThreshold),
		zap.Duration("sleep_window", config.SleepWindow))

	return &CircuitBreaker{
		name:            options.Name,
		config:          config,
		logger:          logger,
		tracer:          tracer,
		state:           Closed,
		failureCount:    0,
		lastFailureTime: time.Time{},
		mutex:           sync.RWMutex{},
	}
}

// Execute executes the given function with circuit breaking
// If the circuit is open, it will return an error immediately
// If the circuit is closed or half-open, it will execute the function
// and update the circuit state based on the result
func Execute[T any](ctx context.Context, cb *CircuitBreaker, operation string, fn func(ctx context.Context) (T, error)) (T, error) {
	var zero T

	if cb == nil {
		// If circuit breaker is disabled, just execute the function
		return fn(ctx)
	}

	// Create a span for the circuit breaker operation
	var span telemetry.Span
	ctx, span = cb.tracer.Start(ctx, "circuit.Execute")
	defer span.End()

	span.SetAttributes(
		attribute.String("circuit.name", cb.name),
		attribute.String("circuit.operation", operation),
		attribute.String("circuit.state", cb.GetState().String()),
	)

	// Check if the circuit is open
	if !cb.allowRequest() {
		err := recovery.ErrCircuitBreakerOpen
		cb.logger.Warn(ctx, "Circuit is open, rejecting request",
			zap.String("circuit", cb.name),
			zap.String("operation", operation))

		span.SetAttributes(attribute.String("circuit.result", "rejected"))
		span.RecordError(err)

		return zero, err
	}

	// Execute the function
	startTime := time.Now()
	result, err := fn(ctx)
	duration := time.Since(startTime)

	span.SetAttributes(
		attribute.Bool("circuit.success", err == nil),
		attribute.Int64("circuit.duration_ms", duration.Milliseconds()),
	)

	// Update the circuit state based on the result
	cb.updateState(err)

	if err != nil {
		span.RecordError(err)
	}

	return result, err
}

// ExecuteWithFallback executes the given function with circuit breaking
// If the circuit is open or the function fails, it will execute the fallback function
func ExecuteWithFallback[T any](ctx context.Context, cb *CircuitBreaker, operation string, fn func(ctx context.Context) (T, error), fallback func(ctx context.Context, err error) (T, error)) (T, error) {
	if cb == nil {
		// If circuit breaker is disabled, just execute the function
		result, err := fn(ctx)
		if err != nil {
			return fallback(ctx, err)
		}
		return result, nil
	}

	// Create a span for the circuit breaker operation
	var span telemetry.Span
	ctx, span = cb.tracer.Start(ctx, "circuit.ExecuteWithFallback")
	defer span.End()

	span.SetAttributes(
		attribute.String("circuit.name", cb.name),
		attribute.String("circuit.operation", operation),
		attribute.String("circuit.state", cb.GetState().String()),
	)

	// Check if the circuit is open
	if !cb.allowRequest() {
		err := recovery.ErrCircuitBreakerOpen
		cb.logger.Warn(ctx, "Circuit is open, using fallback",
			zap.String("circuit", cb.name),
			zap.String("operation", operation))

		span.SetAttributes(attribute.String("circuit.result", "fallback_used"))
		span.RecordError(err)

		return fallback(ctx, err)
	}

	// Execute the function
	startTime := time.Now()
	result, err := fn(ctx)
	duration := time.Since(startTime)

	span.SetAttributes(
		attribute.Bool("circuit.success", err == nil),
		attribute.Int64("circuit.duration_ms", duration.Milliseconds()),
	)

	// Update the circuit state based on the result
	cb.updateState(err)

	// If the function failed, execute the fallback
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("circuit.result", "fallback_used"))
		return fallback(ctx, err)
	}

	span.SetAttributes(attribute.String("circuit.result", "success"))
	return result, nil
}

// allowRequest checks if a request should be allowed through the circuit
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	switch cb.state {
	case Closed:
		return true
	case Open:
		// Check if the sleep window has elapsed
		if time.Since(cb.lastFailureTime) > cb.config.SleepWindow {
			// Transition to half-open state
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			cb.state = HalfOpen
			cb.mutex.Unlock()
			cb.mutex.RLock()
			return true
		}
		return false
	case HalfOpen:
		// In half-open state, allow a limited number of requests through
		return true
	default:
		return true
	}
}

// updateState updates the state of the circuit based on the result of a request
func (cb *CircuitBreaker) updateState(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err != nil {
		// Request failed
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		// Check if the circuit should be opened
		if cb.state == Closed && cb.failureCount >= cb.config.VolumeThreshold && float64(cb.failureCount)/float64(cb.config.VolumeThreshold) >= cb.config.ErrorThreshold {
			cb.logger.Warn(context.Background(), "Opening circuit breaker due to error threshold exceeded",
				zap.String("circuit", cb.name),
				zap.Int("failure_count", cb.failureCount),
				zap.Int("volume_threshold", cb.config.VolumeThreshold),
				zap.Float64("error_threshold", cb.config.ErrorThreshold))
			cb.state = Open
		} else if cb.state == HalfOpen {
			// If a request fails in half-open state, go back to open
			cb.logger.Warn(context.Background(), "Reopening circuit breaker due to failure in half-open state",
				zap.String("circuit", cb.name))
			cb.state = Open
		}
	} else {
		// Request succeeded
		if cb.state == HalfOpen {
			// If a request succeeds in half-open state, close the circuit
			cb.logger.Info(context.Background(), "Closing circuit breaker after successful request in half-open state",
				zap.String("circuit", cb.name))
			cb.state = Closed
			cb.failureCount = 0
		} else if cb.state == Closed {
			// Reset failure count on successful request
			cb.failureCount = 0
		}
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() State {
	if cb == nil {
		return Closed
	}

	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to its initial state
func (cb *CircuitBreaker) Reset() {
	if cb == nil {
		return
	}

	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.state = Closed
	cb.failureCount = 0
	cb.lastFailureTime = time.Time{}
}