# Telemetry Module
The Telemetry Module provides functionality for distributed tracing and metrics in applications. It is designed to be used with OpenTelemetry and Prometheus.


## Overview

The Telemetry module provides comprehensive observability capabilities for Go applications. It serves as a central component in the ServiceLib library for collecting, processing, and exporting metrics and distributed traces. The module implements the OpenTelemetry standard, offering a vendor-neutral approach to observability that can integrate with various backends like Prometheus, Jaeger, Zipkin, and cloud-based observability platforms.

This module is designed to provide insights into application performance, behavior, and health, helping developers identify and diagnose issues in distributed systems. It includes utilities for tracing request flows across service boundaries, measuring operation durations, recording error rates, and monitoring resource usage.

## Features

- Distributed tracing with OpenTelemetry
- Metrics collection with OpenTelemetry and Prometheus
- HTTP instrumentation for tracing requests and responses
- Unified interface for tracing and metrics


## Installation

```bash
go get github.com/abitofhelp/servicelib/telemetry
```


## Quick Start

See the [Initialization example](../EXAMPLES/telemetry/initialization_example.go) for a complete, runnable example of how to use the Telemetry module.


## Configuration

The telemetry package can be configured using the following configuration structure:

```yaml
telemetry:
  service_name: "my-service"
  environment: "production"
  version: "1.0.0"

  tracing:
    enabled: true
    exporter: "jaeger"
    endpoint: "http://jaeger:14268/api/traces"
    sampling_ratio: 0.1

  metrics:
    enabled: true
    exporter: "prometheus"
    endpoint: ":9090"
    reporting_interval: "15s"

  http:
    instrumentation_enabled: true
    request_headers_to_record:
      - "user-agent"
      - "content-type"
      - "x-request-id"
```

## API Documentation

### Initialization

The `NewTelemetryProvider` function creates a new telemetry provider that can be used for tracing and metrics.

#### Creating a Telemetry Provider

See the [Initialization example](../EXAMPLES/telemetry/initialization_example.go) for a complete, runnable example of how to create and configure a telemetry provider.

### Tracing

The telemetry package provides functions for distributed tracing, including starting spans, adding attributes, and recording errors.

#### Using Tracing

See the [Tracing example](../EXAMPLES/telemetry/tracing_example.go) for a complete, runnable example of how to use the tracing functionality.

### HTTP Instrumentation

The telemetry package provides functions for instrumenting HTTP handlers and clients with tracing.

#### Instrumenting HTTP Components

See the [HTTP Instrumentation example](../EXAMPLES/telemetry/http_instrumentation_example.go) for a complete, runnable example of how to instrument HTTP handlers and clients.

### Metrics

The telemetry package provides functions for recording various types of metrics, including HTTP requests, database operations, and errors.

#### Recording Metrics

See the [Metrics example](../EXAMPLES/telemetry/metrics_example.go) for a complete, runnable example of how to record different types of metrics.

### Prometheus Integration

The telemetry package provides functions for integrating with Prometheus for metrics collection and visualization.

#### Integrating with Prometheus

See the [Prometheus Integration example](../EXAMPLES/telemetry/prometheus_integration_example.go) for a complete, runnable example of how to integrate with Prometheus.


### Core Types

The telemetry module provides several core types for distributed tracing and metrics:

#### TelemetryProvider

The `TelemetryProvider` struct is the main entry point for the telemetry module. It encapsulates both tracing and metrics functionality, providing a unified interface for all telemetry operations.

The TelemetryProvider manages the lifecycle of trace providers, metric providers, and exporters, ensuring proper initialization and shutdown.

See the [Initialization example](../EXAMPLES/telemetry/initialization_example.go) for a complete, runnable example of how to create and use a TelemetryProvider.

#### Tracer

The `Tracer` interface provides methods for creating and managing spans in distributed tracing. It wraps the OpenTelemetry Tracer interface, adding convenience methods for common operations.

The Tracer allows you to create spans, add attributes to spans, record errors, and set span status.

See the [Tracing example](../EXAMPLES/telemetry/tracing_example.go) for a complete, runnable example of how to use the Tracer.

#### MetricProvider

The `MetricProvider` interface provides methods for creating and recording metrics. It supports various metric types including counters, gauges, and histograms.

The MetricProvider allows you to create and record metrics with labels, set gauge values, and record histogram measurements.

