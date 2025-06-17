// Copyright (c) 2025 A Bit of Help, Inc.

// Package config provides generic configuration interfaces that can be used across different applications.
package config

import (
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
)

// ConfigInterface is an interface for the config.Config type
type ConfigInterface interface {
	// GetApp returns the App configuration
	GetApp() *config.AppConfig
	
	// GetDatabase returns the Database configuration
	GetDatabase() *config.DatabaseConfig
}