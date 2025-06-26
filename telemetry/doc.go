// Copyright (c) 2025 A Bit of Help, Inc.

// Package telemetry provides functionality for collecting and exporting telemetry data,
// including distributed tracing and metrics.
//
// This package integrates with OpenTelemetry and Prometheus to provide a unified interface
// for telemetry collection. It supports tracing requests across service boundaries,
// collecting metrics about application performance, and exporting this data to various
// backends for analysis and visualization.
//
// Key features:
//   - Distributed tracing with OpenTelemetry
//   - Metrics collection with OpenTelemetry and Prometheus
//   - HTTP middleware for automatic request tracing and metrics
//   - Database operation tracing and metrics
//   - Error tracking and metrics
//   - Configurable exporters for different backends
//
// The package provides several main components:
//   - TelemetryProvider: A unified provider for both tracing and metrics
//   - Tracer: An interface for creating and managing trace spans
//   - Span: An interface representing a single trace span
//   - MetricsProvider: A provider for metrics collection
//
// Example usage:
//
//	// Initialize telemetry provider
//	provider, err := telemetry.NewTelemetryProvider(ctx, logger, config)
//	if err != nil {
//	    logger.Fatal(ctx, "Failed to initialize telemetry", zap.Error(err))
//	}
//	defer provider.Shutdown(ctx)
//
//	// Create a span
//	ctx, span := provider.Tracer().Start(ctx, "operation_name")
//	defer span.End()
//
//	// Add attributes to the span
//	span.SetAttributes(attribute.String("key", "value"))
//
//	// Record an error
//	if err := someOperation(); err != nil {
//	    span.RecordError(err)
//	    telemetry.RecordErrorMetric(ctx, "operation_error", "someOperation")
//	    return err
//	}
//
//	// Record HTTP metrics
//	telemetry.RecordHTTPRequest(ctx, "GET", "/api/resource", 200, duration, responseSize)
//
//	// Record database operation metrics
//	telemetry.RecordDBOperation(ctx, "query", "users_db", "users", duration, nil)
//
// The package also provides utility functions for wrapping operations with spans:
//
//	// Execute a function with a span
//	err := telemetry.WithSpan(ctx, "operation_name", func(ctx context.Context) error {
//	    // Operation code here
//	    return nil
//	})
//
//	// Execute a function with a span and measure duration
//	duration, err := telemetry.WithSpanTimed(ctx, "timed_operation", func(ctx context.Context) error {
//	    // Operation code here
//	    return nil
//	})
//
// For HTTP servers, the package provides middleware for automatic request tracing and metrics:
//
//	// Create HTTP middleware
//	middleware := provider.NewHTTPMiddleware()
//
//	// Apply middleware to HTTP handler
//	http.Handle("/", middleware(myHandler))
//
// The telemetry package is designed to be used throughout the application to provide
// comprehensive visibility into application performance and behavior.
package telemetry