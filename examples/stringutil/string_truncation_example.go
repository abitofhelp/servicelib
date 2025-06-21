//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of string truncation operations
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/stringutil"
)

func main() {
	// Truncate a long string
	longText := "This is a very long string that needs to be truncated to fit in a limited space."
	truncated := stringutil.Truncate(longText, 20)
	fmt.Printf("Original: %q\n", longText)
	fmt.Printf("Truncated to 20 chars: %q\n", truncated) // Output: "This is a very long..."

	// No truncation needed for short strings
	shortText := "Short text"
	truncatedShort := stringutil.Truncate(shortText, 20)
	fmt.Printf("\nOriginal: %q\n", shortText)
	fmt.Printf("Truncated to 20 chars: %q\n", truncatedShort) // Output: "Short text"

	// Handle negative max length
	negativeMax := stringutil.Truncate(longText, -5)
	fmt.Printf("\nNegative max length: %q\n", negativeMax) // Output: "..."

	// Handle zero max length
	zeroMax := stringutil.Truncate(longText, 0)
	fmt.Printf("Zero max length: %q\n", zeroMax) // Output: "..."

	// Various truncation lengths
	fmt.Println("\nVarious truncation lengths:")
	for _, length := range []int{5, 10, 15, 30, 50, 100} {
		result := stringutil.Truncate(longText, length)
		fmt.Printf("Length %2d: %q\n", length, result)
	}

	// Real-world examples
	fmt.Println("\nReal-world examples:")

	// Example 1: Truncating a blog post preview
	blogPost := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
	preview := stringutil.Truncate(blogPost, 50)
	fmt.Printf("Blog preview: %s\n", preview)

	// Example 2: Truncating a filename
	longFilename := "very_long_filename_with_detailed_description_of_contents_and_version_number_v2.1.txt"
	truncatedFilename := stringutil.Truncate(longFilename, 25)
	fmt.Printf("Truncated filename: %s\n", truncatedFilename)

	// Example 3: Truncating a user comment
	userComment := "I really enjoyed this product! It exceeded all my expectations and I would definitely recommend it to others."
	displayComment := stringutil.Truncate(userComment, 35)
	fmt.Printf("User comment: %s\n", displayComment)
}
