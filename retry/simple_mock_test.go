// Copyright (c) 2025 A Bit of Help, Inc.

package retry

import (
	"context"

	"github.com/abitofhelp/servicelib/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// SimpleMockSpan is a simple mock implementation of telemetry.Span
type SimpleMockSpan struct{}

// End implements telemetry.Span
func (s *SimpleMockSpan) End() {}

// SetAttributes implements telemetry.Span
func (s *SimpleMockSpan) SetAttributes(attributes ...attribute.KeyValue) {}

// RecordError implements telemetry.Span
func (s *SimpleMockSpan) RecordError(err error, opts ...trace.EventOption) {}

// SimpleMockTracer is a simple mock implementation of telemetry.Tracer
type SimpleMockTracer struct{}

// Start implements telemetry.Tracer
func (t *SimpleMockTracer) Start(ctx context.Context, name string) (context.Context, telemetry.Span) {
	return ctx, &SimpleMockSpan{}
}
