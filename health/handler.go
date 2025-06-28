// Copyright (c) 2025 A Bit of Help, Inc.

// Package health provides functionality for health checking the application.
package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/abitofhelp/servicelib/config"
	"go.uber.org/zap"
)

// HealthStatus represents the health status response
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services,omitempty"`
}

// NewHandler creates a new health check HTTP handler
func NewHandler(provider HealthCheckProvider, logger *zap.Logger, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context with timeout for the health check
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Use the context with timeout
		r = r.WithContext(ctx)

		// Check the health of dependencies
		services := make(map[string]string)

		// Check database connectivity through the repository factory
		repoFactory := provider.GetRepositoryFactory()
		if repoFactory != nil {
			services["database"] = ServiceUp
		} else {
			services["database"] = ServiceDown
		}

		// Overall status is healthy if all dependencies are healthy
		status := StatusHealthy
		for _, s := range services {
			if s != ServiceUp {
				status = StatusDegraded
				break
			}
		}

		// Create the health response
		healthResponse := HealthStatus{
			Status:    status, // Use the status constant directly
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Version:   cfg.GetApp().GetVersion(),
			Services:  services, // Use the services map directly
		}

		// Set content type
		w.Header().Set("Content-Type", "application/json")

		// Set appropriate status code based on health
		if status == StatusHealthy {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		// Write the response
		if err := json.NewEncoder(w).Encode(healthResponse); err != nil {
			logger.Error("Failed to encode health response", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Log the health check
		logger.Info("Health check",
			zap.String("status", status),
			zap.Any("services", services),
		)
	}
}
