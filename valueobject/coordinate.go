// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Coordinate represents a geographic coordinate (latitude and longitude)
type Coordinate struct {
	latitude  float64
	longitude float64
}

// NewCoordinate creates a new Coordinate with validation
func NewCoordinate(latitude, longitude float64) (Coordinate, error) {
	// Validate latitude range (-90 to 90 degrees)
	if latitude < -90 || latitude > 90 {
		return Coordinate{}, errors.New("latitude must be between -90 and 90 degrees")
	}

	// Validate longitude range (-180 to 180 degrees)
	if longitude < -180 || longitude > 180 {
		return Coordinate{}, errors.New("longitude must be between -180 and 180 degrees")
	}

	// Round to 6 decimal places for precision (approximately 10cm)
	roundedLat := math.Round(latitude*1000000) / 1000000
	roundedLng := math.Round(longitude*1000000) / 1000000

	return Coordinate{
		latitude:  roundedLat,
		longitude: roundedLng,
	}, nil
}

// ParseCoordinate creates a new Coordinate from a string in format "lat,lng"
func ParseCoordinate(s string) (Coordinate, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(s)

	// Empty string is not allowed
	if trimmed == "" {
		return Coordinate{}, errors.New("coordinate string cannot be empty")
	}

	// Split by comma
	parts := strings.Split(trimmed, ",")
	if len(parts) != 2 {
		return Coordinate{}, errors.New("invalid coordinate format, expected 'latitude,longitude'")
	}

	// Parse latitude
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return Coordinate{}, errors.New("invalid latitude format")
	}

	// Parse longitude
	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return Coordinate{}, errors.New("invalid longitude format")
	}

	return NewCoordinate(lat, lng)
}

// String returns the string representation of the Coordinate
func (c Coordinate) String() string {
	return fmt.Sprintf("%f,%f", c.latitude, c.longitude)
}

// Equals checks if two Coordinates are equal
func (c Coordinate) Equals(other Coordinate) bool {
	// Compare with small epsilon to handle floating point precision
	epsilon := 0.0000001
	latDiff := math.Abs(c.latitude - other.latitude)
	lngDiff := math.Abs(c.longitude - other.longitude)
	return latDiff < epsilon && lngDiff < epsilon
}

// IsEmpty checks if the Coordinate is empty (zero value)
func (c Coordinate) IsEmpty() bool {
	return c.latitude == 0 && c.longitude == 0
}

// Latitude returns the latitude value
func (c Coordinate) Latitude() float64 {
	return c.latitude
}

// Longitude returns the longitude value
func (c Coordinate) Longitude() float64 {
	return c.longitude
}

// DistanceTo calculates the distance to another coordinate in kilometers
// using the Haversine formula
func (c Coordinate) DistanceTo(other Coordinate) float64 {
	// Earth's radius in kilometers
	const earthRadius = 6371.0

	// Convert to radians
	lat1 := c.latitude * math.Pi / 180
	lng1 := c.longitude * math.Pi / 180
	lat2 := other.latitude * math.Pi / 180
	lng2 := other.longitude * math.Pi / 180

	// Haversine formula
	dlat := lat2 - lat1
	dlng := lng2 - lng1
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlng/2)*math.Sin(dlng/2)
	angle := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * angle

	return distance
}

// IsNorthOf checks if this coordinate is north of another coordinate
func (c Coordinate) IsNorthOf(other Coordinate) bool {
	return c.latitude > other.latitude
}

// IsSouthOf checks if this coordinate is south of another coordinate
func (c Coordinate) IsSouthOf(other Coordinate) bool {
	return c.latitude < other.latitude
}

// IsEastOf checks if this coordinate is east of another coordinate
func (c Coordinate) IsEastOf(other Coordinate) bool {
	return c.longitude > other.longitude
}

// IsWestOf checks if this coordinate is west of another coordinate
func (c Coordinate) IsWestOf(other Coordinate) bool {
	return c.longitude < other.longitude
}

// Format returns the coordinate in the specified format
// Format options:
// - "dms": Degrees, minutes, seconds (e.g., "41°24'12.2"N, 2°10'26.5"E")
// - "dm": Degrees and decimal minutes (e.g., "41°24.2033'N, 2°10.4417'E")
// - "dd": Decimal degrees (e.g., "41.40339, 2.17403")
func (c Coordinate) Format(format string) string {
	switch format {
	case "dms":
		return c.formatDMS()
	case "dm":
		return c.formatDM()
	case "dd":
		return fmt.Sprintf("%.5f, %.5f", c.latitude, c.longitude)
	default:
		return c.String()
	}
}

// formatDMS formats the coordinate in degrees, minutes, seconds
func (c Coordinate) formatDMS() string {
	// Format latitude
	latAbs := math.Abs(c.latitude)
	latDeg := int(latAbs)
	latMinFloat := (latAbs - float64(latDeg)) * 60
	latMin := int(latMinFloat)
	latSec := (latMinFloat - float64(latMin)) * 60
	latDir := "N"
	if c.latitude < 0 {
		latDir = "S"
	}

	// Format longitude
	lngAbs := math.Abs(c.longitude)
	lngDeg := int(lngAbs)
	lngMinFloat := (lngAbs - float64(lngDeg)) * 60
	lngMin := int(lngMinFloat)
	lngSec := (lngMinFloat - float64(lngMin)) * 60
	lngDir := "E"
	if c.longitude < 0 {
		lngDir = "W"
	}

	return fmt.Sprintf("%d°%d'%.1f\"%s, %d°%d'%.1f\"%s",
		latDeg, latMin, latSec, latDir,
		lngDeg, lngMin, lngSec, lngDir)
}

// formatDM formats the coordinate in degrees and decimal minutes
func (c Coordinate) formatDM() string {
	// Format latitude
	latAbs := math.Abs(c.latitude)
	latDeg := int(latAbs)
	latMin := (latAbs - float64(latDeg)) * 60
	latDir := "N"
	if c.latitude < 0 {
		latDir = "S"
	}

	// Format longitude
	lngAbs := math.Abs(c.longitude)
	lngDeg := int(lngAbs)
	lngMin := (lngAbs - float64(lngDeg)) * 60
	lngDir := "E"
	if c.longitude < 0 {
		lngDir = "W"
	}

	return fmt.Sprintf("%d°%.4f'%s, %d°%.4f'%s",
		latDeg, latMin, latDir,
		lngDeg, lngMin, lngDir)
}
