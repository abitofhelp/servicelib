// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides a comprehensive error handling system for the application.
// It includes error codes, HTTP status mapping, contextual information, and utilities
// for creating, wrapping, and serializing errors.
package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrorCode represents a unique error code for categorizing errors.
// These codes are used for error identification, logging, and mapping to HTTP status codes.
type ErrorCode string

// Standard error codes define all possible error categories in the application.
const (
	// NotFoundCode is used when a resource is not found.
	// Maps to HTTP 404 Not Found.
	NotFoundCode ErrorCode = "NOT_FOUND"

	// InvalidInputCode is used when input validation fails.
	// Maps to HTTP 400 Bad Request.
	InvalidInputCode ErrorCode = "INVALID_INPUT"

	// DatabaseErrorCode is used for database operation failures.
	// Maps to HTTP 500 Internal Server Error.
	DatabaseErrorCode ErrorCode = "DATABASE_ERROR"

	// InternalErrorCode is used for internal server errors.
	// Maps to HTTP 500 Internal Server Error.
	InternalErrorCode ErrorCode = "INTERNAL_ERROR"

	// TimeoutCode is used when an operation times out.
	// Maps to HTTP 504 Gateway Timeout.
	TimeoutCode ErrorCode = "TIMEOUT"

	// CanceledCode is used when an operation is canceled.
	// Maps to HTTP 408 Request Timeout.
	CanceledCode ErrorCode = "CANCELED"

	// AlreadyExistsCode is used when a resource already exists.
	// Maps to HTTP 409 Conflict.
	AlreadyExistsCode ErrorCode = "ALREADY_EXISTS"

	// UnauthorizedCode is used for authentication failures.
	// Maps to HTTP 401 Unauthorized.
	UnauthorizedCode ErrorCode = "UNAUTHORIZED"

	// ForbiddenCode is used for authorization failures.
	// Maps to HTTP 403 Forbidden.
	ForbiddenCode ErrorCode = "FORBIDDEN"

	// ValidationErrorCode is used for domain validation errors.
	// Maps to HTTP 400 Bad Request.
	ValidationErrorCode ErrorCode = "VALIDATION_ERROR"

	// BusinessRuleViolationCode is used when a business rule is violated.
	// Maps to HTTP 422 Unprocessable Entity.
	BusinessRuleViolationCode ErrorCode = "BUSINESS_RULE_VIOLATION"

	// ExternalServiceErrorCode is used when an external service call fails.
	// Maps to HTTP 502 Bad Gateway.
	ExternalServiceErrorCode ErrorCode = "EXTERNAL_SERVICE_ERROR"

	// NetworkErrorCode is used for network-related errors.
	// Maps to HTTP 503 Service Unavailable.
	NetworkErrorCode ErrorCode = "NETWORK_ERROR"

	// ConfigurationErrorCode is used for configuration errors.
	// Maps to HTTP 500 Internal Server Error.
	ConfigurationErrorCode ErrorCode = "CONFIGURATION_ERROR"

	// ResourceExhaustedCode is used when a resource limit is reached.
	// Maps to HTTP 429 Too Many Requests.
	ResourceExhaustedCode ErrorCode = "RESOURCE_EXHAUSTED"

	// DataCorruptionCode is used when data is corrupted.
	// Maps to HTTP 500 Internal Server Error.
	DataCorruptionCode ErrorCode = "DATA_CORRUPTION"

	// ConcurrencyErrorCode is used for concurrency-related errors.
	// Maps to HTTP 409 Conflict.
	ConcurrencyErrorCode ErrorCode = "CONCURRENCY_ERROR"
)

// Standard errors that can be used throughout the application
var (
	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists is returned when a resource already exists
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrInvalidInput is returned when the input to a function is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrInternal is returned when an internal error occurs
	ErrInternal = errors.New("internal error")

	// ErrUnauthorized is returned when a user is not authorized to perform an action
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when a user is forbidden from performing an action
	ErrForbidden = errors.New("forbidden")

	// ErrTimeout is returned when an operation times out
	ErrTimeout = errors.New("operation timed out")

	// ErrCancelled is returned when an operation is cancelled
	ErrCancelled = errors.New("operation cancelled")

	// ErrConflict is returned when there is a conflict with the current state
	ErrConflict = errors.New("conflict with current state")
)

