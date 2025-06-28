// Copyright (c) 2025 A Bit of Help, Inc.

// Package health provides functionality for health checking the application.
package health

// GenericContainerAdapter is a generic adapter for any container that implements RepositoryFactoryProvider
type GenericContainerAdapter[T RepositoryFactoryProvider] struct {
	Container T
}

// NewGenericContainerAdapter creates a new GenericContainerAdapter
func NewGenericContainerAdapter[T RepositoryFactoryProvider](container T) *GenericContainerAdapter[T] {
	return &GenericContainerAdapter[T]{
		Container: container,
	}
}

// GetRepositoryFactory returns the repository factory as an interface{}
func (a *GenericContainerAdapter[T]) GetRepositoryFactory() any {
	return a.Container.GetRepositoryFactory()
}
