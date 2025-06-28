// Copyright (c) 2025 A Bit of Help, Inc.

// Package location provides value objects related to location information.
package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/abitofhelp/servicelib/valueobject/base"
)

// Coordinate represents a geographic coordinate (latitude and longitude)
type Coordinate struct {
	base.BaseStructValueObject
	latitude  float64
	longitude float64
}

// NewCoordinate creates a new Coordinate with validation
func NewCoordinate(latitude float64, longitude float64) (Coordinate, error) {
	// Round to 6 decimal places for precision (approximately 10cm)
	roundedLat := math.Round(latitude*1000000) / 1000000
	roundedLng := math.Round(longitude*1000000) / 1000000

	vo := Coordinate{
		latitude:  roundedLat,
		longitude: roundedLng,
	}

	// Validate the value object
	if err := vo.Validate(); err != nil {
		return Coordinate{}, err
	}

	return vo, nil
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
func (v Coordinate) String() string {
	return fmt.Sprintf("latitude=%f longitude=%f", v.latitude, v.longitude)
}

// Equals checks if two Coordinates are equal
func (v Coordinate) Equals(other Coordinate) bool {
	// Compare with small epsilon to handle floating point precision
	epsilon := 0.0000001
	latDiff := math.Abs(v.latitude - other.latitude)
	lngDiff := math.Abs(v.longitude - other.longitude)
	return latDiff < epsilon && lngDiff < epsilon
}

// IsEmpty checks if the Coordinate is empty (zero value)
func (v Coordinate) IsEmpty() bool {
	return v.latitude == 0 && v.longitude == 0
}

// Validate checks if the Coordinate is valid
func (v Coordinate) Validate() error {
	if v.latitude < -90 || v.latitude > 90 {
		return errors.New("latitude must be between -90 and 90 degrees")
	}

	if v.longitude < -180 || v.longitude > 180 {
		return errors.New("longitude must be between -180 and 180 degrees")
	}

	return nil
}

// ToMap converts the Coordinate to a map[string]interface{}
func (v Coordinate) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"latitude":  v.latitude,
		"longitude": v.longitude,
	}
}

// MarshalJSON implements the json.Marshaler interface
func (v Coordinate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.ToMap())
}

// Latitude returns the latitude value
func (v Coordinate) Latitude() float64 {
	return v.latitude
}

// Longitude returns the longitude value
func (v Coordinate) Longitude() float64 {
	return v.longitude
}

// DistanceTo calculates the distance to another coordinate in kilometers
// using the Haversine formula
func (v Coordinate) DistanceTo(other Coordinate) float64 {
	// Earth's radius in kilometers
	const earthRadius = 6371.0

	// Convert to radians
	lat1 := v.latitude * math.Pi / 180
	lng1 := v.longitude * math.Pi / 180
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
func (v Coordinate) IsNorthOf(other Coordinate) bool {
	return v.latitude > other.latitude
}

// IsSouthOf checks if this coordinate is south of another coordinate
func (v Coordinate) IsSouthOf(other Coordinate) bool {
	return v.latitude < other.latitude
}

// IsEastOf checks if this coordinate is east of another coordinate
func (v Coordinate) IsEastOf(other Coordinate) bool {
	return v.longitude > other.longitude
}

// IsWestOf checks if this coordinate is west of another coordinate
func (v Coordinate) IsWestOf(other Coordinate) bool {
	return v.longitude < other.longitude
}

// Format returns the coordinate in the specified format
// Format options:
// - "dms": Degrees, minutes, seconds (e.g., "41°24'12.2"N, 2°10'26.5"E")
// - "dm": Degrees and decimal minutes (e.g., "41°24.2033'N, 2°10.4417'E")
// - "dd": Decimal degrees (e.g., "41.40339, 2.17403")
func (v Coordinate) Format(format string) string {
	switch format {
	case "dms":
		return v.formatDMS()
	case "dm":
		return v.formatDM()
	case "dd":
		return fmt.Sprintf("%.5f, %.5f", v.latitude, v.longitude)
	default:
		return v.String()
	}
}

// formatDMS formats the coordinate in degrees, minutes, seconds
func (v Coordinate) formatDMS() string {
	// Format latitude
	latAbs := math.Abs(v.latitude)
	latDeg := int(latAbs)
	latMinFloat := (latAbs - float64(latDeg)) * 60
	latMin := int(latMinFloat)
	latSec := (latMinFloat - float64(latMin)) * 60
	latDir := "N"
	if v.latitude < 0 {
		latDir = "S"
	}

	// Format longitude
	lngAbs := math.Abs(v.longitude)
	lngDeg := int(lngAbs)
	lngMinFloat := (lngAbs - float64(lngDeg)) * 60
	lngMin := int(lngMinFloat)
	lngSec := (lngMinFloat - float64(lngMin)) * 60
	lngDir := "E"
	if v.longitude < 0 {
		lngDir = "W"
	}

	return fmt.Sprintf("%d°%d'%.1f\"%s, %d°%d'%.1f\"%s",
		latDeg, latMin, latSec, latDir,
		lngDeg, lngMin, lngSec, lngDir)
}

// formatDM formats the coordinate in degrees and decimal minutes
func (v Coordinate) formatDM() string {
	// Format latitude
	latAbs := math.Abs(v.latitude)
	latDeg := int(latAbs)
	latMin := (latAbs - float64(latDeg)) * 60
	latDir := "N"
	if v.latitude < 0 {
		latDir = "S"
	}

	// Format longitude
	lngAbs := math.Abs(v.longitude)
	lngDeg := int(lngAbs)
	lngMin := (lngAbs - float64(lngDeg)) * 60
	lngDir := "E"
	if v.longitude < 0 {
		lngDir = "W"
	}

	return fmt.Sprintf("%d°%.4f'%s, %d°%.4f'%s",
		latDeg, latMin, latDir,
		lngDeg, lngMin, lngDir)
}
