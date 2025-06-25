# ServiceLib Error Handling
This package provides a comprehensive error handling system for the application. It includes error codes, HTTP status mapping, contextual information, and utilities for creating, wrapping, and serializing errors.

For detailed information about the error handling design, including design principles, error type hierarchy, relationships between packages, and best practices, see [Error Handling Design](error_handling_design.md).

## Package Structure

The error handling package is organized into several sub-packages:
  - `core`: Core error handling functionality (error codes, base error type, utility functions)
  - `domain`: Domain-specific error types (validation errors, business rule errors, not found errors)
  - `infra`: Infrastructure-related error types (database errors, network errors, external service errors)
  - `app`: Application-level error types (configuration errors, authentication errors, authorization errors)
  - `http`: HTTP-related error utilities
  - `log`: Logging integration
  - `metrics`: Metrics integration
  - `trace`: Tracing integration
  - `utils`: Utility functions for error handling


## Overview

Brief description of the errors and its purpose in the ServiceLib library.

## Features

- **Clear Error Type Hierarchy**:
  - BaseError: The foundation for all error types
  - DomainError: Domain-specific errors (validation, business rules, not found)
  - InfrastructureError: Infrastructure-related errors (database, network, external services)
  - ApplicationError: Application-level errors (configuration, authentication, authorization)

- **Consistent Error Creation and Wrapping**:
  - Factory functions for creating different error types
  - Wrap functions for adding context to errors
  - Support for error cause chains

- **Error Codes and HTTP Status Mapping**:
  - Predefined error codes for common error scenarios
  - Automatic mapping of error codes to HTTP status codes

- **Contextual Information**:
  - Operation name
  - Source location (file and line)
  - Additional details as key-value pairs
  - Error cause chain

- **Error Categorization and Type Checking**:
  - Type-specific checking functions (IsValidationError, IsDatabaseError, etc.)
  - Support for standard errors.Is and errors.As functions

- **JSON Serialization**:
  - Convert errors to JSON for API responses
  - Include all contextual information in JSON output

- **Integration with Logging, Metrics, and Tracing**:
  - Log errors with all contextual information
  - Record error metrics with error type and code
  - Add error information to traces


## Installation

```bash
go get github.com/abitofhelp/servicelib/errors
```

## Usage

### Basic Error Creation

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/errors"
)

func main() {
    // Create a simple error
    err := errors.New(errors.InvalidInputCode, "Invalid input provided")
    fmt.Println(err) // Output: Invalid input provided

    // Create a domain error
    err = errors.NewDomainError(errors.NotFoundCode, "User not found", nil)
    fmt.Println(err) // Output: User not found

    // Create a validation error
    err = errors.NewValidationError("Field 'name' is required", "name", nil)
    fmt.Println(err) // Output: Field 'name' is required

    // Create a business rule error
    err = errors.NewBusinessRuleError("User must be at least 18 years old", "MinimumAge", nil)
    fmt.Println(err) // Output: User must be at least 18 years old

    // Create a not found error
    err = errors.NewNotFoundError("User", "123", nil)
    fmt.Println(err) // Output: User with ID 123 not found

    // Create an infrastructure error
    err = errors.NewDatabaseError("Failed to query database", "SELECT", "users", nil)
    fmt.Println(err) // Output: Failed to query database

    // Create an application error
    err = errors.NewConfigurationError("Invalid configuration value", "MAX_CONNECTIONS", "abc", nil)
    fmt.Println(err) // Output: Invalid configuration value
}
```

### Error Wrapping

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/abitofhelp/servicelib/errors"
)

func main() {
    // Wrap a standard error
    originalErr := sql.ErrNoRows
    wrappedErr := errors.Wrap(originalErr, errors.NotFoundCode, "Failed to find user")
    fmt.Println(wrappedErr) // Output: Failed to find user: sql: no rows in result set

    // Wrap with an operation
    wrappedWithOp := errors.WrapWithOperation(originalErr, errors.NotFoundCode, "Failed to find user", "GetUserByID")
    fmt.Println(wrappedWithOp) // Output: GetUserByID: Failed to find user: sql: no rows in result set

    // Wrap with details
    details := map[string]interface{}{
        "user_id": "123",
        "table":   "users",
    }
    wrappedWithDetails := errors.WrapWithDetails(originalErr, errors.NotFoundCode, "Failed to find user", details)
    fmt.Println(wrappedWithDetails) // Output: Failed to find user: sql: no rows in result set

    // Unwrap to get the original error
    unwrapped := errors.Unwrap(wrappedErr)
    fmt.Println(unwrapped) // Output: sql: no rows in result set

    // Check if an error is of a specific type
    if errors.Is(wrappedErr, sql.ErrNoRows) {
        fmt.Println("The original error was sql.ErrNoRows")
    }
}
```

