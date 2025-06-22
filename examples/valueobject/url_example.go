// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the URL value object
package example_valueobject

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/network"
)

func main() {
	// Create a new URL
	url, err := network.NewURL("https://example.com/path?query=value")
	if err != nil {
		// Handle error
		fmt.Println("Error creating URL:", err)
		return
	}

	// Access URL components
	domain, _ := url.Domain()
	path, _ := url.Path()
	query, _ := url.Query()

	fmt.Printf("URL: %s\n", url.String())
	fmt.Printf("Domain: %s\n", domain)
	fmt.Printf("Path: %s\n", path)
	fmt.Printf("Query parameter 'query': %s\n", query.Get("query"))

	// Compare URLs
	otherURL, _ := network.NewURL("https://example.com/path?query=value")
	areEqual := url.Equals(otherURL)
	fmt.Printf("Are URLs equal? %v\n", areEqual) // true

	// Check if URL is empty
	isEmpty := url.IsEmpty()
	fmt.Printf("Is URL empty? %v\n", isEmpty) // false
}
