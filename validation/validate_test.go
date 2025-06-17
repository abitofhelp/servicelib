// Copyright (c) 2025 A Bit of Help, Inc.

package validation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewValidationResult(t *testing.T) {
	result := NewValidationResult()
	assert.NotNil(t, result)
	assert.True(t, result.IsValid())
	assert.Nil(t, result.Error())
}

func TestValidationResult_AddError(t *testing.T) {
	result := NewValidationResult()
	
	// Add an error
	result.AddError("test error", "field1")
	
	// Check that the result is now invalid
	assert.False(t, result.IsValid())
	assert.NotNil(t, result.Error())
}

func TestValidationResult_IsValid(t *testing.T) {
	result := NewValidationResult()
	
	// Initially valid
	assert.True(t, result.IsValid())
	
	// Add an error
	result.AddError("test error", "field1")
	
	// Now invalid
	assert.False(t, result.IsValid())
}

func TestValidationResult_Error(t *testing.T) {
	result := NewValidationResult()
	
	// Initially no error
	assert.Nil(t, result.Error())
	
	// Add an error
	result.AddError("test error", "field1")
	
	// Now has error
	assert.NotNil(t, result.Error())
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		field    string
		expected bool
	}{
		{
			name:     "empty string",
			value:    "",
			field:    "field1",
			expected: false,
		},
		{
			name:     "whitespace only",
			value:    "   ",
			field:    "field1",
			expected: false,
		},
		{
			name:     "non-empty string",
			value:    "value",
			field:    "field1",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewValidationResult()
			Required(tt.value, tt.field, result)
			assert.Equal(t, tt.expected, result.IsValid())
		})
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		min      int
		field    string
		expected bool
	}{
		{
			name:     "empty string",
			value:    "",
			min:      1,
			field:    "field1",
			expected: false,
		},
		{
			name:     "string shorter than min",
			value:    "ab",
			min:      3,
			field:    "field1",
			expected: false,
		},
		{
			name:     "string equal to min",
			value:    "abc",
			min:      3,
			field:    "field1",
			expected: true,
		},
		{
			name:     "string longer than min",
			value:    "abcd",
			min:      3,
			field:    "field1",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewValidationResult()
			MinLength(tt.value, tt.min, tt.field, result)
			assert.Equal(t, tt.expected, result.IsValid())
		})
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		max      int
		field    string
		expected bool
	}{
		{
			name:     "empty string",
			value:    "",
			max:      1,
			field:    "field1",
			expected: true,
		},
		{
			name:     "string shorter than max",
			value:    "ab",
			max:      3,
			field:    "field1",
			expected: true,
		},
		{
			name:     "string equal to max",
			value:    "abc",
			max:      3,
			field:    "field1",
			expected: true,
		},
		{
			name:     "string longer than max",
			value:    "abcd",
			max:      3,
			field:    "field1",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewValidationResult()
			MaxLength(tt.value, tt.max, tt.field, result)
			assert.Equal(t, tt.expected, result.IsValid())
		})
	}
}

func TestPattern(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		pattern  string
		field    string
		expected bool
	}{
		{
			name:     "matching pattern",
			value:    "abc123",
			pattern:  "^[a-z]+[0-9]+$",
			field:    "field1",
			expected: true,
		},
		{
			name:     "non-matching pattern",
			value:    "123abc",
			pattern:  "^[a-z]+[0-9]+$",
			field:    "field1",
			expected: false,
		},
		{
			name:     "empty string",
			value:    "",
			pattern:  "^[a-z]+[0-9]+$",
			field:    "field1",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewValidationResult()
			Pattern(tt.value, tt.pattern, tt.field, result)
			assert.Equal(t, tt.expected, result.IsValid())
		})
	}
}

func TestPastDate(t *testing.T) {
	now := time.Now()
	past := now.Add(-24 * time.Hour)
	future := now.Add(24 * time.Hour)

	tests := []struct {
		name     string
		value    time.Time
		field    string
		expected bool
	}{
		{
			name:     "past date",
			value:    past,
			field:    "field1",
			expected: true,
		},
		{
			name:     "future date",
			value:    future,
			field:    "field1",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewValidationResult()
			PastDate(tt.value, tt.field, result)
			assert.Equal(t, tt.expected, result.IsValid())
		})
	}
}

func TestValidDateRange(t *testing.T) {
	now := time.Now()
	past := now.Add(-24 * time.Hour)
	future := now.Add(24 * time.Hour)
	zeroTime := time.Time{}

	tests := []struct {
		name       string
		start      time.Time
		end        time.Time
		startField string
		endField   string
		expected   bool
	}{
		{
			name:       "start before end",
			start:      past,
			end:        now,
			startField: "startField",
			endField:   "endField",
			expected:   true,
		},
		{
			name:       "start after end",
			start:      future,
			end:        now,
			startField: "startField",
			endField:   "endField",
			expected:   false,
		},
		{
			name:       "end is zero time",
			start:      now,
			end:        zeroTime,
			startField: "startField",
			endField:   "endField",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewValidationResult()
			ValidDateRange(tt.start, tt.end, tt.startField, tt.endField, result)
			assert.Equal(t, tt.expected, result.IsValid())
		})
	}
}

func TestAllTrue(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		predicate func(int) bool
		expected  bool
	}{
		{
			name:      "all items satisfy predicate",
			items:     []int{2, 4, 6, 8},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  true,
		},
		{
			name:      "some items don't satisfy predicate",
			items:     []int{2, 3, 6, 8},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  false,
		},
		{
			name:      "empty slice",
			items:     []int{},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AllTrue(tt.items, tt.predicate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateAll(t *testing.T) {
	tests := []struct {
		name      string
		items     []string
		validator func(string, int, *ValidationResult)
		expected  bool
	}{
		{
			name:  "all items valid",
			items: []string{"a", "b", "c"},
			validator: func(s string, i int, result *ValidationResult) {
				// Do nothing, all valid
			},
			expected: true,
		},
		{
			name:  "some items invalid",
			items: []string{"a", "", "c"},
			validator: func(s string, i int, result *ValidationResult) {
				if s == "" {
					result.AddError("is empty", "item"+string(rune('0'+i)))
				}
			},
			expected: false,
		},
		{
			name:  "empty slice",
			items: []string{},
			validator: func(s string, i int, result *ValidationResult) {
				// Do nothing
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewValidationResult()
			ValidateAll(tt.items, tt.validator, result)
			assert.Equal(t, tt.expected, result.IsValid())
		})
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		field    string
		expected bool
	}{
		{
			name:     "valid ID",
			id:       "123",
			field:    "id",
			expected: true,
		},
		{
			name:     "empty ID",
			id:       "",
			field:    "id",
			expected: false,
		},
		{
			name:     "whitespace ID",
			id:       "   ",
			field:    "id",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewValidationResult()
			ValidateID(tt.id, tt.field, result)
			assert.Equal(t, tt.expected, result.IsValid())
		})
	}
}