### Error Type Checking

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/errors"
)

func main() {
    // Create errors of different types
    notFoundErr := errors.NewNotFoundError("User", "123", nil)
    validationErr := errors.NewValidationError("Email is invalid", "email", nil)
    dbErr := errors.NewDatabaseError("Failed to query database", "SELECT", "users", nil)
    authErr := errors.NewAuthenticationError("Invalid credentials", "john.doe", nil)

    // Check error types
    fmt.Printf("Is not found error: %v\n", errors.IsNotFoundError(notFoundErr))
    fmt.Printf("Is validation error: %v\n", errors.IsValidationError(validationErr))
    fmt.Printf("Is database error: %v\n", errors.IsDatabaseError(dbErr))
    fmt.Printf("Is authentication error: %v\n", errors.IsAuthenticationError(authErr))

    // Use in error handling
    handleError(notFoundErr)
    handleError(validationErr)
    handleError(dbErr)
    handleError(authErr)
}

func handleError(err error) {
    switch {
    case errors.IsNotFoundError(err):
        fmt.Println("Handle not found error:", err)
    case errors.IsValidationError(err):
        fmt.Println("Handle validation error:", err)
    case errors.IsDatabaseError(err):
        fmt.Println("Handle database error:", err)
    case errors.IsAuthenticationError(err):
        fmt.Println("Handle authentication error:", err)
    default:
        fmt.Println("Handle generic error:", err)
    }
}
```

### HTTP Status Mapping

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/errors"
    "net/http"
)

func main() {
    // Create errors of different types
    notFoundErr := errors.NewNotFoundError("User", "123", nil)
    validationErr := errors.NewValidationError("Email is invalid", "email", nil)
    dbErr := errors.NewDatabaseError("Failed to query database", "SELECT", "users", nil)
    authErr := errors.NewAuthenticationError("Invalid credentials", "john.doe", nil)

    // Get HTTP status codes
    fmt.Printf("Not found error HTTP status: %d\n", errors.GetHTTPStatus(notFoundErr))
    fmt.Printf("Validation error HTTP status: %d\n", errors.GetHTTPStatus(validationErr))
    fmt.Printf("Database error HTTP status: %d\n", errors.GetHTTPStatus(dbErr))
    fmt.Printf("Authentication error HTTP status: %d\n", errors.GetHTTPStatus(authErr))

    // Use in HTTP handler
    handleHTTPError(notFoundErr, http.ResponseWriter(nil))
}

func handleHTTPError(err error, w http.ResponseWriter) {
    status := errors.GetHTTPStatus(err)
    if status == 0 {
        status = http.StatusInternalServerError
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write([]byte(errors.ToJSON(err)))
}
```

### JSON Serialization

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/errors"
)

func main() {
    // Create an error with details
    details := map[string]interface{}{
        "user_id": "123",
        "action":  "create_order",
        "status":  "failed",
    }
    err := errors.WrapWithDetails(
        errors.New(errors.ValidationErrorCode, "Invalid input"),
        errors.ValidationErrorCode,
        "Failed to process request",
        details,
    )

    // Convert the error to JSON
    jsonStr := errors.ToJSON(err)
    fmt.Println(jsonStr)

    // Output:
    // {
    //   "message": "Failed to process request: Invalid input",
    //   "code": "VALIDATION_ERROR",
    //   "details": {
    //     "user_id": "123",
    //     "action": "create_order",
    //     "status": "failed"
    //   },
    //   "source": "main.go",
    //   "line": 15
    // }
}
```

### Integration with Logging

```go
package main

import (
    "context"
    "github.com/abitofhelp/servicelib/errors"
    "log"
)

func main() {
    // Create a context
    ctx := context.Background()

    // Create an error
    err := errors.NewDatabaseError("Failed to query database", "SELECT", "users", nil)

    // Log the error with context
    logError(ctx, err)
}

