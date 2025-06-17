# ServiceLib

## Overview

ServiceLib is a comprehensive Go library designed to accelerate the development of robust, production-ready microservices. It provides a collection of reusable components and utilities that address common challenges in service development, allowing developers to focus on business logic rather than infrastructure concerns.

The library follows modern Go practices and design patterns, with a focus on:

- **Modularity**: Each component can be used independently or together with others
- **Testability**: All components are designed with testing in mind
- **Performance**: Optimized for high-throughput microservices
- **Reliability**: Built-in error handling and recovery mechanisms
- **Observability**: Integrated logging, metrics, and tracing

## Features

- **Authentication** - JWT, OAuth2, and OIDC implementations for secure service-to-service and user authentication
- **Configuration** - Flexible configuration management with adapters for various sources (files, environment variables, etc.)
- **Context** - Context utilities for request handling, cancellation, and value propagation
- **Database** - Database connection and transaction management with support for PostgreSQL, SQLite, and MongoDB
- **Dependency Injection** - Container-based DI system for managing service dependencies
- **Error Handling** - Structured error types and handling with rich context information
- **GraphQL** - Utilities for building GraphQL services with gqlgen integration
- **Health Checks** - Health check endpoints and handlers for Kubernetes readiness and liveness probes
- **Logging** - Structured logging with Zap for high-performance logging
- **Middleware** - HTTP middleware components for common cross-cutting concerns
- **Repository Pattern** - Generic repository implementations for data access abstraction
- **Shutdown** - Graceful shutdown utilities for clean service termination
- **Signal Handling** - OS signal handling for responding to system events
- **Telemetry** - Metrics, tracing, and monitoring with Prometheus and OpenTelemetry
- **Validation** - Request and data validation using go-playground/validator

## Installation

```bash
go get github.com/abitofhelp/servicelib
```

## Usage Examples

### Configuration

```go
package main

import (
    "fmt"
    "log"
    "github.com/abitofhelp/servicelib/config"
)

// Example configuration struct
type AppConfig struct {
    Server struct {
        Port    int    `yaml:"port"`
        Host    string `yaml:"host"`
        Timeout int    `yaml:"timeout"`
    } `yaml:"server"`
    Database struct {
        URL      string `yaml:"url"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Pool     int    `yaml:"pool"`
    } `yaml:"database"`
    API struct {
        Key      string `yaml:"key"`
        Endpoint string `yaml:"endpoint"`
        Version  string `yaml:"version"`
    } `yaml:"api"`
    Logging struct {
        Level  string `yaml:"level"`
        Format string `yaml:"format"`
        Path   string `yaml:"path"`
    } `yaml:"logging"`
}

func main() {
    // Create a new configuration
    cfg, err := config.New("config.yaml", "env.yaml")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Get a string value
    apiKey := cfg.GetString("api.key")
    fmt.Println("API Key:", apiKey)

    // Get an int value with default
    port := cfg.GetInt("server.port", 8080)
    fmt.Println("Server Port:", port)

    // Get a nested value
    dbURL := cfg.GetString("database.url")
    fmt.Println("Database URL:", dbURL)

    // Get a boolean value
    debug := cfg.GetBool("logging.debug", false)
    fmt.Println("Debug Mode:", debug)

    // Get a duration value
    timeout := cfg.GetDuration("server.timeout", "30s")
    fmt.Println("Server Timeout:", timeout)

    // Bind configuration to a struct
    var appConfig AppConfig
    if err := cfg.Unmarshal(&appConfig); err != nil {
        log.Fatalf("Failed to unmarshal configuration: %v", err)
    }

    fmt.Printf("Server Configuration: %+v\n", appConfig.Server)
    fmt.Printf("Database Configuration: %+v\n", appConfig.Database)

    // Watch for configuration changes
    cfg.Watch(func() {
        // Reload configuration when changes are detected
        if err := cfg.Unmarshal(&appConfig); err != nil {
            log.Printf("Failed to reload configuration: %v", err)
            return
        }
        log.Println("Configuration reloaded successfully")
    })
}
```

### Health Checks

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/abitofhelp/servicelib/health"
    "github.com/abitofhelp/servicelib/db"
    _ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
    // Create a database connection
    database, err := sql.Open("pgx", "postgres://user:password@localhost:5432/mydb?sslmode=disable")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.Close()

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
    handler.AddReadinessCheck("database", checkDatabaseConnection(database))
    handler.AddReadinessCheck("api", checkExternalAPI("https://api.example.com/health"))

    // Add startup checks (service initialization)
    handler.AddStartupCheck("migrations", checkDatabaseMigrations(database))

    // Create an HTTP server
    server := &http.Server{
        Addr:    ":8080",
        Handler: http.DefaultServeMux,
    }

    // Register health check handlers
    handler.RegisterHandlers(http.DefaultServeMux)

    // Start the server in a goroutine
    go func() {
        log.Println("Starting server on :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop

    // Graceful shutdown
    log.Println("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown failed: %v", err)
    }

    log.Println("Server stopped gracefully")
}

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

func checkDatabaseMigrations(db *sql.DB) health.CheckFunc {
    return func() error {
        // Check if migrations have been applied
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

### Dependency Injection

```go
package main

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "log"

    "github.com/abitofhelp/servicelib/di"
    "github.com/abitofhelp/servicelib/config"
    "github.com/google/uuid"
    _ "github.com/jackc/pgx/v5/stdlib"
)

