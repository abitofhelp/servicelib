// Copyright (c) 2025 A Bit of Help, Inc.

// Package valueobject provides a collection of common Value Objects that can be used in applications.
package valueobject

import (
	"github.com/abitofhelp/servicelib/valueobject/measurement"
)

// Temperature represents a temperature value object
// This is a wrapper around measurement.Temperature for backward compatibility.
// New code should use measurement.Temperature directly.
type Temperature = measurement.Temperature

// TemperatureUnit represents a temperature unit
// This is a wrapper around measurement.TemperatureUnit for backward compatibility.
// New code should use measurement.TemperatureUnit directly.
type TemperatureUnit = measurement.TemperatureUnit

// Temperature unit constants
const (
	Celsius    = measurement.Celsius
	Fahrenheit = measurement.Fahrenheit
	Kelvin     = measurement.Kelvin
)

// NewTemperature creates a new Temperature with validation
// This function is provided for backward compatibility.
// New code should use measurement.NewTemperature directly.
func NewTemperature(value float64, unit TemperatureUnit) (Temperature, error) {
	return measurement.NewTemperature(value, unit)
}

// ParseTemperature creates a new Temperature from a string
// This function is provided for backward compatibility.
// New code should use measurement.ParseTemperature directly.
func ParseTemperature(s string) (Temperature, error) {
	return measurement.ParseTemperature(s)
}
