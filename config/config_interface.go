// Copyright (c) 2025 A Bit of Help, Inc.

// Package config provides generic configuration interfaces that can be used across different applications.
package config

// Note: This file is deprecated. Use the Config interface from config.go instead.
// It's kept for backward compatibility.

// ConfigInterface is an interface for configuration
// Deprecated: Use Config interface instead
type ConfigInterface interface {
	// GetApp returns the App configuration
	GetApp() interface{}

	// GetDatabase returns the Database configuration
	GetDatabase() interface{}
}
