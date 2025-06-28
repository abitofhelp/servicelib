# Validation

## Overview

The Validation component provides utilities for validating data in a structured and reusable way. It offers a collection of validation functions for common validation tasks, such as checking required fields, validating string lengths, matching patterns, validating dates, and more. It also provides a ValidationResult type for collecting and managing validation errors.

## Features

- **String Validation**: Validate required fields, minimum/maximum length, and pattern matching
- **Date Validation**: Validate past dates and date ranges
- **Collection Validation**: Validate all items in a collection or check if all items satisfy a condition
- **Structured Results**: Collect validation errors with field-specific error messages
- **Error Integration**: Seamless integration with the errors package for consistent error handling
- **Composable Design**: Build complex validation logic by combining simple validation rules
- **Generic Support**: Type-safe validation with Go generics

## Installation

```bash
go get github.com/abitofhelp/servicelib/validation
```

## Quick Start

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/validation"
)

type User struct {
    Name     string
    Email    string
    Password string
}

func main() {
    user := User{
        Name:     "John Doe",
        Email:    "invalid-email",
        Password: "123",
    }
    
    // Create a validation result
    result := validation.NewValidationResult()
    
    // Validate the user
    validation.Required(user.Name, "name", result)
    validation.Pattern(user.Email, `^[^@]+@[^@]+\.[^@]+$`, "email", result)
    validation.MinLength(user.Password, 8, "password", result)
    
    // Check if validation passed
    if !result.IsValid() {
        fmt.Printf("Validation failed: %v\n", result.Error())
        return
    }
    
    fmt.Println("Validation passed!")
}
```

## API Documentation

### Core Types

The validation package provides a ValidationResult type for collecting and managing validation errors.

#### ValidationResult

ValidationResult holds the result of a validation operation and provides methods for adding errors and checking validity.

```go
type ValidationResult struct {
    // Contains unexported fields
}

// NewValidationResult creates a new ValidationResult
func NewValidationResult() *ValidationResult

// AddError adds an error to the validation result
func (v *ValidationResult) AddError(msg, field string)

// IsValid returns true if there are no validation errors
func (v *ValidationResult) IsValid() bool

// Error returns the validation errors as an error
func (v *ValidationResult) Error() error
```

### Validation Functions

The validation package provides several validation functions for common validation tasks.

#### String Validation

```go
// Required validates that a string is not empty
func Required(value, field string, result *ValidationResult)

// MinLength validates that a string has a minimum length
func MinLength(value string, min int, field string, result *ValidationResult)

// MaxLength validates that a string has a maximum length
func MaxLength(value string, max int, field string, result *ValidationResult)

// Pattern validates that a string matches a regular expression
func Pattern(value, pattern, field string, result *ValidationResult)

// ValidateID validates that an ID is not empty
func ValidateID(id, field string, result *ValidationResult)
```

#### Date Validation

```go
// PastDate validates that a date is in the past
func PastDate(value time.Time, field string, result *ValidationResult)

// ValidDateRange validates that a start date is before an end date
func ValidDateRange(start, end time.Time, startField, endField string, result *ValidationResult)
```

#### Collection Validation

```go
// AllTrue validates that all items in a slice satisfy a predicate
func AllTrue[T any](items []T, predicate func(T) bool) bool

// ValidateAll validates all items in a slice
func ValidateAll[T any](items []T, validator func(T, int, *ValidationResult), result *ValidationResult)
```

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/validation/basic_usage_example.go) - Shows basic usage of the validation package
- [Custom Validation](../EXAMPLES/validation/custom_validation_example.go) - Shows how to create custom validation rules
- [Collection Validation](../EXAMPLES/validation/collection_validation_example.go) - Shows how to validate collections
- [Integration with Forms](../EXAMPLES/validation/form_validation_example.go) - Shows how to validate form data

## Best Practices

1. **Create a Single ValidationResult**: Create a single ValidationResult for each validation operation to collect all errors
2. **Use Descriptive Field Names**: Use descriptive field names that match your API or UI field names
3. **Compose Validation Rules**: Build complex validation logic by combining simple validation rules
4. **Validate Early**: Validate input data as early as possible in your application flow
5. **Return All Errors**: Return all validation errors at once rather than one at a time for better user experience

## Troubleshooting

### Common Issues

#### Validation Errors Not Being Collected

If validation errors are not being collected properly:

- Ensure you're passing the same ValidationResult instance to all validation functions
- Check that you're calling AddError with the correct field name
- Verify that you're checking IsValid() before proceeding with your logic

#### Regular Expression Patterns Not Matching

If your pattern validation is not working as expected:

- Test your regular expression pattern separately to ensure it's correct
- Remember that Go uses RE2 syntax, which doesn't support some advanced features
- Consider using a more specific validation function for common patterns (email, URL, etc.)

## Related Components

- [Errors](../errors/README.md) - The validation package integrates with the errors package for error handling
- [Model](../model/README.md) - The validation package can be used to validate model objects

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.