See the [Metrics example](../EXAMPLES/telemetry/metrics_example.go) for a complete, runnable example of how to use the MetricProvider.

### Key Methods

The telemetry module provides several key methods for distributed tracing and metrics:

#### NewTelemetryProvider

The `NewTelemetryProvider` function creates a new telemetry provider configured with the provided options.

```
func NewTelemetryProvider(ctx context.Context, config Config) (*TelemetryProvider, error)
```

This function initializes both tracing and metrics providers based on the configuration, setting up exporters and processors as needed.

See the [Initialization example](../EXAMPLES/telemetry/initialization_example.go) for a complete, runnable example of how to use the NewTelemetryProvider function.

#### StartSpan

The `StartSpan` method creates a new span with the given name and options.

```
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
```

This method creates a new span, adds it to the trace in the provided context, and returns a new context containing the span along with the span itself.

See the [Tracing example](../EXAMPLES/telemetry/tracing_example.go) for a complete, runnable example of how to use the StartSpan method.

#### RecordMetric

The `RecordMetric` method records a measurement for the specified metric.

```
func (m *MetricProvider) RecordMetric(ctx context.Context, name string, value float64, labels ...attribute.KeyValue)
```

This method records a measurement for a counter or gauge metric with the specified name, value, and labels.

See the [Metrics example](../EXAMPLES/telemetry/metrics_example.go) for a complete, runnable example of how to use the RecordMetric method.

#### InstrumentHTTPHandler

The `InstrumentHTTPHandler` function wraps an HTTP handler with tracing middleware.

```
func InstrumentHTTPHandler(handler http.Handler, name string, options ...HTTPMiddlewareOption) http.Handler
```

This function creates a middleware that traces HTTP requests, recording span information such as request method, path, status code, and duration.

See the [HTTP Instrumentation example](../EXAMPLES/telemetry/http_instrumentation_example.go) for a complete, runnable example of how to use the InstrumentHTTPHandler function.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/telemetry/basic_usage_example.go) - Shows basic usage of the telemetry
- [Advanced Configuration](../EXAMPLES/telemetry/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/telemetry/error_handling_example.go) - Shows how to handle errors

## Best Practices

1. **Service Name**: Use a consistent service name across all your services
2. **Sampling Ratio**: In production, consider using a sampling ratio less than 1.0 to reduce overhead
3. **Span Naming**: Use descriptive names for spans that include the operation being performed
4. **Attributes**: Add relevant attributes to spans to provide context for debugging
5. **Error Handling**: Always record errors on spans to make debugging easier
6. **Metrics Naming**: Use consistent naming conventions for metrics
7. **Cardinality**: Be careful with high-cardinality labels in metrics

## Troubleshooting

### Common Issues

#### Traces Not Appearing in Backend

**Issue**: Spans are being created in the application, but they don't appear in the tracing backend (Jaeger, Zipkin, etc.).

**Solution**: 
1. Check that the exporter is properly configured with the correct endpoint
2. Verify that the sampling ratio is not set too low (a value of 0.0 means no traces will be exported)
3. Ensure that the backend service is running and accessible from your application
4. Check for any network issues or firewall rules that might be blocking the connection

#### Metrics Not Being Recorded

**Issue**: Metrics are being recorded in the application, but they don't appear in Prometheus or other metrics backends.

**Solution**:
1. Verify that the metrics exporter is properly configured
2. Check that the Prometheus scrape endpoint is accessible
3. Ensure that Prometheus is configured to scrape your application's metrics endpoint
4. Check the reporting interval to ensure it's not set too high

#### High Cardinality Metrics

**Issue**: The application is generating too many unique time series, causing performance issues in Prometheus.

**Solution**:
1. Reduce the number of labels used in metrics
2. Avoid using high-cardinality values like user IDs or request IDs as label values
3. Consider using histograms instead of counters or gauges for certain metrics
4. Implement server-side aggregation where possible

## Related Components

- [Logging](../logging/README.md) - The logging component integrates with telemetry to include trace IDs in log entries, providing correlation between logs and traces.
- [Middleware](../middleware/README.md) - The middleware component uses telemetry for tracing HTTP requests and recording metrics.
- [Health](../health/README.md) - The health component uses telemetry to expose health check metrics.
- [Errors](../errors/README.md) - The errors component integrates with telemetry to record error metrics and add error details to spans.
- [Config](../config/README.md) - The config component is used to configure telemetry providers and exporters.

## Contributing

Contributions to this telemetry are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
