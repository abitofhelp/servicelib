// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"github.com/abitofhelp/servicelib/errors/core"
)

// ErrorCodeCategory represents a category of error codes
type ErrorCodeCategory string

const (
	// ClientError indicates an error caused by client input
	ClientError ErrorCodeCategory = "CLIENT_ERROR"
	// ServerError indicates an error caused by server-side issues
	ServerError ErrorCodeCategory = "SERVER_ERROR"
	// SystemError indicates an error caused by system-level issues
	SystemError ErrorCodeCategory = "SYSTEM_ERROR"
	// ExternalError indicates an error caused by external systems
	ExternalError ErrorCodeCategory = "EXTERNAL_ERROR"
	// SecurityError indicates a security-related error
	SecurityError ErrorCodeCategory = "SECURITY_ERROR"
)

// ErrorCodeCategorizer categorizes error codes
type ErrorCodeCategorizer struct {
	// categorizeFunc returns the category of an error code
	categorizeFunc func(code core.ErrorCode) ErrorCodeCategory
}

var categorizer = &ErrorCodeCategorizer{
	categorizeFunc: func(code core.ErrorCode) ErrorCodeCategory {
		switch code {
		case core.NotFoundCode,
			core.InvalidInputCode,
			core.AlreadyExistsCode,
			core.ValidationErrorCode,
			core.BusinessRuleViolationCode:
			return ClientError
		case core.DatabaseErrorCode,
			core.InternalErrorCode,
			core.DataCorruptionCode,
			core.ConfigurationErrorCode,
			core.ResourceExhaustedCode:
			return ServerError
		case core.TimeoutCode,
			core.CanceledCode,
			core.ConcurrencyErrorCode:
			return SystemError
		case core.ExternalServiceErrorCode,
			core.NetworkErrorCode:
			return ExternalError
		case core.UnauthorizedCode,
			core.ForbiddenCode:
			return SecurityError
		default:
			return ServerError
		}
	},
}

// NewErrorCodeCategorizer creates a new error code categorizer
func NewErrorCodeCategorizer() *ErrorCodeCategorizer {
	return categorizer
}

// Categorize returns the category of an error code
func (c *ErrorCodeCategorizer) Categorize(code core.ErrorCode) ErrorCodeCategory {
	return c.categorizeFunc(code)
}

// IsClientError returns true if the error code is a client error
func (c *ErrorCodeCategorizer) IsClientError(code core.ErrorCode) bool {
	return c.Categorize(code) == ClientError
}

// IsServerError returns true if the error code is a server error
func (c *ErrorCodeCategorizer) IsServerError(code core.ErrorCode) bool {
	return c.Categorize(code) == ServerError
}

// IsSystemError returns true if the error code is a system error
func (c *ErrorCodeCategorizer) IsSystemError(code core.ErrorCode) bool {
	return c.Categorize(code) == SystemError
}

// IsExternalError returns true if the error code is an external error
func (c *ErrorCodeCategorizer) IsExternalError(code core.ErrorCode) bool {
	return c.Categorize(code) == ExternalError
}

// IsSecurityError returns true if the error code is a security error
func (c *ErrorCodeCategorizer) IsSecurityError(code core.ErrorCode) bool {
	return c.Categorize(code) == SecurityError
}
