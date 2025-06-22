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
| Address | Represents a physical address with validation (moved to contact package) |
| Color | Represents a color in various formats (RGB, HEX, etc.) (moved to appearance package) |
| Coordinate | Represents a geographic coordinate (latitude and longitude) |
| DateOfBirth | Represents a person's date of birth with validation |
| DateOfDeath | Represents a person's date of death with validation |
| Duration | Represents a time duration with various units |
| Email | Represents an email address with validation |
| FileSize | Represents a file size with various units using base 10 values (bytes, KB, MB, etc.) per ISO standards (moved to measurement package) |
| Gender | Represents a person's gender |
| ID | Represents a unique identifier |
| IPAddress | Represents an IP address (IPv4 or IPv6) (moved to network package) |
| MemSize | Represents a computer memory size with various units using base 2 binary values (bytes, KiB, MiB, etc.) (in measurement package) |
| Money | Represents a monetary value with currency |
| Name | Represents a person's name |
| Password | Represents a password with validation and security features |
| Percentage | Represents a percentage value |
| Phone | Represents a phone number with validation |
| Rating | Represents a rating value (e.g., 1-5 stars) |
| Temperature | Represents a temperature value with various units |
| URL | Represents a URL with validation (moved to network package) |
| Username | Represents a username with validation |
| Version | Represents a version number (e.g., semantic versioning) (moved to temporal package) |

## Usage Examples

### FileSize

The FileSize value object has been moved to the measurement package and updated to use base 10 values per ISO standards. Use it as follows:

```go
// Example usage of the FileSize value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/measurement"
)

func main() {
	// Create a new file size (using base 10 values per ISO standards)
	fileSize, err := measurement.NewFileSize(1.5, measurement.Gigabytes)
	if err != nil {
		// Handle error
		fmt.Println("Error creating file size:", err)
		return
	}

	// Access the size in different units
	bytes := fileSize.Bytes()
	kb := fileSize.Kilobytes()
	mb := fileSize.Megabytes()
	gb := fileSize.Gigabytes()

	fmt.Printf("File size: %.2f GB\n", gb)
	fmt.Printf("Same size in different units: %d bytes, %.2f KB, %.2f MB\n", 
		bytes, kb, mb)

	// Format with different units
	fmt.Println(fileSize.Format("B"))     // In bytes
	fmt.Println(fileSize.Format("KB"))    // In kilobytes
	fmt.Println(fileSize.Format("MB"))    // In megabytes
	fmt.Println(fileSize.Format("GB"))    // In gigabytes
	fmt.Println(fileSize.Format("auto"))  // Automatic (uses the most appropriate unit)

	// Parse from string
	parsedSize, err := measurement.ParseFileSize("2.5GB")
	if err != nil {
		// Handle error
		fmt.Println("Error parsing file size:", err)
		return
	}
	fmt.Println(parsedSize.String())
}
```

### MemSize

The MemSize value object is in the measurement package and uses base 2 binary values for computer memory. Use it as follows:

```go
// Example usage of the MemSize value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/measurement"
)

func main() {
	// Create a new memory size (using base 2 binary values)
	memSize, err := measurement.NewMemSize(1.5, measurement.Gibibytes)
	if err != nil {
		// Handle error
		fmt.Println("Error creating memory size:", err)
		return
	}

	// Access the size in different units
	bytes := memSize.Bytes()
	kib := memSize.Kibibytes()
	mib := memSize.Mebibytes()
	gib := memSize.Gibibytes()

	fmt.Printf("Memory size: %.2f GiB\n", gib)
	fmt.Printf("Same size in different units: %d bytes, %.2f KiB, %.2f MiB\n", 
		bytes, kib, mib)

	// Format with different units
	fmt.Println(memSize.Format("B"))     // In bytes
	fmt.Println(memSize.Format("KiB"))   // In kibibytes
	fmt.Println(memSize.Format("MiB"))   // In mebibytes
	fmt.Println(memSize.Format("GiB"))   // In gibibytes
	fmt.Println(memSize.Format("auto"))  // Automatic (uses the most appropriate unit)

	// Parse from string
	parsedSize, err := measurement.ParseMemSize("2.5GiB")
	if err != nil {
		// Handle error
		fmt.Println("Error parsing memory size:", err)
		return
	}
	fmt.Println(parsedSize.String())
}
```

See the [FileSize example](../examples/valueobject/filesize_example.go) for a complete, runnable example of how to use the FileSize value object.

### ID

See the [ID example](../examples/valueobject/id_example.go) for a complete, runnable example of how to use the ID value object.

### IPAddress

See the [IPAddress example](../examples/valueobject/ipaddress_example.go) for a complete, runnable example of how to use the IPAddress value object.

### URL

The URL value object has been moved to the network package. Use it as follows:

