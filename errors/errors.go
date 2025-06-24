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

// ErrorCode is an alias for core.ErrorCode
type ErrorCode = core.ErrorCode

// Error code constants
const (
	NotFoundCode              = core.NotFoundCode
	InvalidInputCode          = core.InvalidInputCode
	DatabaseErrorCode         = core.DatabaseErrorCode
	InternalErrorCode         = core.InternalErrorCode
	TimeoutCode               = core.TimeoutCode
	CanceledCode              = core.CanceledCode
	AlreadyExistsCode         = core.AlreadyExistsCode
	UnauthorizedCode          = core.UnauthorizedCode
	ForbiddenCode             = core.ForbiddenCode
	ValidationErrorCode       = core.ValidationErrorCode
	BusinessRuleViolationCode = core.BusinessRuleViolationCode
	ExternalServiceErrorCode  = core.ExternalServiceErrorCode
	NetworkErrorCode          = core.NetworkErrorCode
	ConfigurationErrorCode    = core.ConfigurationErrorCode
	ResourceExhaustedCode     = core.ResourceExhaustedCode
	DataCorruptionCode        = core.DataCorruptionCode
	ConcurrencyErrorCode      = core.ConcurrencyErrorCode
)

// Standard errors
var (
	ErrNotFound      = core.NewBaseError(NotFoundCode, "resource not found", nil)
	ErrInvalidInput  = core.NewBaseError(InvalidInputCode, "invalid input", nil)
	ErrInternal      = core.NewBaseError(InternalErrorCode, "internal error", nil)
	ErrUnauthorized  = core.NewBaseError(UnauthorizedCode, "unauthorized", nil)
	ErrForbidden     = core.NewBaseError(ForbiddenCode, "forbidden", nil)
	ErrTimeout       = core.NewBaseError(TimeoutCode, "operation timed out", nil)
	ErrCanceled      = core.NewBaseError(CanceledCode, "operation canceled", nil)
	ErrAlreadyExists = core.NewBaseError(AlreadyExistsCode, "resource already exists", nil)
)

// Error type aliases
type (
	BaseError            = core.BaseError
	DomainError          = domain.DomainError
	ValidationError      = domain.ValidationError
	ValidationErrors     = domain.ValidationErrors
	BusinessRuleError    = domain.BusinessRuleError
	NotFoundError        = domain.NotFoundError
	InfrastructureError  = infra.InfrastructureError
	DatabaseError        = infra.DatabaseError
	NetworkError         = infra.NetworkError
	ExternalServiceError = infra.ExternalServiceError
	RetryError           = infra.RetryError
	ContextError         = infra.ContextError
	ApplicationError     = app.ApplicationError
	ConfigurationError   = app.ConfigurationError
	AuthenticationError  = app.AuthenticationError
	AuthorizationError   = app.AuthorizationError
)

// New creates a new BaseError.
func New(code ErrorCode, message string) error {
	return core.NewBaseError(code, message, nil)
}

// Wrap wraps an error with additional context.
func Wrap(err error, code ErrorCode, message string) error {
	if err == nil {
		return nil
	}
	return core.NewBaseError(code, message, err)
}

// WrapWithOperation wraps an error with an operation.
func WrapWithOperation(err error, code ErrorCode, message string, operation string) error {
	if err == nil {
		return nil
	}
	return core.NewBaseError(code, message, err).WithOperation(operation)
}

// WrapWithDetails wraps an error with details.
func WrapWithDetails(err error, code ErrorCode, message string, details map[string]interface{}) error {
	if err == nil {
		return nil
	}
	return core.NewBaseError(code, message, err).WithDetails(details)
}

// Domain error creation functions
func NewDomainError(code ErrorCode, message string, cause error) *DomainError {
	return domain.NewDomainError(code, message, cause)
}

func NewValidationError(message string, field string, cause error) *ValidationError {
	return domain.NewValidationError(message, field, cause)
}

func NewValidationErrors(message string, errors ...*ValidationError) *ValidationErrors {
	return domain.NewValidationErrors(message, errors...)
}

func NewBusinessRuleError(message string, rule string, cause error) *BusinessRuleError {
	return domain.NewBusinessRuleError(message, rule, cause)
}

func NewNotFoundError(resourceType string, resourceID string, cause error) *NotFoundError {
	return domain.NewNotFoundError(resourceType, resourceID, cause)
}

// Infrastructure error creation functions
func NewInfrastructureError(code ErrorCode, message string, cause error) *InfrastructureError {
	return infra.NewInfrastructureError(code, message, cause)
}

func NewDatabaseError(message string, operation string, table string, cause error) *DatabaseError {
	return infra.NewDatabaseError(message, operation, table, cause)
}

func NewNetworkError(message string, host string, port string, cause error) *NetworkError {
	return infra.NewNetworkError(message, host, port, cause)
}

func NewExternalServiceError(message string, serviceName string, endpoint string, cause error) *ExternalServiceError {
	return infra.NewExternalServiceError(message, serviceName, endpoint, cause)
}

func NewRetryError(message string, cause error, attempts int, maxAttempts int) *RetryError {
	return infra.NewRetryError(message, cause, attempts, maxAttempts)
}

func NewContextError(message string, cause error) *ContextError {
	return infra.NewContextError(message, cause)
}

// Application error creation functions
func NewApplicationError(code ErrorCode, message string, cause error) *ApplicationError {
	return app.NewApplicationError(code, message, cause)
}

func NewConfigurationError(message string, configKey string, configValue string, cause error) *ConfigurationError {
	return app.NewConfigurationError(message, configKey, configValue, cause)
}

func NewAuthenticationError(message string, username string, cause error) *AuthenticationError {
	return app.NewAuthenticationError(message, username, cause)
}

func NewAuthorizationError(message string, username string, resource string, action string, cause error) *AuthorizationError {
	return app.NewAuthorizationError(message, username, resource, action, cause)
}

// Error checking functions
func Is(err, target error) bool {
	return core.Is(err, target)
}

func As(err error, target interface{}) bool {
	return core.As(err, target)
}

func Unwrap(err error) error {
	return core.Unwrap(err)
}

// GetHTTPStatus returns the HTTP status code for an error.
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

// ToJSON converts an error to a JSON string.
func ToJSON(err error) string {
	return core.ToJSON(err)
}

// Error type checking functions
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

func IsValidationError(err error) bool {
	var e *ValidationError
	return As(err, &e)
}

func IsBusinessRuleError(err error) bool {
	var e *BusinessRuleError
	return As(err, &e)
}

func IsNotFoundError(err error) bool {
	var e *NotFoundError
	return As(err, &e)
}

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

func IsDatabaseError(err error) bool {
	var e *DatabaseError
	return As(err, &e)
}

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

func IsExternalServiceError(err error) bool {
	var e *ExternalServiceError
	return As(err, &e)
}

func IsRetryError(err error) bool {
	var e *RetryError
	return As(err, &e)
}

func IsContextError(err error) bool {
	var e *ContextError
	return As(err, &e)
}

func IsConfigurationError(err error) bool {
	var e *ConfigurationError
	return As(err, &e)
}

func IsAuthenticationError(err error) bool {
	var e *AuthenticationError
	return As(err, &e)
}

func IsAuthorizationError(err error) bool {
	var e *AuthorizationError
	return As(err, &e)
}

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
