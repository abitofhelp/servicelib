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

### FileSize

See the [FileSize example](../examples/valueobject/filesize_example.go) for a complete, runnable example of how to use the FileSize value object.

### ID

See the [ID example](../examples/valueobject/id_example.go) for a complete, runnable example of how to use the ID value object.

### IPAddress

See the [IPAddress example](../examples/valueobject/ipaddress_example.go) for a complete, runnable example of how to use the IPAddress value object.

### URL

See the [URL example](../examples/valueobject/url_example.go) for a complete, runnable example of how to use the URL value object.

### Username

See the [Username example](../examples/valueobject/username_example.go) for a complete, runnable example of how to use the Username value object.

### Coordinate

See the [Coordinate example](../examples/valueobject/coordinate_example.go) for a complete, runnable example of how to use the Coordinate value object.

### Money

See the [Money example](../examples/valueobject/money_example.go) for a complete, runnable example of how to use the Money value object.

### Email

See the [Email example](../examples/valueobject/email_example.go) for a complete, runnable example of how to use the Email value object.

### Address

```go
// Example usage of the Address value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new address
	address, err := valueobject.NewAddress("123 Main St, Anytown, CA 12345")
	if err != nil {
		// Handle error
		fmt.Println("Error creating address:", err)
		return
	}

	// Access the address as string
	fmt.Printf("Address: %s\n", address.String())

	// Check if address is empty
	isEmpty := address.IsEmpty()
	fmt.Printf("Is address empty? %v\n", isEmpty)

	// Compare addresses (case insensitive)
	otherAddress, _ := valueobject.NewAddress("123 Main St, Anytown, CA 12345")
	areEqual := address.Equals(otherAddress)
	fmt.Printf("Are addresses equal? %v\n", areEqual)
}
```

### Color

```go
// Example usage of the Color value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new color
	color, err := valueobject.NewColor("#FF5733")
	if err != nil {
		// Handle error
		fmt.Println("Error creating color:", err)
		return
	}

	// Create from shorthand notation
	shortColor, _ := valueobject.NewColor("#F53")
	fmt.Printf("Expanded color: %s\n", shortColor.String()) // #FF5533

	// Get RGB components
	r, g, b, err := color.RGB()
	if err != nil {
		// Handle error
		fmt.Println("Error getting RGB components:", err)
		return
	}
	fmt.Printf("RGB: (%d, %d, %d)\n", r, g, b)

	// Add alpha component
	colorWithAlpha, _ := color.WithAlpha(128)
	fmt.Printf("Color with alpha: %s\n", colorWithAlpha)

	// Check if color is dark
	isDark, _ := color.IsDark()
	fmt.Printf("Is color dark? %v\n", isDark)

	// Invert color
	inverted, _ := color.Invert()
	fmt.Printf("Inverted color: %s\n", inverted.String())
}
```

### DateOfBirth

```go
// Example usage of the DateOfBirth value object
package main

import (
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new date of birth
	dob, err := valueobject.NewDateOfBirth("1990-01-15")
	if err != nil {
		// Handle error
		fmt.Println("Error creating date of birth:", err)
		return
	}

	// Create from time.Time
	dobTime := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	dob, err = valueobject.NewDateOfBirthFromTime(dobTime)
	if err != nil {
		// Handle error
		fmt.Println("Error creating date of birth from time:", err)
		return
	}

	// Get age
	age := dob.Age()
	fmt.Printf("Age: %d years\n", age)

	// Format as string
	fmt.Printf("Date of birth: %s\n", dob.String())

	// Check if adult
	isAdult := dob.IsAdult()
	fmt.Printf("Is adult? %v\n", isAdult)
}
```

### DateOfDeath

```go
// Example usage of the DateOfDeath value object
package main

import (
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a date of birth first
	dob, _ := valueobject.NewDateOfBirth("1930-01-15")

	// Create a new date of death
	dod, err := valueobject.NewDateOfDeath("2020-05-20", dob)
	if err != nil {
		// Handle error
		fmt.Println("Error creating date of death:", err)
		return
	}

	// Create from time.Time
	dodTime := time.Date(2020, 5, 20, 0, 0, 0, 0, time.UTC)
	dod, err = valueobject.NewDateOfDeathFromTime(dodTime, dob)
	if err != nil {
		// Handle error
		fmt.Println("Error creating date of death from time:", err)
		return
	}

	// Get age at death
	ageAtDeath := dod.AgeAtDeath()
	fmt.Printf("Age at death: %d years\n", ageAtDeath)

	// Format as string
	fmt.Printf("Date of death: %s\n", dod.String())
}
```

