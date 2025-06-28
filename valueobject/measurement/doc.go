// Copyright (c) 2025 A Bit of Help, Inc.

// Package measurement provides value objects related to measurement information.
//
// This package contains value objects that represent different types of measurements
// such as monetary values, file sizes, memory sizes, percentages, ratings, and
// temperatures. These value objects are immutable and follow the Value Object pattern
// from Domain-Driven Design.
//
// Key value objects in this package:
//   - Money: Represents a monetary value with amount and currency
//   - FileSize: Represents a file size with different units (bytes, KB, MB, etc.)
//   - MemSize: Represents a memory size with different units (bytes, KB, MB, etc.)
//   - Percentage: Represents a percentage value
//   - Rating: Represents a rating value (e.g., 1-5 stars)
//   - Temperature: Represents a temperature value with different units
//
// Each value object provides methods for:
//   - Creating and validating instances
//   - String representation
//   - Equality comparison
//   - Conversion between different units or formats
//
// Many value objects also provide arithmetic operations appropriate to their domain:
//   - Money: Addition, subtraction, multiplication, division
//   - Percentage: Addition, subtraction, application to values
//   - Temperature: Conversion between units (Celsius, Fahrenheit, Kelvin)
//
// Example usage:
//
//	// Create a new money value
//	amount := decimal.NewFromFloat(99.99)
//	money, err := measurement.NewMoney(amount, "USD")
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Perform arithmetic operations
//	tax, err := measurement.NewMoney(decimal.NewFromFloat(8.50), "USD")
//	if err != nil {
//	    // Handle validation error
//	}
//	total, err := money.Add(tax)
//	if err != nil {
//	    // Handle error (e.g., currency mismatch)
//	}
//
//	// Create a percentage
//	discount, err := measurement.NewPercentage(15.0)
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Create a temperature
//	temp, err := measurement.NewTemperature(22.5, "C")
//	if err != nil {
//	    // Handle validation error
//	}
//	fahrenheit := temp.ToFahrenheit()
//
// All value objects in this package are designed to be immutable, so they cannot
// be changed after creation. To modify a value object, create a new instance with
// the desired values.
package measurement
