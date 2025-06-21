//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of copying fields between structs
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/model"
)

// Domain model
type User struct {
	ID        string
	Username  string
	Email     string
	FirstName string
	LastName  string
	Age       int
	Active    bool
	private   string // Unexported field
}

// DTO for API responses
type UserResponse struct {
	ID        string
	Username  string
	Email     string
	FirstName string
	LastName  string
	// Note: Age and Active are not included in the response
}

func main() {
	// Create a domain model instance
	user := &User{
		ID:        "123",
		Username:  "johndoe",
		Email:     "john.doe@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
		Active:    true,
		private:   "sensitive data",
	}

	// Create a DTO instance
	response := &UserResponse{}

	// Copy fields from domain model to DTO
	err := model.CopyFields(response, user)
	if err != nil {
		fmt.Printf("Error copying fields: %v\n", err)
		return
	}

	// Print the DTO
	fmt.Printf("User Response: %+v\n", response)
	// Output: User Response: {ID:123 Username:johndoe Email:john.doe@example.com FirstName:John LastName:Doe}
}
