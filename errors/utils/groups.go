// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"github.com/abitofhelp/servicelib/errors/core"
)

// ErrorCodeGroup represents a group of related error codes
type ErrorCodeGroup string

const (
	// ResourceGroup contains resource-related errors
	ResourceGroup ErrorCodeGroup = "RESOURCE"
	// InputGroup contains input validation errors
	InputGroup ErrorCodeGroup = "INPUT"
	// SystemGroup contains system-level errors
	SystemGroup ErrorCodeGroup = "SYSTEM"
	// SecurityGroup contains security-related errors
	SecurityGroup ErrorCodeGroup = "SECURITY"
	// ExternalGroup contains external system errors
	ExternalGroup ErrorCodeGroup = "EXTERNAL"
	// BusinessGroup contains business rule errors
	BusinessGroup ErrorCodeGroup = "BUSINESS"
)

// ErrorCodeGrouper groups error codes
type ErrorCodeGrouper struct {
	// groupFunc returns the group of an error code
	groupFunc func(code core.ErrorCode) ErrorCodeGroup
}

var grouper = &ErrorCodeGrouper{
	groupFunc: func(code core.ErrorCode) ErrorCodeGroup {
		switch code {
		case core.NotFoundCode,
			core.AlreadyExistsCode:
			return ResourceGroup
		case core.InvalidInputCode,
			core.ValidationErrorCode:
			return InputGroup
		case core.DatabaseErrorCode,
			core.InternalErrorCode,
			core.TimeoutCode,
			core.CanceledCode,
			core.ConcurrencyErrorCode,
			core.ConfigurationErrorCode:
			return SystemGroup
		case core.UnauthorizedCode,
			core.ForbiddenCode:
			return SecurityGroup
		case core.ExternalServiceErrorCode,
			core.NetworkErrorCode:
			return ExternalGroup
		case core.ResourceExhaustedCode,
			core.DataCorruptionCode,
			core.BusinessRuleViolationCode:
			return BusinessGroup
		default:
			return SystemGroup
		}
	},
}

// NewErrorCodeGrouper creates a new error code grouper
func NewErrorCodeGrouper() *ErrorCodeGrouper {
	return grouper
}

// Group returns the group of an error code
func (g *ErrorCodeGrouper) Group(code core.ErrorCode) ErrorCodeGroup {
	return g.groupFunc(code)
}

// IsResourceError returns true if the error code is a resource error
func (g *ErrorCodeGrouper) IsResourceError(code core.ErrorCode) bool {
	return g.Group(code) == ResourceGroup
}

// IsInputError returns true if the error code is an input error
func (g *ErrorCodeGrouper) IsInputError(code core.ErrorCode) bool {
	return g.Group(code) == InputGroup
}

// IsSystemError returns true if the error code is a system error
func (g *ErrorCodeGrouper) IsSystemError(code core.ErrorCode) bool {
	return g.Group(code) == SystemGroup
}

// IsSecurityError returns true if the error code is a security error
func (g *ErrorCodeGrouper) IsSecurityError(code core.ErrorCode) bool {
	return g.Group(code) == SecurityGroup
}

// IsExternalError returns true if the error code is an external error
func (g *ErrorCodeGrouper) IsExternalError(code core.ErrorCode) bool {
	return g.Group(code) == ExternalGroup
}

// IsBusinessError returns true if the error code is a business error
func (g *ErrorCodeGrouper) IsBusinessError(code core.ErrorCode) bool {
	return g.Group(code) == BusinessGroup
}
