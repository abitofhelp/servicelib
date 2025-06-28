//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the base value object types
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/abitofhelp/servicelib/valueobject/base"
)

// CustomStringValueObject is a custom string-based value object
type CustomStringValueObject struct {
	base.StringValueObject
}

// NewCustomStringValueObject creates a new CustomStringValueObject
func NewCustomStringValueObject(value string) (CustomStringValueObject, error) {
	// Validate the value
	if strings.TrimSpace(value) == "" {
		return CustomStringValueObject{}, errors.New("value cannot be empty")
	}

	return CustomStringValueObject{
		StringValueObject: base.NewStringValueObject(value),
	}, nil
}

// Equals checks if two CustomStringValueObjects are equal
func (vo CustomStringValueObject) Equals(other CustomStringValueObject) bool {
	return vo.StringValueObject.Equals(other.StringValueObject)
}

// EqualsIgnoreCase checks if two CustomStringValueObjects are equal, ignoring case
func (vo CustomStringValueObject) EqualsIgnoreCase(other CustomStringValueObject) bool {
	return vo.StringValueObject.EqualsIgnoreCase(other.StringValueObject)
}

// CustomStructValueObject is a custom struct-based value object
type CustomStructValueObject struct {
	base.BaseStructValueObject
	name string
	age  int
}

// NewCustomStructValueObject creates a new CustomStructValueObject
func NewCustomStructValueObject(name string, age int) (CustomStructValueObject, error) {
	vo := CustomStructValueObject{
		name: strings.TrimSpace(name),
		age:  age,
	}

	return base.WithValidation(vo, vo.Validate())
}

// Validate implements the Validatable interface
func (vo CustomStructValueObject) Validate() error {
	if vo.name == "" {
		return errors.New("name cannot be empty")
	}
	if vo.age < 0 {
		return errors.New("age cannot be negative")
	}
	return nil
}

// String implements the ValueObject interface
func (vo CustomStructValueObject) String() string {
	return fmt.Sprintf("%s (%d)", vo.name, vo.age)
}

// IsEmpty implements the ValueObject interface
func (vo CustomStructValueObject) IsEmpty() bool {
	return vo.name == "" && vo.age == 0
}

// Equals implements the Equatable interface
func (vo CustomStructValueObject) Equals(other CustomStructValueObject) bool {
	return vo.name == other.name && vo.age == other.age
}

// ToMap implements the StructValueObject interface
func (vo CustomStructValueObject) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name": vo.name,
		"age":  vo.age,
	}
}

func main() {
	// Example with StringValueObject
	fmt.Println("=== String Value Object Example ===")

	// Create a custom string value object
	strVO, err := NewCustomStringValueObject("Hello, World!")
	if err != nil {
		fmt.Println("Error creating string value object:", err)
		return
	}

	// Use the value object
	fmt.Println("Value:", strVO.String())
	fmt.Println("Is empty?", strVO.IsEmpty())

	// Create another value object for comparison
	anotherStrVO, _ := NewCustomStringValueObject("Hello, World!")
	fmt.Println("Equals?", strVO.Equals(anotherStrVO))
	fmt.Println("Equals ignore case?", strVO.EqualsIgnoreCase(anotherStrVO))

	// Example with StructValueObject
	fmt.Println("\n=== Struct Value Object Example ===")

	// Create a custom struct value object
	structVO, err := NewCustomStructValueObject("John Doe", 30)
	if err != nil {
		fmt.Println("Error creating struct value object:", err)
		return
	}

	// Use the value object
	fmt.Println("String representation:", structVO.String())
	fmt.Println("Is empty?", structVO.IsEmpty())

	// Create another value object for comparison
	anotherStructVO, _ := NewCustomStructValueObject("John Doe", 30)
	fmt.Println("Equals?", structVO.Equals(anotherStructVO))

	// Convert to map
	voMap := structVO.ToMap()
	fmt.Printf("As map: %v\n", voMap)

	// JSON marshaling
	jsonBytes, _ := json.Marshal(structVO)
	fmt.Printf("As JSON: %s\n", string(jsonBytes))
}