// Domain models
type User struct {
    ID       string
    Username string
    Email    string
    Active   bool
}

// Repository interfaces
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    GetByUsername(ctx context.Context, username string) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}

// Service interfaces
type UserService interface {
    GetUser(ctx context.Context, id string) (*User, error)
    GetUserByUsername(ctx context.Context, username string) (*User, error)
    CreateUser(ctx context.Context, username, email string) (*User, error)
    DeactivateUser(ctx context.Context, id string) error
}

// Repository implementation
type SQLUserRepository struct {
    db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) UserRepository {
    return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.db.QueryRowContext(ctx, 
        "SELECT id, username, email, active FROM users WHERE id = $1", 
        id,
    ).Scan(&user.ID, &user.Username, &user.Email, &user.Active)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user not found: %w", err)
        }
        return nil, fmt.Errorf("database error: %w", err)
    }

    return &user, nil
}

func (r *SQLUserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
    var user User
    err := r.db.QueryRowContext(ctx, 
        "SELECT id, username, email, active FROM users WHERE username = $1", 
        username,
    ).Scan(&user.ID, &user.Username, &user.Email, &user.Active)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user not found: %w", err)
        }
        return nil, fmt.Errorf("database error: %w", err)
    }

    return &user, nil
}

func (r *SQLUserRepository) Create(ctx context.Context, user *User) error {
    _, err := r.db.ExecContext(ctx,
        "INSERT INTO users (id, username, email, active) VALUES ($1, $2, $3, $4)",
        user.ID, user.Username, user.Email, user.Active,
    )
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    return nil
}

func (r *SQLUserRepository) Update(ctx context.Context, user *User) error {
    _, err := r.db.ExecContext(ctx,
        "UPDATE users SET username = $1, email = $2, active = $3 WHERE id = $4",
        user.Username, user.Email, user.Active, user.ID,
    )
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    return nil
}

func (r *SQLUserRepository) Delete(ctx context.Context, id string) error {
    _, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    return nil
}

// Service implementation
type DefaultUserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
    return &DefaultUserService{repo: repo}
}

func (s *DefaultUserService) GetUser(ctx context.Context, id string) (*User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *DefaultUserService) GetUserByUsername(ctx context.Context, username string) (*User, error) {
    return s.repo.GetByUsername(ctx, username)
}

