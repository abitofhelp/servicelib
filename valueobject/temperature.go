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

// TemperatureUnit represents a temperature unit
type TemperatureUnit string

const (
	// Celsius temperature unit (°C)
	Celsius TemperatureUnit = "C"
	// Fahrenheit temperature unit (°F)
	Fahrenheit TemperatureUnit = "F"
	// Kelvin temperature unit (K)
	Kelvin TemperatureUnit = "K"
)

// Temperature represents a temperature value object
type Temperature struct {
	value float64
	unit  TemperatureUnit
}

// Regular expression for parsing temperature strings
var temperatureRegex = regexp.MustCompile(`^([+-]?\d+(?:\.\d+)?)\s*°?([CFK])$`)

// NewTemperature creates a new Temperature with validation
func NewTemperature(value float64, unit TemperatureUnit) (Temperature, error) {
	// Validate unit
	if unit != Celsius && unit != Fahrenheit && unit != Kelvin {
		return Temperature{}, errors.New("invalid temperature unit, must be C, F, or K")
	}

	// Validate temperature range
	switch unit {
	case Celsius:
		// Absolute zero in Celsius is -273.15°C
		if value < -273.15 {
			return Temperature{}, errors.New("temperature cannot be below absolute zero (-273.15°C)")
		}
	case Fahrenheit:
		// Absolute zero in Fahrenheit is -459.67°F
		if value < -459.67 {
			return Temperature{}, errors.New("temperature cannot be below absolute zero (-459.67°F)")
		}
	case Kelvin:
		// Absolute zero in Kelvin is 0K
		if value < 0 {
			return Temperature{}, errors.New("temperature cannot be below absolute zero (0K)")
		}
	}

	// Round to 2 decimal places for precision
	roundedValue := math.Round(value*100) / 100

	return Temperature{
		value: roundedValue,
		unit:  unit,
	}, nil
}

// ParseTemperature creates a new Temperature from a string
func ParseTemperature(s string) (Temperature, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(s)

	// Empty string is not allowed
	if trimmed == "" {
		return Temperature{}, errors.New("temperature string cannot be empty")
	}

	// Match against regex
	matches := temperatureRegex.FindStringSubmatch(trimmed)
	if matches == nil {
		return Temperature{}, errors.New("invalid temperature format, expected '25°C', '77°F', or '298K'")
	}

	// Parse value
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return Temperature{}, errors.New("invalid temperature value")
	}

	// Parse unit
	unit := TemperatureUnit(matches[2])

	return NewTemperature(value, unit)
}

// String returns the string representation of the Temperature
func (t Temperature) String() string {
	if t.unit == Kelvin {
		return fmt.Sprintf("%.2f K", t.value)
	}
	return fmt.Sprintf("%.2f°%s", t.value, t.unit)
}

// Equals checks if two Temperatures are equal
// This compares the actual temperature values, not just the raw values and units
func (t Temperature) Equals(other Temperature) bool {
	// Convert both to Kelvin for comparison
	tKelvin := t.ToKelvin().value
	otherKelvin := other.ToKelvin().value

	// Compare with small epsilon to handle floating point precision
	epsilon := 0.01
	diff := math.Abs(tKelvin - otherKelvin)
	return diff < epsilon
}

// IsEmpty checks if the Temperature is empty (zero value)
func (t Temperature) IsEmpty() bool {
	return t.value == 0 && t.unit == ""
}

// Value returns the temperature value
func (t Temperature) Value() float64 {
	return t.value
}

// Unit returns the temperature unit
func (t Temperature) Unit() TemperatureUnit {
	return t.unit
}

// ToCelsius converts the temperature to Celsius
func (t Temperature) ToCelsius() Temperature {
	if t.unit == Celsius {
		return t
	}

	var celsiusValue float64
	switch t.unit {
	case Fahrenheit:
		celsiusValue = (t.value - 32) * 5 / 9
	case Kelvin:
		celsiusValue = t.value - 273.15
	default:
		return t // Should never happen due to validation
	}

	result, _ := NewTemperature(celsiusValue, Celsius)
	return result
}