// ErrorContext holds additional context for an error.
// It includes information about the operation that failed, the source location,
// and any additional details that might be useful for debugging or error reporting.
type ErrorContext struct {
	// Operation is the name of the operation that failed
	Operation string `json:"operation,omitempty"`

	// Source is the file and line where the error occurred
	Source string `json:"source,omitempty"`

	// Line is the line number where the error occurred
	Line int `json:"line,omitempty"`

	// Code is the error code
	Code ErrorCode `json:"code,omitempty"`

	// HTTPStatus is the HTTP status code to return for this error
	HTTPStatus int `json:"http_status,omitempty"`

	// Details contains additional information about the error
	Details map[string]interface{} `json:"details,omitempty"`
}

// ContextualError is an error with additional context.
// It wraps another error and adds contextual information like operation name,
// error code, HTTP status, and source location.
type ContextualError struct {
	// Original is the original error that was wrapped
	Original error

	// Context contains additional information about the error
	Context ErrorContext
}

// Error returns the error message with contextual information.
func (e *ContextualError) Error() string {
	var builder strings.Builder

	// Add operation if available
	if e.Context.Operation != "" {
		builder.WriteString(fmt.Sprintf("operation %s: ", e.Context.Operation))
	}

	// Add original error message
	if e.Original != nil {
		builder.WriteString(e.Original.Error())
	} else {
		builder.WriteString("an error occurred")
	}

	// Add source location if available
	if e.Context.Source != "" {
		builder.WriteString(fmt.Sprintf(" (source: %s", e.Context.Source))
		if e.Context.Line > 0 {
			builder.WriteString(fmt.Sprintf(":%d", e.Context.Line))
		}
		builder.WriteString(")")
	}

	return builder.String()
}

// Unwrap returns the original error.
func (e *ContextualError) Unwrap() error {
	return e.Original
}

// Code returns the error code.
func (e *ContextualError) Code() ErrorCode {
	return e.Context.Code
}

// HTTPStatus returns the HTTP status code.
func (e *ContextualError) HTTPStatus() int {
	return e.Context.HTTPStatus
}

// MarshalJSON implements the json.Marshaler interface.
func (e *ContextualError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string       `json:"message"`
		Context ErrorContext `json:"context,omitempty"`
	}{
		Message: e.Error(),
		Context: e.Context,
	})
}

// getCallerInfo returns the file name and line number of the caller.
func getCallerInfo(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "", 0
	}
	return filepath.Base(file), line
}

// withContext wraps an error with contextual information.
// It adds operation name, error code, HTTP status, and source location to the error.
// If the error is already a ContextualError, it updates the context with the new information.
func withContext(err error, operation string, code ErrorCode, httpStatus int, details map[string]interface{}) error {
	// Get caller information
	source, line := getCallerInfo(2)

	// If the error is already a ContextualError, update its context
	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		// Only update operation if it's not already set
		if operation != "" && contextualErr.Context.Operation == "" {
			contextualErr.Context.Operation = operation
		}

		// Only update code if it's not already set
		if code != "" && contextualErr.Context.Code == "" {
			contextualErr.Context.Code = code
		}

		// Only update HTTP status if it's not already set
		if httpStatus != 0 && contextualErr.Context.HTTPStatus == 0 {
			contextualErr.Context.HTTPStatus = httpStatus
		}

		// Merge details if provided
		if details != nil {
			if contextualErr.Context.Details == nil {
				contextualErr.Context.Details = make(map[string]interface{})
			}
			for k, v := range details {
				contextualErr.Context.Details[k] = v
			}
		}

		return contextualErr
	}

	// Create a new ContextualError
	return &ContextualError{
		Original: err,
		Context: ErrorContext{
			Operation:  operation,
			Source:     source,
			Line:       line,
			Code:       code,
			HTTPStatus: httpStatus,
			Details:    details,
		},
	}
}

