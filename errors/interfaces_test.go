// Copyright (c) 2025 A Bit of Help, Inc.

package errors

import (
	"net/http"
	"testing"

	"github.com/abitofhelp/servicelib/errors/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestIsValidationError(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	tests := []struct {
		name     string
		setupErr func() error
		expected bool
	}{
		{
			name: "ValidationError returns true",
			setupErr: func() error {
				mockErr := mocks.NewMockValidationErrorInterface(ctrl)
				mockErr.EXPECT().IsValidationError().Return(true).AnyTimes()
				return mockErr
			},
			expected: true,
		},
		{
			name: "Non-ValidationError returns false",
			setupErr: func() error {
				return ErrNotFound
			},
			expected: false,
		},
		{
			name: "Nil error returns false",
			setupErr: func() error {
				return nil
			},
			expected: false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setupErr()
			result := IsValidationError(err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	tests := []struct {
		name     string
		setupErr func() error
		expected bool
	}{
		{
			name: "NotFoundError returns true",
			setupErr: func() error {
				mockErr := mocks.NewMockNotFoundErrorInterface(ctrl)
				mockErr.EXPECT().IsNotFoundError().Return(true).AnyTimes()
				return mockErr
			},
			expected: true,
		},
		{
			name: "Non-NotFoundError returns false",
			setupErr: func() error {
				return ErrInvalidInput
			},
			expected: false,
		},
		{
			name: "Nil error returns false",
			setupErr: func() error {
				return nil
			},
			expected: false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setupErr()
			result := IsNotFoundError(err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsApplicationError(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	tests := []struct {
		name     string
		setupErr func() error
		expected bool
	}{
		{
			name: "ApplicationError returns true",
			setupErr: func() error {
				mockErr := mocks.NewMockApplicationErrorInterface(ctrl)
				mockErr.EXPECT().IsApplicationError().Return(true).AnyTimes()
				return mockErr
			},
			expected: true,
		},
		{
			name: "Non-ApplicationError returns false",
			setupErr: func() error {
				return ErrInvalidInput
			},
			expected: false,
		},
		{
			name: "Nil error returns false",
			setupErr: func() error {
				return nil
			},
			expected: false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setupErr()
			result := IsApplicationError(err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsRepositoryError(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	tests := []struct {
		name     string
		setupErr func() error
		expected bool
	}{
		{
			name: "RepositoryError returns true",
			setupErr: func() error {
				mockErr := mocks.NewMockRepositoryErrorInterface(ctrl)
				mockErr.EXPECT().IsRepositoryError().Return(true).AnyTimes()
				return mockErr
			},
			expected: true,
		},
		{
			name: "Non-RepositoryError returns false",
			setupErr: func() error {
				return ErrInvalidInput
			},
			expected: false,
		},
		{
			name: "Nil error returns false",
			setupErr: func() error {
				return nil
			},
			expected: false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setupErr()
			result := IsRepositoryError(err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetHTTPStatusFromError(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	tests := []struct {
		name           string
		setupErr       func() error
		expectedStatus int
	}{
		{
			name: "ErrorWithHTTPStatus returns its status",
			setupErr: func() error {
				mockErr := mocks.NewMockErrorWithHTTPStatus(ctrl)
				mockErr.EXPECT().HTTPStatus().Return(http.StatusTeapot).AnyTimes()
				return mockErr
			},
			expectedStatus: http.StatusTeapot,
		},
		{
			name: "ValidationError returns BadRequest",
			setupErr: func() error {
				mockErr := mocks.NewMockValidationErrorInterface(ctrl)
				// First it will check if it implements ErrorWithHTTPStatus
				mockErr.EXPECT().HTTPStatus().Return(http.StatusBadRequest).AnyTimes()
				// Then it will check if it's a ValidationError
				mockErr.EXPECT().IsValidationError().Return(true).AnyTimes()
				return mockErr
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "NotFoundError returns NotFound",
			setupErr: func() error {
				mockErr := mocks.NewMockNotFoundErrorInterface(ctrl)
				// First it will check if it implements ErrorWithHTTPStatus
				mockErr.EXPECT().HTTPStatus().Return(http.StatusNotFound).AnyTimes()
				// Then it will check if it's a NotFoundError
				mockErr.EXPECT().IsNotFoundError().Return(true).AnyTimes()
				return mockErr
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "ApplicationError returns InternalServerError",
			setupErr: func() error {
				mockErr := mocks.NewMockApplicationErrorInterface(ctrl)
				// First it will check if it implements ErrorWithHTTPStatus
				mockErr.EXPECT().HTTPStatus().Return(http.StatusInternalServerError).AnyTimes()
				// Then it will check if it's an ApplicationError
				mockErr.EXPECT().IsApplicationError().Return(true).AnyTimes()
				return mockErr
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "RepositoryError returns InternalServerError",
			setupErr: func() error {
				mockErr := mocks.NewMockRepositoryErrorInterface(ctrl)
				// First it will check if it implements ErrorWithHTTPStatus
				mockErr.EXPECT().HTTPStatus().Return(http.StatusInternalServerError).AnyTimes()
				// Then it will check if it's a RepositoryError
				mockErr.EXPECT().IsRepositoryError().Return(true).AnyTimes()
				return mockErr
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Other error returns InternalServerError",
			setupErr: func() error {
				return ErrInvalidInput
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Nil error returns InternalServerError",
			setupErr: func() error {
				return nil
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setupErr()
			status := GetHTTPStatusFromError(err)
			assert.Equal(t, tt.expectedStatus, status)
		})
	}
}
