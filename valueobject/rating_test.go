// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"testing"
)

func TestNewRating(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		maxValue    float64
		expected    float64
		expectError bool
	}{
		{"Valid Rating", 4.5, 5, 4.5, false},
		{"Valid Rating - Zero", 0, 5, 0, false},
		{"Valid Rating - Max", 5, 5, 5, false},
		{"Valid Rating - Rounding", 4.55, 5, 4.6, false}, // Should round to 4.6
		{"Invalid Rating - Negative", -1, 5, 0, true},
		{"Invalid Rating - Exceeds Max", 6, 5, 0, true},
		{"Invalid Rating - Negative Max", 4, -5, 0, true},
		{"Invalid Rating - Zero Max", 4, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, err := NewRating(tt.value, tt.maxValue)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if rating.value != tt.expected {
					t.Errorf("Expected rating value %.1f, got %.1f", tt.expected, rating.value)
				}

				if rating.maxValue != tt.maxValue {
					t.Errorf("Expected max value %.1f, got %.1f", tt.maxValue, rating.maxValue)
				}
			}
		})
	}
}

func TestParseRating(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedVal float64
		expectedMax float64
		expectError bool
	}{
		{"Valid Rating - Simple", "4", 4, 5, false},
		{"Valid Rating - With Max", "4/5", 4, 5, false},
		{"Valid Rating - Different Max", "8/10", 8, 10, false},
		{"Valid Rating - Decimal", "4.5/5", 4.5, 5, false},
		{"Valid Rating - With Spaces", " 4.5 / 5 ", 4.5, 5, false},
		{"Invalid Rating - Empty", "", 0, 0, true},
		{"Invalid Rating - Wrong Format", "4:5", 0, 0, true},
		{"Invalid Rating - Letters", "four/five", 0, 0, true},
		{"Invalid Rating - Negative", "-4/5", 0, 0, true},
		{"Invalid Rating - Exceeds Max", "6/5", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, err := ParseRating(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if rating.value != tt.expectedVal {
					t.Errorf("Expected rating value %.1f, got %.1f", tt.expectedVal, rating.value)
				}

				if rating.maxValue != tt.expectedMax {
					t.Errorf("Expected max value %.1f, got %.1f", tt.expectedMax, rating.maxValue)
				}
			}
		})
	}
}

func TestRating_String(t *testing.T) {
	tests := []struct {
		name       string
		value      float64
		maxValue   float64
		expected   string
	}{
		{"Whole Number", 4, 5, "4.0/5.0"},
		{"Decimal", 4.5, 5, "4.5/5.0"},
		{"Zero", 0, 5, "0.0/5.0"},
		{"Different Scale", 8, 10, "8.0/10.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, _ := NewRating(tt.value, tt.maxValue)
			if rating.String() != tt.expected {
				t.Errorf("Expected string %s, got %s", tt.expected, rating.String())
			}
		})
	}
}

func TestRating_Equals(t *testing.T) {
	tests := []struct {
		name       string
		r1Value    float64
		r1Max      float64
		r2Value    float64
		r2Max      float64
		shouldEqual bool
	}{
		{"Same Rating", 4, 5, 4, 5, true},
		{"Different Rating", 4, 5, 3, 5, false},
		{"Equivalent Rating - Different Scale", 4, 5, 8, 10, true},
		{"Almost Equivalent - Within Epsilon", 4, 5, 7.99, 10, true},
		{"Almost Equivalent - Outside Epsilon", 4, 5, 7.9, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, _ := NewRating(tt.r1Value, tt.r1Max)
			r2, _ := NewRating(tt.r2Value, tt.r2Max)

			if r1.Equals(r2) != tt.shouldEqual {
				t.Errorf("Expected Equals to return %v for %.1f/%.1f and %.1f/%.1f", 
					tt.shouldEqual, tt.r1Value, tt.r1Max, tt.r2Value, tt.r2Max)
			}
		})
	}
}

func TestRating_IsEmpty(t *testing.T) {
	tests := []struct {
		name       string
		value      float64
		maxValue   float64
		expected   bool
	}{
		{"Empty Rating", 0, 0, true},
		{"Zero Rating", 0, 5, false},
		{"Non-Empty Rating", 4, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rating Rating
			if tt.value == 0 && tt.maxValue == 0 {
				rating = Rating{} // Empty rating
			} else {
				rating, _ = NewRating(tt.value, tt.maxValue)
			}
			
			if rating.IsEmpty() != tt.expected {
				t.Errorf("Expected IsEmpty to return %v for %.1f/%.1f", tt.expected, tt.value, tt.maxValue)
			}
		})
	}
}

func TestRating_Value(t *testing.T) {
	rating, _ := NewRating(4.5, 5)
	if rating.Value() != 4.5 {
		t.Errorf("Expected value 4.5, got %.1f", rating.Value())
	}
}

func TestRating_MaxValue(t *testing.T) {
	rating, _ := NewRating(4.5, 5)
	if rating.MaxValue() != 5 {
		t.Errorf("Expected max value 5, got %.1f", rating.MaxValue())
	}
}

