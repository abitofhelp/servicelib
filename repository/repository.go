// Copyright (c) 2025 A Bit of Help, Inc.

// Package repository provides generic repository interfaces that can be used across different applications.
package repository

import (
	"context"
)

// Repository is a generic repository interface for entity persistence operations.
// This interface represents a port in the Hexagonal Architecture pattern.
// It's defined in the domain layer but implemented in the infrastructure layer.
//
// The generic type parameter T represents the entity type that the repository manages.
// This allows for type-safe repository operations without code duplication.
type Repository[T any] interface {
	// GetByID retrieves an entity by its ID.
	// It returns the entity if found, or an error if the entity doesn't exist
	// or if there was a problem accessing the data store.
	//
	// The context parameter can be used to control cancellation and timeouts.
	GetByID(ctx context.Context, id string) (T, error)

	// GetAll retrieves all entities of type T.
	// It returns a slice of entities, which may be empty if no entities exist,
	// or an error if there was a problem accessing the data store.
	//
	// The context parameter can be used to control cancellation and timeouts.
	GetAll(ctx context.Context) ([]T, error)

	// Save persists an entity to the data store.
	// For new entities, this typically creates a new record.
	// For existing entities, this updates the existing record.
	//
	// It returns an error if there was a problem saving the entity.
	// The context parameter can be used to control cancellation and timeouts.
	Save(ctx context.Context, entity T) error
}

// RepositoryFactory is an interface for creating repositories.
// This interface follows the Factory pattern and is used to abstract
// the creation of repository instances, allowing for dependency injection
// and easier testing.
type RepositoryFactory interface {
	// GetRepository returns a repository instance.
	// The returned repository should be cast to the appropriate type
	// by the caller, typically Repository[T] for some entity type T.
	//
	// Example:
	//   factory := NewMyRepositoryFactory()
	//   repo, ok := factory.GetRepository().(Repository[User])
	//   if !ok {
	//       return errors.New("invalid repository type")
	//   }
	GetRepository() any
}
