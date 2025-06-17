// Copyright (c) 2025 A Bit of Help, Inc.

// Package health provides functionality for health checking applications.
package health

// Status constants for health check responses
const (
	// StatusHealthy represents a healthy system status
	StatusHealthy = "Healthy"

	// StatusDegraded represents a degraded system status
	StatusDegraded = "Degraded"
)

// Service status constants for health check responses
const (
	// ServiceUp represents a service that is up and running
	ServiceUp = "Up"

	// ServiceDown represents a service that is down or unavailable
	ServiceDown = "Down"
)