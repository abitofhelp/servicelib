// Copyright (c) 2025 A Bit of Help, Inc.

package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/stretchr/testify/assert"
)

// TestWriteError tests the WriteError function
func TestWriteError(t *testing.T) {
	testCases := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   string
		expectedMsg    string
	}{
		{
			name:           "Standard error",
			err:            fmt.Errorf("standard error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "",
			expectedMsg:    "standard error",
		},
		{
			name:           "Not found error",
			err:            errors.NewNotFoundError("User", "123", nil),
			expectedStatus: http.StatusNotFound,
			expectedCode:   string(core.NotFoundCode),
			expectedMsg:    "User with ID 123 not found",
		},
		{
			name:           "Validation error",
			err:            errors.NewValidationError("Invalid input", "field", nil),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   string(core.ValidationErrorCode),
			expectedMsg:    "Invalid input",
		},
		{
			name:           "Database error",
			err:            errors.NewDatabaseError("Database error", "SELECT", "users", nil),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   string(core.DatabaseErrorCode),
			expectedMsg:    "Database error",
		},
		{
			name:           "Authentication error",
			err:            errors.NewAuthenticationError("Invalid credentials", "user", nil),
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   string(core.UnauthorizedCode),
			expectedMsg:    "Invalid credentials",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the function
			WriteError(rr, tc.err)

			// Check status code
			assert.Equal(t, tc.expectedStatus, rr.Code)

			// Check content type
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			// Parse the response body
			var response ErrorResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check response fields
			assert.Equal(t, tc.expectedCode, response.Code)
			assert.Contains(t, response.Message, tc.expectedMsg)
		})
	}
}

// TestWriteErrorWithStatus tests the WriteErrorWithStatus function
func TestWriteErrorWithStatus(t *testing.T) {
	// Create a test error
	err := errors.NewValidationError("Invalid input", "field", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the function with a custom status code
	customStatus := http.StatusTeapot // 418 I'm a teapot
	WriteErrorWithStatus(rr, err, customStatus)

	// Check status code
	assert.Equal(t, customStatus, rr.Code)

	// Check content type
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Parse the response body
	var response ErrorResponse
	unmarshalErr := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, unmarshalErr)

	// Check response fields
	assert.Equal(t, string(core.ValidationErrorCode), response.Code)
	assert.Contains(t, response.Message, "Invalid input")
}

// TestGetErrorFromRequest tests the GetErrorFromRequest function
func TestGetErrorFromRequest(t *testing.T) {
	testCases := []struct {
		name           string
		statusCode     int
		body           string
		expectedCode   core.ErrorCode
		expectedMsg    string
	}{
		{
			name:           "Not found",
			statusCode:     http.StatusNotFound,
			body:           "Resource not found",
			expectedCode:   core.NotFoundCode,
			expectedMsg:    "Resource not found",
		},
		{
			name:           "Bad request",
			statusCode:     http.StatusBadRequest,
			body:           "Invalid input",
			expectedCode:   core.InvalidInputCode,
			expectedMsg:    "Invalid input",
		},
		{
			name:           "Unauthorized",
			statusCode:     http.StatusUnauthorized,
			body:           "Unauthorized",
			expectedCode:   core.UnauthorizedCode,
			expectedMsg:    "Unauthorized",
		},
		{
			name:           "Forbidden",
			statusCode:     http.StatusForbidden,
			body:           "Forbidden",
			expectedCode:   core.ForbiddenCode,
			expectedMsg:    "Forbidden",
		},
		{
			name:           "Request timeout",
			statusCode:     http.StatusRequestTimeout,
			body:           "Request timeout",
			expectedCode:   core.TimeoutCode,
			expectedMsg:    "Request timeout",
		},
		{
			name:           "Conflict",
			statusCode:     http.StatusConflict,
			body:           "Conflict",
			expectedCode:   core.AlreadyExistsCode,
			expectedMsg:    "Conflict",
		},
		{
			name:           "Too many requests",
			statusCode:     http.StatusTooManyRequests,
			body:           "Too many requests",
			expectedCode:   core.ResourceExhaustedCode,
			expectedMsg:    "Too many requests",
		},
		{
			name:           "Bad gateway",
			statusCode:     http.StatusBadGateway,
			body:           "Bad gateway",
			expectedCode:   core.ExternalServiceErrorCode,
			expectedMsg:    "Bad gateway",
		},
		{
			name:           "Service unavailable",
			statusCode:     http.StatusServiceUnavailable,
			body:           "Service unavailable",
			expectedCode:   core.NetworkErrorCode,
			expectedMsg:    "Service unavailable",
		},
		{
			name:           "Gateway timeout",
			statusCode:     http.StatusGatewayTimeout,
			body:           "Gateway timeout",
			expectedCode:   core.TimeoutCode,
			expectedMsg:    "Gateway timeout",
		},
		{
			name:           "Internal server error",
			statusCode:     http.StatusInternalServerError,
			body:           "Internal server error",
			expectedCode:   core.InternalErrorCode,
			expectedMsg:    "Internal server error",
		},
		{
			name:           "Empty body",
			statusCode:     http.StatusNotFound,
			body:           "",
			expectedCode:   core.NotFoundCode,
			expectedMsg:    "Not Found", // Default message from http.StatusText
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock HTTP response
			resp := &http.Response{
				StatusCode: tc.statusCode,
			}

			// Call the function
			err := GetErrorFromRequest(resp, []byte(tc.body))

			// Check that the error is not nil
			assert.NotNil(t, err)

			// Check error code
			var baseErr *core.BaseError
			assert.True(t, errors.As(err, &baseErr))
			assert.Equal(t, tc.expectedCode, baseErr.GetCode())

			// Check error message
			assert.Equal(t, tc.expectedMsg, baseErr.GetMessage())
		})
	}
}
