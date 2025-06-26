// Copyright (c) 2025 A Bit of Help, Inc.

// Package telemetry provides functionality for collecting and exporting telemetry data,
// including distributed tracing and metrics.
//
// This file contains the tracing-specific components of the telemetry package,
// including interfaces and implementations for creating and managing trace spans.
// It provides abstractions over the underlying OpenTelemetry tracing implementation,
// making it easier to use tracing throughout the application and to mock tracing
// in tests.
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

// noopSpan is a no-op implementation of Span.
// This implementation satisfies the Span interface but does not perform any actual
// tracing operations. It's used when tracing is disabled or in testing scenarios
// where actual tracing is not desired.
type noopSpan struct{}

// End implements the Span.End method but does nothing.
// This method is called when the operation represented by the span is finished.
func (s *noopSpan) End() {}

// SetAttributes implements the Span.SetAttributes method but does nothing.
// This method would normally set attributes on the span to provide additional
// context about the operation being traced.
//
// Parameters:
//   - attributes: Key-value pairs that would be added as attributes to the span.
func (s *noopSpan) SetAttributes(attributes ...attribute.KeyValue) {}

// RecordError implements the Span.RecordError method but does nothing.
// This method would normally mark the span as having encountered an error and
// add error details to the span.
//
// Parameters:
//   - err: The error that would be recorded.
//   - opts: Additional options for the error event.
func (s *noopSpan) RecordError(err error, opts ...trace.EventOption) {}

// noopTracer is a no-op implementation of Tracer.
// This implementation satisfies the Tracer interface but does not perform any actual
// tracing operations. It's used when tracing is disabled or in testing scenarios
// where actual tracing is not desired.
type noopTracer struct{}

// Start implements the Tracer.Start method but returns a no-op span.
// This method would normally create a new span with the given name and return
// a context containing the span and the span itself.
//
// Parameters:
//   - ctx: The parent context that would normally be used to create a child span.
//   - name: The name of the operation that would be traced.
//
// Returns:
//   - context.Context: The original context, unchanged.
//   - Span: A no-op implementation of the Span interface.
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

// otelSpan is an implementation of Span that wraps an OpenTelemetry span.
// This implementation delegates all operations to the underlying OpenTelemetry span,
// providing a bridge between our Span interface and the OpenTelemetry trace.Span type.
// This allows the rest of the application to use our simplified Span interface while
// still leveraging the full capabilities of OpenTelemetry.
type otelSpan struct {
	// span is the underlying OpenTelemetry span that this otelSpan wraps
	span trace.Span
}

// End completes the span by calling End on the underlying OpenTelemetry span.
// This should be called when the operation represented by the span is finished,
// typically in a defer statement after creating a span.
func (s *otelSpan) End() {
	s.span.End()
}

// SetAttributes sets attributes on the underlying OpenTelemetry span.
// Attributes provide additional context about the operation being traced,
// such as request parameters, response codes, or other relevant information.
//
// Parameters:
//   - attributes: Key-value pairs to add as attributes to the span.
func (s *otelSpan) SetAttributes(attributes ...attribute.KeyValue) {
	s.span.SetAttributes(attributes...)
}

// RecordError records an error on the underlying OpenTelemetry span.
// This marks the span as having encountered an error and adds error details
// to the span, making it easier to identify and debug failures in traces.
//
// Parameters:
//   - err: The error to record.
//   - opts: Additional options for the error event, such as timestamps or attributes.
func (s *otelSpan) RecordError(err error, opts ...trace.EventOption) {
	s.span.RecordError(err, opts...)
}

// otelTracer is an implementation of Tracer that uses OpenTelemetry.
// This implementation delegates all operations to the underlying OpenTelemetry tracer,
// providing a bridge between our Tracer interface and the OpenTelemetry trace.Tracer type.
// This allows the rest of the application to use our simplified Tracer interface while
// still leveraging the full capabilities of OpenTelemetry.
type otelTracer struct {
	// tracer is the underlying OpenTelemetry tracer that this otelTracer wraps
	tracer trace.Tracer
}

// Start creates a new span with the given name using the underlying OpenTelemetry tracer.
// It returns a new context containing the span and the span itself wrapped in our Span interface.
// The returned context should be passed to downstream operations to maintain the trace
// context across function calls.
//
// Parameters:
//   - ctx: The parent context, which may contain a parent span.
//   - name: The name of the operation being traced, which should be descriptive of the operation.
//
// Returns:
//   - context.Context: A new context containing the span.
//   - Span: The newly created span wrapped in our Span interface.
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
