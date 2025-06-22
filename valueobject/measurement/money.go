// Copyright (c) 2025 A Bit of Help, Inc.

// Package measurement provides value objects related to measurement information.
package measurement

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"regexp"
	"strings"

	"github.com/abitofhelp/servicelib/valueobject/base"
)

// Money represents a monetary value object with amount and currency
type Money struct {
	base.BaseStructValueObject

	amount decimal.Decimal

	currency string
}

// NewMoney creates a new Money with validation
func NewMoney(amount decimal.Decimal, currency string) (Money, error) {
	vo := Money{

		amount: amount,

		currency: currency,
	}

	// Validate the value object
	if err := vo.Validate(); err != nil {
		return Money{}, err
	}

	return vo, nil
}

// String returns the string representation of the Money
func (v Money) String() string {
	return fmt.Sprintf("amount=%v currency=%v", v.amount, v.currency)
}

// Equals checks if two Moneys are equal
func (v Money) Equals(other Money) bool {

	if !v.amount.Equal(other.amount) {
		return false
	}

	if !base.StringsEqualFold(v.currency, other.currency) {
		return false
	}

	return true
}

// IsEmpty checks if the Money is empty (zero value)
func (v Money) IsEmpty() bool {
	return v.amount.IsZero() && v.currency == ""
}

// Validate checks if the Money is valid
func (v Money) Validate() error {

	// Trim whitespace from currency
	trimmedCurrency := strings.TrimSpace(v.currency)

	// Currency is required
	if trimmedCurrency == "" {
		return errors.New("currency cannot be empty")
	}

	// Basic currency code validation (assuming 3-letter ISO currency codes)
	if len(trimmedCurrency) != 3 {
		return errors.New("currency must be a 3-letter ISO code")
	}

	return nil
}

// ToMap converts the Money to a map[string]interface{}
func (v Money) ToMap() map[string]interface{} {
	return map[string]interface{}{

		"amount": v.amount,

		"currency": v.currency,
	}
}

// MarshalJSON implements the json.Marshaler interface
func (v Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.ToMap())
}

// Amount returns the amount as a decimal.Decimal
func (v Money) Amount() decimal.Decimal {
	return v.amount
}

// Currency returns the currency code
func (v Money) Currency() string {
	return v.currency
}

// AmountInCents returns the amount in the smallest currency unit (e.g., cents)
func (v Money) AmountInCents() int64 {
	// Multiply by 100 and convert to int64
	cents := v.amount.Mul(decimal.NewFromInt(100)).IntPart()
	return cents
}

// Add adds another Money value and returns a new Money
// Both Money values must have the same currency
func (v Money) Add(other Money) (Money, error) {
	if v.currency != other.currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}

	return Money{
		amount:   v.amount.Add(other.amount),
		currency: v.currency,
	}, nil
}

// Subtract subtracts another Money value and returns a new Money
// Both Money values must have the same currency
func (v Money) Subtract(other Money) (Money, error) {
	if v.currency != other.currency {
		return Money{}, errors.New("cannot subtract money with different currencies")
	}

	return Money{
		amount:   v.amount.Sub(other.amount),
		currency: v.currency,
	}, nil
}

// Multiply multiplies the money amount by a factor (e.g., for calculating interest)
// and returns a new Money object
func (v Money) Multiply(factor decimal.Decimal) Money {
	return Money{
		amount:   v.amount.Mul(factor),
		currency: v.currency,
	}
}

// Divide divides the money amount by a divisor (e.g., for splitting payments)
// and returns a new Money object
// Returns an error if the divisor is zero
func (v Money) Divide(divisor decimal.Decimal) (Money, error) {
	if divisor.IsZero() {
		return Money{}, errors.New("cannot divide by zero")
	}

	return Money{
		amount:   v.amount.Div(divisor),
		currency: v.currency,
	}, nil
}

// Abs returns the absolute value of the Money
func (v Money) Abs() Money {
	return Money{
		amount:   v.amount.Abs(),
		currency: v.currency,
	}
}

// Negate returns the negated value of the Money
func (v Money) Negate() Money {
	return Money{
		amount:   v.amount.Neg(),
		currency: v.currency,
	}
}

// IsPositive checks if the Money amount is positive
func (v Money) IsPositive() bool {
	return v.amount.IsPositive()
}

// IsNegative checks if the Money amount is negative
func (v Money) IsNegative() bool {
	return v.amount.IsNegative()
}

// IsGreaterThan checks if the Money amount is greater than another Money amount
// Returns an error if the currencies are different
func (v Money) IsGreaterThan(other Money) (bool, error) {
	if v.currency != other.currency {
		return false, errors.New("cannot compare money with different currencies")
	}

	return v.amount.GreaterThan(other.amount), nil
}

// IsLessThan checks if the Money amount is less than another Money amount
// Returns an error if the currencies are different
func (v Money) IsLessThan(other Money) (bool, error) {
	if v.currency != other.currency {
		return false, errors.New("cannot compare money with different currencies")
	}

	return v.amount.LessThan(other.amount), nil
}

// NewMoneyFromFloat64 creates a new Money with validation from a float64 amount
func NewMoneyFromFloat64(amount float64, currency string) (Money, error) {
	// Convert float64 to decimal
	// This avoids floating-point precision issues
	decimalAmount := decimal.NewFromFloat(amount)

	return NewMoney(decimalAmount, currency)
}

// NewMoneyFromString creates a new Money with validation from a string amount
// This is more precise than using a float64 for monetary values
func NewMoneyFromString(amount string, currency string) (Money, error) {
	// Parse the amount as decimal
	decimalAmount, err := decimal.NewFromString(amount)
	if err != nil {
		return Money{}, errors.New("invalid amount format")
	}

	return NewMoney(decimalAmount, currency)
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
