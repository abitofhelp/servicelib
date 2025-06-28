// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abitofhelp/servicelib/errors/core"
)

func TestErrorCodeCategorizer(t *testing.T) {
	categorizer := NewErrorCodeCategorizer()
	tests := []struct {
		name    string
		code    core.ErrorCode
		wantCat ErrorCodeCategory
	}{
		{"NotFound", core.NotFoundCode, ClientError},
		{"InvalidInput", core.InvalidInputCode, ClientError},
		{"AlreadyExists", core.AlreadyExistsCode, ClientError},
		{"ResourceExhausted", core.ResourceExhaustedCode, ServerError},
		{"Validation", core.ValidationErrorCode, ClientError},
		{"BusinessRule", core.BusinessRuleViolationCode, ClientError},
		{"Database", core.DatabaseErrorCode, ServerError},
		{"Internal", core.InternalErrorCode, ServerError},
		{"DataCorruption", core.DataCorruptionCode, ServerError},
		{"Configuration", core.ConfigurationErrorCode, ServerError},
		{"Timeout", core.TimeoutCode, SystemError},
		{"Canceled", core.CanceledCode, SystemError},
		{"Concurrency", core.ConcurrencyErrorCode, SystemError},
		{"ExternalService", core.ExternalServiceErrorCode, ExternalError},
		{"Network", core.NetworkErrorCode, ExternalError},
		{"Unauthorized", core.UnauthorizedCode, SecurityError},
		{"Forbidden", core.ForbiddenCode, SecurityError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCat := categorizer.Categorize(tt.code)
			assert.Equal(t, tt.wantCat, gotCat)
		})
	}

	// Test category-specific check functions
	assert.True(t, categorizer.IsClientError(core.NotFoundCode))
	assert.True(t, categorizer.IsServerError(core.DatabaseErrorCode))
	assert.True(t, categorizer.IsSystemError(core.TimeoutCode))
	assert.True(t, categorizer.IsExternalError(core.ExternalServiceErrorCode))
	assert.True(t, categorizer.IsSecurityError(core.UnauthorizedCode))

	assert.False(t, categorizer.IsClientError(core.DatabaseErrorCode))
	assert.False(t, categorizer.IsServerError(core.NotFoundCode))
	assert.False(t, categorizer.IsSystemError(core.UnauthorizedCode))
	assert.False(t, categorizer.IsExternalError(core.TimeoutCode))
	assert.False(t, categorizer.IsSecurityError(core.ExternalServiceErrorCode))
}
