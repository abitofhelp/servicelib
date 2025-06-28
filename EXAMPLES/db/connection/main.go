//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of connecting to different database types
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/db"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

func main() {
	// Create a context
	// Note: ctx and contextLogger variables are only used in the commented-out code below
	// that would be used in a real application
	ctx := context.Background()

	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Sync()

	// Create a context logger
	contextLogger := logging.NewContextLogger(logger)

	// Suppress unused variable warnings
	_ = ctx
	_ = contextLogger

	// Example 1: Connect to PostgreSQL
	fmt.Println("=== PostgreSQL Connection Example ===")
	postgresURI := "postgres://postgres:password@localhost:5432/testdb"

	// In a real application, you would use this code:
	// postgresPool, err := db.InitPostgresPool(ctx, postgresURI, db.DefaultTimeout)
	// if err != nil {
	//     fmt.Printf("Failed to connect to PostgreSQL: %v\n", err)
	// } else {
	//     defer postgresPool.Close()
	//     db.LogDatabaseConnection(ctx, contextLogger, "PostgreSQL")
	//     fmt.Println("Successfully connected to PostgreSQL")
	// }

	// For the example, we'll just print what would happen
	fmt.Printf("Would connect to PostgreSQL with URI: %s\n", postgresURI)
	fmt.Printf("Using default timeout: %v\n", db.DefaultTimeout)
	fmt.Println("On success, would log connection and store pool for later use")

	// Example 2: Connect to SQLite
	fmt.Println("\n=== SQLite Connection Example ===")
	sqliteURI := "file:test.db?cache=shared&mode=memory"

	// In a real application, you would use this code:
	// sqliteDB, err := db.InitSQLiteDB(
	//     ctx,
	//     sqliteURI,
	//     db.DefaultTimeout,
	//     5*time.Minute, // Connection max lifetime
	//     10,            // Max open connections
	//     5,             // Max idle connections
	// )
	// if err != nil {
	//     fmt.Printf("Failed to connect to SQLite: %v\n", err)
	// } else {
	//     defer sqliteDB.Close()
	//     db.LogDatabaseConnection(ctx, contextLogger, "SQLite")
	//     fmt.Println("Successfully connected to SQLite")
	// }

	// For the example, we'll just print what would happen
	fmt.Printf("Would connect to SQLite with URI: %s\n", sqliteURI)
	fmt.Printf("Using default timeout: %v\n", db.DefaultTimeout)
	fmt.Println("Connection pool settings:")
	fmt.Println("- Max open connections: 10")
	fmt.Println("- Max idle connections: 5")
	fmt.Println("- Connection max lifetime: 5 minutes")
	fmt.Println("On success, would log connection and store DB for later use")

	// Example 3: Connect to MongoDB
	fmt.Println("\n=== MongoDB Connection Example ===")
	mongoURI := "mongodb://localhost:27017"

	// In a real application, you would use this code:
	// mongoClient, err := db.InitMongoClient(ctx, mongoURI, db.DefaultTimeout)
	// if err != nil {
	//     fmt.Printf("Failed to connect to MongoDB: %v\n", err)
	// } else {
	//     defer mongoClient.Disconnect(ctx)
	//     db.LogDatabaseConnection(ctx, contextLogger, "MongoDB")
	//     fmt.Println("Successfully connected to MongoDB")
	// }

	// For the example, we'll just print what would happen
	fmt.Printf("Would connect to MongoDB with URI: %s\n", mongoURI)
	fmt.Printf("Using default timeout: %v\n", db.DefaultTimeout)
	fmt.Println("On success, would log connection and store client for later use")

	// Example 4: Connection with custom timeout
	fmt.Println("\n=== Custom Timeout Connection Example ===")
	customTimeout := 10 * time.Second

	fmt.Printf("Would connect to database with custom timeout: %v\n", customTimeout)
	fmt.Println("This is useful for environments with slower network connections")

	// Expected output:
	// === PostgreSQL Connection Example ===
	// Would connect to PostgreSQL with URI: postgres://postgres:password@localhost:5432/testdb
	// Using default timeout: 30s
	// On success, would log connection and store pool for later use
	//
	// === SQLite Connection Example ===
	// Would connect to SQLite with URI: file:test.db?cache=shared&mode=memory
	// Using default timeout: 30s
	// Connection pool settings:
	// - Max open connections: 10
	// - Max idle connections: 5
	// - Connection max lifetime: 5 minutes
	// On success, would log connection and store DB for later use
	//
	// === MongoDB Connection Example ===
	// Would connect to MongoDB with URI: mongodb://localhost:27017
	// Using default timeout: 30s
	// On success, would log connection and store client for later use
	//
	// === Custom Timeout Connection Example ===
	// Would connect to database with custom timeout: 10s
	// This is useful for environments with slower network connections
}