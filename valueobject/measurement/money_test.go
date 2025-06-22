// Copyright (c) 2025 A Bit of Help, Inc.

package measurement

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name string

		amount decimal.Decimal

		currency string

		expectError bool
	}{
		{
			name: "Valid Value",

			amount: 42,

			currency: "valid",

			expectError: false,
		},
		// Add more test cases specific to this value object
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := NewMoney(tt.amount, tt.currency)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.amount, value.amount)

				assert.Equal(t, tt.currency, value.currency)

			}
		})
	}
}

func TestMoney_String(t *testing.T) {
	value := Money{

		amount: 42,

		currency: "test",
	}

	// The exact string representation depends on the fields
	assert.Contains(t, value.String(), "amount")
}

func TestMoney_Equals(t *testing.T) {
	value1 := Money{

		amount: 42,

		currency: "test",
	}

	value2 := Money{

		amount: 42,

		currency: "test",
	}

	value3 := Money{

		amount: 43,

		currency: "other",
	}

	assert.True(t, value1.Equals(value2))
	assert.False(t, value1.Equals(value3))
}

func TestMoney_IsEmpty(t *testing.T) {
	emptyValue := Money{}
	assert.True(t, emptyValue.IsEmpty())

	nonEmptyValue := Money{

		amount: 42,

		currency: "test",
	}
	assert.False(t, nonEmptyValue.IsEmpty())
}

func TestMoney_Validate(t *testing.T) {
	tests := []struct {
		name        string
		value       Money
		expectError bool
	}{
		{
			name: "Valid Value",
			value: Money{

				amount: 42,

				currency: "valid",
			},
			expectError: false,
		},
		// Add more test cases specific to this value object
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.value.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMoney_ToMap(t *testing.T) {
	value := Money{

		amount: 42,

		currency: "test",
	}

	expected := map[string]interface{}{

		"amount": 42,

		"currency": "test",
	}

	assert.Equal(t, expected, value.ToMap())
}

func TestMoney_MarshalJSON(t *testing.T) {
	value := Money{

		amount: 42,

		currency: "test",
	}

	data, err := value.MarshalJSON()
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	expected := map[string]interface{}{

		"amount": float64(42),

		"currency": "test",
	}

	assert.Equal(t, expected, result)
}
