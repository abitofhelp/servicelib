# Middleware Package

The `middleware` package provides HTTP middleware components for common cross-cutting concerns in Go applications. It includes middleware for authentication, logging, metrics collection, tracing, and more.

## Features

- **Authentication**: JWT authentication middleware
- **Logging**: Request/response logging
- **Metrics**: Request metrics collection
- **Tracing**: Distributed tracing
- **Recovery**: Panic recovery
- **CORS**: Cross-Origin Resource Sharing
- **Rate Limiting**: Request rate limiting

## Installation

```bash
go get github.com/abitofhelp/servicelib/middleware
```

## Usage

### Middleware Chain

```go
package main

import (
    "log"
    "net/http"
    
    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/middleware"
)

func main() {
    // Create a logger
    logger, _ := logging.NewLogger("info", false)
    defer logger.Sync()
    
    // Create a simple HTTP handler
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    // Create middleware chain
    handler := middleware.Chain(
        mux,
        middleware.RequestID(),
        middleware.Logging(logger),
        middleware.Metrics(),
        middleware.Tracing("my-service"),
        middleware.Recovery(logger),
    )
    
    // Start server with middleware
    log.Println("Starting server on :8080")
    http.ListenAndServe(":8080", handler)
}
```

### Individual Middleware Components

#### Request ID Middleware

```go
// Add request ID middleware
handler := middleware.RequestID()(yourHandler)

// Customize request ID header
handler := middleware.RequestIDWithConfig(middleware.RequestIDConfig{
    Header: "X-Custom-Request-ID",
    Generator: func() string {
        return uuid.New().String()
    },
})(yourHandler)
```

#### Logging Middleware

```go
// Create a logger
logger, _ := logging.NewLogger("info", false)

// Add logging middleware
handler := middleware.Logging(logger)(yourHandler)

// Customize logging middleware
handler := middleware.LoggingWithConfig(middleware.LoggingConfig{
    Logger: logger,
    SkipPaths: []string{"/health", "/metrics"},
    LogRequestHeaders: true,
    LogResponseHeaders: true,
    LogRequestBody: false,
    LogResponseBody: false,
})(yourHandler)
```

#### Metrics Middleware

```go
// Add metrics middleware
handler := middleware.Metrics()(yourHandler)

// Customize metrics middleware
handler := middleware.MetricsWithConfig(middleware.MetricsConfig{
    Subsystem: "api",
    SkipPaths: []string{"/health", "/metrics"},
    LabelNames: []string{"method", "path", "status"},
})(yourHandler)
```

#### Tracing Middleware

```go
// Add tracing middleware
handler := middleware.Tracing("my-service")(yourHandler)

// Customize tracing middleware
handler := middleware.TracingWithConfig(middleware.TracingConfig{
    ServiceName: "my-service",
    SkipPaths: []string{"/health", "/metrics"},
    TracePropagationHeaders: true,
})(yourHandler)
```

#### Recovery Middleware

```go
// Add recovery middleware
handler := middleware.Recovery(logger)(yourHandler)

// Customize recovery middleware
handler := middleware.RecoveryWithConfig(middleware.RecoveryConfig{
    Logger: logger,
    PrintStack: true,
    LogAllRequests: false,
    PanicHandler: func(w http.ResponseWriter, r *http.Request, err interface{}) {
        // Custom panic handling
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Something went wrong"))
    },
})(yourHandler)
```

#### CORS Middleware

```go
// Add CORS middleware
handler := middleware.CORS()(yourHandler)

// Customize CORS middleware
handler := middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"https://example.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Content-Type", "Authorization"},
    ExposeHeaders: []string{"X-Request-ID"},
    AllowCredentials: true,
    MaxAge: 86400,
})(yourHandler)
```

#### Rate Limiting Middleware

```go
// Add rate limiting middleware
handler := middleware.RateLimit(10, 1*time.Second)(yourHandler)

// Customize rate limiting middleware
handler := middleware.RateLimitWithConfig(middleware.RateLimitConfig{
    Rate: 10,
    Burst: 20,
    Period: 1 * time.Second,
    StoreType: "memory", // or "redis"
    RedisURL: "redis://localhost:6379",
    KeyFunc: func(r *http.Request) string {
        return r.RemoteAddr // Rate limit by IP
    },
})(yourHandler)
```

### Creating Custom Middleware

```go
package main

import (
    "net/http"
    "time"
)

// Create a custom middleware function
func TimingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Call the next handler
        next.ServeHTTP(w, r)
        
        // Calculate request duration
        duration := time.Since(start)
        
        // Add a custom header with the request duration
        w.Header().Set("X-Request-Duration", duration.String())
    })
}

func main() {
    // Create a simple HTTP handler
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    // Add your custom middleware
    handler := TimingMiddleware(mux)
    
    // Start server with middleware
    http.ListenAndServe(":8080", handler)
}
```

## Best Practices

1. **Order Matters**: The order of middleware in the chain is important. For example, recovery middleware should be first to catch panics in other middleware.

   ```go
   handler := middleware.Chain(
       mux,
       middleware.Recovery(logger),  // First to catch panics
       middleware.RequestID(),       // Second to add request ID for logging
       middleware.Logging(logger),   // Third to log with request ID
       middleware.Metrics(),         // Fourth to collect metrics
       middleware.Tracing("service") // Fifth to add tracing
   )
   ```

2. **Performance**: Be mindful of middleware performance, especially for high-traffic services.

3. **Error Handling**: Ensure middleware properly handles errors and doesn't swallow them.

4. **Context Values**: Use context to pass values between middleware and handlers.

5. **Middleware Configuration**: Use configuration options to customize middleware behavior.

6. **Testing**: Test middleware in isolation and as part of the chain.

7. **Logging**: Include relevant information in logs, but be careful not to log sensitive data.

## License

This project is licensed under the MIT License - see the LICENSE file for details.