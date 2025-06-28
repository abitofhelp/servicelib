// Copyright (c) 2025 A Bit of Help, Inc.

// Package model provides utilities for working with domain models and DTOs (Data Transfer Objects).
//
// This package offers tools for managing the conversion and manipulation of data models
// across different layers of an application. It focuses on simplifying the process of
// copying data between domain models and DTOs, which is a common requirement in
// applications following Clean Architecture or Hexagonal Architecture patterns.
//
// In layered architectures, different representations of the same data are often needed:
//   - Domain models: Rich objects with behavior used in the domain layer
//   - DTOs: Simple data containers used for API responses or persistence
//   - View models: Data structures tailored for specific UI requirements
//
// Converting between these representations manually can be error-prone and repetitive.
// This package provides reflection-based utilities to automate these conversions while
// maintaining type safety and handling complex nested structures.
//
// Key features:
//   - Field copying between structs based on field names
//   - Deep copying of complex objects with nested structures
//   - Support for various types including pointers, slices, and maps
//   - Type-safe operations with comprehensive error handling
//
// Example usage for copying fields between different struct types:
//
//	// Domain model
//	type User struct {
//	    ID        string
//	    FirstName string
//	    LastName  string
//	    Email     string
//	    Password  string // Sensitive field
//	    CreatedAt time.Time
//	}
//
//	// DTO for API responses
//	type UserDTO struct {
//	    ID        string
//	    FirstName string
//	    LastName  string
//	    Email     string
//	    // Password is omitted for security
//	    CreatedAt time.Time
//	}
//
//	// Copy fields from domain model to DTO
//	user := &User{
//	    ID:        "123",
//	    FirstName: "John",
//	    LastName:  "Doe",
//	    Email:     "john@example.com",
//	    Password:  "secret",
//	    CreatedAt: time.Now(),
//	}
//
//	userDTO := &UserDTO{}
//	err := model.CopyFields(userDTO, user)
//	if err != nil {
//	    log.Fatalf("Failed to copy fields: %v", err)
//	}
//
//	// userDTO now contains all matching fields from user (except Password)
//
// Example usage for creating a deep copy of an object:
//
//	// Create a deep copy of a complex object
//	originalUser := &User{
//	    ID:        "123",
//	    FirstName: "John",
//	    LastName:  "Doe",
//	    Email:     "john@example.com",
//	    Password:  "secret",
//	    CreatedAt: time.Now(),
//	}
//
//	copiedUser := &User{}
//	err := model.DeepCopy(copiedUser, originalUser)
//	if err != nil {
//	    log.Fatalf("Failed to create deep copy: %v", err)
//	}
//
//	// copiedUser is now a complete copy of originalUser
//	// Modifying copiedUser will not affect originalUser
//
// The package is designed to be used in service layers or adapters where
// conversion between different data representations is required.
package model
