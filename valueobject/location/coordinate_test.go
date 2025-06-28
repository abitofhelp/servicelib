// Copyright (c) 2025 A Bit of Help, Inc.

package location

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCoordinate(t *testing.T) {
	tests := []struct {
		name string

		latitude float64

		longitude float64

		expectError bool
	}{
		{
			name: "Valid Value",

			latitude: 42.0,

			longitude: 42.0,

			expectError: false,
		},
		// Add more test cases specific to this value object
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := NewCoordinate(tt.latitude, tt.longitude)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.latitude, value.latitude)

				assert.Equal(t, tt.longitude, value.longitude)

			}
		})
	}
}

func TestCoordinate_String(t *testing.T) {
	value := Coordinate{

		latitude: 42.0,

		longitude: 42.0,
	}

	// The exact string representation depends on the fields
	assert.Contains(t, value.String(), "latitude")
}

func TestCoordinate_Equals(t *testing.T) {
	value1 := Coordinate{

		latitude: 42.0,

		longitude: 42.0,
	}

	value2 := Coordinate{

		latitude: 42.0,

		longitude: 42.0,
	}

	value3 := Coordinate{

		latitude: 43.0,

		longitude: 43.0,
	}

	assert.True(t, value1.Equals(value2))
	assert.False(t, value1.Equals(value3))
}

func TestCoordinate_IsEmpty(t *testing.T) {
	emptyValue := Coordinate{}
	assert.True(t, emptyValue.IsEmpty())

	nonEmptyValue := Coordinate{

		latitude: 42.0,

		longitude: 42.0,
	}
	assert.False(t, nonEmptyValue.IsEmpty())
}

func TestCoordinate_Validate(t *testing.T) {
	tests := []struct {
		name        string
		value       Coordinate
		expectError bool
	}{
		{
			name: "Valid Value",
			value: Coordinate{

				latitude: 42.0,

				longitude: 42.0,
			},
			expectError: false,
		},
		// Add more test cases specific to this value object
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.value.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCoordinate_ToMap(t *testing.T) {
	value := Coordinate{

		latitude: 42.0,

		longitude: 42.0,
	}

	expected := map[string]interface{}{

		"latitude": 42.0,

		"longitude": 42.0,
	}

	assert.Equal(t, expected, value.ToMap())
}

func TestCoordinate_MarshalJSON(t *testing.T) {
	value := Coordinate{

		latitude: 42.0,

		longitude: 42.0,
	}

	data, err := value.MarshalJSON()
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	expected := map[string]interface{}{

		"latitude": 42.0,

		"longitude": 42.0,
	}

	assert.Equal(t, expected, result)
}
