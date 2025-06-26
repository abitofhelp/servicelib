# ServiceLib
[![codecov](https://codecov.io/gh/abitofhelp/servicelib/graph/badge.svg)](https://codecov.io/gh/abitofhelp/servicelib)
[![Go Report Card](https://goreportcard.com/badge/github.com/abitofhelp/servicelib)](https://goreportcard.com/report/github.com/abitofhelp/servicelib)
[![GoDoc](https://godoc.org/github.com/abitofhelp/servicelib?status.svg)](https://godoc.org/github.com/abitofhelp/servicelib)
## Testing Coverage
### Top Level represents the entire project; The bottom level represents individual files.
### Click the image for more details.
[![codecov](https://codecov.io/gh/abitofhelp/servicelib/graphs/icicle.svg)](https://codecov.io/gh/abitofhelp/servicelib)

## Overview

ServiceLib is a comprehensive Go library designed to provide a robust foundation for building scalable, maintainable, and production-ready microservices. It offers a collection of packages that address common challenges in service development, from authentication and configuration to error handling and telemetry.

## Features

- **Authentication & Authorization**: Secure your services with JWT, OIDC, and role-based access control
- **Configuration Management**: Flexible configuration with support for multiple sources and formats
- **Error Handling**: Structured error types with context, stack traces, and categorization
- **Database Access**: Connection management, health checks, and transaction support
- **Dependency Injection**: Simple yet powerful DI container for managing service dependencies
- **Telemetry**: Integrated logging, metrics, and distributed tracing
- **Health Checks**: Standardized health check endpoints and status reporting
- **Middleware**: Common HTTP middleware for logging, error handling, and more
- **Retry & Circuit Breaking**: Resilience patterns for handling transient failures
- **Validation**: Comprehensive input validation utilities

## Installation

```bash
go get github.com/abitofhelp/servicelib
```

## Quick Start

### Basic HTTP Server

Here's a simple example of creating an HTTP server with ServiceLib:

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/middleware"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, err := logging.NewLogger("info", true)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Create a context logger
	contextLogger := logging.NewContextLogger(logger)

	// Create a simple HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, ServiceLib!"))
	})

	// Add middleware for logging, metrics, and recovery
	handler := middleware.Chain(
		mux,
		middleware.WithRequestID(context.Background()),
		middleware.Logging(contextLogger),
		middleware.Recovery(contextLogger),
	)

	// Start the server
	contextLogger.Info(context.Background(), "Starting server", zap.String("address", ":8080"))
	if err := http.ListenAndServe(":8080", handler); err != nil {
		contextLogger.Fatal(context.Background(), "Server failed", zap.Error(err))
	}
}
```

### Error Handling

ServiceLib provides a comprehensive error handling system:

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/abitofhelp/servicelib/errors"
	errhttp "github.com/abitofhelp/servicelib/errors/http"
	errlog "github.com/abitofhelp/servicelib/errors/log"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// GetUserByID retrieves a user by ID with proper error handling
func GetUserByID(ctx context.Context, logger *logging.ContextLogger, id string) (*User, error) {
	// Validate input
	if id == "" {
		// Create a validation error
		err := errors.NewValidationError("User ID is required", "id", nil)

		// Log the error
		errlog.LogError(ctx, logger, err)

		return nil, err
	}

	// Simulate a database error
	if id == "db-error" {
		dbErr := sql.ErrNoRows
		err := errors.NewDatabaseError("Failed to query user", "SELECT", "users", dbErr)

		// Log the error
		errlog.LogError(ctx, logger, err)

		return nil, err
	}

	// Return a user
	return &User{ID: id, Name: "John Doe"}, nil
}

// HTTP handler with error handling
func GetUserHandler(w http.ResponseWriter, r *http.Request, service *UserService) {
	// Get user ID from request
	id := r.URL.Query().Get("id")

	// Get user from service
	user, err := service.GetUserByID(r.Context(), id)
	if err != nil {
		// Write error response with appropriate status code
		errhttp.WriteError(w, err)
		return
	}

	// Write success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"id":"%s","name":"%s"}`, user.ID, user.Name)
}
```

For more examples, see the [Examples directory](./EXAMPLES/README.md).

## Packages

ServiceLib is organized into the following packages:

- [auth](./auth/README.md) - Authentication and authorization
- [cache](./cache/README.md) - Caching utilities
- [circuit](./circuit/README.md) - Circuit breaker implementation
- [config](./config/README.md) - Configuration management
- [context](./context/README.md) - Context utilities
- [date](./date/README.md) - Date and time utilities
- [db](./db/README.md) - Database access and management
- [di](./di/README.md) - Dependency injection
- [env](./env/README.md) - Environment variable utilities
- [errors](./errors/README.md) - Error handling, management, and recovery patterns
- [graphql](./graphql/README.md) - GraphQL utilities
- [health](./health/README.md) - Health check utilities
- [logging](./logging/README.md) - Structured logging
- [middleware](./middleware/README.md) - HTTP middleware
- [model](./model/README.md) - Model utilities
- [rate](./rate/README.md) - Rate limiting
- [repository](./repository/README.md) - Repository pattern implementation
- [retry](./retry/README.md) - Retry utilities
- [shutdown](./shutdown/README.md) - Graceful shutdown utilities
- [signal](./signal/README.md) - Signal handling
- [stringutil](./stringutil/README.md) - String utilities
- [telemetry](./telemetry/README.md) - Telemetry (metrics, tracing)
- [transaction](./transaction/README.md) - Transaction management
- [validation](./validation/README.md) - Input validation
- [valueobject](./valueobject/README.md) - Value object implementations

## Examples

For complete, runnable examples of each component, see the [EXAMPLES](./EXAMPLES/README.md) directory.

## Architecture Diagrams

The following UML diagrams provide a visual representation of the ServiceLib architecture:

- [Architecture Overview](./DOCS/diagrams/svg/Architecture%20Overview.svg) - High-level overview of the ServiceLib architecture based on Clean Architecture, DDD, and Hexagonal Architecture
- [Package Diagram](./DOCS/diagrams/svg/Package%20Diagram.svg) - Diagram showing the packages in ServiceLib and their relationships
- [HTTP Request Sequence](./DOCS/diagrams/svg/HTTP%20Request%20Sequence.svg) - Sequence diagram illustrating the HTTP request processing flow
- [Errors Package Class Diagram](./DOCS/diagrams/svg/Errors%20Package%20Class%20Diagram.svg) - Class diagram for the errors package

## Contributing

Contributions to ServiceLib are welcome! Please see the [Contributing Guide](./CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
