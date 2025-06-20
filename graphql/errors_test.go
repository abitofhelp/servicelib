// Copyright (c) 2025 A Bit of Help, Inc.

package graphql

import (
	"context"
	"errors"
	"testing"

	myerrors "github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

// TestHandleError tests the HandleError function
func TestHandleError(t *testing.T) {
	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Test cases
	tests := []struct {
		name          string
		err           error
		operation     string
		expectedCode  string
		expectedMsg   string
		checkExactMsg bool
	}{
		{
			name:          "Context canceled",
			err:           context.Canceled,
			operation:     "TestOperation",
			expectedCode:  "CANCELED",
			expectedMsg:   "Operation canceled",
			checkExactMsg: true,
		},
		{
			name:          "Context deadline exceeded",
			err:           context.DeadlineExceeded,
			operation:     "TestOperation",
			expectedCode:  "TIMEOUT",
			expectedMsg:   "Operation timed out",
			checkExactMsg: true,
		},
		{
			name:          "Validation error",
			err:           myerrors.NewValidationError("Invalid input"),
			operation:     "TestOperation",
			expectedCode:  "VALIDATION_ERROR",
			expectedMsg:   "validation error: Invalid input",
			checkExactMsg: true,
		},
		{
			name:          "Not found error",
			err:           myerrors.NewNotFoundError("User", "123"),
			operation:     "TestOperation",
			expectedCode:  "NOT_FOUND",
			expectedMsg:   "User with ID 123 not found",
			checkExactMsg: true,
		},
		{
			name:          "Domain error",
			err:           myerrors.NewDomainError(errors.New("business rule violation"), "Cannot delete active user", "ACTIVE_USER"),
			operation:     "TestOperation",
			expectedCode:  "ACTIVE_USER",
			expectedMsg:   "Cannot delete active user",
			checkExactMsg: true,
		},
		{
			name:          "Application error",
			err:           myerrors.NewApplicationError(errors.New("internal error"), "Application error", "APP_ERROR"),
			operation:     "TestOperation",
			expectedCode:  "APP_ERROR",
			expectedMsg:   "Application error",
			checkExactMsg: true,
		},
		{
			name:          "Repository error",
			err:           myerrors.NewRepositoryError(errors.New("database error"), "Failed to query database", "DB_ERROR"),
			operation:     "TestOperation",
			expectedCode:  "INTERNAL_ERROR",
			expectedMsg:   "An internal error occurred",
			checkExactMsg: true,
		},
		{
			name:          "Generic error",
			err:           errors.New("some error"),
			operation:     "TestOperation",
			expectedCode:  "INTERNAL_ERROR",
			expectedMsg:   "An unexpected error occurred",
			checkExactMsg: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function
			result := HandleError(context.Background(), tt.err, tt.operation, logger)

			// Check that the result is a GraphQL error
			gqlErr, ok := result.(*gqlerror.Error)
			assert.True(t, ok, "Result should be a GraphQL error")

			// Check the error code
			code, ok := gqlErr.Extensions["code"].(string)
			assert.True(t, ok, "Extensions should contain a code")
			assert.Equal(t, tt.expectedCode, code)

			// Check the error message
			if tt.checkExactMsg {
				assert.Equal(t, tt.expectedMsg, gqlErr.Message)
			} else {
				assert.Contains(t, gqlErr.Message, tt.expectedMsg)
			}
		})
	}
}
