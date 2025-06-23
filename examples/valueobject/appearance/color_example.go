// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the Color value object
package appearance

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject/appearance"
)

func main() {
	// Create a new color
	color, err := appearance.NewColor("#FF5733")
	if err != nil {
		// Handle error
		fmt.Println("Error creating color:", err)
		return
	}

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

	// Create color from short hex
	shortColor, err := appearance.NewColor("#F00")
	if err != nil {
		fmt.Println("Error creating color from short hex:", err)
		return
	}
	fmt.Printf("Color from #F00: %s\n", shortColor)

	// Compare colors
	anotherColor, _ := appearance.NewColor("#ff5733")
	areEqual := color.Equals(anotherColor)
	fmt.Printf("Colors are equal? %v\n", areEqual) // true (case insensitive)
}