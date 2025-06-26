// Copyright (c) 2025 A Bit of Help, Inc.

// Package repository provides generic repository interfaces for data persistence operations.
//
// This package implements the Repository pattern, which abstracts the data access layer
// and provides a consistent interface for working with different data sources.
// It aligns with the Hexagonal Architecture (Ports and Adapters) pattern by defining
// ports (interfaces) that can be implemented by various adapters.
//
// The package is designed to be used with any entity type through Go generics,
// allowing for type-safe repository operations without code duplication.
//
// Key components:
//   - Repository: A generic interface for CRUD operations on entities
//   - RepositoryFactory: An interface for creating repositories
//
// The Repository pattern provides several benefits:
//   - Decouples business logic from data access implementation
//   - Simplifies testing through mocking
//   - Enables switching between different data sources with minimal code changes
//   - Centralizes data access logic
//
// Example usage:
//
//	// Define an entity
//	type User struct {
//	    ID   string
//	    Name string
//	    Age  int
//	}
//
//	// Use the repository interface
//	func ProcessUser(ctx context.Context, repo repository.Repository[User], userID string) error {
//	    // Get user from repository
//	    user, err := repo.GetByID(ctx, userID)
//	    if err != nil {
//	        return err
//	    }
//	
//	    // Process user...
//	    user.Age++
//	
//	    // Save changes
//	    return repo.Save(ctx, user)
//	}
//
//	// Implementation would be provided by an adapter in the infrastructure layer
//	type MongoUserRepository struct {
//	    // MongoDB client and collection
//	}
//
//	func (r *MongoUserRepository) GetByID(ctx context.Context, id string) (User, error) {
//	    // Implementation using MongoDB
//	}
//
//	func (r *MongoUserRepository) GetAll(ctx context.Context) ([]User, error) {
//	    // Implementation using MongoDB
//	}
//
//	func (r *MongoUserRepository) Save(ctx context.Context, user User) error {
//	    // Implementation using MongoDB
//	}
package repository