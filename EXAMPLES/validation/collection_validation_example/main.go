// Copyright (c) 2025 A Bit of Help, Inc.

// Example of collection validation using the validation package
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/validation"
)

func main() {
	// Create a validation result
	result := validation.NewValidationResult()

	// Validate a slice of strings (all must be non-empty)
	tags := []string{"tag1", "tag2", ""}
	fmt.Println("Validating tags:", tags)

	// Using AllTrue to check if all tags are non-empty
	allNonEmpty := validation.AllTrue(tags, func(tag string) bool {
		return tag != ""
	})

	if !allNonEmpty {
		result.AddError("all tags must be non-empty", "tags")
	}

	// Using ValidateAll to validate each tag individually
	validation.ValidateAll(tags, func(tag string, index int, result *validation.ValidationResult) {
		if tag == "" {
			result.AddError(fmt.Sprintf("tag at index %d is empty", index), "tags")
		}
	}, result)

	// Check if validation passed
	if !result.IsValid() {
		fmt.Printf("Validation failed: %v\n", result.Error())
	} else {
		fmt.Println("Validation passed!")
	}

	// Create a new validation result for a valid example
	validResult := validation.NewValidationResult()

	// Validate a valid slice of strings (all non-empty)
	validTags := []string{"tag1", "tag2", "tag3"}
	fmt.Println("\nValidating valid tags:", validTags)

	// Using AllTrue to check if all tags are non-empty
	allValidNonEmpty := validation.AllTrue(validTags, func(tag string) bool {
		return tag != ""
	})

	if !allValidNonEmpty {
		validResult.AddError("all tags must be non-empty", "validTags")
	}

	// Using ValidateAll to validate each tag individually
	validation.ValidateAll(validTags, func(tag string, index int, result *validation.ValidationResult) {
		if tag == "" {
			result.AddError(fmt.Sprintf("tag at index %d is empty", index), "validTags")
		}
	}, validResult)

	// Validate a slice of numbers (all must be positive)
	numbers := []int{1, 2, 3, -1, 5}
	fmt.Println("Validating numbers:", numbers)

	// Using ValidateAll to validate each number individually
	validation.ValidateAll(numbers, func(num int, index int, result *validation.ValidationResult) {
		if num <= 0 {
			result.AddError(fmt.Sprintf("number at index %d must be positive", index), "numbers")
		}
	}, validResult)

	// Check if validation passed
	if !validResult.IsValid() {
		fmt.Printf("Validation failed: %v\n", validResult.Error())
	} else {
		fmt.Println("Validation passed!")
	}

	// Expected output:
	// Validating tags: [tag1 tag2 ]
	// Validation failed: validation errors: tags: all tags must be non-empty, tags: tag at index 2 is empty
	//
	// Validating valid tags: [tag1 tag2 tag3]
	// Validating numbers: [1 2 3 -1 5]
	// Validation failed: validation errors: numbers: number at index 3 must be positive
}
