// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"fmt"
	"github.com/abitofhelp/servicelib/valueobject"
	"github.com/shopspring/decimal"
)

func main() {
	// Create a new Money object
	money, err := valueobject.NewMoneyFromFloat64(10.99, "USD")
	if err != nil {
		fmt.Printf("Error creating money: %v\n", err)
		return
	}
	fmt.Printf("Money: %s\n", money.String())

	// Create another Money object from string
	moneyFromString, err := valueobject.NewMoneyFromString("5.99", "USD")
	if err != nil {
		fmt.Printf("Error creating money from string: %v\n", err)
		return
	}
	fmt.Printf("Money from string: %s\n", moneyFromString.String())

	// Add money
	sum, err := money.Add(moneyFromString)
	if err != nil {
		fmt.Printf("Error adding money: %v\n", err)
		return
	}
	fmt.Printf("Sum: %s\n", sum.String())

	// Subtract money
	diff, err := money.Subtract(moneyFromString)
	if err != nil {
		fmt.Printf("Error subtracting money: %v\n", err)
		return
	}
	fmt.Printf("Difference: %s\n", diff.String())

	// Multiply money
	factor := decimal.NewFromFloat(1.1) // 10% increase
	product := money.Multiply(factor)
	fmt.Printf("Product: %s\n", product.String())

	// Divide money
	divisor := decimal.NewFromInt(2) // Split in half
	quotient, err := money.Divide(divisor)
	if err != nil {
		fmt.Printf("Error dividing money: %v\n", err)
		return
	}
	fmt.Printf("Quotient: %s\n", quotient.String())

	// Compare money
	isGreater, err := money.IsGreaterThan(moneyFromString)
	if err != nil {
		fmt.Printf("Error comparing money: %v\n", err)
		return
	}
	fmt.Printf("Is %s greater than %s? %t\n", money.String(), moneyFromString.String(), isGreater)

	// Parse money from string
	parsedMoney, err := valueobject.Parse("15.75 EUR")
	if err != nil {
		fmt.Printf("Error parsing money: %v\n", err)
		return
	}
	fmt.Printf("Parsed money: %s\n", parsedMoney.String())

	// Demonstrate precision with decimal
	// This would have precision issues with float64
	preciseA, _ := valueobject.NewMoneyFromString("0.1", "USD")
	preciseB, _ := valueobject.NewMoneyFromString("0.2", "USD")
	preciseSum, _ := preciseA.Add(preciseB)
	fmt.Printf("0.1 + 0.2 = %s (with decimal precision)\n", preciseSum.String())

	// Demonstrate rounding
	roundedMoney := money.Round(1) // Round to 1 decimal place
	fmt.Printf("Rounded money: %s\n", roundedMoney.String())
}
