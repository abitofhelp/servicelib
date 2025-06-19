# Money Value Object

## Overview
The Money value object in this library has been updated to use the `decimal.Decimal` type from the `github.com/shopspring/decimal` package instead of using `float64` or `int64` for monetary values. This change provides several benefits:

1. **Precise decimal arithmetic**: Avoids floating-point precision issues common with `float64` (e.g., 0.1 + 0.2 = 0.30000000000000004)
2. **Arbitrary precision**: Can handle any number of decimal places without loss of precision
3. **Specialized methods**: The decimal package provides methods for rounding, scaling, and other operations useful for monetary values
4. **Safer financial calculations**: Critical for applications dealing with money where precision is essential

## Installation
The implementation requires the `github.com/shopspring/decimal` package:

```bash
go get github.com/shopspring/decimal
```

## Usage Examples

### Creating Money Objects
```go
// From float64
money, err := valueobject.NewMoney(10.99, "USD")
if err != nil {
    // Handle error
}

// From string (more precise)
moneyFromString, err := valueobject.NewMoneyFromString("5.99", "USD")
if err != nil {
    // Handle error
}

// From string representation
parsedMoney, err := valueobject.Parse("15.75 EUR")
if err != nil {
    // Handle error
}
```

### Basic Operations
```go
// Addition
sum, err := money.Add(otherMoney)
if err != nil {
    // Handle error (e.g., different currencies)
}

// Subtraction
difference, err := money.Subtract(otherMoney)
if err != nil {
    // Handle error (e.g., different currencies)
}

// Multiplication
factor := decimal.NewFromFloat(1.1) // 10% increase
product := money.Multiply(factor)

// Division
divisor := decimal.NewFromInt(2) // Split in half
quotient, err := money.Divide(divisor)
if err != nil {
    // Handle error (e.g., division by zero)
}
```

### Comparison
```go
// Greater than
isGreater, err := money.IsGreaterThan(otherMoney)
if err != nil {
    // Handle error (e.g., different currencies)
}

// Less than
isLess, err := money.IsLessThan(otherMoney)
if err != nil {
    // Handle error (e.g., different currencies)
}

// Equality
isEqual := money.Equals(otherMoney) // Also checks currency equality
```

### Other Operations
```go
// Rounding
roundedMoney := money.Round(2) // Round to 2 decimal places

// Currency conversion
exchangeRate := decimal.NewFromFloat(1.2) // 1 USD = 1.2 EUR
convertedMoney, err := money.Scale(exchangeRate, "EUR")
if err != nil {
    // Handle error
}

// Absolute value
absoluteMoney := money.Abs()

// Negation
negatedMoney := money.Negate()

// Check if positive/negative
isPositive := money.IsPositive()
isNegative := money.IsNegative()
```

### Accessing Components
```go
// Get amount as float64 (for backward compatibility)
amount := money.Amount()

// Get amount in cents
cents := money.AmountInCents()

// Get currency
currency := money.Currency()

// String representation
str := money.String() // e.g., "10.99 USD"
```

### JSON Marshaling/Unmarshaling
The Money type implements the `json.Marshaler` and `json.Unmarshaler` interfaces:

```go
// Marshal to JSON
jsonBytes, err := json.Marshal(money)
if err != nil {
    // Handle error
}

// Unmarshal from JSON
var unmarshaledMoney valueobject.Money
err = json.Unmarshal(jsonBytes, &unmarshaledMoney)
if err != nil {
    // Handle error
}
```

## Benefits of Using Decimal for Money

### Precision Example
```go
// With float64, this would result in precision errors
preciseA, _ := valueobject.NewMoneyFromString("0.1", "USD")
preciseB, _ := valueobject.NewMoneyFromString("0.2", "USD")
preciseSum, _ := preciseA.Add(preciseB)
fmt.Println(preciseSum.String()) // Correctly outputs "0.30 USD"
```

### Rounding Control
```go
// Round to different decimal places
roundedToWhole := money.Round(0)      // e.g., "11 USD"
roundedToTenth := money.Round(1)      // e.g., "10.9 USD"
roundedToHundredth := money.Round(2)  // e.g., "10.99 USD"
```

## Best Practices

1. **Use string input when possible**: When creating Money objects from user input or external sources, prefer `NewMoneyFromString` over `NewMoney` to avoid float64 precision issues.

2. **Always check errors**: All operations that could fail (e.g., adding different currencies) return errors that should be checked.

3. **Use appropriate rounding**: For display purposes, use `Round(2)` for standard currency display, but maintain full precision for calculations.

4. **Validate currencies**: The implementation enforces 3-letter ISO currency codes, but you may need additional validation for specific use cases.

5. **Consider using JSON marshaling**: The implementation provides JSON marshaling/unmarshaling that preserves full decimal precision.