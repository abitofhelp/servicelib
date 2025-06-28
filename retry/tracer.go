// Copyright (c) 2025 A Bit of Help, Inc.

package retry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Span represents a tracing span
type Span interface {
	// End completes the span
	End()
	// SetAttributes sets attributes on the span
	SetAttributes(attributes ...attribute.KeyValue)
	// RecordError records an error on the span
	RecordError(err error)
}

// Tracer is an interface for creating spans
type Tracer interface {
	// Start creates a new span
	Start(ctx context.Context, name string) (context.Context, Span)
}

// noopSpan is a no-op implementation of Span
type noopSpan struct{}

func (s *noopSpan) End() {}

func (s *noopSpan) SetAttributes(attributes ...attribute.KeyValue) {}

func (s *noopSpan) RecordError(err error) {}

// noopTracer is a no-op implementation of Tracer
type noopTracer struct{}

func (t *noopTracer) Start(ctx context.Context, name string) (context.Context, Span) {
	return ctx, &noopSpan{}
}

// NewNoopTracer creates a new no-op tracer
func NewNoopTracer() Tracer {
	return &noopTracer{}
}

// otelTracer is an implementation of Tracer that uses OpenTelemetry
type otelTracer struct {
	tracer trace.Tracer
}

// otelSpan is an implementation of Span that wraps an OpenTelemetry span
type otelSpan struct {
	span trace.Span
}

func (s *otelSpan) End() {
	s.span.End()
}

func (s *otelSpan) SetAttributes(attributes ...attribute.KeyValue) {
	s.span.SetAttributes(attributes...)
}

func (s *otelSpan) RecordError(err error) {
	s.span.RecordError(err)
}

func (t *otelTracer) Start(ctx context.Context, name string) (context.Context, Span) {
	ctx, span := t.tracer.Start(ctx, name)
	return ctx, &otelSpan{span: span}
}

// NewOtelTracer creates a new OpenTelemetry tracer
func NewOtelTracer(tracer trace.Tracer) Tracer {
	return &otelTracer{tracer: tracer}
}