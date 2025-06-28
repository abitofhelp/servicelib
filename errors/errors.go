// Copyright (c) 2024 A Bit of Help, Inc.

// Package errors provides a comprehensive error handling system for the application.
// It includes error codes, HTTP status mapping, contextual information, and utilities
// for creating, wrapping, and serializing errors.
package errors

import (
	"context"
	"strings"

	"github.com/abitofhelp/servicelib/errors/app"
	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/abitofhelp/servicelib/errors/domain"
	"github.com/abitofhelp/servicelib/errors/infra"
)

// ErrorCode is an alias for core.ErrorCode that represents a categorized error type.
// It is used to classify errors and map them to appropriate HTTP status codes.
type ErrorCode = core.ErrorCode

// Error code constants define the standard error categories used throughout the application.
const (
	// NotFoundCode indicates that a requested resource could not be found.
	NotFoundCode = core.NotFoundCode

	// InvalidInputCode indicates that the provided input is invalid.
	InvalidInputCode = core.InvalidInputCode

	// DatabaseErrorCode indicates an error occurred during a database operation.
	DatabaseErrorCode = core.DatabaseErrorCode

	// InternalErrorCode indicates an unexpected internal error.
	InternalErrorCode = core.InternalErrorCode

	// TimeoutCode indicates that an operation timed out.
	TimeoutCode = core.TimeoutCode

	// CanceledCode indicates that an operation was canceled.
	CanceledCode = core.CanceledCode

	// AlreadyExistsCode indicates that a resource already exists.
	AlreadyExistsCode = core.AlreadyExistsCode

	// UnauthorizedCode indicates that authentication is required but was not provided.
	UnauthorizedCode = core.UnauthorizedCode

	// ForbiddenCode indicates that the authenticated user does not have permission.
	ForbiddenCode = core.ForbiddenCode

	// ValidationErrorCode indicates that input validation failed.
	ValidationErrorCode = core.ValidationErrorCode

	// BusinessRuleViolationCode indicates that a business rule was violated.
	BusinessRuleViolationCode = core.BusinessRuleViolationCode

	// ExternalServiceErrorCode indicates an error occurred in an external service.
	ExternalServiceErrorCode = core.ExternalServiceErrorCode

	// NetworkErrorCode indicates a network-related error.
	NetworkErrorCode = core.NetworkErrorCode

	// ConfigurationErrorCode indicates an error in the application configuration.
	ConfigurationErrorCode = core.ConfigurationErrorCode

	// ResourceExhaustedCode indicates that a resource limit has been reached.
	ResourceExhaustedCode = core.ResourceExhaustedCode

	// DataCorruptionCode indicates that data is corrupted or invalid.
	DataCorruptionCode = core.DataCorruptionCode

	// ConcurrencyErrorCode indicates a concurrency-related error.
	ConcurrencyErrorCode = core.ConcurrencyErrorCode
)

// Standard errors provides commonly used error instances that can be used directly or wrapped.
var (
	// ErrNotFound represents a generic "resource not found" error.
	ErrNotFound = core.NewBaseError(NotFoundCode, "resource not found", nil)

	// ErrInvalidInput represents a generic "invalid input" error.
	ErrInvalidInput = core.NewBaseError(InvalidInputCode, "invalid input", nil)

	// ErrInternal represents a generic "internal error" for unexpected system issues.
	ErrInternal = core.NewBaseError(InternalErrorCode, "internal error", nil)

	// ErrUnauthorized represents a generic "unauthorized" error for authentication failures.
	ErrUnauthorized = core.NewBaseError(UnauthorizedCode, "unauthorized", nil)

	// ErrForbidden represents a generic "forbidden" error for authorization failures.
	ErrForbidden = core.NewBaseError(ForbiddenCode, "forbidden", nil)

	// ErrTimeout represents a generic "operation timed out" error.
	ErrTimeout = core.NewBaseError(TimeoutCode, "operation timed out", nil)

	// ErrCanceled represents a generic "operation canceled" error.
	ErrCanceled = core.NewBaseError(CanceledCode, "operation canceled", nil)

	// ErrAlreadyExists represents a generic "resource already exists" error.
	ErrAlreadyExists = core.NewBaseError(AlreadyExistsCode, "resource already exists", nil)
)

