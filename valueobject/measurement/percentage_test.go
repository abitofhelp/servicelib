// Copyright (c) 2025 A Bit of Help, Inc.

package measurement

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPercentage(t *testing.T) {
	tests := []struct {
		name string

		value float64

		expectError bool
	}{
		{
			name: "Valid Value",

			value: 42.0,

			expectError: false,
		},
		// Add more test cases specific to this value object
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := NewPercentage(tt.value)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.value, value.value)

			}
		})
	}
}

func TestPercentage_String(t *testing.T) {
	value := Percentage{

		value: 42.0,
	}

	// The exact string representation depends on the fields
	assert.Contains(t, value.String(), "value")
}

func TestPercentage_Equals(t *testing.T) {
	value1 := Percentage{

		value: 42.0,
	}

	value2 := Percentage{

		value: 42.0,
	}

	value3 := Percentage{

		value: 43.0,
	}

	assert.True(t, value1.Equals(value2))
	assert.False(t, value1.Equals(value3))
}

func TestPercentage_IsEmpty(t *testing.T) {
	emptyValue := Percentage{}
	assert.True(t, emptyValue.IsEmpty())

	nonEmptyValue := Percentage{

		value: 42.0,
	}
	assert.False(t, nonEmptyValue.IsEmpty())
}

func TestPercentage_Validate(t *testing.T) {
	tests := []struct {
		name        string
		value       Percentage
		expectError bool
	}{
		{
			name: "Valid Value",
			value: Percentage{

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

func TestPercentage_ToMap(t *testing.T) {
	value := Percentage{

		value: 42.0,
	}

	expected := map[string]interface{}{

		"value": 42.0,
	}

	assert.Equal(t, expected, value.ToMap())
}

func TestPercentage_MarshalJSON(t *testing.T) {
	value := Percentage{

		value: 42.0,
	}

	data, err := value.MarshalJSON()
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	expected := map[string]interface{}{

		"value": 42.0,
	}

	assert.Equal(t, expected, result)
}
