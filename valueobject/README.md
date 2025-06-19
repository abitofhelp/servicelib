# Value Objects

## Overview

Value Objects are immutable objects that represent a concept in your domain with no identity. They are defined by their attributes rather than by a unique identifier. Value Objects are used to encapsulate and validate domain values, ensuring that they are always in a valid state.

This package provides a collection of common Value Objects that can be used in your applications.

## Characteristics of Value Objects

- **Immutability**: Once created, Value Objects cannot be changed. Any operation that would change a Value Object returns a new instance.
- **Equality by value**: Two Value Objects are equal if all their attributes are equal.
- **Self-validation**: Value Objects validate their attributes during creation and reject invalid values.
- **Encapsulation**: Value Objects encapsulate domain rules and behaviors related to the values they represent.

## Available Value Objects

| Value Object | Description |
|--------------|-------------|
| Address | Represents a physical address with validation |
| Color | Represents a color in various formats (RGB, HEX, etc.) |
| Coordinate | Represents a geographic coordinate (latitude and longitude) |
| DateOfBirth | Represents a person's date of birth with validation |
| DateOfDeath | Represents a person's date of death with validation |
| Duration | Represents a time duration with various units |
| Email | Represents an email address with validation |
| FileSize | Represents a file size with various units (bytes, KB, MB, etc.) |
| Gender | Represents a person's gender |
| ID | Represents a unique identifier |
| IPAddress | Represents an IP address (IPv4 or IPv6) |
| Money | Represents a monetary value with currency |
| Name | Represents a person's name |
| Password | Represents a password with validation and security features |
| Percentage | Represents a percentage value |
| Phone | Represents a phone number with validation |
| Rating | Represents a rating value (e.g., 1-5 stars) |
| Temperature | Represents a temperature value with various units |
| URL | Represents a URL with validation |
| Username | Represents a username with validation |
| Version | Represents a version number (e.g., semantic versioning) |

## Usage Examples

### Coordinate

```go
// Example usage of the Coordinate value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new coordinate
	coord, err := valueobject.NewCoordinate(40.7128, -74.0060)
	if err != nil {
		// Handle error
		fmt.Println("Error creating coordinate:", err)
		return
	}

	// Access values
	lat := coord.Latitude()
	lng := coord.Longitude()
	fmt.Printf("Coordinate: %f, %f\n", lat, lng)

	// Calculate distance to another coordinate
	otherCoord, _ := valueobject.NewCoordinate(34.0522, -118.2437)
	distance := coord.DistanceTo(otherCoord) // in kilometers
	fmt.Printf("Distance: %.2f km\n", distance)

	// Format in different ways
	fmt.Println(coord.Format("dms")) // 40°42'46.1"N, 74°0'21.6"W
	fmt.Println(coord.Format("dm"))  // 40°42.7683'N, 74°0.3600'W
	fmt.Println(coord.Format("dd"))  // 40.71280, -74.00600
}
```

### Money

```go
// Example usage of the Money value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
	"github.com/shopspring/decimal"
)

func main() {
	// Create a new money value
	money, err := valueobject.NewMoney(10.99, "USD")
	if err != nil {
		// Handle error
		fmt.Println("Error creating money:", err)
		return
	}

	// Create from string for better precision
	money, err = valueobject.NewMoneyFromString("10.99", "USD")
	if err != nil {
		// Handle error
		fmt.Println("Error creating money from string:", err)
		return
	}

	// Access values
	amount := money.Amount()
	currency := money.Currency()
	fmt.Printf("Money: %.2f %s\n", amount, currency)

	// Perform calculations
	otherMoney, _ := valueobject.NewMoney(5.99, "USD")
	sum, err := money.Add(otherMoney)
	if err != nil {
		// Handle error (e.g., different currencies)
		fmt.Println("Error adding money:", err)
		return
	}
	fmt.Printf("Sum: %s\n", sum.String())

	// Multiply by a factor
	factor := decimal.NewFromFloat(1.1) // 10% increase
	multiplied := money.Multiply(factor)
	fmt.Printf("After 10%% increase: %s\n", multiplied.String())

	// Format as string
	fmt.Println(money.String()) // "10.99 USD"

	// Parse from string
	parsed, err := valueobject.Parse("15.75 EUR")
	if err != nil {
		// Handle error
		fmt.Println("Error parsing money:", err)
		return
	}
	fmt.Printf("Parsed: %s\n", parsed.String())
}
```

### Email

```go
// Example usage of the Email value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new email
	email, err := valueobject.NewEmail("user@example.com")
	if err != nil {
		// Handle error (invalid email format)
		fmt.Println("Error creating email:", err)
		return
	}

	// Access values
	address := email.Address()
	domain := email.Domain()
	fmt.Printf("Email address: %s, Domain: %s\n", address, domain)

	// Check if it's a specific domain
	isGmail := email.IsDomain("gmail.com")
	fmt.Printf("Is Gmail: %v\n", isGmail)

	// Format as string
	fmt.Println(email.String()) // "user@example.com"
}
```

## Common Patterns

All Value Objects in this package follow these common patterns:

1. **Constructor Functions**: Each Value Object has one or more constructor functions (e.g., `New<ObjectName>`) that validate inputs and return errors for invalid values.

2. **Parser Functions**: Many Value Objects have parser functions to create objects from string representations.

3. **Getter Methods**: Value Objects provide getter methods to access their attributes.

4. **String Method**: All Value Objects implement the `String()` method for string representation.

5. **Equals Method**: Value Objects provide an `Equals()` method to compare with other instances.

6. **Domain-Specific Methods**: Value Objects include methods specific to their domain (e.g., `DistanceTo()` for Coordinate).

7. **JSON Marshaling/Unmarshaling**: Many Value Objects implement JSON marshaling and unmarshaling for serialization.

## Best Practices

- Always check for errors when creating Value Objects.
- Treat Value Objects as immutable; never modify their internal state.
- Use the appropriate Value Object for the domain concept you're modeling.
- Create custom Value Objects for domain-specific concepts not covered by this package.
- Use Value Objects to encapsulate validation logic and business rules related to values.

## Contributing

If you need a Value Object that isn't provided by this package, consider creating it following the patterns established here and contributing it back to the package.
