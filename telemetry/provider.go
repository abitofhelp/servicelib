// Copyright (c) 2025 A Bit of Help, Inc.

// Package telemetry provides functionality for monitoring and tracing application behavior.
// It offers a unified interface for both metrics collection and distributed tracing
// using OpenTelemetry and Prometheus.
package telemetry

import (
	"context"
	"net/http"
	"time"

	"github.com/knadh/koanf/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/abitofhelp/servicelib/logging"
)

// TelemetryProvider is a unified provider for tracing and metrics
type TelemetryProvider struct {
	// metricsProvider is the metrics provider
	metricsProvider *MetricsProvider

	// tracingProvider is the tracing provider
	tracingProvider *TracingProvider

	// logger is the logger for the telemetry provider
	logger *logging.ContextLogger
}

// TracingProvider holds the tracing provider and tracer
type TracingProvider struct {
	// tracer is the OpenTelemetry tracer
	tracer trace.Tracer

	// logger is the logger for the tracing provider
	logger *logging.ContextLogger
}

// NewTelemetryProvider creates a new telemetry provider
//
// Parameters:
//   - ctx: The context for the operation
//   - logger: The logger to use for logging telemetry events
//   - k: The koanf instance to load configuration from
//
// Returns:
//   - *TelemetryProvider: The telemetry provider
//   - error: An error if the telemetry provider creation fails
func NewTelemetryProvider(ctx context.Context, logger *logging.ContextLogger, k *koanf.Koanf) (*TelemetryProvider, error) {
	config := LoadConfig(k)
	if !config.Enabled {
		logger.Info(ctx, "Telemetry is disabled")
		return nil, nil
	}

	// Create metrics provider
	metricsProvider, err := NewMetricsProvider(ctx, logger, k)
	if err != nil {
		return nil, err
	}

	// Create tracing provider
	tracingProvider := &TracingProvider{
		tracer: otel.Tracer("github.com/abitofhelp/servicelib/telemetry"),
		logger: logger,
	}

	logger.Info(ctx, "Telemetry provider initialized",
		zap.String("service", config.ServiceName),
		zap.String("environment", config.Environment),
		zap.Bool("metrics_enabled", config.Metrics.Enabled),
		zap.Bool("tracing_enabled", config.Tracing.Enabled),
	)

	return &TelemetryProvider{
		metricsProvider: metricsProvider,
		tracingProvider: tracingProvider,
		logger:          logger,
	}, nil
}

// Shutdown shuts down the telemetry provider
//
// Parameters:
//   - ctx: The context for the operation
//
// Returns:
//   - error: An error if the shutdown fails
func (tp *TelemetryProvider) Shutdown(ctx context.Context) error {
	if tp.metricsProvider != nil {
		if err := tp.metricsProvider.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Meter returns the meter
//
// Returns:
//   - metric.Meter: The OpenTelemetry meter
func (tp *TelemetryProvider) Meter() metric.Meter {
	if tp.metricsProvider != nil {
		return tp.metricsProvider.Meter()
	}
	return otel.Meter("github.com/abitofhelp/servicelib/telemetry")
}

// Tracer returns the tracer
//
// Returns:
//   - trace.Tracer: The OpenTelemetry tracer
func (tp *TelemetryProvider) Tracer() trace.Tracer {
	if tp.tracingProvider != nil {
		return tp.tracingProvider.tracer
	}
	return otel.Tracer("github.com/abitofhelp/servicelib/telemetry")
}

// CreatePrometheusHandler creates a handler for the Prometheus metrics endpoint
//
// Returns:
//   - http.Handler: The Prometheus metrics handler
func (tp *TelemetryProvider) CreatePrometheusHandler() http.Handler {
	return CreatePrometheusHandler()
}

// InstrumentHandler wraps an http.Handler with OpenTelemetry instrumentation
//
// Parameters:
//   - handler: The HTTP handler to instrument
//   - operation: The name of the operation for tracing
//
// Returns:
//   - http.Handler: The instrumented HTTP handler
func (tp *TelemetryProvider) InstrumentHandler(handler http.Handler, operation string) http.Handler {
	return InstrumentHandler(handler, operation)
}

// NewHTTPMiddleware creates a new middleware for HTTP request tracing
//
// Returns:
//   - func(http.Handler) http.Handler: The middleware function
func (tp *TelemetryProvider) NewHTTPMiddleware() func(http.Handler) http.Handler {
	return NewHTTPMiddleware(tp.logger)
}

// WithSpan wraps a function with a span.
// This is useful for tracing functions.
//
// Parameters:
//   - ctx: The context to start the span from
//   - name: The name of the span
//   - fn: The function to wrap
//
// Returns:
//   - error: The error returned by the function
func WithSpan(ctx context.Context, name string, fn func(context.Context) error) error {
	ctx, span := StartSpan(ctx, name)
	defer span.End()

	return fn(ctx)
}

// WithSpanTimed wraps a function with a span and records the execution time.
// This is useful for tracing functions and measuring their execution time.
//
// Parameters:
//   - ctx: The context to start the span from
//   - name: The name of the span
//   - fn: The function to wrap
//
// Returns:
//   - time.Duration: The execution time of the function
//   - error: The error returned by the function
func WithSpanTimed(ctx context.Context, name string, fn func(context.Context) error) (time.Duration, error) {
	ctx, span := StartSpan(ctx, name)
	defer span.End()

	start := time.Now()
	err := fn(ctx)
	duration := time.Since(start)

	span.SetAttributes(attribute.Float64("duration_ms", float64(duration.Milliseconds())))
	if err != nil {
		RecordErrorSpan(ctx, err)
	}

	return duration, err
}
