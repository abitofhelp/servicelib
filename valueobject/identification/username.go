// Copyright (c) 2025 A Bit of Help, Inc.

// Package identification provides value objects related to identification.
package identification

import (
	"errors"
	"regexp"
	"strings"
)

// Username represents a username value object
type Username string

// Regular expression for validating username format
// Allows alphanumeric characters, underscores, hyphens, and periods
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)

// NewUsername creates a new Username with validation
func NewUsername(username string) (Username, error) {
	// Trim whitespace
	trimmedUsername := strings.TrimSpace(username)

	// Username is required
	if trimmedUsername == "" {
		return "", errors.New("username cannot be empty")
	}

	// Check minimum length
	if len(trimmedUsername) < 3 {
		return "", errors.New("username must be at least 3 characters long")
	}

	// Check maximum length
	if len(trimmedUsername) > 30 {
		return "", errors.New("username cannot exceed 30 characters")
	}

	// Validate username format
	if !usernameRegex.MatchString(trimmedUsername) {
		return "", errors.New("username can only contain letters, numbers, underscores, hyphens, and periods")
	}

	// Ensure username doesn't start or end with a special character
	if strings.HasPrefix(trimmedUsername, "_") || strings.HasPrefix(trimmedUsername, "-") || strings.HasPrefix(trimmedUsername, ".") ||
		strings.HasSuffix(trimmedUsername, "_") || strings.HasSuffix(trimmedUsername, "-") || strings.HasSuffix(trimmedUsername, ".") {
		return "", errors.New("username cannot start or end with an underscore, hyphen, or period")
	}

	// Ensure username doesn't contain consecutive special characters
	if strings.Contains(trimmedUsername, "__") || strings.Contains(trimmedUsername, "--") ||
		strings.Contains(trimmedUsername, "..") || strings.Contains(trimmedUsername, "_.") ||
		strings.Contains(trimmedUsername, "_-") || strings.Contains(trimmedUsername, ".-") ||
		strings.Contains(trimmedUsername, "._") {
		return "", errors.New("username cannot contain consecutive special characters")
	}

	return Username(trimmedUsername), nil
}

// String returns the string representation of the Username
func (u Username) String() string {
	return string(u)
}

// Equals checks if two Usernames are equal (case insensitive)
func (u Username) Equals(other Username) bool {
	return strings.EqualFold(string(u), string(other))
}

// IsEmpty checks if the Username is empty
func (u Username) IsEmpty() bool {
	return u == ""
}

// ToLower returns the lowercase version of the username
func (u Username) ToLower() Username {
	return Username(strings.ToLower(string(u)))
}

// Length returns the length of the username
func (u Username) Length() int {
	return len(string(u))
}

// ContainsSubstring checks if the username contains the given substring (case insensitive)
func (u Username) ContainsSubstring(substring string) bool {
	return strings.Contains(
		strings.ToLower(string(u)),
		strings.ToLower(substring),
	)
}
