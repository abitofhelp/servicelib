# model field_copying Example

## Overview

This example demonstrates how to use the model.CopyFields function from the ServiceLib model package to copy fields between structs. It shows how to copy data from a domain model to a DTO (Data Transfer Object) while handling different field sets.

## Features

- **Struct Field Copying**: Copy fields between different struct types
- **Field Matching**: Only copy fields with matching names and compatible types
- **Selective Copying**: Copy only a subset of fields from the source struct
- **Error Handling**: Handle errors that might occur during the copying process

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Struct Definitions

The example defines two struct types: a domain model (User) and a DTO (UserResponse):

```go
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
```

### Field Copying

The example demonstrates how to use model.CopyFields to copy fields from the domain model to the DTO:

```go
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
```

### Error Handling

The example shows how to handle errors that might occur during the copying process:

```go
err := model.CopyFields(response, user)
if err != nil {
    fmt.Printf("Error copying fields: %v\n", err)
    return
}
```

## Expected Output

When you run the example, you should see the following output:

```
User Response: {ID:123 Username:johndoe Email:john.doe@example.com FirstName:John LastName:Doe}
```

Note that only the fields that exist in both structs are copied. The Age and Active fields from the User struct are not copied to the UserResponse struct because they don't exist in the target struct. The private field is also not copied because it's unexported.

## Related Examples


- [deep_copy](../deep_copy/README.md) - Related example for deep_copy
- [error_handling](../error_handling/README.md) - Related example for error_handling

## Related Components

- [model Package](../../../model/README.md) - The model package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