func TestRating_Normalized(t *testing.T) {
	tests := []struct {
		name       string
		value      float64
		maxValue   float64
		expected   float64
	}{
		{"Full Rating", 5, 5, 1.0},
		{"Half Rating", 2.5, 5, 0.5},
		{"Zero Rating", 0, 5, 0.0},
		{"Different Scale", 8, 10, 0.8},
		{"Empty Rating", 0, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rating Rating
			if tt.value == 0 && tt.maxValue == 0 {
				rating = Rating{} // Empty rating
			} else {
				rating, _ = NewRating(tt.value, tt.maxValue)
			}
			
			normalized := rating.Normalized()
			if normalized != tt.expected {
				t.Errorf("Expected normalized value %.2f, got %.2f", tt.expected, normalized)
			}
		})
	}
}

func TestRating_ToScale(t *testing.T) {
	tests := []struct {
		name       string
		value      float64
		maxValue   float64
		newMax     float64
		expectedVal float64
		expectError bool
	}{
		{"5 to 10 Scale", 4, 5, 10, 8, false},
		{"10 to 5 Scale", 8, 10, 5, 4, false},
		{"Same Scale", 4, 5, 5, 4, false},
		{"Invalid New Scale", 4, 5, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, _ := NewRating(tt.value, tt.maxValue)
			newRating, err := rating.ToScale(tt.newMax)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				
				if newRating.value != tt.expectedVal {
					t.Errorf("Expected new value %.1f, got %.1f", tt.expectedVal, newRating.value)
				}
				
				if newRating.maxValue != tt.newMax {
					t.Errorf("Expected new max %.1f, got %.1f", tt.newMax, newRating.maxValue)
				}
			}
		})
	}
}

func TestRating_IsHigher(t *testing.T) {
	tests := []struct {
		name       string
		r1Value    float64
		r1Max      float64
		r2Value    float64
		r2Max      float64
		expected   bool
	}{
		{"Higher Rating", 4, 5, 3, 5, true},
		{"Lower Rating", 3, 5, 4, 5, false},
		{"Equal Rating", 4, 5, 4, 5, false},
		{"Higher Rating - Different Scale", 8, 10, 3, 5, true},
		{"Lower Rating - Different Scale", 3, 5, 8, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, _ := NewRating(tt.r1Value, tt.r1Max)
			r2, _ := NewRating(tt.r2Value, tt.r2Max)
			
			if r1.IsHigher(r2) != tt.expected {
				t.Errorf("Expected IsHigher to return %v for %.1f/%.1f and %.1f/%.1f", 
					tt.expected, tt.r1Value, tt.r1Max, tt.r2Value, tt.r2Max)
			}
		})
	}
}

func TestRating_IsLower(t *testing.T) {
	tests := []struct {
		name       string
		r1Value    float64
		r1Max      float64
		r2Value    float64
		r2Max      float64
		expected   bool
	}{
		{"Lower Rating", 3, 5, 4, 5, true},
		{"Higher Rating", 4, 5, 3, 5, false},
		{"Equal Rating", 4, 5, 4, 5, false},
		{"Lower Rating - Different Scale", 3, 5, 8, 10, true},
		{"Higher Rating - Different Scale", 8, 10, 3, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, _ := NewRating(tt.r1Value, tt.r1Max)
			r2, _ := NewRating(tt.r2Value, tt.r2Max)
			
			if r1.IsLower(r2) != tt.expected {
				t.Errorf("Expected IsLower to return %v for %.1f/%.1f and %.1f/%.1f", 
					tt.expected, tt.r1Value, tt.r1Max, tt.r2Value, tt.r2Max)
			}
		})
	}
}

func TestRating_Percentage(t *testing.T) {
	tests := []struct {
		name       string
		value      float64
		maxValue   float64
		expected   float64
	}{
		{"Full Rating", 5, 5, 100.0},
		{"Half Rating", 2.5, 5, 50.0},
		{"Zero Rating", 0, 5, 0.0},
		{"Different Scale", 8, 10, 80.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, _ := NewRating(tt.value, tt.maxValue)
			percentage := rating.Percentage()
			if percentage != tt.expected {
				t.Errorf("Expected percentage %.1f, got %.1f", tt.expected, percentage)
			}
		})
	}
}

func TestRating_Stars(t *testing.T) {
	tests := []struct {
		name       string
		value      float64
		maxValue   float64
		expected   string
	}{
		{"5 Stars", 5, 5, "★★★★★"},
		{"4.5 Stars", 4.5, 5, "★★★★½"},
		{"4 Stars", 4, 5, "★★★★☆"},
		{"3.5 Stars", 3.5, 5, "★★★½☆"},
		{"0 Stars", 0, 5, "☆☆☆☆☆"},
		{"Different Scale - 8/10", 8, 10, "★★★★☆"}, // Should convert to 4/5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, _ := NewRating(tt.value, tt.maxValue)
			stars := rating.Stars()
			if stars != tt.expected {
				t.Errorf("Expected stars %s, got %s", tt.expected, stars)
			}
		})
	}
}

func TestRating_Format(t *testing.T) {
	tests := []struct {
		name       string
		value      float64
		maxValue   float64
		format     string
		expected   string
	}{
		{"Decimal Format", 4.5, 5, "decimal", "4.5/5.0"},
		{"Percentage Format", 4.5, 5, "percentage", "90%"},
		{"Stars Format", 4.5, 5, "stars", "★★★★½"},
		{"Fraction Format", 4.5, 5, "fraction", "9/10"},
		{"Default Format", 4.5, 5, "unknown", "4.5/5.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, _ := NewRating(tt.value, tt.maxValue)
			formatted := rating.Format(tt.format)
			if formatted != tt.expected {
				t.Errorf("Expected format %s, got %s", tt.expected, formatted)
			}
		})
	}
}