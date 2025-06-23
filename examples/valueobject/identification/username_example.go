// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the Username value object
package identification

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/identification"
)

func main() {
	// Create a new username
	username, err := identification.NewUsername("john.doe")
	if err != nil {
		// Handle error
		fmt.Println("Error creating username:", err)
		return
	}

	// Access username properties
	fmt.Printf("Username: %s\n", username.String())
	fmt.Printf("Length: %d\n", username.Length())

	// Convert to lowercase
	lowerUsername := username.ToLower()
	fmt.Printf("Lowercase: %s\n", lowerUsername)

	// Check if username contains a substring
	containsJohn := username.ContainsSubstring("john")
	fmt.Printf("Contains 'john'? %v\n", containsJohn) // true

	// Compare usernames (case insensitive)
	otherUsername, _ := identification.NewUsername("John.Doe")
	areEqual := username.Equals(otherUsername)
	fmt.Printf("Are usernames equal? %v\n", areEqual) // true
}