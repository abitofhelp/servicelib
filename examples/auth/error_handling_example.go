// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating error handling in the auth module
package example_auth

import (
	"errors"
	"fmt"
)

// This is a simplified example to demonstrate error handling concepts
// In a real application, you would use the auth module's error handling

// Define error types
var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrMissingClaims    = errors.New("missing claims")
)

// ContextError represents an error with context information
type ContextError struct {
	Op      string            // Operation that caused the error
	Err     error             // Original error
	Message string            // Human-readable error message
	Context map[string]string // Additional context information
}

// Error implements the error interface
func (e *ContextError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// Unwrap returns the wrapped error
func (e *ContextError) Unwrap() error {
	return e.Err
}

// NewError creates a new context error
func NewError(op string, err error, message string) *ContextError {
	return &ContextError{
		Op:      op,
		Err:     err,
		Message: message,
		Context: make(map[string]string),
	}
}

// WithContext adds context information to the error
func (e *ContextError) WithContext(key, value string) *ContextError {
	e.Context[key] = value
	return e
}

// GetContext gets context information from an error
func GetContext(err error, key string) (string, bool) {
	var contextErr *ContextError
	if errors.As(err, &contextErr) {
		value, ok := contextErr.Context[key]
		return value, ok
	}
	return "", false
}

// GetOp gets the operation that caused the error
func GetOp(err error) (string, bool) {
	var contextErr *ContextError
	if errors.As(err, &contextErr) {
		return contextErr.Op, true
	}
	return "", false
}

// GetMessage gets the error message
func GetMessage(err error) (string, bool) {
	var contextErr *ContextError
	if errors.As(err, &contextErr) {
		return contextErr.Message, true
	}
	return "", false
}

// SimulateTokenValidation simulates token validation with error handling
func SimulateTokenValidation(token string) error {
	// Simulate token validation
	if token == "" {
		return NewError(
			"ValidateToken",
			ErrInvalidToken,
			"Token cannot be empty",
		).WithContext("token_length", "0")
	}

	if token == "expired" {
		return NewError(
			"ValidateToken",
			ErrTokenExpired,
			"Token has expired",
		).WithContext("token", token).WithContext("expiry", "2023-01-01")
	}

	if token == "invalid_signature" {
		return NewError(
			"ValidateToken",
			ErrInvalidSignature,
			"Token has an invalid signature",
		).WithContext("token", token)
	}

	if token == "missing_claims" {
		return NewError(
			"ValidateToken",
			ErrMissingClaims,
			"Token is missing required claims",
		).WithContext("token", token).WithContext("missing", "sub,exp")
	}

	return nil
}

func main() {
	// Test with valid token
	err := SimulateTokenValidation("valid_token")
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
	} else {
		fmt.Println("Token is valid")
	}

	// Test with empty token
	err = SimulateTokenValidation("")
	if err != nil {
		// Check for specific error types
		if errors.Is(err, ErrInvalidToken) {
			fmt.Println("Error: Invalid token")
		}

		// Get context information
		if tokenLength, ok := GetContext(err, "token_length"); ok {
			fmt.Printf("Token length: %s\n", tokenLength)
		}

		// Get operation
		if op, ok := GetOp(err); ok {
			fmt.Printf("Operation: %s\n", op)
		}

		// Get message
		if message, ok := GetMessage(err); ok {
			fmt.Printf("Message: %s\n", message)
		}
	}

	// Test with expired token
	err = SimulateTokenValidation("expired")
	if err != nil {
		if errors.Is(err, ErrTokenExpired) {
			fmt.Println("Error: Token expired")
		}

		if expiry, ok := GetContext(err, "expiry"); ok {
			fmt.Printf("Expiry date: %s\n", expiry)
		}
	}

	// Test with invalid signature
	err = SimulateTokenValidation("invalid_signature")
	if err != nil {
		if errors.Is(err, ErrInvalidSignature) {
			fmt.Println("Error: Invalid signature")
		}
	}

	// Test with missing claims
	err = SimulateTokenValidation("missing_claims")
	if err != nil {
		if errors.Is(err, ErrMissingClaims) {
			fmt.Println("Error: Missing claims")
		}

		if missing, ok := GetContext(err, "missing"); ok {
			fmt.Printf("Missing claims: %s\n", missing)
		}
	}
}
