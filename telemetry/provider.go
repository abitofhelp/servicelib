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

// TelemetryProvider is a unified provider for tracing and metrics.
// It encapsulates both metrics and tracing functionality, providing a single entry point
// for all telemetry operations in the application. The provider manages the lifecycle
// of telemetry components, including initialization and shutdown.
//
// TelemetryProvider implements a facade pattern, hiding the complexity of the underlying
// telemetry implementations and providing a simple, consistent interface for the application.
type TelemetryProvider struct {
	// metricsProvider is the metrics provider for collecting and exporting metrics
	metricsProvider *MetricsProvider

	// tracingProvider is the tracing provider for distributed tracing
	tracingProvider *TracingProvider

	// logger is the logger for the telemetry provider
	logger *logging.ContextLogger
}

// TracingProvider holds the tracing provider and tracer.
// It encapsulates the OpenTelemetry tracer and provides a way to create
// and manage trace spans. This struct is used internally by TelemetryProvider
// to handle the tracing aspect of telemetry.
type TracingProvider struct {
	// tracer is the OpenTelemetry tracer used to create spans
	tracer trace.Tracer

	// logger is the logger for the tracing provider
	logger *logging.ContextLogger
}

// NewTelemetryProvider creates a new telemetry provider.
// This function initializes both metrics and tracing components based on the provided
// configuration. If telemetry is disabled in the configuration, it returns nil without
// an error. The provider manages the lifecycle of all telemetry components and provides
// a unified interface for telemetry operations.
//
// Parameters:
//   - ctx: The context for the operation, which can be used to cancel initialization.
//   - logger: The logger to use for logging telemetry events and errors.
//   - k: The koanf instance to load configuration from. This should contain telemetry
//     configuration under the "telemetry" key.
//
// Returns:
//   - *TelemetryProvider: The initialized telemetry provider, or nil if telemetry is disabled.
//   - error: An error if the telemetry provider creation fails, such as if metrics
//     initialization fails.
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

