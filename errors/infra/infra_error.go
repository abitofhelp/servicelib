// Copyright (c) 2025 A Bit of Help, Inc.

// Package infra provides infrastructure-related error types for the application.
package infra

import (
	"context"
	"github.com/abitofhelp/servicelib/errors/core"
)

// InfrastructureError represents an infrastructure-related error.
// It extends BaseError with infrastructure-specific information.
type InfrastructureError struct {
	*core.BaseError
}

// NewInfrastructureError creates a new InfrastructureError.
func NewInfrastructureError(code core.ErrorCode, message string, cause error) *InfrastructureError {
	return &InfrastructureError{
		BaseError: core.NewBaseError(code, message, cause),
	}
}

// IsInfrastructureError identifies this as an infrastructure error.
func (e *InfrastructureError) IsInfrastructureError() bool {
	return true
}

// RetryError represents an error that occurred during a retry operation.
// It extends InfrastructureError with retry-specific information.
type RetryError struct {
	*InfrastructureError
	Attempts    int `json:"attempts,omitempty"`
	MaxAttempts int `json:"max_attempts,omitempty"`
}

// NewRetryError creates a new RetryError.
func NewRetryError(message string, cause error, attempts int, maxAttempts int) *RetryError {
	return &RetryError{
		InfrastructureError: NewInfrastructureError(core.InternalErrorCode, message, cause),
		Attempts:            attempts,
		MaxAttempts:         maxAttempts,
	}
}

// IsRetryError identifies this as a retry error.
func (e *RetryError) IsRetryError() bool {
	return true
}

// ContextError represents an error that occurred due to a context cancellation or timeout during a retry operation.
// It extends InfrastructureError with context-specific information.
type ContextError struct {
	*InfrastructureError
}

// NewContextError creates a new ContextError.
func NewContextError(message string, cause error) *ContextError {
	var code core.ErrorCode
	if cause == context.Canceled {
		code = core.CanceledCode
	} else if cause == context.DeadlineExceeded {
		code = core.TimeoutCode
	} else {
		code = core.InternalErrorCode
	}

	return &ContextError{
		InfrastructureError: NewInfrastructureError(code, message, cause),
	}
}

// IsContextError identifies this as a context error.
func (e *ContextError) IsContextError() bool {
	return true
}

// DatabaseError represents a database operation failure.
// It extends InfrastructureError with database-specific information.
type DatabaseError struct {
	*InfrastructureError
	Operation string `json:"operation,omitempty"`
	Table     string `json:"table,omitempty"`
}

// NewDatabaseError creates a new DatabaseError.
func NewDatabaseError(message string, operation string, table string, cause error) *DatabaseError {
	return &DatabaseError{
		InfrastructureError: NewInfrastructureError(core.DatabaseErrorCode, message, cause),
		Operation:           operation,
		Table:               table,
	}
}

// IsDatabaseError identifies this as a database error.
func (e *DatabaseError) IsDatabaseError() bool {
	return true
}

// NetworkError represents a network-related error.
// It extends InfrastructureError with network-specific information.
type NetworkError struct {
	*InfrastructureError
	Host string `json:"host,omitempty"`
	Port string `json:"port,omitempty"`
}

// NewNetworkError creates a new NetworkError.
func NewNetworkError(message string, host string, port string, cause error) *NetworkError {
	return &NetworkError{
		InfrastructureError: NewInfrastructureError(core.NetworkErrorCode, message, cause),
		Host:                host,
		Port:                port,
	}
}

// IsNetworkError identifies this as a network error.
func (e *NetworkError) IsNetworkError() bool {
	return true
}

// ExternalServiceError represents an external service call failure.
// It extends InfrastructureError with service-specific information.
type ExternalServiceError struct {
	*InfrastructureError
	ServiceName string `json:"service_name,omitempty"`
	Endpoint    string `json:"endpoint,omitempty"`
}

// NewExternalServiceError creates a new ExternalServiceError.
func NewExternalServiceError(message string, serviceName string, endpoint string, cause error) *ExternalServiceError {
	return &ExternalServiceError{
		InfrastructureError: NewInfrastructureError(core.ExternalServiceErrorCode, message, cause),
		ServiceName:         serviceName,
		Endpoint:            endpoint,
	}
}

// IsExternalServiceError identifies this as an external service error.
func (e *ExternalServiceError) IsExternalServiceError() bool {
	return true
}
