// Copyright (c) 2025 A Bit of Help, Inc.

// Package location provides value objects related to location information.
//
// This package contains value objects that represent geographical locations
// and coordinates. These value objects are immutable and follow the Value Object
// pattern from Domain-Driven Design.
//
// Key value objects in this package:
//   - Coordinate: Represents a geographic coordinate (latitude and longitude)
//
// The Coordinate value object provides methods for:
//   - Creating and validating coordinates
//   - Parsing coordinates from strings
//   - String representation and formatting in different formats (decimal, DMS, DM)
//   - Equality comparison
//   - Distance calculation between coordinates
//   - Directional comparison (north, south, east, west)
//   - JSON marshaling and conversion to maps
//
// Example usage:
//
//	// Create a new coordinate
//	coord, err := location.NewCoordinate(37.7749, -122.4194)
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Parse a coordinate from a string
//	coord, err := location.ParseCoordinate("37.7749, -122.4194")
//	if err != nil {
//	    // Handle parsing error
//	}
//
//	// Calculate distance between coordinates
//	distance := coord.DistanceTo(otherCoord)
//
//	// Format coordinate in different ways
//	decimalFormat := coord.Format("decimal")  // "37.7749, -122.4194"
//	dmsFormat := coord.Format("dms")          // "37° 46' 29.64" N, 122° 25' 9.84" W"
//
// All value objects in this package are designed to be immutable, so they cannot
// be changed after creation. To modify a value object, create a new instance with
// the desired values.
package location