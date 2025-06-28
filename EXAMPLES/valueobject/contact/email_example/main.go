//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the Email value object
package main

import (
	"fmt"
	"strings"

	"github.com/abitofhelp/servicelib/valueobject/contact"
)

func main() {
	// Create a new email
	email, err := contact.NewEmail("user@example.com")
	if err != nil {
		// Handle error (invalid email format)
		fmt.Println("Error creating email:", err)
		return
	}

	// Access values
	address := email.String()
	parts := strings.Split(address, "@")
	domain := ""
	if len(parts) > 1 {
		domain = parts[1]
	}
	fmt.Printf("Email address: %s, Domain: %s\n", address, domain)

	// Check if it's a specific domain
	isGmail := domain == "gmail.com"
	fmt.Printf("Is Gmail: %v\n", isGmail)

	// Format as string
	fmt.Println(email.String()) // "user@example.com"
}