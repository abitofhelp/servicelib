# Model

## Overview

The Model component provides utilities for working with domain models and Data Transfer Objects (DTOs). It includes functions for copying fields between structs and creating deep copies of objects, which are essential operations when working with domain models in a clean architecture.

## Features

- **Field Copying**: Copy fields between structs based on field names with type safety
- **Deep Object Copying**: Create complete deep copies of objects, including nested structs, slices, and maps
- **Type Compatibility Checking**: Automatically verify type compatibility during field copying
- **Reflection-Based Operations**: Utilize Go's reflection capabilities for dynamic object manipulation
- **Pointer Handling**: Properly handle nil pointers and create new instances during deep copying

## Installation

```bash
go get github.com/abitofhelp/servicelib/model
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/model"
)

func main() {
    // Define source and destination structs
    type Source struct {
        Name    string
        Age     int
        Address string
    }
    
    type Destination struct {
        Name    string
        Age     int
        Address string
        Extra   string
    }
    
    // Create instances
    src := &Source{
        Name:    "John Doe",
        Age:     30,
        Address: "123 Main St",
    }
    dst := &Destination{}
    
    // Copy fields from source to destination
    err := model.CopyFields(dst, src)
    if err != nil {
        fmt.Printf("Error copying fields: %v\n", err)
        return
    }
    
    // Print the result
    fmt.Printf("Destination after copy: %+v\n", dst)
}
```

## API Documentation

### Core Types

The model package primarily provides utility functions rather than types.

### Key Methods

#### CopyFields

Copies fields from source to destination based on field names.

```go
func CopyFields(dst, src interface{}) error
```

Both source and destination must be pointers to structs. Fields are copied only if they have the same name and compatible types.

#### DeepCopy

Creates a deep copy of the source object.

```go
func DeepCopy(dst, src interface{}) error
```

Both source and destination must be pointers to structs of the same type. This function creates a complete copy of the object, including all nested structs, slices, and maps.

## Examples

Currently, there are no dedicated examples for the model package in the EXAMPLES directory. The following code snippets demonstrate common usage patterns:

### Field Copying Example

```go
// Define source and destination structs
type User struct {
    ID    int
    Name  string
    Email string
}

type UserDTO struct {
    ID    int
    Name  string
    Email string
    Role  string
}

// Create instances
user := &User{ID: 1, Name: "John Doe", Email: "john@example.com"}
userDTO := &UserDTO{}

// Copy fields from user to userDTO
err := model.CopyFields(userDTO, user)
if err != nil {
    // Handle error
}

// userDTO now contains: {ID: 1, Name: "John Doe", Email: "john@example.com", Role: ""}
```

### Deep Copy Example

```go
// Define a struct with nested elements
type Department struct {
    Name     string
    Location string
}

type Employee struct {
    ID         int
    Name       string
    Department *Department
    Skills     []string
    Metadata   map[string]string
}

// Create source instance
src := &Employee{
    ID:   1,
    Name: "Jane Smith",
    Department: &Department{
        Name:     "Engineering",
        Location: "Building A",
    },
    Skills: []string{"Go", "Python", "SQL"},
    Metadata: map[string]string{
        "hire_date": "2022-01-15",
        "status":    "active",
    },
}

// Create destination instance
dst := &Employee{}

// Perform deep copy
err := model.DeepCopy(dst, src)
if err != nil {
    // Handle error
}

// Modify source to demonstrate independence
src.Name = "Jane Brown"
src.Department.Location = "Building B"
src.Skills[0] = "Rust"
src.Metadata["status"] = "on leave"

// dst still contains the original values
```

## Best Practices

1. **Use Pointers**: Always pass pointers to structs when using CopyFields and DeepCopy
2. **Check Error Returns**: Always check the error return values from these functions
3. **Type Compatibility**: Be aware of type compatibility when copying fields between different struct types
4. **Performance Considerations**: Reflection-based operations are slower than direct assignments, so use these functions judiciously in performance-critical code
5. **Unexported Fields**: Remember that unexported (lowercase) fields will not be copied

## Troubleshooting

### Common Issues

#### Type Mismatch Errors

If you're copying between structs with fields of different types, those fields will be skipped. Ensure your struct fields have compatible types.

#### Nil Pointer Errors

Ensure both source and destination are valid pointers to structs before calling CopyFields or DeepCopy.

## Related Components

- [Validation](../validation/README.md) - Validation utilities for model objects
- [Repository](../repository/README.md) - Repository pattern implementation for storing and retrieving models

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
