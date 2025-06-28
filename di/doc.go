// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides a comprehensive dependency injection framework for Go applications.
//
// This package implements a flexible dependency injection system that helps manage
// application dependencies and promotes loose coupling between components. It supports
// various application architectures, with special focus on Domain-Driven Design (DDD)
// and Clean Architecture patterns.
//
// Dependency injection is a design pattern that allows the separation of object creation
// from object usage, making code more modular, testable, and maintainable. This package
// provides containers that manage the lifecycle of dependencies and their relationships.
//
// Key features:
//   - Generic containers that work with any configuration type
//   - Support for layered architecture (repositories, domain services, application services)
//   - Type-safe dependency management using Go generics
//   - Integration with logging, validation, and context management
//   - Resource lifecycle management with proper cleanup
//   - Thread-safe operation for concurrent applications
//
// The package provides several container types for different use cases:
//   - BaseContainer: A foundational container with basic dependencies
//   - Container: A backward-compatible container using interface{} types
//   - GenericAppContainer: A type-safe container for layered applications
//   - RepositoryContainer: A specialized container for repository management
//   - ServiceContainer: A container focused on service management
//
// Example usage for a simple application:
//
//	// Create a container with application configuration
//	container, err := di.NewContainer(ctx, logger, appConfig)
//	if err != nil {
//	    log.Fatalf("Failed to create container: %v", err)
//	}
//	defer container.Close()
//
//	// Use the container to access dependencies
//	validator := container.GetValidator()
//	logger := container.GetLogger()
//
// Example usage for a layered application:
//
//	// Create a generic container with typed dependencies
//	container, err := di.NewGenericAppContainer(
//	    ctx,
//	    logger,
//	    appConfig,
//	    connectionString,
//	    initUserRepository,
//	    initUserDomainService,
//	    initUserApplicationService,
//	)
//	if err != nil {
//	    log.Fatalf("Failed to create container: %v", err)
//	}
//	defer container.Close()
//
//	// Access typed dependencies
//	userRepo := container.GetRepository()
//	userService := container.GetApplicationService()
//
// The package is designed to be extended for specific application needs.
// Custom containers can be created by embedding the provided containers
// and adding application-specific dependencies and initialization logic.
package di
