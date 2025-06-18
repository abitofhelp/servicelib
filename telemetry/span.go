// Copyright (c) 2025 A Bit of Help, Inc.

// Package telemetry provides functionality for monitoring and tracing application behavior.
package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// GetSpanFromContext retrieves the current span from the context.
// This is a helper function for testing and debugging.
//
// Parameters:
//   - ctx: The context containing the span
//
// Returns:
//   - trace.Span: The span from the context, or nil if no span is present
func GetSpanFromContext(ctx context.Context) trace.Span {
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)

	// Check if the span is the default no-op span (which is returned when no span exists in the context)
	if !span.SpanContext().IsValid() {
		return nil
	}

	return span
}
