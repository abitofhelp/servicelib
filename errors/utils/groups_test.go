// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abitofhelp/servicelib/errors/core"
)

func TestErrorCodeGrouper(t *testing.T) {
	grouper := NewErrorCodeGrouper()
	tests := []struct {
		name    string
		code    core.ErrorCode
		wantCat ErrorCodeGroup
	}{
		{"Resource", core.NotFoundCode, ResourceGroup},
		{"Input", core.InvalidInputCode, InputGroup},
		{"System", core.TimeoutCode, SystemGroup},
		{"Security", core.UnauthorizedCode, SecurityGroup},
		{"External", core.ExternalServiceErrorCode, ExternalGroup},
		{"Business", core.ResourceExhaustedCode, BusinessGroup},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCat := grouper.Group(tt.code)
			assert.Equal(t, tt.wantCat, gotCat)
		})
	}

	// Test group-specific check functions
	assert.True(t, grouper.IsResourceError(core.NotFoundCode))
	assert.True(t, grouper.IsInputError(core.InvalidInputCode))
	assert.True(t, grouper.IsSystemError(core.TimeoutCode))
	assert.True(t, grouper.IsSecurityError(core.UnauthorizedCode))
	assert.True(t, grouper.IsExternalError(core.ExternalServiceErrorCode))
	assert.True(t, grouper.IsBusinessError(core.ResourceExhaustedCode))

	assert.False(t, grouper.IsResourceError(core.InvalidInputCode))
	assert.False(t, grouper.IsInputError(core.TimeoutCode))
	assert.False(t, grouper.IsSystemError(core.UnauthorizedCode))
	assert.False(t, grouper.IsSecurityError(core.ExternalServiceErrorCode))
	assert.False(t, grouper.IsExternalError(core.ResourceExhaustedCode))
	assert.False(t, grouper.IsBusinessError(core.NotFoundCode))
}

func TestErrorCodeGrouper_IsExternalError(t *testing.T) {
	grouper := NewErrorCodeGrouper()
	assert.True(t, grouper.IsExternalError(core.ExternalServiceErrorCode))
	assert.False(t, grouper.IsExternalError(core.NotFoundCode))
}

func TestErrorCodeGrouper_IsBusinessError(t *testing.T) {
	grouper := NewErrorCodeGrouper()
	assert.True(t, grouper.IsBusinessError(core.BusinessRuleViolationCode))
	assert.False(t, grouper.IsBusinessError(core.NotFoundCode))
	assert.False(t, grouper.IsBusinessError(core.NotFoundCode))
}
