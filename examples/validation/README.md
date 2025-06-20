# Validation Package Examples

This directory contains examples demonstrating how to use the `validation` package, which provides utilities for validating data in Go applications. The package offers a structured way to perform validations, collect validation errors, and report them in a consistent format.

## Examples

### 1. Basic Validation Example

[basic_validation_example.go](basic_validation_example.go)

Demonstrates basic validation functions for strings.

Key concepts:
- Creating a validation result
- Using `Required` to validate that a value is not empty
- Using `MinLength` and `MaxLength` to validate string length
- Using `Pattern` to validate that a string matches a regular expression
- Checking if validation passed with `IsValid()`

### 2. Date Validation Example

[date_validation_example.go](date_validation_example.go)

Shows how to validate dates and date ranges.

Key concepts:
- Using `PastDate` to validate that a date is in the past
- Using `ValidDateRange` to validate that a start date is before an end date
- Handling validation errors for dates

### 3. Collection Validation Example

[collection_validation_example.go](collection_validation_example.go)

Demonstrates validation of collections (slices).

Key concepts:
- Using `AllTrue` to check if all items in a collection satisfy a condition
- Using `ValidateAll` to validate each item in a collection individually
- Handling validation errors for collections
- Validating different types of collections (strings, numbers)

### 4. Custom Validation Example

[custom_validation_example.go](custom_validation_example.go)

Shows how to create custom validation for a struct.

Key concepts:
- Creating a validation function for a custom struct
- Combining multiple validators for complex validation rules
- Using `ValidateID` for ID validation
- Adding custom validation logic
- Returning validation errors as a standard error

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/validation/basic_validation_example.go
```

## Additional Resources

For more information about the validation package, see the [validation package documentation](../../validation/README.md).