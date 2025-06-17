# Telemetry Package

The `telemetry` package provides functionality for distributed tracing and metrics in applications. It is designed to be used with OpenTelemetry and Prometheus.

## Features

- Distributed tracing with OpenTelemetry
- Metrics collection with OpenTelemetry and Prometheus
- HTTP instrumentation for tracing requests and responses
- Unified interface for tracing and metrics

## Usage

### Initialization

```go
import (
    "context"
    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/telemetry"
    "github.com/knadh/koanf/v2"
)

func main() {
    ctx := context.Background()
    logger := logging.NewContextLogger(/* ... */)
    
    // Load configuration
    k := koanf.New(".")
    // ... load configuration into koanf
    
    // Create telemetry provider
    telemetryProvider, err := telemetry.NewTelemetryProvider(ctx, logger, k)
    if err != nil {
        logger.Error(ctx, "Failed to create telemetry provider", err)
        return
    }
    defer telemetryProvider.Shutdown(ctx)
    
    // Use telemetry provider
    // ...
}
```

### Tracing

```go
// Start a span
ctx, span := telemetry.StartSpan(ctx, "operation-name")
defer span.End()

// Add attributes to a span
telemetry.AddSpanAttributes(ctx, 
    attribute.String("key", "value"),
    attribute.Int("count", 42),
)

// Record an error
if err != nil {
    telemetry.RecordErrorSpan(ctx, err)
}

// Wrap a function with a span
err := telemetry.WithSpan(ctx, "operation-name", func(ctx context.Context) error {
    // Do something
    return nil
})

// Wrap a function with a span and measure execution time
duration, err := telemetry.WithSpanTimed(ctx, "operation-name", func(ctx context.Context) error {
    // Do something
    return nil
})
```

### HTTP Instrumentation

```go
// Instrument an HTTP handler
handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Handle request
})
instrumentedHandler := telemetry.InstrumentHandler(handler, "http-handler")

// Instrument an HTTP client
client := &http.Client{}
instrumentedClient := telemetry.InstrumentClient(client)

// Create an HTTP middleware for tracing
middleware := telemetry.NewHTTPMiddleware(logger)
```

### Metrics

```go
// Record HTTP request metrics
telemetry.RecordHTTPRequest(ctx, "GET", "/api/users", 200, 50*time.Millisecond, 1024)

// Record database operation metrics
telemetry.RecordDBOperation(ctx, "query", "postgres", "users", 10*time.Millisecond, nil)

// Record error metrics
telemetry.RecordErrorMetric(ctx, "validation", "create-user")

// Update database connection count
telemetry.UpdateDBConnections(ctx, "postgres", 1) // +1 connection
telemetry.UpdateDBConnections(ctx, "postgres", -1) // -1 connection
```

### Prometheus Integration

```go
// Create a Prometheus handler
prometheusHandler := telemetryProvider.CreatePrometheusHandler()

// Register the handler
http.Handle("/metrics", prometheusHandler)
```

## Configuration

The telemetry package can be configured using the following configuration structure:

```yaml
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