// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Color represents a color value object in hexadecimal format (#RRGGBB)
type Color string

// Regular expression for validating color format
var colorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)

// NewColor creates a new Color with validation
func NewColor(hexColor string) (Color, error) {
	// Trim whitespace
	trimmedColor := strings.TrimSpace(hexColor)

	// Empty color is allowed (optional field)
	if trimmedColor == "" {
		return "", nil
	}

	// Ensure color starts with #
	if !strings.HasPrefix(trimmedColor, "#") {
		trimmedColor = "#" + trimmedColor
	}

	// Validate color format
	if !colorRegex.MatchString(trimmedColor) {
		return "", errors.New("invalid color format, expected #RRGGBB or #RGB")
	}

	// Convert 3-digit hex to 6-digit hex if needed
	if len(trimmedColor) == 4 {
		r := string(trimmedColor[1])
		g := string(trimmedColor[2])
		b := string(trimmedColor[3])
		trimmedColor = "#" + r + r + g + g + b + b
	}

	return Color(strings.ToUpper(trimmedColor)), nil
}

// String returns the string representation of the Color
func (c Color) String() string {
	return string(c)
}

// Equals checks if two Colors are equal (case insensitive)
func (c Color) Equals(other Color) bool {
	return strings.EqualFold(string(c), string(other))
}

// IsEmpty checks if the Color is empty
func (c Color) IsEmpty() bool {
	return c == ""
}

// RGB returns the RGB components of the color
func (c Color) RGB() (r, g, b int, err error) {
	if c.IsEmpty() {
		return 0, 0, 0, errors.New("color is empty")
	}

	hexColor := string(c)[1:] // Remove #
	
	// Parse red component
	r64, err := strconv.ParseInt(hexColor[0:2], 16, 0)
	if err != nil {
		return 0, 0, 0, errors.New("invalid red component")
	}
	
	// Parse green component
	g64, err := strconv.ParseInt(hexColor[2:4], 16, 0)
	if err != nil {
		return 0, 0, 0, errors.New("invalid green component")
	}
	
	// Parse blue component
	b64, err := strconv.ParseInt(hexColor[4:6], 16, 0)
	if err != nil {
		return 0, 0, 0, errors.New("invalid blue component")
	}
	
	return int(r64), int(g64), int(b64), nil
}

// WithAlpha returns a new color with the specified alpha component (0-255)
func (c Color) WithAlpha(alpha int) (string, error) {
	if c.IsEmpty() {
		return "", errors.New("color is empty")
	}
	
	if alpha < 0 || alpha > 255 {
		return "", errors.New("alpha must be between 0 and 255")
	}
	
	return fmt.Sprintf("%s%02X", string(c), alpha), nil
}

// IsDark returns true if the color is considered dark
func (c Color) IsDark() (bool, error) {
	r, g, b, err := c.RGB()
	if err != nil {
		return false, err
	}
	
	// Calculate perceived brightness using the formula:
	// (0.299*R + 0.587*G + 0.114*B)
	brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	
	// If brightness is less than 128, the color is considered dark
	return brightness < 128, nil
}

// Invert returns the inverted color
func (c Color) Invert() (Color, error) {
	if c.IsEmpty() {
		return "", errors.New("color is empty")
	}
	
	r, g, b, err := c.RGB()
	if err != nil {
		return "", err
	}
	
	// Invert each component
	r = 255 - r
	g = 255 - g
	b = 255 - b
	
	// Create new color
	return Color(fmt.Sprintf("#%02X%02X%02X", r, g, b)), nil
}