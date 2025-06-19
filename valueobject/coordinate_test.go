// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"math"
	"strings"
	"testing"
)

func TestNewCoordinate(t *testing.T) {
	tests := []struct {
		name        string
		latitude    float64
		longitude   float64
		expectError bool
	}{
		{"Valid Coordinate", 40.7128, -74.0060, false},
		{"Zero Coordinate", 0, 0, false},
		{"Max Values", 90, 180, false},
		{"Min Values", -90, -180, false},
		{"Invalid Latitude High", 91, 0, true},
		{"Invalid Latitude Low", -91, 0, true},
		{"Invalid Longitude High", 0, 181, true},
		{"Invalid Longitude Low", 0, -181, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coord, err := NewCoordinate(tt.latitude, tt.longitude)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Check rounding to 6 decimal places
				expectedLat := math.Round(tt.latitude*1000000) / 1000000
				expectedLng := math.Round(tt.longitude*1000000) / 1000000

				if coord.latitude != expectedLat {
					t.Errorf("Expected latitude %f, got %f", expectedLat, coord.latitude)
				}
				if coord.longitude != expectedLng {
					t.Errorf("Expected longitude %f, got %f", expectedLng, coord.longitude)
				}
			}
		})
	}
}

func TestParseCoordinate(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedLat float64
		expectedLng float64
		expectError bool
	}{
		{"Valid Coordinate", "40.7128,-74.0060", 40.7128, -74.0060, false},
		{"Valid Coordinate with Spaces", " 40.7128 , -74.0060 ", 40.7128, -74.0060, false},
		{"Zero Coordinate", "0,0", 0, 0, false},
		{"Empty String", "", 0, 0, true},
		{"Invalid Format", "40.7128", 0, 0, true},
		{"Invalid Format with Extra Parts", "40.7128,-74.0060,0", 0, 0, true},
		{"Invalid Latitude", "abc,-74.0060", 0, 0, true},
		{"Invalid Longitude", "40.7128,abc", 0, 0, true},
		{"Invalid Range", "91,0", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coord, err := ParseCoordinate(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Check rounding to 6 decimal places
				expectedLat := math.Round(tt.expectedLat*1000000) / 1000000
				expectedLng := math.Round(tt.expectedLng*1000000) / 1000000

				if coord.latitude != expectedLat {
					t.Errorf("Expected latitude %f, got %f", expectedLat, coord.latitude)
				}
				if coord.longitude != expectedLng {
					t.Errorf("Expected longitude %f, got %f", expectedLng, coord.longitude)
				}
			}
		})
	}
}

func TestCoordinate_String(t *testing.T) {
	coord, _ := NewCoordinate(40.7128, -74.0060)
	expected := "40.712800,-74.006000"

	if coord.String() != expected {
		t.Errorf("Expected string %s, got %s", expected, coord.String())
	}
}

func TestCoordinate_Equals(t *testing.T) {
	coord1, _ := NewCoordinate(40.7128, -74.0060)
	coord2, _ := NewCoordinate(40.7128, -74.0060)
	coord3, _ := NewCoordinate(40.7129, -74.0060) // Slightly different latitude
	coord4, _ := NewCoordinate(40.7128, -74.0061) // Slightly different longitude

	if !coord1.Equals(coord2) {
		t.Errorf("Expected coordinates to be equal")
	}

	if coord1.Equals(coord3) {
		t.Errorf("Expected coordinates with different latitudes to not be equal")
	}

	if coord1.Equals(coord4) {
		t.Errorf("Expected coordinates with different longitudes to not be equal")
	}

	// Test with very small differences (within epsilon)
	coordSmallDiff, _ := NewCoordinate(40.7128000001, -74.0060000001)
	if !coord1.Equals(coordSmallDiff) {
		t.Errorf("Expected coordinates with very small differences to be equal")
	}
}

func TestCoordinate_IsEmpty(t *testing.T) {
	emptyCoord := Coordinate{}
	coord, _ := NewCoordinate(40.7128, -74.0060)

	if !emptyCoord.IsEmpty() {
		t.Errorf("Expected empty coordinate to be empty")
	}

	if coord.IsEmpty() {
		t.Errorf("Expected non-empty coordinate to not be empty")
	}
}

