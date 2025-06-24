// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using transactions with different database types
package main

import (
	"context"
	"fmt"

	// These imports are used in the commented-out code below
	// that would be used in a real application
	// Using blank imports to suppress unused import warnings
	_ "database/sql"
	_ "github.com/abitofhelp/servicelib/db"
	_ "github.com/jackc/pgx/v5"
)

func main() {
	// Create a context
	// Note: ctx is only used in the commented-out code below
	// that would be used in a real application
	ctx := context.Background()

	// Suppress unused variable warnings
	_ = ctx

	// Example 1: PostgreSQL Transaction
	fmt.Println("=== PostgreSQL Transaction Example ===")

	// In a real application, you would use this code:
	// postgresPool, err := db.InitPostgresPool(ctx, "postgres://postgres:password@localhost:5432/testdb", db.DefaultTimeout)
	// if err != nil {
	//     fmt.Printf("Failed to connect to PostgreSQL: %v\n", err)
	//     return
	// }
	// defer postgresPool.Close()
	//
	// // Execute a transaction
	// err = db.ExecutePostgresTransaction(ctx, postgresPool, func(tx pgx.Tx) error {
	//     // First operation: Insert a user
	//     _, err := tx.Exec(ctx, "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
	//         "user123", "John Doe", "john.doe@example.com")
	//     if err != nil {
	//         return fmt.Errorf("failed to insert user: %w", err)
	//     }
	//
	//     // Second operation: Insert user preferences
	//     _, err = tx.Exec(ctx, "INSERT INTO user_preferences (user_id, theme, notifications) VALUES ($1, $2, $3)",
	//         "user123", "dark", true)
	//     if err != nil {
	//         return fmt.Errorf("failed to insert user preferences: %w", err)
	//     }
	//
	//     return nil
	// })
	//
	// if err != nil {
	//     fmt.Printf("Transaction failed: %v\n", err)
	// } else {
	//     fmt.Println("Transaction completed successfully")
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would execute a PostgreSQL transaction with the following operations:")
	fmt.Println("1. Insert a user into the users table")
	fmt.Println("2. Insert user preferences into the user_preferences table")
	fmt.Println("If any operation fails, the entire transaction is rolled back")
	fmt.Println("This ensures data consistency across related tables")

	// Example 2: SQLite Transaction
	fmt.Println("\n=== SQLite Transaction Example ===")

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
	// // Execute a transaction
	// err = db.ExecuteSQLTransaction(ctx, sqliteDB, func(tx *sql.Tx) error {
	//     // First operation: Create a product
	//     _, err := tx.Exec("INSERT INTO products (id, name, price) VALUES (?, ?, ?)",
	//         "prod456", "Smartphone", 999.99)
	//     if err != nil {
	//         return fmt.Errorf("failed to insert product: %w", err)
	//     }
	//
	//     // Second operation: Update inventory
	//     _, err = tx.Exec("UPDATE inventory SET quantity = quantity - ? WHERE product_id = ?",
	//         1, "prod456")
	//     if err != nil {
	//         return fmt.Errorf("failed to update inventory: %w", err)
	//     }
	//
	//     return nil
	// })
	//
	// if err != nil {
	//     fmt.Printf("Transaction failed: %v\n", err)
	// } else {
	//     fmt.Println("Transaction completed successfully")
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would execute a SQLite transaction with the following operations:")
	fmt.Println("1. Insert a product into the products table")
	fmt.Println("2. Update inventory quantity for the product")
	fmt.Println("SQLite transactions provide ACID guarantees for local database operations")

	// Example 3: Transaction with Error Handling
	fmt.Println("\n=== Transaction with Error Handling Example ===")

	// In a real application, you would use this code:
	// err = db.ExecuteSQLTransaction(ctx, sqliteDB, func(tx *sql.Tx) error {
	//     // Check if product exists
	//     var count int
	//     err := tx.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", "prod789").Scan(&count)
	//     if err != nil {
	//         return fmt.Errorf("failed to check product existence: %w", err)
	//     }
	//
	//     if count == 0 {
	//         return fmt.Errorf("product not found: prod789")
	//     }
	//
	//     // Update product price
	//     _, err = tx.Exec("UPDATE products SET price = ? WHERE id = ?", 1299.99, "prod789")
	//     if err != nil {
	//         return fmt.Errorf("failed to update product price: %w", err)
	//     }
	//
	//     return nil
	// })
	//
	// if err != nil {
	//     fmt.Printf("Transaction failed: %v\n", err)
	// } else {
	//     fmt.Println("Transaction completed successfully")
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would execute a transaction with error handling:")
	fmt.Println("1. Check if product exists")
	fmt.Println("2. If product exists, update its price")
	fmt.Println("3. If product doesn't exist, return an error")
	fmt.Println("Error handling within transactions ensures data integrity")

	// Example 4: Nested Transactions (Not Supported in Most Databases)
	fmt.Println("\n=== Nested Transactions Example ===")
	fmt.Println("Most SQL databases don't support true nested transactions")
	fmt.Println("Instead, you can use savepoints in databases that support them")
	fmt.Println("For PostgreSQL, you might use:")
	fmt.Println("tx.Exec(ctx, \"SAVEPOINT my_savepoint\")")
	fmt.Println("// ... operations ...")
	fmt.Println("// If error: tx.Exec(ctx, \"ROLLBACK TO SAVEPOINT my_savepoint\")")
	fmt.Println("// If success: tx.Exec(ctx, \"RELEASE SAVEPOINT my_savepoint\")")

	// Example 5: Transaction Isolation Levels
	fmt.Println("\n=== Transaction Isolation Levels Example ===")
	fmt.Println("SQL databases support different transaction isolation levels:")
	fmt.Println("1. READ UNCOMMITTED - Lowest isolation, allows dirty reads")
	fmt.Println("2. READ COMMITTED - Prevents dirty reads")
	fmt.Println("3. REPEATABLE READ - Prevents dirty and non-repeatable reads")
	fmt.Println("4. SERIALIZABLE - Highest isolation, prevents all concurrency issues")

	fmt.Println("\nIn a real application, you would set isolation level like this:")
	fmt.Println("tx, err := db.BeginTx(ctx, &sql.TxOptions{")
	fmt.Println("    Isolation: sql.LevelSerializable,")
	fmt.Println("    ReadOnly: false,")
	fmt.Println("})")

	// Expected output:
	// === PostgreSQL Transaction Example ===
	// Would execute a PostgreSQL transaction with the following operations:
	// 1. Insert a user into the users table
	// 2. Insert user preferences into the user_preferences table
	// If any operation fails, the entire transaction is rolled back
	// This ensures data consistency across related tables
	//
	// === SQLite Transaction Example ===
	// Would execute a SQLite transaction with the following operations:
	// 1. Insert a product into the products table
	// 2. Update inventory quantity for the product
	// SQLite transactions provide ACID guarantees for local database operations
	//
	// === Transaction with Error Handling Example ===
	// Would execute a transaction with error handling:
	// 1. Check if product exists
	// 2. If product exists, update its price
	// 3. If product doesn't exist, return an error
	// Error handling within transactions ensures data integrity
	//
	// === Nested Transactions Example ===
	// Most SQL databases don't support true nested transactions
	// Instead, you can use savepoints in databases that support them
	// For PostgreSQL, you might use:
	// tx.Exec(ctx, "SAVEPOINT my_savepoint")
	// // ... operations ...
	// // If error: tx.Exec(ctx, "ROLLBACK TO SAVEPOINT my_savepoint")
	// // If success: tx.Exec(ctx, "RELEASE SAVEPOINT my_savepoint")
	//
	// === Transaction Isolation Levels Example ===
	// SQL databases support different transaction isolation levels:
	// 1. READ UNCOMMITTED - Lowest isolation, allows dirty reads
	// 2. READ COMMITTED - Prevents dirty reads
	// 3. REPEATABLE READ - Prevents dirty and non-repeatable reads
	// 4. SERIALIZABLE - Highest isolation, prevents all concurrency issues
	//
	// In a real application, you would set isolation level like this:
	// tx, err := db.BeginTx(ctx, &sql.TxOptions{
	//     Isolation: sql.LevelSerializable,
	//     ReadOnly: false,
	// })
}