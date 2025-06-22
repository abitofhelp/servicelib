// Copyright (c) 2025 A Bit of Help, Inc.

package temporal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDuration(t *testing.T) {
	tests := []struct {
		name string

		duration time.Duration

		expectError bool
	}{
		{
			name: "Valid Value",

			duration: 42,

			expectError: false,
		},
		// Add more test cases specific to this value object
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := NewDuration(tt.duration)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.duration, value.duration)

			}
		})
	}
}

func TestDuration_String(t *testing.T) {
	value := Duration{

		duration: 42,
	}

	// The exact string representation depends on the fields
	assert.Contains(t, value.String(), "duration")
}

func TestDuration_Equals(t *testing.T) {
	value1 := Duration{

		duration: 42,
	}

	value2 := Duration{

		duration: 42,
	}

	value3 := Duration{

		duration: 43,
	}

	assert.True(t, value1.Equals(value2))
	assert.False(t, value1.Equals(value3))
}

func TestDuration_IsEmpty(t *testing.T) {
	emptyValue := Duration{}
	assert.True(t, emptyValue.IsEmpty())

	nonEmptyValue := Duration{

		duration: 42,
	}
	assert.False(t, nonEmptyValue.IsEmpty())
}

func TestDuration_Validate(t *testing.T) {
	tests := []struct {
		name        string
		value       Duration
		expectError bool
	}{
		{
			name: "Valid Value",
			value: Duration{

				duration: 42,
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

func TestDuration_ToMap(t *testing.T) {
	value := Duration{

		duration: 42,
	}

	expected := map[string]interface{}{

		"duration": 42,
	}

	assert.Equal(t, expected, value.ToMap())
}

func TestDuration_MarshalJSON(t *testing.T) {
	value := Duration{

		duration: 42,
	}

	data, err := value.MarshalJSON()
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	expected := map[string]interface{}{

		"duration": float64(42),
	}

	assert.Equal(t, expected, result)
}