```go
// Example usage of the URL value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/network"
)

func main() {
	// Create a new URL
	url, err := network.NewURL("https://example.com/path?param=value")
	if err != nil {
		// Handle error
		fmt.Println("Error creating URL:", err)
		return
	}

	// Get components
	domain, _ := url.Domain()
	path, _ := url.Path()
	query, _ := url.Query()

	fmt.Printf("Domain: %s\n", domain)
	fmt.Printf("Path: %s\n", path)
	fmt.Printf("Query parameters: %v\n", query)

	// Check if URL is empty
	isEmpty := url.IsEmpty()
	fmt.Printf("Is URL empty? %v\n", isEmpty)

	// Compare URLs
	otherURL, _ := network.NewURL("https://example.com/path?param=value")
	areEqual := url.Equals(otherURL)
	fmt.Printf("Are URLs equal? %v\n", areEqual)
}
```

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

The Address value object has been moved to the contact package. Use it as follows:

```go
// Example usage of the Address value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/contact"
)

func main() {
	// Create a new address
	address, err := contact.NewAddress("123 Main St, Anytown, CA 12345")
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
	otherAddress, _ := contact.NewAddress("123 Main St, Anytown, CA 12345")
	areEqual := address.Equals(otherAddress)
	fmt.Printf("Are addresses equal? %v\n", areEqual)

	// Validate the address
	err = address.Validate()
	if err != nil {
		fmt.Println("Address validation failed:", err)
	}
}
```

### Color

The Color value object has been moved to the appearance package. Use it as follows:

```go
// Example usage of the Color value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/appearance"
)

func main() {
	// Create a new color
	color, err := appearance.NewColor("#FF5733")
	if err != nil {
		// Handle error
		fmt.Println("Error creating color:", err)
		return
	}

	// Create from shorthand notation
	shortColor, _ := appearance.NewColor("#F53")
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

The DateOfBirth value object has been moved to the identification package. Use it as follows:

```go
// Example usage of the DateOfBirth value object
package main

import (
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/valueobject/identification"
)

func main() {
	// Create a new date of birth
	dob, err := identification.NewDateOfBirth("1990-01-15")
	if err != nil {
		// Handle error
		fmt.Println("Error creating date of birth:", err)
		return
	}

	// Create from time.Time
	dobTime := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	dob, err = identification.NewDateOfBirthFromTime(dobTime)
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

The DateOfDeath value object has been moved to the identification package. Use it as follows:

```go
// Example usage of the DateOfDeath value object
package main

import (
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/valueobject/identification"
)

func main() {
	// Create a date of birth first
	dob, _ := identification.NewDateOfBirth("1930-01-15")

	// Create a new date of death
	dod, err := identification.NewDateOfDeath("2020-05-20", dob)
	if err != nil {
		// Handle error
		fmt.Println("Error creating date of death:", err)
		return
	}

	// Create from time.Time
	dodTime := time.Date(2020, 5, 20, 0, 0, 0, 0, time.UTC)
	dod, err = identification.NewDateOfDeathFromTime(dodTime, dob)
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

The Duration value object has been moved to the temporal package. Use it as follows:

```go
// Example usage of the Duration value object
package main

import (
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/valueobject/temporal"
)

func main() {
	// Create a new duration
	duration, err := temporal.NewDuration(3, 30, 15)
	if err != nil {
		// Handle error
		fmt.Println("Error creating duration:", err)
		return
	}

	// Create from time.Duration
	timeDuration := 2*time.Hour + 45*time.Minute + 30*time.Second
	duration, err = temporal.NewDurationFromTimeDuration(timeDuration)
	if err != nil {
		// Handle error
		fmt.Println("Error creating duration from time.Duration:", err)
		return
	}

	// Parse from string
	duration, err = temporal.ParseDuration("1h30m45s")
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

The Gender value object has been moved to the identification package. Use it as follows:

```go
// Example usage of the Gender value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/identification"
)

func main() {
	// Create a new gender
	gender, err := identification.NewGender("male")
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
	gender, _ = identification.NewGender("FEMALE")
	fmt.Printf("Normalized gender: %s\n", gender.String()) // "female"
}
```

### Name

The Name value object has been moved to the identification package. Use it as follows:

```go
// Example usage of the Name value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/identification"
)

func main() {
	// Create a new name
	name, err := identification.NewName("John", "Doe")
	if err != nil {
		// Handle error
		fmt.Println("Error creating name:", err)
		return
	}

	// Create with middle name
	nameWithMiddle, _ := identification.NewNameWithMiddle("Jane", "Marie", "Smith")

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

The Password value object has been moved to the identification package. Use it as follows:

```go
// Example usage of the Password value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/identification"
)

func main() {
	// Create a new password
	password, err := identification.NewPassword("P@ssw0rd123!")
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

The Percentage value object has been moved to the measurement package. Use it as follows:

```go
// Example usage of the Percentage value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/measurement"
)

func main() {
	// Create a new percentage
	percentage, err := measurement.NewPercentage(75.5)
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

The Phone value object has been moved to the contact package. Use it as follows:

```go
// Example usage of the Phone value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/contact"
)

func main() {
	// Create a new phone number
	phone, err := contact.NewPhone("+1-555-123-4567")
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

The Rating value object has been moved to the measurement package. Use it as follows:

```go
// Example usage of the Rating value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/measurement"
)

func main() {
	// Create a new rating (1-5 scale)
	rating, err := measurement.NewRating(4.5, 5)
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

The Temperature value object has been moved to the measurement package. Use it as follows:

```go
// Example usage of the Temperature value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/measurement"
)

func main() {
	// Create a new temperature in Celsius
	temp, err := measurement.NewTemperature(25.5, measurement.Celsius)
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
	tempF, _ := measurement.NewTemperature(98.6, measurement.Fahrenheit)
	celsius := tempF.ToCelsius()
	fmt.Printf("%.1f°F = %.1f°C\n", tempF.Value(), celsius)

	// Format with unit
	fmt.Println(temp.String())                     // "25.5°C"
	fmt.Println(temp.Format(measurement.Fahrenheit)) // "77.9°F"
	fmt.Println(temp.Format(measurement.Kelvin))     // "298.7K"
}
```

### Version

The Version value object has been moved to the temporal package. Use it as follows:

```go
// Example usage of the Version value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/temporal"
)

