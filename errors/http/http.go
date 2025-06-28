// Copyright (c) 2025 A Bit of Help, Inc.

// Package http provides HTTP-related error utilities for the application.
package http

import (
	"encoding/json"
	"net/http"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/errors/core"
)

// ErrorResponse represents an error response to be sent to clients.
type ErrorResponse struct {
	// Code is the error code
	Code string `json:"code,omitempty"`

	// Message is the error message
	Message string `json:"message,omitempty"`

	// Details contains additional information about the error
	Details map[string]interface{} `json:"details,omitempty"`
}

// WriteError writes an error response to the HTTP response writer.
// It sets the appropriate status code and content type, and writes the error as JSON.
func WriteError(w http.ResponseWriter, err error) {
	// Get the HTTP status code
	status := errors.GetHTTPStatus(err)
	if status == 0 {
		status = http.StatusInternalServerError
	}

	// Create the error response
	response := ErrorResponse{
		Message: err.Error(),
	}

	// Add error code if available
	if e, ok := err.(interface{ GetCode() core.ErrorCode }); ok {
		response.Code = string(e.GetCode())
	}

	// Add details if available
	if e, ok := err.(interface{ GetDetails() map[string]interface{} }); ok {
		response.Details = e.GetDetails()
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// WriteErrorWithStatus writes an error response with a specific status code.
// This is useful when you want to override the default status code mapping.
func WriteErrorWithStatus(w http.ResponseWriter, err error, status int) {
	// Create the error response
	response := ErrorResponse{
		Message: err.Error(),
	}

	// Add error code if available
	if e, ok := err.(interface{ GetCode() core.ErrorCode }); ok {
		response.Code = string(e.GetCode())
	}

	// Add details if available
	if e, ok := err.(interface{ GetDetails() map[string]interface{} }); ok {
		response.Details = e.GetDetails()
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// GetErrorFromRequest extracts an error from a failed HTTP request.
// It creates an appropriate error type based on the response status code.
func GetErrorFromRequest(resp *http.Response, body []byte) error {
	// Create a message from the body if available
	message := string(body)
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}

	// Map the status code to an error code
	var code core.ErrorCode
	switch resp.StatusCode {
	case http.StatusNotFound:
		code = core.NotFoundCode
	case http.StatusBadRequest:
		code = core.InvalidInputCode
	case http.StatusUnauthorized:
		code = core.UnauthorizedCode
	case http.StatusForbidden:
		code = core.ForbiddenCode
	case http.StatusRequestTimeout:
		code = core.TimeoutCode
	case http.StatusConflict:
		code = core.AlreadyExistsCode
	case http.StatusTooManyRequests:
		code = core.ResourceExhaustedCode
	case http.StatusBadGateway:
		code = core.ExternalServiceErrorCode
	case http.StatusServiceUnavailable:
		code = core.NetworkErrorCode
	case http.StatusGatewayTimeout:
		code = core.TimeoutCode
	default:
		code = core.InternalErrorCode
	}

	// Create the error
	return errors.New(code, message)
}
