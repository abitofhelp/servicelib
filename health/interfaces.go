// Copyright (c) 2025 A Bit of Help, Inc.

// Package health provides functionality for health checking applications.
package health

// RepositoryFactoryProvider defines the interface for getting a repository factory
type RepositoryFactoryProvider interface {
	GetRepositoryFactory() any // Using 'any' as a more modern alias for interface{}
}

// HealthCheckProvider combines the interfaces needed for health checking
type HealthCheckProvider interface {
	RepositoryFactoryProvider
}

// VersionProvider defines the interface for getting the application version
type VersionProvider interface {
	GetVersion() string
}

// HealthConfig defines the interface for health check configuration
type HealthConfig interface {
	// GetVersion returns the application version
	GetVersion() string

	// GetName returns the application name
	GetName() string

	// GetEnvironment returns the application environment
	GetEnvironment() string

	// GetTimeout returns the timeout for health checks
	GetTimeout() int
}
