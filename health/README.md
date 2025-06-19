# Health Check Package

The `health` package provides components for implementing health check endpoints in Go applications. It helps you create health check handlers for Kubernetes readiness and liveness probes, as well as custom health checks for your services.

## Features

- **Check Types**:
  - Liveness checks (basic service health)
  - Readiness checks (service ability to handle requests)
  - Startup checks (service initialization)

- **Features**:
  - Configurable check intervals
  - Automatic registration with HTTP server
  - Detailed health status reporting
  - Integration with Kubernetes probes

## Installation

```bash
go get github.com/abitofhelp/servicelib/health
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/abitofhelp/servicelib/health"
    _ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
    // Create a health handler with custom configuration
    config := health.Config{
        Name:           "my-service",
        Version:        "1.0.0",
        CheckTimeout:   5 * time.Second,
        CheckInterval:  30 * time.Second,
        ShutdownDelay:  10 * time.Second,
        ReadinessPath:  "/ready",
        LivenessPath:   "/live",
        StartupPath:    "/startup",
    }

    handler := health.NewHandler(config)

    // Add liveness checks (basic service health)
    handler.AddLivenessCheck("goroutines", health.GoroutineCountCheck(1000))
    handler.AddLivenessCheck("memory", health.MemoryUsageCheck(85.0))

    // Add readiness checks (service ability to handle requests)
    // Assuming you have a database connection
    db, err := sql.Open("pgx", "postgres://user:password@localhost:5432/mydb?sslmode=disable")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    handler.AddReadinessCheck("database", checkDatabaseConnection(db))
    handler.AddReadinessCheck("api", checkExternalAPI("https://api.example.com/health"))

    // Add startup checks (service initialization)
    handler.AddStartupCheck("migrations", checkDatabaseMigrations(db))

    // Create an HTTP server
    server := &http.Server{
        Addr:    ":8080",
        Handler: http.DefaultServeMux,
    }

    // Register health check handlers
    handler.RegisterHandlers(http.DefaultServeMux)

    // Start the server
    log.Println("Starting server on :8080")
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Failed to start server: %v", err)
    }
}

// Check function for database connection
func checkDatabaseConnection(db *sql.DB) health.CheckFunc {
    return func() error {
        ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
        defer cancel()

        if err := db.PingContext(ctx); err != nil {
            return fmt.Errorf("database ping failed: %w", err)
        }
        return nil
    }
}

// Check function for external API
func checkExternalAPI(url string) health.CheckFunc {
    return func() error {
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()

        req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        if err != nil {
            return fmt.Errorf("failed to create request: %w", err)
        }

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            return fmt.Errorf("API request failed: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
        }

        return nil
    }
}

// Check function for database migrations
func checkDatabaseMigrations(db *sql.DB) health.CheckFunc {
    return func() error {
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()

        var count int
        err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations").Scan(&count)
        if err != nil {
            return fmt.Errorf("failed to check migrations: %w", err)
        }

        if count == 0 {
            return fmt.Errorf("no migrations have been applied")
        }

        return nil
    }
}
```

### Built-in Checks

The health package provides several built-in checks:

```go
// Check if the number of goroutines is below a threshold
handler.AddLivenessCheck("goroutines", health.GoroutineCountCheck(1000))

// Check if memory usage is below a percentage threshold
handler.AddLivenessCheck("memory", health.MemoryUsageCheck(85.0))

// Check if disk usage is below a percentage threshold
handler.AddLivenessCheck("disk", health.DiskUsageCheck("/", 90.0))

// Check if CPU usage is below a percentage threshold
handler.AddLivenessCheck("cpu", health.CPUUsageCheck(90.0))

// Check if a TCP connection can be established
handler.AddReadinessCheck("tcp", health.TCPDialCheck("example.com:80", 1*time.Second))

// Check if an HTTP endpoint returns a 200 status
handler.AddReadinessCheck("http", health.HTTPGetCheck("https://example.com/health", 2*time.Second))
```

### Custom Health Status

You can create a custom health status response:

```go
package main

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/abitofhelp/servicelib/health"
)

// Custom health status
type CustomHealthStatus struct {
    Status    string            `json:"status"`
    Timestamp string            `json:"timestamp"`
    Version   string            `json:"version"`
    Checks    map[string]string `json:"checks"`
    Uptime    string            `json:"uptime"`
    Memory    string            `json:"memory"`
}

func main() {
    // Create a health handler
    handler := health.NewHandler(health.Config{
        Name:    "my-service",
        Version: "1.0.0",
    })

    // Add some checks
    handler.AddLivenessCheck("memory", health.MemoryUsageCheck(85.0))
    handler.AddReadinessCheck("api", health.HTTPGetCheck("https://api.example.com/health", 2*time.Second))

    // Create a custom handler that uses the health checks
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        // Run the checks
        status := handler.RunChecks()

        // Create a custom response
        customStatus := CustomHealthStatus{
            Status:    status.Status,
            Timestamp: time.Now().Format(time.RFC3339),
            Version:   status.Version,
            Checks:    status.Checks,
            Uptime:    time.Since(startTime).String(),
            Memory:    fmt.Sprintf("%.2f MB", float64(memStats.Alloc)/1024/1024),
        }

        // Write the response
        w.Header().Set("Content-Type", "application/json")
        if status.Status != "healthy" {
            w.WriteHeader(http.StatusServiceUnavailable)
        }
        json.NewEncoder(w).Encode(customStatus)
    })

    // Start the server
    http.ListenAndServe(":8080", nil)
}

var startTime = time.Now()
var memStats runtime.MemStats
```

### Integration with Kubernetes

The health package is designed to work well with Kubernetes probes:

```yaml
# Kubernetes deployment example
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-service
  template:
    metadata:
      labels:
        app: my-service
    spec:
      containers:
      - name: my-service
        image: my-service:1.0.0
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 2
        startupProbe:
          httpGet:
            path: /startup
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 30
```

## Health Check Types

### Liveness Checks

Liveness checks determine if your application is running properly. If a liveness check fails, Kubernetes will restart your container. Use liveness checks for conditions that require a restart to resolve.

Examples:
- Deadlocks
- Memory leaks
- Too many goroutines

### Readiness Checks

Readiness checks determine if your application can handle requests. If a readiness check fails, Kubernetes will stop sending traffic to your container until it passes again. Use readiness checks for temporary conditions that don't require a restart.

Examples:
- Database connections
- External API dependencies
- Cache warming

### Startup Checks

Startup checks determine if your application has successfully initialized. Kubernetes will not start liveness and readiness checks until the startup check passes. Use startup checks for initialization tasks.

Examples:
- Database migrations
- Configuration loading
- Initial data caching

## Best Practices

1. **Keep Checks Lightweight**: Health checks should be fast and not consume significant resources.

2. **Appropriate Timeouts**: Set appropriate timeouts for health checks to prevent hanging.

3. **Separate Concerns**: Use different check types for different concerns:
   - Liveness: Internal application health
   - Readiness: External dependency health
   - Startup: Initialization tasks

4. **Avoid False Positives**: Don't make checks too sensitive to prevent unnecessary restarts.

5. **Include Version Information**: Include version information in health status to help with debugging.

6. **Detailed Status**: Provide detailed status information for each check to aid in troubleshooting.

7. **Monitoring Integration**: Integrate health checks with your monitoring system.

## License

This project is licensed under the MIT License - see the LICENSE file for details.