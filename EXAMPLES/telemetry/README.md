# Telemetry Package Examples

This directory contains examples demonstrating how to use the `telemetry` package, which provides utilities for application monitoring, metrics collection, and distributed tracing in Go applications. The package integrates with OpenTelemetry and supports various backends for collecting and visualizing telemetry data.

## Examples

### 1. Initialization Example

[initialization_example.go](initialization_example.go)

Demonstrates how to initialize the telemetry package.

Key concepts:
- Creating a telemetry provider with configuration
- Setting up configuration for tracing, metrics, and Prometheus
- Initializing the telemetry provider
- Handling graceful shutdown of telemetry
- Using context for telemetry operations

### 2. HTTP Instrumentation Example

[http_instrumentation_example.go](http_instrumentation_example.go)

Shows how to instrument HTTP servers and clients.

Key concepts:
- Adding tracing to HTTP handlers
- Instrumenting HTTP client requests
- Propagating trace context across HTTP calls
- Collecting HTTP metrics
- Analyzing HTTP request performance

### 3. Metrics Example

[metrics_example.go](metrics_example.go)

Demonstrates how to collect and report metrics.

Key concepts:
- Creating and using counters, gauges, and histograms
- Recording metric values
- Adding dimensions/labels to metrics
- Setting up metric views
- Configuring metric exporters

### 4. Prometheus Integration Example

[prometheus_integration_example.go](prometheus_integration_example.go)

Shows how to integrate with Prometheus for metrics collection.

Key concepts:
- Setting up a Prometheus metrics endpoint
- Configuring Prometheus scraping
- Creating custom Prometheus metrics
- Exposing application metrics to Prometheus
- Visualizing metrics in Prometheus

### 5. Tracing Example

[tracing_example.go](tracing_example.go)

Demonstrates how to implement distributed tracing.

Key concepts:
- Creating and managing spans
- Adding attributes to spans
- Propagating trace context
- Recording span events
- Handling trace sampling

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/telemetry/initialization_example.go
```

## Additional Resources

For more information about the telemetry package, see the [telemetry package documentation](../../telemetry/README.md).