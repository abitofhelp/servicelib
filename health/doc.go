// Copyright (c) 2025 A Bit of Help, Inc.

// Package health provides functionality for health checking applications.
//
// This package offers a comprehensive solution for implementing health checks
// in web applications and services. It includes HTTP handlers for health check
// endpoints, interfaces for health status reporting, and utilities for checking
// the health of various dependencies like databases and external services.
//
// Health checks are essential for monitoring the status of applications in
// production environments, enabling automated systems to detect and respond
// to issues, and providing visibility into the application's operational state.
//
// Key features:
//   - HTTP handlers for health check endpoints
//   - Standardized health status representation
//   - Service dependency health checking
//   - Configurable health check timeouts
//   - Version information inclusion in health responses
//   - Support for different health status levels (healthy, degraded)
//
// The package defines several interfaces that applications can implement to
// integrate with the health checking system:
//   - HealthCheckProvider: Combines interfaces needed for health checking
//   - RepositoryFactoryProvider: For checking database connectivity
//   - VersionProvider: For including version information in health responses
//   - HealthConfig: For configuring health check behavior
//
// Example usage:
//
//	// Create a health check handler
//	healthHandler := health.NewHandler(
//	    app,           // Implements HealthCheckProvider
//	    logger,        // For logging health check results
//	    appConfig,     // For configuration
//	)
//
//	// Register the health check endpoint
//	http.HandleFunc("/health", healthHandler)
//
// The health check endpoint returns a JSON response with the following structure:
//
//	{
//	    "status": "Healthy",
//	    "timestamp": "2023-06-25T12:34:56Z",
//	    "version": "1.0.0",
//	    "services": {
//	        "database": "Up",
//	        "cache": "Up",
//	        "external_api": "Down"
//	    }
//	}
//
// The status field can be "Healthy" or "Degraded", depending on the health of
// the application and its dependencies. The services field contains the status
// of each dependency, which can be "Up" or "Down".
//
// The package also provides adapters for common dependencies, making it easy
// to integrate with existing code:
//   - ConfigAdapter: Adapts a configuration object to the HealthConfig interface
//   - RepositoryAdapter: Adapts a repository factory to the RepositoryFactoryProvider interface
//
// Health checks are designed to be lightweight and fast, with configurable
// timeouts to prevent them from impacting application performance.
package health
