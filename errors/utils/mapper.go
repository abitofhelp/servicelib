// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"github.com/abitofhelp/servicelib/errors/core"
)

// ErrorCodeMapper maps error codes to various representations
var ErrorCodeMapper = struct {
	// ToHTTPStatus returns the HTTP status code for an error code
	ToHTTPStatus func(code core.ErrorCode) int
	// ToCategory returns the category for an error code
	ToCategory func(code core.ErrorCode) string
	// ToMessage returns a user-friendly message for an error code
	ToMessage func(code core.ErrorCode) string
}{
	ToHTTPStatus: func(code core.ErrorCode) int {
		return core.GetHTTPStatus(code)
	},
	ToCategory: func(code core.ErrorCode) string {
		switch code {
		case core.NotFoundCode,
			core.InvalidInputCode,
			core.AlreadyExistsCode,
			core.ResourceExhaustedCode,
			core.ValidationErrorCode,
			core.BusinessRuleViolationCode:
			return "Client"
		case core.DatabaseErrorCode,
			core.InternalErrorCode,
			core.DataCorruptionCode,
			core.ConfigurationErrorCode:
			return "Server"
		case core.TimeoutCode,
			core.CanceledCode,
			core.ConcurrencyErrorCode:
			return "System"
		case core.ExternalServiceErrorCode,
			core.NetworkErrorCode:
			return "External"
		case core.UnauthorizedCode,
			core.ForbiddenCode:
			return "Security"
		default:
			return "Unknown"
		}
	},
	ToMessage: func(code core.ErrorCode) string {
		switch code {
		case core.NotFoundCode:
			return "The requested resource was not found"
		case core.InvalidInputCode:
			return "The provided input is invalid"
		case core.DatabaseErrorCode:
			return "A database operation failed"
		case core.InternalErrorCode:
			return "An internal server error occurred"
		case core.TimeoutCode:
			return "The operation timed out"
		case core.CanceledCode:
			return "The operation was canceled"
		case core.AlreadyExistsCode:
			return "The resource already exists"
		case core.UnauthorizedCode:
			return "Authentication failed"
		case core.ForbiddenCode:
			return "Access is forbidden"
		case core.ValidationErrorCode:
			return "Input validation failed"
		case core.BusinessRuleViolationCode:
			return "A business rule was violated"
		case core.ExternalServiceErrorCode:
			return "An external service call failed"
		case core.NetworkErrorCode:
			return "A network error occurred"
		case core.ConfigurationErrorCode:
			return "Configuration error"
		case core.ResourceExhaustedCode:
			return "Resource limit reached"
		case core.DataCorruptionCode:
			return "Data corruption detected"
		case core.ConcurrencyErrorCode:
			return "Concurrency violation"
		default:
			return "An unknown error occurred"
		}
	},
}
