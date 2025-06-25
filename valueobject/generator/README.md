# Value Object Generator

## Overview

The Value Object Generator is a tool for generating value objects in the ServiceLib library. It generates code for both string-based and struct-based value objects, along with their corresponding test files.

## Usage

### Command-Line Tool

The generator can be used as a command-line tool:

```bash
go run cmd/generate/main.go -config path/to/config.json -templates path/to/templates -output path/to/output
```

Options:
- `-config`: Path to the configuration file (required)
- `-templates`: Directory containing the templates (default: "templates")
- `-output`: Directory where the generated code will be written (default: ".")

### Configuration File

The configuration file is a JSON file that contains an array of value object configurations. Each configuration has the following fields:

- `Name`: The name of the value object (e.g., "Email")
- `Package`: The package name (e.g., "contact")
- `Type`: The type of value object ("string" or "struct")
- `BaseType`: The base type for string-based value objects (e.g., "string")
- `Description`: A description of the value object
- `Fields`: A map of field names to field types (for struct-based value objects)
- `Validations`: A map of field names to validation functions
- `Imports`: A list of additional imports

Example configuration file:

```json
[
  {
    "Name": "PostalCode",
    "Package": "contact",
    "Type": "string",
    "BaseType": "string",
    "Description": "a postal code value object",
    "Validations": {
      "PostalCode": "if len(trimmedValue) < 5 {\n\t\treturn \"\", errors.New(\"postal code is too short\")\n\t}\n\n\tif len(trimmedValue) > 10 {\n\t\treturn \"\", errors.New(\"postal code is too long\")\n\t}"
    },
    "Imports": []
  },
  {
    "Name": "Currency",
    "Package": "measurement",
    "Type": "struct",
    "Description": "a currency value object",
    "Fields": {
      "Code": "string",
      "Symbol": "string",
      "Name": "string"
    },
    "Validations": {
      "Code": "if v.Code == \"\" {\n\t\treturn errors.New(\"currency code cannot be empty\")\n\t}\n\n\tif len(v.Code) != 3 {\n\t\treturn errors.New(\"currency code must be 3 characters\")\n\t}",
      "Symbol": "if v.Symbol == \"\" {\n\t\treturn errors.New(\"currency symbol cannot be empty\")\n\t}"
    },
    "Imports": [
      "errors"
    ]
  }
]
```

### Programmatic Usage

The generator can also be used programmatically:

```go
import "github.com/abitofhelp/servicelib/valueobject/generator"

func main() {
    // Create a generator
    gen := generator.NewGenerator("templates", "output")

    // Create a value object configuration
    config := generator.ValueObjectConfig{
        Name:        "Email",
        Package:     "contact",
        Type:        generator.StringBased,
        BaseType:    "string",
        Description: "an email address value object",
        Validations: map[string]string{
            "Email": `if !strings.Contains(trimmedValue, "@") {
                return "", errors.New("invalid email format: missing @ symbol")
            }`,
        },
        Imports: []string{"errors"},
    }

    // Generate the value object
    if err := gen.Generate(config); err != nil {
        panic(err)
    }
}
```

## Templates

The generator uses the following templates:

- `string_value_object.tmpl`: Template for string-based value objects
- `string_value_object_test.tmpl`: Template for string-based value object tests
- `struct_value_object.tmpl`: Template for struct-based value objects
- `struct_value_object_test.tmpl`: Template for struct-based value object tests

These templates are located in the `templates` directory.

## Generated Code

The generator generates the following files for each value object:

- `<package>/<name>.go`: The value object implementation
- `<package>/<name>_test.go`: The value object tests

The generated code follows the ServiceLib coding standards and includes:

- A constructor function that validates the input and returns a new instance
- Methods for String(), Equals(), IsEmpty(), and Validate()
- Additional methods for struct-based value objects: ToMap() and MarshalJSON()
- Comprehensive tests for all methods

## Features

- **Feature 1**: Description of feature 1
- **Feature 2**: Description of feature 2
- **Feature 3**: Description of feature 3

## Installation

```bash
go get github.com/abitofhelp/servicelib/generator
```

## Quick Start

See the [Quick Start example](../EXAMPLES/generator/quickstart_example.go) for a complete, runnable example of how to use the generator.

## Configuration

See the [Configuration example](../EXAMPLES/generator/configuration_example.go) for a complete, runnable example of how to configure the generator.

## API Documentation


### Core Types

Description of the main types provided by the generator.

#### Type 1

Description of Type 1 and its purpose.

See the [Type 1 example](../EXAMPLES/generator/type1_example.go) for a complete, runnable example of how to use Type 1.

### Key Methods

Description of the key methods provided by the generator.

#### Method 1

Description of Method 1 and its purpose.

See the [Method 1 example](../EXAMPLES/generator/method1_example.go) for a complete, runnable example of how to use Method 1.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/generator/basic_usage_example.go) - Shows basic usage of the generator
- [Advanced Configuration](../EXAMPLES/generator/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/generator/error_handling_example.go) - Shows how to handle errors

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

- [Component 1](../generator1/README.md) - Description of how this generator relates to Component 1
- [Component 2](../generator2/README.md) - Description of how this generator relates to Component 2

## Contributing

Contributions to this generator are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
