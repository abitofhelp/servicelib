// Copyright (c) 2025 A Bit of Help, Inc.

// Package health provides functionality for health checking applications.
package health

import (
	"github.com/abitofhelp/servicelib/repository"
)

// RepositoryAdapter adapts a repository.RepositoryFactory to the RepositoryFactoryProvider interface
type RepositoryAdapter struct {
	factory repository.RepositoryFactory
}

// NewRepositoryAdapter creates a new RepositoryAdapter
func NewRepositoryAdapter(factory repository.RepositoryFactory) *RepositoryAdapter {
	return &RepositoryAdapter{
		factory: factory,
	}
}

// GetRepositoryFactory returns the repository factory
func (a *RepositoryAdapter) GetRepositoryFactory() any {
	return a.factory.GetRepository()
}