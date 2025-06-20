# ServiceLib Examples

This directory contains complete example applications that demonstrate how to use ServiceLib in real-world scenarios.

## Available Examples

### Value Object Examples

The following examples demonstrate how to use various value objects from the valueobject package:

- **[coordinate_example.go](valueobject/coordinate_example.go)** - Demonstrates how to use the Coordinate value object, including creating coordinates, accessing latitude and longitude, calculating distances, and checking hemisphere information.
- **[email_example.go](valueobject/email_example.go)** - Demonstrates how to use the Email value object, including validation, accessing address and domain components, and domain checking.
- **[filesize_example.go](valueobject/filesize_example.go)** - Demonstrates how to use the FileSize value object, including creating file sizes, converting between units, formatting, and performing calculations.
- **[id_example.go](valueobject/id_example.go)** - Demonstrates how to use the ID value object, including generating random IDs, creating IDs from strings, and comparing IDs.
- **[ipaddress_example.go](valueobject/ipaddress_example.go)** - Demonstrates how to use the IPAddress value object, including creating IPv4 and IPv6 addresses, checking address types, and comparing addresses.
- **[money_example.go](valueobject/money_example.go)** - Demonstrates how to use the Money value object, including creating Money objects, performing arithmetic operations, comparing values, and handling currency.
- **[url_example.go](valueobject/url_example.go)** - Demonstrates how to use the URL value object, including creating URLs, accessing components, and comparing URLs.
- **[username_example.go](valueobject/username_example.go)** - Demonstrates how to use the Username value object, including validation, case conversion, and substring checking.

To run any of these examples:

```bash
go run valueobject/example_name.go
```

For example:

```bash
go run valueobject/money_example.go
```

### Auth Examples

The following examples demonstrate how to use various features of the auth package:

- **[quickstart_example.go](auth/quickstart_example.go)** - Demonstrates a quick start guide for using the Auth module.
- **[configuration_example.go](auth/configuration_example.go)** - Demonstrates how to configure the Auth module.
- **[auth_instance_example.go](auth/auth_instance_example.go)** - Demonstrates how to create an Auth instance.
- **[middleware_example.go](auth/middleware_example.go)** - Demonstrates how to use the Auth middleware.
- **[token_handling_example.go](auth/token_handling_example.go)** - Demonstrates how to generate and validate tokens.
- **[authorization_example.go](auth/authorization_example.go)** - Demonstrates how to perform authorization checks.
- **[user_info_example.go](auth/user_info_example.go)** - Demonstrates how to get user information from the context.
- **[context_utilities_example.go](auth/context_utilities_example.go)** - Demonstrates how to use context utilities.
- **[error_handling_example.go](auth/error_handling_example.go)** - Demonstrates how to handle errors.

To run any of these examples:

```bash
go run auth/example_name.go
```

For example:

```bash
go run auth/quickstart_example.go
```

### GraphQL Examples

The following examples demonstrate how to use various features of the graphql package:

- **[directive_registration_example.go](graphql/directive_registration_example.go)** - Demonstrates how to register the @isAuthorized directive in a GraphQL server.
- **[auth_configuration_example.go](graphql/auth_configuration_example.go)** - Demonstrates how to configure the auth service for GraphQL.
- **[auth_middleware_example.go](graphql/auth_middleware_example.go)** - Demonstrates how to apply the auth middleware to a GraphQL handler.
- **[resolver_authorization_example.go](graphql/resolver_authorization_example.go)** - Demonstrates how to check authorization in a GraphQL resolver.
- **[jwt_token_generation_example.go](graphql/jwt_token_generation_example.go)** - Demonstrates how to generate JWT tokens for testing GraphQL RBAC.

To run any of these examples:

```bash
go run graphql/example_name.go
```

For example:

```bash
go run graphql/resolver_authorization_example.go
```

### Package-Specific Examples

Examples for individual packages can be found in their respective README.md files:

- [Configuration Examples](../config/README.md)
- [Database Examples](../db/README.md)
- [Dependency Injection Examples](../di/README.md)
- [Health Check Examples](../health/README.md)
- [Logging Examples](../logging/README.md)
- [Telemetry Examples](../telemetry/README.md)
- [Transaction Examples](../transaction/README.md)

## Contributing Examples

If you'd like to contribute an example application, please follow these guidelines:

1. Create a new directory with a descriptive name for your example
2. Include a README.md that explains what the example demonstrates
3. Keep the example focused on demonstrating a specific use case
4. Include comments in the code to explain key concepts
5. Ensure the example follows best practices for Go code
6. Make sure the example can be run with minimal setup

## Running Examples

Each example should include instructions for running it in its README.md file.
