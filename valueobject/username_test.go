// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"testing"
)

func TestNewUsername(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		expectError bool
		errorMsg    string
	}{
		{"Valid Username", "johndoe", false, ""},
		{"Valid Username with Numbers", "john123", false, ""},
		{"Valid Username with Underscore", "john_doe", false, ""},
		{"Valid Username with Hyphen", "john-doe", false, ""},
		{"Valid Username with Period", "john.doe", false, ""},
		{"Valid Username with Mixed Characters", "john.doe_123-xyz", false, ""},
		{"Empty Username", "", true, "username cannot be empty"},
		{"Username with Spaces", "  ", true, "username cannot be empty"},
		{"Username Too Short", "jo", true, "username must be at least 3 characters long"},
		{"Username Too Long", "abcdefghijklmnopqrstuvwxyz12345", true, "username cannot exceed 30 characters"},
		{"Username with Invalid Characters", "john@doe", true, "username can only contain letters, numbers, underscores, hyphens, and periods"},
		{"Username with Space", "john doe", true, "username can only contain letters, numbers, underscores, hyphens, and periods"},
		{"Username Starting with Underscore", "_johndoe", true, "username cannot start or end with an underscore, hyphen, or period"},
		{"Username Ending with Underscore", "johndoe_", true, "username cannot start or end with an underscore, hyphen, or period"},
		{"Username Starting with Hyphen", "-johndoe", true, "username cannot start or end with an underscore, hyphen, or period"},
		{"Username Ending with Hyphen", "johndoe-", true, "username cannot start or end with an underscore, hyphen, or period"},
		{"Username Starting with Period", ".johndoe", true, "username cannot start or end with an underscore, hyphen, or period"},
		{"Username Ending with Period", "johndoe.", true, "username cannot start or end with an underscore, hyphen, or period"},
		{"Username with Consecutive Underscores", "john__doe", true, "username cannot contain consecutive special characters"},
		{"Username with Consecutive Hyphens", "john--doe", true, "username cannot contain consecutive special characters"},
		{"Username with Consecutive Periods", "john..doe", true, "username cannot contain consecutive special characters"},
		{"Username with Consecutive Mixed Specials (._)", "john._doe", true, "username cannot contain consecutive special characters"},
		{"Username with Consecutive Mixed Specials (_-)", "john_-doe", true, "username cannot contain consecutive special characters"},
		{"Username with Consecutive Mixed Specials (.-)", "john.-doe", true, "username cannot contain consecutive special characters"},
		{"Username with Whitespace", "john doe", true, "username can only contain letters, numbers, underscores, hyphens, and periods"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, err := NewUsername(tt.username)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if username.String() != tt.username {
					t.Errorf("Expected username %s, got %s", tt.username, username.String())
				}
			}
		})
	}
}

func TestUsername_String(t *testing.T) {
	username, _ := NewUsername("johndoe")

	if username.String() != "johndoe" {
		t.Errorf("Expected string 'johndoe', got '%s'", username.String())
	}
}

func TestUsername_Equals(t *testing.T) {
	username1, _ := NewUsername("johndoe")
	username2, _ := NewUsername("JohnDoe")
	username3, _ := NewUsername("janedoe")

	// Test case insensitive equality
	if !username1.Equals(username2) {
		t.Errorf("Expected 'johndoe' to equal 'JohnDoe' (case insensitive)")
	}

	if username1.Equals(username3) {
		t.Errorf("Expected 'johndoe' to not equal 'janedoe'")
	}
}

func TestUsername_IsEmpty(t *testing.T) {
	// We can't create an empty Username with NewUsername, so we'll create it directly
	var emptyUsername Username = ""
	username, _ := NewUsername("johndoe")

	if !emptyUsername.IsEmpty() {
		t.Errorf("Expected empty username to be empty")
	}

	if username.IsEmpty() {
		t.Errorf("Expected non-empty username to not be empty")
	}
}

func TestUsername_ToLower(t *testing.T) {
	username, _ := NewUsername("JohnDoe")
	lowercased := username.ToLower()

	if lowercased.String() != "johndoe" {
		t.Errorf("Expected lowercase 'johndoe', got '%s'", lowercased.String())
	}

	// Original should be unchanged
	if username.String() != "JohnDoe" {
		t.Errorf("Expected original 'JohnDoe' to be unchanged, got '%s'", username.String())
	}
}

func TestUsername_Length(t *testing.T) {
	tests := []struct {
		username string
		expected int
	}{
		{"johndoe", 7},
		{"a.b-c_d", 7},
		{"user123", 7},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			username, _ := NewUsername(tt.username)
			length := username.Length()

			if length != tt.expected {
				t.Errorf("Expected length %d, got %d", tt.expected, length)
			}
		})
	}
}

func TestUsername_ContainsSubstring(t *testing.T) {
	username, _ := NewUsername("JohnDoe")

	tests := []struct {
		substring string
		expected  bool
	}{
		{"john", true},
		{"doe", true},
		{"JOHN", true},
		{"DOE", true},
		{"hnD", true},
		{"xyz", false},
		{"", true}, // Empty string is contained in any string
	}

	for _, tt := range tests {
		t.Run(tt.substring, func(t *testing.T) {
			contains := username.ContainsSubstring(tt.substring)

			if contains != tt.expected {
				t.Errorf("Expected ContainsSubstring('%s') to be %v, got %v", 
					tt.substring, tt.expected, contains)
			}
		})
	}
}

func TestUsername_Trimming(t *testing.T) {
	// Test that whitespace is trimmed
	username, err := NewUsername("  johndoe  ")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if username.String() != "johndoe" {
		t.Errorf("Expected whitespace to be trimmed, got '%s'", username.String())
	}
}
