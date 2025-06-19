// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"testing"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		currency    string
		expectError bool
	}{
		{"Valid Money", 100.50, "USD", false},
		{"Zero Amount", 0, "EUR", false},
		{"Negative Amount", -10.25, "GBP", false},
		{"Empty Currency", 100, "", true},
		{"Invalid Currency Length", 100, "USDD", true},
		{"Currency With Spaces", 100, " USD ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoney(tt.amount, tt.currency)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				amount, _ := money.amount.Float64()
				if amount != tt.amount {
					t.Errorf("Expected amount %f, got %f", tt.amount, amount)
				}
				if money.currency != tt.currency && money.currency != "USD" {
					t.Errorf("Expected currency %s, got %s", tt.currency, money.currency)
				}
			}
		})
	}
}

func TestNewMoneyFromString(t *testing.T) {
	tests := []struct {
		name        string
		amount      string
		currency    string
		expectError bool
	}{
		{"Valid Money", "100.50", "USD", false},
		{"Zero Amount", "0", "EUR", false},
		{"Negative Amount", "-10.25", "GBP", false},
		{"Empty Currency", "100", "", true},
		{"Invalid Currency Length", "100", "USDD", true},
		{"Invalid Amount Format", "abc", "USD", true},
		{"Currency With Spaces", "100", " USD ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoneyFromString(tt.amount, tt.currency)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				expectedAmount, _ := decimal.NewFromString(tt.amount)
				if !money.amount.Equal(expectedAmount) {
					t.Errorf("Expected amount %s, got %s", expectedAmount, money.amount)
				}
				if money.currency != tt.currency && money.currency != "USD" {
					t.Errorf("Expected currency %s, got %s", tt.currency, money.currency)
				}
			}
		})
	}
}

func TestMoney_AmountInCents(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected int64
	}{
		{"Whole Number", 100, 10000},
		{"Decimal Number", 100.50, 10050},
		{"Zero", 0, 0},
		{"Negative", -10.25, -1025},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, _ := NewMoney(tt.amount, "USD")
			cents := money.AmountInCents()
			if cents != tt.expected {
				t.Errorf("Expected %d cents, got %d", tt.expected, cents)
			}
		})
	}
}

func TestMoney_Amount(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected float64
	}{
		{"Whole Number", 100, 100},
		{"Decimal Number", 100.50, 100.50},
		{"Zero", 0, 0},
		{"Negative", -10.25, -10.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, _ := NewMoney(tt.amount, "USD")
			amount := money.Amount()
			if amount != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, amount)
			}
		})
	}
}

func TestMoney_Currency(t *testing.T) {
	money, _ := NewMoney(100, "USD")
	if money.Currency() != "USD" {
		t.Errorf("Expected USD, got %s", money.Currency())
	}
}

func TestMoney_String(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		currency string
		expected string
	}{
		{"Whole Number", 100, "USD", "100.00 USD"},
		{"Decimal Number", 100.50, "EUR", "100.50 EUR"},
		{"Zero", 0, "GBP", "0.00 GBP"},
		{"Negative", -10.25, "JPY", "-10.25 JPY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, _ := NewMoney(tt.amount, tt.currency)
			str := money.String()
			if str != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, str)
			}
		})
	}
}

func TestMoney_Equals(t *testing.T) {
	money1, _ := NewMoney(100, "USD")
	money2, _ := NewMoney(100, "USD")
	money3, _ := NewMoney(200, "USD")
	money4, _ := NewMoney(100, "EUR")

	if !money1.Equals(money2) {
		t.Errorf("Expected money1 to equal money2")
	}

	if money1.Equals(money3) {
		t.Errorf("Expected money1 to not equal money3")
	}

	if money1.Equals(money4) {
		t.Errorf("Expected money1 to not equal money4")
	}
}