// ToFahrenheit converts the temperature to Fahrenheit
func (t Temperature) ToFahrenheit() Temperature {
	if t.unit == Fahrenheit {
		return t
	}

	var fahrenheitValue float64
	switch t.unit {
	case Celsius:
		fahrenheitValue = (t.value * 9 / 5) + 32
	case Kelvin:
		fahrenheitValue = (t.value-273.15)*9/5 + 32
	default:
		return t // Should never happen due to validation
	}

	result, _ := NewTemperature(fahrenheitValue, Fahrenheit)
	return result
}

// ToKelvin converts the temperature to Kelvin
func (t Temperature) ToKelvin() Temperature {
	if t.unit == Kelvin {
		return t
	}

	var kelvinValue float64
	switch t.unit {
	case Celsius:
		kelvinValue = t.value + 273.15
	case Fahrenheit:
		kelvinValue = (t.value-32)*5/9 + 273.15
	default:
		return t // Should never happen due to validation
	}

	result, _ := NewTemperature(kelvinValue, Kelvin)
	return result
}

// Add adds another temperature and returns a new Temperature
func (t Temperature) Add(other Temperature) Temperature {
	// Convert other to the same unit as this temperature
	converted := other
	switch t.unit {
	case Celsius:
		converted = other.ToCelsius()
	case Fahrenheit:
		converted = other.ToFahrenheit()
	case Kelvin:
		converted = other.ToKelvin()
	}

	// Add the values
	result, _ := NewTemperature(t.value+converted.value, t.unit)
	return result
}

// Subtract subtracts another temperature and returns a new Temperature
func (t Temperature) Subtract(other Temperature) (Temperature, error) {
	// Convert other to the same unit as this temperature
	converted := other
	switch t.unit {
	case Celsius:
		converted = other.ToCelsius()
	case Fahrenheit:
		converted = other.ToFahrenheit()
	case Kelvin:
		converted = other.ToKelvin()
	}

	// Calculate the new value
	newValue := t.value - converted.value

	// Check if the result would be below absolute zero
	switch t.unit {
	case Celsius:
		if newValue < -273.15 {
			return Temperature{}, errors.New("result would be below absolute zero (-273.15°C)")
		}
	case Fahrenheit:
		if newValue < -459.67 {
			return Temperature{}, errors.New("result would be below absolute zero (-459.67°F)")
		}
	case Kelvin:
		if newValue < 0 {
			return Temperature{}, errors.New("result would be below absolute zero (0K)")
		}
	}

	// Subtract the values
	return NewTemperature(newValue, t.unit)
}

// IsFreezing checks if the temperature is at or below the freezing point of water
func (t Temperature) IsFreezing() bool {
	celsius := t.ToCelsius()
	return celsius.value <= 0
}

// IsBoiling checks if the temperature is at or above the boiling point of water
func (t Temperature) IsBoiling() bool {
	celsius := t.ToCelsius()
	return celsius.value >= 100
}

// Format returns the temperature in the specified format
// Format options:
// - "short": "25°C"
// - "long": "25 degrees Celsius"
// - "scientific": "298.15K"
func (t Temperature) Format(format string) string {
	switch format {
	case "short":
		if t.unit == Kelvin {
			return fmt.Sprintf("%.2fK", t.value)
		}
		return fmt.Sprintf("%.2f°%s", t.value, t.unit)

	case "long":
		unitName := ""
		switch t.unit {
		case Celsius:
			unitName = "degrees Celsius"
		case Fahrenheit:
			unitName = "degrees Fahrenheit"
		case Kelvin:
			unitName = "Kelvin"
		}
		return fmt.Sprintf("%.2f %s", t.value, unitName)

	case "scientific":
		kelvin := t.ToKelvin()
		return fmt.Sprintf("%.2fK", kelvin.value)

	default:
		return t.String()
	}
}
