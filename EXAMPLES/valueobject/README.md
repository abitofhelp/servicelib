# ValueObject Package Examples

This directory contains examples demonstrating how to use the `valueobject` package, which provides implementations of various value objects for Go applications. Value objects are immutable objects that represent concepts in the domain and encapsulate validation and behavior related to those concepts.

## Examples

### 1. Coordinate Example

[location/coordinate_example/main.go](location/coordinate_example/main.go)

Demonstrates how to use the Coordinate value object for geographic coordinates.

Key concepts:
- Creating coordinate objects with latitude and longitude
- Validating coordinate values
- Calculating distance between coordinates
- Formatting coordinates for display
- Using coordinates in geospatial applications

### 2. Email Example

[contact/email_example/main.go](contact/email_example/main.go)

Shows how to use the Email value object for email addresses.

Key concepts:
- Creating email objects with validation
- Extracting domain and local parts
- Normalizing email addresses
- Comparing email addresses
- Handling invalid email formats

### 3. FileSize Example

[measurement/filesize_example/main.go](measurement/filesize_example/main.go)

Demonstrates how to use the FileSize value object for representing file sizes.

Key concepts:
- Creating file size objects with different units
- Converting between size units (bytes, KB, MB, GB)
- Formatting file sizes for display
- Comparing file sizes
- Performing arithmetic on file sizes

### 4. ID Example

[identification/id_example/main.go](identification/id_example/main.go)

Shows how to use the ID value object for entity identifiers.

Key concepts:
- Creating ID objects with different formats
- Validating ID formats
- Generating new IDs
- Comparing IDs
- Using IDs in domain entities

### 5. IPAddress Example

[network/ipaddress_example/main.go](network/ipaddress_example/main.go)

Demonstrates how to use the IPAddress value object.

Key concepts:
- Creating IP address objects
- Validating IPv4 and IPv6 addresses
- Checking if an IP is in a subnet
- Determining IP address types
- Formatting IP addresses

### 6. Money Example

[measurement/money_example/main.go](measurement/money_example/main.go)

Shows how to use the Money value object for monetary values.

Key concepts:
- Creating money objects with amount and currency
- Performing arithmetic operations (add, subtract, multiply, divide)
- Comparing money values
- Handling currency conversion
- Preventing operations between different currencies

### 7. URL Example

[network/url_example/main.go](network/url_example/main.go)

Demonstrates how to use the URL value object.

Key concepts:
- Creating URL objects with validation
- Extracting URL components (scheme, host, path)
- Building URLs from components
- Comparing URLs
- Handling query parameters

### 8. Username Example

[identification/username_example/main.go](identification/username_example/main.go)

Shows how to use the Username value object.

Key concepts:
- Creating username objects with validation
- Enforcing username rules
- Normalizing usernames
- Comparing usernames
- Handling invalid username formats

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run measurement/money_example/main.go
```

## Additional Resources

For more information about the valueobject package, see the [valueobject package documentation](../../valueobject/README.md).
