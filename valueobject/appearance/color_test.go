// Copyright (c) 2025 A Bit of Help, Inc.

package appearance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewColor(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{"Valid 6-digit Hex", "#FF5733", "#FF5733", false},
		{"Valid 6-digit Hex Lowercase", "#ff5733", "#FF5733", false},
		{"Valid 3-digit Hex", "#F57", "#FF5577", false},
		{"Valid 3-digit Hex Lowercase", "#f57", "#FF5577", false},
		{"Valid Hex without #", "FF5733", "#FF5733", false},
		{"Valid 3-digit Hex without #", "F57", "#FF5577", false},
		{"Empty Color", "", "", false}, // Empty is allowed
		{"Invalid Format", "ZZZZZZ", "", true},
		{"Invalid Length", "#FFFF", "", true},
		{"Invalid Characters", "#FFZ733", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := NewColor(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, color.String())
			}
		})
	}
}

func TestColor_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Regular Color", "#FF5733", "#FF5733"},
		{"Empty Color", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, _ := NewColor(tt.input)
			assert.Equal(t, tt.expected, color.String())
		})
	}
}

func TestColor_Equals(t *testing.T) {
	tests := []struct {
		name     string
		color1   string
		color2   string
		expected bool
	}{
		{"Same Color", "#FF5733", "#FF5733", true},
		{"Different Case", "#ff5733", "#FF5733", true},
		{"Different Color", "#FF5733", "#00FF00", false},
		{"Empty Colors", "", "", true},
		{"One Empty Color", "#FF5733", "", false},
		{"3-digit vs 6-digit Equivalent", "#F57", "#FF5577", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color1, _ := NewColor(tt.color1)
			color2, _ := NewColor(tt.color2)

			assert.Equal(t, tt.expected, color1.Equals(color2))
		})
	}
}

func TestColor_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		expected bool
	}{
		{"Empty Color", "", true},
		{"Non-Empty Color", "#FF5733", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, _ := NewColor(tt.color)
			assert.Equal(t, tt.expected, color.IsEmpty())
		})
	}
}

func TestColor_Validate(t *testing.T) {
	tests := []struct {
		name        string
		color       Color
		expectError bool
	}{
		{"Valid Color", Color("#FF5733"), false},
		{"Empty Color", Color(""), false}, // Empty is allowed
		{"Invalid Format", Color("#ZZZZZZ"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.color.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestColor_RGB(t *testing.T) {
	tests := []struct {
		name        string
		color       string
		expectedR   int
		expectedG   int
		expectedB   int
		expectError bool
	}{
		{"Valid Color", "#FF5733", 255, 87, 51, false},
		{"Black", "#000000", 0, 0, 0, false},
		{"White", "#FFFFFF", 255, 255, 255, false},
		{"Empty Color", "", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, _ := NewColor(tt.color)
			r, g, b, err := color.RGB()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedR, r)
				assert.Equal(t, tt.expectedG, g)
				assert.Equal(t, tt.expectedB, b)
			}
		})
	}
}

func TestColor_WithAlpha(t *testing.T) {
	tests := []struct {
		name        string
		color       string
		alpha       int
		expected    string
		expectError bool
	}{
		{"Valid Alpha", "#FF5733", 128, "#FF573380", false},
		{"Zero Alpha", "#FF5733", 0, "#FF573300", false},
		{"Full Alpha", "#FF5733", 255, "#FF5733FF", false},
		{"Negative Alpha", "#FF5733", -1, "", true},
		{"Too High Alpha", "#FF5733", 256, "", true},
		{"Empty Color", "", 128, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, _ := NewColor(tt.color)
			result, err := color.WithAlpha(tt.alpha)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestColor_IsDark(t *testing.T) {
	tests := []struct {
		name        string
		color       string
		expected    bool
		expectError bool
	}{
		{"Dark Color", "#000000", true, false},
		{"Light Color", "#FFFFFF", false, false},
		{"Medium Dark Color", "#555555", true, false},
		{"Medium Light Color", "#AAAAAA", false, false},
		{"Empty Color", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, _ := NewColor(tt.color)
			result, err := color.IsDark()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestColor_Invert(t *testing.T) {
	tests := []struct {
		name        string
		color       string
		expected    string
		expectError bool
	}{
		{"Black to White", "#000000", "#FFFFFF", false},
		{"White to Black", "#FFFFFF", "#000000", false},
		{"Invert Red", "#FF0000", "#00FFFF", false},
		{"Invert Green", "#00FF00", "#FF00FF", false},
		{"Invert Blue", "#0000FF", "#FFFF00", false},
		{"Empty Color", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, _ := NewColor(tt.color)
			result, err := color.Invert()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result.String())
			}
		})
	}
}
