// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides generic database initializers that can be used across different applications.
package di

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/db"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// MongoDBInitializer initializes a MongoDB client and collection
func MongoDBInitializer(
	ctx context.Context,
	uri string,
	databaseName string,
	collectionName string,
	logger *zap.Logger,
) (*mongo.Collection, error) {
	// Check if logger is nil
	if logger == nil {
		return nil, fmt.Errorf("failed to initialize MongoDB client: context canceled")
	}

	// Create a context logger
	contextLogger := logging.NewContextLogger(logger)

	// Initialize MongoDB client using the db package
	client, err := db.InitMongoClient(ctx, uri, DefaultTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MongoDB client: %w", err)
	}

	// Log successful connection
	db.LogDatabaseConnection(ctx, contextLogger, "MongoDB")

	// Get the collection
	collection := client.Database(databaseName).Collection(collectionName)

	return collection, nil
}

// PostgresInitializer initializes a PostgreSQL connection pool
func PostgresInitializer(
	ctx context.Context,
	dsn string,
	logger *zap.Logger,
) (*pgxpool.Pool, error) {
	// Check if logger is nil
	if logger == nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL connection pool: context canceled")
	}

	// Create a context logger
	contextLogger := logging.NewContextLogger(logger)

	// Initialize PostgreSQL connection pool using the db package
	pool, err := db.InitPostgresPool(ctx, dsn, DefaultTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL connection pool: %w", err)
	}

	// Log successful connection
	db.LogDatabaseConnection(ctx, contextLogger, "PostgreSQL")

	return pool, nil
}

// SQLiteInitializer initializes a SQLite database connection
func SQLiteInitializer(
	ctx context.Context,
	uri string,
	logger *zap.Logger,
) (*sql.DB, error) {
	// Check if logger is nil
	if logger == nil {
		return nil, fmt.Errorf("failed to initialize SQLite database connection: context canceled")
	}

	// Create a context logger
	contextLogger := logging.NewContextLogger(logger)

	// Initialize SQLite database connection using the db package
	sqliteDB, err := db.InitSQLiteDB(ctx, uri, DefaultTimeout, time.Hour, 10, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQLite database connection: %w", err)
	}

	// Log successful connection
	db.LogDatabaseConnection(ctx, contextLogger, "SQLite")

	return sqliteDB, nil
}
