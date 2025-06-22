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