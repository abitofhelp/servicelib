//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of case-insensitive string operations
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/stringutil"
)

func main() {
	// Check if a string starts with a prefix (case-insensitive)
	hasPrefix := stringutil.HasPrefixIgnoreCase("Hello, World!", "hello")
	fmt.Printf("Has prefix 'hello': %v\n", hasPrefix) // Output: true

	// Check if a string contains a substring (case-insensitive)
	contains := stringutil.ContainsIgnoreCase("Hello, World!", "WORLD")
	fmt.Printf("Contains 'WORLD': %v\n", contains) // Output: true

	// Convert to lowercase (wrapper around strings.ToLower)
	lower := stringutil.ToLowerCase("Hello, World!")
	fmt.Printf("Lowercase: %s\n", lower) // Output: hello, world!

	// Examples with different cases
	fmt.Println("\nMore examples:")

	// HasPrefixIgnoreCase examples
	examples := []struct {
		str    string
		prefix string
	}{
		{"UPPERCASE", "upper"},
		{"lowercase", "LOWER"},
		{"MixedCase", "mixed"},
		{"No match", "something"},
	}

	for _, ex := range examples {
		result := stringutil.HasPrefixIgnoreCase(ex.str, ex.prefix)
		fmt.Printf("HasPrefixIgnoreCase(%q, %q) = %v\n", ex.str, ex.prefix, result)
	}

	// ContainsIgnoreCase examples
	phrases := []string{
		"The Quick Brown Fox",
		"A simple example",
		"No matching text here",
	}

	searchTerms := []string{"quick", "EXAMPLE", "missing"}

	for _, phrase := range phrases {
		for _, term := range searchTerms {
			result := stringutil.ContainsIgnoreCase(phrase, term)
			fmt.Printf("ContainsIgnoreCase(%q, %q) = %v\n", phrase, term, result)
		}
	}
}