func main() {
	// Create a new semantic version
	version, err := temporal.NewVersion(1, 2, 3, "", "")
	if err != nil {
		// Handle error
		fmt.Println("Error creating version:", err)
		return
	}

	// Create with pre-release and build metadata
	versionWithMeta, _ := temporal.NewVersion(2, 0, 0, "alpha.1", "build.123")

	// Access components
	major := version.Major()
	minor := version.Minor()
	patch := version.Patch()
	fmt.Printf("Version: %d.%d.%d\n", major, minor, patch)

	// Get pre-release and build info
	preRelease := versionWithMeta.PreRelease()
	build := versionWithMeta.Build()
	fmt.Printf("Pre-release: %s, Build metadata: %s\n", preRelease, build)

	// Compare versions
	otherVersion, _ := temporal.NewVersion(1, 3, 0, "", "")
	comparison := otherVersion.CompareTo(version)
	fmt.Printf("Is 1.3.0 greater than 1.2.3? %v\n", comparison > 0) // true

	// Format as string
	fmt.Println(version.String())         // "1.2.3"
	fmt.Println(versionWithMeta.String()) // "2.0.0-alpha.1+build.123"

	// Parse from string
	parsedVersion, err := temporal.ParseVersion("1.2.3-beta+build.456")
	if err != nil {
		// Handle error
		fmt.Println("Error parsing version:", err)
		return
	}
	fmt.Println(parsedVersion.String()) // "1.2.3-beta+build.456"
}
```

## Package Structure

The value object package is organized into several sub-packages:
  - `base`: Common interfaces and utilities for value objects
    - Provides core interfaces like `ValueObject`, `Equatable`, `Comparable`, etc.
    - Contains validation utilities for common value types
    - Includes helper functions for string and numeric comparisons
    - Provides base types like `StringValueObject` and `BaseStructValueObject` that can be embedded in specific value objects
  - `appearance`: Value objects related to visual appearance and styling
    - Color: Represents and validates colors in various formats (RGB, HEX, etc.)
    - More appearance-related value objects will be added in the future
  - `contact`: Value objects related to contact information
    - Address: Represents and validates physical addresses
    - Email: Represents and validates email addresses
    - Phone: Represents and validates phone numbers with formatting options
    - More contact-related value objects will be added in the future
  - `measurement`: Value objects related to measurements and units
    - FileSize: Represents and validates file sizes with various units using base 10 values (ISO standard)
    - MemSize: Represents and validates computer memory sizes with various units using base 2 binary values
    - Money: Represents and validates monetary values with currency
    - Percentage: Represents and validates percentage values
    - Rating: Represents and validates rating values
    - Temperature: Represents and validates temperature values with unit conversion
  - `identification`: Value objects related to identification
    - ID: Represents and validates unique identifiers
    - Username: Represents and validates usernames
    - Password: Represents and validates passwords with security features
    - Name: Represents and validates person names
    - Gender: Represents and validates gender information
    - DateOfBirth: Represents and validates dates of birth
    - DateOfDeath: Represents and validates dates of death
  - `network`: Value objects related to network and internet
    - URL: Represents and validates URLs with utility methods for parsing components
    - IPAddress: Represents and validates IP addresses (IPv4 and IPv6)
    - More network-related value objects will be added in the future
  - `temporal`: Value objects related to time and versioning
    - Duration: Represents and validates time durations
    - Version: Represents and validates semantic version numbers
    - More time-related value objects will be added in the future
  - `generator`: Code generation tools for value objects
    - Provides a generator for creating new value objects
    - Includes templates for string-based and struct-based value objects
    - Offers a command-line tool for generating value objects from a configuration file

## Generic Value Object Framework

The value object package provides a generic framework for creating value objects. This framework consists of interfaces, base types, and utilities that can be used to create new value objects with minimal code.

### Interfaces

The `base` package provides several interfaces that define the common behaviors of value objects:

- `ValueObject`: The base interface for all value objects, with methods for String() and IsEmpty()
- `Equatable`: Interface for value objects that can be compared for equality
- `Comparable`: Interface for value objects that can be compared for ordering
- `Validatable`: Interface for value objects that can be validated
- `JSONMarshallable` and `JSONUnmarshallable`: Interfaces for value objects that can be marshalled to and unmarshalled from JSON

### Base Types

The `base` package also provides base types that can be embedded in specific value objects:

- `StringValueObject`: A base type for string-based value objects
- `BaseStructValueObject`: A base type for struct-based value objects

These base types provide default implementations of the common methods, which can be overridden by specific value objects.

### Creating a New Value Object

There are three ways to create a new value object:

1. **Using the Generator**: The easiest way is to use the generator, which can create both string-based and struct-based value objects from a configuration file. See the [Generator README](generator/README.md) for details.

2. **Embedding a Base Type**: You can create a new value object by embedding one of the base types:

```go
// StringValueObject example
type Email struct {
    base.StringValueObject
}

