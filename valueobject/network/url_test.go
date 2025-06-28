// Copyright (c) 2025 A Bit of Help, Inc.

package network

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewURL(t *testing.T) {
	tests := []struct {
		name        string
		rawURL      string
		expectError bool
	}{
		{"Valid HTTP URL", "http://example.com", false},
		{"Valid HTTPS URL", "https://example.com", false},
		{"Valid URL with Path", "https://example.com/path", false},
		{"Valid URL with Query", "https://example.com?param=value", false},
		{"Valid URL with Path and Query", "https://example.com/path?param=value", false},
		{"Empty URL", "", false},
		{"URL with Whitespace", "  https://example.com  ", false},
		{"Missing Scheme", "example.com", true},
		{"Invalid Scheme", "ftp://example.com", true},
		{"Missing Host", "http://", true},
		{"Invalid URL Format", "http://[invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewURL(tt.rawURL)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.rawURL != "" {
					assert.Equal(t, strings.TrimSpace(tt.rawURL), u.String())
				}
			}
		})
	}
}

func TestURL_String(t *testing.T) {
	urlStr := "https://example.com/path?param=value"
	u, _ := NewURL(urlStr)

	assert.Equal(t, urlStr, u.String())
}

func TestURL_Equals(t *testing.T) {
	u1, _ := NewURL("https://example.com")
	u2, _ := NewURL("https://example.com")
	u3, _ := NewURL("https://example.org")

	// Test with invalid URLs for the fallback case
	var u4 URL = "https://example.com"
	var u5 URL = "https://[invalid"

	assert.True(t, u1.Equals(u2))
	assert.False(t, u1.Equals(u3))

	// Test the fallback case
	assert.True(t, u4.Equals(u4))
	assert.False(t, u4.Equals(u5))
}

func TestURL_IsEmpty(t *testing.T) {
	emptyURL := URL("")
	u, _ := NewURL("https://example.com")

	assert.True(t, emptyURL.IsEmpty())
	assert.False(t, u.IsEmpty())
}

func TestURL_Domain(t *testing.T) {
	tests := []struct {
		name        string
		rawURL      string
		expected    string
		expectError bool
	}{
		{"Valid URL", "https://example.com", "example.com", false},
		{"URL with Port", "https://example.com:8080", "example.com:8080", false},
		{"Empty URL", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u URL
			var err error

			if tt.rawURL != "" {
				u, err = NewURL(tt.rawURL)
				if err != nil {
					t.Fatalf("Failed to create URL: %v", err)
				}
			}

			domain, err := u.Domain()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, domain)
			}
		})
	}

	// Test with invalid URL format
	var invalidURL URL = "https://[invalid"
	_, err := invalidURL.Domain()
	assert.Error(t, err)
}

func TestURL_Path(t *testing.T) {
	tests := []struct {
		name        string
		rawURL      string
		expected    string
		expectError bool
	}{
		{"URL with Path", "https://example.com/path", "/path", false},
		{"URL without Path", "https://example.com", "", false},
		{"URL with Complex Path", "https://example.com/path/to/resource", "/path/to/resource", false},
		{"Empty URL", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u URL
			var err error

			if tt.rawURL != "" {
				u, err = NewURL(tt.rawURL)
				if err != nil {
					t.Fatalf("Failed to create URL: %v", err)
				}
			}

			path, err := u.Path()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, path)
			}
		})
	}

	// Test with invalid URL format
	var invalidURL URL = "https://[invalid"
	_, err := invalidURL.Path()
	assert.Error(t, err)
}

func TestURL_Query(t *testing.T) {
	tests := []struct {
		name        string
		rawURL      string
		queryParams map[string]string
		expectError bool
	}{
		{
			name:        "URL with Query",
			rawURL:      "https://example.com?param1=value1&param2=value2",
			queryParams: map[string]string{"param1": "value1", "param2": "value2"},
			expectError: false,
		},
		{
			name:        "URL without Query",
			rawURL:      "https://example.com",
			queryParams: map[string]string{},
			expectError: false,
		},
		{
			name:        "Empty URL",
			rawURL:      "",
			queryParams: nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u URL
			var err error

			if tt.rawURL != "" {
				u, err = NewURL(tt.rawURL)
				if err != nil {
					t.Fatalf("Failed to create URL: %v", err)
				}
			}

			query, err := u.Query()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Check that all expected query parameters are present
				for key, expectedValue := range tt.queryParams {
					values := query[key]
					if len(values) == 0 {
						t.Errorf("Expected query parameter %s not found", key)
					} else {
						assert.Equal(t, expectedValue, values[0])
					}
				}

				// Check that there are no unexpected query parameters
				for key := range query {
					_, exists := tt.queryParams[key]
					assert.True(t, exists, "Unexpected query parameter %s", key)
				}
			}
		})
	}

	// Test with invalid URL format
	var invalidURL URL = "https://[invalid"
	_, err := invalidURL.Query()
	assert.Error(t, err)
}

func TestURL_Validate(t *testing.T) {
	tests := []struct {
		name        string
		rawURL      string
		expectError bool
	}{
		{"Valid HTTP URL", "http://example.com", false},
		{"Valid HTTPS URL", "https://example.com", false},
		{"Empty URL", "", false},
		{"Missing Scheme", "example.com", true},
		{"Invalid Scheme", "ftp://example.com", true},
		{"Missing Host", "http://", true},
		{"Invalid URL Format", "http://[invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u URL
			if tt.rawURL != "" {
				u = URL(tt.rawURL)
			}

			err := u.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}