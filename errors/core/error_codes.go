// Copyright (c) 2024 A Bit of Help, Inc.

// Package core provides the core error handling functionality for the errors package.
// It includes error codes, error context, and basic error types.
package core

import (
	"fmt"
	"net/http"
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
	ErrNotFound = fmt.Errorf("resource not found")

	// ErrAlreadyExists is returned when a resource already exists
	ErrAlreadyExists = fmt.Errorf("resource already exists")

	// ErrInvalidInput is returned when the input to a function is invalid
	ErrInvalidInput = fmt.Errorf("invalid input")

	// ErrInternal is returned when an internal error occurs
	ErrInternal = fmt.Errorf("internal error")

	// ErrUnauthorized is returned when a user is not authorized to perform an action
	ErrUnauthorized = fmt.Errorf("unauthorized")

	// ErrForbidden is returned when a user is forbidden from performing an action
	ErrForbidden = fmt.Errorf("forbidden")

	// ErrTimeout is returned when an operation times out
	ErrTimeout = fmt.Errorf("operation timed out")

	// ErrCancelled is returned when an operation is cancelled
	ErrCancelled = fmt.Errorf("operation cancelled")

	// ErrConflict is returned when there is a conflict with the current state
	ErrConflict = fmt.Errorf("conflict with current state")
)

// Map of error codes to HTTP status codes
var ErrorCodeToHTTPStatus = map[ErrorCode]int{
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