// AppError is a generic error type that can be used for different error categories
type AppError[T ~string] struct {
	Err     error
	Message string
	Code    string
	Type    T
}

func (e *AppError[T]) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AppError[T]) Unwrap() error {
	return e.Err
}

func (e *AppError[T]) ErrorType() T {
	return e.Type
}

// NewAppError creates a new AppError
func NewAppError[T ~string](err error, message, code string, errorType T) *AppError[T] {
	return &AppError[T]{
		Err:     err,
		Message: message,
		Code:    code,
		Type:    errorType,
	}
}

// Map of error codes to HTTP status codes
var errorCodeToHTTPStatus = map[ErrorCode]int{
	NotFoundCode:              http.StatusNotFound,
	InvalidInputCode:          http.StatusBadRequest,
	DatabaseErrorCode:         http.StatusInternalServerError,
	InternalErrorCode:         http.StatusInternalServerError,
	TimeoutCode:               http.StatusGatewayTimeout,
	CanceledCode:              http.StatusRequestTimeout,
	AlreadyExistsCode:         http.StatusConflict,
	UnauthorizedCode:          http.StatusUnauthorized,
	ForbiddenCode:             http.StatusForbidden,
	ValidationErrorCode:       http.StatusBadRequest,
	BusinessRuleViolationCode: http.StatusUnprocessableEntity,
	ExternalServiceErrorCode:  http.StatusBadGateway,
	NetworkErrorCode:          http.StatusServiceUnavailable,
	ConfigurationErrorCode:    http.StatusInternalServerError,
	ResourceExhaustedCode:     http.StatusTooManyRequests,
	DataCorruptionCode:        http.StatusInternalServerError,
	ConcurrencyErrorCode:      http.StatusConflict,
}

// NotFound creates a new error for when a resource is not found.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with NotFoundCode and HTTP 404 status
func NotFound(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", NotFoundCode, http.StatusNotFound, nil)
}

// InvalidInput creates a new error for when input validation fails.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with InvalidInputCode and HTTP 400 status
func InvalidInput(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", InvalidInputCode, http.StatusBadRequest, nil)
}

// DatabaseOperation creates a new error for database operation failures.
// Parameters:
//   - err: The original error
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with DatabaseErrorCode and HTTP 500 status
func DatabaseOperation(err error, format string, args ...interface{}) error {
	if err == nil {
		err = fmt.Errorf(format, args...)
	} else {
		err = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
	}
	return withContext(err, "database", DatabaseErrorCode, http.StatusInternalServerError, nil)
}

// Internal creates a new error for internal server errors.
// Parameters:
//   - err: The original error
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with InternalErrorCode and HTTP 500 status
func Internal(err error, format string, args ...interface{}) error {
	if err == nil {
		err = fmt.Errorf(format, args...)
	} else {
		err = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
	}
	return withContext(err, "", InternalErrorCode, http.StatusInternalServerError, nil)
}

// Timeout creates a new error for when an operation times out.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with TimeoutCode and HTTP 504 status
func Timeout(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", TimeoutCode, http.StatusGatewayTimeout, nil)
}

// Canceled creates a new error for when an operation is canceled.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with CanceledCode and HTTP 408 status
func Canceled(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", CanceledCode, http.StatusRequestTimeout, nil)
}

// AlreadyExists creates a new error for when a resource already exists.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with AlreadyExistsCode and HTTP 409 status
func AlreadyExists(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", AlreadyExistsCode, http.StatusConflict, nil)
}

// Unauthorized creates a new error for authentication failures.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with UnauthorizedCode and HTTP 401 status
func Unauthorized(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", UnauthorizedCode, http.StatusUnauthorized, nil)
}

// Forbidden creates a new error for authorization failures.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with ForbiddenCode and HTTP 403 status
func Forbidden(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", ForbiddenCode, http.StatusForbidden, nil)
}

