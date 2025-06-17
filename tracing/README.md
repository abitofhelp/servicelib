# Tracing Package

The `tracing` package provides functionality for distributed tracing and metrics in applications. It is designed to be used with OpenTelemetry and Prometheus.

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
    "github.com/abitofhelp/servicelib/tracing"
)

func main() {
    ctx := context.Background()
    logger := logging.NewContextLogger(/* ... */)

    // Configure telemetry
    config := tracing.TelemetryConfig{
        Enabled:     true,
        ServiceName: "my-service",
        Environment: "development",
        Version:     "1.0.0",
        Tracing: struct {
            Enabled         bool
            SamplingRatio   float64
            PropagationKeys []string
            OTLPEndpoint    string
            OTLPInsecure    bool
            OTLPTimeout     int
            ShutdownTimeout int
        }{
            Enabled:         true,
            SamplingRatio:   1.0,
            PropagationKeys: []string{"traceparent", "tracestate", "baggage"},
            OTLPEndpoint:    "localhost:4317",
            OTLPInsecure:    true,
            OTLPTimeout:     5,
            ShutdownTimeout: 5,
        },
        Metrics: struct {
            Enabled    bool
            Prometheus struct {
                Enabled bool
                Path    string
            }
        }{
            Enabled: true,
            Prometheus: struct {
                Enabled bool
                Path    string
            }{
                Enabled: true,
                Path:    "/metrics",
            },
        },
    }

    // Create telemetry provider
    telemetry, err := tracing.NewTelemetryProvider(ctx, logger, config)
    if err != nil {
        logger.Error(ctx, "Failed to create telemetry provider", err)
        return
    }
    defer telemetry.Shutdown(ctx)

    // Use telemetry provider
    // ...
}
```

### Tracing

```go
// Start a span
ctx, span := tracing.StartSpan(ctx, "operation-name")
defer span.End()

// Add attributes to a span
tracing.AddSpanAttributes(ctx, 
    attribute.String("key", "value"),
    attribute.Int("count", 42),
)

// Record an error
if err != nil {
    tracing.RecordError(ctx, err)
}

// Wrap a function with a span
err := tracing.WithSpan(ctx, "operation-name", func(ctx context.Context) error {
    // Do something
    return nil
})

// Wrap a function with a span and measure execution time
duration, err := tracing.WithSpanTimed(ctx, "operation-name", func(ctx context.Context) error {
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
instrumentedHandler := tracing.InstrumentHandler(handler, "http-handler")

// Instrument an HTTP client
client := &http.Client{}
instrumentedClient := tracing.InstrumentClient(client)

// Create an HTTP middleware for tracing
middleware := tracing.NewHTTPMiddleware(logger)
```

### Metrics

```go
// Create a counter
counter, err := telemetry.Meter().Int64Counter(
    "counter-name",
    metric.WithDescription("Counter description"),
    metric.WithUnit("{count}"),
)

// Increment a counter
tracing.IncrementCounter(ctx, counter, 1,
    attribute.String("key", "value"),
)

// Create a histogram
histogram, err := telemetry.Meter().Float64Histogram(
    "histogram-name",
    metric.WithDescription("Histogram description"),
    metric.WithUnit("s"),
)

// Record a duration
tracing.RecordDuration(ctx, histogram, duration,
    attribute.String("key", "value"),
)

// Create an up/down counter
upDownCounter, err := telemetry.Meter().Int64UpDownCounter(
    "counter-name",
    metric.WithDescription("Counter description"),
    metric.WithUnit("{count}"),
)

// Update an up/down counter
tracing.UpdateUpDownCounter(ctx, upDownCounter, 1,
    attribute.String("key", "value"),
)
```

### Prometheus

```go
// Create a Prometheus handler
prometheusHandler := tracing.CreatePrometheusHandler()

// Register the handler
http.Handle("/metrics", prometheusHandler)
```

## Configuration with Docker Compose

The package is designed to work with Prometheus and Grafana for monitoring. Here's an example Docker Compose configuration:

```yaml
services:
  prometheus:
    image: "bitnami/prometheus:3.4.0"
    container_name: "prometheus"
    ports:
      - "9090:9090"
    volumes:
      - prometheus_data:/bitnami/prometheus
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - mynetwork
    restart: always
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"

  grafana:
    image: "bitnami/grafana:12.0.1"
    container_name: "grafana"
    depends_on:
      - prometheus
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/bitnami/grafana
      - ./grafana.ini:/etc/grafana/grafana.ini
    networks:
      - mynetwork
    restart: always
```

## Prometheus Configuration

Example prometheus.yml:

```yaml
global:
  scrape_interval:     15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: "my-service"
    metrics_path: "/metrics"
    scrape_interval: 15s
    static_configs:
      - targets: ["my-service:8080"]
```

## Grafana Configuration

Example grafana.ini:

```ini
[paths]
provisioning = /etc/grafana/provisioning

[server]
http_port = 3000

[datasources]
path = /etc/grafana/provisioning/datasources

[security]
admin_user = ${GF_SECURITY_ADMIN_USER}
admin_password = ${GF_SECURITY_ADMIN_PASSWORD}

[dashboards]
default_home_dashboard_path = /etc/grafana/dashboards/my_service_dashboard.json
```