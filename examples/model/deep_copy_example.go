//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of creating deep copies of objects
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/model"
)

// Complex struct with nested objects
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
	Country string
}

type Person struct {
	ID        string
	Name      string
	Age       int
	Addresses []Address
	Metadata  map[string]string
	Parent    *Person
}

func main() {
	// Create a complex object
	parent := &Person{
		ID:   "parent123",
		Name: "Parent",
		Age:  55,
	}

	original := &Person{
		ID:   "person123",
		Name: "John Doe",
		Age:  30,
		Addresses: []Address{
			{
				Street:  "123 Main St",
				City:    "Anytown",
				State:   "CA",
				ZipCode: "12345",
				Country: "USA",
			},
			{
				Street:  "456 Oak Ave",
				City:    "Othertown",
				State:   "NY",
				ZipCode: "67890",
				Country: "USA",
			},
		},
		Metadata: map[string]string{
			"created": "2023-01-15",
			"source":  "registration",
		},
		Parent: parent,
	}

	// Create a deep copy
	copy := &Person{}
	err := model.DeepCopy(copy, original)
	if err != nil {
		fmt.Printf("Error creating deep copy: %v\n", err)
		return
	}

	// Verify the copy is independent from the original
	fmt.Printf("Original: %+v\n", original)
	fmt.Printf("Copy: %+v\n", copy)

	// Modify the copy
	copy.Name = "Jane Doe"
	copy.Age = 28
	copy.Addresses[0].Street = "789 Pine St"
	copy.Metadata["updated"] = "2023-02-20"

	// Verify the original is unchanged
	fmt.Printf("Original after modification: %+v\n", original)
	fmt.Printf("Copy after modification: %+v\n", copy)
}
