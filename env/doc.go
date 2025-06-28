// Copyright (c) 2025 A Bit of Help, Inc.

// Package env provides utilities for working with environment variables.
//
// This package offers a simple and consistent way to access environment variables
// with fallback values, making it easier to configure applications through the
// environment. Using environment variables for configuration is a best practice
// for cloud-native applications, following the principles of the Twelve-Factor App
// methodology.
//
// The package is designed to be lightweight and focused on a single responsibility:
// retrieving environment variables with sensible defaults. This helps prevent
// application crashes due to missing environment variables and reduces the need
// for conditional checks throughout the codebase.
//
// Key features:
//   - Retrieval of environment variables with fallback values
//   - Simple, consistent API for environment access
//   - Zero external dependencies
//
// Example usage:
//
//	// Get database connection string with a default value
//	dbConnString := env.GetEnv("DATABASE_URL", "postgres://localhost:5432/mydb")
//
//	// Get application port with a default value
//	port := env.GetEnv("PORT", "8080")
//
//	// Get log level with a default value
//	logLevel := env.GetEnv("LOG_LEVEL", "info")
//
// The package is typically used during application startup to load configuration
// values from the environment. It can be used directly or as part of a more
// comprehensive configuration system.
package env
