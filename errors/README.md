# Error Handling Package

The `errors` package provides structured error types and handling with rich context information for Go applications. It extends Go's standard error handling with features like error categorization, error wrapping with context, and error codes.

## Features

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
    err := errors.New("something went wrong")
    fmt.Println(err) // Output: something went wrong

    // Create an error with a code
    err = errors.NewWithCode("invalid_input", "Invalid input provided")
    fmt.Println(err) // Output: invalid_input: Invalid input provided

    // Create a domain error
    err = errors.NewDomainError("user_not_found", "User not found")
    fmt.Println(err) // Output: domain error: user_not_found: User not found

    // Create an infrastructure error
    err = errors.NewInfrastructureError("database_connection", "Failed to connect to database")
    fmt.Println(err) // Output: infrastructure error: database_connection: Failed to connect to database

    // Create an application error
    err = errors.NewApplicationError("invalid_state", "Application is in an invalid state")
    fmt.Println(err) // Output: application error: invalid_state: Application is in an invalid state

    // Create a validation error
    err = errors.NewValidationError("required_field", "Field 'name' is required")
    fmt.Println(err) // Output: validation error: required_field: Field 'name' is required
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
    wrappedErr := errors.Wrap(originalErr, "failed to find user")
    fmt.Println(wrappedErr) // Output: failed to find user: sql: no rows in result set

    // Wrap with a code
    wrappedWithCode := errors.WrapWithCode(originalErr, "user_not_found", "failed to find user")
    fmt.Println(wrappedWithCode) // Output: user_not_found: failed to find user: sql: no rows in result set

    // Wrap multiple times
    deeplyWrapped := errors.Wrap(wrappedErr, "error in GetUserByID")
    fmt.Println(deeplyWrapped) // Output: error in GetUserByID: failed to find user: sql: no rows in result set

    // Unwrap to get the original error
    unwrapped := errors.Unwrap(deeplyWrapped)
    fmt.Println(unwrapped) // Output: failed to find user: sql: no rows in result set

    // Check if an error is of a specific type
    if errors.Is(deeplyWrapped, sql.ErrNoRows) {
        fmt.Println("The original error was sql.ErrNoRows")
    }
}
```

### Error Context

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/errors"
)

func main() {
    // Create an error with context
    err := errors.NewWithContext("failed to process request", map[string]interface{}{
        "user_id": "123",
        "action":  "create_order",
        "status":  "failed",
    })

    // Get context from the error
    if ctx, ok := errors.GetContext(err); ok {
        fmt.Printf("Error context: %v\n", ctx)
    }

    // Get a specific context value
    if userID, ok := errors.GetContextValue(err, "user_id"); ok {
        fmt.Printf("User ID: %v\n", userID)
    }

    // Add more context to an existing error
    enrichedErr := errors.AddContext(err, "request_id", "req-456")
    
    // Get the enriched context
    if ctx, ok := errors.GetContext(enrichedErr); ok {
        fmt.Printf("Enriched context: %v\n", ctx)
    }
}
```

### Error Categorization

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/errors"
)

func main() {
    // Create errors of different types
    notFoundErr := errors.NewDomainError("user_not_found", "User not found")
    validationErr := errors.NewValidationError("invalid_email", "Email is invalid")
    authErr := errors.NewApplicationError("unauthorized", "User is not authorized")
    dbErr := errors.NewInfrastructureError("db_connection", "Failed to connect to database")

    // Check error categories
    fmt.Printf("Is not found: %v\n", errors.IsNotFound(notFoundErr))
    fmt.Printf("Is validation: %v\n", errors.IsValidation(validationErr))
    fmt.Printf("Is unauthorized: %v\n", errors.IsUnauthorized(authErr))
    fmt.Printf("Is infrastructure: %v\n", errors.IsInfrastructure(dbErr))

    // Use in error handling
    handleError(notFoundErr)
    handleError(validationErr)
    handleError(authErr)
    handleError(dbErr)
}

func handleError(err error) {
    switch {
    case errors.IsNotFound(err):
        fmt.Println("Handle not found error:", err)
    case errors.IsValidation(err):
        fmt.Println("Handle validation error:", err)
    case errors.IsUnauthorized(err):
        fmt.Println("Handle unauthorized error:", err)
    case errors.IsInfrastructure(err):
        fmt.Println("Handle infrastructure error:", err)
    default:
        fmt.Println("Handle generic error:", err)
    }
}
```

### Stack Traces

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/errors"
)

func main() {
    // Create an error with a stack trace
    err := errors.NewWithStack("something went wrong")

    // Print the error with stack trace
    fmt.Println(errors.StackTrace(err))

    // Create a function that returns an error
    err = someFunction()

    // Print the stack trace
    fmt.Println(errors.StackTrace(err))
}

func someFunction() error {
    return anotherFunction()
}

func anotherFunction() error {
    return errors.NewWithStack("error in anotherFunction")
}
```

### Error Codes

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/errors"
)

func main() {
    // Create an error with a code
    err := errors.NewWithCode("invalid_input", "Invalid input provided")

    // Get the error code
    if code, ok := errors.GetCode(err); ok {
        fmt.Printf("Error code: %s\n", code)
    }

    // Check if an error has a specific code
    if errors.HasCode(err, "invalid_input") {
        fmt.Println("This is an invalid input error")
    }

    // Create a wrapped error with a code
    originalErr := fmt.Errorf("value out of range")
    wrappedErr := errors.WrapWithCode(originalErr, "validation_error", "Input validation failed")

    // Get the code from the wrapped error
    if code, ok := errors.GetCode(wrappedErr); ok {
        fmt.Printf("Wrapped error code: %s\n", code)
    }
}
```

### Localized Error Messages

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/errors"
)

func main() {
    // Create an error with localized messages
    err := errors.NewLocalized("invalid_input", map[string]string{
        "en": "Invalid input provided",
        "es": "Entrada no válida proporcionada",
        "fr": "Entrée invalide fournie",
    })

    // Get the default message (usually English)
    fmt.Println(err) // Output: invalid_input: Invalid input provided

    // Get a localized message
    if msg, ok := errors.GetLocalizedMessage(err, "es"); ok {
        fmt.Printf("Spanish error message: %s\n", msg)
    }

    // Get a localized message with fallback
    msg := errors.GetLocalizedMessageWithFallback(err, "de", "en")
    fmt.Printf("German message (fallback to English): %s\n", msg)
}
```

## Best Practices

1. **Use Structured Errors**: Create structured errors with codes and context information to make debugging easier.

   ```go
   return errors.NewWithContext("failed to process order", map[string]interface{}{
       "order_id": orderID,
       "user_id": userID,
       "status": "failed",
   })
   ```

2. **Error Categorization**: Categorize errors to handle them appropriately at the API boundary.

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

3. **Wrap Errors**: Wrap errors to add context as they propagate up the call stack.

   ```go
   if err := repository.GetUser(id); err != nil {
       return errors.Wrap(err, "failed to get user")
   }
   ```

4. **Error Codes**: Use consistent error codes across your application.

   ```go
   const (
       ErrNotFound      = "not_found"
       ErrInvalidInput  = "invalid_input"
       ErrUnauthorized  = "unauthorized"
       ErrInternal      = "internal_error"
   )
   ```

5. **Stack Traces**: Include stack traces for unexpected errors to aid debugging.

   ```go
   if err != nil && !isExpectedError(err) {
       return errors.NewWithStack(fmt.Sprintf("unexpected error: %v", err))
   }
   ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.