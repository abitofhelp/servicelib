# Error Handling Design

## Overview

This document explains the design of the error handling system in the ServiceLib library. It covers the design principles, error type hierarchy, relationships between different error packages, and best practices for error handling.

## Design Principles

The error handling system is designed with the following principles in mind:

1. **Simplicity**: The error handling system should be simple to use and understand.
2. **Consistency**: Error creation, wrapping, and handling should be consistent throughout the codebase.
3. **Extensibility**: The system should be extensible to accommodate domain-specific error types.
4. **Integration**: The system should integrate well with logging, metrics, and tracing.
5. **Performance**: The system should be efficient and not introduce significant overhead.

## Error Type Hierarchy

The error handling system has a clear hierarchy of error types:

```
BaseError (implements error interface)
├── DomainError (domain-specific errors)
│   ├── ValidationError (validation failures)
│   ├── BusinessRuleError (business rule violations)
│   └── NotFoundError (resource not found)
├── InfrastructureError (infrastructure-related errors)
│   ├── DatabaseError (database operation failures)
│   ├── NetworkError (network-related errors)
│   └── ExternalServiceError (external service call failures)
└── ApplicationError (application-level errors)
    ├── ConfigurationError (configuration issues)
    ├── AuthenticationError (authentication failures)
    └── AuthorizationError (authorization failures)
```

This hierarchy allows for specific error types for different scenarios while maintaining a consistent interface and behavior.

## Package Structure

The error handling system is organized into the following packages:

- **core**: Core error handling functionality
  - BaseError: The foundation for all error types
  - ErrorCode: Error codes for categorizing errors
  - HTTP status mapping: Mapping error codes to HTTP status codes
  - Utility functions: Functions for error handling

- **domain**: Domain-specific error types
  - DomainError: Base type for domain-specific errors
  - ValidationError: For validation failures
  - BusinessRuleError: For business rule violations
  - NotFoundError: For resource not found errors

- **app**: Application-level error types
  - ApplicationError: Base type for application-level errors
  - ConfigurationError: For configuration issues
  - AuthenticationError: For authentication failures
  - AuthorizationError: For authorization failures

- **infra**: Infrastructure-related error types
  - InfrastructureError: Base type for infrastructure-related errors
  - DatabaseError: For database operation failures
  - NetworkError: For network-related errors
  - ExternalServiceError: For external service call failures

## Relationships Between Packages

The error packages have the following relationships:

1. **core** is the foundation for all error types. It provides the BaseError type, error codes, and utility functions that are used by all other error packages.

2. **domain**, **app**, and **infra** all depend on **core** to extend the BaseError type with specific error types for different scenarios.

3. **domain**, **app**, and **infra** are independent of each other. They don't depend on each other, which allows for a clean separation of concerns.

4. Other packages in the codebase depend on these error packages to create and handle errors. For example, the **auth** package depends on the **errors** package to create authentication and authorization errors.

## Error Handling Approach

The error handling system uses the following approach:

1. **Error Creation**: Errors are created using constructor functions for specific error types. For example, `NewValidationError` creates a validation error, `NewDatabaseError` creates a database error, etc.

2. **Error Wrapping**: Errors can be wrapped with additional context using the `WithOperation` and `WithDetails` methods. This allows for adding operation names and additional details to errors.

3. **Error Type Checking**: Errors can be checked for specific types using type assertion or the `Is` and `As` functions from the standard errors package. Each error type also provides a type check method (e.g., `IsValidationError`, `IsDatabaseError`, etc.).

4. **HTTP Status Mapping**: Error codes are mapped to HTTP status codes for API error handling. This mapping is defined in the core package and can be accessed using the `GetHTTPStatus` method.

5. **Error Serialization**: Errors can be serialized to JSON for API responses using the `MarshalJSON` method.

## Best Practices

Here are some best practices for using the error handling system:

1. **Use the Appropriate Error Type**: Choose the most specific error type for your situation. For example, use `ValidationError` for validation failures, `DatabaseError` for database operation failures, etc.

2. **Add Context to Errors**: Wrap errors with additional context using the `WithOperation` and `WithDetails` methods. This helps with debugging and provides more information about the error.

3. **Check Error Types**: Use the type check methods (e.g., `IsValidationError`, `IsDatabaseError`, etc.) to handle different error types appropriately.

4. **Map Errors to HTTP Status Codes**: Use the `GetHTTPStatus` method to map errors to HTTP status codes for API error handling.

5. **Provide Detailed Error Messages**: Include detailed error messages that help identify the issue. For example, instead of "Invalid input", use "Email must be a valid email address".

6. **Use Standard Error Codes**: Use the standard error codes defined in the core package (e.g., `NotFoundCode`, `ValidationErrorCode`, etc.) for consistency.

7. **Include Source Information**: The error handling system automatically includes source file and line information, which helps with debugging.

## Error Detection and Handling

The error handling system provides comprehensive error detection and handling throughout the codebase:

1. **Error Detection**: Errors are detected at the appropriate level (domain, application, infrastructure) and wrapped with additional context as they propagate up the call stack.

2. **Error Handling**: Errors are handled at the appropriate level, with specific handling for different error types. For example, validation errors might be handled by returning a 400 Bad Request response, while database errors might be logged and retried.

3. **Retries and Timeouts**: The error handling system supports retries and timeouts for operations that might fail temporarily. This is especially important for infrastructure-related operations like database queries and network requests.

4. **Logging and Metrics**: Errors are logged with all contextual information and recorded as metrics for monitoring and alerting.

5. **Tracing**: Error information is added to traces for distributed tracing systems like OpenTelemetry.

## Conclusion

The error handling system in ServiceLib provides a comprehensive approach to error handling, with a clear hierarchy of error types, consistent error creation and handling, and integration with logging, metrics, and tracing. By following the best practices outlined in this document, you can ensure that errors are handled appropriately throughout your application.