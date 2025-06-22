// Copyright (c) 2025 A Bit of Help, Inc.

package app

import (
	"errors"
	"testing"

	"github.com/abitofhelp/servicelib/errors/core"
	"github.com/stretchr/testify/assert"
)

func TestNewApplicationError(t *testing.T) {
	// Test creating a new ApplicationError
	err := NewApplicationError(core.ConfigurationErrorCode, "Configuration error", nil)

	// Check that the error is not nil
	assert.NotNil(t, err)

	// Check that the error code is set correctly
	assert.Equal(t, core.ConfigurationErrorCode, err.BaseError.GetCode())

	// Check that the message is set correctly
	assert.Equal(t, "Configuration error", err.BaseError.GetMessage())

	// Check that IsApplicationError returns true
	assert.True(t, err.IsApplicationError())
}

func TestNewConfigurationError(t *testing.T) {
	// Test creating a new ConfigurationError
	err := NewConfigurationError("Invalid configuration value", "MAX_CONNECTIONS", "abc", nil)

	// Check that the error is not nil
	assert.NotNil(t, err)

	// Check that the error code is set correctly
	assert.Equal(t, core.ConfigurationErrorCode, err.BaseError.GetCode())

	// Check that the message is set correctly
	assert.Equal(t, "Invalid configuration value", err.BaseError.GetMessage())

	// Check that the config key and value are set correctly
	assert.Equal(t, "MAX_CONNECTIONS", err.ConfigKey)
	assert.Equal(t, "abc", err.ConfigValue)

	// Check that IsConfigurationError returns true
	assert.True(t, err.IsConfigurationError())

	// Check that IsApplicationError returns true (inheritance)
	assert.True(t, err.IsApplicationError())
}

func TestNewAuthenticationError(t *testing.T) {
	// Test creating a new AuthenticationError
	err := NewAuthenticationError("Invalid credentials", "john.doe", nil)

	// Check that the error is not nil
	assert.NotNil(t, err)

	// Check that the error code is set correctly
	assert.Equal(t, core.UnauthorizedCode, err.BaseError.GetCode())

	// Check that the message is set correctly
	assert.Equal(t, "Invalid credentials", err.BaseError.GetMessage())

	// Check that the username is set correctly
	assert.Equal(t, "john.doe", err.Username)

	// Check that IsAuthenticationError returns true
	assert.True(t, err.IsAuthenticationError())

	// Check that IsApplicationError returns true (inheritance)
	assert.True(t, err.IsApplicationError())
}

func TestNewAuthorizationError(t *testing.T) {
	// Test creating a new AuthorizationError
	err := NewAuthorizationError("Access denied", "john.doe", "users", "delete", nil)

	// Check that the error is not nil
	assert.NotNil(t, err)

	// Check that the error code is set correctly
	assert.Equal(t, core.ForbiddenCode, err.BaseError.GetCode())

	// Check that the message is set correctly
	assert.Equal(t, "Access denied", err.BaseError.GetMessage())

	// Check that the username, resource, and action are set correctly
	assert.Equal(t, "john.doe", err.Username)
	assert.Equal(t, "users", err.Resource)
	assert.Equal(t, "delete", err.Action)

	// Check that IsAuthorizationError returns true
	assert.True(t, err.IsAuthorizationError())

	// Check that IsApplicationError returns true (inheritance)
	assert.True(t, err.IsApplicationError())
}

func TestAppErrorsAs(t *testing.T) {
	// Create errors of different types
	appErr := NewApplicationError(core.ConfigurationErrorCode, "Application error", nil)
	configErr := NewConfigurationError("Invalid configuration value", "MAX_CONNECTIONS", "abc", nil)
	authErr := NewAuthenticationError("Invalid credentials", "john.doe", nil)
	authzErr := NewAuthorizationError("Access denied", "john.doe", "users", "delete", nil)

	// Test errors.As with ApplicationError
	var ae *ApplicationError
	assert.True(t, errors.As(appErr, &ae))

	// Test errors.As with ConfigurationError
	var ce *ConfigurationError
	assert.True(t, errors.As(configErr, &ce))

	// Test errors.As with AuthenticationError
	var authe *AuthenticationError
	assert.True(t, errors.As(authErr, &authe))

	// Test errors.As with AuthorizationError
	var authze *AuthorizationError
	assert.True(t, errors.As(authzErr, &authze))

	// Test inheritance using IsApplicationError
	// ConfigurationError should be recognized as an ApplicationError
	assert.True(t, configErr.IsApplicationError())

	// AuthenticationError should be recognized as an ApplicationError
	assert.True(t, authErr.IsApplicationError())

	// AuthorizationError should be recognized as an ApplicationError
	assert.True(t, authzErr.IsApplicationError())
}

func TestAppErrorWithCause(t *testing.T) {
	// Create a cause error
	cause := errors.New("original error")

	// Create an ApplicationError with the cause
	err := NewApplicationError(core.ConfigurationErrorCode, "Configuration error", cause)

	// Check that the cause is set correctly
	assert.Equal(t, cause, err.BaseError.GetCause())

	// Check that the error message includes the cause
	assert.Contains(t, err.BaseError.Error(), "original error")
}
