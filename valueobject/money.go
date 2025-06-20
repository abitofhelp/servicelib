// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"regexp"
	"strings"
)

// Money represents a monetary value object with amount and currency
// Uses decimal.Decimal to store amount with arbitrary precision
// to avoid floating-point precision issues
//
// This implementation uses the github.com/shopspring/decimal package instead of int64
// for the following reasons:
// 1. Decimal arithmetic is more precise than floating-point arithmetic for monetary values
// 2. It can handle arbitrary precision, which is important for financial calculations
// 3. It provides methods for rounding, scaling, and other operations that are useful for monetary values
// 4. It avoids common floating-point errors like 0.1 + 0.2 != 0.3
//
// Note: This implementation requires the github.com/shopspring/decimal package to be added
// to the project's dependencies with: go get github.com/shopspring/decimal
type Money struct {
	// amount stored as a decimal value
	amount   decimal.Decimal
	currency string
}

// NewMoneyFromFloat64 creates a new Money with validation from a float64 amount
func NewMoneyFromFloat64(amount float64, currency string) (Money, error) {
	// Trim whitespace from currency
	trimmedCurrency := strings.TrimSpace(currency)

	// Currency is required
	if trimmedCurrency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}

	// Basic currency code validation (assuming 3-letter ISO currency codes)
	if len(trimmedCurrency) != 3 {
		return Money{}, errors.New("currency must be a 3-letter ISO code")
	}

	// Convert currency to uppercase
	trimmedCurrency = strings.ToUpper(trimmedCurrency)

	// Convert float64 to decimal
	// This avoids floating-point precision issues
	decimalAmount := decimal.NewFromFloat(amount)

	return Money{
		amount:   decimalAmount,
		currency: trimmedCurrency,
	}, nil
}

// NewMoneyFromString creates a new Money with validation from a string amount
// This is more precise than using a float64 for monetary values
func NewMoneyFromString(amount string, currency string) (Money, error) {
	// Trim whitespace from currency
	trimmedCurrency := strings.TrimSpace(currency)

	// Currency is required
	if trimmedCurrency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}

	// Basic currency code validation (assuming 3-letter ISO currency codes)
	if len(trimmedCurrency) != 3 {
		return Money{}, errors.New("currency must be a 3-letter ISO code")
	}

	// Convert currency to uppercase
	trimmedCurrency = strings.ToUpper(trimmedCurrency)

	// Parse the amount as decimal
	decimalAmount, err := decimal.NewFromString(amount)
	if err != nil {
		return Money{}, errors.New("invalid amount format")
	}

	return Money{
		amount:   decimalAmount,
		currency: trimmedCurrency,
	}, nil
}

// NewMoneyFromUint64 creates a new Money with validation from a uint64 amount
// This is useful for representing monetary values in the smallest currency unit (e.g., cents)
func NewMoneyFromUint64(amount uint64, currency string) (Money, error) {
	// Trim whitespace from currency
	trimmedCurrency := strings.TrimSpace(currency)

	// Currency is required
	if trimmedCurrency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}

	// Basic currency code validation (assuming 3-letter ISO currency codes)
	if len(trimmedCurrency) != 3 {
		return Money{}, errors.New("currency must be a 3-letter ISO code")
	}

	// Convert currency to uppercase
	trimmedCurrency = strings.ToUpper(trimmedCurrency)

	// Check if the uint64 value is too large to fit in an int64
	if amount > uint64(9223372036854775807) { // Max int64 value
		// Handle as overflow case by setting amount to -1
		return Money{
			amount:   decimal.NewFromInt(-1),
			currency: trimmedCurrency,
		}, nil
	}

	// Convert uint64 to decimal
	decimalAmount := decimal.NewFromUint64(amount)

	return Money{
		amount:   decimalAmount,
		currency: trimmedCurrency,
	}, nil
}

// AmountInCents returns the amount in the smallest currency unit (e.g., cents)
func (m Money) AmountInCents() int64 {
	// Multiply by 100 and convert to int64
	cents := m.amount.Mul(decimal.NewFromInt(100)).IntPart()
	return cents
}

// Amount returns the amount as a float64 in the standard currency unit (e.g., dollars)
// This is for backward compatibility with the previous API
func (m Money) Amount() float64 {
	// Convert decimal to float64
	amount, _ := m.amount.Float64()
	return amount
}

// Currency returns the currency code
func (m Money) Currency() string {
	return m.currency
}

// String returns the string representation of the Money
func (m Money) String() string {
	// Format with 2 decimal places
	return fmt.Sprintf("%s %s", m.amount.StringFixed(2), m.currency)
}

// Equals checks if two Money values are equal
func (m Money) Equals(other Money) bool {
	return m.amount.Equal(other.amount) && m.currency == other.currency
}

// Add adds another Money value and returns a new Money
// Both Money values must have the same currency
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}

	return Money{
		amount:   m.amount.Add(other.amount),
		currency: m.currency,
	}, nil
}

// Subtract subtracts another Money value and returns a new Money
// Both Money values must have the same currency
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot subtract money with different currencies")
	}

	return Money{
		amount:   m.amount.Sub(other.amount),
		currency: m.currency,
	}, nil
}

