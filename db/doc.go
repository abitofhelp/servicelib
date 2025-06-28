// Copyright (c) 2025 A Bit of Help, Inc.

// Package db provides utilities for working with various database systems.
//
// This package offers a unified interface for connecting to, querying, and managing
// different types of databases, including PostgreSQL, MongoDB, and SQLite. It provides
// functions for initializing database connections, executing transactions, checking
// database health, and handling common database errors.
//
// The package is designed to work with the context package for proper cancellation
// and timeout handling, the telemetry package for tracing database operations,
// and the retry package for handling transient errors.
//
// Key features:
//   - Database initialization for PostgreSQL, MongoDB, and SQLite
//   - Transaction management with automatic retries for transient errors
//   - Health check functions for monitoring database connectivity
//   - Error handling with classification of transient vs. permanent errors
//   - Telemetry integration for tracing database operations
//   - Logging of database operations and errors
//
// Example usage for PostgreSQL:
//
//	// Configure PostgreSQL connection
//	config := db.PostgresConfig{
//	    ConnectionString: "postgres://user:password@localhost:5432/mydb",
//	    MaxConnections:   10,
//	    ConnectTimeout:   5 * time.Second,
//	    MaxConnLifetime:  30 * time.Minute,
//	    MaxConnIdleTime:  5 * time.Minute,
//	}
//
//	// Initialize PostgreSQL connection pool
//	pool, err := db.InitPostgresPool(ctx, config)
//	if err != nil {
//	    log.Fatalf("Failed to connect to database: %v", err)
//	}
//	defer pool.Close()
//
//	// Execute a transaction
//	err = db.ExecutePostgresTransaction(ctx, pool, func(tx pgx.Tx) error {
//	    // Perform database operations within the transaction
//	    _, err := tx.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", "John", "john@example.com")
//	    return err
//	})
//
// Example usage for MongoDB:
//
//	// Initialize MongoDB client
//	client, err := db.InitMongoClient(ctx, "mongodb://localhost:27017", 5*time.Second)
//	if err != nil {
//	    log.Fatalf("Failed to connect to MongoDB: %v", err)
//	}
//	defer client.Disconnect(ctx)
//
//	// Check MongoDB health
//	if err := db.CheckMongoHealth(ctx, client); err != nil {
//	    log.Printf("MongoDB health check failed: %v", err)
//	}
//
// Example usage for SQLite:
//
//	// Initialize SQLite database
//	sqliteDB, err := db.InitSQLiteDB(ctx, "file:mydb.sqlite", 5*time.Second, 30*time.Minute, 10, 5)
//	if err != nil {
//	    log.Fatalf("Failed to connect to SQLite: %v", err)
//	}
//	defer sqliteDB.Close()
//
//	// Execute a transaction
//	err = db.ExecuteSQLTransaction(ctx, sqliteDB, func(tx *sql.Tx) error {
//	    // Perform database operations within the transaction
//	    _, err := tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "John", "john@example.com")
//	    return err
//	})
//
// The package is designed to be used throughout the application to provide
// consistent database access patterns and error handling.
package db