// Validation creates a new error for domain validation errors.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with ValidationErrorCode and HTTP 400 status
func Validation(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", ValidationErrorCode, http.StatusBadRequest, nil)
}

// BusinessRuleViolation creates a new error for when a business rule is violated.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with BusinessRuleViolationCode and HTTP 422 status
func BusinessRuleViolation(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", BusinessRuleViolationCode, http.StatusUnprocessableEntity, nil)
}

// ExternalService creates a new error for external service call failures.
// Parameters:
//   - err: The original error
//   - service: The name of the external service
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with ExternalServiceErrorCode and HTTP 502 status
func ExternalService(err error, service string, format string, args ...interface{}) error {
	if err == nil {
		err = fmt.Errorf(format, args...)
	} else {
		err = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
	}
	return withContext(err, fmt.Sprintf("external_service_%s", service), ExternalServiceErrorCode, http.StatusBadGateway, nil)
}

// Network creates a new error for network-related errors.
// Parameters:
//   - err: The original error
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with NetworkErrorCode and HTTP 503 status
func Network(err error, format string, args ...interface{}) error {
	if err == nil {
		err = fmt.Errorf(format, args...)
	} else {
		err = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
	}
	return withContext(err, "network", NetworkErrorCode, http.StatusServiceUnavailable, nil)
}

// Configuration creates a new error for configuration errors.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with ConfigurationErrorCode and HTTP 500 status
func Configuration(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "configuration", ConfigurationErrorCode, http.StatusInternalServerError, nil)
}

// ResourceExhausted creates a new error for when a resource limit is reached.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with ResourceExhaustedCode and HTTP 429 status
func ResourceExhausted(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", ResourceExhaustedCode, http.StatusTooManyRequests, nil)
}

// DataCorruption creates a new error for when data is corrupted.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with DataCorruptionCode and HTTP 500 status
func DataCorruption(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", DataCorruptionCode, http.StatusInternalServerError, nil)
}

// Concurrency creates a new error for concurrency-related errors.
// Parameters:
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error with ConcurrencyErrorCode and HTTP 409 status
func Concurrency(format string, args ...interface{}) error {
	return withContext(fmt.Errorf(format, args...), "", ConcurrencyErrorCode, http.StatusConflict, nil)
}

// Is checks if an error matches a target error.
// This is a wrapper around errors.Is that adds support for ContextualError.
// Parameters:
//   - err: The error to check
//   - target: The target error to match against
//
// Returns:
//   - bool: True if the error matches the target
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target.
// This is a wrapper around errors.As that adds support for ContextualError.
// Parameters:
//   - err: The error to check
//   - target: A pointer to the error type to match against
//
// Returns:
//   - bool: True if a match was found
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap returns the underlying error.
// This is a wrapper around errors.Unwrap that adds support for ContextualError.
// Parameters:
//   - err: The error to unwrap
//
// Returns:
//   - error: The underlying error, or nil if there is none
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// WrapWithOperation wraps an error with an operation name and message.
// It preserves the error chain and adds source location information.
// Parameters:
//   - err: The error to wrap
//   - operation: The name of the operation that failed
//   - format: The error message format string
//   - args: Arguments for the format string
//
// Returns:
//   - error: A new error that wraps the original error
func WrapWithOperation(err error, operation string, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	// Get caller information
	source, line := getCallerInfo(1)

	// If the error is already a ContextualError, update its context
	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		// Create a new ContextualError with the updated message and operation
		return &ContextualError{
			Original: contextualErr.Original,
			Context: ErrorContext{
				Operation:  operation,
				Source:     source,
				Line:       line,
				Code:       contextualErr.Context.Code,
				HTTPStatus: contextualErr.Context.HTTPStatus,
				Details:    contextualErr.Context.Details,
			},
		}
	}

	// Create a new ContextualError
	return &ContextualError{
		Original: err,
		Context: ErrorContext{
			Operation: operation,
			Source:    source,
			Line:      line,
		},
	}
}

