// Copyright (c) 2025 A Bit of Help, Inc.

// Package stringutil provides additional string manipulation utilities
// beyond what's available in the standard library.
//
// This package offers a collection of helper functions for common string operations
// that are not directly available in the standard library's strings package.
// It includes functions for case-insensitive comparisons, string joining with
// grammatical conventions, whitespace handling, and path normalization.
//
// Key features:
//   - Case-insensitive string operations (HasPrefixIgnoreCase, ContainsIgnoreCase)
//   - Multiple prefix checking (HasAnyPrefix)
//   - Grammatically correct string joining (JoinWithAnd)
//   - Whitespace detection and removal (IsEmpty, IsNotEmpty, RemoveWhitespace)
//   - String truncation with ellipsis (Truncate)
//   - Path normalization (ForwardSlashPath)
//
// Example usage:
//
//	// Case-insensitive prefix check
//	if stringutil.HasPrefixIgnoreCase(userInput, "help") {
//	    // Handle help command
//	}
//
//	// Join items with proper grammar
//	items := []string{"apples", "oranges", "bananas"}
//	fmt.Println("Shopping list: " + stringutil.JoinWithAnd(items, true))
//	// Output: "Shopping list: apples, oranges, and bananas"
//
//	// Normalize file paths
//	path := "C:\\Users\\username\\Documents"
//	normalizedPath := stringutil.ForwardSlashPath(path)
//	// Result: "C:/Users/username/Documents"
package stringutil