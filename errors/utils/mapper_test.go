// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abitofhelp/servicelib/errors/core"
)

func TestErrorCodeMapper_ToHTTPStatus(t *testing.T) {
	tests := []struct {
		name   string
		code   core.ErrorCode
		status int
	}{
		{
			name:   "not found",
			code:   core.NotFoundCode,
			status: 404,
		},
		{
			name:   "invalid input",
			code:   core.InvalidInputCode,
			status: 400,
		},
		{
			name:   "internal error",
			code:   core.InternalErrorCode,
			status: 500,
		},
		{
			name:   "timeout",
			code:   core.TimeoutCode,
			status: 504,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.status, ErrorCodeMapper.ToHTTPStatus(tt.code))
		})
	}
}

func TestErrorCodeMapper_ToCategory(t *testing.T) {
	tests := []struct {
		name     string
		code     core.ErrorCode
		category string
	}{
		{
			name:     "client error",
			code:     core.NotFoundCode,
			category: "Client",
		},
		{
			name:     "server error",
			code:     core.InternalErrorCode,
			category: "Server",
		},
		{
			name:     "system error",
			code:     core.TimeoutCode,
			category: "System",
		},
		{
			name:     "security error",
			code:     core.UnauthorizedCode,
			category: "Security",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.category, ErrorCodeMapper.ToCategory(tt.code))
		})
	}
}

func TestErrorCodeMapper_ToMessage(t *testing.T) {
	tests := []struct {
		name   string
		code   core.ErrorCode
		msg    string
	}{
		{
			name: "not found",
			code: core.NotFoundCode,
			msg:  "The requested resource was not found",
		},
		{
			name: "invalid input",
			code: core.InvalidInputCode,
			msg:  "The provided input is invalid",
		},
		{
			name: "internal error",
			code: core.InternalErrorCode,
			msg:  "An internal server error occurred",
		},
		{
			name: "timeout",
			code: core.TimeoutCode,
			msg:  "The operation timed out",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.msg, ErrorCodeMapper.ToMessage(tt.code))
		})
	}
}
