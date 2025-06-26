# Value Object

## Overview

The Value Object component provides a comprehensive framework for implementing domain value objects in Go applications. Value objects are immutable objects that represent concepts or entities in the domain model. They are defined by their attributes rather than their identity, meaning two value objects with the same attributes are considered equal regardless of their memory location.

This package implements the Value Object pattern from Domain-Driven Design (DDD) and provides a foundation for creating robust, type-safe domain models.

## Features

- **Immutable Objects**: Value objects are immutable, ensuring data integrity and thread safety
- **Equality by Value**: Two value objects with the same attributes are considered equal
- **Validation**: Built-in validation mechanisms to ensure value objects maintain invariants
- **Type Safety**: Strong typing for domain concepts to prevent errors
- **Serialization Support**: Support for persistence and data transfer
- **Domain-Specific Types**: Implementations for common domain concepts
- **Composability**: Value objects can be composed to create complex domain models

## Installation

```bash
go get github.com/abitofhelp/servicelib/valueobject
```

## Quick Start

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/valueobject/contact"
    "github.com/abitofhelp/servicelib/valueobject/identification"
)

func main() {
    // Create an email address value object
    email, err := contact.NewEmail("user@example.com")
    if err != nil {
        fmt.Printf("Error creating email: %v\n", err)
        return
    }
    
    // Create a UUID identifier
    id, err := identification.NewUUID()
    if err != nil {
        fmt.Printf("Error creating UUID: %v\n", err)
        return
    }
    
    fmt.Printf("Email: %s\n", email)
    fmt.Printf("UUID: %s\n", id)
}
```

## API Documentation

### Core Types

The Value Object package is organized into several subpackages, each containing value objects for specific domain concepts:

#### Base

Core interfaces and base implementations for value objects.

```go
// ValueObject is the base interface for all value objects
type ValueObject interface {
    // Equals checks if two value objects are equal
    Equals(other interface{}) bool
    // Validate checks if the value object is valid
    Validate() error
}
```

#### Appearance

Value objects related to visual appearance (color, style, etc.).

```go
// Color represents a color value object in hexadecimal format (#RRGGBB)
type Color string

// NewColor creates a new Color with validation
func NewColor(hexColor string) (Color, error)
```

#### Contact

Value objects for contact information (email, phone, address, etc.).

```go
// Email represents an email address value object
type Email string

// NewEmail creates a new Email with validation
func NewEmail(email string) (Email, error)
```

#### Identification

Value objects for identifiers (UUID, custom IDs, etc.).

```go
// UUID represents a UUID value object
type UUID string

// NewUUID creates a new UUID
func NewUUID() (UUID, error)
```

### Key Methods

Most value objects implement the following methods:

#### New Constructor

Creates a new value object with validation.

```go
func NewX(value Type) (X, error)
```

#### Equals

Checks if two value objects are equal.

```go
func (x X) Equals(other X) bool
```

#### Validate

Validates the value object.

```go
func (x X) Validate() error
```

#### String

Returns the string representation of the value object.

```go
func (x X) String() string
```

## Examples

For complete, runnable examples, see the following code snippets:

### Creating and Validating Value Objects

```go
// Create an email address value object
email, err := contact.NewEmail("user@example.com")
if err != nil {
    // Handle validation error
}

// Create a UUID identifier
id, err := identification.NewUUID()
if err != nil {
    // Handle error
}

// Create a date range
dateRange, err := temporal.NewDateRange(startDate, endDate)
if err != nil {
    // Handle validation error
}
```

### Comparing Value Objects

```go
// Create two email addresses
email1, _ := contact.NewEmail("user@example.com")
email2, _ := contact.NewEmail("user@example.com")
email3, _ := contact.NewEmail("other@example.com")

// Compare them
fmt.Println(email1.Equals(email2)) // true
fmt.Println(email1.Equals(email3)) // false
```

### Using Value Objects in Domain Models

```go
type User struct {
    ID       identification.UUID
    Email    contact.Email
    Address  location.Address
    Birthday temporal.Date
}

func NewUser(id identification.UUID, email contact.Email, address location.Address, birthday temporal.Date) (*User, error) {
    user := &User{
        ID:       id,
        Email:    email,
        Address:  address,
        Birthday: birthday,
    }
    
    // Additional validation if needed
    
    return user, nil
}
```

## Best Practices

1. **Keep Value Objects Immutable**: Never modify a value object after creation; create a new one instead
2. **Validate at Creation Time**: Perform all validation in the constructor to ensure objects are always valid
3. **Use Type-Specific Methods**: Add domain-specific methods to value objects to encapsulate behavior
4. **Prefer Value Objects Over Primitives**: Use value objects instead of primitive types to add domain meaning
5. **Compose Value Objects**: Build complex domain models by composing value objects

## Troubleshooting

### Common Issues

#### Invalid Value Objects

If you're getting validation errors when creating value objects, check that the input values meet the requirements:

```go
// This will fail validation
email, err := contact.NewEmail("invalid-email")
if err != nil {
    fmt.Printf("Error: %v\n", err) // Will print validation error
}

// This will pass validation
email, err = contact.NewEmail("valid@example.com")
if err != nil {
    // Should not happen
}
```

#### Modifying Value Objects

Remember that value objects are immutable. To "change" a value object, create a new one:

```go
// Incorrect - trying to modify a value object
// color.SetHex("#FF0000") // This method doesn't exist

// Correct - create a new value object
oldColor, _ := appearance.NewColor("#000000")
newColor, _ := appearance.NewColor("#FF0000")
```

## Related Components

- [Model](../model/README.md) - Domain model components that use value objects
- [Validation](../validation/README.md) - Validation utilities used by value objects
- [Errors](../errors/README.md) - Error handling for value object validation

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.