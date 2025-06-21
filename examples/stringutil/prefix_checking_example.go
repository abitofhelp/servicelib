//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of prefix checking operations
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/stringutil"
)

func main() {
	// Check if a string starts with any of the given prefixes
	hasAnyPrefix := stringutil.HasAnyPrefix("https://example.com", "http://", "https://")
	fmt.Printf("Has any prefix: %v\n", hasAnyPrefix) // Output: true

	// Check multiple strings for prefixes
	urls := []string{
		"https://example.com",
		"http://example.org",
		"ftp://example.net",
		"example.io",
	}

	for _, url := range urls {
		if stringutil.HasAnyPrefix(url, "http://", "https://") {
			fmt.Printf("%s is a web URL\n", url)
		} else {
			fmt.Printf("%s is not a web URL\n", url)
		}
	}

	// More complex example with file extensions
	filenames := []string{
		"document.pdf",
		"image.jpg",
		"image.png",
		"spreadsheet.xlsx",
		"text.txt",
		"archive.zip",
	}

	// Group files by type
	var documents, images, others []string

	for _, filename := range filenames {
		switch {
		case stringutil.HasAnyPrefix(filename, "document."):
			documents = append(documents, filename)
		case stringutil.HasAnyPrefix(filename, "image."):
			images = append(images, filename)
		default:
			others = append(others, filename)
		}
	}

	fmt.Println("\nGrouped files:")
	fmt.Printf("Documents: %v\n", documents)
	fmt.Printf("Images: %v\n", images)
	fmt.Printf("Others: %v\n", others)

	// Example with empty prefixes
	emptyResult := stringutil.HasAnyPrefix("test string")
	fmt.Printf("\nHasAnyPrefix with no prefixes: %v\n", emptyResult) // Output: false

	// Example with empty string
	emptyStringResult := stringutil.HasAnyPrefix("", "prefix")
	fmt.Printf("Empty string HasAnyPrefix: %v\n", emptyStringResult) // Output: false
}