// WithDetails adds details to an error.
// It preserves the error chain and adds additional context information.
// Parameters:
//   - err: The error to enhance
//   - details: A map of additional details to add to the error
//
// Returns:
//   - error: A new error with the added details
func WithDetails(err error, details map[string]interface{}) error {
	if err == nil {
		return nil
	}

	// Get caller information
	source, line := getCallerInfo(1)

	// If the error is already a ContextualError, update its context
	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		// Merge the details
		newDetails := make(map[string]interface{})
		if contextualErr.Context.Details != nil {
			for k, v := range contextualErr.Context.Details {
				newDetails[k] = v
			}
		}
		for k, v := range details {
			newDetails[k] = v
		}

		// Create a new ContextualError with the merged details
		return &ContextualError{
			Original: contextualErr.Original,
			Context: ErrorContext{
				Operation:  contextualErr.Context.Operation,
				Source:     source,
				Line:       line,
				Code:       contextualErr.Context.Code,
				HTTPStatus: contextualErr.Context.HTTPStatus,
				Details:    newDetails,
			},
		}
	}

	// Create a new ContextualError
	return &ContextualError{
		Original: err,
		Context: ErrorContext{
			Source:  source,
			Line:    line,
			Details: details,
		},
	}
}

// GetCode returns the error code from an error.
// If the error is a ContextualError, it returns the code from the context.
// Otherwise, it returns an empty string.
// Parameters:
//   - err: The error to get the code from
//
// Returns:
//   - ErrorCode: The error code, or an empty string if not available
func GetCode(err error) ErrorCode {
	if err == nil {
		return ""
	}

	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		return contextualErr.Context.Code
	}

	return ""
}

// GetHTTPStatus returns the HTTP status code from an error.
// If the error is a ContextualError, it returns the HTTP status from the context.
// Otherwise, it returns 0.
// Parameters:
//   - err: The error to get the HTTP status from
//
// Returns:
//   - int: The HTTP status code, or 0 if not available
func GetHTTPStatus(err error) int {
	if err == nil {
		return 0
	}

	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		return contextualErr.Context.HTTPStatus
	}

	return 0
}

// ToJSON converts an error to a JSON string.
// If the error is a ContextualError, it uses the MarshalJSON method.
// Otherwise, it creates a simple JSON object with the error message.
// Parameters:
//   - err: The error to convert to JSON
//
// Returns:
//   - string: The JSON representation of the error
func ToJSON(err error) string {
	if err == nil {
		return "{}"
	}

	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		jsonBytes, jsonErr := json.Marshal(contextualErr)
		if jsonErr != nil {
			return fmt.Sprintf(`{"message":"Error marshaling error to JSON: %s"}`, jsonErr.Error())
		}
		return string(jsonBytes)
	}

	// For non-ContextualError errors, create a simple JSON object
	jsonBytes, jsonErr := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	})
	if jsonErr != nil {
		return fmt.Sprintf(`{"message":"Error marshaling error to JSON: %s"}`, jsonErr.Error())
	}
	return string(jsonBytes)
}

// DomainErrorType represents the type of domain error
type DomainErrorType string

// Domain error type constants
const (
	DomainErrorGeneral DomainErrorType = "DOMAIN_ERROR"
)

// DomainError represents an error that occurred in the domain layer
type DomainError = AppError[DomainErrorType]

// NewDomainError creates a new DomainError
func NewDomainError(err error, message, code string) *DomainError {
	return NewAppError(err, message, code, DomainErrorGeneral)
}

// ValidationError represents a validation error
type ValidationError struct {
	Msg   string
	Field string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error: %s (field: %s)", e.Msg, e.Field)
	}
	return fmt.Sprintf("validation error: %s", e.Msg)
}

// NewValidationError creates a new ValidationError
func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Msg: msg}
}

// NewFieldValidationError creates a new ValidationError with a field name
func NewFieldValidationError(msg, field string) *ValidationError {
	return &ValidationError{Msg: msg, Field: field}
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []*ValidationError
}

