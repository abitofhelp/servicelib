// Copyright (c) 2025 A Bit of Help, Inc.

package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Span represents a tracing span that can be used across packages
type Span interface {
	// End completes the span
	End()
	// SetAttributes sets attributes on the span
	SetAttributes(attributes ...attribute.KeyValue)
	// RecordError records an error on the span
	RecordError(err error, opts ...trace.EventOption)
}

// Tracer is an interface for creating spans that can be used across packages
type Tracer interface {
	// Start creates a new span
	Start(ctx context.Context, name string) (context.Context, Span)
}

// noopSpan is a no-op implementation of Span
type noopSpan struct{}

func (s *noopSpan) End() {}

func (s *noopSpan) SetAttributes(attributes ...attribute.KeyValue) {}

func (s *noopSpan) RecordError(err error, opts ...trace.EventOption) {}

// noopTracer is a no-op implementation of Tracer
type noopTracer struct{}

func (t *noopTracer) Start(ctx context.Context, name string) (context.Context, Span) {
	return ctx, &noopSpan{}
}

// NewNoopTracer creates a new no-op tracer
func NewNoopTracer() Tracer {
	return &noopTracer{}
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

func (s *otelSpan) RecordError(err error, opts ...trace.EventOption) {
	s.span.RecordError(err, opts...)
}

// otelTracer is an implementation of Tracer that uses OpenTelemetry
type otelTracer struct {
	tracer trace.Tracer
}

func (t *otelTracer) Start(ctx context.Context, name string) (context.Context, Span) {
	ctx, span := t.tracer.Start(ctx, name)
	return ctx, &otelSpan{span: span}
}

// NewOtelTracer creates a new OpenTelemetry tracer
func NewOtelTracer(tracer trace.Tracer) Tracer {
	return &otelTracer{tracer: tracer}
}

// GetTracer returns a tracer that can be used for tracing operations
// If a tracer is provided, it will be used; otherwise, a new tracer will be created
func GetTracer(tracer trace.Tracer) Tracer {
	if tracer == nil {
		return NewOtelTracer(otel.Tracer("github.com/abitofhelp/servicelib/telemetry"))
	}
	return NewOtelTracer(tracer)
}