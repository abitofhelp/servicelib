// Copyright (c) 2025 A Bit of Help, Inc.

// Package health provides functionality for health checking the application.
package health

import (
	"github.com/abitofhelp/family-service/cmd/server/graphql/di"
)

// ContainerAdapter adapts the di.Container to implement HealthCheckProvider
type ContainerAdapter struct {
	Container *di.Container
}

// NewContainerAdapter creates a new ContainerAdapter
func NewContainerAdapter(container *di.Container) *ContainerAdapter {
	return &ContainerAdapter{
		Container: container,
	}
}

// GetRepositoryFactory returns the repository factory as an interface{}
func (a *ContainerAdapter) GetRepositoryFactory() any {
	return a.Container.GetRepositoryFactory()
}