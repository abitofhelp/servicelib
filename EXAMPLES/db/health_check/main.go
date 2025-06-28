// Copyright (c) 2025 A Bit of Help, Inc.

// Example of checking database health
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/db"
)

func main() {
	// Create a context
	// Note: ctx is only used in the commented-out code below
	// that would be used in a real application
	ctx := context.Background()

	// Suppress unused variable warnings
	_ = ctx

	// Example 1: Check PostgreSQL health
	fmt.Println("=== PostgreSQL Health Check Example ===")

	// In a real application, you would use this code:
	// postgresPool, err := db.InitPostgresPool(ctx, "postgres://postgres:password@localhost:5432/testdb", db.DefaultTimeout)
	// if err != nil {
	//     fmt.Printf("Failed to connect to PostgreSQL: %v\n", err)
	//     return
	// }
	// defer postgresPool.Close()
	//
	// // Check health
	// if err := db.CheckPostgresHealth(ctx, postgresPool); err != nil {
	//     fmt.Printf("PostgreSQL health check failed: %v\n", err)
	// } else {
	//     fmt.Println("PostgreSQL health check passed")
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would check PostgreSQL health by pinging the database")
	fmt.Printf("Using default timeout: %v\n", db.DefaultTimeout)
	fmt.Println("Health check helps ensure the database is responsive")

	// Example 2: Check SQLite health
	fmt.Println("\n=== SQLite Health Check Example ===")

	// In a real application, you would use this code:
	// sqliteDB, err := db.InitSQLiteDB(
	//     ctx,
	//     "file:test.db?cache=shared&mode=memory",
	//     db.DefaultTimeout,
	//     5*time.Minute,
	//     10,
	//     5,
	// )
	// if err != nil {
	//     fmt.Printf("Failed to connect to SQLite: %v\n", err)
	//     return
	// }
	// defer sqliteDB.Close()
	//
	// // Check health
	// if err := db.CheckSQLiteHealth(ctx, sqliteDB); err != nil {
	//     fmt.Printf("SQLite health check failed: %v\n", err)
	// } else {
	//     fmt.Println("SQLite health check passed")
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would check SQLite health by pinging the database")
	fmt.Printf("Using default timeout: %v\n", db.DefaultTimeout)
	fmt.Println("Health check verifies the database file is accessible and not corrupted")

	// Example 3: Check MongoDB health
	fmt.Println("\n=== MongoDB Health Check Example ===")

	// In a real application, you would use this code:
	// mongoClient, err := db.InitMongoClient(ctx, "mongodb://localhost:27017", db.DefaultTimeout)
	// if err != nil {
	//     fmt.Printf("Failed to connect to MongoDB: %v\n", err)
	//     return
	// }
	// defer mongoClient.Disconnect(ctx)
	//
	// // Check health
	// if err := db.CheckMongoHealth(ctx, mongoClient); err != nil {
	//     fmt.Printf("MongoDB health check failed: %v\n", err)
	// } else {
	//     fmt.Println("MongoDB health check passed")
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would check MongoDB health by pinging the database")
	fmt.Printf("Using default timeout: %v\n", db.DefaultTimeout)
	fmt.Println("Health check confirms the MongoDB server is running and responsive")

	// Example 4: Periodic health checks
	fmt.Println("\n=== Periodic Health Check Example ===")
	fmt.Println("In a production application, you might want to run health checks periodically")

	// In a real application, you would use this code:
	// ticker := time.NewTicker(1 * time.Minute)
	// defer ticker.Stop()
	//
	// go func() {
	//     for {
	//         select {
	//         case <-ticker.C:
	//             // Check database health
	//             if err := db.CheckPostgresHealth(ctx, postgresPool); err != nil {
	//                 fmt.Printf("PostgreSQL health check failed: %v\n", err)
	//                 // Take action, e.g., reconnect or alert
	//             }
	//         case <-ctx.Done():
	//             return
	//         }
	//     }
	// }()

	// For the example, we'll just print what would happen
	fmt.Println("Would set up a ticker to check health every minute")
	fmt.Println("Health check would run in a separate goroutine")
	fmt.Println("If health check fails, application could reconnect or alert administrators")

	// Example 5: Health check with custom timeout
	fmt.Println("\n=== Custom Timeout Health Check Example ===")
	customTimeout := 5 * time.Second

	// In a real application, you would use this code:
	// healthCtx, cancel := context.WithTimeout(ctx, customTimeout)
	// defer cancel()
	//
	// if err := db.CheckPostgresHealth(healthCtx, postgresPool); err != nil {
	//     fmt.Printf("Health check with custom timeout failed: %v\n", err)
	// } else {
	//     fmt.Println("Health check with custom timeout passed")
	// }

	// For the example, we'll just print what would happen
	fmt.Printf("Would check health with custom timeout: %v\n", customTimeout)
	fmt.Println("Custom timeout is useful for health checks in different environments")
	fmt.Println("For example, shorter timeouts for critical services, longer for non-critical")

	// Expected output:
	// === PostgreSQL Health Check Example ===
	// Would check PostgreSQL health by pinging the database
	// Using default timeout: 30s
	// Health check helps ensure the database is responsive
	//
	// === SQLite Health Check Example ===
	// Would check SQLite health by pinging the database
	// Using default timeout: 30s
	// Health check verifies the database file is accessible and not corrupted
	//
	// === MongoDB Health Check Example ===
	// Would check MongoDB health by pinging the database
	// Using default timeout: 30s
	// Health check confirms the MongoDB server is running and responsive
	//
	// === Periodic Health Check Example ===
	// In a production application, you might want to run health checks periodically
	// Would set up a ticker to check health every minute
	// Health check would run in a separate goroutine
	// If health check fails, application could reconnect or alert administrators
	//
	// === Custom Timeout Health Check Example ===
	// Would check health with custom timeout: 5s
	// Custom timeout is useful for health checks in different environments
	// For example, shorter timeouts for critical services, longer for non-critical
}