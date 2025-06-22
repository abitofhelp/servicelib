// Copyright (c) 2025 A Bit of Help, Inc.

package base

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStructValueObject is a test implementation of StructValueObject
type TestStructValueObject struct {
	BaseStructValueObject
	Name  string
	Value int
}

func (vo TestStructValueObject) IsEmpty() bool {
	return vo.Name == "" && vo.Value == 0
}

func (vo TestStructValueObject) String() string {
	return vo.Name
}

func (vo TestStructValueObject) Equals(other TestStructValueObject) bool {
	return vo.Name == other.Name && vo.Value == other.Value
}

func TestBaseStructValueObject_IsEmpty(t *testing.T) {
	vo := BaseStructValueObject{}
	assert.False(t, vo.IsEmpty(), "Default IsEmpty should return false")
}

func TestBaseStructValueObject_String(t *testing.T) {
	vo := BaseStructValueObject{}
	assert.Contains(t, vo.String(), "BaseStructValueObject", "String should contain the type name")
}

func TestBaseStructValueObject_ToMap(t *testing.T) {
	// Create a custom implementation of ToMap for TestStructValueObject
	type CustomTestStructValueObject struct {
		BaseStructValueObject
		Name  string
		Value int
	}

	// Override the ToMap method for this test
	customVO := CustomTestStructValueObject{
		Name:  "test",
		Value: 42,
	}

	// The base implementation returns an empty map
	expected := map[string]interface{}{}

	assert.Equal(t, expected, customVO.ToMap())
}

func TestBaseStructValueObject_MarshalJSON(t *testing.T) {
	// Create a custom implementation for this test
	type CustomTestStructValueObject struct {
		BaseStructValueObject
		Name  string
		Value int
	}

	// Override the ToMap method for this test
	customVO := CustomTestStructValueObject{
		Name:  "test",
		Value: 42,
	}

	data, err := customVO.MarshalJSON()
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// The base implementation returns an empty map
	expected := map[string]interface{}{}

	assert.Equal(t, expected, result)
}

func TestWithValidation(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		validateErr error
		expectErr   bool
	}{
		{"Valid value", "test", nil, false},
		{"Invalid value", "invalid", errors.New("invalid value"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test value object
			vo := TestStructValueObject{
				Name:  tt.value,
				Value: 42,
			}

			// Use WithValidation
			result, err := WithValidation(vo, tt.validateErr)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, TestStructValueObject{}, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, vo, result)
			}
		})
	}
}

func TestTestStructValueObject_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		vo       TestStructValueObject
		expected bool
	}{
		{"Empty", TestStructValueObject{}, true},
		{"Name only", TestStructValueObject{Name: "test"}, false},
		{"Value only", TestStructValueObject{Value: 42}, false},
		{"Both fields", TestStructValueObject{Name: "test", Value: 42}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.vo.IsEmpty())
		})
	}
}

func TestTestStructValueObject_String(t *testing.T) {
	tests := []struct {
		name     string
		vo       TestStructValueObject
		expected string
	}{
		{"Empty", TestStructValueObject{}, ""},
		{"With name", TestStructValueObject{Name: "test"}, "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.vo.String())
		})
	}
}

func TestTestStructValueObject_Equals(t *testing.T) {
	tests := []struct {
		name     string
		vo1      TestStructValueObject
		vo2      TestStructValueObject
		expected bool
	}{
		{"Both empty", TestStructValueObject{}, TestStructValueObject{}, true},
		{"Same values", TestStructValueObject{Name: "test", Value: 42}, TestStructValueObject{Name: "test", Value: 42}, true},
		{"Different names", TestStructValueObject{Name: "test", Value: 42}, TestStructValueObject{Name: "other", Value: 42}, false},
		{"Different values", TestStructValueObject{Name: "test", Value: 42}, TestStructValueObject{Name: "test", Value: 43}, false},
		{"Completely different", TestStructValueObject{Name: "test", Value: 42}, TestStructValueObject{Name: "other", Value: 43}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.vo1.Equals(tt.vo2))
		})
	}
}
