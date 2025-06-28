// Copyright (c) 2025 A Bit of Help, Inc.

// Package measurement provides value objects related to measurement information.
package measurement

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/abitofhelp/servicelib/valueobject/base"
)

// Rating represents a rating value object (e.g., 1-5 stars)
type Rating struct {
	base.BaseStructValueObject

	maxValue float64

	value float64
}

// NewRating creates a new Rating with validation
func NewRating(value float64, maxValue float64) (Rating, error) {
	vo := Rating{
		value:    value,
		maxValue: maxValue,
	}

	// Validate the value object
	if err := vo.Validate(); err != nil {
		return Rating{}, err
	}

	return vo, nil
}

// String returns the string representation of the Rating
func (v Rating) String() string {
	return fmt.Sprintf("value=%.1f maxValue=%.1f", v.value, v.maxValue)
}

// Equals checks if two Ratings are equal
func (v Rating) Equals(other Rating) bool {

	if !base.FloatsEqual(v.maxValue, other.maxValue) {
		return false
	}

	if !base.FloatsEqual(v.value, other.value) {
		return false
	}

	return true
}

// IsEmpty checks if the Rating is empty (zero value)
func (v Rating) IsEmpty() bool {
	return v.maxValue == float64(0) && v.value == float64(0)
}

// Validate checks if the Rating is valid
func (v Rating) Validate() error {

	// Validate max value
	if v.maxValue <= 0 {
		return errors.New("maximum rating value must be positive")
	}

	// Validate value range
	if v.value < 0 {
		return errors.New("rating cannot be negative")
	}
	if v.value > v.maxValue {
		return fmt.Errorf("rating cannot exceed maximum value of %.2f", v.maxValue)
	}

	return nil
}

// ToMap converts the Rating to a map[string]interface{}
func (v Rating) ToMap() map[string]interface{} {
	return map[string]interface{}{

		"maxValue": v.maxValue,

		"value": v.value,
	}
}

// MarshalJSON implements the json.Marshaler interface
func (v Rating) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.ToMap())
}

// Value returns the rating value
func (v Rating) Value() float64 {
	return v.value
}

// MaxValue returns the maximum rating value
func (v Rating) MaxValue() float64 {
	return v.maxValue
}

// Normalized returns the rating normalized to a 0-1 scale
func (v Rating) Normalized() float64 {
	if v.maxValue == 0 {
		return 0
	}
	return v.value / v.maxValue
}

// ToScale converts the rating to a different scale
func (v Rating) ToScale(newMaxValue float64) (Rating, error) {
	if newMaxValue <= 0 {
		return Rating{}, errors.New("new maximum value must be positive")
	}

	// Convert to the new scale
	newValue := v.Normalized() * newMaxValue

	return NewRating(newValue, newMaxValue)
}

// IsHigher checks if this rating is higher than another rating
// This compares the normalized values to handle different scales
func (v Rating) IsHigher(other Rating) bool {
	return v.Normalized() > other.Normalized()
}

// IsLower checks if this rating is lower than another rating
// This compares the normalized values to handle different scales
func (v Rating) IsLower(other Rating) bool {
	return v.Normalized() < other.Normalized()
}

// Percentage returns the rating as a percentage of the maximum value
func (v Rating) Percentage() float64 {
	return v.Normalized() * 100
}

// Stars returns the rating as a string of stars (★) and empty stars (☆)
// For example, a rating of 3.5/5 would return "★★★½☆"
func (v Rating) Stars() string {
	// Only works well for 5-star scale
	if v.maxValue != 5 {
		// Convert to 5-star scale first
		fiveStar, _ := v.ToScale(5)
		return fiveStar.Stars()
	}

	// Build the star string
	var result strings.Builder
	fullStars := int(v.value)
	hasHalfStar := v.value-float64(fullStars) >= 0.5

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
	for i := fullStars; i < int(v.maxValue); i++ {
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
func (v Rating) Format(format string) string {
	switch format {
	case "decimal":
		return fmt.Sprintf("%.1f/%.1f", v.value, v.maxValue)
	case "percentage":
		return fmt.Sprintf("%.0f%%", v.Percentage())
	case "stars":
		return v.Stars()
	case "fraction":
		// Convert to a fraction with denominator 10 for readability
		rating, _ := v.ToScale(10)
		return fmt.Sprintf("%d/10", int(math.Round(rating.value)))
	default:
		return v.String()
	}
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
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*(?:\/\s*(\d+(?:\.\d+)?))?$`)
	matches := re.FindStringSubmatch(trimmed)
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
