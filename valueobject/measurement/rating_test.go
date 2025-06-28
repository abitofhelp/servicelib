// Copyright (c) 2025 A Bit of Help, Inc.

package measurement

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRating(t *testing.T) {
	tests := []struct {
		name string

		maxvalue float64

		value float64

		expectError bool
	}{
		{
			name: "Valid Value",

			maxvalue: 42.0,

			value: 42.0,

			expectError: false,
		},
		// Add more test cases specific to this value object
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := NewRating(tt.maxvalue, tt.value)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.maxvalue, value.maxValue)

				assert.Equal(t, tt.value, value.value)

			}
		})
	}
}

func TestRating_String(t *testing.T) {
	value := Rating{

		maxValue: 42.0,

		value: 42.0,
	}

	// The exact string representation depends on the fields
	assert.Contains(t, value.String(), "maxValue")
}

func TestRating_Equals(t *testing.T) {
	value1 := Rating{

		maxValue: 42.0,

		value: 42.0,
	}

	value2 := Rating{

		maxValue: 42.0,

		value: 42.0,
	}

	value3 := Rating{

		maxValue: 43.0,

		value: 43.0,
	}

	assert.True(t, value1.Equals(value2))
	assert.False(t, value1.Equals(value3))
}

func TestRating_IsEmpty(t *testing.T) {
	emptyValue := Rating{}
	assert.True(t, emptyValue.IsEmpty())

	nonEmptyValue := Rating{

		maxValue: 42.0,

		value: 42.0,
	}
	assert.False(t, nonEmptyValue.IsEmpty())
}

func TestRating_Validate(t *testing.T) {
	tests := []struct {
		name        string
		value       Rating
		expectError bool
	}{
		{
			name: "Valid Value",
			value: Rating{

				maxValue: 42.0,

				value: 42.0,
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

func TestRating_ToMap(t *testing.T) {
	value := Rating{

		maxValue: 42.0,

		value: 42.0,
	}

	expected := map[string]interface{}{

		"maxValue": 42.0,

		"value": 42.0,
	}

	assert.Equal(t, expected, value.ToMap())
}

func TestRating_MarshalJSON(t *testing.T) {
	value := Rating{

		maxValue: 42.0,

		value: 42.0,
	}

	data, err := value.MarshalJSON()
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	expected := map[string]interface{}{

		"maxValue": 42.0,

		"value": 42.0,
	}

	assert.Equal(t, expected, result)
}
