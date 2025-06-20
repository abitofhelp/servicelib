# Error Handling Examples

This directory contains examples demonstrating how to use the `errors` package, which provides structured error types and handling with rich context information for Go applications. The package extends Go's standard error handling with features like error categorization, error wrapping with context, and error codes.

## Examples

### 1. Basic Error Creation Example

[basic_error_creation_example.go](basic_error_creation_example.go)

Demonstrates how to create different types of errors.

Key concepts:
- Creating simple errors with `errors.New`
- Creating errors with codes using `errors.NewWithCode`
- Creating domain-specific errors with `errors.NewDomainError`
- Creating infrastructure errors with `errors.NewInfrastructureError`
- Creating application errors with `errors.NewApplicationError`
- Creating validation errors with `errors.NewValidationError`

### 2. Error Wrapping Example

[error_wrapping_example.go](error_wrapping_example.go)

Shows how to wrap errors to add context as they propagate up the call stack.

Key concepts:
- Wrapping errors with `errors.Wrap`
- Wrapping errors with codes using `errors.WrapWithCode`
- Unwrapping errors with `errors.Unwrap`
- Checking error types with `errors.Is`

### 3. Error Context Example

[error_context_example.go](error_context_example.go)

Demonstrates how to add and retrieve context information from errors.

Key concepts:
- Creating errors with context using `errors.NewWithContext`
- Getting context from errors with `errors.GetContext`
- Getting specific context values with `errors.GetContextValue`
- Adding context to existing errors with `errors.AddContext`

### 4. Error Categorization Example

[error_categorization_example.go](error_categorization_example.go)

Shows how to categorize and handle different types of errors.

Key concepts:
- Checking error categories with `errors.IsNotFound`, `errors.IsValidation`, etc.
- Handling errors based on their category
- Creating custom error categories

### 5. Stack Traces Example

[stack_traces_example.go](stack_traces_example.go)

Demonstrates how to include and retrieve stack traces from errors.

Key concepts:
- Creating errors with stack traces using `errors.NewWithStack`
- Retrieving stack traces with `errors.StackTrace`
- Understanding stack trace output

### 6. Localized Error Messages Example

[localized_messages_example.go](localized_messages_example.go)

Shows how to create and retrieve localized error messages.

Key concepts:
- Creating errors with localized messages using `errors.NewLocalized`
- Retrieving localized messages with `errors.GetLocalizedMessage`
- Using fallback languages with `errors.GetLocalizedMessageWithFallback`

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/errors/basic_error_creation_example.go
```

## Additional Resources

For more information about the errors package, see the [errors package documentation](../../errors/README.md).