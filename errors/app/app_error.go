// Copyright (c) 2025 A Bit of Help, Inc.

// Package app provides application-level error types for the application.
package app

import (
	"github.com/abitofhelp/servicelib/errors/core"
)

// ApplicationError represents an application-level error.
// It extends BaseError with application-specific information.
type ApplicationError struct {
	*core.BaseError
}

// NewApplicationError creates a new ApplicationError.
func NewApplicationError(code core.ErrorCode, message string, cause error) *ApplicationError {
	return &ApplicationError{
		BaseError: core.NewBaseError(code, message, cause),
	}
}

// IsApplicationError identifies this as an application error.
func (e *ApplicationError) IsApplicationError() bool {
	return true
}

// Unwrap returns the underlying error.
func (e *ApplicationError) Unwrap() error {
	return e.BaseError
}

// As implements the errors.As interface for ApplicationError.
func (e *ApplicationError) As(target interface{}) bool {
	// Check if target is nil
	if target == nil {
		return false
	}

	// Check if target is *ApplicationError
	if t, ok := target.(**ApplicationError); ok {
		*t = e
		return true
	}

	// Check if target is *BaseError
	if t, ok := target.(**core.BaseError); ok {
		*t = e.BaseError
		return true
	}

	// Delegate to BaseError.As for other types
	return e.BaseError.As(target)
}

// ConfigurationError represents a configuration error.
// It extends ApplicationError with configuration-specific information.
type ConfigurationError struct {
	*ApplicationError
	ConfigKey   string `json:"config_key,omitempty"`
	ConfigValue string `json:"config_value,omitempty"`
}

// NewConfigurationError creates a new ConfigurationError.
func NewConfigurationError(message string, configKey string, configValue string, cause error) *ConfigurationError {
	return &ConfigurationError{
		ApplicationError: NewApplicationError(core.ConfigurationErrorCode, message, cause),
		ConfigKey:        configKey,
		ConfigValue:      configValue,
	}
}

// IsConfigurationError identifies this as a configuration error.
func (e *ConfigurationError) IsConfigurationError() bool {
	return true
}

// Unwrap returns the underlying error.
func (e *ConfigurationError) Unwrap() error {
	return e.ApplicationError
}

// As implements the errors.As interface for ConfigurationError.
func (e *ConfigurationError) As(target interface{}) bool {
	// Check if target is nil
	if target == nil {
		return false
	}

	// Check if target is *ConfigurationError
	if t, ok := target.(**ConfigurationError); ok {
		*t = e
		return true
	}

	// Check if target is *ApplicationError
	if t, ok := target.(**ApplicationError); ok {
		*t = e.ApplicationError
		return true
	}

	// Check if target is *BaseError
	if t, ok := target.(**core.BaseError); ok {
		*t = e.ApplicationError.BaseError
		return true
	}

	// Delegate to ApplicationError.As for other types
	return e.ApplicationError.As(target)
}

// AuthenticationError represents an authentication failure.
// It extends ApplicationError with authentication-specific information.
type AuthenticationError struct {
	*ApplicationError
	Username string `json:"username,omitempty"`
}

// NewAuthenticationError creates a new AuthenticationError.
func NewAuthenticationError(message string, username string, cause error) *AuthenticationError {
	return &AuthenticationError{
		ApplicationError: NewApplicationError(core.UnauthorizedCode, message, cause),
		Username:         username,
	}
}

// IsAuthenticationError identifies this as an authentication error.
func (e *AuthenticationError) IsAuthenticationError() bool {
	return true
}

// Unwrap returns the underlying error.
func (e *AuthenticationError) Unwrap() error {
	return e.ApplicationError
}

// As implements the errors.As interface for AuthenticationError.
func (e *AuthenticationError) As(target interface{}) bool {
	// Check if target is nil
	if target == nil {
		return false
	}

	// Check if target is *AuthenticationError
	if t, ok := target.(**AuthenticationError); ok {
		*t = e
		return true
	}

	// Check if target is *ApplicationError
	if t, ok := target.(**ApplicationError); ok {
		*t = e.ApplicationError
		return true
	}

	// Check if target is *BaseError
	if t, ok := target.(**core.BaseError); ok {
		*t = e.ApplicationError.BaseError
		return true
	}

	// Delegate to ApplicationError.As for other types
	return e.ApplicationError.As(target)
}

// AuthorizationError represents an authorization failure.
// It extends ApplicationError with authorization-specific information.
type AuthorizationError struct {
	*ApplicationError
	Username string `json:"username,omitempty"`
	Resource string `json:"resource,omitempty"`
	Action   string `json:"action,omitempty"`
}

// NewAuthorizationError creates a new AuthorizationError.
func NewAuthorizationError(message string, username string, resource string, action string, cause error) *AuthorizationError {
	return &AuthorizationError{
		ApplicationError: NewApplicationError(core.ForbiddenCode, message, cause),
		Username:         username,
		Resource:         resource,
		Action:           action,
	}
}

// IsAuthorizationError identifies this as an authorization error.
func (e *AuthorizationError) IsAuthorizationError() bool {
	return true
}

// Unwrap returns the underlying error.
func (e *AuthorizationError) Unwrap() error {
	return e.ApplicationError
}

// As implements the errors.As interface for AuthorizationError.
func (e *AuthorizationError) As(target interface{}) bool {
	// Check if target is nil
	if target == nil {
		return false
	}

	// Check if target is *AuthorizationError
	if t, ok := target.(**AuthorizationError); ok {
		*t = e
		return true
	}

	// Check if target is *ApplicationError
	if t, ok := target.(**ApplicationError); ok {
		*t = e.ApplicationError
		return true
	}

	// Check if target is *BaseError
	if t, ok := target.(**core.BaseError); ok {
		*t = e.ApplicationError.BaseError
		return true
	}

	// Delegate to ApplicationError.As for other types
	return e.ApplicationError.As(target)
}