func (s *DefaultUserService) CreateUser(ctx context.Context, username, email string) (*User, error) {
    // Check if user already exists
    existing, err := s.repo.GetByUsername(ctx, username)
    if err == nil && existing != nil {
        return nil, fmt.Errorf("username already taken")
    }

    // Create new user
    user := &User{
        ID:       uuid.New().String(),
        Username: username,
        Email:    email,
        Active:   true,
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *DefaultUserService) DeactivateUser(ctx context.Context, id string) error {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    user.Active = false
    return s.repo.Update(ctx, user)
}

func main() {
    // Load configuration
    cfg, err := config.New("config.yaml")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Create database connection
    db, err := sql.Open("pgx", cfg.GetString("database.url"))
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Create a new DI container
    container := di.NewContainer()

    // Register dependencies
    container.Register("config", cfg)
    container.Register("db", db)
    container.Register("userRepository", func(c di.Container) (interface{}, error) {
        db, err := c.Get("db")
        if err != nil {
            return nil, err
        }
        return NewSQLUserRepository(db.(*sql.DB)), nil
    })
    container.Register("userService", func(c di.Container) (interface{}, error) {
        repo, err := c.Get("userRepository")
        if err != nil {
            return nil, err
        }
        return NewUserService(repo.(UserRepository)), nil
    })

    // Resolve dependencies
    service, err := container.Get("userService")
    if err != nil {
        log.Fatalf("Failed to resolve dependencies: %v", err)
    }

    // Use the service
    userService := service.(UserService)

    // Create a new user
    ctx := context.Background()
    user, err := userService.CreateUser(ctx, "johndoe", "john.doe@example.com")
    if err != nil {
        log.Fatalf("Failed to create user: %v", err)
    }

    fmt.Printf("Created user: %+v\n", user)

    // Get the user by ID
    retrievedUser, err := userService.GetUser(ctx, user.ID)
    if err != nil {
        log.Fatalf("Failed to get user: %v", err)
    }

    fmt.Printf("Retrieved user: %+v\n", retrievedUser)
}
```

### Logging and Telemetry

```go
package main

import (
    "context"
    "net/http"
    "os"
    "time"

    "github.com/abitofhelp/servicelib/logging"
    "github.com/abitofhelp/servicelib/telemetry"
    "github.com/abitofhelp/servicelib/middleware"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

func main() {
    // Initialize logger
    logger, err := logging.NewLogger(logging.Config{
        Level:      "info",
        Format:     "json",
        OutputPath: "stdout",
        ErrorPath:  "stderr",
    })
    if err != nil {
        panic("Failed to initialize logger: " + err.Error())
    }
    defer logger.Sync()

    // Initialize telemetry
    telemetryConfig := telemetry.Config{
        ServiceName:    "example-service",
        ServiceVersion: "1.0.0",
        Environment:    "development",
        TracingEnabled: true,
        MetricsEnabled: true,
    }

    tp, mp, err := telemetry.InitTelemetry(telemetryConfig)
    if err != nil {
        logger.Fatal("Failed to initialize telemetry", "error", err)
    }
    defer func() {
        if err := tp.Shutdown(context.Background()); err != nil {
            logger.Error("Error shutting down tracer provider", "error", err)
        }
        if err := mp.Shutdown(context.Background()); err != nil {
            logger.Error("Error shutting down meter provider", "error", err)
        }
    }()

    // Create HTTP server with middleware
    mux := http.NewServeMux()

    // Add routes
    mux.HandleFunc("/", handleHome(logger))
    mux.HandleFunc("/api/users", handleUsers(logger))

    // Add Prometheus metrics endpoint
    mux.Handle("/metrics", promhttp.Handler())

    // Create middleware chain
    handler := middleware.Chain(
        mux,
        middleware.RequestID(),
        middleware.Logging(logger),
        middleware.Metrics(),
        middleware.Tracing("example-service"),
        middleware.Recovery(logger),
    )

    // Start server
    addr := ":8080"
    logger.Info("Starting server", "address", addr)

    server := &http.Server{
        Addr:         addr,
        Handler:      handler,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  30 * time.Second,
    }

    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        logger.Fatal("Server failed", "error", err)
    }
}

func handleHome(logger logging.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        logger.Info("Handling home request", "path", r.URL.Path, "method", r.Method)

        // Create a span for this handler
        ctx, span := otel.Tracer("example").Start(ctx, "handleHome")
        defer span.End()

        // Add attributes to the span
        span.SetAttributes(
            attribute.String("http.path", r.URL.Path),
            attribute.String("http.method", r.Method),
        )

        // Simulate some work
        time.Sleep(10 * time.Millisecond)

        // Create a child span
        _, childSpan := otel.Tracer("example").Start(ctx, "processHomeRequest")
        childSpan.SetAttributes(attribute.String("processing.type", "home"))
        time.Sleep(5 * time.Millisecond)
        childSpan.End()

        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Welcome to the example service!"))
    }
}

