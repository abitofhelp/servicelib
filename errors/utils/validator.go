// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"fmt"
	"github.com/abitofhelp/servicelib/errors/core"
)

// ErrorCodeValidator validates error codes
type ErrorCodeValidator struct {
}

// NewErrorCodeValidator creates a new error code validator
func NewErrorCodeValidator() *ErrorCodeValidator {
	return &ErrorCodeValidator{}
}

// Validate validates an error code
func (v *ErrorCodeValidator) Validate(code core.ErrorCode) error {
	if !v.IsValid(code) {
		return fmt.Errorf("invalid error code: %v", code)
	}
	return nil
}

// GetCodeName returns the human-readable name of an error code
func (v *ErrorCodeValidator) GetCodeName(code core.ErrorCode) string {
	return getHumanReadableName(code)
}

// IsValid returns true if the error code is valid
func (v *ErrorCodeValidator) IsValid(code core.ErrorCode) bool {
	return IsValid(code)
}

// IsValid returns true if the error code is valid
func IsValid(code core.ErrorCode) bool {
	switch code {
	case core.NotFoundCode,
		core.InvalidInputCode,
		core.DatabaseErrorCode,
		core.InternalErrorCode,
		core.TimeoutCode,
		core.CanceledCode,
		core.AlreadyExistsCode,
		core.UnauthorizedCode,
		core.ForbiddenCode,
		core.ValidationErrorCode,
		core.BusinessRuleViolationCode,
		core.ExternalServiceErrorCode,
		core.NetworkErrorCode,
		core.ConfigurationErrorCode,
		core.ResourceExhaustedCode,
		core.DataCorruptionCode,
		core.ConcurrencyErrorCode:
		return true
	default:
		return false
	}
}

// getHumanReadableName returns the human-readable name of an error code
func getHumanReadableName(code core.ErrorCode) string {
	switch code {
	case core.NotFoundCode:
		return "NotFound"
	case core.InvalidInputCode:
		return "InvalidInput"
	case core.DatabaseErrorCode:
		return "Database"
	case core.InternalErrorCode:
		return "Internal"
	case core.TimeoutCode:
		return "Timeout"
	case core.CanceledCode:
		return "Canceled"
	case core.AlreadyExistsCode:
		return "AlreadyExists"
	case core.UnauthorizedCode:
		return "Unauthorized"
	case core.ForbiddenCode:
		return "Forbidden"
	case core.ValidationErrorCode:
		return "Validation"
	case core.BusinessRuleViolationCode:
		return "BusinessRuleViolation"
	case core.ExternalServiceErrorCode:
		return "ExternalService"
	case core.NetworkErrorCode:
		return "Network"
	case core.ConfigurationErrorCode:
		return "Configuration"
	case core.ResourceExhaustedCode:
		return "ResourceExhausted"
	case core.DataCorruptionCode:
		return "DataCorruption"
	case core.ConcurrencyErrorCode:
		return "Concurrency"
	default:
		return "Unknown"
	}
}
