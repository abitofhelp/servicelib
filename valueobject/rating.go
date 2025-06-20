// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Rating represents a rating value object (e.g., 1-5 stars)
type Rating struct {
	value    float64
	maxValue float64
}

// Regular expression for parsing rating strings
var ratingRegex = regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*(?:\/\s*(\d+(?:\.\d+)?))?$`)

// NewRating creates a new Rating with validation
func NewRating(value float64, maxValue float64) (Rating, error) {
	// Validate max value
	if maxValue <= 0 {
		return Rating{}, errors.New("maximum rating value must be positive")
	}

	// Validate value range
	if value < 0 {
		return Rating{}, errors.New("rating cannot be negative")
	}
	if value > maxValue {
		return Rating{}, fmt.Errorf("rating cannot exceed maximum value of %.2f", maxValue)
	}

	// Round to 1 decimal place for precision
	roundedValue := math.Round(value*10) / 10

	return Rating{
		value:    roundedValue,
		maxValue: maxValue,
	}, nil
}

// ParseRating creates a new Rating from a string
// Formats supported:
// - "4" (assumes max value of 5)
// - "4/5" (explicit max value)
// - "8/10" (explicit max value)
func ParseRating(s string) (Rating, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(s)

	// Empty string is not allowed
	if trimmed == "" {
		return Rating{}, errors.New("rating string cannot be empty")
	}

	// Match against regex
	matches := ratingRegex.FindStringSubmatch(trimmed)
	if matches == nil {
		return Rating{}, errors.New("invalid rating format, expected '4' or '4/5'")
	}

	// Parse value
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return Rating{}, errors.New("invalid rating value")
	}

	// Parse max value (default to 5 if not specified)
	maxValue := 5.0
	if len(matches) > 2 && matches[2] != "" {
		maxValue, err = strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return Rating{}, errors.New("invalid maximum rating value")
		}
	}

	return NewRating(value, maxValue)
}

// String returns the string representation of the Rating
func (r Rating) String() string {
	return fmt.Sprintf("%.1f/%.1f", r.value, r.maxValue)
}

// Equals checks if two Ratings are equal
// This compares the normalized values (value/maxValue)
func (r Rating) Equals(other Rating) bool {
	// Normalize both ratings to a 0-1 scale for comparison
	normalizedThis := r.value / r.maxValue
	normalizedOther := other.value / other.maxValue

	// Compare with small epsilon to handle floating point precision
	epsilon := 0.001
	diff := math.Abs(normalizedThis - normalizedOther)
	return diff < epsilon
}

// IsEmpty checks if the Rating is empty (zero value)
func (r Rating) IsEmpty() bool {
	return r.value == 0 && r.maxValue == 0
}

// Value returns the rating value
func (r Rating) Value() float64 {
	return r.value
}

// MaxValue returns the maximum rating value
func (r Rating) MaxValue() float64 {
	return r.maxValue
}

// Normalized returns the rating normalized to a 0-1 scale
func (r Rating) Normalized() float64 {
	if r.maxValue == 0 {
		return 0
	}
	return r.value / r.maxValue
}

// ToScale converts the rating to a different scale
func (r Rating) ToScale(newMaxValue float64) (Rating, error) {
	if newMaxValue <= 0 {
		return Rating{}, errors.New("new maximum value must be positive")
	}

	// Convert to the new scale
	newValue := r.Normalized() * newMaxValue

	return NewRating(newValue, newMaxValue)
}

// IsHigher checks if this rating is higher than another rating
// This compares the normalized values to handle different scales
func (r Rating) IsHigher(other Rating) bool {
	return r.Normalized() > other.Normalized()
}

// IsLower checks if this rating is lower than another rating
// This compares the normalized values to handle different scales
func (r Rating) IsLower(other Rating) bool {
	return r.Normalized() < other.Normalized()
}

// Percentage returns the rating as a percentage of the maximum value
func (r Rating) Percentage() float64 {
	return r.Normalized() * 100
}

// Stars returns the rating as a string of stars (★) and empty stars (☆)
// For example, a rating of 3.5/5 would return "★★★½☆"
func (r Rating) Stars() string {
	// Only works well for 5-star scale
	if r.maxValue != 5 {
		// Convert to 5-star scale first
		fiveStar, _ := r.ToScale(5)
		return fiveStar.Stars()
	}

	// Build the star string
	var result strings.Builder
	fullStars := int(r.value)
	hasHalfStar := r.value-float64(fullStars) >= 0.5

	// Add full stars
	for i := 0; i < fullStars; i++ {
		result.WriteString("★")
	}

	// Add half star if needed
	if hasHalfStar {
		result.WriteString("½")
		fullStars++
	}

	// Add empty stars
	for i := fullStars; i < int(r.maxValue); i++ {
		result.WriteString("☆")
	}

	return result.String()
}

// Format returns the rating in the specified format
// Format options:
// - "decimal": "4.5/5.0"
// - "percentage": "90%"
// - "stars": "★★★★½"
// - "fraction": "9/10"
func (r Rating) Format(format string) string {
	switch format {
	case "decimal":
		return fmt.Sprintf("%.1f/%.1f", r.value, r.maxValue)
	case "percentage":
		return fmt.Sprintf("%.0f%%", r.Percentage())
	case "stars":
		return r.Stars()
	case "fraction":
		// Convert to a fraction with denominator 10 for readability
		rating, _ := r.ToScale(10)
		return fmt.Sprintf("%d/10", int(math.Round(rating.value)))
	default:
		return r.String()
	}
}
