# Telemetry Module

The Telemetry Module provides functionality for distributed tracing and metrics in applications. It is designed to be used with OpenTelemetry and Prometheus.

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

See the [Initialization example](../examples/telemetry/initialization_example.go) for a complete, runnable example of how to use the Telemetry module.

## API Documentation

### Initialization

The `NewTelemetryProvider` function creates a new telemetry provider that can be used for tracing and metrics.

#### Creating a Telemetry Provider

See the [Initialization example](../examples/telemetry/initialization_example.go) for a complete, runnable example of how to create and configure a telemetry provider.

### Tracing

The telemetry package provides functions for distributed tracing, including starting spans, adding attributes, and recording errors.

#### Using Tracing

See the [Tracing example](../examples/telemetry/tracing_example.go) for a complete, runnable example of how to use the tracing functionality.

### HTTP Instrumentation

The telemetry package provides functions for instrumenting HTTP handlers and clients with tracing.

#### Instrumenting HTTP Components

See the [HTTP Instrumentation example](../examples/telemetry/http_instrumentation_example.go) for a complete, runnable example of how to instrument HTTP handlers and clients.

### Metrics

The telemetry package provides functions for recording various types of metrics, including HTTP requests, database operations, and errors.

#### Recording Metrics

See the [Metrics example](../examples/telemetry/metrics_example.go) for a complete, runnable example of how to record different types of metrics.

### Prometheus Integration

The telemetry package provides functions for integrating with Prometheus for metrics collection and visualization.

#### Integrating with Prometheus

See the [Prometheus Integration example](../examples/telemetry/prometheus_integration_example.go) for a complete, runnable example of how to integrate with Prometheus.

## Configuration

The telemetry package can be configured using the following configuration structure:

```yaml
# Example configuration for the telemetry module
telemetry:
  enabled: true
  service_name: "my-service"
  environment: "development"
  version: "1.0.0"
  shutdown_timeout: 5

  otlp:
    endpoint: "localhost:4317"
    insecure: true
    timeout_seconds: 5

  tracing:
    enabled: true
    sampling_ratio: 1.0
    propagation_keys:
      - "traceparent"
      - "tracestate"
      - "baggage"

  metrics:
    enabled: true
    reporting_frequency_seconds: 15
    prometheus:
      enabled: true
      listen: "0.0.0.0:8089"
      path: "/metrics"

  http:
    tracing_enabled: true
```

## Best Practices

1. **Service Name**: Use a consistent service name across all your services
2. **Sampling Ratio**: In production, consider using a sampling ratio less than 1.0 to reduce overhead
3. **Span Naming**: Use descriptive names for spans that include the operation being performed
4. **Attributes**: Add relevant attributes to spans to provide context for debugging
5. **Error Handling**: Always record errors on spans to make debugging easier
6. **Metrics Naming**: Use consistent naming conventions for metrics
7. **Cardinality**: Be careful with high-cardinality labels in metrics