func NewEmail(email string) (Email, error) {
    // Validate the email
    if err := base.ValidateEmail(email); err != nil {
        return Email{}, err
    }

    // Create a new StringValueObject
    vo := base.NewStringValueObject(email)

    // Return the Email value object
    return Email{StringValueObject: vo}, nil
}

// Override methods as needed
func (e Email) Validate() error {
    return base.ValidateEmail(e.Value())
}
```

```go
// BaseStructValueObject example
type Currency struct {
    base.BaseStructValueObject
    Code   string
    Symbol string
    Name   string
}

func NewCurrency(code, symbol, name string) (Currency, error) {
    currency := Currency{
        Code:   code,
        Symbol: symbol,
        Name:   name,
    }

    // Validate the currency
    if err := currency.Validate(); err != nil {
        return Currency{}, err
    }

    return currency, nil
}

func (c Currency) Validate() error {
    if c.Code == "" {
        return errors.New("currency code cannot be empty")
    }

    if len(c.Code) != 3 {
        return errors.New("currency code must be 3 characters")
    }

    if c.Symbol == "" {
        return errors.New("currency symbol cannot be empty")
    }

    return nil
}

func (c Currency) String() string {
    return fmt.Sprintf("%s (%s)", c.Name, c.Code)
}

func (c Currency) IsEmpty() bool {
    return c.Code == "" && c.Symbol == "" && c.Name == ""
}

func (c Currency) Equals(other Currency) bool {
    return c.Code == other.Code && c.Symbol == other.Symbol && c.Name == other.Name
}
```

3. **Implementing the Interfaces Directly**: You can also create a new value object by implementing the interfaces directly:

```go
type Email string

func NewEmail(email string) (Email, error) {
    // Validate the email
    if err := base.ValidateEmail(email); err != nil {
        return "", err
    }

    return Email(email), nil
}

func (e Email) String() string {
    return string(e)
}

func (e Email) IsEmpty() bool {
    return e == ""
}

func (e Email) Equals(other Email) bool {
    return base.StringsEqualFold(string(e), string(other))
}

func (e Email) Validate() error {
    return base.ValidateEmail(string(e))
}
```

### Using the New Structure

For backward compatibility, the main value object types are still available through the top-level package:

```go
import "github.com/abitofhelp/servicelib/valueobject"

// Create a new email
email, err := valueobject.NewEmail("user@example.com")

// Create a new phone number
phone, err := valueobject.NewPhone("+1-555-123-4567")
```

However, for new code, it's recommended to use the subpackages directly:

```go
import (
    "github.com/abitofhelp/servicelib/valueobject/appearance"
    "github.com/abitofhelp/servicelib/valueobject/contact"
)

// Create a new email
email, err := contact.NewEmail("user@example.com")

// Create a new phone number
phone, err := contact.NewPhone("+1-555-123-4567")

// Create a new color
color, err := appearance.NewColor("#FF5733")
```

This provides access to additional functionality and a more organized structure.

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