### Duration

```go
// Example usage of the Duration value object
package main

import (
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new duration
	duration, err := valueobject.NewDuration(3, 30, 15)
	if err != nil {
		// Handle error
		fmt.Println("Error creating duration:", err)
		return
	}

	// Create from time.Duration
	timeDuration := 2*time.Hour + 45*time.Minute + 30*time.Second
	duration, err = valueobject.NewDurationFromTimeDuration(timeDuration)
	if err != nil {
		// Handle error
		fmt.Println("Error creating duration from time.Duration:", err)
		return
	}

	// Parse from string
	duration, err = valueobject.ParseDuration("1h30m45s")
	if err != nil {
		// Handle error
		fmt.Println("Error parsing duration:", err)
		return
	}

	// Access components
	hours := duration.Hours()
	minutes := duration.Minutes()
	seconds := duration.Seconds()
	fmt.Printf("Duration: %d hours, %d minutes, %d seconds\n", hours, minutes, seconds)

	// Format as string
	fmt.Printf("Formatted duration: %s\n", duration.String())

	// Convert to time.Duration
	timeDur := duration.ToTimeDuration()
	fmt.Printf("As time.Duration: %v\n", timeDur)
}
```

### Gender

```go
// Example usage of the Gender value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new gender
	gender, err := valueobject.NewGender("male")
	if err != nil {
		// Handle error
		fmt.Println("Error creating gender:", err)
		return
	}

	// Check gender type
	isMale := gender.IsMale()
	isFemale := gender.IsFemale()
	isOther := gender.IsOther()
	fmt.Printf("Is male: %v, Is female: %v, Is other: %v\n", isMale, isFemale, isOther)

	// Format as string
	fmt.Printf("Gender: %s\n", gender.String())

	// Create with different case
	gender, _ = valueobject.NewGender("FEMALE")
	fmt.Printf("Normalized gender: %s\n", gender.String()) // "female"
}
```

### Name

```go
// Example usage of the Name value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new name
	name, err := valueobject.NewName("John", "Doe")
	if err != nil {
		// Handle error
		fmt.Println("Error creating name:", err)
		return
	}

	// Create with middle name
	nameWithMiddle, _ := valueobject.NewNameWithMiddle("Jane", "Marie", "Smith")

	// Access components
	firstName := name.FirstName()
	lastName := name.LastName()
	fmt.Printf("First name: %s, Last name: %s\n", firstName, lastName)

	// Get full name
	fullName := name.FullName()
	fmt.Printf("Full name: %s\n", fullName)

	// Get initials
	initials := name.Initials()
	fmt.Printf("Initials: %s\n", initials)

	// Format with middle name
	middleInitial := nameWithMiddle.MiddleInitial()
	fmt.Printf("Name with middle initial: %s %s. %s\n", 
		nameWithMiddle.FirstName(), middleInitial, nameWithMiddle.LastName())
}
```

### Password

```go
// Example usage of the Password value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new password
	password, err := valueobject.NewPassword("P@ssw0rd123!")
	if err != nil {
		// Handle error
		fmt.Println("Error creating password:", err)
		return
	}

	// Check password strength
	strength := password.Strength()
	fmt.Printf("Password strength: %d/5\n", strength)

	// Check if password meets requirements
	hasUppercase := password.HasUppercase()
	hasLowercase := password.HasLowercase()
	hasDigit := password.HasDigit()
	hasSpecial := password.HasSpecialChar()
	fmt.Printf("Has uppercase: %v, Has lowercase: %v, Has digit: %v, Has special: %v\n",
		hasUppercase, hasLowercase, hasDigit, hasSpecial)

	// Verify a password
	isMatch := password.Verify("P@ssw0rd123!")
	fmt.Printf("Password matches: %v\n", isMatch)

	// Masked representation for logging
	fmt.Printf("Masked password: %s\n", password.Masked())
}
```

### Percentage