func handleUsers(logger logging.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        logger.Info("Handling users request", "path", r.URL.Path, "method", r.Method)

        // Create a span for this handler
        ctx, span := otel.Tracer("example").Start(ctx, "handleUsers")
        defer span.End()

        // Add attributes to the span
        span.SetAttributes(
            attribute.String("http.path", r.URL.Path),
            attribute.String("http.method", r.Method),
            attribute.String("user.id", r.URL.Query().Get("id")),
        )

        // Record an event
        span.AddEvent("Processing user request", trace.WithAttributes(
            attribute.String("request.id", r.Header.Get("X-Request-ID")),
        ))

        // Simulate some work
        time.Sleep(20 * time.Millisecond)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"users": [{"id": "1", "name": "John"}, {"id": "2", "name": "Jane"}]}`))
    }
}
```

## Component Documentation

### Authentication

The `auth` package provides implementations for JWT, OAuth2, and OIDC authentication methods:

- **JWT**: JSON Web Token implementation for stateless authentication
  - Token generation and validation
  - Support for standard claims and custom claims
  - Configurable signing methods (HMAC, RSA, ECDSA)

- **OAuth2**: OAuth 2.0 client implementation
  - Authorization Code, Client Credentials, and Password grant types
  - Token refresh and validation
  - State management for CSRF protection

- **OIDC**: OpenID Connect implementation
  - Discovery of provider configuration
  - ID token validation
  - User info retrieval

### Configuration

The `config` package provides a flexible configuration system that supports multiple sources and formats:

- **Multiple Sources**:
  - YAML and JSON files
  - Environment variables
  - Command-line flags
  - In-memory values

- **Features**:
  - Hierarchical configuration with dot notation
  - Default values
  - Type conversion
  - Configuration reloading
  - Validation

- **Adapters**: Easily create custom adapters for different configuration sources

### Context

The `context` package extends Go's standard context package with additional utilities:

- **Value Management**: Strongly typed context values
- **Timeout Management**: Utilities for working with context deadlines
- **Cancellation**: Simplified cancellation patterns
- **Propagation**: Utilities for propagating context values across service boundaries

### Database

The `db` package provides utilities for database connection management and operations:

- **Connection Management**:
  - Connection pooling
  - Automatic reconnection
  - Health checks

- **Supported Databases**:
  - PostgreSQL (via pgx)
  - SQLite
  - MongoDB

- **Features**:
  - Transaction management
  - Query execution with retries
  - Result mapping
  - Migrations

### Date

The `date` package provides utilities for working with dates and times:

- **Formatting**: Consistent date/time formatting
- **Parsing**: Robust date/time parsing with error handling
- **Comparison**: Date comparison utilities
- **Timezone**: Timezone conversion and management

### Dependency Injection

The `di` package provides a container-based dependency injection system:

- **Container Types**:
  - Base container
  - Service container
  - Repository container
  - Generic container

- **Features**:
  - Constructor injection
  - Singleton instances
  - Lazy initialization
  - Scoped instances
  - Circular dependency detection

### Error Handling

The `errors` package provides structured error types and handling:

- **Error Types**:
  - Domain errors
  - Infrastructure errors
  - Application errors
  - Validation errors

- **Features**:
  - Error wrapping with context
  - Error codes
  - Localized error messages
  - Stack traces
  - Error categorization

### GraphQL

The `graphql` package provides utilities for building GraphQL services:

- **Integration**: Integration with gqlgen
- **Error Handling**: Structured error handling for GraphQL
- **Middleware**: GraphQL-specific middleware
- **Resolvers**: Utilities for implementing resolvers

### Health Checks

The `health` package provides components for implementing health check endpoints:

- **Check Types**:
  - Liveness checks
  - Readiness checks
  - Startup checks

