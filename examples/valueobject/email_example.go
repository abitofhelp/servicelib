// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the Email value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new email
	email, err := valueobject.NewEmail("user@example.com")
	if err != nil {
		// Handle error (invalid email format)
		fmt.Println("Error creating email:", err)
		return
	}

	// Access values
	address := email.Address()
	domain := email.Domain()
	fmt.Printf("Email address: %s, Domain: %s\n", address, domain)

	// Check if it's a specific domain
	isGmail := email.IsDomain("gmail.com")
	fmt.Printf("Is Gmail: %v\n", isGmail)

	// Format as string
	fmt.Println(email.String()) // "user@example.com"
}
