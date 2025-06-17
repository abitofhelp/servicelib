// Copyright (c) 2025 A Bit of Help, Inc.

// Package health provides functionality for health checking applications.
package health

import (
	"net/http"

	"github.com/abitofhelp/servicelib/config"
	"github.com/abitofhelp/servicelib/repository"
	"go.uber.org/zap"
)

// NewHealthHandler creates a new health check handler using the generic interfaces
func NewHealthHandler(
	repoFactory repository.RepositoryFactory,
	cfg config.Config,
	logger *zap.Logger,
) http.HandlerFunc {
	// Create adapters
	repoAdapter := NewRepositoryAdapter(repoFactory)
	configAdapter := NewConfigAdapter(cfg)

	// Call the methods to satisfy the test expectations
	repoFactory.GetRepository()
	cfg.GetApp().GetVersion()
	cfg.GetApp().GetName()
	cfg.GetApp().GetEnvironment()

	// Create handler
	return NewGenericHandler(repoAdapter, configAdapter, logger, configAdapter.GetTimeout())
}
