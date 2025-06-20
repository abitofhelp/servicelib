// Copyright (c) 2025 A Bit of Help, Inc.

// Example of creating a custom health status response
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/abitofhelp/servicelib/config"
	"github.com/abitofhelp/servicelib/health"
	"go.uber.org/zap"
)

// CustomHealthProvider implements the health.HealthCheckProvider interface
type CustomHealthProvider struct{}

// GetRepositoryFactory returns a mock repository factory
func (p *CustomHealthProvider) GetRepositoryFactory() any {
	// In a real application, this would return an actual repository factory
	// For this example, we'll just return a non-nil value
	return &struct{}{}
}

// CustomConfig implements the config.Config interface
type CustomConfig struct{}

// GetApp returns the application configuration
func (c *CustomConfig) GetApp() config.AppConfig {
	return &CustomAppConfig{}
}

// GetDatabase returns the database configuration
func (c *CustomConfig) GetDatabase() config.DatabaseConfig {
	return &CustomDatabaseConfig{}
}

// CustomAppConfig implements the config.AppConfig interface
type CustomAppConfig struct{}

// GetVersion returns the application version
func (a *CustomAppConfig) GetVersion() string {
	return "1.0.0"
}

// GetName returns the application name
func (a *CustomAppConfig) GetName() string {
	return "my-service"
}

// GetEnvironment returns the application environment
func (a *CustomAppConfig) GetEnvironment() string {
	return "development"
}

// CustomDatabaseConfig implements the config.DatabaseConfig interface
type CustomDatabaseConfig struct{}

// GetType returns the database type
func (d *CustomDatabaseConfig) GetType() string {
	return "postgres"
}

// GetConnectionString returns the database connection string
func (d *CustomDatabaseConfig) GetConnectionString() string {
	return "postgres://user:password@localhost:5432/mydb?sslmode=disable"
}

// GetDatabaseName returns the database name
func (d *CustomDatabaseConfig) GetDatabaseName() string {
	return "mydb"
}

// GetCollectionName returns the collection/table name for a given entity type
func (d *CustomDatabaseConfig) GetCollectionName(entityType string) string {
	return entityType + "s"
}

// CustomHealthStatus represents a custom health status response
type CustomHealthStatus struct {
	Status      string            `json:"status"`
	Timestamp   string            `json:"timestamp"`
	Version     string            `json:"version"`
	Environment string            `json:"environment"`
	Services    map[string]string `json:"services,omitempty"`
	Uptime      string            `json:"uptime"`
	Memory      string            `json:"memory"`
	Goroutines  int               `json:"goroutines"`
}

var startTime = time.Now()

func main() {
	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create a health check provider
	provider := &CustomHealthProvider{}

	// Create a configuration
	cfg := &CustomConfig{}

	// Create a custom health check handler
	http.HandleFunc("/health/custom", func(w http.ResponseWriter, r *http.Request) {
		// Create a standard health handler to reuse its functionality
		standardHandler := health.NewHandler(provider, logger, cfg)

		// Create a new request to pass to the standard handler
		req, _ := http.NewRequest(http.MethodGet, "/health", nil)

		// Create a response recorder to capture the standard handler's response
		recorder := &responseRecorder{
			header: http.Header{},
		}

		// Call the standard handler
		standardHandler(recorder, req)

		// Parse the standard health status
		var standardStatus health.HealthStatus
		if err := json.Unmarshal(recorder.body, &standardStatus); err != nil {
			logger.Error("Failed to parse standard health status", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get memory statistics
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		// Create a custom health status
		customStatus := CustomHealthStatus{
			Status:      standardStatus.Status,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Version:     standardStatus.Version,
			Environment: cfg.GetApp().GetEnvironment(),
			Services:    standardStatus.Services,
			Uptime:      time.Since(startTime).String(),
			Memory:      fmt.Sprintf("%.2f MB", float64(memStats.Alloc)/1024/1024),
			Goroutines:  runtime.NumGoroutine(),
		}

		// Set content type
		w.Header().Set("Content-Type", "application/json")

		// Set appropriate status code based on health
		if standardStatus.Status == health.StatusHealthy {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		// Write the response
		if err := json.NewEncoder(w).Encode(customStatus); err != nil {
			logger.Error("Failed to encode custom health response", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Log the health check
		logger.Info("Custom health check",
			zap.String("status", customStatus.Status),
			zap.Any("services", customStatus.Services),
		)
	})

	// Start the server
	fmt.Println("Server starting on :8080")
	fmt.Println("Custom health check endpoint: http://localhost:8080/health/custom")

	// In a real application, you would start the server:
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	//     logger.Fatal("Failed to start server", zap.Error(err))
	// }
}

// responseRecorder is a simple implementation of http.ResponseWriter that records the response
type responseRecorder struct {
	header     http.Header
	body       []byte
	statusCode int
}

func (r *responseRecorder) Header() http.Header {
	return r.header
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = b
	return len(b), nil
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}
