// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides a comprehensive error handling system for applications.
//
// This package implements a structured approach to error handling, offering:
//   - A hierarchy of error types for different domains (application, domain, infrastructure)
//   - Error codes with HTTP status mapping
//   - Contextual information for debugging
//   - Stack traces for error origin tracking
//   - Utilities for creating, wrapping, and serializing errors
//   - Error type checking and categorization
//
// The error types are organized into several categories:
//   - Domain errors: Represent business rule violations and domain-specific issues
//   - Infrastructure errors: Represent issues with external systems and resources
//   - Application errors: Represent issues with the application itself
//
// Each error type provides specific contextual information relevant to its category,
// making it easier to understand, log, and respond to errors appropriately.
//
// Example usage:
//
//	// Create a validation error
//	err := errors.NewValidationError("Invalid email format", "email", nil)
//
//	// Check error type
//	if errors.IsValidationError(err) {
//	    // Handle validation error
//	}
//
//	// Create a database error
//	dbErr := errors.NewDatabaseError("Failed to insert record", "INSERT", "users", err)
//
//	// Get HTTP status code for error
//	statusCode := errors.GetHTTPStatus(dbErr)
package errors
