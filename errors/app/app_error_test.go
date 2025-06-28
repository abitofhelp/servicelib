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

// TestApplicationErrorAsMethod tests the As method of ApplicationError
func TestApplicationErrorAsMethod(t *testing.T) {
	// Create an ApplicationError
	appErr := NewApplicationError(core.ConfigurationErrorCode, "Application error", nil)

	// Test As with ApplicationError target
	var targetAppErr *ApplicationError
	assert.True(t, errors.As(appErr, &targetAppErr), "errors.As should return true for ApplicationError target")
	assert.Equal(t, appErr, targetAppErr, "Target should be set to the original error")

	// Test As with BaseError target
	var targetBaseErr *core.BaseError
	assert.True(t, errors.As(appErr, &targetBaseErr), "errors.As should return true for BaseError target")
	assert.Equal(t, appErr.BaseError, targetBaseErr, "Target should be set to the BaseError")

	// Test As with wrong target type
	var wrongTarget *ConfigurationError
	assert.False(t, errors.As(appErr, &wrongTarget), "errors.As should return false for wrong target type")

	// Test direct call to As method with ApplicationError target
	var directTargetAppErr *ApplicationError
	assert.True(t, appErr.As(&directTargetAppErr), "As should return true for ApplicationError target")
	assert.Equal(t, appErr, directTargetAppErr, "Target should be set to the original error")

	// Test direct call to As method with BaseError target
	var directTargetBaseErr *core.BaseError
	assert.True(t, appErr.As(&directTargetBaseErr), "As should return true for BaseError target")
	assert.Equal(t, appErr.BaseError, directTargetBaseErr, "Target should be set to the BaseError")

	// Test direct call to As method with wrong target type
	var directWrongTarget *ConfigurationError
	assert.False(t, appErr.As(&directWrongTarget), "As should return false for wrong target type")

	// Test with nil target
	assert.False(t, appErr.As(nil), "As should return false for nil target")
}

// TestConfigurationErrorAsMethod tests the As method of ConfigurationError
func TestConfigurationErrorAsMethod(t *testing.T) {
	// Create a ConfigurationError
	configErr := NewConfigurationError("Invalid configuration", "MAX_CONNECTIONS", "abc", nil)

	// Test As with ConfigurationError target
	var targetConfigErr *ConfigurationError
	assert.True(t, errors.As(configErr, &targetConfigErr), "errors.As should return true for ConfigurationError target")
	assert.Equal(t, configErr, targetConfigErr, "Target should be set to the original error")

	// Test As with ApplicationError target
	var targetAppErr *ApplicationError
	assert.True(t, errors.As(configErr, &targetAppErr), "errors.As should return true for ApplicationError target")
	assert.Equal(t, configErr.ApplicationError, targetAppErr, "Target should be set to the ApplicationError")

	// Test As with BaseError target
	var targetBaseErr *core.BaseError
	assert.True(t, errors.As(configErr, &targetBaseErr), "errors.As should return true for BaseError target")
	assert.Equal(t, configErr.ApplicationError.BaseError, targetBaseErr, "Target should be set to the BaseError")

	// Test As with wrong target type
	var wrongTarget *AuthenticationError
	assert.False(t, errors.As(configErr, &wrongTarget), "errors.As should return false for wrong target type")

	// Test direct call to As method with ConfigurationError target
	var directTargetConfigErr *ConfigurationError
	assert.True(t, configErr.As(&directTargetConfigErr), "As should return true for ConfigurationError target")
	assert.Equal(t, configErr, directTargetConfigErr, "Target should be set to the original error")

	// Test direct call to As method with ApplicationError target
	var directTargetAppErr *ApplicationError
	assert.True(t, configErr.As(&directTargetAppErr), "As should return true for ApplicationError target")
	assert.Equal(t, configErr.ApplicationError, directTargetAppErr, "Target should be set to the ApplicationError")

	// Test direct call to As method with BaseError target
	var directTargetBaseErr *core.BaseError
	assert.True(t, configErr.As(&directTargetBaseErr), "As should return true for BaseError target")
	assert.Equal(t, configErr.ApplicationError.BaseError, directTargetBaseErr, "Target should be set to the BaseError")

	// Test direct call to As method with wrong target type
	var directWrongTarget *AuthenticationError
	assert.False(t, configErr.As(&directWrongTarget), "As should return false for wrong target type")

	// Test with nil target
	assert.False(t, configErr.As(nil), "As should return false for nil target")
}

