# Contact Value Objects

## Overview

The `contact` package provides value objects related to contact information, such as email addresses, phone numbers, and physical addresses. These value objects encapsulate validation logic and provide methods for common operations on contact information.

## Value Objects

### Email

The `Email` value object represents an email address with validation.

```go
// Create a new email
email, err := contact.NewEmail("user@example.com")
if err != nil {
    // Handle error
}

// Check if email is empty
isEmpty := email.IsEmpty()

// Compare emails (case insensitive)
otherEmail, _ := contact.NewEmail("USER@EXAMPLE.COM")
areEqual := email.Equals(otherEmail) // true

// Get string representation
emailStr := email.String() // "user@example.com"

// Validate an existing email
err = email.Validate()
```

### Address

The `Address` value object represents a postal address with validation.

```go
// Create a new address
address, err := contact.NewAddress("123 Main St, Anytown, CA 12345")
if err != nil {
    // Handle error
}

// Check if address is empty
isEmpty := address.IsEmpty()

// Compare addresses (case insensitive)
otherAddress, _ := contact.NewAddress("123 MAIN ST, ANYTOWN, CA 12345")
areEqual := address.Equals(otherAddress) // true

// Get string representation
addressStr := address.String() // "123 Main St, Anytown, CA 12345"

// Validate an existing address
err = address.Validate()
```

### Phone

The `Phone` value object represents a phone number with validation and formatting options.

```go
// Create a new phone number
phone, err := contact.NewPhone("+1-555-123-4567")
if err != nil {
    // Handle error
}

// Format in different ways
e164 := phone.Format("e164")         // "+15551234567"
national := phone.Format("national")      // "(555) 123-4567"
international := phone.Format("international") // "+1 555 123 4567"

// Get country code
countryCode := phone.CountryCode() // "1"

// Check if valid for country
isValidUS := phone.IsValidForCountry("1") // true

// Get normalized version (E.164 format)
normalized := phone.Normalized() // "+15551234567"

// Compare phone numbers (ignoring formatting)
otherPhone, _ := contact.NewPhone("555-123-4567")
areEqual := phone.Equals(otherPhone) // true
```

## Usage

To use the contact value objects, import the package:

```go
import "github.com/abitofhelp/servicelib/valueobject/contact"
```

Then create and use the value objects as shown in the examples above.

## Backward Compatibility

For backward compatibility, Email and Phone value objects are also available through the main `valueobject` package. However, the Address value object is now only available in the contact package:

```go
import "github.com/abitofhelp/servicelib/valueobject"
import "github.com/abitofhelp/servicelib/valueobject/contact"

// Create a new email
email, err := valueobject.NewEmail("user@example.com")

// Create a new phone number
phone, err := valueobject.NewPhone("+1-555-123-4567")

// Create a new address (only available in contact package)
address, err := contact.NewAddress("123 Main St, Anytown, CA 12345")
```

However, new code should use the `contact` package directly to take advantage of the more organized structure and additional functionality.

## Features

- **Feature 1**: Description of feature 1
- **Feature 2**: Description of feature 2
- **Feature 3**: Description of feature 3

## Installation

```bash
go get github.com/abitofhelp/servicelib/contact
```

## Quick Start

See the [Quick Start example](../EXAMPLES/contact/quickstart_example.go) for a complete, runnable example of how to use the contact.

## Configuration

See the [Configuration example](../EXAMPLES/contact/configuration_example.go) for a complete, runnable example of how to configure the contact.

## API Documentation


### Core Types

Description of the main types provided by the contact.

#### Type 1

Description of Type 1 and its purpose.

See the [Type 1 example](../EXAMPLES/contact/type1_example.go) for a complete, runnable example of how to use Type 1.

### Key Methods

Description of the key methods provided by the contact.

#### Method 1

Description of Method 1 and its purpose.

See the [Method 1 example](../EXAMPLES/contact/method1_example.go) for a complete, runnable example of how to use Method 1.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/contact/basic_usage_example.go) - Shows basic usage of the contact
- [Advanced Configuration](../EXAMPLES/contact/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/contact/error_handling_example.go) - Shows how to handle errors

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

- [Component 1](../contact1/README.md) - Description of how this contact relates to Component 1
- [Component 2](../contact2/README.md) - Description of how this contact relates to Component 2

## Contributing

Contributions to this contact are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