// Error type aliases provide convenient access to error types from the core, domain, infra, and app packages.
type (
	// BaseError is the foundation error type that all other error types extend.
	// It provides common functionality like error codes, messages, and cause tracking.
	BaseError = core.BaseError

	// DomainError represents errors that occur due to domain logic violations.
	// Use this for errors related to business rules and domain constraints.
	DomainError = domain.DomainError

	// ValidationError represents an error that occurs when input validation fails.
	// It includes the field name that failed validation.
	ValidationError = domain.ValidationError

	// ValidationErrors represents a collection of ValidationError instances.
	// Use this when multiple validation errors occur simultaneously.
	ValidationErrors = domain.ValidationErrors

	// BusinessRuleError represents an error that occurs when a business rule is violated.
	// It includes the name of the rule that was violated.
	BusinessRuleError = domain.BusinessRuleError

	// NotFoundError represents an error that occurs when a requested resource cannot be found.
	// It includes the resource type and ID that was not found.
	NotFoundError = domain.NotFoundError

	// InfrastructureError represents errors that occur in the infrastructure layer.
	// Use this for errors related to external systems, databases, networks, etc.
	InfrastructureError = infra.InfrastructureError

	// DatabaseError represents an error that occurs during a database operation.
	// It includes the operation and table name where the error occurred.
	DatabaseError = infra.DatabaseError

	// NetworkError represents an error that occurs during network communication.
	// It includes the host and port where the error occurred.
	NetworkError = infra.NetworkError

	// ExternalServiceError represents an error that occurs when calling an external service.
	// It includes the service name and endpoint where the error occurred.
	ExternalServiceError = infra.ExternalServiceError

	// RetryError represents an error that occurs when all retry attempts have been exhausted.
	// It includes the number of attempts made and the maximum allowed attempts.
	RetryError = infra.RetryError

	// ContextError represents an error that occurs due to context cancellation or timeout.
	// Use this when an operation fails because its context was canceled or timed out.
	ContextError = infra.ContextError

	// ApplicationError represents errors that occur in the application layer.
	// Use this for errors related to application logic and configuration.
	ApplicationError = app.ApplicationError

	// ConfigurationError represents an error that occurs due to invalid configuration.
	// It includes the configuration key and value that caused the error.
	ConfigurationError = app.ConfigurationError

	// AuthenticationError represents an error that occurs during authentication.
	// It includes the username that failed authentication.
	AuthenticationError = app.AuthenticationError

	// AuthorizationError represents an error that occurs during authorization.
	// It includes the username, resource, and action that failed authorization.
	AuthorizationError = app.AuthorizationError
)

// New creates a new BaseError with the specified error code and message.
//
// Parameters:
//   - code: The error code that categorizes this error.
//   - message: A human-readable description of the error.
//
// Returns:
//   - An error instance with the specified code and message.
func New(code ErrorCode, message string) error {
	return core.NewBaseError(code, message, nil)
}

// Wrap wraps an existing error with additional context, including an error code and message.
// If the input error is nil, it returns nil.
//
// Parameters:
//   - err: The original error to wrap.
//   - code: The error code to assign to the wrapped error.
//   - message: A human-readable description of the error context.
//
// Returns:
//   - A new error that wraps the original error with additional context.
func Wrap(err error, code ErrorCode, message string) error {
	if err == nil {
		return nil
	}
	return core.NewBaseError(code, message, err)
}

// WrapWithOperation wraps an error with additional context, including an operation name.
// This is useful for providing context about what operation was being performed when the error occurred.
// If the input error is nil, it returns nil.
//
// Parameters:
//   - err: The original error to wrap.
//   - code: The error code to assign to the wrapped error.
//   - message: A human-readable description of the error context.
//   - operation: The name of the operation that was being performed.
//
// Returns:
//   - A new error that wraps the original error with additional context.
func WrapWithOperation(err error, code ErrorCode, message string, operation string) error {
	if err == nil {
		return nil
	}
	return core.NewBaseError(code, message, err).WithOperation(operation)
}

// WrapWithDetails wraps an error with additional context, including a map of details.
// This is useful for providing structured data about the error context.
// If the input error is nil, it returns nil.
//
// Parameters:
//   - err: The original error to wrap.
//   - code: The error code to assign to the wrapped error.
//   - message: A human-readable description of the error context.
//   - details: A map of key-value pairs providing additional context.
//
// Returns:
//   - A new error that wraps the original error with additional context.
func WrapWithDetails(err error, code ErrorCode, message string, details map[string]interface{}) error {
	if err == nil {
		return nil
	}
	return core.NewBaseError(code, message, err).WithDetails(details)
}

// Domain error creation functions

