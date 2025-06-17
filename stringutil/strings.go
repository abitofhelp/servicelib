// Copyright (c) 2025 A Bit of Help, Inc.

// Package stringutil provides additional string manipulation utilities
// beyond what's available in the standard library.
package stringutil

import (
	"regexp"
	"strings"
)

// HasPrefixIgnoreCase checks if the string s begins with the specified prefix, ignoring case.
// It is safe for UTF-8 encoded strings and performs a case-insensitive comparison.
// Parameters:
//   - s: The string to check
//   - prefix: The prefix to look for at the beginning of s
//
// Returns:
//   - bool: true if s starts with prefix (ignoring case), false otherwise
func HasPrefixIgnoreCase(s, prefix string) bool {
	r := len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
	return r
}

// ContainsIgnoreCase checks if substr is within s, ignoring case.
// Both s and substr are converted to lowercase before comparison.
// Parameters:
//   - s: The string to search in
//   - substr: The substring to search for
//
// Returns:
//   - bool: true if substr is found in s (ignoring case), false otherwise
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// HasAnyPrefix checks if the string s begins with any of the specified prefixes.
// It iterates through the provided prefixes and returns as soon as a match is found.
// Parameters:
//   - s: The string to check
//   - prefixes: One or more prefixes to look for at the beginning of s
//
// Returns:
//   - bool: true if s starts with any of the prefixes, false otherwise
func HasAnyPrefix(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// ToLowerCase converts a string to lowercase.
// It's a simple wrapper around strings.ToLower for consistency within the package.
// Parameters:
//   - s: The string to convert to lowercase
//
// Returns:
//   - string: The lowercase version of the input string
func ToLowerCase(s string) string {
	return strings.ToLower(s)
}

// JoinWithAnd joins a slice of strings with commas and "and".
// It handles different list lengths appropriately and supports Oxford comma usage.
// Parameters:
//   - items: The slice of strings to join
//   - useOxfordComma: Whether to include a comma before "and" for lists of 3 or more items
//
// Returns:
//   - string: The joined string, or an empty string if items is empty
func JoinWithAnd(items []string, useOxfordComma bool) string {
	length := len(items)

	if length == 0 {
		return ""
	}

	if length == 1 {
		return items[0]
	}

	if length == 2 {
		return items[0] + " and " + items[1]
	}

	// For 3 or more items
	var result string
	for i := 0; i < length-1; i++ {
		result += items[i] + ", "
	}

	// Remove the trailing comma and space
	result = result[:len(result)-2]

	// Add the final part with or without Oxford comma
	if useOxfordComma {
		result += ", and " + items[length-1]
	} else {
		result += " and " + items[length-1]
	}

	return result
}

// IsEmpty checks if a string is empty or contains only whitespace.
// It uses strings.TrimSpace to remove all leading and trailing whitespace.
// Parameters:
//   - s: The string to check
//
// Returns:
//   - bool: true if the string is empty or contains only whitespace, false otherwise
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty checks if a string is not empty and contains non-whitespace characters.
// It's the logical opposite of IsEmpty.
// Parameters:
//   - s: The string to check
//
// Returns:
//   - bool: true if the string contains non-whitespace characters, false otherwise
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// Truncate truncates a string to a specified maximum length.
// If the string is longer than maxLength, it adds "..." to the end.
// If maxLength is negative, it's treated as 0.
// Parameters:
//   - s: The string to truncate
//   - maxLength: The maximum length of the returned string before adding "..."
//
// Returns:
//   - string: The truncated string, or the original string if it's not longer than maxLength
func Truncate(s string, maxLength int) string {
	// Handle negative maxLength
	if maxLength < 0 {
		maxLength = 0
	}

	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// RemoveWhitespace removes all whitespace characters from a string.
// It uses a regular expression to match and remove spaces, tabs, newlines, and other whitespace.
// Parameters:
//   - s: The string to process
//
// Returns:
//   - string: The string with all whitespace characters removed
func RemoveWhitespace(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, "")
}

// ForwardSlashPath converts backslashes in a path to forward slashes.
// This is useful for normalizing file paths across different operating systems.
// Parameters:
//   - path: The file path to normalize
//
// Returns:
//   - string: The path with all backslashes replaced by forward slashes
func ForwardSlashPath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
