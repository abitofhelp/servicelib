# Value Object Command Line Tools

## Overview

The Value Object Command Line Tools component provides command-line utilities for working with value objects in the ServiceLib library.

## Features

- **Code Generation**: Tools for generating value object code
- **Validation**: Tools for validating value object implementations
- **Documentation**: Tools for generating documentation for value objects
- **Testing**: Tools for generating tests for value objects

## Installation

```bash
go get github.com/abitofhelp/servicelib/valueobject/cmd
```

## Quick Start

See the [Quick Start example](../../EXAMPLES/valueobject/cmd/README.md) for a complete, runnable example of how to use the value object command line tools.

## API Documentation

### Core Types

Description of the main types provided by the component.

#### Command

Represents a command-line command.

```
type Command struct {
    // Fields
}
```

#### Option

Represents a command-line option.

```
type Option struct {
    // Fields
}
```

### Key Methods

Description of the key methods provided by the component.

#### NewCommand

Creates a new command.

```
func NewCommand(name string, options ...Option) (*Command, error)
```

#### Execute

Executes a command.

```
func (c *Command) Execute(args []string) error
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Generate](../../EXAMPLES/valueobject/cmd/generate/README.md) - Shows how to generate value objects
- [Validate](../../EXAMPLES/valueobject/cmd/validate/README.md) - Shows how to validate value objects

## Best Practices

1. **Use Code Generation**: Use the code generation tools to ensure consistent value object implementations
2. **Validate Generated Code**: Always validate generated code before using it
3. **Document Value Objects**: Use the documentation tools to generate documentation for value objects
4. **Test Value Objects**: Use the testing tools to generate tests for value objects
5. **Customize Templates**: Customize the code generation templates to match your project's style

## Troubleshooting

### Common Issues

#### Code Generation Errors

If you encounter errors during code generation, check that the input schema is valid and that the templates are correctly formatted.

#### Command Execution Errors

If commands fail to execute, check that the command-line arguments are correct and that the required files exist.

## Related Components

- [Value Object](../README.md) - The main value object component
- [Value Object Generator](../generator/README.md) - The value object generator component

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.