// NewDomainError creates a new DomainError with the specified error code, message, and cause.
// Use this for errors that occur due to domain logic violations.
//
// Parameters:
//   - code: The error code that categorizes this error.
//   - message: A human-readable description of the error.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new DomainError instance.
func NewDomainError(code ErrorCode, message string, cause error) *DomainError {
	return domain.NewDomainError(code, message, cause)
}

// NewValidationError creates a new ValidationError for a specific field.
// Use this when input validation fails for a specific field.
//
// Parameters:
//   - message: A human-readable description of the validation error.
//   - field: The name of the field that failed validation.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new ValidationError instance.
func NewValidationError(message string, field string, cause error) *ValidationError {
	return domain.NewValidationError(message, field, cause)
}

// NewValidationErrors creates a new ValidationErrors collection with the specified errors.
// Use this when multiple validation errors occur simultaneously.
//
// Parameters:
//   - message: A human-readable description of the validation errors.
//   - errors: A list of ValidationError instances to include in the collection.
//
// Returns:
//   - A new ValidationErrors instance containing all the specified validation errors.
func NewValidationErrors(message string, errors ...*ValidationError) *ValidationErrors {
	return domain.NewValidationErrors(message, errors...)
}

// NewBusinessRuleError creates a new BusinessRuleError for a specific business rule.
// Use this when a business rule is violated.
//
// Parameters:
//   - message: A human-readable description of the business rule violation.
//   - rule: The name of the business rule that was violated.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new BusinessRuleError instance.
func NewBusinessRuleError(message string, rule string, cause error) *BusinessRuleError {
	return domain.NewBusinessRuleError(message, rule, cause)
}

// NewNotFoundError creates a new NotFoundError for a specific resource.
// Use this when a requested resource cannot be found.
//
// Parameters:
//   - resourceType: The type of resource that was not found (e.g., "user", "product").
//   - resourceID: The identifier of the resource that was not found.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new NotFoundError instance.
func NewNotFoundError(resourceType string, resourceID string, cause error) *NotFoundError {
	return domain.NewNotFoundError(resourceType, resourceID, cause)
}

// Infrastructure error creation functions

// NewInfrastructureError creates a new InfrastructureError with the specified error code, message, and cause.
// Use this for errors that occur in the infrastructure layer.
//
// Parameters:
//   - code: The error code that categorizes this error.
//   - message: A human-readable description of the error.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new InfrastructureError instance.
func NewInfrastructureError(code ErrorCode, message string, cause error) *InfrastructureError {
	return infra.NewInfrastructureError(code, message, cause)
}

// NewDatabaseError creates a new DatabaseError for a specific database operation.
// Use this when a database operation fails.
//
// Parameters:
//   - message: A human-readable description of the database error.
//   - operation: The database operation that failed (e.g., "SELECT", "INSERT").
//   - table: The name of the database table involved in the operation.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new DatabaseError instance.
func NewDatabaseError(message string, operation string, table string, cause error) *DatabaseError {
	return infra.NewDatabaseError(message, operation, table, cause)
}

// NewNetworkError creates a new NetworkError for a specific network connection.
// Use this when a network operation fails.
//
// Parameters:
//   - message: A human-readable description of the network error.
//   - host: The hostname or IP address of the remote host.
//   - port: The port number of the remote host.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new NetworkError instance.
func NewNetworkError(message string, host string, port string, cause error) *NetworkError {
	return infra.NewNetworkError(message, host, port, cause)
}

// NewExternalServiceError creates a new ExternalServiceError for a specific external service call.
// Use this when a call to an external service fails.
//
// Parameters:
//   - message: A human-readable description of the external service error.
//   - serviceName: The name of the external service.
//   - endpoint: The specific endpoint or API that was called.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new ExternalServiceError instance.
func NewExternalServiceError(message string, serviceName string, endpoint string, cause error) *ExternalServiceError {
	return infra.NewExternalServiceError(message, serviceName, endpoint, cause)
}

// NewRetryError creates a new RetryError when all retry attempts have been exhausted.
// Use this when an operation fails after multiple retry attempts.
//
// Parameters:
//   - message: A human-readable description of the retry error.
//   - cause: The underlying error that caused the retries to fail.
//   - attempts: The number of attempts that were made.
//   - maxAttempts: The maximum number of attempts that were allowed.
//
// Returns:
//   - A new RetryError instance.
func NewRetryError(message string, cause error, attempts int, maxAttempts int) *RetryError {
	return infra.NewRetryError(message, cause, attempts, maxAttempts)
}

