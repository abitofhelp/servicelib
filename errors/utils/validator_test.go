// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
)

func TestErrorCodeValidator(t *testing.T) {
	v := NewErrorCodeValidator()
	tests := []struct {
		name    string
		code    core.ErrorCode
		wantErr bool
	}{
		{"Valid code", core.NotFoundCode, false},
		{"Invalid code", core.ErrorCode("INVALID"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(tt.code); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// Test code name mapping
	if got := v.GetCodeName(core.NotFoundCode); got != "NotFound" {
		t.Errorf("GetCodeName(NotFoundCode) = %v, want NotFound", got)
	}
	if got := v.GetCodeName(core.InvalidInputCode); got != "InvalidInput" {
		t.Errorf("GetCodeName(InvalidInputCode) = %v, want InvalidInput", got)
	}
	if got := v.GetCodeName(core.AlreadyExistsCode); got != "AlreadyExists" {
		t.Errorf("GetCodeName(AlreadyExistsCode) = %v, want AlreadyExists", got)
	}
	if got := v.GetCodeName(core.ResourceExhaustedCode); got != "ResourceExhausted" {
		t.Errorf("GetCodeName(ResourceExhaustedCode) = %v, want ResourceExhausted", got)
	}
	if got := v.GetCodeName(core.ValidationErrorCode); got != "Validation" {
		t.Errorf("GetCodeName(ValidationErrorCode) = %v, want Validation", got)
	}
	if got := v.GetCodeName(core.BusinessRuleViolationCode); got != "BusinessRuleViolation" {
		t.Errorf("GetCodeName(BusinessRuleViolationCode) = %v, want BusinessRuleViolation", got)
	}
	if got := v.GetCodeName(core.DatabaseErrorCode); got != "Database" {
		t.Errorf("GetCodeName(DatabaseErrorCode) = %v, want Database", got)
	}
	if got := v.GetCodeName(core.InternalErrorCode); got != "Internal" {
		t.Errorf("GetCodeName(InternalErrorCode) = %v, want Internal", got)
	}
	if got := v.GetCodeName(core.TimeoutCode); got != "Timeout" {
		t.Errorf("GetCodeName(TimeoutCode) = %v, want Timeout", got)
	}
	if got := v.GetCodeName(core.CanceledCode); got != "Canceled" {
		t.Errorf("GetCodeName(CanceledCode) = %v, want Canceled", got)
	}
	if got := v.GetCodeName(core.ConcurrencyErrorCode); got != "Concurrency" {
		t.Errorf("GetCodeName(ConcurrencyErrorCode) = %v, want Concurrency", got)
	}
	if got := v.GetCodeName(core.NetworkErrorCode); got != "Network" {
		t.Errorf("GetCodeName(NetworkErrorCode) = %v, want Network", got)
	}
	if got := v.GetCodeName(core.ErrorCode("UNKNOWN")); got != "Unknown" {
		t.Errorf("GetCodeName(UNKNOWN) = %v, want Unknown", got)
	}
	if got := v.GetCodeName(core.ExternalServiceErrorCode); got != "ExternalService" {
		t.Errorf("GetCodeName(ExternalServiceErrorCode) = %v, want ExternalService", got)
	}
	if got := v.GetCodeName(core.NetworkErrorCode); got != "Network" {
		t.Errorf("GetCodeName(NetworkErrorCode) = %v, want Network", got)
	}
	if got := v.GetCodeName(core.UnauthorizedCode); got != "Unauthorized" {
		t.Errorf("GetCodeName(UnauthorizedCode) = %v, want Unauthorized", got)
	}
	if got := v.GetCodeName(core.ForbiddenCode); got != "Forbidden" {
		t.Errorf("GetCodeName(ForbiddenCode) = %v, want Forbidden", got)
	}
	if got := v.GetCodeName(core.ConfigurationErrorCode); got != "Configuration" {
		t.Errorf("GetCodeName(ConfigurationErrorCode) = %v, want Configuration", got)
	}
	if got := v.GetCodeName(core.DataCorruptionCode); got != "DataCorruption" {
		t.Errorf("GetCodeName(DataCorruptionCode) = %v, want DataCorruption", got)
	}
	// Test IsValid
	if !v.IsValid(core.NotFoundCode) {
		t.Errorf("IsValid(NotFoundCode) = false, want true")
	}
	if !v.IsValid(core.InvalidInputCode) {
		t.Errorf("IsValid(InvalidInputCode) = false, want true")
	}
	if v.IsValid("INVALID_CODE") {
		t.Errorf("IsValid(INVALID_CODE) = true, want false")
	}
}
