// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides generic repository interfaces that can be used across different applications.
package di

// Repository is a generic interface for repositories
type Repository interface {
	// GetID returns the repository ID
	GetID() string
}

// DomainService is a generic interface for domain services
type DomainService interface {
	// GetID returns the domain service ID
	GetID() string
}

// ApplicationService is a generic interface for application services
type ApplicationService interface {
	// GetID returns the application service ID
	GetID() string
}