// NewContextError creates a new ContextError when an operation fails due to context cancellation or timeout.
// Use this when an operation fails because its context was canceled or timed out.
//
// Parameters:
//   - message: A human-readable description of the context error.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new ContextError instance.
func NewContextError(message string, cause error) *ContextError {
	return infra.NewContextError(message, cause)
}

// Application error creation functions

// NewApplicationError creates a new ApplicationError with the specified error code, message, and cause.
// Use this for errors that occur in the application layer.
//
// Parameters:
//   - code: The error code that categorizes this error.
//   - message: A human-readable description of the error.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new ApplicationError instance.
func NewApplicationError(code ErrorCode, message string, cause error) *ApplicationError {
	return app.NewApplicationError(code, message, cause)
}

// NewConfigurationError creates a new ConfigurationError for a specific configuration key.
// Use this when an error occurs due to invalid or missing configuration.
//
// Parameters:
//   - message: A human-readable description of the configuration error.
//   - configKey: The configuration key that caused the error.
//   - configValue: The value of the configuration key that caused the error.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new ConfigurationError instance.
func NewConfigurationError(message string, configKey string, configValue string, cause error) *ConfigurationError {
	return app.NewConfigurationError(message, configKey, configValue, cause)
}

// NewAuthenticationError creates a new AuthenticationError for a specific user.
// Use this when authentication fails for a user.
//
// Parameters:
//   - message: A human-readable description of the authentication error.
//   - username: The username of the user who failed authentication.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new AuthenticationError instance.
func NewAuthenticationError(message string, username string, cause error) *AuthenticationError {
	return app.NewAuthenticationError(message, username, cause)
}

// NewAuthorizationError creates a new AuthorizationError for a specific user, resource, and action.
// Use this when a user is not authorized to perform an action on a resource.
//
// Parameters:
//   - message: A human-readable description of the authorization error.
//   - username: The username of the user who is not authorized.
//   - resource: The resource the user is trying to access.
//   - action: The action the user is trying to perform.
//   - cause: The underlying error that caused this error, if any.
//
// Returns:
//   - A new AuthorizationError instance.
func NewAuthorizationError(message string, username string, resource string, action string, cause error) *AuthorizationError {
	return app.NewAuthorizationError(message, username, resource, action, cause)
}

// Error checking functions

// Is reports whether any error in err's chain matches target.
// This is a wrapper around the standard errors.Is function.
//
// Parameters:
//   - err: The error to check.
//   - target: The target error to compare against.
//
// Returns:
//   - true if err or any error in its chain matches target, false otherwise.
func Is(err, target error) bool {
	return core.Is(err, target)
}

// As finds the first error in err's chain that matches the type of target,
// and if so, sets target to that error value and returns true.
// This is a wrapper around the standard errors.As function.
//
// Parameters:
//   - err: The error to check.
//   - target: A pointer to a variable of the error type to check for.
//
// Returns:
//   - true if an error of the target type was found in err's chain, false otherwise.
func As(err error, target interface{}) bool {
	return core.As(err, target)
}

// Unwrap returns the underlying error of err, if any.
// This is a wrapper around the standard errors.Unwrap function.
//
// Parameters:
//   - err: The error to unwrap.
//
// Returns:
//   - The underlying error of err, or nil if err has no underlying error.
func Unwrap(err error) error {
	return core.Unwrap(err)
}

// GetHTTPStatus returns the HTTP status code for an error.
// This function maps error codes to appropriate HTTP status codes.
// If the error is nil, it returns 0.
//
// Parameters:
//   - err: The error to get the HTTP status code for.
//
// Returns:
//   - The HTTP status code corresponding to the error, or 0 if the error is nil or has no mapping.
func GetHTTPStatus(err error) int {
	if err == nil {
		return 0
	}

	// Check if the error has a GetHTTPStatus method
	if e, ok := err.(interface{ GetHTTPStatus() int }); ok {
		return e.GetHTTPStatus()
	}

	// Check if the error has a Code method
	if e, ok := err.(interface{ GetCode() ErrorCode }); ok {
		return core.GetHTTPStatus(e.GetCode())
	}

	return 0
}

// ToJSON converts an error to a JSON string representation.
// This is useful for logging errors or returning them in API responses.
//
// Parameters:
//   - err: The error to convert to JSON.
//
// Returns:
//   - A JSON string representation of the error.
func ToJSON(err error) string {
	return core.ToJSON(err)
}

// Error type checking functions

