package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWithDetailsComprehensive tests the WithDetails function more comprehensively
func TestWithDetailsComprehensive(t *testing.T) {
	// Test with a regular error
	t.Run("RegularError", func(t *testing.T) {
		err := fmt.Errorf("regular error")
		details := map[string]interface{}{
			"key1": "value1",
			"key2": 123,
			"key3": true,
		}
		result := WithDetails(err, details)
		assert.NotNil(t, result)

		// Check that it's a ContextualError
		contextualErr, ok := result.(*ContextualError)
		assert.True(t, ok)

		// Check that the original error is preserved
		assert.Equal(t, err, contextualErr.Original)

		// Check that the details were added
		assert.Equal(t, "value1", contextualErr.Context.Details["key1"])
		assert.Equal(t, 123, contextualErr.Context.Details["key2"])
		assert.Equal(t, true, contextualErr.Context.Details["key3"])
	})

	// Test with an existing ContextualError with no details
	t.Run("ContextualErrorNoDetails", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		contextualErr := &ContextualError{
			Original: originalErr,
			Context: ErrorContext{
				Operation: "TestOperation",
				Source:    "test_file.go",
				Line:      42,
				Code:      "test_code",
			},
		}

		details := map[string]interface{}{
			"key1": "value1",
		}

		result := WithDetails(contextualErr, details)
		assert.NotNil(t, result)

		// Check that it's a ContextualError
		resultErr, ok := result.(*ContextualError)
		assert.True(t, ok)

		// Check that the original error is preserved
		assert.Equal(t, originalErr, resultErr.Original)

		// Check that the operation is preserved
		assert.Equal(t, "TestOperation", resultErr.Context.Operation)

		// Check that the details were added
		assert.Equal(t, "value1", resultErr.Context.Details["key1"])
	})

	// Test with an existing ContextualError with existing details
	t.Run("ContextualErrorWithDetails", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		contextualErr := &ContextualError{
			Original: originalErr,
			Context: ErrorContext{
				Operation: "TestOperation",
				Source:    "test_file.go",
				Line:      42,
				Code:      "test_code",
				Details: map[string]interface{}{
					"existing_key": "existing_value",
				},
			},
		}

		details := map[string]interface{}{
			"new_key": "new_value",
		}

		result := WithDetails(contextualErr, details)
		assert.NotNil(t, result)

		// Check that it's a ContextualError
		resultErr, ok := result.(*ContextualError)
		assert.True(t, ok)

		// Check that the original error is preserved
		assert.Equal(t, originalErr, resultErr.Original)

		// Check that the operation is preserved
		assert.Equal(t, "TestOperation", resultErr.Context.Operation)

		// Check that both the existing and new details are present
		assert.Equal(t, "existing_value", resultErr.Context.Details["existing_key"])
		assert.Equal(t, "new_value", resultErr.Context.Details["new_key"])
	})

	// Test with overlapping keys
	t.Run("OverlappingKeys", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		contextualErr := &ContextualError{
			Original: originalErr,
			Context: ErrorContext{
				Operation: "TestOperation",
				Source:    "test_file.go",
				Line:      42,
				Code:      "test_code",
				Details: map[string]interface{}{
					"key": "old_value",
				},
			},
		}

		details := map[string]interface{}{
			"key": "new_value",
		}

		result := WithDetails(contextualErr, details)
		assert.NotNil(t, result)

		// Check that it's a ContextualError
		resultErr, ok := result.(*ContextualError)
		assert.True(t, ok)

		// Check that the overlapping key has the new value
		assert.Equal(t, "new_value", resultErr.Context.Details["key"])
	})
}

// TestGetHTTPStatusFromErrorComprehensive tests the GetHTTPStatusFromError function more comprehensively
func TestGetHTTPStatusFromErrorComprehensive(t *testing.T) {
	// Test with a custom error that implements ErrorWithHTTPStatus
	t.Run("CustomErrorWithHTTPStatus", func(t *testing.T) {
		// Create a custom error that implements ErrorWithHTTPStatus
		err := &mockErrorWithHTTPStatus{
			msg:    "custom error",
			status: 418,
		}

		status := GetHTTPStatusFromError(err)
		assert.Equal(t, 418, status)
	})
}