- **Features**:
  - Configurable check intervals
  - Automatic registration with HTTP server
  - Detailed health status reporting
  - Integration with Kubernetes probes

### Logging

The `logging` package provides structured logging with Zap:

- **Log Levels**: Debug, Info, Warn, Error, Fatal
- **Structured Logging**: Key-value pairs for better searchability
- **Output Formats**: JSON, console
- **Integration**: Context-aware logging
- **Performance**: High-performance logging with minimal allocations

### Middleware

The `middleware` package provides HTTP middleware components:

- **Authentication**: JWT authentication middleware
- **Logging**: Request/response logging
- **Metrics**: Request metrics collection
- **Tracing**: Distributed tracing
- **Recovery**: Panic recovery
- **CORS**: Cross-Origin Resource Sharing
- **Rate Limiting**: Request rate limiting

### Repository Pattern

The `repository` package provides generic repository implementations:

- **Generic Repository**: Type-safe repository implementation
- **CRUD Operations**: Create, Read, Update, Delete
- **Query Building**: Fluent query building
- **Pagination**: Offset and cursor-based pagination
- **Sorting**: Multi-field sorting

### Shutdown

The `shutdown` package provides graceful shutdown utilities:

- **Graceful Shutdown**: Orderly shutdown of services
- **Timeout Management**: Configurable shutdown timeouts
- **Dependency Order**: Shutdown in the correct order
- **Resource Cleanup**: Ensure all resources are properly released

### Signal Handling

The `signal` package provides OS signal handling:

- **Signal Types**: SIGINT, SIGTERM, SIGHUP
- **Custom Handlers**: Register custom signal handlers
- **Graceful Shutdown**: Integration with shutdown package

### String Utilities

The `stringutil` package provides string manipulation utilities:

- **Formatting**: String formatting utilities
- **Validation**: String validation
- **Transformation**: Case conversion, trimming, etc.
- **Generation**: Random string generation

### Telemetry

The `telemetry` package provides utilities for metrics, tracing, and monitoring:

- **Metrics**:
  - Prometheus integration
  - Counter, gauge, histogram, and summary metrics
  - Default metrics for HTTP, gRPC, and database

- **Tracing**:
  - OpenTelemetry integration
  - Span creation and management
  - Context propagation
  - Automatic instrumentation for HTTP and gRPC

- **Monitoring**:
  - Health check integration
  - Alerting utilities
  - Dashboard templates

### Transaction

The `transaction` package provides utilities for managing distributed transactions:

- **Saga Pattern**: Implementation of the Saga pattern for distributed transactions
- **Compensation**: Transaction compensation for rollback
- **Coordination**: Transaction coordination across services

### Validation

The `validation` package provides request and data validation:

- **Integration**: Integration with go-playground/validator
- **Custom Validators**: Define custom validation rules
- **Validation Middleware**: HTTP request validation
- **Error Handling**: Structured validation errors

## Building and Testing

ServiceLib uses Go modules for dependency management and Make for build automation.

### Prerequisites

- Go 1.24 or higher
- Make (optional, for using the Makefile)

### Build Commands

```bash
# Build the library
make build

# Run tests
make test

# Run tests with coverage
make coverage

# Run linter
make lint

# Format code
make fmt

# Check for security vulnerabilities
make security
```

## Troubleshooting

### Common Issues

#### Connection Pooling

**Issue**: Database connections are not being properly released, leading to connection pool exhaustion.

**Solution**: Ensure that all database operations properly close their resources, especially in error cases. Use the `defer` statement to ensure connections are returned to the pool:

```go
conn, err := db.GetConnection(ctx)
if err != nil {
    return err
}
defer conn.Close() // This ensures the connection is returned to the pool
```

#### Memory Leaks

**Issue**: Memory usage grows over time, indicating potential memory leaks.

**Solution**: Use the telemetry package to monitor memory usage and identify leaks. Common causes include:

- Forgetting to close response bodies
- Goroutines that never terminate
- Large objects stored in context values

#### Circular Dependencies

**Issue**: Dependency injection container fails with circular dependency errors.

**Solution**: Restructure your dependencies to break the cycle. Consider:

