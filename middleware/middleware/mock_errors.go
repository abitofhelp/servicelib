// Copyright (c) 2025 A Bit of Help, Inc.

package middleware

import (
	"net/http"
)

// MockValidationError is a mock implementation of the ValidationErrorInterface
type MockValidationError struct {
	Msg string
}

func (e *MockValidationError) Error() string {
	return e.Msg
}

func (e *MockValidationError) Code() string {
	return "VALIDATION_ERROR"
}

func (e *MockValidationError) HTTPStatus() int {
	return http.StatusBadRequest
}

func (e *MockValidationError) IsValidationError() bool {
	return true
}

// MockNotFoundError is a mock implementation of the NotFoundErrorInterface
type MockNotFoundError struct {
	ResourceType string
	ID           string
}

func (e *MockNotFoundError) Error() string {
	return e.ResourceType + " with ID " + e.ID + " not found"
}

func (e *MockNotFoundError) Code() string {
	return "NOT_FOUND"
}

func (e *MockNotFoundError) HTTPStatus() int {
	return http.StatusNotFound
}

func (e *MockNotFoundError) IsNotFoundError() bool {
	return true
}

// MockApplicationError is a mock implementation of the ApplicationErrorInterface
type MockApplicationError struct {
	Msg      string
	CodeName string
}

func (e *MockApplicationError) Error() string {
	return e.Msg
}

func (e *MockApplicationError) Code() string {
	return e.CodeName
}

func (e *MockApplicationError) HTTPStatus() int {
	return http.StatusInternalServerError
}

func (e *MockApplicationError) IsApplicationError() bool {
	return true
}

// MockRepositoryError is a mock implementation of the RepositoryErrorInterface
type MockRepositoryError struct {
	Msg      string
	CodeName string
}

func (e *MockRepositoryError) Error() string {
	return e.Msg
}

func (e *MockRepositoryError) Code() string {
	return e.CodeName
}

func (e *MockRepositoryError) HTTPStatus() int {
	return http.StatusInternalServerError
}

func (e *MockRepositoryError) IsRepositoryError() bool {
	return true
}
