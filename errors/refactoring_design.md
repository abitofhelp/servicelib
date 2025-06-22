# Error Handling System Redesign

## Overview

This document outlines the design for a complete redesign of the error handling system in the servicelib package. The redesign addresses the issues identified in the refactoring recommendations, including multiple overlapping error types, inconsistent wrap function behavior, and documentation gaps.

## Design Principles

1. **Simplicity**: The error handling system should be simple to use and understand.
2. **Consistency**: Error creation, wrapping, and handling should be consistent throughout the codebase.
3. **Extensibility**: The system should be extensible to accommodate domain-specific error types.
4. **Integration**: The system should integrate well with logging, metrics, and tracing.
5. **Performance**: The system should be efficient and not introduce significant overhead.

## Error Type Hierarchy

The new error handling system will have a clear hierarchy of error types:

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

## Error Structure

All errors will have a consistent structure with the following fields:

1. **Code**: A unique error code for categorizing errors.
2. **Message**: A human-readable error message.
3. **Operation**: The operation that failed (optional).
4. **Details**: Additional context information (optional).
5. **Cause**: The underlying error that caused this error (optional).
6. **Stack**: A stack trace for debugging (optional).

## Error Creation

Error creation will be standardized with a set of factory functions:

```
// Create a new error
New(code ErrorCode, message string) error

// Create a domain error
NewDomainError(code ErrorCode, message string) error
NewValidationError(code ErrorCode, message string, field string) error
NewBusinessRuleError(code ErrorCode, message string) error
NewNotFoundError(code ErrorCode, message string, resourceType, resourceID string) error

// Create an infrastructure error
NewInfrastructureError(code ErrorCode, message string) error
NewDatabaseError(code ErrorCode, message string) error
NewNetworkError(code ErrorCode, message string) error
NewExternalServiceError(code ErrorCode, message string, serviceName string) error

// Create an application error
NewApplicationError(code ErrorCode, message string) error
NewConfigurationError(code ErrorCode, message string) error
NewAuthenticationError(code ErrorCode, message string) error
NewAuthorizationError(code ErrorCode, message string) error
```

## Error Wrapping

Error wrapping will be standardized with a single approach:

```
// Wrap an error with additional context
Wrap(err error, code ErrorCode, message string) error

// Wrap an error with an operation
WrapWithOperation(err error, code ErrorCode, message string, operation string) error

// Wrap an error with details
WrapWithDetails(err error, code ErrorCode, message string, details map[string]interface{}) error
```

## Error Handling

Error handling utilities will include:

```
// Check if an error is of a specific type
Is(err, target error) bool
As(err error, target interface{}) bool

// Get information from an error
GetCode(err error) ErrorCode
GetMessage(err error) string
GetOperation(err error) string
GetDetails(err error) map[string]interface{}
GetCause(err error) error
GetStack(err error) string

// Convert an error to a specific format
ToJSON(err error) string
ToMap(err error) map[string]interface{}
```

## HTTP Status Mapping

Errors will be mapped to HTTP status codes based on their type and code:

```
// Map an error to an HTTP status code
GetHTTPStatus(err error) int
```

## Integration with Logging, Metrics, and Tracing

The error handling system will integrate with logging, metrics, and tracing:

```
// Log an error
LogError(ctx context.Context, err error)

// Record an error metric
RecordErrorMetric(err error)

// Add error information to a trace
AddErrorToTrace(ctx context.Context, err error)
```

## Package Structure

The new error handling system will be organized into the following packages:

```
errors/
├── core/       # Core error types and interfaces
├── domain/     # Domain-specific error types
├── infra/      # Infrastructure-related error types
├── app/        # Application-level error types
├── http/       # HTTP-related error utilities
├── log/        # Logging integration
├── metrics/    # Metrics integration
├── trace/      # Tracing integration
└── utils/      # Utility functions
```

## Migration Strategy

Since there is no requirement for backward compatibility, we will implement the new error handling system in a single step:

1. Create the new error handling system
2. Update all code to use the new system
3. Remove the old error handling system

## Documentation

The new error handling system will be thoroughly documented:

1. Package-level documentation explaining the error handling approach
2. Function-level documentation with examples
3. Usage examples for common error handling scenarios
4. Guidelines for error handling best practices
