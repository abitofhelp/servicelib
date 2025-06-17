// Copyright (c) 2025 A Bit of Help, Inc.

// Package tracing provides functionality for distributed tracing in applications.
package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	return otelhttp.NewHandler(handler, operation, opts...)
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

	client.Transport = otelhttp.NewTransport(
		client.Transport,
		opts...,
	)

	return client
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
	tracer := otel.Tracer("github.com/abitofhelp/servicelib/tracing.http")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			// Log request with trace ID
			logger.Info(ctx, "HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("trace_id", traceID),
			)

			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
