// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the Money value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
	"github.com/shopspring/decimal"
)

func main() {
	// Create a new money value
	money, err := valueobject.NewMoneyFromFloat64(10.99, "USD")
	if err != nil {
		// Handle error
		fmt.Println("Error creating money:", err)
		return
	}

	// Access values
	amount := money.Amount()
	currency := money.Currency()
	fmt.Printf("Money: %s %s\n", amount, currency)

	// Create another money value
	otherMoney, _ := valueobject.NewMoneyFromFloat64(5.99, "USD")

	// Add money values
	sum, err := money.Add(otherMoney)
	if err != nil {
		// Handle error (different currencies)
		fmt.Println("Error adding money:", err)
		return
	}
	fmt.Printf("Sum: %s\n", sum)

	// Subtract money values
	diff, err := money.Subtract(otherMoney)
	if err != nil {
		// Handle error (different currencies)
		fmt.Println("Error subtracting money:", err)
		return
	}
	fmt.Printf("Difference: %s\n", diff)

	// Multiply by a factor
	factor := decimal.NewFromFloat(1.1) // 10% increase
	product := money.Multiply(factor)
	fmt.Printf("Product: %s\n", product)

	// Divide by a factor
	divisor := decimal.NewFromInt(2) // Half
	quotient, err := money.Divide(divisor)
	if err != nil {
		// Handle error (division by zero)
		fmt.Println("Error dividing money:", err)
		return
	}
	fmt.Printf("Quotient: %s\n", quotient)

	// Compare money values
	isEqual := money.Equals(otherMoney)
	fmt.Printf("Are equal? %v\n", isEqual)

	isGreater, err := money.IsGreaterThan(otherMoney)
	if err != nil {
		// Handle error (different currencies)
		fmt.Println("Error comparing money:", err)
		return
	}
	fmt.Printf("Is greater? %v\n", isGreater)
}
