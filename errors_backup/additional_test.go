// Copyright (c) 2025 A Bit of Help, Inc.

package errors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWithDetailsEdgeCases tests edge cases for the WithDetails function
func TestWithDetailsEdgeCases(t *testing.T) {
	// Test with nil error
	t.Run("NilError", func(t *testing.T) {
		result := WithDetails(nil, map[string]interface{}{"key": "value"})
		assert.Nil(t, result)
	})

	// Test with nil details
	t.Run("NilDetails", func(t *testing.T) {
		err := fmt.Errorf("test error")
		result := WithDetails(err, nil)
		assert.NotNil(t, result)
		assert.Contains(t, result.Error(), "test error")
	})

	// Test with empty details
	t.Run("EmptyDetails", func(t *testing.T) {
		err := fmt.Errorf("test error")
		result := WithDetails(err, map[string]interface{}{})
		assert.NotNil(t, result)
		assert.Contains(t, result.Error(), "test error")
	})
}

// TestToJSONEdgeCases tests edge cases for the ToJSON function
func TestToJSONEdgeCases(t *testing.T) {
	// Test with nil error
	t.Run("NilError", func(t *testing.T) {
		result := ToJSON(nil)
		assert.Equal(t, "{}", result)
	})

	// Test with non-ContextualError
	t.Run("NonContextualError", func(t *testing.T) {
		err := fmt.Errorf("simple error")
		result := ToJSON(err)
		assert.Contains(t, result, "simple error")
	})
}

// TestWrapEdgeCases tests edge cases for the Wrap function
func TestWrapEdgeCases(t *testing.T) {
	// Test with nil error
	t.Run("NilError", func(t *testing.T) {
		result := Wrap(nil, "operation", "message")
		assert.Nil(t, result)
	})

	// Test with domain error
	t.Run("DomainError", func(t *testing.T) {
		originalErr := &Error{
			Original: fmt.Errorf("original error"),
			Code:     "original_code",
			Message:  "original message",
			Op:       "original_op",
		}
		result := Wrap(originalErr, "new_op", "new message")
		assert.NotNil(t, result)

		// Check that it's a domain error
		domainErr, ok := result.(*Error)
		assert.True(t, ok)

		// Check that the original error is preserved
		assert.Equal(t, originalErr.Original, domainErr.Original)

		// Check that the code is preserved
		assert.Equal(t, "original_code", domainErr.Code)

		// Check that the message and op are updated
		assert.Equal(t, "new message", domainErr.Message)
		assert.Equal(t, "new_op", domainErr.Op)
	})
}

// TestGetCodeEdgeCases tests edge cases for the GetCode function
func TestGetCodeEdgeCases(t *testing.T) {
	// Test with nil error
	t.Run("NilError", func(t *testing.T) {
		result := GetCode(nil)
		assert.Equal(t, ErrorCode(""), result)
	})

	// Test with non-ContextualError
	t.Run("NonContextualError", func(t *testing.T) {
		err := fmt.Errorf("simple error")
		result := GetCode(err)
		assert.Equal(t, ErrorCode(""), result)
	})
}

// TestGetHTTPStatusEdgeCases tests edge cases for the GetHTTPStatus function
func TestGetHTTPStatusEdgeCases(t *testing.T) {
	// Test with nil error
	t.Run("NilError", func(t *testing.T) {
		result := GetHTTPStatus(nil)
		assert.Equal(t, 0, result)
	})

	// Test with non-ContextualError
	t.Run("NonContextualError", func(t *testing.T) {
		err := fmt.Errorf("simple error")
		result := GetHTTPStatus(err)
		assert.Equal(t, 0, result)
	})
}

// TestRepositoryErrorMethods tests the methods of RepositoryError
func TestRepositoryErrorMethods(t *testing.T) {
	// Test Error method with message
	t.Run("ErrorWithMessage", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		err := NewRepositoryError(originalErr, "repository error", "repo_error")
		assert.Equal(t, "repository error", err.Error())
	})

	// Test Error method without message
	t.Run("ErrorWithoutMessage", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		err := NewRepositoryError(originalErr, "", "repo_error")
		assert.Equal(t, "original error", err.Error())
	})

	// Test Unwrap method
	t.Run("Unwrap", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		err := NewRepositoryError(originalErr, "repository error", "repo_error")
		unwrapped := err.Unwrap()
		assert.Equal(t, originalErr, unwrapped)
	})
}

