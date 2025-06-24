# Model Package Examples

This directory contains examples demonstrating how to use the `model` package, which provides utilities for working with data models in Go applications. The package offers functionality for deep copying objects, handling errors in model operations, and copying fields between structs.

## Examples

### 1. Deep Copy Example

[deep_copy_example.go](deep_copy_example.go)

Demonstrates how to create deep copies of complex objects.

Key concepts:
- Creating deep copies of structs with nested objects
- Handling slices, maps, and pointers during copying
- Verifying independence between original and copied objects
- Maintaining object relationships in copied structures
- Error handling during the deep copy process

### 2. Error Handling Example

[error_handling_example.go](error_handling_example.go)

Shows how to handle errors in model operations.

Key concepts:
- Handling validation errors
- Working with domain-specific errors
- Error wrapping and unwrapping
- Structured error information
- Error classification and handling strategies

### 3. Field Copying Example

[field_copying_example.go](field_copying_example.go)

Demonstrates how to copy fields between structs.

Key concepts:
- Copying fields between different struct types
- Handling field name mapping
- Selective field copying
- Type conversion during copying
- Handling nested structures

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/model/deep_copy_example.go
```

## Additional Resources

For more information about the model package, see the [model package documentation](../../model/README.md).