func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}

	msg := fmt.Sprintf("%d validation errors:", len(e.Errors))
	for i, err := range e.Errors {
		msg += fmt.Sprintf("\n  %d. %s", i+1, err.Error())
	}
	return msg
}

// NewValidationErrors creates a new ValidationErrors
func NewValidationErrors(errors ...*ValidationError) *ValidationErrors {
	return &ValidationErrors{Errors: errors}
}

// AddError adds a validation error to the collection
func (e *ValidationErrors) AddError(err *ValidationError) {
	e.Errors = append(e.Errors, err)
}

// HasErrors returns true if there are any validation errors
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// RepositoryError represents an error that occurred in the repository layer
type RepositoryError struct {
	Err     error
	Message string
	Code    string
}

func (e *RepositoryError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// NewRepositoryError creates a new RepositoryError
func NewRepositoryError(err error, message, code string) *RepositoryError {
	return &RepositoryError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	ResourceType string
	ID           string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.ResourceType, e.ID)
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(resourceType, id string) *NotFoundError {
	return &NotFoundError{
		ResourceType: resourceType,
		ID:           id,
	}
}

// Is implements the errors.Is interface for NotFoundError
func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

// ApplicationError represents an error that occurred in the application layer
type ApplicationError struct {
	Err     error
	Message string
	Code    string
}

func (e *ApplicationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *ApplicationError) Unwrap() error {
	return e.Err
}

// NewApplicationError creates a new ApplicationError
func NewApplicationError(err error, message, code string) *ApplicationError {
	return &ApplicationError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// Error represents a domain error with operation context
type Error struct {
	// Original is the original error
	Original error

	// Code is a machine-readable error code
	Code string

	// Message is a human-readable error message
	Message string

	// Op is the operation that caused the error
	Op string

	// Param is the parameter that caused the error
	Param string
}

// Error returns a string representation of the error
func (e *Error) Error() string {
	if e.Original != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Original)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

// Unwrap returns the original error
func (e *Error) Unwrap() error {
	return e.Original
}

// Is reports whether the error is of the given target type
func (e *Error) Is(target error) bool {
	if target == nil {
		return e == nil
	}

	t, ok := target.(*Error)
	if !ok {
		return errors.Is(e.Original, target)
	}

	return e.Code == t.Code
}

// New creates a new Error
func New(op, code, message string, original error) error {
	return &Error{
		Original: original,
		Code:     code,
		Message:  message,
		Op:       op,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, op, message string) error {
	if err == nil {
		return nil
	}

	// If it's already a domain error, just update the op and message
	if e, ok := err.(*Error); ok {
		return &Error{
			Original: e.Original,
			Code:     e.Code,
			Message:  message,
			Op:       op,
		}
	}

	// Otherwise, create a new domain error
	return &Error{
		Original: err,
		Message:  message,
		Op:       op,
	}
}

// IsNotFound returns true if the error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsInvalidInput returns true if the error is an invalid input error
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsUnauthorized returns true if the error is an unauthorized error
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden returns true if the error is a forbidden error
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsInternal returns true if the error is an internal error
func IsInternal(err error) bool {
	return errors.Is(err, ErrInternal)
}

// IsTimeout returns true if the error is a timeout error
func IsTimeout(err error) bool {
	return errors.Is(err, ErrTimeout)
}

// IsCancelled returns true if the error is a cancelled error
func IsCancelled(err error) bool {
	return errors.Is(err, ErrCancelled)
}

// IsConflict returns true if the error is a conflict error
func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

// GenericError is a generic error type that can be used for different error categories
type GenericError[T any] struct {
	Err      error
	Message  string
	Code     string
	Category T
}

func (e *GenericError[T]) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *GenericError[T]) Unwrap() error {
	return e.Err
}

// NewGenericError creates a new GenericError
func NewGenericError[T any](err error, message, code string, category T) *GenericError[T] {
	return &GenericError[T]{
		Err:      err,
		Message:  message,
		Code:     code,
		Category: category,
	}
}
