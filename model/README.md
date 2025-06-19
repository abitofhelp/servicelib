# Model Package

The `model` package provides utilities for working with domain models and Data Transfer Objects (DTOs) in Go applications. It includes functions for copying fields between structs and creating deep copies of objects.

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

## Usage

### Copying Fields Between Structs

```go
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
```

### Creating Deep Copies of Objects

```go
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
```

### Error Handling

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/model"
)

type Source struct {
    Name string
    Age  int
}

type Destination struct {
    Name    string
    Address string // Different field type
}

func main() {
    // Create source and destination
    src := &Source{Name: "John", Age: 30}
    dst := &Destination{}
    
    // Copy fields
    err := model.CopyFields(dst, src)
    if err != nil {
        fmt.Printf("Error copying fields: %v\n", err)
        return
    }
    
    // Print result
    fmt.Printf("Destination: %+v\n", dst)
    // Output: Destination: {Name:John Address:}
    
    // Try to copy between incompatible types
    type NotAStruct string
    srcNotStruct := NotAStruct("test")
    dstNotStruct := NotAStruct("")
    
    err = model.CopyFields(&dstNotStruct, &srcNotStruct)
    if err != nil {
        fmt.Printf("Expected error: %v\n", err)
        // Output: Expected error: both source and destination must be pointers to structs
    }
}
```

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