// Copyright (c) 2025 A Bit of Help, Inc.

// Example of string joining with natural language formatting
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/stringutil"
)

func main() {
	// Join strings with "and" (no Oxford comma)
	fruits := []string{"apples", "bananas", "oranges"}
	joined := stringutil.JoinWithAnd(fruits, false)
	fmt.Printf("Without Oxford comma: %s\n", joined) // Output: apples, bananas and oranges

	// Join strings with "and" (with Oxford comma)
	joinedOxford := stringutil.JoinWithAnd(fruits, true)
	fmt.Printf("With Oxford comma: %s\n", joinedOxford) // Output: apples, bananas, and oranges

	// Handle different list lengths
	fmt.Println("\nDifferent list lengths:")
	fmt.Printf("Empty list: %q\n", stringutil.JoinWithAnd([]string{}, false))                   // Output: ""
	fmt.Printf("One item: %q\n", stringutil.JoinWithAnd([]string{"apples"}, false))             // Output: "apples"
	fmt.Printf("Two items: %q\n", stringutil.JoinWithAnd([]string{"apples", "bananas"}, false)) // Output: "apples and bananas"

	// Real-world examples
	fmt.Println("\nReal-world examples:")

	// Example 1: Listing ingredients
	ingredients := []string{"flour", "sugar", "eggs", "butter", "vanilla extract"}
	fmt.Printf("Ingredients: %s\n", stringutil.JoinWithAnd(ingredients, true))

	// Example 2: Listing participants
	participants := []string{"John", "Alice"}
	fmt.Printf("Participants: %s\n", stringutil.JoinWithAnd(participants, false))

	// Example 3: Listing error messages
	errors := []string{"Connection timeout", "Invalid credentials", "Permission denied"}
	fmt.Printf("Errors encountered: %s\n", stringutil.JoinWithAnd(errors, true))

	// Example 4: Listing options
	options := []string{"Option A"}
	fmt.Printf("Available options: %s\n", stringutil.JoinWithAnd(options, false))
}
