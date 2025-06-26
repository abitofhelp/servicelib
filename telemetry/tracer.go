// Copyright (c) 2025 A Bit of Help, Inc.

package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Span represents a tracing span that can be used across packages.
// It provides methods for ending the span, setting attributes, and recording errors.
// This interface abstracts away the underlying tracing implementation,
// allowing for easier testing and flexibility in choosing tracing backends.
type Span interface {
	// End completes the span.
	// This should be called when the operation represented by the span is finished.
	// It is typically used in a defer statement after creating a span.
	End()

	// SetAttributes sets attributes on the span.
	// Attributes provide additional context about the operation being traced.
	//
	// Parameters:
	//   - attributes: Key-value pairs to add as attributes to the span.
	SetAttributes(attributes ...attribute.KeyValue)

	// RecordError records an error on the span.
	// This marks the span as having encountered an error and adds error details.
	//
	// Parameters:
	//   - err: The error to record.
	//   - opts: Additional options for the error event.
	RecordError(err error, opts ...trace.EventOption)
}

// Tracer is an interface for creating spans that can be used across packages.
// It provides a method for starting new spans to trace operations.
// This interface abstracts away the underlying tracing implementation,
// allowing for easier testing and flexibility in choosing tracing backends.
type Tracer interface {
	// Start creates a new span with the given name.
	// It returns a new context containing the span and the span itself.
	// The returned context should be passed to downstream operations to maintain
	// the trace context across function calls.
	//
	// Parameters:
	//   - ctx: The parent context.
	//   - name: The name of the operation being traced.
	//
	// Returns:
	//   - context.Context: A new context containing the span.
	//   - Span: The newly created span.
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

// NewNoopTracer creates a new no-op tracer.
// The no-op tracer implements the Tracer interface but does not perform any actual tracing.
// This is useful for testing, when tracing is disabled, or when you want to avoid the
// overhead of tracing in certain environments.
//
// Returns:
//   - Tracer: A no-op implementation of the Tracer interface.
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

// NewOtelTracer creates a new OpenTelemetry tracer.
// This function wraps an OpenTelemetry trace.Tracer in our Tracer interface,
// allowing it to be used with the rest of our telemetry system.
//
// Parameters:
//   - tracer: The OpenTelemetry trace.Tracer to wrap.
//
// Returns:
//   - Tracer: An implementation of the Tracer interface that uses OpenTelemetry.
func NewOtelTracer(tracer trace.Tracer) Tracer {
	return &otelTracer{tracer: tracer}
}

// GetTracer returns a tracer that can be used for tracing operations.
// This function provides a convenient way to get a tracer, either using a provided
// OpenTelemetry tracer or creating a default one if none is provided.
//
// Parameters:
//   - tracer: An optional OpenTelemetry trace.Tracer. If nil, a default tracer will be created.
//
// Returns:
//   - Tracer: An implementation of the Tracer interface that can be used for tracing operations.
//     If tracer is nil, a new tracer with the name "github.com/abitofhelp/servicelib/telemetry" will be created.
func GetTracer(tracer trace.Tracer) Tracer {
	if tracer == nil {
		return NewOtelTracer(otel.Tracer("github.com/abitofhelp/servicelib/telemetry"))
	}
	return NewOtelTracer(tracer)
}
