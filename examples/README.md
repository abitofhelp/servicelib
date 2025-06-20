# ServiceLib Examples

This directory contains complete example applications that demonstrate how to use ServiceLib in real-world scenarios.

## Available Examples

### Standalone Examples

- **money_example.go** - Demonstrates how to use the Money value object from the valueobject package, including creating Money objects, performing arithmetic operations, comparing values, parsing from strings, and handling precision issues.

To run the money example:

```bash
go run money_example.go
```

### Package-Specific Examples

Examples for individual packages can be found in their respective README.md files:

- [Authentication Examples](../auth/README.md)
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