// TestApplicationErrorMethods tests the methods of ApplicationError
func TestApplicationErrorMethods(t *testing.T) {
	// Test Error method with message
	t.Run("ErrorWithMessage", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		err := NewApplicationError(originalErr, "application error", "app_error")
		assert.Equal(t, "application error", err.Error())
	})

	// Test Error method without message
	t.Run("ErrorWithoutMessage", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		err := NewApplicationError(originalErr, "", "app_error")
		assert.Equal(t, "original error", err.Error())
	})

	// Test Unwrap method
	t.Run("Unwrap", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		err := NewApplicationError(originalErr, "application error", "app_error")
		unwrapped := err.Unwrap()
		assert.Equal(t, originalErr, unwrapped)
	})
}

// TestErrorIsMethods tests the Is method of Error
func TestErrorIsMethods(t *testing.T) {
	// Test Is method with nil target
	t.Run("IsWithNilTarget", func(t *testing.T) {
		err := New("op", "code", "message", nil)
		assert.False(t, err.(*Error).Is(nil))
	})

	// Test Is method with non-Error target
	t.Run("IsWithNonErrorTarget", func(t *testing.T) {
		err := New("op", "code", "message", fmt.Errorf("original"))
		target := fmt.Errorf("target")
		assert.False(t, err.(*Error).Is(target))
	})

	// Test Is method with Error target with different code
	t.Run("IsWithDifferentCode", func(t *testing.T) {
		err := New("op", "code1", "message", nil)
		target := New("op", "code2", "message", nil)
		assert.False(t, err.(*Error).Is(target))
	})

	// Test Is method with Error target with same code
	t.Run("IsWithSameCode", func(t *testing.T) {
		err := New("op", "code", "message", nil)
		target := New("op", "code", "message", nil)
		assert.True(t, err.(*Error).Is(target))
	})
}

// TestWrapWithOperationEdgeCases tests edge cases for the WrapWithOperation function
func TestWrapWithOperationEdgeCases(t *testing.T) {
	// Test with nil error
	t.Run("NilError", func(t *testing.T) {
		result := WrapWithOperation(nil, "operation", "format")
		assert.Nil(t, result)
	})
}

// TestDatabaseOperationEdgeCases tests edge cases for the DatabaseOperation function
func TestDatabaseOperationEdgeCases(t *testing.T) {
	// Test with nil error
	t.Run("NilError", func(t *testing.T) {
		err := DatabaseOperation(nil, "failed to query %s", "users")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to query users")
	})
}

// TestInternalEdgeCases tests edge cases for the Internal function
func TestInternalEdgeCases(t *testing.T) {
	// Test with nil error
	t.Run("NilError", func(t *testing.T) {
		err := Internal(nil, "internal error: %s", "service failure")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "internal error: service failure")
	})
}

// TestExternalServiceEdgeCases tests edge cases for the ExternalService function
func TestExternalServiceEdgeCases(t *testing.T) {
	// Test with nil error
	t.Run("NilError", func(t *testing.T) {
		err := ExternalService(nil, "payment-service", "failed to process payment: %s", "timeout")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to process payment: timeout")
	})
}

// TestNetworkEdgeCases tests edge cases for the Network function
func TestNetworkEdgeCases(t *testing.T) {
	// Test with nil error
	t.Run("NilError", func(t *testing.T) {
		err := Network(nil, "network error: %s", "connection reset")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "network error: connection reset")
	})
}

// TestGetHTTPStatusFromErrorEdgeCases tests edge cases for the GetHTTPStatusFromError function
func TestGetHTTPStatusFromErrorEdgeCases(t *testing.T) {
	// Test with nil error
	t.Run("NilError", func(t *testing.T) {
		status := GetHTTPStatusFromError(nil)
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}
