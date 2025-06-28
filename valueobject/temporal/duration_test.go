// Copyright (c) 2025 A Bit of Help, Inc.

package temporal

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDuration(t *testing.T) {
	tests := []struct {
		name        string
		duration    time.Duration
		expectError bool
	}{
		{
			name:        "Valid Positive Duration",
			duration:    time.Hour + 30*time.Minute,
			expectError: false,
		},
		{
			name:        "Zero Duration",
			duration:    0,
			expectError: false,
		},
		{
			name:        "Negative Duration",
			duration:    -time.Hour,
			expectError: true,
		},
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

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    time.Duration
		expectError bool
	}{
		{
			name:        "Valid Duration",
			input:       "1h30m45s",
			expected:    time.Hour + 30*time.Minute + 45*time.Second,
			expectError: false,
		},
		{
			name:        "Zero Duration",
			input:       "0s",
			expected:    0,
			expectError: false,
		},
		{
			name:        "Empty String",
			input:       "",
			expectError: true,
		},
		{
			name:        "Invalid Format",
			input:       "invalid",
			expectError: true,
		},
		{
			name:        "Negative Duration",
			input:       "-1h",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := ParseDuration(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, value.Value())
			}
		})
	}
}

func TestDuration_String(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "Hour and Minutes",
			duration: time.Hour + 30*time.Minute,
			expected: "1h30m0s",
		},
		{
			name:     "Zero Duration",
			duration: 0,
			expected: "0s",
		},
		{
			name:     "Only Seconds",
			duration: 45 * time.Second,
			expected: "45s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := NewDuration(tt.duration)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, value.String())
		})
	}
}

func TestDuration_Equals(t *testing.T) {
	d1, _ := NewDuration(time.Hour)
	d2, _ := NewDuration(time.Hour)
	d3, _ := NewDuration(2 * time.Hour)

	assert.True(t, d1.Equals(d2))
	assert.False(t, d1.Equals(d3))
}

func TestDuration_IsEmpty(t *testing.T) {
	emptyDuration, _ := NewDuration(0)
	nonEmptyDuration, _ := NewDuration(time.Hour)

	assert.True(t, emptyDuration.IsEmpty())
	assert.False(t, nonEmptyDuration.IsEmpty())
}

func TestDuration_Value(t *testing.T) {
	expected := 2*time.Hour + 30*time.Minute
	d, _ := NewDuration(expected)
	assert.Equal(t, expected, d.Value())
}

func TestDuration_TimeUnits(t *testing.T) {
	d, _ := NewDuration(2*time.Hour + 30*time.Minute + 45*time.Second + 500*time.Millisecond)

	assert.Equal(t, 2.5125, d.Hours())
	assert.Equal(t, 150.75, d.Minutes())
	assert.Equal(t, 9045.5, d.Seconds())
	assert.Equal(t, int64(9045500), d.Milliseconds())
}

func TestDuration_Add(t *testing.T) {
	d1, _ := NewDuration(time.Hour)
	d2, _ := NewDuration(30 * time.Minute)

	result, err := d1.Add(d2)
	require.NoError(t, err)
	assert.Equal(t, time.Hour+30*time.Minute, result.Value())
}

func TestDuration_Subtract(t *testing.T) {
	tests := []struct {
		name     string
		d1       time.Duration
		d2       time.Duration
		expected time.Duration
	}{
		{
			name:     "Normal Subtraction",
			d1:       2 * time.Hour,
			d2:       30 * time.Minute,
			expected: 90 * time.Minute,
		},
		{
			name:     "Zero Result",
			d1:       time.Hour,
			d2:       time.Hour,
			expected: 0,
		},
		{
			name:     "Would Be Negative",
			d1:       time.Hour,
			d2:       2 * time.Hour,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d1, _ := NewDuration(tt.d1)
			d2, _ := NewDuration(tt.d2)

			result, err := d1.Subtract(d2)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.Value())
		})
	}
}

func TestDuration_Format(t *testing.T) {
	d, _ := NewDuration(2*time.Hour + 30*time.Minute + 45*time.Second)

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "Short Format",
			format:   "short",
			expected: "2h 30m 45s",
		},
		{
			name:     "Long Format",
			format:   "long",
			expected: "2 hours 30 minutes 45 seconds",
		},
		{
			name:     "Compact Format",
			format:   "compact",
			expected: "2:30:45",
		},
		{
			name:     "Default Format",
			format:   "default",
			expected: "2h30m45s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, d.Format(tt.format))
		})
	}

	// Test singular forms
	d1, _ := NewDuration(time.Hour + time.Minute + time.Second)
	assert.Equal(t, "1 hour 1 minute 1 second", d1.Format("long"))

	// Test zero duration
	d2, _ := NewDuration(0)
	assert.Equal(t, "0 seconds", d2.Format("long"))
}

func TestDuration_ToMap(t *testing.T) {
	d, _ := NewDuration(time.Hour)
	expected := map[string]interface{}{
		"duration": time.Hour,
	}
	assert.Equal(t, expected, d.ToMap())
}

func TestDuration_MarshalJSON(t *testing.T) {
	d, _ := NewDuration(time.Hour)
	data, err := d.MarshalJSON()
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"duration": float64(time.Hour),
	}
	assert.Equal(t, expected, result)
}
