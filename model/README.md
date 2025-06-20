# Model Module

The Model Module provides utilities for working with domain models and Data Transfer Objects (DTOs) in Go applications. It includes functions for copying fields between structs and creating deep copies of objects.

## Features

- **Field Copying**: Copy fields between structs based on field names
- **Deep Copying**: Create complete deep copies of objects, including nested structures
- **Reflection-Based**: Works with any struct types without requiring specific interfaces
- **Type Safety**: Ensures type compatibility when copying fields
- **Comprehensive Support**: Handles various types including pointers, structs, slices, and maps

## Installation

```bash
go get github.com/abitofhelp/servicelib/model
```

## Quick Start

See the [Field Copying example](../examples/model/field_copying_example.go) for a complete, runnable example of how to use the Model module.

## API Documentation

### Field Copying

The `CopyFields` function copies fields from source to destination based on field names.

#### Copying Fields Between Structs

See the [Field Copying example](../examples/model/field_copying_example.go) for a complete, runnable example of how to copy fields between structs.

### Deep Copying

The `DeepCopy` function creates a deep copy of the source object.

#### Creating Deep Copies of Objects

See the [Deep Copy example](../examples/model/deep_copy_example.go) for a complete, runnable example of how to create deep copies of objects.

### Error Handling

The model package provides comprehensive error handling for various error scenarios.

#### Error Handling Examples

See the [Error Handling example](../examples/model/error_handling_example.go) for a complete, runnable example of how to handle errors when using the model package.

## Best Practices

1. **Type Safety**: Always check for errors when copying fields or creating deep copies.

2. **Performance Considerations**: Reflection-based operations are slower than direct assignments. Use these utilities for convenience in non-performance-critical code paths.

3. **DTO Mapping**: Use CopyFields for mapping between domain models and DTOs to avoid manual field-by-field copying.

4. **Immutability**: Use DeepCopy to create immutable copies of objects when needed.

5. **Field Naming**: Keep field names consistent between related structs to maximize the effectiveness of CopyFields.

6. **Unexported Fields**: Be aware that unexported (lowercase) fields will be skipped during copying.

7. **Complex Objects**: For very complex objects with custom copying logic, consider implementing your own copy methods instead of relying on reflection.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
