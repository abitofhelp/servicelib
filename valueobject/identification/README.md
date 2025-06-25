# Identification Value Objects

## Overview

The `identification` package provides value objects related to identification, such as IDs, usernames, passwords, names, genders, dates of birth, and dates of death. These value objects encapsulate validation logic and provide methods for common operations on identification information.

## Value Objects

### ID

The `ID` value object represents a unique identifier.

```go
// Create a new ID
id, err := identification.NewID("123e4567-e89b-12d3-a456-426614174000")
if err != nil {
    // Handle error
}

// Check if ID is empty
isEmpty := id.IsEmpty()

// Compare IDs
otherID, _ := identification.NewID("123e4567-e89b-12d3-a456-426614174000")
areEqual := id.Equals(otherID) // true

// Get string representation
idStr := id.String() // "123e4567-e89b-12d3-a456-426614174000"

// Validate an existing ID
err = id.Validate()
```

### Username

The `Username` value object represents a username with validation.

```go
// Create a new username
username, err := identification.NewUsername("johndoe")
if err != nil {
    // Handle error
}

// Check if username is empty
isEmpty := username.IsEmpty()

// Compare usernames (case insensitive)
otherUsername, _ := identification.NewUsername("JohnDoe")
areEqual := username.Equals(otherUsername) // true

// Get string representation
usernameStr := username.String() // "johndoe"

// Validate an existing username
err = username.Validate()
```

### Password

The `Password` value object represents a password with validation and security features.

```go
// Create a new password
password, err := identification.NewPassword("P@ssw0rd123!")
if err != nil {
    // Handle error
}

// Check password strength
strength := password.Strength()

// Check if password meets requirements
hasUppercase := password.HasUppercase()
hasLowercase := password.HasLowercase()
hasDigit := password.HasDigit()
hasSpecial := password.HasSpecialChar()

// Verify a password
isMatch := password.Verify("P@ssw0rd123!")

// Get masked representation for logging
maskedPassword := password.Masked() // "********"
```

### Name

The `Name` value object represents a person's name.

```go
// Create a new name
name, err := identification.NewName("John", "Doe")
if err != nil {
    // Handle error
}

// Create with middle name
nameWithMiddle, _ := identification.NewNameWithMiddle("Jane", "Marie", "Smith")

// Access components
firstName := name.FirstName()
lastName := name.LastName()

// Get full name
fullName := name.FullName() // "John Doe"

// Get initials
initials := name.Initials() // "JD"

// Format with middle name
middleInitial := nameWithMiddle.MiddleInitial()
formattedName := fmt.Sprintf("%s %s. %s", nameWithMiddle.FirstName(), middleInitial, nameWithMiddle.LastName())
```

### Gender

The `Gender` value object represents a person's gender.

```go
// Create a new gender
gender, err := identification.NewGender("male")
if err != nil {
    // Handle error
}

// Check gender type
isMale := gender.IsMale()
isFemale := gender.IsFemale()
isOther := gender.IsOther()

// Format as string
genderStr := gender.String() // "male"

// Create with different case
gender, _ = identification.NewGender("FEMALE")
normalizedGender := gender.String() // "female"
```

### DateOfBirth

The `DateOfBirth` value object represents a person's date of birth with validation.

```go
// Create a new date of birth
dob, err := identification.NewDateOfBirth(1990, 1, 15)
if err != nil {
    // Handle error
}

// Create from time.Time
dobTime := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
dob, err = identification.NewDateOfBirthFromTime(dobTime)
if err != nil {
    // Handle error
}

// Get age
age := dob.Age()

// Format as string
dobStr := dob.String() // "1990-01-15"

// Check if adult
isAdult := dob.IsAdult()
```

### DateOfDeath

The `DateOfDeath` value object represents a person's date of death with validation.

```go
// Create a date of birth first
dob, _ := identification.NewDateOfBirth(1930, 1, 15)

// Create a new date of death
dod, err := identification.NewDateOfDeath(2020, 5, 20, dob)
if err != nil {
    // Handle error
}

// Create from time.Time
dodTime := time.Date(2020, 5, 20, 0, 0, 0, 0, time.UTC)
dod, err = identification.NewDateOfDeathFromTime(dodTime, dob)
if err != nil {
    // Handle error
}

// Get age at death
ageAtDeath := dod.AgeAtDeath()

// Format as string
dodStr := dod.String() // "2020-05-20"
```

## Usage

To use the identification value objects, import the package:

```go
import "github.com/abitofhelp/servicelib/valueobject/identification"
```

Then create and use the value objects as shown in the examples above.

## Backward Compatibility

For backward compatibility, these value objects are also available through the main `valueobject` package:

```go
import "github.com/abitofhelp/servicelib/valueobject"

// Create a new ID
id, err := valueobject.NewID("123e4567-e89b-12d3-a456-426614174000")

// Create a new username
username, err := valueobject.NewUsername("johndoe")
```

However, new code should use the `identification` package directly to take advantage of the more organized structure and additional functionality.

## Features

- **Feature 1**: Description of feature 1
- **Feature 2**: Description of feature 2
- **Feature 3**: Description of feature 3

## Installation

```bash
go get github.com/abitofhelp/servicelib/identification
```

## Quick Start

See the [Quick Start example](../EXAMPLES/identification/quickstart_example.go) for a complete, runnable example of how to use the identification.

## Configuration

See the [Configuration example](../EXAMPLES/identification/configuration_example.go) for a complete, runnable example of how to configure the identification.

## API Documentation


### Core Types

Description of the main types provided by the identification.

#### Type 1

Description of Type 1 and its purpose.

See the [Type 1 example](../EXAMPLES/identification/type1_example.go) for a complete, runnable example of how to use Type 1.

### Key Methods

Description of the key methods provided by the identification.

#### Method 1

Description of Method 1 and its purpose.

See the [Method 1 example](../EXAMPLES/identification/method1_example.go) for a complete, runnable example of how to use Method 1.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/identification/basic_usage_example.go) - Shows basic usage of the identification
- [Advanced Configuration](../EXAMPLES/identification/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/identification/error_handling_example.go) - Shows how to handle errors

## Best Practices

1. **Best Practice 1**: Description of best practice 1
2. **Best Practice 2**: Description of best practice 2
3. **Best Practice 3**: Description of best practice 3

## Troubleshooting

### Common Issues

#### Issue 1

Description of issue 1 and how to resolve it.

#### Issue 2

Description of issue 2 and how to resolve it.

## Related Components

- [Component 1](../identification1/README.md) - Description of how this identification relates to Component 1
- [Component 2](../identification2/README.md) - Description of how this identification relates to Component 2

## Contributing

Contributions to this identification are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
