// Copyright (c) 2025 A Bit of Help, Inc.

package identification

import (
	"testing"
)

func TestNewGender(t *testing.T) {
	tests := []struct {
		name        string
		gender      string
		expected    Gender
		expectError bool
	}{
		{"Male", "male", GenderMale, false},
		{"Female", "female", GenderFemale, false},
		{"Other", "other", GenderOther, false},
		{"Male with Spaces", " male ", GenderMale, false},
		{"Female with Uppercase", "FEMALE", GenderFemale, false},
		{"Other with Mixed Case", "OtHeR", GenderOther, false},
		{"Empty String", "", "", true},
		{"Invalid Gender", "unknown", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gender, err := NewGender(tt.gender)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if gender != tt.expected {
					t.Errorf("Expected gender %s, got %s", tt.expected, gender)
				}
			}
		})
	}
}

func TestGender_String(t *testing.T) {
	tests := []struct {
		name     string
		gender   Gender
		expected string
	}{
		{"Male", GenderMale, "male"},
		{"Female", GenderFemale, "female"},
		{"Other", GenderOther, "other"},
		{"Empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.gender.String() != tt.expected {
				t.Errorf("Expected string %s, got %s", tt.expected, tt.gender.String())
			}
		})
	}
}

func TestGender_Equals(t *testing.T) {
	tests := []struct {
		name     string
		gender1  Gender
		gender2  Gender
		expected bool
	}{
		{"Same Gender - Male", GenderMale, GenderMale, true},
		{"Same Gender - Female", GenderFemale, GenderFemale, true},
		{"Same Gender - Other", GenderOther, GenderOther, true},
		{"Different Gender - Male vs Female", GenderMale, GenderFemale, false},
		{"Different Gender - Male vs Other", GenderMale, GenderOther, false},
		{"Different Gender - Female vs Other", GenderFemale, GenderOther, false},
		{"Empty vs Male", "", GenderMale, false},
		{"Empty vs Empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.gender1.Equals(tt.gender2) != tt.expected {
				t.Errorf("Expected Equals to return %v for %s and %s", tt.expected, tt.gender1, tt.gender2)
			}
		})
	}
}

func TestGender_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		gender   Gender
		expected bool
	}{
		{"Male", GenderMale, false},
		{"Female", GenderFemale, false},
		{"Other", GenderOther, false},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.gender.IsEmpty() != tt.expected {
				t.Errorf("Expected IsEmpty to return %v for %s", tt.expected, tt.gender)
			}
		})
	}
}
