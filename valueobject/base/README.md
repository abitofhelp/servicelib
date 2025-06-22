# Base Value Object Package

## Overview

The `base` package provides common interfaces and utilities for value objects. It defines the core interfaces that all value objects should implement and provides utility functions for validation and comparison.

## Interfaces

### ValueObject

The `ValueObject` interface defines the common methods that all value objects should implement:

```go
type ValueObject interface {
    // String returns the string representation of the value object.
    String() string

    // IsEmpty checks if the value object is empty (zero value).
    IsEmpty() bool
}
```

### Equatable

The `Equatable` interface defines methods for value objects that can be compared for equality:

```go
type Equatable[T any] interface {
    // Equals checks if two value objects are equal.
    Equals(other T) bool
}
```

### Comparable

The `Comparable` interface defines methods for value objects that can be compared for ordering:

```go
type Comparable[T any] interface {
    // CompareTo compares this value object with another and returns:
    // -1 if this < other
    //  0 if this == other
    //  1 if this > other
    CompareTo(other T) int
}
```

### Validatable

The `Validatable` interface defines methods for value objects that can be validated:

```go
type Validatable interface {
    // Validate checks if the value object is valid.
    // Returns nil if valid, otherwise returns an error.
    Validate() error
}
```

### JSONMarshallable and JSONUnmarshallable

These interfaces define methods for value objects that can be marshalled to and unmarshalled from JSON:

```go
type JSONMarshallable interface {
    // MarshalJSON implements the json.Marshaler interface.
    MarshalJSON() ([]byte, error)
}

type JSONUnmarshallable interface {
    // UnmarshalJSON implements the json.Unmarshaler interface.
    UnmarshalJSON(data []byte) error
}
```

## Utilities

### Validation

The `validation.go` file provides common validation functions for value objects:

```go
// ValidateEmail validates an email address.
func ValidateEmail(email string) error

// ValidateURL validates a URL.
func ValidateURL(url string) error

// ValidatePhone validates a phone number.
func ValidatePhone(phone string) error

// ValidateUsername validates a username.
func ValidateUsername(username string) error

// ValidateCurrencyCode validates an ISO currency code.
func ValidateCurrencyCode(code string) error
```

### Comparison

The `utils.go` file provides utility functions for comparing values:

```go
// CompareStrings compares two strings case-insensitively.
func CompareStrings(a, b string) int

// StringsEqualFold checks if two strings are equal, ignoring case.
func StringsEqualFold(a, b string) bool

// CompareFloats compares two float64 values with a small epsilon.
func CompareFloats(a, b float64) int

// FloatsEqual checks if two float64 values are equal, with a small epsilon.
func FloatsEqual(a, b float64) bool
```

## Usage

To use the base package, import it in your value object implementation:

```go
import "github.com/abitofhelp/servicelib/valueobject/base"
```

Then implement the appropriate interfaces and use the utility functions as needed.