# Value Object Generator

## Overview

The Value Object Generator is a command-line tool for generating value objects based on templates and configuration files. It automates the creation of value objects following the Value Object pattern from Domain-Driven Design.

## Features

- **Template-Based Generation**: Generate value objects based on predefined templates
- **Configuration-Driven**: Customize value objects using JSON configuration files
- **Consistent Patterns**: Ensure all value objects follow the same patterns and conventions
- **Boilerplate Reduction**: Reduce repetitive code when creating new value objects
- **Development Tool**: Streamline the development process for creating value objects

## Installation

```bash
go get github.com/abitofhelp/servicelib/valueobject/cmd/generate
```

## Quick Start

To generate a value object, run the command:

```bash
go run main.go [flags]
```

## API Documentation

### Usage

The tool reads configuration from JSON files that specify the properties and behavior of the value objects to be generated. It then creates Go files with the appropriate code for these value objects.

```
go run main.go [flags]
```

## Examples

For examples of value objects generated with this tool, see the following directories:

- [Value Object Examples](../../../EXAMPLES/valueobject/README.md) - Examples of various value objects

## Best Practices

1. **Use Configuration Files**: Store your value object configurations in JSON files
2. **Follow Naming Conventions**: Use consistent naming for your value objects
3. **Generate Related Objects Together**: Generate related value objects in the same package
4. **Review Generated Code**: Always review the generated code before committing
5. **Keep Templates Updated**: Update templates when patterns or conventions change

## Troubleshooting

### Common Issues

#### Generation Failures

If generation fails, check that your configuration file is valid JSON and contains all required fields.

## Related Components

- [Value Object Package](../../README.md) - The parent value object package
- [Value Object Generator](../../generator/README.md) - The generator package used by this tool
- [Value Object Templates](../../generator/templates/README.md) - Templates used for code generation

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.