// TestAuthenticationErrorAsMethod tests the As method of AuthenticationError
func TestAuthenticationErrorAsMethod(t *testing.T) {
	// Create an AuthenticationError
	authErr := NewAuthenticationError("Invalid credentials", "john.doe", nil)

	// Test As with AuthenticationError target
	var targetAuthErr *AuthenticationError
	assert.True(t, errors.As(authErr, &targetAuthErr), "errors.As should return true for AuthenticationError target")
	assert.Equal(t, authErr, targetAuthErr, "Target should be set to the original error")

	// Test As with ApplicationError target
	var targetAppErr *ApplicationError
	assert.True(t, errors.As(authErr, &targetAppErr), "errors.As should return true for ApplicationError target")
	assert.Equal(t, authErr.ApplicationError, targetAppErr, "Target should be set to the ApplicationError")

	// Test As with BaseError target
	var targetBaseErr *core.BaseError
	assert.True(t, errors.As(authErr, &targetBaseErr), "errors.As should return true for BaseError target")
	assert.Equal(t, authErr.ApplicationError.BaseError, targetBaseErr, "Target should be set to the BaseError")

	// Test As with wrong target type
	var wrongTarget *ConfigurationError
	assert.False(t, errors.As(authErr, &wrongTarget), "errors.As should return false for wrong target type")

	// Test direct call to As method with AuthenticationError target
	var directTargetAuthErr *AuthenticationError
	assert.True(t, authErr.As(&directTargetAuthErr), "As should return true for AuthenticationError target")
	assert.Equal(t, authErr, directTargetAuthErr, "Target should be set to the original error")

	// Test direct call to As method with ApplicationError target
	var directTargetAppErr *ApplicationError
	assert.True(t, authErr.As(&directTargetAppErr), "As should return true for ApplicationError target")
	assert.Equal(t, authErr.ApplicationError, directTargetAppErr, "Target should be set to the ApplicationError")

	// Test direct call to As method with BaseError target
	var directTargetBaseErr *core.BaseError
	assert.True(t, authErr.As(&directTargetBaseErr), "As should return true for BaseError target")
	assert.Equal(t, authErr.ApplicationError.BaseError, directTargetBaseErr, "Target should be set to the BaseError")

	// Test direct call to As method with wrong target type
	var directWrongTarget *ConfigurationError
	assert.False(t, authErr.As(&directWrongTarget), "As should return false for wrong target type")

	// Test with nil target
	assert.False(t, authErr.As(nil), "As should return false for nil target")
}

// TestAuthorizationErrorAsMethod tests the As method of AuthorizationError
func TestAuthorizationErrorAsMethod(t *testing.T) {
	// Create an AuthorizationError
	authzErr := NewAuthorizationError("Access denied", "john.doe", "users", "delete", nil)

	// Test As with AuthorizationError target
	var targetAuthzErr *AuthorizationError
	assert.True(t, errors.As(authzErr, &targetAuthzErr), "errors.As should return true for AuthorizationError target")
	assert.Equal(t, authzErr, targetAuthzErr, "Target should be set to the original error")

	// Test As with ApplicationError target
	var targetAppErr *ApplicationError
	assert.True(t, errors.As(authzErr, &targetAppErr), "errors.As should return true for ApplicationError target")
	assert.Equal(t, authzErr.ApplicationError, targetAppErr, "Target should be set to the ApplicationError")

	// Test As with BaseError target
	var targetBaseErr *core.BaseError
	assert.True(t, errors.As(authzErr, &targetBaseErr), "errors.As should return true for BaseError target")
	assert.Equal(t, authzErr.ApplicationError.BaseError, targetBaseErr, "Target should be set to the BaseError")

	// Test As with wrong target type
	var wrongTarget *ConfigurationError
	assert.False(t, errors.As(authzErr, &wrongTarget), "errors.As should return false for wrong target type")

	// Test direct call to As method with AuthorizationError target
	var directTargetAuthzErr *AuthorizationError
	assert.True(t, authzErr.As(&directTargetAuthzErr), "As should return true for AuthorizationError target")
	assert.Equal(t, authzErr, directTargetAuthzErr, "Target should be set to the original error")

	// Test direct call to As method with ApplicationError target
	var directTargetAppErr *ApplicationError
	assert.True(t, authzErr.As(&directTargetAppErr), "As should return true for ApplicationError target")
	assert.Equal(t, authzErr.ApplicationError, directTargetAppErr, "Target should be set to the ApplicationError")

	// Test direct call to As method with BaseError target
	var directTargetBaseErr *core.BaseError
	assert.True(t, authzErr.As(&directTargetBaseErr), "As should return true for BaseError target")
	assert.Equal(t, authzErr.ApplicationError.BaseError, directTargetBaseErr, "Target should be set to the BaseError")

	// Test direct call to As method with wrong target type
	var directWrongTarget *ConfigurationError
	assert.False(t, authzErr.As(&directWrongTarget), "As should return false for wrong target type")

	// Test with nil target
	assert.False(t, authzErr.As(nil), "As should return false for nil target")
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
