# Validation Package

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

## Usage

### Basic Validation

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/validation"
)

func main() {
    // Create a validation result
    result := validation.NewValidationResult()
    
    // Validate a username
    username := "john"
    validation.Required(username, "username", result)
    validation.MinLength(username, 5, "username", result)
    validation.MaxLength(username, 50, "username", result)
    
    // Validate an email
    email := "john@example.com"
    validation.Required(email, "email", result)
    validation.Pattern(email, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "email", result)
    
    // Check if validation passed
    if !result.IsValid() {
        fmt.Printf("Validation failed: %v\n", result.Error())
        return
    }
    
    fmt.Println("Validation passed!")
}
```

### Date Validation

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/abitofhelp/servicelib/validation"
)

func main() {
    // Create a validation result
    result := validation.NewValidationResult()
    
    // Validate a birth date (must be in the past)
    birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
    validation.PastDate(birthDate, "birthDate", result)
    
    // Validate a date range
    startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
    endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
    validation.ValidDateRange(startDate, endDate, "startDate", "endDate", result)
    
    // Check if validation passed
    if !result.IsValid() {
        fmt.Printf("Validation failed: %v\n", result.Error())
        return
    }
    
    fmt.Println("Validation passed!")
}
```

### Collection Validation

```go
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
    
    // Using AllTrue
    allNonEmpty := validation.AllTrue(tags, func(tag string) bool {
        return tag != ""
    })
    
    if !allNonEmpty {
        result.AddError("all tags must be non-empty", "tags")
    }
    
    // Using ValidateAll
    validation.ValidateAll(tags, func(tag string, index int, result *validation.ValidationResult) {
        if tag == "" {
            result.AddError(fmt.Sprintf("tag at index %d is empty", index), "tags")
        }
    }, result)
    
    // Check if validation passed
    if !result.IsValid() {
        fmt.Printf("Validation failed: %v\n", result.Error())
        return
    }
    
    fmt.Println("Validation passed!")
}
```

### Custom Validation

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/validation"
)

// User represents a user in the system
type User struct {
    ID       string
    Username string
    Email    string
    Age      int
    Role     string
}

// ValidateUser validates a user
func ValidateUser(user User) error {
    result := validation.NewValidationResult()
    
    // Validate ID
    validation.ValidateID(user.ID, "id", result)
    
    // Validate username
    validation.Required(user.Username, "username", result)
    validation.MinLength(user.Username, 3, "username", result)
    validation.MaxLength(user.Username, 50, "username", result)
    
    // Validate email
    validation.Required(user.Email, "email", result)
    validation.Pattern(user.Email, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "email", result)
    
    // Validate age
    if user.Age < 18 {
        result.AddError("must be at least 18 years old", "age")
    }
    
    // Validate role
    if user.Role != "admin" && user.Role != "user" && user.Role != "guest" {
        result.AddError("must be one of: admin, user, guest", "role")
    }
    
    return result.Error()
}

func main() {
    // Create a user
    user := User{
        ID:       "123",
        Username: "jo", // Too short
        Email:    "invalid-email", // Invalid format
        Age:      16, // Too young
        Role:     "superuser", // Invalid role
    }
    
    // Validate the user
    if err := ValidateUser(user); err != nil {
        fmt.Printf("User validation failed: %v\n", err)
        return
    }
    
    fmt.Println("User validation passed!")
}
```

## Best Practices

1. **Create a Validation Result Early**: Create a validation result at the beginning of your validation function and pass it to all validators.

   ```go
   func ValidateRequest(req Request) error {
       result := validation.NewValidationResult()
       
       validation.Required(req.Name, "name", result)
       validation.Required(req.Email, "email", result)
       
       return result.Error()
   }
   ```

2. **Use Field Names Consistently**: Use consistent field names across your application to make error messages more understandable.

   ```go
   // Good - consistent field names
   validation.Required(user.FirstName, "firstName", result)
   validation.Required(user.LastName, "lastName", result)
   
   // Avoid - inconsistent field names
   validation.Required(user.FirstName, "first_name", result)
   validation.Required(user.LastName, "last", result)
   ```

3. **Combine Validators**: Combine multiple validators for complex validation rules.

   ```go
   // Validate a password
   validation.Required(password, "password", result)
   validation.MinLength(password, 8, "password", result)
   validation.Pattern(password, `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).*$`, "password", result)
   ```

4. **Validate Collections Properly**: Use ValidateAll for collections where each item needs validation.

   ```go
   // Validate a list of addresses
   validation.ValidateAll(user.Addresses, func(address Address, index int, result *validation.ValidationResult) {
       validation.Required(address.Street, fmt.Sprintf("addresses[%d].street", index), result)
       validation.Required(address.City, fmt.Sprintf("addresses[%d].city", index), result)
       validation.Required(address.Country, fmt.Sprintf("addresses[%d].country", index), result)
   }, result)
   ```

5. **Return Early for Critical Validations**: For performance reasons, return early if critical validations fail.

   ```go
   func ValidateOrder(order Order) error {
       result := validation.NewValidationResult()
       
       // Critical validation - return early if ID is missing
       if order.ID == "" {
           result.AddError("is required", "id")
           return result.Error()
       }
       
       // Continue with other validations
       validation.Required(order.CustomerID, "customerId", result)
       // ...
       
       return result.Error()
   }
   ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.