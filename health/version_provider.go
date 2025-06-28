// Copyright (c) 2025 A Bit of Help, Inc.

// Package health provides functionality for health checking applications.
package health

// ConfigVersionAdapter adapts a config to the VersionProvider interface
type ConfigVersionAdapter struct {
	version string
}

// NewConfigVersionAdapter creates a new ConfigVersionAdapter
func NewConfigVersionAdapter(version string) *ConfigVersionAdapter {
	return &ConfigVersionAdapter{
		version: version,
	}
}

// GetVersion returns the application version
func (a *ConfigVersionAdapter) GetVersion() string {
	return a.version
}