// mockErrorWithHTTPStatus is a simple implementation of ErrorWithHTTPStatus for testing
type mockErrorWithHTTPStatus struct {
	msg    string
	status int
}

func (m *mockErrorWithHTTPStatus) Error() string {
	return m.msg
}

func (m *mockErrorWithHTTPStatus) HTTPStatus() int {
	return m.status
}

// TestToJSONComprehensive tests the ToJSON function more comprehensively
func TestToJSONComprehensive(t *testing.T) {
	// Test with a ContextualError that has complex details
	t.Run("ComplexDetails", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		details := map[string]interface{}{
			"string": "value",
			"number": 123,
			"bool":   true,
			"array":  []string{"a", "b", "c"},
			"nested": map[string]interface{}{
				"key": "value",
			},
		}

		err := &ContextualError{
			Original: originalErr,
			Context: ErrorContext{
				Operation:  "TestOperation",
				Source:     "test_file.go",
				Line:       42,
				Code:       "test_code",
				HTTPStatus: 400,
				Details:    details,
			},
		}

		jsonStr := ToJSON(err)
		assert.Contains(t, jsonStr, "original error")
		assert.Contains(t, jsonStr, "TestOperation")
		assert.Contains(t, jsonStr, "test_code")
		assert.Contains(t, jsonStr, "value")
		assert.Contains(t, jsonStr, "123")
		assert.Contains(t, jsonStr, "true")
		assert.Contains(t, jsonStr, "a")
		assert.Contains(t, jsonStr, "b")
		assert.Contains(t, jsonStr, "c")
	})
}

// TestWrapWithOperationComprehensive tests the WrapWithOperation function more comprehensively
func TestWrapWithOperationComprehensive(t *testing.T) {
	// Test with a ContextualError
	t.Run("ContextualError", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		contextualErr := &ContextualError{
			Original: originalErr,
			Context: ErrorContext{
				Operation:  "OriginalOperation",
				Source:     "original_file.go",
				Line:       100,
				Code:       "original_code",
				HTTPStatus: 500,
				Details: map[string]interface{}{
					"key": "value",
				},
			},
		}

		result := WrapWithOperation(contextualErr, "NewOperation", "new format %s", "arg")
		assert.NotNil(t, result)

		// Check that it's a ContextualError
		resultErr, ok := result.(*ContextualError)
		assert.True(t, ok)

		// Check that the original error is preserved
		assert.Equal(t, originalErr, resultErr.Original)

		// Check that the operation is updated
		assert.Equal(t, "NewOperation", resultErr.Context.Operation)

		// Check that the code and HTTP status are preserved
		assert.Equal(t, ErrorCode("original_code"), resultErr.Context.Code)
		assert.Equal(t, 500, resultErr.Context.HTTPStatus)

		// Check that the details are preserved
		assert.Equal(t, "value", resultErr.Context.Details["key"])
	})
}

// TestGetCallerInfoComprehensive tests the getCallerInfo function more comprehensively
func TestGetCallerInfoComprehensive(t *testing.T) {
	// Test with different skip values
	t.Run("DifferentSkipValues", func(t *testing.T) {
		// Skip 0 (current function)
		file0, line0 := getCallerInfo(0)
		assert.Equal(t, "coverage_test.go", file0)
		assert.Greater(t, line0, 0)

		// Skip 1 (parent function)
		file1, line1 := func() (string, int) {
			return getCallerInfo(1)
		}()
		assert.Equal(t, "coverage_test.go", file1)
		assert.Greater(t, line1, 0)

		// Skip 2 (grandparent function)
		file2, line2 := func() (string, int) {
			return func() (string, int) {
				return getCallerInfo(2)
			}()
		}()
		assert.Equal(t, "coverage_test.go", file2)
		assert.Greater(t, line2, 0)
	})

	// Test with an invalid skip value
	t.Run("InvalidSkipValue", func(t *testing.T) {
		file, line := getCallerInfo(100) // Too high, should be out of call stack
		assert.Equal(t, "", file)
		assert.Equal(t, 0, line)
	})
}
