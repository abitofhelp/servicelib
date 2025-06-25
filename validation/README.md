# Validation Package

## Overview

The `validation` package provides utilities for validating data in Go applications. It offers a structured way to perform validations, collect validation errors, and report them in a consistent format. This package is useful for validating user input, request data, or any other data that needs to conform to certain rules or constraints.


## Features

- **Structured Validation**: Collect and report validation errors in a structured way
- **Common Validators**: Built-in validators for common validation scenarios
- **Generic Support**: Type-safe validation with Go generics
- **Composable Validation**: Combine multiple validators for complex validation rules
- **Field-Level Validation**: Validate individual fields with specific error messages
- **Collection Validation**: Validate all items in a collection


## Installation

```bash
go get github.com/abitofhelp/servicelib/validation
```


## Quick Start

See the [Basic Validation Example](../EXAMPLES/validation/basic_validation_example.go) for a complete, runnable example of how to use the validation package.


## Configuration

The validation package does not require any specific configuration. It can be used directly without any setup.


## API Documentation


### Core Types

#### ValidationResult

The `ValidationResult` struct is the main entry point for the validation package. It provides methods for collecting and reporting validation errors.

See the [Basic Validation Example](../EXAMPLES/validation/basic_validation_example.go) for a complete, runnable example of how to use the ValidationResult struct.

#### ValidationError

The `ValidationError` struct represents a validation error with a message and a field name.


### Key Methods

#### NewValidationResult

The `NewValidationResult` function creates a new validation result.

#### AddError

The `AddError` method adds an error to the validation result.

#### IsValid

The `IsValid` method checks if the validation result is valid (has no errors).

#### Error

The `Error` method returns an error if the validation result is not valid.

### Common Validators

The validation package provides several common validators:

- `Required`: Validates that a value is not empty
- `MinLength`: Validates that a string has a minimum length
- `MaxLength`: Validates that a string has a maximum length
- `Pattern`: Validates that a string matches a regular expression pattern
- `PastDate`: Validates that a date is in the past
- `ValidDateRange`: Validates that a date range is valid
- `ValidateID`: Validates that an ID is valid
- `AllTrue`: Validates that all items in a collection satisfy a condition
- `ValidateAll`: Validates all items in a collection

See the [Date Validation Example](../EXAMPLES/validation/date_validation_example.go) for a complete, runnable example of how to use date validators.


## Examples

For complete, runnable examples, see the following files in the examples directory:

- [Basic Validation Example](../EXAMPLES/validation/basic_validation_example.go) - Shows how to use basic validators
- [Collection Validation Example](../EXAMPLES/validation/collection_validation_example.go) - Shows how to validate collections
- [Custom Validation Example](../EXAMPLES/validation/custom_validation_example.go) - Shows how to create custom validators
- [Date Validation Example](../EXAMPLES/validation/date_validation_example.go) - Shows how to validate dates


## Best Practices

1. **Create a Validation Result Early**: Create a validation result at the beginning of your validation function and pass it to all validators.

2. **Use Field Names Consistently**: Use consistent field names across your application to make error messages more understandable.

3. **Combine Validators**: Combine multiple validators for complex validation rules.

4. **Validate Collections Properly**: Use ValidateAll for collections where each item needs validation.

5. **Return Early for Critical Validations**: For performance reasons, return early if critical validations fail.

See the [Custom Validation Example](../EXAMPLES/validation/custom_validation_example.go) for a complete, runnable example of how to implement these best practices.


## Troubleshooting

### Common Issues

#### Validation Errors Not Being Reported

**Issue**: Validation errors are not being reported even though validation is failing.

**Solution**: Ensure that you are checking the result of `IsValid()` or using the `Error()` method to get the validation errors.

#### Inconsistent Error Messages

**Issue**: Error messages are inconsistent or not descriptive enough.

**Solution**: Use consistent field names and descriptive error messages. Consider creating custom validators for complex validation rules.

#### Performance Issues with Large Collections

**Issue**: Validation is slow for large collections.

**Solution**: Consider validating only a subset of the collection or implementing pagination. For critical validations, return early if they fail.


## Related Components

- [Errors](../errors/README.md) - The errors component is used by the validation package for error handling.
- [Context](../context/README.md) - The context component can be used with the validation package for context-aware validation.
- [Logging](../logging/README.md) - The logging component can be used with the validation package for logging validation errors.


## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.


## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
