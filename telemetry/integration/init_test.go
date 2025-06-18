// Copyright (c) 2025 A Bit of Help, Inc.

//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	// Initialize the OpenTelemetry SDK with a text map propagator
	// This is required for trace context propagation to work correctly
	
	// Create a tracer provider
	tp := sdktrace.NewTracerProvider()
	
	// Set the global tracer provider
	otel.SetTracerProvider(tp)
	
	// Set the global propagator to propagate trace context in headers
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

// InitTestTracer initializes a test tracer and returns a cleanup function
func InitTestTracer(t *testing.T) (trace.Tracer, func()) {
	// Create a tracer
	tracer := otel.Tracer("test")
	
	// Return the tracer and a cleanup function
	return tracer, func() {
		// Nothing to clean up in this case
	}
}

// WithTestSpan creates a new span for testing and returns the updated context
func WithTestSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	if ctx == nil {
		ctx = context.Background()
	}
	
	// Create a new span
	return otel.Tracer("test").Start(ctx, name)
}