// IsDomainError checks if an error is a domain error.
// Domain errors include DomainError, ValidationError, BusinessRuleError, and NotFoundError.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is a domain error, false otherwise.
func IsDomainError(err error) bool {
	// Check if the error is a DomainError
	var domainErr *DomainError
	if As(err, &domainErr) {
		return true
	}

	// Check if the error is a ValidationError
	var validationErr *ValidationError
	if As(err, &validationErr) {
		return true
	}

	// Check if the error is a BusinessRuleError
	var businessRuleErr *BusinessRuleError
	if As(err, &businessRuleErr) {
		return true
	}

	// Check if the error is a NotFoundError
	var notFoundErr *NotFoundError
	if As(err, &notFoundErr) {
		return true
	}

	return false
}

// IsValidationError checks if an error is a ValidationError.
// Use this to determine if an error is related to input validation.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is a ValidationError, false otherwise.
func IsValidationError(err error) bool {
	var e *ValidationError
	return As(err, &e)
}

// IsBusinessRuleError checks if an error is a BusinessRuleError.
// Use this to determine if an error is related to a business rule violation.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is a BusinessRuleError, false otherwise.
func IsBusinessRuleError(err error) bool {
	var e *BusinessRuleError
	return As(err, &e)
}

// IsNotFoundError checks if an error is a NotFoundError.
// Use this to determine if an error is related to a resource not being found.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is a NotFoundError, false otherwise.
func IsNotFoundError(err error) bool {
	var e *NotFoundError
	return As(err, &e)
}

// IsInfrastructureError checks if an error is an infrastructure error.
// Infrastructure errors include InfrastructureError, DatabaseError, NetworkError, and ExternalServiceError.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is an infrastructure error, false otherwise.
func IsInfrastructureError(err error) bool {
	// Check if the error is an InfrastructureError
	var infraErr *InfrastructureError
	if As(err, &infraErr) {
		return true
	}

	// Check if the error is a DatabaseError
	var dbErr *DatabaseError
	if As(err, &dbErr) {
		return true
	}

	// Check if the error is a NetworkError
	var networkErr *NetworkError
	if As(err, &networkErr) {
		return true
	}

	// Check if the error is an ExternalServiceError
	var externalErr *ExternalServiceError
	if As(err, &externalErr) {
		return true
	}

	return false
}

// IsDatabaseError checks if an error is a DatabaseError.
// Use this to determine if an error is related to a database operation.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is a DatabaseError, false otherwise.
func IsDatabaseError(err error) bool {
	var e *DatabaseError
	return As(err, &e)
}

// IsNetworkError checks if an error is a network-related error.
// This function checks for NetworkError types from this package,
// as well as errors that implement the net.Error interface,
// timeout errors, and errors with common network error messages.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is a network-related error, false otherwise.
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Check if the error is a known network error type from our errors package
	var networkErr *NetworkError
	if As(err, &networkErr) {
		return true
	}

	// Check if the error implements net.Error interface with Timeout() or Temporary() methods
	if netErr, ok := err.(interface {
		Timeout() bool
		Temporary() bool
	}); ok {
		return netErr.Temporary() || netErr.Timeout()
	}

	// Check if the error is a timeout error
	if IsTimeout(err) {
		return true
	}

	// Fall back to string matching for errors that don't implement specific interfaces
	// This is less reliable but provides backward compatibility
	errMsg := err.Error()
	networkErrorStrings := []string{
		"connection refused",
		"connection reset",
		"connection closed",
		"no route to host",
		"network is unreachable",
		"broken pipe",
	}

	for _, s := range networkErrorStrings {
		if strings.Contains(strings.ToLower(errMsg), s) {
			return true
		}
	}

	return false
}

// IsExternalServiceError checks if an error is an ExternalServiceError.
// Use this to determine if an error is related to an external service call.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is an ExternalServiceError, false otherwise.
func IsExternalServiceError(err error) bool {
	var e *ExternalServiceError
	return As(err, &e)
}

// IsRetryError checks if an error is a RetryError.
// Use this to determine if an error occurred after exhausting all retry attempts.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is a RetryError, false otherwise.
func IsRetryError(err error) bool {
	var e *RetryError
	return As(err, &e)
}

// IsContextError checks if an error is a ContextError.
// Use this to determine if an error is related to context cancellation or timeout.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is a ContextError, false otherwise.
func IsContextError(err error) bool {
	var e *ContextError
	return As(err, &e)
}

// IsConfigurationError checks if an error is a ConfigurationError.
// Use this to determine if an error is related to invalid or missing configuration.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is a ConfigurationError, false otherwise.
func IsConfigurationError(err error) bool {
	var e *ConfigurationError
	return As(err, &e)
}

