// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"math"
	"testing"
)

func TestNewTemperature(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		unit        TemperatureUnit
		expected    float64
		expectError bool
	}{
		{"Valid Celsius", 25.5, Celsius, 25.5, false},
		{"Valid Fahrenheit", 77.3, Fahrenheit, 77.3, false},
		{"Valid Kelvin", 298.15, Kelvin, 298.15, false},
		{"Below Absolute Zero Celsius", -274, Celsius, 0, true},
		{"Below Absolute Zero Fahrenheit", -460, Fahrenheit, 0, true},
		{"Below Absolute Zero Kelvin", -1, Kelvin, 0, true},
		{"Invalid Unit", 25, TemperatureUnit("X"), 0, true},
		{"Rounding Test", 25.567, Celsius, 25.57, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp, err := NewTemperature(tt.value, tt.unit)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if temp.value != tt.expected {
					t.Errorf("Expected value %f, got %f", tt.expected, temp.value)
				}
				if temp.unit != tt.unit {
					t.Errorf("Expected unit %s, got %s", tt.unit, temp.unit)
				}
			}
		})
	}
}

func TestParseTemperature(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedVal  float64
		expectedUnit TemperatureUnit
		expectError  bool
	}{
		{"Valid Celsius", "25°C", 25, Celsius, false},
		{"Valid Fahrenheit", "77°F", 77, Fahrenheit, false},
		{"Valid Kelvin", "298K", 298, Kelvin, false},
		{"Valid Negative Celsius", "-10°C", -10, Celsius, false},
		{"Valid Decimal", "25.5°C", 25.5, Celsius, false},
		{"No Degree Symbol", "25C", 25, Celsius, false},
		{"With Space", "25 °C", 25, Celsius, false},
		{"Empty String", "", 0, "", true},
		{"Invalid Format", "25 Celsius", 0, "", true},
		{"Invalid Value", "abc°C", 0, "", true},
		{"Invalid Unit", "25°X", 0, "", true},
		{"Below Absolute Zero", "-274°C", 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp, err := ParseTemperature(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if temp.value != tt.expectedVal {
					t.Errorf("Expected value %f, got %f", tt.expectedVal, temp.value)
				}
				if temp.unit != tt.expectedUnit {
					t.Errorf("Expected unit %s, got %s", tt.expectedUnit, temp.unit)
				}
			}
		})
	}
}

