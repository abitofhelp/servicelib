// Copyright (c) 2024 A Bit of Help, Inc.

package core

import "net/http"

// Map of error codes to HTTP status codes
var errorCodeToHTTPStatus = map[ErrorCode]int{
	NotFoundCode:              http.StatusNotFound,
	InvalidInputCode:          http.StatusBadRequest,
	DatabaseErrorCode:         http.StatusInternalServerError,
	InternalErrorCode:         http.StatusInternalServerError,
	TimeoutCode:              http.StatusGatewayTimeout,
	CanceledCode:             http.StatusRequestTimeout,
	AlreadyExistsCode:        http.StatusConflict,
	UnauthorizedCode:         http.StatusUnauthorized,
	ForbiddenCode:            http.StatusForbidden,
	ValidationErrorCode:      http.StatusBadRequest,
	BusinessRuleViolationCode: http.StatusUnprocessableEntity,
	ExternalServiceErrorCode: http.StatusBadGateway,
	NetworkErrorCode:         http.StatusServiceUnavailable,
	ConfigurationErrorCode:   http.StatusInternalServerError,
	ResourceExhaustedCode:    http.StatusTooManyRequests,
	DataCorruptionCode:      http.StatusInternalServerError,
	ConcurrencyErrorCode:     http.StatusConflict,
}

// GetHTTPStatus returns the HTTP status code for an error code
func GetHTTPStatus(code ErrorCode) int {
	if status, ok := errorCodeToHTTPStatus[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}
