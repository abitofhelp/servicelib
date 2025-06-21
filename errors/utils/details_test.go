// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abitofhelp/servicelib/errors/core"
)

func TestAddDetail(t *testing.T) {
	originalErr := errors.New("original error")
	err := AddDetail(originalErr, "field", "value")
	assert.NotNil(t, err)

	if ce, ok := err.(*core.ContextualError); ok {
		assert.Equal(t, "value", ce.Context.Details["field"])
	}
}

func TestAddDetails(t *testing.T) {
	originalErr := errors.New("original error")
	details := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}
	err := AddDetails(originalErr, details)
	assert.NotNil(t, err)

	if ce, ok := err.(*core.ContextualError); ok {
		assert.Equal(t, details, ce.Context.Details)
	}
}

func TestGetDetail(t *testing.T) {
	originalErr := errors.New("original error")
	err := AddDetail(originalErr, "field", "value")
	assert.NotNil(t, err)

	value, ok := GetDetail(err, "field")
	assert.True(t, ok)
	assert.Equal(t, "value", value)

	_, ok = GetDetail(err, "nonexistent")
	assert.False(t, ok)
}

func TestGetDetails(t *testing.T) {
	originalErr := errors.New("original error")
	details := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}
	err := AddDetails(originalErr, details)
	assert.NotNil(t, err)

	result := GetDetails(err)
	assert.Equal(t, details, result)
}
