// Copyright (c) 2025 A Bit of Help, Inc.

// Package health provides functionality for health checking applications.
package health

import (
	"github.com/abitofhelp/servicelib/config"
)

// ConfigAdapter adapts a config.Config to the HealthConfig interface
type ConfigAdapter struct {
	config config.Config
}

// NewConfigAdapter creates a new ConfigAdapter
func NewConfigAdapter(cfg config.Config) *ConfigAdapter {
	return &ConfigAdapter{
		config: cfg,
	}
}

// GetVersion returns the application version
func (a *ConfigAdapter) GetVersion() string {
	return a.config.GetApp().GetVersion()
}

// GetName returns the application name
func (a *ConfigAdapter) GetName() string {
	return a.config.GetApp().GetName()
}

// GetEnvironment returns the application environment
func (a *ConfigAdapter) GetEnvironment() string {
	return a.config.GetApp().GetEnvironment()
}

// GetTimeout returns the timeout for health checks
func (a *ConfigAdapter) GetTimeout() int {
	// Default to 5 seconds if not specified
	return 5
}
