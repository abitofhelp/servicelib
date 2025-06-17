// Copyright (c) 2025 A Bit of Help, Inc.

// Package tracing provides functionality for distributed tracing in applications.
package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/abitofhelp/servicelib/logging"
)

// TracingProvider holds the tracing provider and tracer
type TracingProvider struct {
	// provider is the OpenTelemetry tracer provider
	provider *sdktrace.TracerProvider

	// tracer is the OpenTelemetry tracer
	tracer trace.Tracer

	// logger is the logger for the tracing provider
	logger *logging.ContextLogger
}

// Config holds configuration for tracing
type Config struct {
	// Enabled indicates whether tracing is enabled
	Enabled bool

	// ServiceName is the name of the service
	ServiceName string

	// Environment is the environment the service is running in
	Environment string

	// Version is the version of the service
	Version string

	// SamplingRatio is the ratio of traces to sample
	SamplingRatio float64

	// PropagationKeys are the keys to propagate in trace context
	PropagationKeys []string

	// OTLPEndpoint is the OTLP endpoint
	OTLPEndpoint string

	// OTLPInsecure indicates whether to use insecure connections
	OTLPInsecure bool

	// OTLPTimeout is the timeout for OTLP operations in seconds
	OTLPTimeout int

	// ShutdownTimeout is the timeout for shutting down tracing in seconds
	ShutdownTimeout int
}

// NewTracingProvider creates a new tracing provider
//
// Parameters:
//   - ctx: The context for the operation
//   - logger: The logger to use for logging tracing events
//   - config: The tracing configuration
//
// Returns:
//   - *TracingProvider: The tracing provider
//   - error: An error if the tracing provider creation fails
func NewTracingProvider(ctx context.Context, logger *logging.ContextLogger, config Config) (*TracingProvider, error) {
	if !config.Enabled {
		logger.Info(ctx, "Tracing is disabled")
		return nil, nil
	}

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.Version),
			attribute.String("environment", config.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Configure OTLP exporter
	secureOption := otlptracegrpc.WithTLSCredentials(insecure.NewCredentials())
	if !config.OTLPInsecure {
		secureOption = otlptracegrpc.WithTLSCredentials(nil)
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(config.OTLPEndpoint),
		secureOption,
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create tracer provider
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(config.SamplingRatio)),
	)

	// Set global tracer provider
	otel.SetTracerProvider(provider)

	// Set global propagator
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	// Create tracer
	tracer := provider.Tracer("github.com/abitofhelp/servicelib/tracing")

	logger.Info(ctx, "Tracing provider initialized",
		zap.String("service", config.ServiceName),
		zap.String("environment", config.Environment),
		zap.String("otlp_endpoint", config.OTLPEndpoint),
		zap.Float64("sampling_ratio", config.SamplingRatio),
	)

	return &TracingProvider{
		provider: provider,
		tracer:   tracer,
		logger:   logger,
	}, nil
}

// Shutdown shuts down the tracing provider
//
// Parameters:
//   - ctx: The context for the operation
//
// Returns:
//   - error: An error if the shutdown fails
func (tp *TracingProvider) Shutdown(ctx context.Context) error {
	if tp.provider == nil {
		return nil
	}

	if err := tp.provider.Shutdown(ctx); err != nil {
		tp.logger.Error(ctx, "Failed to shutdown tracing provider", zap.Error(err))
		return err
	}

	tp.logger.Info(ctx, "Tracing provider shut down")
	return nil
}

// Tracer returns the tracer
//
// Returns:
//   - trace.Tracer: The OpenTelemetry tracer
func (tp *TracingProvider) Tracer() trace.Tracer {
	return tp.tracer
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
	tracer := otel.Tracer("github.com/abitofhelp/servicelib/tracing")
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

// RecordError records an error on the current span in context.
// This is useful for recording errors that occur during a traced operation.
//
// Parameters:
//   - ctx: The context containing the span
//   - err: The error to record
//   - opts: Additional options for the error event
func RecordError(ctx context.Context, err error, opts ...trace.EventOption) {
	if err != nil {
		span := trace.SpanFromContext(ctx)
		span.RecordError(err, opts...)
	}
}

// WithSpan wraps a function with a span.
// This is useful for tracing functions.
//
// Parameters:
//   - ctx: The context to start the span from
//   - name: The name of the span
//   - fn: The function to wrap
//
// Returns:
//   - error: The error returned by the function
func WithSpan(ctx context.Context, name string, fn func(context.Context) error) error {
	ctx, span := StartSpan(ctx, name)
	defer span.End()

	return fn(ctx)
}

// WithSpanTimed wraps a function with a span and records the execution time.
// This is useful for tracing functions and measuring their execution time.
//
// Parameters:
//   - ctx: The context to start the span from
//   - name: The name of the span
//   - fn: The function to wrap
//
// Returns:
//   - time.Duration: The execution time of the function
//   - error: The error returned by the function
func WithSpanTimed(ctx context.Context, name string, fn func(context.Context) error) (time.Duration, error) {
	ctx, span := StartSpan(ctx, name)
	defer span.End()

	start := time.Now()
	err := fn(ctx)
	duration := time.Since(start)

	span.SetAttributes(attribute.Float64("duration_ms", float64(duration.Milliseconds())))
	if err != nil {
		RecordError(ctx, err)
	}

	return duration, err
}