```go
// Example usage of the Percentage value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new percentage
	percentage, err := valueobject.NewPercentage(75.5)
	if err != nil {
		// Handle error
		fmt.Println("Error creating percentage:", err)
		return
	}

	// Access value
	value := percentage.Value()
	fmt.Printf("Percentage value: %.2f%%\n", value)

	// Format as string
	fmt.Printf("Formatted: %s\n", percentage.String()) // "75.50%"

	// Convert to decimal (0-1 range)
	decimal := percentage.ToDecimal()
	fmt.Printf("As decimal: %.4f\n", decimal) // 0.7550

	// Calculate percentage of a value
	amount := 200.0
	result := percentage.Of(amount)
	fmt.Printf("%.2f%% of %.2f = %.2f\n", value, amount, result) // 151.00
}
```

### Phone

```go
// Example usage of the Phone value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new phone number
	phone, err := valueobject.NewPhone("+1-555-123-4567")
	if err != nil {
		// Handle error
		fmt.Println("Error creating phone:", err)
		return
	}

	// Format in different ways
	fmt.Println(phone.String())                // "+1-555-123-4567"
	fmt.Println(phone.Format("E.164"))         // "+15551234567"
	fmt.Println(phone.Format("national"))      // "(555) 123-4567"
	fmt.Println(phone.Format("international")) // "+1 555 123 4567"

	// Get country code
	countryCode := phone.CountryCode()
	fmt.Printf("Country code: %s\n", countryCode)

	// Check if valid for country
	isValidUS := phone.IsValidForCountry("US")
	fmt.Printf("Valid for US: %v\n", isValidUS)

	// Get normalized version (E.164 format)
	normalized := phone.Normalized()
	fmt.Printf("Normalized: %s\n", normalized)
}
```

### Rating

```go
// Example usage of the Rating value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new rating (1-5 scale)
	rating, err := valueobject.NewRating(4.5, 5)
	if err != nil {
		// Handle error
		fmt.Println("Error creating rating:", err)
		return
	}

	// Access values
	value := rating.Value()
	scale := rating.Scale()
	fmt.Printf("Rating: %.1f out of %d\n", value, scale)

	// Format as string
	fmt.Printf("Formatted: %s\n", rating.String()) // "4.5/5"

	// Convert to percentage
	percentage := rating.ToPercentage()
	fmt.Printf("As percentage: %.1f%%\n", percentage) // 90.0%

	// Convert to different scale
	tenScale := rating.ToScale(10)
	fmt.Printf("On 10-point scale: %.1f\n", tenScale) // 9.0
}
```

### Temperature

```go
// Example usage of the Temperature value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new temperature in Celsius
	temp, err := valueobject.NewTemperature(25.5, valueobject.Celsius)
	if err != nil {
		// Handle error
		fmt.Println("Error creating temperature:", err)
		return
	}

	// Convert to different units
	fahrenheit := temp.ToFahrenheit()
	kelvin := temp.ToKelvin()
	fmt.Printf("%.1f°C = %.1f°F = %.1fK\n", temp.Value(), fahrenheit, kelvin)

	// Create from Fahrenheit
	tempF, _ := valueobject.NewTemperature(98.6, valueobject.Fahrenheit)
	celsius := tempF.ToCelsius()
	fmt.Printf("%.1f°F = %.1f°C\n", tempF.Value(), celsius)

	// Format with unit
	fmt.Println(temp.String())                     // "25.5°C"
	fmt.Println(temp.Format(valueobject.Fahrenheit)) // "77.9°F"
	fmt.Println(temp.Format(valueobject.Kelvin))     // "298.7K"
}
```

### Version

```go
// Example usage of the Version value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new semantic version
	version, err := valueobject.NewVersion("1.2.3")
	if err != nil {
		// Handle error
		fmt.Println("Error creating version:", err)
		return
	}

	// Create with pre-release and build metadata
	versionWithMeta, _ := valueobject.NewVersion("2.0.0-alpha.1+build.123")

	// Access components
	major := version.Major()
	minor := version.Minor()
	patch := version.Patch()
	fmt.Printf("Version: %d.%d.%d\n", major, minor, patch)

	// Get pre-release and build info
	preRelease := versionWithMeta.PreRelease()
	buildMeta := versionWithMeta.BuildMetadata()
	fmt.Printf("Pre-release: %s, Build metadata: %s\n", preRelease, buildMeta)

	// Compare versions
	otherVersion, _ := valueobject.NewVersion("1.3.0")
	isGreater := otherVersion.IsGreaterThan(version)
	fmt.Printf("Is 1.3.0 greater than 1.2.3? %v\n", isGreater) // true

	// Format as string
	fmt.Println(version.String())         // "1.2.3"
	fmt.Println(versionWithMeta.String()) // "2.0.0-alpha.1+build.123"
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
