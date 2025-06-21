package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		key      string
		value    string
		fallback string
		setup    func()
		cleanup  func()
		expected string
	}{
		{
			name:     "environment variable is set",
			key:      "TEST_ENV_VAR",
			value:    "test_value",
			fallback: "fallback_value",
			setup: func() {
				os.Setenv("TEST_ENV_VAR", "test_value")
			},
			cleanup: func() {
				os.Unsetenv("TEST_ENV_VAR")
			},
			expected: "test_value",
		},
		{
			name:     "environment variable is not set",
			key:      "NONEXISTENT_ENV_VAR",
			fallback: "fallback_value",
			setup: func() {
				os.Unsetenv("NONEXISTENT_ENV_VAR")
			},
			cleanup: func() {
				// No cleanup needed
			},
			expected: "fallback_value",
		},
		{
			name:     "environment variable is set to empty string",
			key:      "EMPTY_ENV_VAR",
			value:    "",
			fallback: "fallback_value",
			setup: func() {
				os.Setenv("EMPTY_ENV_VAR", "")
			},
			cleanup: func() {
				os.Unsetenv("EMPTY_ENV_VAR")
			},
			expected: "",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setup != nil {
				tt.setup()
			}

			// Cleanup
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()

			// Test
			result := GetEnv(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}
}