// IsAuthenticationError checks if an error is an AuthenticationError.
// Use this to determine if an error is related to authentication failure.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is an AuthenticationError, false otherwise.
func IsAuthenticationError(err error) bool {
	var e *AuthenticationError
	return As(err, &e)
}

// IsAuthorizationError checks if an error is an AuthorizationError.
// Use this to determine if an error is related to authorization failure.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is an AuthorizationError, false otherwise.
func IsAuthorizationError(err error) bool {
	var e *AuthorizationError
	return As(err, &e)
}

// IsApplicationError checks if an error is an application error.
// Application errors include ApplicationError, ConfigurationError, AuthenticationError, and AuthorizationError.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is an application error, false otherwise.
func IsApplicationError(err error) bool {
	// Check if the error is an ApplicationError
	var appErr *ApplicationError
	if As(err, &appErr) {
		return true
	}

	// Check if the error is a ConfigurationError
	var configErr *ConfigurationError
	if As(err, &configErr) {
		return true
	}

	// Check if the error is an AuthenticationError
	var authErr *AuthenticationError
	if As(err, &authErr) {
		return true
	}

	// Check if the error is an AuthorizationError
	var authzErr *AuthorizationError
	if As(err, &authzErr) {
		return true
	}

	return false
}

// IsCancelled checks if an error is a cancellation error.
// This function checks for errors that indicate an operation was canceled.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error indicates cancellation, false otherwise.
func IsCancelled(err error) bool {
	if err == nil {
		return false
	}

	// Check if the error is ErrContextCanceled from recovery package
	if Is(err, ErrCanceled) {
		return true
	}

	// Check if the error has a GetCode method and the code is CanceledCode
	if e, ok := err.(interface{ GetCode() ErrorCode }); ok {
		return e.GetCode() == CanceledCode
	}

	return false
}

// IsTimeout checks if an error is a timeout error.
// This function checks for errors that indicate an operation timed out.
// It checks for timeout errors from this package, context deadline exceeded errors,
// errors that implement the net.Error interface with a Timeout method,
// and errors with common timeout error messages.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error indicates a timeout, false otherwise.
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}

	// Check if the error is ErrTimeout
	if Is(err, ErrTimeout) {
		return true
	}

	// Check if the error has a GetCode method and the code is TimeoutCode
	if e, ok := err.(interface{ GetCode() ErrorCode }); ok {
		return e.GetCode() == TimeoutCode
	}

	// Check if the error is a context deadline exceeded error
	if err == context.DeadlineExceeded {
		return true
	}

	// Check if the error implements net.Error interface with Timeout() method
	if netErr, ok := err.(interface {
		Timeout() bool
	}); ok {
		return netErr.Timeout()
	}

	// Fall back to string matching for errors that don't implement specific interfaces
	// This is less reliable but provides backward compatibility
	errMsg := err.Error()
	timeoutErrorStrings := []string{
		"timeout",
		"timed out",
		"deadline exceeded",
	}

	for _, s := range timeoutErrorStrings {
		if strings.Contains(strings.ToLower(errMsg), s) {
			return true
		}
	}

	return false
}

// IsTransientError checks if an error is a transient error that may be resolved by retrying.
// Transient errors include network errors, timeout errors, and errors that explicitly
// indicate they are temporary. This function is useful for determining if an operation
// should be retried.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is likely transient and may be resolved by retrying, false otherwise.
func IsTransientError(err error) bool {
	if err == nil {
		return false
	}

	// Check if the error is a known transient error type from our errors package
	// such as RetryError, NetworkError, or ExternalServiceError
	if IsRetryError(err) || IsNetworkError(err) || IsExternalServiceError(err) {
		return true
	}

	// Check if the error is a network error or timeout error
	if IsNetworkError(err) || IsTimeout(err) {
		return true
	}

	// Check if the error implements a Temporary() method (common in Go standard library)
	if tempErr, ok := err.(interface {
		Temporary() bool
	}); ok {
		return tempErr.Temporary()
	}

	// Fall back to string matching for errors that don't implement specific interfaces
	// This is less reliable but provides backward compatibility
	errMsg := err.Error()
	transientErrorStrings := []string{
		"server is busy",
		"too many requests",
		"rate limit exceeded",
		"try again",
		"temporary",
		"transient",
	}

	for _, s := range transientErrorStrings {
		if strings.Contains(strings.ToLower(errMsg), s) {
			return true
		}
	}

	return false
}
