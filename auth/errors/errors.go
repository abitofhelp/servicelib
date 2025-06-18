// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides comprehensive error handling for the auth module.
// It defines error types and provides context-aware error handling.
package errors

import (
	"errors"
	"fmt"
)

// Standard error types
var (
	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrExpiredToken is returned when a token has expired
	ErrExpiredToken = errors.New("token expired")

	// ErrMissingToken is returned when a token is missing
	ErrMissingToken = errors.New("token missing")

	// ErrInvalidSignature is returned when a token has an invalid signature
	ErrInvalidSignature = errors.New("invalid token signature")

	// ErrInvalidClaims is returned when a token has invalid claims
	ErrInvalidClaims = errors.New("invalid token claims")

	// ErrUnauthorized is returned when a user is not authorized to perform an operation
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when a user is forbidden from performing an operation
	ErrForbidden = errors.New("forbidden")

	// ErrInvalidConfig is returned when the configuration is invalid
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrInternal is returned when an internal error occurs
	ErrInternal = errors.New("internal error")

	// ErrNotImplemented is returned when a feature is not implemented
	ErrNotImplemented = errors.New("not implemented")
)

// AuthError represents an authentication or authorization error with context
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
	return &AuthError{
		Err:     err,
		Op:      op,
		Message: message,
		Context: context,
	}
}

// WithContext adds context to an error
func WithContext(err error, key string, value interface{}) error {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		if authErr.Context == nil {
			authErr.Context = make(map[string]interface{})
		}
		authErr.Context[key] = value
		return authErr
	}
	return &AuthError{
		Err:     err,
		Context: map[string]interface{}{key: value},
	}
}

// WithOp adds an operation to an error
func WithOp(err error, op string) error {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		authErr.Op = op
		return authErr
	}
	return &AuthError{
		Err: err,
		Op:  op,
	}
}

// WithMessage adds a message to an error
func WithMessage(err error, message string) error {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		authErr.Message = message
		return authErr
	}
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
	var authErr *AuthError
	if errors.As(err, &authErr) && authErr.Context != nil {
		value, ok := authErr.Context[key]
		return value, ok
	}
	return nil, false
}

// GetOp gets the operation from an error
func GetOp(err error) (string, bool) {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Op, authErr.Op != ""
	}
	return "", false
}

// GetMessage gets the message from an error
func GetMessage(err error) (string, bool) {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Message, authErr.Message != ""
	}
	return "", false
}