- Using interfaces to break direct dependencies
- Introducing a mediator or facade
- Using lazy initialization for some dependencies

### Debugging

#### Enabling Debug Logging

To enable debug logging for troubleshooting:

```go
logger, _ := logging.NewLogger(logging.Config{
    Level: "debug",
    Format: "console", // More readable for debugging
})
```

#### Tracing Requests

For detailed request tracing:

1. Enable the tracing middleware
2. Set the sampling rate to 1.0 (100%)
3. Use the OpenTelemetry UI or Jaeger to view traces

## Best Practices

### Service Structure

- **Layered Architecture**: Organize your service with clear separation between:
  - API/Transport layer (HTTP, gRPC)
  - Service layer (business logic)
  - Repository layer (data access)

- **Dependency Injection**: Use the DI container to manage dependencies and make testing easier

- **Configuration**: Externalize all configuration and use environment variables for deployment-specific settings

### Error Handling

- **Structured Errors**: Use the errors package to create structured errors with context

```go
if err != nil {
    return errors.NewInfrastructureError("database_error", "Failed to query database", err)
}
```

- **Error Categorization**: Categorize errors to handle them appropriately at the API boundary

```go
switch {
case errors.IsNotFound(err):
    return http.StatusNotFound, errorResponse(err)
case errors.IsValidation(err):
    return http.StatusBadRequest, errorResponse(err)
case errors.IsUnauthorized(err):
    return http.StatusUnauthorized, errorResponse(err)
default:
    return http.StatusInternalServerError, errorResponse(err)
}
```

### Performance Optimization

- **Connection Pooling**: Configure database connection pools based on expected load

```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(5 * time.Minute)
```

- **Caching**: Use caching for frequently accessed, rarely changed data

- **Pagination**: Always implement pagination for endpoints that return collections

### Testing

- **Unit Tests**: Test each component in isolation using mocks

- **Integration Tests**: Test the integration between components

- **End-to-End Tests**: Test the complete service flow

- **Load Tests**: Test performance under load to identify bottlenecks

## Architecture and Design

ServiceLib is designed with the following architectural principles:

### Modularity

Each package in ServiceLib is designed to be used independently or together with other packages. This allows you to use only the components you need without bringing in unnecessary dependencies.

### Hexagonal Architecture

The library encourages a hexagonal (ports and adapters) architecture:

- **Core Domain**: Business logic independent of external concerns
- **Ports**: Interfaces defining how the core interacts with the outside world
- **Adapters**: Implementations of ports for specific technologies

This architecture makes it easier to:
- Replace implementations without changing business logic
- Test business logic in isolation
- Adapt to changing requirements

### Design Patterns

ServiceLib implements several design patterns:

- **Repository Pattern**: Abstracts data access behind interfaces
- **Dependency Injection**: Manages dependencies and facilitates testing
- **Factory Pattern**: Creates complex objects with consistent configuration
- **Adapter Pattern**: Converts interfaces to work with different systems
- **Observer Pattern**: Implements event-based communication

## Compatibility and Versioning

### Version Compatibility

ServiceLib follows semantic versioning (SemVer):

- **Major version** (X.y.z): Incompatible API changes
- **Minor version** (x.Y.z): Backwards-compatible functionality
- **Patch version** (x.y.Z): Backwards-compatible bug fixes

### Go Version Compatibility

- **Minimum Go version**: 1.24
- **Tested Go versions**: 1.24, 1.25

### Dependencies

ServiceLib has the following major dependencies:

- **zap**: Structured logging
- **prometheus**: Metrics collection
- **opentelemetry**: Distributed tracing
- **validator**: Request validation
- **pgx**: PostgreSQL driver
- **gqlgen**: GraphQL implementation

### Backward Compatibility Guarantees

- No breaking changes will be introduced in minor or patch releases
- Deprecated features will be marked with `Deprecated` in the documentation
- Deprecated features will be removed only in major version releases
- Migration guides will be provided for major version upgrades

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

### Coding Standards

- Follow Go best practices and style guidelines
- Write tests for new functionality
- Document public APIs
- Keep backward compatibility in mind

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
