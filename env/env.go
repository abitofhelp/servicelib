// Copyright (c) 2025 A Bit of Help, Inc.

// Package env provides utilities for working with environment variables.
package env

import "os"

// GetEnv retrieves the value of the environment variable named by the key.
// If the variable is not present, the fallback value is returned.
// Parameters:
//   - key: The name of the environment variable
//   - fallback: The value to return if the environment variable is not set
//
// Returns:
//   - string: The value of the environment variable or the fallback value
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