func TestTemperature_String(t *testing.T) {
	tests := []struct {
		name     string
		temp     Temperature
		expected string
	}{
		{"Celsius", Temperature{25, Celsius}, "25.00°C"},
		{"Fahrenheit", Temperature{77, Fahrenheit}, "77.00°F"},
		{"Kelvin", Temperature{298.15, Kelvin}, "298.15 K"},
		{"Negative", Temperature{-10, Celsius}, "-10.00°C"},
		{"Decimal", Temperature{25.5, Celsius}, "25.50°C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.temp.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestTemperature_Equals(t *testing.T) {
	temp1, _ := NewTemperature(25, Celsius)
	temp2, _ := NewTemperature(77, Fahrenheit) // 25°C = 77°F
	temp3, _ := NewTemperature(298.15, Kelvin) // 25°C = 298.15K
	temp4, _ := NewTemperature(30, Celsius)

	if !temp1.Equals(temp2) {
		t.Errorf("Expected 25°C to equal 77°F")
	}

	if !temp1.Equals(temp3) {
		t.Errorf("Expected 25°C to equal 298.15K")
	}

	if temp1.Equals(temp4) {
		t.Errorf("Expected 25°C to not equal 30°C")
	}
}

func TestTemperature_IsEmpty(t *testing.T) {
	emptyTemp := Temperature{}
	temp, _ := NewTemperature(25, Celsius)

	if !emptyTemp.IsEmpty() {
		t.Errorf("Expected empty temperature to be empty")
	}

	if temp.IsEmpty() {
		t.Errorf("Expected non-empty temperature to not be empty")
	}
}

func TestTemperature_Value(t *testing.T) {
	temp, _ := NewTemperature(25.5, Celsius)
	if temp.Value() != 25.5 {
		t.Errorf("Expected value 25.5, got %f", temp.Value())
	}
}

func TestTemperature_Unit(t *testing.T) {
	temp, _ := NewTemperature(25, Celsius)
	if temp.Unit() != Celsius {
		t.Errorf("Expected unit Celsius, got %s", temp.Unit())
	}
}

func TestTemperature_Conversions(t *testing.T) {
	// Test Celsius to other units
	celsiusTemp, _ := NewTemperature(25, Celsius)

	fahrenheitFromC := celsiusTemp.ToFahrenheit()
	expectedF := 77.0
	if math.Abs(fahrenheitFromC.value-expectedF) > 0.01 {
		t.Errorf("Expected 25°C to be %.2f°F, got %.2f°F", expectedF, fahrenheitFromC.value)
	}

	kelvinFromC := celsiusTemp.ToKelvin()
	expectedK := 298.15
	if math.Abs(kelvinFromC.value-expectedK) > 0.01 {
		t.Errorf("Expected 25°C to be %.2fK, got %.2fK", expectedK, kelvinFromC.value)
	}

	// Test Fahrenheit to other units
	fahrenheitTemp, _ := NewTemperature(77, Fahrenheit)

	celsiusFromF := fahrenheitTemp.ToCelsius()
	expectedC := 25.0
	if math.Abs(celsiusFromF.value-expectedC) > 0.01 {
		t.Errorf("Expected 77°F to be %.2f°C, got %.2f°C", expectedC, celsiusFromF.value)
	}

	kelvinFromF := fahrenheitTemp.ToKelvin()
	if math.Abs(kelvinFromF.value-expectedK) > 0.01 {
		t.Errorf("Expected 77°F to be %.2fK, got %.2fK", expectedK, kelvinFromF.value)
	}

	// Test Kelvin to other units
	kelvinTemp, _ := NewTemperature(298.15, Kelvin)

	celsiusFromK := kelvinTemp.ToCelsius()
	if math.Abs(celsiusFromK.value-expectedC) > 0.01 {
		t.Errorf("Expected 298.15K to be %.2f°C, got %.2f°C", expectedC, celsiusFromK.value)
	}

	fahrenheitFromK := kelvinTemp.ToFahrenheit()
	if math.Abs(fahrenheitFromK.value-expectedF) > 0.01 {
		t.Errorf("Expected 298.15K to be %.2f°F, got %.2f°F", expectedF, fahrenheitFromK.value)
	}

	// Test conversion to same unit
	celsiusSame := celsiusTemp.ToCelsius()
	if celsiusSame.value != celsiusTemp.value || celsiusSame.unit != celsiusTemp.unit {
		t.Errorf("Converting to same unit should return the same temperature")
	}
}

func TestTemperature_Add(t *testing.T) {
	temp1, _ := NewTemperature(25, Celsius)
	temp2, _ := NewTemperature(10, Celsius)
	temp3, _ := NewTemperature(50, Fahrenheit)

	// Add same unit
	result1 := temp1.Add(temp2)
	if result1.value != 35 || result1.unit != Celsius {
		t.Errorf("Expected 25°C + 10°C = 35°C, got %.2f°%s", result1.value, result1.unit)
	}

	// Add different unit
	result2 := temp1.Add(temp3)
	if math.Abs(result2.value-35) > 0.01 || result2.unit != Celsius {
		t.Errorf("Expected 25°C + 50°F = 35°C, got %.2f°%s", result2.value, result2.unit)
	}
}

func TestTemperature_Subtract(t *testing.T) {
	temp1, _ := NewTemperature(25, Celsius)
	temp2, _ := NewTemperature(10, Celsius)
	temp3, _ := NewTemperature(50, Fahrenheit)

	// Subtract same unit
	result1, err1 := temp1.Subtract(temp2)
	if err1 != nil {
		t.Errorf("Unexpected error: %v", err1)
	}
	if result1.value != 15 || result1.unit != Celsius {
		t.Errorf("Expected 25°C - 10°C = 15°C, got %.2f°%s", result1.value, result1.unit)
	}

	// Subtract different unit
	result2, err2 := temp1.Subtract(temp3)
	if err2 != nil {
		t.Errorf("Unexpected error: %v", err2)
	}
	if math.Abs(result2.value-15) > 0.01 || result2.unit != Celsius {
		t.Errorf("Expected 25°C - 50°F = 15°C, got %.2f°%s", result2.value, result2.unit)
	}

	// Create a temperature close to absolute zero
	tempNearZero, _ := NewTemperature(-273, Celsius)
	tempHigh, _ := NewTemperature(1, Celsius)

	// Subtract to below absolute zero
	_, err3 := tempNearZero.Subtract(tempHigh)
	if err3 == nil {
		t.Errorf("Expected error when subtracting to below absolute zero")
	}
}

func TestTemperature_IsFreezing(t *testing.T) {
	tests := []struct {
		name     string
		temp     Temperature
		expected bool
	}{
		{"Below Freezing Celsius", Temperature{-5, Celsius}, true},
		{"At Freezing Celsius", Temperature{0, Celsius}, true},
		{"Above Freezing Celsius", Temperature{5, Celsius}, false},
		{"Below Freezing Fahrenheit", Temperature{30, Fahrenheit}, true},
		{"At Freezing Fahrenheit", Temperature{32, Fahrenheit}, true},
		{"Above Freezing Fahrenheit", Temperature{40, Fahrenheit}, false},
		{"Below Freezing Kelvin", Temperature{270, Kelvin}, true},
		{"At Freezing Kelvin", Temperature{273.15, Kelvin}, true},
		{"Above Freezing Kelvin", Temperature{280, Kelvin}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.temp.IsFreezing()
			if result != tt.expected {
				t.Errorf("Expected IsFreezing() to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTemperature_IsBoiling(t *testing.T) {
	tests := []struct {
		name     string
		temp     Temperature
		expected bool
	}{
		{"Below Boiling Celsius", Temperature{95, Celsius}, false},
		{"At Boiling Celsius", Temperature{100, Celsius}, true},
		{"Above Boiling Celsius", Temperature{105, Celsius}, true},
		{"Below Boiling Fahrenheit", Temperature{210, Fahrenheit}, false},
		{"At Boiling Fahrenheit", Temperature{212, Fahrenheit}, true},
		{"Above Boiling Fahrenheit", Temperature{220, Fahrenheit}, true},
		{"Below Boiling Kelvin", Temperature{370, Kelvin}, false},
		{"At Boiling Kelvin", Temperature{373.15, Kelvin}, true},
		{"Above Boiling Kelvin", Temperature{380, Kelvin}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.temp.IsBoiling()
			if result != tt.expected {
				t.Errorf("Expected IsBoiling() to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTemperature_Format(t *testing.T) {
	temp, _ := NewTemperature(25, Celsius)

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{"Short Format Celsius", "short", "25.00°C"},
		{"Long Format Celsius", "long", "25.00 degrees Celsius"},
		{"Scientific Format Celsius", "scientific", "298.15K"},
		{"Default Format Celsius", "invalid", "25.00°C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := temp.Format(tt.format)
			if result != tt.expected {
				t.Errorf("Expected Format(%s) to be %s, got %s", tt.format, tt.expected, result)
			}
		})
	}

	// Test other units
	tempF, _ := NewTemperature(77, Fahrenheit)
	if tempF.Format("short") != "77.00°F" {
		t.Errorf("Expected Format(short) for Fahrenheit to be 77.00°F, got %s", tempF.Format("short"))
	}
	if tempF.Format("long") != "77.00 degrees Fahrenheit" {
		t.Errorf("Expected Format(long) for Fahrenheit to be 77.00 degrees Fahrenheit, got %s", tempF.Format("long"))
	}

	tempK, _ := NewTemperature(298.15, Kelvin)
	if tempK.Format("short") != "298.15K" {
		t.Errorf("Expected Format(short) for Kelvin to be 298.15K, got %s", tempK.Format("short"))
	}
	if tempK.Format("long") != "298.15 Kelvin" {
		t.Errorf("Expected Format(long) for Kelvin to be 298.15 Kelvin, got %s", tempK.Format("long"))
	}
}
