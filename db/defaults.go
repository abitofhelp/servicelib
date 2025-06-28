// Copyright (c) 2025 A Bit of Help, Inc.

package db

import "time"

// DefaultPostgresConfig returns a default configuration for a PostgreSQL connection pool.
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		URI:               "postgres://postgres:postgres@localhost:5432/postgres",
		Timeout:           10 * time.Second,
		MaxConns:          10,
		MinConns:          2,
		MaxConnLifetime:   1 * time.Hour,
		MaxConnIdleTime:   30 * time.Minute,
		HealthCheckPeriod: 1 * time.Minute,
	}
}
