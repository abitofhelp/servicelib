// Copyright (c) 2025 A Bit of Help, Inc.

// Example of executing queries with retries
package main

import (
	"context"
	"fmt"

	// These imports are used in the commented-out code below
	// that would be used in a real application
	// Using blank imports to suppress unused import warnings
	_ "database/sql"
	_ "github.com/abitofhelp/servicelib/db"
	_ "time"
)

func main() {
	// Create a context
	// Note: ctx is only used in the commented-out code below
	// that would be used in a real application
	ctx := context.Background()

	// Suppress unused variable warnings
	_ = ctx

	// Example 1: Basic Query with Retries
	fmt.Println("=== Basic Query with Retries Example ===")

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
	// // Configure retry options
	// options := db.RetryOptions{
	//     MaxRetries:   3,
	//     InitialDelay: 100 * time.Millisecond,
	//     MaxDelay:     1 * time.Second,
	//     Multiplier:   2.0,
	// }
	//
	// // Execute query with retries
	// var count int
	// err = db.QueryWithRetries(ctx, sqliteDB, options, func(ctx context.Context, db *sql.DB) error {
	//     return db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	// })
	//
	// if err != nil {
	//     fmt.Printf("Query failed after retries: %v\n", err)
	// } else {
	//     fmt.Printf("User count: %d\n", count)
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would execute a query with retries to count users:")
	fmt.Println("SQL: SELECT COUNT(*) FROM users")
	fmt.Println("Retry configuration:")
	fmt.Println("- Maximum retries: 3")
	fmt.Println("- Initial delay: 100ms")
	fmt.Println("- Maximum delay: 1s")
	fmt.Println("- Backoff multiplier: 2.0")
	fmt.Println("This means:")
	fmt.Println("- First retry after 100ms")
	fmt.Println("- Second retry after 200ms")
	fmt.Println("- Third retry after 400ms")
	fmt.Println("Retries help handle transient database errors")

	// Example 2: Customizing Retry Behavior
	fmt.Println("\n=== Customizing Retry Behavior Example ===")

	// In a real application, you would use this code:
	// // Configure custom retry options
	// customOptions := db.RetryOptions{
	//     MaxRetries:   5,
	//     InitialDelay: 50 * time.Millisecond,
	//     MaxDelay:     2 * time.Second,
	//     Multiplier:   1.5,
	// }
	//
	// // Execute query with custom retry options
	// var productCount int
	// err = db.QueryWithRetries(ctx, sqliteDB, customOptions, func(ctx context.Context, db *sql.DB) error {
	//     return db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products").Scan(&productCount)
	// })
	//
	// if err != nil {
	//     fmt.Printf("Query failed after retries: %v\n", err)
	// } else {
	//     fmt.Printf("Product count: %d\n", productCount)
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would execute a query with custom retry options:")
	fmt.Println("Custom retry configuration:")
	fmt.Println("- Maximum retries: 5")
	fmt.Println("- Initial delay: 50ms")
	fmt.Println("- Maximum delay: 2s")
	fmt.Println("- Backoff multiplier: 1.5")
	fmt.Println("This means:")
	fmt.Println("- First retry after 50ms")
	fmt.Println("- Second retry after 75ms")
	fmt.Println("- Third retry after 112.5ms")
	fmt.Println("- Fourth retry after 168.75ms")
	fmt.Println("- Fifth retry after 253.125ms")
	fmt.Println("Customize retry options based on the operation's importance and expected duration")

	// Example 3: Retrying Transactions
	fmt.Println("\n=== Retrying Transactions Example ===")

	// In a real application, you would use this code:
	// // Configure retry options for transaction
	// txRetryOptions := db.RetryOptions{
	//     MaxRetries:   2,
	//     InitialDelay: 200 * time.Millisecond,
	//     MaxDelay:     1 * time.Second,
	//     Multiplier:   2.0,
	// }
	//
	// // Execute transaction with retries
	// err = db.TransactionWithRetries(ctx, sqliteDB, txRetryOptions, func(ctx context.Context, tx *sql.Tx) error {
	//     // First operation
	//     _, err := tx.ExecContext(ctx, "INSERT INTO users (id, username) VALUES (?, ?)", "user456", "newuser")
	//     if err != nil {
	//         return err
	//     }
	//
	//     // Second operation
	//     _, err = tx.ExecContext(ctx, "UPDATE user_counts SET count = count + 1")
	//     if err != nil {
	//         return err
	//     }
	//
	//     return nil
	// })
	//
	// if err != nil {
	//     fmt.Printf("Transaction failed after retries: %v\n", err)
	// } else {
	//     fmt.Println("Transaction completed successfully")
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would execute a transaction with retries:")
	fmt.Println("Transaction operations:")
	fmt.Println("1. Insert a new user")
	fmt.Println("2. Update user count")
	fmt.Println("Retry configuration:")
	fmt.Println("- Maximum retries: 2")
	fmt.Println("- Initial delay: 200ms")
	fmt.Println("- Maximum delay: 1s")
	fmt.Println("- Backoff multiplier: 2.0")
	fmt.Println("Retrying transactions is useful for handling deadlocks or lock timeouts")

	// Example 4: Handling Specific Errors
	fmt.Println("\n=== Handling Specific Errors Example ===")

	// In a real application, you would use this code:
	// // Configure retry options with custom retry condition
	// customRetryOptions := db.RetryOptionsWithCondition{
	//     RetryOptions: db.RetryOptions{
	//         MaxRetries:   3,
	//         InitialDelay: 100 * time.Millisecond,
	//         MaxDelay:     1 * time.Second,
	//         Multiplier:   2.0,
	//     },
	//     ShouldRetry: func(err error) bool {
	//         // Only retry on specific errors like connection issues or deadlocks
	//         return db.IsConnectionError(err) || db.IsDeadlockError(err)
	//     },
	// }
	//
	// // Execute query with custom retry condition
	// var result string
	// err = db.QueryWithRetriesAndCondition(ctx, sqliteDB, customRetryOptions, func(ctx context.Context, db *sql.DB) error {
	//     return db.QueryRowContext(ctx, "SELECT value FROM settings WHERE key = ?", "app_name").Scan(&result)
	// })
	//
	// if err != nil {
	//     fmt.Printf("Query failed after retries: %v\n", err)
	// } else {
	//     fmt.Printf("Setting value: %s\n", result)
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would execute a query with custom retry condition:")
	fmt.Println("Custom retry condition:")
	fmt.Println("- Only retry on connection issues or deadlocks")
	fmt.Println("- Don't retry on other errors (e.g., syntax errors)")
	fmt.Println("This is more efficient than retrying on all errors")
	fmt.Println("Some errors (like syntax errors) will never succeed on retry")

	// Example 5: Exponential Backoff with Jitter
	fmt.Println("\n=== Exponential Backoff with Jitter Example ===")

	// In a real application, you would use this code:
	// // Configure retry options with jitter
	// jitterOptions := db.RetryOptionsWithJitter{
	//     RetryOptions: db.RetryOptions{
	//         MaxRetries:   3,
	//         InitialDelay: 100 * time.Millisecond,
	//         MaxDelay:     1 * time.Second,
	//         Multiplier:   2.0,
	//     },
	//     JitterFactor: 0.2, // 20% jitter
	// }
	//
	// // Execute query with jitter
	// var userID string
	// err = db.QueryWithRetriesAndJitter(ctx, sqliteDB, jitterOptions, func(ctx context.Context, db *sql.DB) error {
	//     return db.QueryRowContext(ctx, "SELECT id FROM users WHERE username = ?", "admin").Scan(&userID)
	// })
	//
	// if err != nil {
	//     fmt.Printf("Query failed after retries: %v\n", err)
	// } else {
	//     fmt.Printf("User ID: %s\n", userID)
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would execute a query with exponential backoff and jitter:")
	fmt.Println("Jitter adds randomness to retry delays to prevent thundering herd problems")
	fmt.Println("With 20% jitter and 100ms initial delay:")
	fmt.Println("- First retry might be after 80-120ms (100ms ± 20%)")
	fmt.Println("- Second retry might be after 160-240ms (200ms ± 20%)")
	fmt.Println("- Third retry might be after 320-480ms (400ms ± 20%)")
	fmt.Println("Jitter is important in high-concurrency environments")

	// Expected output:
	// === Basic Query with Retries Example ===
	// Would execute a query with retries to count users:
	// SQL: SELECT COUNT(*) FROM users
	// Retry configuration:
	// - Maximum retries: 3
	// - Initial delay: 100ms
	// - Maximum delay: 1s
	// - Backoff multiplier: 2.0
	// This means:
	// - First retry after 100ms
	// - Second retry after 200ms
	// - Third retry after 400ms
	// Retries help handle transient database errors
	//
	// === Customizing Retry Behavior Example ===
	// Would execute a query with custom retry options:
	// Custom retry configuration:
	// - Maximum retries: 5
	// - Initial delay: 50ms
	// - Maximum delay: 2s
	// - Backoff multiplier: 1.5
	// This means:
	// - First retry after 50ms
	// - Second retry after 75ms
	// - Third retry after 112.5ms
	// - Fourth retry after 168.75ms
	// - Fifth retry after 253.125ms
	// Customize retry options based on the operation's importance and expected duration
	//
	// === Retrying Transactions Example ===
	// Would execute a transaction with retries:
	// Transaction operations:
	// 1. Insert a new user
	// 2. Update user count
	// Retry configuration:
	// - Maximum retries: 2
	// - Initial delay: 200ms
	// - Maximum delay: 1s
	// - Backoff multiplier: 2.0
	// Retrying transactions is useful for handling deadlocks or lock timeouts
	//
	// === Handling Specific Errors Example ===
	// Would execute a query with custom retry condition:
	// Custom retry condition:
	// - Only retry on connection issues or deadlocks
	// - Don't retry on other errors (e.g., syntax errors)
	// This is more efficient than retrying on all errors
	// Some errors (like syntax errors) will never succeed on retry
	//
	// === Exponential Backoff with Jitter Example ===
	// Would execute a query with exponential backoff and jitter:
	// Jitter adds randomness to retry delays to prevent thundering herd problems
	// With 20% jitter and 100ms initial delay:
	// - First retry might be after 80-120ms (100ms ± 20%)
	// - Second retry might be after 160-240ms (200ms ± 20%)
	// - Third retry might be after 320-480ms (400ms ± 20%)
	// Jitter is important in high-concurrency environments
}