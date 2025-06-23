// Copyright (c) 2025 A Bit of Help, Inc.

// Package config provides generic configuration interfaces that can be used across different applications.
package config

import (
	"github.com/abitofhelp/servicelib/config/interfaces"
)

// AppConfig is an alias for interfaces.AppConfig for backward compatibility
type AppConfig = interfaces.AppConfig

// DatabaseConfig is an alias for interfaces.DatabaseConfig for backward compatibility
type DatabaseConfig = interfaces.DatabaseConfig

// Config is an alias for interfaces.Config for backward compatibility
type Config = interfaces.Config
