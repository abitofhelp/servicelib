// Copyright (c) 2025 A Bit of Help, Inc.

// Package validation provides utilities for validating data in a structured and reusable way.
//
// This package offers a collection of validation functions for common validation tasks,
// such as checking required fields, validating string lengths, matching patterns,
// validating dates, and more. It also provides a ValidationResult type for collecting
// and managing validation errors.
//
// The validation functions are designed to be composable, allowing you to build complex
// validation logic by combining simple validation rules. The package integrates with
// the errors package to provide detailed validation error information.
//
// Key features:
//   - String validation (required, min/max length, pattern matching)
//   - Date validation (past dates, date ranges)
//   - Collection validation (validate all items, check if all items satisfy a condition)
//   - Structured validation results with field-specific error messages
//   - Integration with the errors package for consistent error handling
//
// Example usage:
//
//	// Create a validation result
//	result := validation.NewValidationResult()
//
//	// Validate a user struct
//	validation.Required(user.Name, "name", result)
//	validation.MinLength(user.Password, 8, "password", result)
//	validation.Pattern(user.Email, `^[^@]+@[^@]+\.[^@]+$`, "email", result)
//
//	// Check if validation passed
//	if !result.IsValid() {
//	    return result.Error()
//	}
//
//	// Validate a collection
//	validation.ValidateAll(user.Addresses, func(addr Address, i int, result *validation.ValidationResult) {
//	    validation.Required(addr.Street, fmt.Sprintf("addresses[%d].street", i), result)
//	    validation.Required(addr.City, fmt.Sprintf("addresses[%d].city", i), result)
//	    validation.Required(addr.Country, fmt.Sprintf("addresses[%d].country", i), result)
//	}, result)
//
// The package is designed to be used throughout the application to ensure data integrity
// and provide consistent validation error messages to users.
package validation
