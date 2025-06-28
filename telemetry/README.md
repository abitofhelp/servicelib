# Telemetry

## Overview

The Telemetry component provides functionality for collecting and exporting telemetry data, including distributed tracing and metrics. It integrates with OpenTelemetry and Prometheus to provide a unified interface for telemetry collection, supporting tracing requests across service boundaries, collecting metrics about application performance, and exporting this data to various backends for analysis and visualization.

## Features

- **Distributed Tracing**: Implement distributed tracing with OpenTelemetry to track requests across service boundaries
- **Metrics Collection**: Collect application metrics with OpenTelemetry and Prometheus
- **HTTP Middleware**: Automatically trace HTTP requests and collect metrics
- **Database Monitoring**: Track database operations with timing and error metrics
- **Error Tracking**: Record and monitor errors with detailed context
- **Configurable Exporters**: Export telemetry data to various backends (Jaeger, Zipkin, Prometheus, etc.)
- **Context Propagation**: Maintain trace context across service boundaries
- **Low Overhead**: Designed for minimal performance impact in production environments

## Installation

```bash
go get github.com/abitofhelp/servicelib/telemetry
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "net/http"
    
    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/telemetry"
    "github.com/knadh/koanf/v2"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    logger, _ := zap.NewProduction()
    contextLogger := logging.NewContextLogger(logger)
    
    // Create a configuration
    config := koanf.New(".")
    
    // Initialize telemetry provider
    ctx := context.Background()
    provider, err := telemetry.NewTelemetryProvider(ctx, contextLogger, config)
    if err != nil {
        log.Fatalf("Failed to initialize telemetry: %v", err)
    }
    defer provider.Shutdown(ctx)
    
    // Create an HTTP handler with telemetry
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    // Apply telemetry middleware
    instrumentedHandler := provider.InstrumentHandler(handler, "hello-world")
    
    // Start the server
    http.ListenAndServe(":8080", instrumentedHandler)
}
```

## Configuration

The telemetry component can be configured using the koanf configuration library. Here's an example configuration:

```yaml
telemetry:
  service_name: "my-service"
  enabled: true
  tracing:
    enabled: true
    exporter: "jaeger"
    jaeger:
      endpoint: "http://jaeger:14268/api/traces"
    sampling:
      ratio: 0.1
  metrics:
    enabled: true
    exporter: "prometheus"
    prometheus:
      port: 9090
```

## API Documentation

### Core Types

The telemetry package provides several core types for collecting and exporting telemetry data.

#### TelemetryProvider

The TelemetryProvider is the main entry point for the telemetry package. It provides access to tracers and meters, and methods for instrumenting HTTP handlers.

```go
type TelemetryProvider struct {
    // Contains unexported fields
}

// NewTelemetryProvider creates a new telemetry provider
func NewTelemetryProvider(ctx context.Context, logger *logging.ContextLogger, k *koanf.Koanf) (*TelemetryProvider, error)

// Shutdown shuts down the telemetry provider
func (tp *TelemetryProvider) Shutdown(ctx context.Context) error

// Meter returns the meter for collecting metrics
func (tp *TelemetryProvider) Meter() metric.Meter

// Tracer returns the tracer for creating spans
func (tp *TelemetryProvider) Tracer() trace.Tracer

// CreatePrometheusHandler creates an HTTP handler for Prometheus metrics
func (tp *TelemetryProvider) CreatePrometheusHandler() http.Handler

// InstrumentHandler instruments an HTTP handler with tracing and metrics
func (tp *TelemetryProvider) InstrumentHandler(handler http.Handler, operation string) http.Handler

// NewHTTPMiddleware creates middleware for HTTP handlers
func (tp *TelemetryProvider) NewHTTPMiddleware() func(http.Handler) http.Handler
```

#### Tracer

The Tracer interface provides methods for creating and managing trace spans.

```go
type Tracer interface {
    // Start creates a new span with the given name
    Start(ctx context.Context, name string) (context.Context, Span)
}

// NewNoopTracer creates a new no-op tracer
func NewNoopTracer() Tracer

// NewOtelTracer creates a new OpenTelemetry tracer
func NewOtelTracer(tracer trace.Tracer) Tracer

// GetTracer returns a tracer that can be used for tracing operations
func GetTracer(tracer trace.Tracer) Tracer
```

#### Span

The Span interface represents a single trace span.

```go
type Span interface {
    // End completes the span
    End()
    
    // SetAttributes sets attributes on the span
    SetAttributes(attributes ...attribute.KeyValue)
    
    // RecordError records an error on the span
    RecordError(err error, opts ...trace.EventOption)
}
```

### Key Methods

The telemetry package provides several utility methods for working with traces and spans.

#### WithSpan

The WithSpan function executes a function within a new span.

```go
func WithSpan(ctx context.Context, name string, fn func(context.Context) error) error
```

#### WithSpanTimed

The WithSpanTimed function executes a function within a new span and returns the duration.

```go
func WithSpanTimed(ctx context.Context, name string, fn func(context.Context) error) (time.Duration, error)
```

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/telemetry/basic_usage_example.go) - Shows basic usage of the telemetry package
- [HTTP Tracing](../EXAMPLES/telemetry/http_tracing_example.go) - Shows how to trace HTTP requests
- [Metrics Collection](../EXAMPLES/telemetry/metrics_example.go) - Shows how to collect metrics
- [Database Tracing](../EXAMPLES/telemetry/database_example.go) - Shows how to trace database operations
- [Error Tracking](../EXAMPLES/telemetry/error_tracking_example.go) - Shows how to track errors

## Best Practices

1. **Initialize Early**: Initialize the telemetry provider at the start of your application
2. **Use Context Propagation**: Always pass context through your application to maintain trace context
3. **Name Spans Clearly**: Use descriptive names for spans to make traces easier to understand
4. **Add Relevant Attributes**: Include attributes that provide context about the operation
5. **Record Errors**: Always record errors on spans to make debugging easier
6. **Set Appropriate Sampling**: Configure sampling based on your traffic volume and needs
7. **Monitor Resource Usage**: Keep an eye on the resource usage of your telemetry system

## Troubleshooting

### Common Issues

#### High Cardinality

If you're experiencing high cardinality in your metrics (too many unique time series):

- Reduce the number of attributes you're adding to metrics
- Use fewer unique values for attributes
- Consider using a lower sampling rate

#### Missing Spans

If spans are missing from your traces:

- Ensure context is being properly propagated through your application
- Check that your sampling configuration isn't filtering out too many spans
- Verify that your exporter is properly configured and connected to your backend

## Related Components

- [Logging](../logging/README.md) - The telemetry package integrates with the logging package for error reporting
- [Errors](../errors/README.md) - The telemetry package uses the errors package for error handling
- [Context](../context/README.md) - The telemetry package relies on context for propagating trace context

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.