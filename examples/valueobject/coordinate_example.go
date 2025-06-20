// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the Coordinate value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new coordinate
	coord, err := valueobject.NewCoordinate(40.7128, -74.0060)
	if err != nil {
		// Handle error
		fmt.Println("Error creating coordinate:", err)
		return
	}

	// Access latitude and longitude
	lat := coord.Latitude()
	lng := coord.Longitude()
	fmt.Printf("Latitude: %.4f, Longitude: %.4f\n", lat, lng)

	// Calculate distance to another coordinate
	otherCoord, _ := valueobject.NewCoordinate(34.0522, -118.2437)
	distance := coord.DistanceTo(otherCoord)
	fmt.Printf("Distance to Los Angeles: %.2f km\n", distance)

	// Check if coordinate is in northern hemisphere
	isNorthern := coord.IsNorthernHemisphere()
	fmt.Printf("Is in northern hemisphere? %v\n", isNorthern)

	// Check if coordinate is in western hemisphere
	isWestern := coord.IsWesternHemisphere()
	fmt.Printf("Is in western hemisphere? %v\n", isWestern)

	// Format as string
	fmt.Println(coord.String()) // "40.7128, -74.0060"
}
