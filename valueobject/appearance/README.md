# Appearance Value Objects
This package contains value objects related to visual appearance and styling.

## Value Objects

### Color

The `Color` value object represents a color in hexadecimal format (#RRGGBB). It provides methods for:

- Converting to RGB components
- Adding an alpha channel
- Determining if a color is dark
- Inverting a color

## Usage

```go
// Create a new color
color, err := appearance.NewColor("#FF5733")
if err != nil {
    // Handle error
}

// Get RGB components
r, g, b, err := color.RGB()
if err != nil {
    // Handle error
}

// Check if color is dark
isDark, err := color.IsDark()
if err != nil {
    // Handle error
}

// Invert the color
inverted, err := color.Invert()
if err != nil {
    // Handle error
}

// Add alpha channel
colorWithAlpha, err := color.WithAlpha(128) // 50% opacity
if err != nil {
    // Handle error
}
```

## Overview

Brief description of the appearance and its purpose in the ServiceLib library.

## Features

- **Feature 1**: Description of feature 1
- **Feature 2**: Description of feature 2
- **Feature 3**: Description of feature 3

## Installation

```bash
go get github.com/abitofhelp/servicelib/appearance
```

## Quick Start

See the [Quick Start example](../EXAMPLES/appearance/quickstart_example.go) for a complete, runnable example of how to use the appearance.

## Configuration

See the [Configuration example](../EXAMPLES/appearance/configuration_example.go) for a complete, runnable example of how to configure the appearance.

## API Documentation


### Core Types

Description of the main types provided by the appearance.

#### Type 1

Description of Type 1 and its purpose.

See the [Type 1 example](../EXAMPLES/appearance/type1_example.go) for a complete, runnable example of how to use Type 1.

### Key Methods

Description of the key methods provided by the appearance.

#### Method 1

Description of Method 1 and its purpose.

See the [Method 1 example](../EXAMPLES/appearance/method1_example.go) for a complete, runnable example of how to use Method 1.

## Examples

For complete, runnable examples, see the following files in the EXAMPLES directory:

- [Basic Usage](../EXAMPLES/appearance/basic_usage_example.go) - Shows basic usage of the appearance
- [Advanced Configuration](../EXAMPLES/appearance/advanced_configuration_example.go) - Shows advanced configuration options
- [Error Handling](../EXAMPLES/appearance/error_handling_example.go) - Shows how to handle errors

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

- [Component 1](../appearance1/README.md) - Description of how this appearance relates to Component 1
- [Component 2](../appearance2/README.md) - Description of how this appearance relates to Component 2

## Contributing

Contributions to this appearance are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