// IsZero checks if the Money amount is zero
func (m Money) IsZero() bool {
	return m.amount.IsZero()
}

// Round rounds the money amount to the specified number of decimal places
// and returns a new Money object
func (m Money) Round(places int32) Money {
	return Money{
		amount:   m.amount.Round(places),
		currency: m.currency,
	}
}

// Scale scales the money amount by a factor (e.g., for currency conversion)
// and returns a new Money with the specified currency
func (m Money) Scale(factor decimal.Decimal, newCurrency string) (Money, error) {
	// Validate currency
	trimmedCurrency := strings.TrimSpace(newCurrency)
	if trimmedCurrency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	if len(trimmedCurrency) != 3 {
		return Money{}, errors.New("currency must be a 3-letter ISO code")
	}

	// Convert currency to uppercase
	trimmedCurrency = strings.ToUpper(trimmedCurrency)

	// Scale the amount
	scaledAmount := m.amount.Mul(factor)

	return Money{
		amount:   scaledAmount,
		currency: trimmedCurrency,
	}, nil
}

// Multiply multiplies the money amount by a factor (e.g., for calculating interest)
// and returns a new Money object
func (m Money) Multiply(factor decimal.Decimal) Money {
	return Money{
		amount:   m.amount.Mul(factor),
		currency: m.currency,
	}
}

// Divide divides the money amount by a divisor (e.g., for splitting payments)
// and returns a new Money object
// Returns an error if the divisor is zero
func (m Money) Divide(divisor decimal.Decimal) (Money, error) {
	if divisor.IsZero() {
		return Money{}, errors.New("cannot divide by zero")
	}

	return Money{
		amount:   m.amount.Div(divisor),
		currency: m.currency,
	}, nil
}

// Min returns the smaller of two Money values
// Both Money values must have the same currency
func Min(a, b Money) (Money, error) {
	if a.currency != b.currency {
		return Money{}, errors.New("cannot compare money with different currencies")
	}

	if a.amount.LessThan(b.amount) {
		return a, nil
	}
	return b, nil
}

// Max returns the larger of two Money values
// Both Money values must have the same currency
func Max(a, b Money) (Money, error) {
	if a.currency != b.currency {
		return Money{}, errors.New("cannot compare money with different currencies")
	}

	if a.amount.GreaterThan(b.amount) {
		return a, nil
	}
	return b, nil
}

// Abs returns the absolute value of the Money
func (m Money) Abs() Money {
	return Money{
		amount:   m.amount.Abs(),
		currency: m.currency,
	}
}

// Negate returns the negated value of the Money
func (m Money) Negate() Money {
	return Money{
		amount:   m.amount.Neg(),
		currency: m.currency,
	}
}

// IsPositive checks if the Money amount is positive
func (m Money) IsPositive() bool {
	return m.amount.IsPositive()
}

// IsNegative checks if the Money amount is negative
func (m Money) IsNegative() bool {
	return m.amount.IsNegative()
}

// IsGreaterThan checks if the Money amount is greater than another Money amount
// Returns an error if the currencies are different
func (m Money) IsGreaterThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, errors.New("cannot compare money with different currencies")
	}

	return m.amount.GreaterThan(other.amount), nil
}

// IsLessThan checks if the Money amount is less than another Money amount
// Returns an error if the currencies are different
func (m Money) IsLessThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, errors.New("cannot compare money with different currencies")
	}

	return m.amount.LessThan(other.amount), nil
}

// MoneyJSON is a struct used for JSON marshaling and unmarshaling
type MoneyJSON struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// MarshalJSON implements the json.Marshaler interface
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(MoneyJSON{
		Amount:   m.amount.String(),
		Currency: m.currency,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (m *Money) UnmarshalJSON(data []byte) error {
	var moneyJSON MoneyJSON
	if err := json.Unmarshal(data, &moneyJSON); err != nil {
		return err
	}

	// Parse the amount
	decimalAmount, err := decimal.NewFromString(moneyJSON.Amount)
	if err != nil {
		return errors.New("invalid amount format in JSON")
	}

	// Validate currency
	trimmedCurrency := strings.TrimSpace(moneyJSON.Currency)
	if trimmedCurrency == "" {
		return errors.New("currency cannot be empty")
	}
	if len(trimmedCurrency) != 3 {
		return errors.New("currency must be a 3-letter ISO code")
	}

	// Set the values
	m.amount = decimalAmount
	m.currency = strings.ToUpper(trimmedCurrency)

	return nil
}

// Parse creates a Money object from a string representation like "10.99 USD"
func Parse(s string) (Money, error) {
	// Use regular expression to extract amount and currency
	re := regexp.MustCompile(`^\s*(-?\d+(?:\.\d+)?)\s+([A-Za-z]{3})\s*$`)
	matches := re.FindStringSubmatch(s)

	if matches == nil || len(matches) != 3 {
		return Money{}, errors.New("invalid money format, expected 'amount currency'")
	}

	// Parse the amount as decimal
	decimalAmount, err := decimal.NewFromString(matches[1])
	if err != nil {
		return Money{}, errors.New("invalid amount format")
	}

	return Money{
		amount:   decimalAmount,
		currency: strings.ToUpper(matches[2]),
	}, nil
}
