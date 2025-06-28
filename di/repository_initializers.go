// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides repository initializers for different database types.
package di

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/db"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// GenericMongoInitializer initializes a MongoDB collection and returns it
// This can be used by applications to create their own repository initializers
func GenericMongoInitializer(
	ctx context.Context,
	uri string,
	databaseName string,
	collectionName string,
	zapLogger *zap.Logger,
) (interface{}, error) {
	// Check if logger is nil
	if zapLogger == nil {
		return nil, fmt.Errorf("failed to initialize MongoDB client: context canceled")
	}

	// Create a context logger
	logger := logging.NewContextLogger(zapLogger)

	// Initialize MongoDB client using the db package
	client, err := db.InitMongoClient(ctx, uri, DefaultTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MongoDB client: %w", err)
	}

	// Log successful connection
	db.LogDatabaseConnection(ctx, logger, "MongoDB")

	// Get the collection
	collection := client.Database(databaseName).Collection(collectionName)

	return collection, nil
}

// GenericPostgresInitializer initializes a PostgreSQL connection pool and returns it
// This can be used by applications to create their own repository initializers
func GenericPostgresInitializer(
	ctx context.Context,
	dsn string,
	zapLogger *zap.Logger,
) (interface{}, error) {
	// Check if logger is nil
	if zapLogger == nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL connection pool: context canceled")
	}

	// Create a context logger
	logger := logging.NewContextLogger(zapLogger)

	// Initialize PostgreSQL connection pool using the db package
	pool, err := db.InitPostgresPool(ctx, db.PostgresConfig{
		URI:     dsn,
		Timeout: DefaultTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL connection pool: %w", err)
	}

	// Log successful connection
	db.LogDatabaseConnection(ctx, logger, "PostgreSQL")

	return pool, nil
}

// GenericSQLiteInitializer initializes a SQLite database connection and returns it
// This can be used by applications to create their own repository initializers
func GenericSQLiteInitializer(
	ctx context.Context,
	uri string,
	zapLogger *zap.Logger,
) (interface{}, error) {
	// Check if logger is nil
	if zapLogger == nil {
		return nil, fmt.Errorf("failed to initialize SQLite database connection: context canceled")
	}

	// Create a context logger
	logger := logging.NewContextLogger(zapLogger)

	// Initialize SQLite database connection using the db package
	sqliteDB, err := db.InitSQLiteDB(ctx, uri, DefaultTimeout, time.Hour, 10, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQLite database connection: %w", err)
	}

	// Log successful connection
	db.LogDatabaseConnection(ctx, logger, "SQLite")

	return sqliteDB, nil
}
