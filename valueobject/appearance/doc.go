// Copyright (c) 2025 A Bit of Help, Inc.

// Package appearance provides value objects related to appearance information.
//
// This package contains value objects that represent visual appearance attributes
// such as colors, styles, and other visual properties. These value objects are
// immutable and follow the Value Object pattern from Domain-Driven Design.
//
// Key value objects in this package:
//   - Color: Represents a color in hexadecimal format (#RRGGBB)
//
// The Color value object provides methods for:
//   - Creating and validating colors in hexadecimal format
//   - Converting between different color formats
//   - Extracting RGB components
//   - Adding alpha channel information
//   - Determining if a color is dark or light
//   - Inverting colors
//
// Example usage:
//
//	// Create a new color
//	color, err := appearance.NewColor("#FF5500")
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Get RGB components
//	r, g, b, err := color.RGB()
//	if err != nil {
//	    // Handle error
//	}
//
//	// Check if color is dark
//	isDark, err := color.IsDark()
//	if err != nil {
//	    // Handle error
//	}
//
//	// Invert the color
//	invertedColor, err := color.Invert()
//	if err != nil {
//	    // Handle error
//	}
//
// All value objects in this package are designed to be immutable, so they cannot
// be changed after creation. To modify a value object, create a new instance with
// the desired values.
package appearance