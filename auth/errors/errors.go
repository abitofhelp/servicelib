// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides comprehensive error handling for the auth module.
// It defines error types and provides context-aware error handling.
//
// IMPORTANT: This package is now implemented using the main errors framework
// while maintaining backward compatibility with existing code. New code should
// use the main errors framework directly instead of this package.
//
// The main errors framework provides a more comprehensive error handling system
// with better integration with the hybrid architecture. It includes error types
// for different layers (domain, application, infrastructure) and better support
// for error context, HTTP status mapping, and error recovery.
package errors

import (
	"fmt"

	"github.com/abitofhelp/servicelib/errors"
)

// Standard error types
var (
	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.NewAuthenticationError("invalid token", "", nil)

	// ErrExpiredToken is returned when a token has expired
	ErrExpiredToken = errors.NewAuthenticationError("token expired", "", nil)

	// ErrMissingToken is returned when a token is missing
	ErrMissingToken = errors.NewAuthenticationError("token missing", "", nil)

	// ErrInvalidSignature is returned when a token has an invalid signature
	ErrInvalidSignature = errors.NewAuthenticationError("invalid token signature", "", nil)

	// ErrInvalidClaims is returned when a token has invalid claims
	ErrInvalidClaims = errors.NewAuthenticationError("invalid token claims", "", nil)

	// ErrUnauthorized is returned when a user is not authorized to perform an operation
	ErrUnauthorized = errors.ErrUnauthorized

	// ErrForbidden is returned when a user is forbidden from performing an operation
	ErrForbidden = errors.ErrForbidden

	// ErrInvalidConfig is returned when the configuration is invalid
	ErrInvalidConfig = errors.NewConfigurationError("invalid configuration", "", "", nil)

	// ErrInternal is returned when an internal error occurs
	ErrInternal = errors.ErrInternal

	// ErrNotImplemented is returned when a feature is not implemented
	ErrNotImplemented = errors.New(errors.InternalErrorCode, "not implemented")
)

// AuthError represents an authentication or authorization error with context
// It's now a wrapper around the new errors framework's AuthenticationError
type AuthError struct {
	// Err is the underlying error
	Err error

	// Op is the operation that caused the error
	Op string

	// Message is a human-readable message
	Message string

	// Context contains additional context for the error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *AuthError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// Unwrap returns the underlying error
func (e *AuthError) Unwrap() error {
	return e.Err
}

// Is reports whether the error is of the target type
func (e *AuthError) Is(target error) bool {
	return errors.Is(e.Err, target)
}

// NewAuthError creates a new AuthError
func NewAuthError(err error, op string, message string, context map[string]interface{}) *AuthError {
	// Create a new AuthError that wraps the provided error
	return &AuthError{
		Err:     err,
		Op:      op,
		Message: message,
		Context: context,
	}
}

// getErrorCode returns an appropriate ErrorCode for the given error
func getErrorCode(err error) errors.ErrorCode {
	// Try to get the code from the error if it has a GetCode method
	if e, ok := err.(interface{ GetCode() errors.ErrorCode }); ok {
		return e.GetCode()
	}

	// Map HTTP status codes to error codes
	switch errors.GetHTTPStatus(err) {
	case 400:
		return errors.InvalidInputCode
	case 401:
		return errors.UnauthorizedCode
	case 403:
		return errors.ForbiddenCode
	case 404:
		return errors.NotFoundCode
	case 409:
		return errors.AlreadyExistsCode
	case 422:
		return errors.ValidationErrorCode
	case 429:
		return errors.ResourceExhaustedCode
	case 500:
		return errors.InternalErrorCode
	case 502, 503, 504:
		return errors.ExternalServiceErrorCode
	default:
		return errors.InternalErrorCode
	}
}

// WithContext adds context to an error
func WithContext(err error, key string, value interface{}) error {
	details := map[string]interface{}{key: value}

	// If it's already an AuthError, update its context
	var authErr *AuthError
	if errors.As(err, &authErr) {
		if authErr.Context == nil {
			authErr.Context = make(map[string]interface{})
		}
		authErr.Context[key] = value
		return authErr
	}

	// Otherwise, create a new AuthError
	return &AuthError{
		Err:     err,
		Context: details,
	}
}

// WithOp adds an operation to an error
func WithOp(err error, op string) error {
	// If it's already an AuthError, update its operation
	var authErr *AuthError
	if errors.As(err, &authErr) {
		authErr.Op = op
		return authErr
	}

	// Otherwise, create a new AuthError
	return &AuthError{
		Err: err,
		Op:  op,
	}
}

// WithMessage adds a message to an error
func WithMessage(err error, message string) error {
	// If it's already an AuthError, update its message
	var authErr *AuthError
	if errors.As(err, &authErr) {
		authErr.Message = message
		return authErr
	}

	// Otherwise, create a new AuthError
	return &AuthError{
		Err:     err,
		Message: message,
	}
}

// Wrap wraps an error with a message
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// GetContext gets a context value from an error
func GetContext(err error, key string) (interface{}, bool) {
	// Try to get context from AuthError
	var authErr *AuthError
	if errors.As(err, &authErr) && authErr.Context != nil {
		value, ok := authErr.Context[key]
		return value, ok
	}

	// Try to get details from the new errors framework
	if e, ok := err.(interface{ GetDetails() map[string]interface{} }); ok {
		details := e.GetDetails()
		if details != nil {
			if value, exists := details[key]; exists {
				return value, true
			}
		}
	}

	return nil, false
}

// GetOp gets the operation from an error
func GetOp(err error) (string, bool) {
	// Try to get operation from AuthError
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Op, authErr.Op != ""
	}

	// Try to get operation from the new errors framework
	if e, ok := err.(interface{ GetOperation() string }); ok {
		op := e.GetOperation()
		return op, op != ""
	}

	return "", false
}

// GetMessage gets the message from an error
func GetMessage(err error) (string, bool) {
	// Try to get message from AuthError
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Message, authErr.Message != ""
	}

	// For non-AuthError errors, return false and empty string for backward compatibility
	return "", false
}
