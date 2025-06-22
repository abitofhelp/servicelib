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