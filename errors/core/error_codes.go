// Copyright (c) 2024 A Bit of Help, Inc.

package core

import "fmt"

// ErrorCode represents a unique error code for categorizing errors.
type ErrorCode string

// Standard error codes define all possible error categories in the application.
const (
	// NotFoundCode is used when a resource is not found.
	NotFoundCode ErrorCode = "NOT_FOUND"

	// InvalidInputCode is used when input validation fails.
	InvalidInputCode ErrorCode = "INVALID_INPUT"

	// DatabaseErrorCode is used for database operation failures.
	DatabaseErrorCode ErrorCode = "DATABASE_ERROR"

	// InternalErrorCode is used for internal server errors.
	InternalErrorCode ErrorCode = "INTERNAL_ERROR"

	// TimeoutCode is used when an operation times out.
	TimeoutCode ErrorCode = "TIMEOUT"

	// CanceledCode is used when an operation is canceled.
	CanceledCode ErrorCode = "CANCELED"

	// AlreadyExistsCode is used when a resource already exists.
	AlreadyExistsCode ErrorCode = "ALREADY_EXISTS"

	// UnauthorizedCode is used for authentication failures.
	UnauthorizedCode ErrorCode = "UNAUTHORIZED"

	// ForbiddenCode is used for authorization failures.
	ForbiddenCode ErrorCode = "FORBIDDEN"

	// ValidationErrorCode is used for domain validation errors.
	ValidationErrorCode ErrorCode = "VALIDATION_ERROR"

	// BusinessRuleViolationCode is used when a business rule is violated.
	BusinessRuleViolationCode ErrorCode = "BUSINESS_RULE_VIOLATION"

	// ExternalServiceErrorCode is used when an external service call fails.
	ExternalServiceErrorCode ErrorCode = "EXTERNAL_SERVICE_ERROR"

	// NetworkErrorCode is used for network-related errors.
	NetworkErrorCode ErrorCode = "NETWORK_ERROR"

	// ConfigurationErrorCode is used for configuration errors.
	ConfigurationErrorCode ErrorCode = "CONFIGURATION_ERROR"

	// ResourceExhaustedCode is used when a resource limit is reached.
	ResourceExhaustedCode ErrorCode = "RESOURCE_EXHAUSTED"

	// DataCorruptionCode is used when data is corrupted.
	DataCorruptionCode ErrorCode = "DATA_CORRUPTION"

 // ConcurrencyErrorCode is used for concurrency-related errors.
	ConcurrencyErrorCode ErrorCode = "CONCURRENCY_ERROR"
)

// Standard error variables
var (
	// ErrNotFound is used when a resource is not found.
	ErrNotFound = fmt.Errorf("not found: %s", NotFoundCode)

	// ErrInvalidInput is used when input validation fails.
	ErrInvalidInput = fmt.Errorf("invalid input: %s", InvalidInputCode)
)
