// Copyright (c) 2025 A Bit of Help, Inc.

// Package circuit provides functionality for circuit breaking on external dependencies.
//
// This package implements the circuit breaker pattern to protect against
// cascading failures when external dependencies are unavailable.
package circuit

import (
	"context"
	"sync"
	"time"

	"github.com/abitofhelp/servicelib/errors/recovery"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// State represents the state of the circuit breaker.
// It defines the three possible states that a circuit breaker can be in:
// closed (normal operation), open (failing, rejecting requests), or half-open (testing recovery).
type State int

const (
	// Closed indicates that the circuit is closed and requests are allowed through normally.
	// This is the default state and represents normal operation.
	Closed State = iota

	// Open indicates that the circuit is open and requests are immediately rejected.
	// This state is entered when the error threshold is exceeded, preventing further
	// requests to a failing component.
	Open

	// HalfOpen indicates that the circuit is allowing a limited number of requests through
	// to test if the dependency has recovered. This state is entered after the sleep window
	// has elapsed since the circuit was opened.
	HalfOpen
)

// String returns a string representation of the circuit breaker state.
// This method implements the fmt.Stringer interface, allowing State values
// to be easily formatted in log messages and error reports.
//
// Returns:
//   - A human-readable string representing the state: "Closed", "Open", "HalfOpen", or "Unknown".
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

// Config contains circuit breaker configuration parameters.
// It defines the behavior of the circuit breaker, including thresholds for
// tripping the circuit, timeouts, and recovery behavior.
type Config struct {
	// Enabled determines if the circuit breaker is enabled.
	// If set to false, the circuit breaker becomes a no-op and all requests are allowed through.
	Enabled bool

	// Timeout is the maximum time allowed for a request.
	// Requests that exceed this timeout are considered failures.
	Timeout time.Duration

	// MaxConcurrent is the maximum number of concurrent requests allowed.
	// This helps prevent resource exhaustion during high load.
	MaxConcurrent int

	// ErrorThreshold is the percentage of errors that will trip the circuit (0.0-1.0).
	// When the error rate exceeds this threshold, the circuit will open.
	// For example, 0.5 means the circuit will open when 50% or more of requests fail.
	ErrorThreshold float64

	// VolumeThreshold is the minimum number of requests before the error threshold is checked.
	// This prevents the circuit from opening due to a small number of failures.
	VolumeThreshold int

	// SleepWindow is the time to wait before allowing a single request through in half-open state.
	// After this duration has elapsed since the circuit opened, a test request will be allowed
	// to determine if the dependency has recovered.
	SleepWindow time.Duration
}

// DefaultConfig returns a default circuit breaker configuration with reasonable values.
// The default configuration includes:
//   - Enabled: true (circuit breaker is enabled)
//   - Timeout: 5 seconds (requests that take longer are considered failures)
//   - MaxConcurrent: 100 (maximum of 100 concurrent requests)
//   - ErrorThreshold: 0.5 (circuit opens when 50% or more of requests fail)
//   - VolumeThreshold: 10 (minimum of 10 requests before checking error threshold)
//   - SleepWindow: 1 second (wait 1 second before testing if dependency has recovered)
//
// Returns:
//   - A Config instance with default values.
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

// WithEnabled sets whether the circuit breaker is enabled.
// If enabled is set to false, the circuit breaker becomes a no-op and all requests are allowed through.
//
// Parameters:
//   - enabled: A boolean indicating whether the circuit breaker should be enabled.
//
// Returns:
//   - A new Config instance with the updated Enabled value.
func (c Config) WithEnabled(enabled bool) Config {
	c.Enabled = enabled
	return c
}

// WithTimeout sets the maximum time allowed for a request.
// Requests that exceed this timeout are considered failures.
// If a non-positive value is provided, it will be set to 1 millisecond.
//
// Parameters:
//   - timeout: The maximum duration allowed for a request.
//
// Returns:
//   - A new Config instance with the updated Timeout value.
func (c Config) WithTimeout(timeout time.Duration) Config {
	if timeout <= 0 {
		timeout = 1 * time.Millisecond
	}
	c.Timeout = timeout
	return c
}

// WithMaxConcurrent sets the maximum number of concurrent requests allowed.
// This helps prevent resource exhaustion during high load.
// If a non-positive value is provided, it will be set to 1.
//
// Parameters:
//   - maxConcurrent: The maximum number of concurrent requests allowed.
//
// Returns:
//   - A new Config instance with the updated MaxConcurrent value.
func (c Config) WithMaxConcurrent(maxConcurrent int) Config {
	if maxConcurrent <= 0 {
		maxConcurrent = 1
	}
	c.MaxConcurrent = maxConcurrent
	return c
}

// WithErrorThreshold sets the percentage of errors that will trip the circuit.
// When the error rate exceeds this threshold, the circuit will open.
// The value should be between 0 and 1, where 0 means any error will trip the circuit,
// and 1 means the circuit will never trip. If a value outside this range is provided,
// it will be clamped to the valid range.
//
// Parameters:
//   - errorThreshold: The error threshold as a decimal between 0 and 1.
//
// Returns:
//   - A new Config instance with the updated ErrorThreshold value.
func (c Config) WithErrorThreshold(errorThreshold float64) Config {
	if errorThreshold < 0 {
		errorThreshold = 0
	} else if errorThreshold > 1 {
		errorThreshold = 1
	}
	c.ErrorThreshold = errorThreshold
	return c
}

// WithVolumeThreshold sets the minimum number of requests before the error threshold is checked.
// This prevents the circuit from opening due to a small number of failures.
// If a non-positive value is provided, it will be set to 1.
//
// Parameters:
//   - volumeThreshold: The minimum number of requests required before checking the error threshold.
//
// Returns:
//   - A new Config instance with the updated VolumeThreshold value.
func (c Config) WithVolumeThreshold(volumeThreshold int) Config {
	if volumeThreshold <= 0 {
		volumeThreshold = 1
	}
	c.VolumeThreshold = volumeThreshold
	return c
}

// WithSleepWindow sets the time to wait before allowing a single request through in half-open state.
// After this duration has elapsed since the circuit opened, a test request will be allowed
// to determine if the dependency has recovered. If the test request succeeds, the circuit
// will close; if it fails, the circuit will remain open for another sleep window.
// If a non-positive value is provided, it will be set to 1 millisecond.
//
// Parameters:
//   - sleepWindow: The duration to wait before testing if the dependency has recovered.
//
// Returns:
//   - A new Config instance with the updated SleepWindow value.
func (c Config) WithSleepWindow(sleepWindow time.Duration) Config {
	if sleepWindow <= 0 {
		sleepWindow = 1 * time.Millisecond
	}
	c.SleepWindow = sleepWindow
	return c
}

// Options contains additional options for the circuit breaker.
// These options are not directly related to the circuit breaker behavior itself,
// but provide additional functionality like logging, tracing, and identification.
type Options struct {
	// Logger is used for logging circuit breaker operations.
	// If nil, a no-op logger will be used.
	Logger *logging.ContextLogger

	// Tracer is used for tracing circuit breaker operations.
	// It provides integration with OpenTelemetry for distributed tracing.
	Tracer telemetry.Tracer

	// Name is the name of the circuit breaker.
	// This is useful for identifying the circuit breaker in logs and traces,
	// especially when multiple circuit breakers are used in the same application.
	Name string
}

// DefaultOptions returns default options for circuit breaker operations.
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

// WithLogger sets the logger for the circuit breaker.
// The logger is used to log circuit breaker operations, such as state transitions,
// request successes and failures, and initialization events.
//
// Parameters:
//   - logger: A ContextLogger instance for logging circuit breaker operations.
//
// Returns:
//   - A new Options instance with the updated Logger value.
func (o Options) WithLogger(logger *logging.ContextLogger) Options {
	o.Logger = logger
	return o
}

// WithOtelTracer returns Options with an OpenTelemetry tracer.
// This allows users to opt-in to OpenTelemetry tracing if they need it.
// The tracer is used to create spans for circuit breaker operations, which can be
// viewed in a distributed tracing system to understand the behavior and performance
// of the circuit breaker.
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

// WithName sets the name of the circuit breaker.
// The name is used to identify the circuit breaker in logs and traces,
// which is especially useful when multiple circuit breakers are used in the same application
// to protect different resources or operations.
//
// Parameters:
//   - name: A string identifier for the circuit breaker.
//
// Returns:
//   - A new Options instance with the updated Name value.
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
