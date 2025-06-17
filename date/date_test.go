// Copyright (c) 2025 A Bit of Help, Inc.

package date

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name        string
		dateStr     string
		expected    time.Time
		expectError bool
	}{
		{
			name:        "valid date",
			dateStr:     "2023-01-02T15:04:05Z",
			expected:    time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "invalid date format",
			dateStr:     "2023-01-02",
			expectError: true,
		},
		{
			name:        "empty string",
			dateStr:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDate(tt.dateStr)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseOptionalDate(t *testing.T) {
	validDateStr := "2023-01-02T15:04:05Z"
	invalidDateStr := "2023-01-02"
	emptyStr := ""
	
	expectedTime := time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC)

	tests := []struct {
		name        string
		dateStr     *string
		expected    *time.Time
		expectError bool
	}{
		{
			name:        "nil date string",
			dateStr:     nil,
			expected:    nil,
			expectError: false,
		},
		{
			name:        "valid date",
			dateStr:     &validDateStr,
			expected:    &expectedTime,
			expectError: false,
		},
		{
			name:        "invalid date format",
			dateStr:     &invalidDateStr,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "empty string",
			dateStr:     &emptyStr,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseOptionalDate(tt.dateStr)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expected == nil {
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
					assert.Equal(t, *tt.expected, *result)
				}
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "standard date",
			date:     time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC),
			expected: "2023-01-02T15:04:05Z",
		},
		{
			name:     "zero date",
			date:     time.Time{},
			expected: "0001-01-01T00:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDate(tt.date)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatOptionalDate(t *testing.T) {
	date := time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC)
	zeroDate := time.Time{}
	expectedStr := "2023-01-02T15:04:05Z"
	expectedZeroStr := "0001-01-01T00:00:00Z"

	tests := []struct {
		name     string
		date     *time.Time
		expected *string
	}{
		{
			name:     "nil date",
			date:     nil,
			expected: nil,
		},
		{
			name:     "standard date",
			date:     &date,
			expected: &expectedStr,
		},
		{
			name:     "zero date",
			date:     &zeroDate,
			expected: &expectedZeroStr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatOptionalDate(tt.date)

			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}