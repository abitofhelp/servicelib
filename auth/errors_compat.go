// Copyright (c) 2025 A Bit of Help, Inc.

// Package auth provides authentication and authorization functionality.
// This file contains compatibility functions for the old auth/errors package.
package auth

import (
	"github.com/abitofhelp/servicelib/errors"
)

// Error constants for backward compatibility
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

// WithContext adds context to an error for backward compatibility
func WithContext(err error, key string, value interface{}) error {
	details := map[string]interface{}{key: value}
	return errors.WrapWithDetails(err, getErrorCode(err), "", details)
}

// WithOp adds an operation to an error for backward compatibility
func WithOp(err error, op string) error {
	return errors.WrapWithOperation(err, getErrorCode(err), "", op)
}

// WithMessage adds a message to an error for backward compatibility
func WithMessage(err error, message string) error {
	return errors.Wrap(err, getErrorCode(err), message)
}

// Wrap wraps an error with a message for backward compatibility
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, getErrorCode(err), message)
}

// GetContext gets a context value from an error for backward compatibility
func GetContext(err error, key string) (interface{}, bool) {
	// Try to get details from various error types
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

// GetOp gets the operation from an error for backward compatibility
func GetOp(err error) (string, bool) {
	// Try to get operation from various error types
	if e, ok := err.(interface{ GetOperation() string }); ok {
		op := e.GetOperation()
		return op, op != ""
	}
	return "", false
}

// GetMessage gets the message from an error for backward compatibility
func GetMessage(err error) (string, bool) {
	// Try to get message from various error types
	if e, ok := err.(interface{ Error() string }); ok {
		msg := e.Error()
		return msg, msg != ""
	}
	return "", false
}