// Shutdown shuts down the telemetry provider.
// This method gracefully shuts down all telemetry components, ensuring that
// any buffered telemetry data is flushed to its destination. It should be called
// when the application is shutting down to prevent data loss.
//
// If the telemetry provider is nil (which happens when telemetry is disabled),
// this method is a no-op and returns nil.
//
// Parameters:
//   - ctx: The context for the operation, which can be used to set a timeout
//     for the shutdown process or cancel it.
//
// Returns:
//   - error: An error if the shutdown fails, such as if the metrics provider
//     fails to shut down properly.
func (tp *TelemetryProvider) Shutdown(ctx context.Context) error {
	if tp.metricsProvider != nil {
		if err := tp.metricsProvider.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Meter returns the OpenTelemetry meter for creating and recording metrics.
// This method provides access to the meter instance that can be used to create
// various types of metrics (counters, gauges, histograms) and record measurements.
//
// If the metrics provider is nil or its meter is nil (which can happen when metrics
// are disabled), this method returns a default meter from the global OpenTelemetry provider.
//
// Returns:
//   - metric.Meter: The OpenTelemetry meter that can be used to create and record metrics.
//     This will never be nil, even if metrics are disabled.
func (tp *TelemetryProvider) Meter() metric.Meter {
	if tp.metricsProvider != nil {
		meter := tp.metricsProvider.Meter()
		if meter != nil {
			return meter
		}
	}
	return otel.Meter("github.com/abitofhelp/servicelib/telemetry")
}

// Tracer returns the OpenTelemetry tracer for creating spans.
// This method provides access to the tracer instance that can be used to create
// spans for tracing operations across the application. The tracer is used to
// create spans that track the execution of operations and record events, attributes,
// and errors.
//
// If the tracing provider is nil or its tracer is nil (which can happen when tracing
// is disabled), this method returns a default tracer from the global OpenTelemetry provider.
//
// Returns:
//   - trace.Tracer: The OpenTelemetry tracer that can be used to create spans.
//     This will never be nil, even if tracing is disabled.
func (tp *TelemetryProvider) Tracer() trace.Tracer {
	if tp.tracingProvider != nil && tp.tracingProvider.tracer != nil {
		return tp.tracingProvider.tracer
	}
	return otel.Tracer("github.com/abitofhelp/servicelib/telemetry")
}

// CreatePrometheusHandler creates a handler for the Prometheus metrics endpoint.
// This method returns an HTTP handler that exposes metrics in the Prometheus format.
// The handler can be registered with an HTTP server to provide a metrics endpoint
// that can be scraped by Prometheus.
//
// This is a convenience method that delegates to the global CreatePrometheusHandler function.
// It's provided as a method on TelemetryProvider to maintain a consistent interface
// for telemetry operations.
//
// Returns:
//   - http.Handler: The Prometheus metrics handler that can be registered with an HTTP server.
//     This handler will expose all metrics that have been registered with the OpenTelemetry
//     meter provider.
func (tp *TelemetryProvider) CreatePrometheusHandler() http.Handler {
	return CreatePrometheusHandler()
}

// InstrumentHandler wraps an http.Handler with OpenTelemetry instrumentation.
// This method adds tracing and metrics collection to an HTTP handler, allowing
// requests to be traced and metrics to be collected automatically. The instrumented
// handler will create a span for each request, record request duration, and track
// response status codes.
//
// This is a convenience method that delegates to the global InstrumentHandler function.
// It's provided as a method on TelemetryProvider to maintain a consistent interface
// for telemetry operations.
//
// Parameters:
//   - handler: The HTTP handler to instrument. This is the handler that will be wrapped
//     with telemetry instrumentation.
//   - operation: The name of the operation for tracing. This will be used as the span name
//     and will appear in traces and metrics.
//
// Returns:
//   - http.Handler: The instrumented HTTP handler that can be used in place of the original
//     handler to automatically collect telemetry data for HTTP requests.
func (tp *TelemetryProvider) InstrumentHandler(handler http.Handler, operation string) http.Handler {
	return InstrumentHandler(handler, operation)
}

// NewHTTPMiddleware creates a new middleware for HTTP request tracing and metrics.
// This method returns a middleware function that can be used to add telemetry
// instrumentation to all HTTP handlers in an application. The middleware will
// create a span for each request, record request duration, track response status
// codes, and collect other HTTP-related metrics.
//
// This middleware is particularly useful when using a router or framework that
// supports middleware, as it allows you to add telemetry to all routes with a
// single middleware registration.
//
// This is a convenience method that delegates to the global NewHTTPMiddleware function.
// It's provided as a method on TelemetryProvider to maintain a consistent interface
// for telemetry operations.
//
// Returns:
//   - func(http.Handler) http.Handler: A middleware function that takes an HTTP handler
//     and returns a new handler with telemetry instrumentation. This function can be
//     passed to router or framework middleware registration functions.
func (tp *TelemetryProvider) NewHTTPMiddleware() func(http.Handler) http.Handler {
	return NewHTTPMiddleware(tp.logger)
}

// WithSpan wraps a function with a span for tracing.
// This utility function creates a new span with the given name, executes the provided
// function with the context containing the span, and automatically ends the span when
// the function completes. It's a convenient way to add tracing to any function without
// having to manually create and end spans.
//
// The span will be created as a child of any existing span in the provided context,
// allowing for proper trace context propagation and hierarchical span relationships.
//
// Example usage:
//
//	err := telemetry.WithSpan(ctx, "database_operation", func(ctx context.Context) error {
//	    return db.ExecuteQuery(ctx, query)
//	})
//
// Parameters:
//   - ctx: The context to start the span from. This should contain any parent span.
//   - name: The name of the span, which should describe the operation being performed.
//   - fn: The function to execute within the span. This function receives the context
//     containing the span and should pass it to any operations that should be included
//     in the span.
//
// Returns:
//   - error: The error returned by the function, if any. This allows for transparent
//     error propagation from the wrapped function.
func WithSpan(ctx context.Context, name string, fn func(context.Context) error) error {
	ctx, span := StartSpan(ctx, name)
	defer span.End()

	return fn(ctx)
}

// WithSpanTimed wraps a function with a span and records the execution time.
// This utility function is similar to WithSpan, but it also measures and returns
// the execution time of the wrapped function. It creates a new span with the given name,
// executes the provided function with the context containing the span, records the
// execution time as a span attribute, and automatically ends the span when the function
// completes.
//
// The execution time is recorded as a "duration_ms" attribute on the span, making it
// visible in trace visualizations. If the function returns an error, it is also recorded
// on the span.
//
// Example usage:
//
//	duration, err := telemetry.WithSpanTimed(ctx, "expensive_operation", func(ctx context.Context) error {
//	    return performExpensiveOperation(ctx)
//	})
//	logger.Info(ctx, "Operation completed", zap.Duration("duration", duration))
//
// Parameters:
//   - ctx: The context to start the span from. This should contain any parent span.
//   - name: The name of the span, which should describe the operation being performed.
//   - fn: The function to execute within the span. This function receives the context
//     containing the span and should pass it to any operations that should be included
//     in the span.
//
// Returns:
//   - time.Duration: The execution time of the function, which can be used for logging
//     or other purposes outside the span.
//   - error: The error returned by the function, if any. This allows for transparent
//     error propagation from the wrapped function.
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
