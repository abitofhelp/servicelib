// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the ID value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/identification"
)

func main() {
	// Generate a new random ID
	id := identification.GenerateID()
	fmt.Println("Generated ID:", id.String())

	// Create an ID from a string
	idStr := "123e4567-e89b-12d3-a456-426614174000"
	parsedID, err := identification.NewID(idStr)
	if err != nil {
		// Handle error
		fmt.Println("Error creating ID:", err)
		return
	}

	// Check if IDs are equal
	areEqual := id.Equals(parsedID)
	fmt.Printf("Are IDs equal? %v\n", areEqual)

	// Check if ID is empty
	isEmpty := id.IsEmpty()
	fmt.Printf("Is ID empty? %v\n", isEmpty)
}
