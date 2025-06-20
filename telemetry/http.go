// Copyright (c) 2025 A Bit of Help, Inc.

// Package telemetry provides functionality for monitoring and tracing application behavior.
package telemetry

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/abitofhelp/servicelib/logging"
)

// InstrumentHandler wraps an http.Handler with OpenTelemetry instrumentation.
// This adds tracing to all HTTP requests handled by the provided handler.
//
// Parameters:
//   - handler: The HTTP handler to instrument
//   - operation: The name of the operation for tracing
//   - opts: Additional options for the instrumentation
//
// Returns:
//   - http.Handler: The instrumented HTTP handler
func InstrumentHandler(handler http.Handler, operation string, opts ...otelhttp.Option) http.Handler {
	// Ensure we have the default propagator to propagate trace context in headers
	defaultOpts := []otelhttp.Option{
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	}

	// Combine default options with user-provided options
	allOpts := append(defaultOpts, opts...)

	return otelhttp.NewHandler(handler, operation, allOpts...)
}

// InstrumentClient wraps an http.Client with OpenTelemetry instrumentation.
// This adds tracing to all HTTP requests made by the provided client.
//
// Parameters:
//   - client: The HTTP client to instrument
//   - opts: Additional options for the instrumentation
//
// Returns:
//   - *http.Client: The instrumented HTTP client
func InstrumentClient(client *http.Client, opts ...otelhttp.Option) *http.Client {
	if client == nil {
		client = http.DefaultClient
	}

	// Create a new client to avoid modifying the original
	newClient := &http.Client{
		CheckRedirect: client.CheckRedirect,
		Jar:           client.Jar,
		Timeout:       client.Timeout,
	}

	// Get the base transport
	baseTransport := client.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}

	// Ensure we have the default propagator to propagate trace context in headers
	defaultOpts := []otelhttp.Option{
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	}

	// Add additional default options to ensure trace headers are propagated
	defaultOpts = append(defaultOpts,
		otelhttp.WithSpanOptions(trace.WithAttributes(attribute.String("http.client", "instrumented"))),
		otelhttp.WithFilter(func(r *http.Request) bool {
			// Always trace all requests
			return true
		}),
	)

	// Combine default options with user-provided options
	allOpts := append(defaultOpts, opts...)

	// Create a new transport with OpenTelemetry instrumentation
	newClient.Transport = otelhttp.NewTransport(
		baseTransport,
		allOpts...,
	)

	return newClient
}

// NewHTTPMiddleware creates a new middleware for HTTP request tracing.
// This middleware adds tracing to all HTTP requests and logs request information.
//
// Parameters:
//   - logger: The logger to use for logging request information
//
// Returns:
//   - func(http.Handler) http.Handler: The middleware function
func NewHTTPMiddleware(logger *logging.ContextLogger) func(http.Handler) http.Handler {
	// Use the global tracer provider
	tracer := otel.Tracer("infrastructure.telemetry.http")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Start a new span for this request
			ctx, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path)
			defer span.End()

			// Add common attributes to the span
			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.host", r.Host),
				attribute.String("http.user_agent", r.UserAgent()),
			)

			// Add trace ID to response headers for debugging
			traceID := span.SpanContext().TraceID().String()
			w.Header().Set("X-Trace-ID", traceID)

			// Add trace context to request headers for propagation
			otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))

			// Log request with trace ID
			logger.Info(ctx, "HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("trace_id", traceID),
			)

			// Create a new request with the updated context
			r = r.WithContext(ctx)

			// Call the next handler with the updated request
			next.ServeHTTP(w, r)
		})
	}
}

// StartSpan is a helper function to start a new span from a context.
// This is useful for tracing operations within a request.
//
// Parameters:
//   - ctx: The context to start the span from
//   - name: The name of the span
//
// Returns:
//   - context.Context: The context with the span
//   - trace.Span: The span
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	tracer := otel.Tracer("infrastructure.telemetry")
	return tracer.Start(ctx, name)
}

// AddSpanAttributes adds attributes to the current span in context.
// This is useful for adding additional information to a span.
//
// Parameters:
//   - ctx: The context containing the span
//   - attrs: The attributes to add to the span
func AddSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// RecordErrorSpan records an error on the current span in context.
// This is useful for recording errors that occur during a traced operation.
//
// Parameters:
//   - ctx: The context containing the span
//   - err: The error to record
//   - opts: Additional options for the error event
func RecordErrorSpan(ctx context.Context, err error, opts ...trace.EventOption) {
	if err != nil {
		span := trace.SpanFromContext(ctx)
		span.RecordError(err, opts...)
	}
}
