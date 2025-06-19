// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"strings"
	"testing"
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
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.rawURL != "" && u.String() != strings.TrimSpace(tt.rawURL) {
					t.Errorf("Expected URL %s, got %s", strings.TrimSpace(tt.rawURL), u.String())
				}
			}
		})
	}
}

func TestURL_String(t *testing.T) {
	urlStr := "https://example.com/path?param=value"
	u, _ := NewURL(urlStr)

	if u.String() != urlStr {
		t.Errorf("Expected string %s, got %s", urlStr, u.String())
	}
}

func TestURL_Equals(t *testing.T) {
	u1, _ := NewURL("https://example.com")
	u2, _ := NewURL("https://example.com")
	u3, _ := NewURL("https://example.org")

	// Test with invalid URLs for the fallback case
	var u4 URL = "https://example.com"
	var u5 URL = "https://[invalid"

	if !u1.Equals(u2) {
		t.Errorf("Expected u1 to equal u2")
	}

	if u1.Equals(u3) {
		t.Errorf("Expected u1 to not equal u3")
	}

	// Test the fallback case
	if !u4.Equals(u4) {
		t.Errorf("Expected u4 to equal itself")
	}

	if u4.Equals(u5) {
		t.Errorf("Expected u4 to not equal u5")
	}
}

func TestURL_IsEmpty(t *testing.T) {
	emptyURL := URL("")
	u, _ := NewURL("https://example.com")

	if !emptyURL.IsEmpty() {
		t.Errorf("Expected empty URL to be empty")
	}

	if u.IsEmpty() {
		t.Errorf("Expected non-empty URL to not be empty")
	}
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
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if domain != tt.expected {
					t.Errorf("Expected domain %s, got %s", tt.expected, domain)
				}
			}
		})
	}

	// Test with invalid URL format
	var invalidURL URL = "https://[invalid"
	_, err := invalidURL.Domain()
	if err == nil {
		t.Errorf("Expected error with invalid URL format")
	}
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
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if path != tt.expected {
					t.Errorf("Expected path %s, got %s", tt.expected, path)
				}
			}
		})
	}

	// Test with invalid URL format
	var invalidURL URL = "https://[invalid"
	_, err := invalidURL.Path()
	if err == nil {
		t.Errorf("Expected error with invalid URL format")
	}
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
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Check that all expected query parameters are present
				for key, expectedValue := range tt.queryParams {
					values := query[key]
					if len(values) == 0 {
						t.Errorf("Expected query parameter %s not found", key)
					} else if values[0] != expectedValue {
						t.Errorf("Expected query parameter %s to be %s, got %s", key, expectedValue, values[0])
					}
				}

				// Check that there are no unexpected query parameters
				for key := range query {
					if _, exists := tt.queryParams[key]; !exists {
						t.Errorf("Unexpected query parameter %s", key)
					}
				}
			}
		})
	}

	// Test with invalid URL format
	var invalidURL URL = "https://[invalid"
	_, err := invalidURL.Query()
	if err == nil {
		t.Errorf("Expected error with invalid URL format")
	}
}
