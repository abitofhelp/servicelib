// Copyright (c) 2025 A Bit of Help, Inc.

// Example of executing queries with different database types
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

// User represents a user in the system
type User struct {
	ID       string
	Username string
	Email    string
	Active   bool
}

// Product represents a product in the system
type Product struct {
	ID    string
	Name  string
	Price float64
}

func main() {
	// Create a context
	// Note: ctx is only used in the commented-out code below
	// that would be used in a real application
	ctx := context.Background()

	// Suppress unused variable warnings
	_ = ctx

	// Example 1: Basic Query with PostgreSQL
	fmt.Println("=== Basic Query with PostgreSQL Example ===")

	// In a real application, you would use this code:
	// postgresPool, err := db.InitPostgresPool(ctx, "postgres://postgres:password@localhost:5432/testdb", db.DefaultTimeout)
	// if err != nil {
	//     fmt.Printf("Failed to connect to PostgreSQL: %v\n", err)
	//     return
	// }
	// defer postgresPool.Close()
	//
	// // Execute a query
	// rows, err := postgresPool.Query(ctx, "SELECT id, username, email, active FROM users WHERE active = $1", true)
	// if err != nil {
	//     fmt.Printf("Failed to execute query: %v\n", err)
	//     return
	// }
	// defer rows.Close()
	//
	// // Process results
	// var users []User
	// for rows.Next() {
	//     var user User
	//     if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Active); err != nil {
	//         fmt.Printf("Failed to scan row: %v\n", err)
	//         return
	//     }
	//     users = append(users, user)
	// }
	//
	// if err := rows.Err(); err != nil {
	//     fmt.Printf("Error iterating rows: %v\n", err)
	//     return
	// }
	//
	// fmt.Printf("Found %d active users\n", len(users))

	// For the example, we'll just print what would happen
	fmt.Println("Would execute a query to find active users:")
	fmt.Println("SQL: SELECT id, username, email, active FROM users WHERE active = $1")
	fmt.Println("Parameters: true")
	fmt.Println("Would process results by scanning each row into a User struct")
	fmt.Println("Would handle any errors during query execution or result processing")

	// Example 2: Query with SQLite
	fmt.Println("\n=== Query with SQLite Example ===")

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
	// // Execute a query
	// rows, err := sqliteDB.QueryContext(ctx, "SELECT id, name, price FROM products WHERE price > ?", 100.0)
	// if err != nil {
	//     fmt.Printf("Failed to execute query: %v\n", err)
	//     return
	// }
	// defer rows.Close()
	//
	// // Process results
	// var products []Product
	// for rows.Next() {
	//     var product Product
	//     if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
	//         fmt.Printf("Failed to scan row: %v\n", err)
	//         return
	//     }
	//     products = append(products, product)
	// }
	//
	// if err := rows.Err(); err != nil {
	//     fmt.Printf("Error iterating rows: %v\n", err)
	//     return
	// }
	//
	// fmt.Printf("Found %d products with price > $100\n", len(products))

	// For the example, we'll just print what would happen
	fmt.Println("Would execute a query to find products with price > $100:")
	fmt.Println("SQL: SELECT id, name, price FROM products WHERE price > ?")
	fmt.Println("Parameters: 100.0")
	fmt.Println("Would process results by scanning each row into a Product struct")
	fmt.Println("SQLite uses ? as parameter placeholder instead of $1")

	// Example 3: Query Single Row
	fmt.Println("\n=== Query Single Row Example ===")

	// In a real application, you would use this code:
	// var user User
	// err = sqliteDB.QueryRowContext(ctx, "SELECT id, username, email, active FROM users WHERE id = ?", "user123").
	//     Scan(&user.ID, &user.Username, &user.Email, &user.Active)
	// if err != nil {
	//     if err == sql.ErrNoRows {
	//         fmt.Println("User not found")
	//     } else {
	//         fmt.Printf("Failed to query user: %v\n", err)
	//     }
	//     return
	// }
	//
	// fmt.Printf("Found user: %s (%s)\n", user.Username, user.Email)

	// For the example, we'll just print what would happen
	fmt.Println("Would execute a query to find a single user by ID:")
	fmt.Println("SQL: SELECT id, username, email, active FROM users WHERE id = ?")
	fmt.Println("Parameters: \"user123\"")
	fmt.Println("Would use QueryRowContext for queries that return a single row")
	fmt.Println("Would handle sql.ErrNoRows specifically to detect when no rows are found")

	// Example 4: Parameterized Queries
	fmt.Println("\n=== Parameterized Queries Example ===")
	fmt.Println("Always use parameterized queries to prevent SQL injection:")

	// Safe:
	fmt.Println("Safe: db.QueryContext(ctx, \"SELECT * FROM users WHERE username = ?\", userInput)")

	// Unsafe:
	fmt.Println("Unsafe: db.QueryContext(ctx, \"SELECT * FROM users WHERE username = '\" + userInput + \"'\")")

	fmt.Println("Parameterized queries ensure user input is properly escaped")

	// Example 5: Working with NULL Values
	fmt.Println("\n=== Working with NULL Values Example ===")

	// In a real application, you would use this code:
	// type UserWithNullable struct {
	//     ID       string
	//     Username string
	//     Email    sql.NullString
	//     LastLogin sql.NullTime
	// }
	//
	// var user UserWithNullable
	// err = sqliteDB.QueryRowContext(ctx, "SELECT id, username, email, last_login FROM users WHERE id = ?", "user123").
	//     Scan(&user.ID, &user.Username, &user.Email, &user.LastLogin)
	// if err != nil {
	//     fmt.Printf("Failed to query user: %v\n", err)
	//     return
	// }
	//
	// // Check if email is NULL
	// if user.Email.Valid {
	//     fmt.Printf("Email: %s\n", user.Email.String)
	// } else {
	//     fmt.Println("Email is NULL")
	// }
	//
	// // Check if last login is NULL
	// if user.LastLogin.Valid {
	//     fmt.Printf("Last login: %v\n", user.LastLogin.Time)
	// } else {
	//     fmt.Println("Last login is NULL")
	// }

	// For the example, we'll just print what would happen
	fmt.Println("Would use sql.NullString, sql.NullInt64, sql.NullFloat64, sql.NullTime, etc.")
	fmt.Println("These types handle NULL values from the database")
	fmt.Println("Each has a Valid field to check if the value is NULL")
	fmt.Println("And a field of the appropriate type to access the value if not NULL")

	// Expected output:
	// === Basic Query with PostgreSQL Example ===
	// Would execute a query to find active users:
	// SQL: SELECT id, username, email, active FROM users WHERE active = $1
	// Parameters: true
	// Would process results by scanning each row into a User struct
	// Would handle any errors during query execution or result processing
	//
	// === Query with SQLite Example ===
	// Would execute a query to find products with price > $100:
	// SQL: SELECT id, name, price FROM products WHERE price > ?
	// Parameters: 100.0
	// Would process results by scanning each row into a Product struct
	// SQLite uses ? as parameter placeholder instead of $1
	//
	// === Query Single Row Example ===
	// Would execute a query to find a single user by ID:
	// SQL: SELECT id, username, email, active FROM users WHERE id = ?
	// Parameters: "user123"
	// Would use QueryRowContext for queries that return a single row
	// Would handle sql.ErrNoRows specifically to detect when no rows are found
	//
	// === Parameterized Queries Example ===
	// Always use parameterized queries to prevent SQL injection:
	// Safe: db.QueryContext(ctx, "SELECT * FROM users WHERE username = ?", userInput)
	// Unsafe: db.QueryContext(ctx, "SELECT * FROM users WHERE username = '" + userInput + "'")
	// Parameterized queries ensure user input is properly escaped
	//
	// === Working with NULL Values Example ===
	// Would use sql.NullString, sql.NullInt64, sql.NullFloat64, sql.NullTime, etc.
	// These types handle NULL values from the database
	// Each has a Valid field to check if the value is NULL
	// And a field of the appropriate type to access the value if not NULL
}