func TestCoordinate_Latitude(t *testing.T) {
	coord, _ := NewCoordinate(40.7128, -74.0060)

	if coord.Latitude() != 40.7128 {
		t.Errorf("Expected latitude 40.7128, got %f", coord.Latitude())
	}
}

func TestCoordinate_Longitude(t *testing.T) {
	coord, _ := NewCoordinate(40.7128, -74.0060)

	if coord.Longitude() != -74.0060 {
		t.Errorf("Expected longitude -74.0060, got %f", coord.Longitude())
	}
}

func TestCoordinate_DistanceTo(t *testing.T) {
	// New York City
	nyc, _ := NewCoordinate(40.7128, -74.0060)

	// Los Angeles
	la, _ := NewCoordinate(34.0522, -118.2437)

	// Distance between NYC and LA is approximately 3935 km
	distance := nyc.DistanceTo(la)
	expectedDistance := 3935.0

	// Allow for some margin of error
	if math.Abs(distance-expectedDistance) > 10 {
		t.Errorf("Expected distance around %f km, got %f km", expectedDistance, distance)
	}

	// Distance to self should be 0
	selfDistance := nyc.DistanceTo(nyc)
	if selfDistance != 0 {
		t.Errorf("Expected distance to self to be 0, got %f", selfDistance)
	}
}

func TestCoordinate_Directions(t *testing.T) {
	// Reference point
	ref, _ := NewCoordinate(0, 0)

	// North
	north, _ := NewCoordinate(1, 0)
	if !north.IsNorthOf(ref) {
		t.Errorf("Expected (1,0) to be north of (0,0)")
	}
	if north.IsSouthOf(ref) {
		t.Errorf("Expected (1,0) to not be south of (0,0)")
	}

	// South
	south, _ := NewCoordinate(-1, 0)
	if !south.IsSouthOf(ref) {
		t.Errorf("Expected (-1,0) to be south of (0,0)")
	}
	if south.IsNorthOf(ref) {
		t.Errorf("Expected (-1,0) to not be north of (0,0)")
	}

	// East
	east, _ := NewCoordinate(0, 1)
	if !east.IsEastOf(ref) {
		t.Errorf("Expected (0,1) to be east of (0,0)")
	}
	if east.IsWestOf(ref) {
		t.Errorf("Expected (0,1) to not be west of (0,0)")
	}

	// West
	west, _ := NewCoordinate(0, -1)
	if !west.IsWestOf(ref) {
		t.Errorf("Expected (0,-1) to be west of (0,0)")
	}
	if west.IsEastOf(ref) {
		t.Errorf("Expected (0,-1) to not be east of (0,0)")
	}
}

func TestCoordinate_Format(t *testing.T) {
	// New York City
	nyc, _ := NewCoordinate(40.7128, -74.0060)

	tests := []struct {
		name     string
		format   string
		contains string // Using contains instead of exact match due to floating point precision
	}{
		{"Decimal Degrees", "dd", "40.71280, -74.00600"},
		{"Default Format", "", "40.712800,-74.006000"},
		{"Invalid Format", "invalid", "40.712800,-74.006000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nyc.Format(tt.format)
			if !contains(result, tt.contains) {
				t.Errorf("Expected format %s to contain %s, got %s", tt.format, tt.contains, result)
			}
		})
	}

	// Test DMS format
	dmsFmt := nyc.Format("dms")
	if !contains(dmsFmt, "40°42'") && !contains(dmsFmt, "N") && !contains(dmsFmt, "74°0'") && !contains(dmsFmt, "W") {
		t.Errorf("Expected DMS format to contain degrees, minutes, seconds and direction, got %s", dmsFmt)
	}

	// Test DM format
	dmFmt := nyc.Format("dm")
	if !contains(dmFmt, "40°42") && !contains(dmFmt, "N") && !contains(dmFmt, "74°0") && !contains(dmFmt, "W") {
		t.Errorf("Expected DM format to contain degrees, decimal minutes and direction, got %s", dmFmt)
	}
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