func TestMoney_Add(t *testing.T) {
	money1, _ := NewMoney(100, "USD")
	money2, _ := NewMoney(50, "USD")
	money3, _ := NewMoney(50, "EUR")

	result, err := money1.Add(money2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.Amount() != 150 {
		t.Errorf("Expected 150, got %f", result.Amount())
	}
	if result.Currency() != "USD" {
		t.Errorf("Expected USD, got %s", result.Currency())
	}

	_, err = money1.Add(money3)
	if err == nil {
		t.Errorf("Expected error when adding different currencies")
	}
}

func TestMoney_Subtract(t *testing.T) {
	money1, _ := NewMoney(100, "USD")
	money2, _ := NewMoney(50, "USD")
	money3, _ := NewMoney(50, "EUR")

	result, err := money1.Subtract(money2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.Amount() != 50 {
		t.Errorf("Expected 50, got %f", result.Amount())
	}
	if result.Currency() != "USD" {
		t.Errorf("Expected USD, got %s", result.Currency())
	}

	_, err = money1.Subtract(money3)
	if err == nil {
		t.Errorf("Expected error when subtracting different currencies")
	}
}

func TestMoney_IsZero(t *testing.T) {
	money1, _ := NewMoney(0, "USD")
	money2, _ := NewMoney(100, "USD")

	if !money1.IsZero() {
		t.Errorf("Expected money1 to be zero")
	}

	if money2.IsZero() {
		t.Errorf("Expected money2 to not be zero")
	}
}

func TestMoney_Round(t *testing.T) {
	money, _ := NewMoney(100.567, "USD")
	rounded := money.Round(2)

	if rounded.Amount() != 100.57 {
		t.Errorf("Expected 100.57, got %f", rounded.Amount())
	}
}

func TestMoney_Scale(t *testing.T) {
	money, _ := NewMoney(100, "USD")
	factor := decimal.NewFromFloat(1.5)

	scaled, err := money.Scale(factor, "EUR")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if scaled.Amount() != 150 {
		t.Errorf("Expected 150, got %f", scaled.Amount())
	}
	if scaled.Currency() != "EUR" {
		t.Errorf("Expected EUR, got %s", scaled.Currency())
	}

	_, err = money.Scale(factor, "")
	if err == nil {
		t.Errorf("Expected error with empty currency")
	}

	_, err = money.Scale(factor, "USDD")
	if err == nil {
		t.Errorf("Expected error with invalid currency length")
	}
}

func TestMoney_Multiply(t *testing.T) {
	money, _ := NewMoney(100, "USD")
	factor := decimal.NewFromFloat(1.5)

	result := money.Multiply(factor)
	if result.Amount() != 150 {
		t.Errorf("Expected 150, got %f", result.Amount())
	}
	if result.Currency() != "USD" {
		t.Errorf("Expected USD, got %s", result.Currency())
	}
}

func TestMoney_Divide(t *testing.T) {
	money, _ := NewMoney(100, "USD")
	divisor := decimal.NewFromFloat(2)
	zeroDivisor := decimal.NewFromFloat(0)

	result, err := money.Divide(divisor)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.Amount() != 50 {
		t.Errorf("Expected 50, got %f", result.Amount())
	}
	if result.Currency() != "USD" {
		t.Errorf("Expected USD, got %s", result.Currency())
	}

	_, err = money.Divide(zeroDivisor)
	if err == nil {
		t.Errorf("Expected error when dividing by zero")
	}
}

func TestMin(t *testing.T) {
	money1, _ := NewMoney(100, "USD")
	money2, _ := NewMoney(50, "USD")
	money3, _ := NewMoney(50, "EUR")

	min, err := Min(money1, money2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if min.Amount() != 50 {
		t.Errorf("Expected 50, got %f", min.Amount())
	}

	_, err = Min(money1, money3)
	if err == nil {
		t.Errorf("Expected error when comparing different currencies")
	}
}

func TestMax(t *testing.T) {
	money1, _ := NewMoney(100, "USD")
	money2, _ := NewMoney(50, "USD")
	money3, _ := NewMoney(50, "EUR")

	max, err := Max(money1, money2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if max.Amount() != 100 {
		t.Errorf("Expected 100, got %f", max.Amount())
	}

	_, err = Max(money1, money3)
	if err == nil {
		t.Errorf("Expected error when comparing different currencies")
	}
}

func TestMoney_Abs(t *testing.T) {
	money, _ := NewMoney(-100, "USD")
	abs := money.Abs()

	if abs.Amount() != 100 {
		t.Errorf("Expected 100, got %f", abs.Amount())
	}
}

func TestMoney_Negate(t *testing.T) {
	money, _ := NewMoney(100, "USD")
	negated := money.Negate()

	if negated.Amount() != -100 {
		t.Errorf("Expected -100, got %f", negated.Amount())
	}
}

func TestMoney_IsPositive(t *testing.T) {
	money1, _ := NewMoney(100, "USD")
	money2, _ := NewMoney(0, "USD")
	money3, _ := NewMoney(-100, "USD")

	if !money1.IsPositive() {
		t.Errorf("Expected money1 to be positive")
	}
	if money2.IsPositive() {
		t.Errorf("Expected money2 to not be positive")
	}
	if money3.IsPositive() {
		t.Errorf("Expected money3 to not be positive")
	}
}

func TestMoney_IsNegative(t *testing.T) {
	money1, _ := NewMoney(100, "USD")
	money2, _ := NewMoney(0, "USD")
	money3, _ := NewMoney(-100, "USD")

	if money1.IsNegative() {
		t.Errorf("Expected money1 to not be negative")
	}
	if money2.IsNegative() {
		t.Errorf("Expected money2 to not be negative")
	}
	if !money3.IsNegative() {
		t.Errorf("Expected money3 to be negative")
	}
}

func TestMoney_IsGreaterThan(t *testing.T) {
	money1, _ := NewMoney(100, "USD")
	money2, _ := NewMoney(50, "USD")
	money3, _ := NewMoney(50, "EUR")

	result, err := money1.IsGreaterThan(money2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected money1 to be greater than money2")
	}

	_, err = money1.IsGreaterThan(money3)
	if err == nil {
		t.Errorf("Expected error when comparing different currencies")
	}
}

func TestMoney_IsLessThan(t *testing.T) {
	money1, _ := NewMoney(50, "USD")
	money2, _ := NewMoney(100, "USD")
	money3, _ := NewMoney(50, "EUR")

	result, err := money1.IsLessThan(money2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected money1 to be less than money2")
	}

	_, err = money1.IsLessThan(money3)
	if err == nil {
		t.Errorf("Expected error when comparing different currencies")
	}
}

func TestMoney_MarshalJSON(t *testing.T) {
	money, _ := NewMoney(100.50, "USD")
	
	data, err := json.Marshal(money)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	var moneyJSON MoneyJSON
	err = json.Unmarshal(data, &moneyJSON)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if moneyJSON.Amount != "100.5" {
		t.Errorf("Expected amount 100.5, got %s", moneyJSON.Amount)
	}
	if moneyJSON.Currency != "USD" {
		t.Errorf("Expected currency USD, got %s", moneyJSON.Currency)
	}
}

func TestMoney_UnmarshalJSON(t *testing.T) {
	jsonData := []byte(`{"amount":"100.5","currency":"USD"}`)
	
	var money Money
	err := json.Unmarshal(jsonData, &money)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if money.Amount() != 100.5 {
		t.Errorf("Expected amount 100.5, got %f", money.Amount())
	}
	if money.Currency() != "USD" {
		t.Errorf("Expected currency USD, got %s", money.Currency())
	}
	
	// Test invalid JSON
	invalidJSON := []byte(`{"amount":"abc","currency":"USD"}`)
	err = json.Unmarshal(invalidJSON, &money)
	if err == nil {
		t.Errorf("Expected error with invalid amount")
	}
	
	// Test invalid currency
	invalidCurrency := []byte(`{"amount":"100.5","currency":""}`)
	err = json.Unmarshal(invalidCurrency, &money)
	if err == nil {
		t.Errorf("Expected error with empty currency")
	}
	
	invalidCurrencyLength := []byte(`{"amount":"100.5","currency":"USDD"}`)
	err = json.Unmarshal(invalidCurrencyLength, &money)
	if err == nil {
		t.Errorf("Expected error with invalid currency length")
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		amount      string
		currency    string
		expectError bool
	}{
		{"Valid Format", "100.50 USD", "100.50", "USD", false},
		{"Zero Amount", "0 EUR", "0", "EUR", false},
		{"Negative Amount", "-10.25 GBP", "-10.25", "GBP", false},
		{"Invalid Format", "100USD", "", "", true},
		{"Invalid Amount", "abc USD", "", "", true},
		{"With Extra Spaces", "  100.50  USD  ", "100.50", "USD", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := Parse(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				expectedAmount, _ := decimal.NewFromString(tt.amount)
				if !money.amount.Equal(expectedAmount) {
					t.Errorf("Expected amount %s, got %s", expectedAmount, money.amount)
				}
				if money.currency != tt.currency {
					t.Errorf("Expected currency %s, got %s", tt.currency, money.currency)
				}
			}
		})
	}
}