# Color Value Object Example

## Overview

This example demonstrates the Color value object functionality of the ServiceLib valueobject/appearance package. It shows how to create, manipulate, and compare color values in a type-safe manner.

## Features

- **Color Creation**: Creating color objects from hex strings
- **Color Manipulation**: Converting, inverting, and adding alpha to colors
- **Color Properties**: Accessing RGB components and checking if a color is dark
- **Color Comparison**: Comparing colors for equality

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Creating Colors

The example demonstrates how to create Color value objects from hex strings:

```
// Create a new color from a standard hex string
color, err := appearance.NewColor("#FF5733")
if err != nil {
    // Handle error
    fmt.Println("Error creating color:", err)
    return
}

// Create color from short hex
shortColor, err := appearance.NewColor("#F00")
if err != nil {
    fmt.Println("Error creating color from short hex:", err)
    return
}
```

### Accessing Color Properties

The example shows how to access various properties of a color:

```
// Access color as string
fmt.Printf("Color: %s\n", color.String())

// Get RGB components
r, g, b, err := color.RGB()
if err != nil {
    fmt.Println("Error getting RGB components:", err)
    return
}
fmt.Printf("RGB: (%d, %d, %d)\n", r, g, b)

// Check if color is dark
isDark, err := color.IsDark()
if err != nil {
    fmt.Println("Error checking if color is dark:", err)
    return
}
fmt.Printf("Is dark? %v\n", isDark)
```

### Manipulating Colors

The example demonstrates how to manipulate colors:

```
// Add alpha component
colorWithAlpha, err := color.WithAlpha(128)
if err != nil {
    fmt.Println("Error adding alpha component:", err)
    return
}
fmt.Printf("Color with alpha: %s\n", colorWithAlpha)

// Invert color
invertedColor, err := color.Invert()
if err != nil {
    fmt.Println("Error inverting color:", err)
    return
}
fmt.Printf("Inverted color: %s\n", invertedColor)
```

## Expected Output

```
Color: #FF5733
RGB: (255, 87, 51)
Is dark? false
Color with alpha: #FF573380
Inverted color: #00A8CC
Color from #F00: #FF0000
Colors are equal? true
```

## Related Examples

- [Email Example](../../contact/email_example/README.md) - Example of the Email value object
- [UUID Example](../../identification/uuid_example/README.md) - Example of the UUID value object
- [Money Example](../../measurement/money_example/README.md) - Example of the Money value object

## Related Components

- [ValueObject Package](../../../../valueobject/README.md) - The value object package documentation
- [Validation Package](../../../../validation/README.md) - The validation package used for validating value objects
- [Errors Package](../../../../errors/README.md) - The errors package used for error handling

## License

This project is licensed under the MIT License - see the [LICENSE](../../../../LICENSE) file for details.