func logError(ctx context.Context, err error) {
    // Get error details
    code := ""
    if e, ok := err.(interface{ GetCode() errors.ErrorCode }); ok {
        code = string(e.GetCode())
    }

    // Log the error with context
    log.Printf(
        "Error occurred: code=%s, message=%s, http_status=%d",
        code,
        err.Error(),
        errors.GetHTTPStatus(err),
    )
}
```


## Quick Start

See the [Quick Start example](../EXAMPLES/errors/quickstart_example.go) for a complete, runnable example of how to use the errors.

## Configuration

See the [Configuration example](../EXAMPLES/errors/configuration_example.go) for a complete, runnable example of how to configure the errors.

## API Documentation


### Core Types

Description of the main types provided by the errors.

#### Type 1

Description of Type 1 and its purpose.

See the [Type 1 example](../EXAMPLES/errors/type1_example.go) for a complete, runnable example of how to use Type 1.

### Key Methods

Description of the key methods provided by the errors.

#### Method 1

Description of Method 1 and its purpose.

See the [Method 1 example](../EXAMPLES/errors/method1_example.go) for a complete, runnable example of how to use Method 1.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/errors/basic_usage_example.go) - Shows basic usage of the errors
- [Advanced Configuration](../EXAMPLES/errors/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/errors/error_handling_example.go) - Shows how to handle errors

## Best Practices

### 1. Use the Appropriate Error Type

Choose the most specific error type for your situation:

```go
// For domain validation errors
err := errors.NewValidationError("Email is invalid", "email", nil)

// For not found errors
err := errors.NewNotFoundError("User", "123", nil)

// For database errors
err := errors.NewDatabaseError("Failed to query database", "SELECT", "users", nil)

// For authentication errors
err := errors.NewAuthenticationError("Invalid credentials", "john.doe", nil)
```

### 2. Add Context to Errors

Wrap errors with additional context:

```go
// Wrap with operation name
err = errors.WrapWithOperation(err, errors.DatabaseErrorCode, "Database query failed", "GetUserByID")

// Wrap with details
details := map[string]interface{}{
    "user_id": "123",
    "query":   "SELECT * FROM users WHERE id = ?",
}
err = errors.WrapWithDetails(err, errors.DatabaseErrorCode, "Database query failed", details)
```

### 3. Check Error Types

Use the type checking functions to handle different error types:

```go
switch {
case errors.IsNotFoundError(err):
    // Handle not found error
case errors.IsValidationError(err):
    // Handle validation error
case errors.IsDatabaseError(err):
    // Handle database error
default:
    // Handle other errors
}
```
```

### 4. Map Errors to HTTP Status Codes

Use the GetHTTPStatus function to map errors to HTTP status codes:

```go
func handleHTTPError(w http.ResponseWriter, err error) {
    status := errors.GetHTTPStatus(err)
    if status == 0 {
        status = http.StatusInternalServerError
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write([]byte(errors.ToJSON(err)))
}
```

### 5. Provide Detailed Error Messages

Include detailed error messages that help identify the issue:

```go
// Instead of this:
err := errors.New(errors.ValidationErrorCode, "Invalid input")

// Do this:
err := errors.NewValidationError("Email must be a valid email address", "email", nil)
```

### 6. Use Standard Error Codes

Use the standard error codes defined in the package:

```go
// Standard error codes
errors.NotFoundCode
errors.InvalidInputCode
errors.DatabaseErrorCode
errors.InternalErrorCode
errors.UnauthorizedCode
errors.ForbiddenCode
errors.ValidationErrorCode
errors.BusinessRuleViolationCode
```

### 7. Include Source Information

The error handling system automatically includes source file and line information, which helps with debugging:

```go
// The error will include the source file and line number
err := errors.New(errors.InternalErrorCode, "Something went wrong")
```


## Troubleshooting

### Common Issues

#### Issue 1

Description of issue 1 and how to resolve it.

#### Issue 2

Description of issue 2 and how to resolve it.

## Related Components

- [Component 1](../errors1/README.md) - Description of how this errors relates to Component 1
- [Component 2](../errors2/README.md) - Description of how this errors relates to Component 2

## Contributing

Contributions to this errors are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
