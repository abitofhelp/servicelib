// Copyright (c) 2025 A Bit of Help, Inc.

// Package domain provides domain-specific error types for the application.
package domain

import (
	"github.com/abitofhelp/servicelib/errors/core"
)

// DomainError represents a domain-specific error.
// It extends BaseError with domain-specific information.
type DomainError struct {
	*core.BaseError
}

// NewDomainError creates a new DomainError.
func NewDomainError(code core.ErrorCode, message string, cause error) *DomainError {
	return &DomainError{
		BaseError: core.NewBaseError(code, message, cause),
	}
}

// IsDomainError identifies this as a domain error.
func (e *DomainError) IsDomainError() bool {
	return true
}

// As implements the errors.As interface for DomainError.
func (e *DomainError) As(target interface{}) bool {
	// Check if target is *DomainError
	if t, ok := target.(*DomainError); ok {
		*t = *e
		return true
	}

	// Delegate to BaseError.As for other types
	return e.BaseError.As(target)
}

// ValidationError represents a validation error.
// It extends DomainError with field-specific information.
type ValidationError struct {
	*DomainError
	Field string `json:"field,omitempty"`
}

// NewValidationError creates a new ValidationError.
func NewValidationError(message string, field string, cause error) *ValidationError {
	return &ValidationError{
		DomainError: NewDomainError(core.ValidationErrorCode, message, cause),
		Field:       field,
	}
}

// IsValidationError identifies this as a validation error.
func (e *ValidationError) IsValidationError() bool {
	return true
}

// As implements the errors.As interface for ValidationError.
func (e *ValidationError) As(target interface{}) bool {
	// Debug print
	println("ValidationError.As called with target type:", target)

	// Check if target is *ValidationError
	_, isValidationErr := target.(*ValidationError)
	println("Is target *ValidationError?", isValidationErr)
	if isValidationErr {
		println("Target is *ValidationError")
		*target.(*ValidationError) = *e
		return true
	}

	// Check if target is *DomainError
	_, isDomainErr := target.(*DomainError)
	println("Is target *DomainError?", isDomainErr)
	if isDomainErr {
		println("Target is *DomainError")
		*target.(*DomainError) = *e.DomainError
		return true
	}

	// Delegate to DomainError.As for other types
	println("Delegating to DomainError.As")
	return e.DomainError.As(target)
}

// ValidationErrors represents multiple validation errors.
type ValidationErrors struct {
	*DomainError
	Errors []*ValidationError `json:"errors,omitempty"`
}

// NewValidationErrors creates a new ValidationErrors.
func NewValidationErrors(message string, errors ...*ValidationError) *ValidationErrors {
	return &ValidationErrors{
		DomainError: NewDomainError(core.ValidationErrorCode, message, nil),
		Errors:      errors,
	}
}

// AddError adds a validation error to the collection.
func (e *ValidationErrors) AddError(err *ValidationError) {
	e.Errors = append(e.Errors, err)
}

// HasErrors returns true if there are any validation errors.
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// As implements the errors.As interface for ValidationErrors.
func (e *ValidationErrors) As(target interface{}) bool {
	// Check if target is *ValidationErrors
	if t, ok := target.(*ValidationErrors); ok {
		*t = *e
		return true
	}

	// Check if target is *DomainError
	if t, ok := target.(*DomainError); ok {
		*t = *e.DomainError
		return true
	}

	// Delegate to DomainError.As for other types
	return e.DomainError.As(target)
}

// BusinessRuleError represents a business rule violation.
// It extends DomainError with rule-specific information.
type BusinessRuleError struct {
	*DomainError
	Rule string `json:"rule,omitempty"`
}

// NewBusinessRuleError creates a new BusinessRuleError.
func NewBusinessRuleError(message string, rule string, cause error) *BusinessRuleError {
	return &BusinessRuleError{
		DomainError: NewDomainError(core.BusinessRuleViolationCode, message, cause),
		Rule:        rule,
	}
}

// IsBusinessRuleError identifies this as a business rule error.
func (e *BusinessRuleError) IsBusinessRuleError() bool {
	return true
}

// As implements the errors.As interface for BusinessRuleError.
func (e *BusinessRuleError) As(target interface{}) bool {
	// Check if target is *BusinessRuleError
	if t, ok := target.(*BusinessRuleError); ok {
		*t = *e
		return true
	}

	// Check if target is *DomainError
	if t, ok := target.(*DomainError); ok {
		*t = *e.DomainError
		return true
	}

	// Delegate to DomainError.As for other types
	return e.DomainError.As(target)
}

// NotFoundError represents a resource not found error.
// It extends DomainError with resource-specific information.
type NotFoundError struct {
	*DomainError
	ResourceType string `json:"resource_type,omitempty"`
	ResourceID   string `json:"resource_id,omitempty"`
}

// NewNotFoundError creates a new NotFoundError.
func NewNotFoundError(resourceType string, resourceID string, cause error) *NotFoundError {
	message := resourceType + " with ID " + resourceID + " not found"
	return &NotFoundError{
		DomainError:  NewDomainError(core.NotFoundCode, message, cause),
		ResourceType: resourceType,
		ResourceID:   resourceID,
	}
}

// IsNotFoundError identifies this as a not found error.
func (e *NotFoundError) IsNotFoundError() bool {
	return true
}

// As implements the errors.As interface for NotFoundError.
func (e *NotFoundError) As(target interface{}) bool {
	// Check if target is *NotFoundError
	if t, ok := target.(*NotFoundError); ok {
		*t = *e
		return true
	}

	// Check if target is *DomainError
	if t, ok := target.(*DomainError); ok {
		*t = *e.DomainError
		return true
	}

	// Delegate to DomainError.As for other types
	return e.DomainError.As(target)
}
