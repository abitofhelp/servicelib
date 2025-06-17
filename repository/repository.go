// Copyright (c) 2025 A Bit of Help, Inc.

// Package repository provides generic repository interfaces that can be used across different applications.
package repository

import (
	"context"
)

// Repository is a generic repository interface for entity persistence operations.
// This interface represents a port in the Hexagonal Architecture pattern.
// It's defined in the pkg layer but implemented in the infrastructure layer.
type Repository[T any] interface {
	// GetByID retrieves an entity by its ID
	GetByID(ctx context.Context, id string) (T, error)

	// GetAll retrieves all entities
	GetAll(ctx context.Context) ([]T, error)

	// Save persists an entity
	Save(ctx context.Context, entity T) error
}

// RepositoryFactory is an interface for creating repositories
type RepositoryFactory interface {
	// GetRepository returns a repository for the given entity type
	GetRepository() any
}