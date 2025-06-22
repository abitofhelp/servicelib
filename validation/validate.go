// Copyright (c) 2025 A Bit of Help, Inc.

package validation

import (
	"regexp"
	"strings"
	"time"

	"github.com/abitofhelp/servicelib/errors"
)

// ValidationResult holds the result of a validation operation
type ValidationResult struct {
	errors *errors.ValidationErrors
}

// NewValidationResult creates a new ValidationResult
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		errors: errors.NewValidationErrors("Validation failed"),
	}
}

// AddError adds an error to the validation result
func (v *ValidationResult) AddError(msg, field string) {
	v.errors.AddError(errors.NewValidationError(msg, field, nil))
}

// IsValid returns true if there are no validation errors
func (v *ValidationResult) IsValid() bool {
	return !v.errors.HasErrors()
}

// Error returns the validation errors as an error
func (v *ValidationResult) Error() error {
	if v.IsValid() {
		return nil
	}
	return v.errors
}

// Required validates that a string is not empty
func Required(value, field string, result *ValidationResult) {
	if strings.TrimSpace(value) == "" {
		result.AddError("is required", field)
	}
}

// MinLength validates that a string has a minimum length
func MinLength(value string, min int, field string, result *ValidationResult) {
	if len(value) < min {
		result.AddError(
			"must be at least "+string(rune('0'+min))+" characters long",
			field,
		)
	}
}

// MaxLength validates that a string has a maximum length
func MaxLength(value string, max int, field string, result *ValidationResult) {
	if len(value) > max {
		result.AddError(
			"must be at most "+string(rune('0'+max))+" characters long",
			field,
		)
	}
}

// Pattern validates that a string matches a regular expression
func Pattern(value, pattern, field string, result *ValidationResult) {
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		result.AddError("has an invalid format", field)
	}
}

// PastDate validates that a date is in the past
func PastDate(value time.Time, field string, result *ValidationResult) {
	if value.After(time.Now()) {
		result.AddError("must be in the past", field)
	}
}

// ValidDateRange validates that a start date is before an end date
func ValidDateRange(start, end time.Time, startField, endField string, result *ValidationResult) {
	if !end.IsZero() && start.After(end) {
		result.AddError("must be before "+endField, startField)
	}
}

// AllTrue validates that all items in a slice satisfy a predicate
func AllTrue[T any](items []T, predicate func(T) bool) bool {
	for _, item := range items {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// ValidateAll validates all items in a slice
func ValidateAll[T any](items []T, validator func(T, int, *ValidationResult), result *ValidationResult) {
	for i, item := range items {
		validator(item, i, result)
	}
}

// ValidateID validates that an ID is not empty
func ValidateID(id, field string, result *ValidationResult) {
	if strings.TrimSpace(id) == "" {
		result.AddError("is required", field)
	}
}
