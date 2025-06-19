// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"testing"
)

func TestNewPercentage(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		expected    float64
		expectError bool
	}{
		{"Valid Percentage", 75.5, 75.5, false},
		{"Zero Percentage", 0, 0, false},
		{"Max Percentage", 100, 100, false},
		{"Negative Percentage", -10, 0, true},
		{"Exceeds Max", 101, 0, true},
		{"Rounding Test", 75.555, 75.56, false}, // Should round to 75.56
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percentage, err := NewPercentage(tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if float64(percentage) != tt.expected {
					t.Errorf("Expected percentage %.2f, got %.2f", tt.expected, float64(percentage))
				}
			}
		})
	}
}

func TestParsePercentage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    float64
		expectError bool
	}{
		{"Valid Percentage", "75.5", 75.5, false},
		{"Valid Percentage with %", "75.5%", 75.5, false},
		{"Valid Percentage with Spaces", " 75.5% ", 75.5, false},
		{"Zero Percentage", "0", 0, false},
		{"Zero Percentage with %", "0%", 0, false},
		{"Max Percentage", "100", 100, false},
		{"Empty String", "", 0, true},
		{"Invalid Format", "abc", 0, true},
		{"Negative Percentage", "-10", 0, true},
		{"Exceeds Max", "101", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percentage, err := ParsePercentage(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if float64(percentage) != tt.expected {
					t.Errorf("Expected percentage %.2f, got %.2f", tt.expected, float64(percentage))
				}
			}
		})
	}
}

func TestPercentage_String(t *testing.T) {
	tests := []struct {
		name       string
		percentage float64
		expected   string
	}{
		{"Whole Number", 75, "75.00%"},
		{"Decimal", 75.5, "75.50%"},
		{"Zero", 0, "0.00%"},
		{"Max", 100, "100.00%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percentage, _ := NewPercentage(tt.percentage)
			if percentage.String() != tt.expected {
				t.Errorf("Expected string %s, got %s", tt.expected, percentage.String())
			}
		})
	}
}

func TestPercentage_Equals(t *testing.T) {
	tests := []struct {
		name       string
		p1         float64
		p2         float64
		shouldEqual bool
	}{
		{"Same Value", 75.5, 75.5, true},
		{"Different Values", 75.5, 80, false},
		{"Close Values Within Epsilon", 75.5, 75.5001, true},
		{"Close Values Outside Epsilon", 75.5, 75.6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1, _ := NewPercentage(tt.p1)
			p2, _ := NewPercentage(tt.p2)

			if p1.Equals(p2) != tt.shouldEqual {
				t.Errorf("Expected Equals to return %v for %.4f and %.4f", tt.shouldEqual, tt.p1, tt.p2)
			}
		})
	}
}

func TestPercentage_IsEmpty(t *testing.T) {
	tests := []struct {
		name       string
		percentage float64
		expected   bool
	}{
		{"Zero Percentage", 0, true},
		{"Non-Zero Percentage", 75.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percentage, _ := NewPercentage(tt.percentage)
			if percentage.IsEmpty() != tt.expected {
				t.Errorf("Expected IsEmpty to return %v for %.2f", tt.expected, tt.percentage)
			}
		})
	}
}

func TestPercentage_Value(t *testing.T) {
	tests := []struct {
		name       string
		percentage float64
	}{
		{"Whole Number", 75},
		{"Decimal", 75.5},
		{"Zero", 0},
		{"Max", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percentage, _ := NewPercentage(tt.percentage)
			if percentage.Value() != tt.percentage {
				t.Errorf("Expected value %.2f, got %.2f", tt.percentage, percentage.Value())
			}
		})
	}
}

func TestPercentage_AsDecimal(t *testing.T) {
	tests := []struct {
		name       string
		percentage float64
		expected   float64
	}{
		{"Whole Number", 75, 0.75},
		{"Decimal", 75.5, 0.755},
		{"Zero", 0, 0},
		{"Max", 100, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percentage, _ := NewPercentage(tt.percentage)
			if percentage.AsDecimal() != tt.expected {
				t.Errorf("Expected decimal %.4f, got %.4f", tt.expected, percentage.AsDecimal())
			}
		})
	}
}

func TestPercentage_Of(t *testing.T) {
	tests := []struct {
		name       string
		percentage float64
		value      float64
		expected   float64
	}{
		{"75% of 200", 75, 200, 150},
		{"50% of 100", 50, 100, 50},
		{"25.5% of 200", 25.5, 200, 51},
		{"0% of 200", 0, 200, 0},
		{"100% of 200", 100, 200, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percentage, _ := NewPercentage(tt.percentage)
			result := percentage.Of(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, result)
			}
		})
	}
}

func TestPercentage_Inverse(t *testing.T) {
	tests := []struct {
		name       string
		percentage float64
		expected   float64
	}{
		{"Inverse of 75%", 75, 25},
		{"Inverse of 25%", 25, 75},
		{"Inverse of 0%", 0, 100},
		{"Inverse of 100%", 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percentage, _ := NewPercentage(tt.percentage)
			inverse, err := percentage.Inverse()
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if float64(inverse) != tt.expected {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, float64(inverse))
			}
		})
	}
}

func TestPercentage_Add(t *testing.T) {
	tests := []struct {
		name       string
		p1         float64
		p2         float64
		expected   float64
		expectError bool
	}{
		{"Normal Addition", 25, 50, 75, false},
		{"Addition to Zero", 0, 50, 50, false},
		{"Addition of Zero", 50, 0, 50, false},
		{"Result at Max", 50, 50, 100, false},
		{"Result Exceeds Max", 75, 50, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1, _ := NewPercentage(tt.p1)
			p2, _ := NewPercentage(tt.p2)
			
			result, err := p1.Add(p2)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				
				if float64(result) != tt.expected {
					t.Errorf("Expected %.2f, got %.2f", tt.expected, float64(result))
				}
			}
		})
	}
}

func TestPercentage_Subtract(t *testing.T) {
	tests := []struct {
		name       string
		p1         float64
		p2         float64
		expected   float64
	}{
		{"Normal Subtraction", 75, 25, 50},
		{"Subtraction from Zero", 0, 25, 0},
		{"Subtraction of Zero", 75, 0, 75},
		{"Result at Zero", 25, 25, 0},
		{"Result Would Be Negative", 25, 50, 0}, // Should floor at 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1, _ := NewPercentage(tt.p1)
			p2, _ := NewPercentage(tt.p2)
			
			result, err := p1.Subtract(p2)
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if float64(result) != tt.expected {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, float64(result))
			}
		})
	}
}