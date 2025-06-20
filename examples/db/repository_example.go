// Copyright (c) 2025 A Bit of Help, Inc.

// Example of implementing the repository pattern with the db package
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

// UserEntity represents a user in the system
type UserEntity struct {
	ID       string
	Username string
	Email    string
	Active   bool
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	FindActive(ctx context.Context) ([]User, error)
}

// SQLUserRepository implements UserRepository for SQL databases
type SQLUserRepository struct {
	// In a real application, this would be:
	// db *sql.DB
}

// NewSQLUserRepository creates a new SQLUserRepository
func NewSQLUserRepository() *SQLUserRepository {
	// In a real application, this would be:
	// return &SQLUserRepository{db: db}
	return &SQLUserRepository{}
}

// GetByID retrieves a user by ID
func (r *SQLUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	// In a real application, this would be:
	// var user User
	// err := r.db.QueryRowContext(ctx, 
	//     "SELECT id, username, email, active FROM users WHERE id = ?", 
	//     id,
	// ).Scan(&user.ID, &user.Username, &user.Email, &user.Active)
	//
	// if err != nil {
	//     if err == sql.ErrNoRows {
	//         return nil, nil // Not found
	//     }
	//     return nil, fmt.Errorf("failed to get user: %w", err)
	// }
	//
	// return &user, nil

	// For the example, we'll just print what would happen
	fmt.Printf("Would retrieve user with ID: %s\n", id)
	fmt.Println("SQL: SELECT id, username, email, active FROM users WHERE id = ?")
	fmt.Println("Would return user object if found, nil if not found, or error if query fails")

	// Return a mock user for demonstration
	return &User{
		ID:       id,
		Username: "mockuser",
		Email:    "mock@example.com",
		Active:   true,
	}, nil
}

// Create inserts a new user
func (r *SQLUserRepository) Create(ctx context.Context, user *User) error {
	// In a real application, this would be:
	// _, err := r.db.ExecContext(ctx,
	//     "INSERT INTO users (id, username, email, active) VALUES (?, ?, ?, ?)",
	//     user.ID, user.Username, user.Email, user.Active,
	// )
	// if err != nil {
	//     return fmt.Errorf("failed to create user: %w", err)
	// }
	// return nil

	// For the example, we'll just print what would happen
	fmt.Println("Would create a new user with the following data:")
	fmt.Printf("- ID: %s\n", user.ID)
	fmt.Printf("- Username: %s\n", user.Username)
	fmt.Printf("- Email: %s\n", user.Email)
	fmt.Printf("- Active: %v\n", user.Active)
	fmt.Println("SQL: INSERT INTO users (id, username, email, active) VALUES (?, ?, ?, ?)")

	return nil
}

// Update updates an existing user
func (r *SQLUserRepository) Update(ctx context.Context, user *User) error {
	// In a real application, this would be:
	// result, err := r.db.ExecContext(ctx,
	//     "UPDATE users SET username = ?, email = ?, active = ? WHERE id = ?",
	//     user.Username, user.Email, user.Active, user.ID,
	// )
	// if err != nil {
	//     return fmt.Errorf("failed to update user: %w", err)
	// }
	//
	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	//     return fmt.Errorf("failed to get rows affected: %w", err)
	// }
	//
	// if rowsAffected == 0 {
	//     return fmt.Errorf("user not found: %s", user.ID)
	// }
	//
	// return nil

	// For the example, we'll just print what would happen
	fmt.Println("Would update user with the following data:")
	fmt.Printf("- ID: %s\n", user.ID)
	fmt.Printf("- Username: %s\n", user.Username)
	fmt.Printf("- Email: %s\n", user.Email)
	fmt.Printf("- Active: %v\n", user.Active)
	fmt.Println("SQL: UPDATE users SET username = ?, email = ?, active = ? WHERE id = ?")
	fmt.Println("Would check rows affected to ensure user exists")

	return nil
}

// Delete removes a user by ID
func (r *SQLUserRepository) Delete(ctx context.Context, id string) error {
	// In a real application, this would be:
	// result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	// if err != nil {
	//     return fmt.Errorf("failed to delete user: %w", err)
	// }
	//
	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	//     return fmt.Errorf("failed to get rows affected: %w", err)
	// }
	//
	// if rowsAffected == 0 {
	//     return fmt.Errorf("user not found: %s", id)
	// }
	//
	// return nil

	// For the example, we'll just print what would happen
	fmt.Printf("Would delete user with ID: %s\n", id)
	fmt.Println("SQL: DELETE FROM users WHERE id = ?")
	fmt.Println("Would check rows affected to ensure user exists")

	return nil
}

// FindActive finds all active users
func (r *SQLUserRepository) FindActive(ctx context.Context) ([]User, error) {
	// In a real application, this would be:
	// rows, err := r.db.QueryContext(ctx, "SELECT id, username, email, active FROM users WHERE active = ?", true)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to query active users: %w", err)
	// }
	// defer rows.Close()
	//
	// var users []User
	// for rows.Next() {
	//     var user User
	//     if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Active); err != nil {
	//         return nil, fmt.Errorf("failed to scan user: %w", err)
	//     }
	//     users = append(users, user)
	// }
	//
	// if err := rows.Err(); err != nil {
	//     return nil, fmt.Errorf("error iterating rows: %w", err)
	// }
	//
	// return users, nil

	// For the example, we'll just print what would happen
	fmt.Println("Would find all active users")
	fmt.Println("SQL: SELECT id, username, email, active FROM users WHERE active = ?")
	fmt.Println("Parameters: true")
	fmt.Println("Would scan results into User structs")

	// Return mock users for demonstration
	return []User{
		{ID: "user1", Username: "alice", Email: "alice@example.com", Active: true},
		{ID: "user2", Username: "bob", Email: "bob@example.com", Active: true},
	}, nil
}

// MongoUserRepository implements UserRepository for MongoDB
type MongoUserRepository struct {
	// In a real application, this would be:
	// client *mongo.Client
	// database string
	// collection string
}

// NewMongoUserRepository creates a new MongoUserRepository
func NewMongoUserRepository() *MongoUserRepository {
	// In a real application, this would be:
	// return &MongoUserRepository{
	//     client: client,
	//     database: "myapp",
	//     collection: "users",
	// }
	return &MongoUserRepository{}
}

// GetByID retrieves a user by ID from MongoDB
func (r *MongoUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	// In a real application, this would be:
	// collection := r.client.Database(r.database).Collection(r.collection)
	// 
	// var user User
	// filter := bson.M{"_id": id}
	// err := collection.FindOne(ctx, filter).Decode(&user)
	// if err != nil {
	//     if err == mongo.ErrNoDocuments {
	//         return nil, nil // Not found
	//     }
	//     return nil, fmt.Errorf("failed to get user: %w", err)
	// }
	//
	// return &user, nil

	// For the example, we'll just print what would happen
	fmt.Printf("Would retrieve user with ID: %s from MongoDB\n", id)
	fmt.Println("Filter: {\"_id\": \"" + id + "\"}")
	fmt.Println("Would return user object if found, nil if not found, or error if query fails")

	// Return a mock user for demonstration
	return &User{
		ID:       id,
		Username: "mongomockuser",
		Email:    "mongomock@example.com",
		Active:   true,
	}, nil
}

// Other MongoDB repository methods would be implemented similarly

func main() {
	// Create a context
	// Note: ctx is only used in the commented-out code below
	// that would be used in a real application
	ctx := context.Background()

	// Suppress unused variable warnings
	_ = ctx

	// Example 1: Using SQL Repository
	fmt.Println("=== SQL Repository Example ===")

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
	// // Create repository
	// userRepo := NewSQLUserRepository(sqliteDB)

	// For the example, we'll create a repository without a real database
	userRepo := NewSQLUserRepository()

	// Use the repository
	fmt.Println("\nRetrieving user:")
	user, err := userRepo.GetByID(ctx, "user123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else if user == nil {
		fmt.Println("User not found")
	} else {
		fmt.Printf("Found user: %s (%s)\n", user.Username, user.Email)
	}

	fmt.Println("\nCreating user:")
	newUser := &User{
		ID:       "user456",
		Username: "janedoe",
		Email:    "jane.doe@example.com",
		Active:   true,
	}
	if err := userRepo.Create(ctx, newUser); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\nUpdating user:")
	newUser.Email = "jane.updated@example.com"
	if err := userRepo.Update(ctx, newUser); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\nFinding active users:")
	activeUsers, err := userRepo.FindActive(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d active users\n", len(activeUsers))
		for _, u := range activeUsers {
			fmt.Printf("- %s (%s)\n", u.Username, u.Email)
		}
	}

	fmt.Println("\nDeleting user:")
	if err := userRepo.Delete(ctx, "user456"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Example 2: Using MongoDB Repository
	fmt.Println("\n=== MongoDB Repository Example ===")

	// In a real application, you would use this code:
	// mongoClient, err := db.InitMongoClient(ctx, "mongodb://localhost:27017", db.DefaultTimeout)
	// if err != nil {
	//     fmt.Printf("Failed to connect to MongoDB: %v\n", err)
	//     return
	// }
	// defer mongoClient.Disconnect(ctx)
	//
	// // Create repository
	// mongoUserRepo := NewMongoUserRepository(mongoClient)

	// For the example, we'll create a repository without a real database
	mongoUserRepo := NewMongoUserRepository()

	// Use the repository
	fmt.Println("\nRetrieving user from MongoDB:")
	mongoUser, err := mongoUserRepo.GetByID(ctx, "user789")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else if mongoUser == nil {
		fmt.Println("User not found")
	} else {
		fmt.Printf("Found user: %s (%s)\n", mongoUser.Username, mongoUser.Email)
	}

	// Example 3: Repository Pattern Benefits
	fmt.Println("\n=== Repository Pattern Benefits ===")
	fmt.Println("1. Abstraction: Hides database implementation details")
	fmt.Println("2. Testability: Easy to mock for unit testing")
	fmt.Println("3. Flexibility: Can switch database types without changing business logic")
	fmt.Println("4. Separation of Concerns: Keeps data access logic separate from business logic")
	fmt.Println("5. Consistency: Provides a consistent interface for data access")

	// Expected output:
	// === SQL Repository Example ===
	//
	// Retrieving user:
	// Would retrieve user with ID: user123
	// SQL: SELECT id, username, email, active FROM users WHERE id = ?
	// Would return user object if found, nil if not found, or error if query fails
	// Found user: mockuser (mock@example.com)
	//
	// Creating user:
	// Would create a new user with the following data:
	// - ID: user456
	// - Username: janedoe
	// - Email: jane.doe@example.com
	// - Active: true
	// SQL: INSERT INTO users (id, username, email, active) VALUES (?, ?, ?, ?)
	//
	// Updating user:
	// Would update user with the following data:
	// - ID: user456
	// - Username: janedoe
	// - Email: jane.updated@example.com
	// - Active: true
	// SQL: UPDATE users SET username = ?, email = ?, active = ? WHERE id = ?
	// Would check rows affected to ensure user exists
	//
	// Finding active users:
	// Would find all active users
	// SQL: SELECT id, username, email, active FROM users WHERE active = ?
	// Parameters: true
	// Would scan results into User structs
	// Found 2 active users
	// - alice (alice@example.com)
	// - bob (bob@example.com)
	//
	// Deleting user:
	// Would delete user with ID: user456
	// SQL: DELETE FROM users WHERE id = ?
	// Would check rows affected to ensure user exists
	//
	// === MongoDB Repository Example ===
	//
	// Retrieving user from MongoDB:
	// Would retrieve user with ID: user789 from MongoDB
	// Filter: {"_id": "user789"}
	// Would return user object if found, nil if not found, or error if query fails
	// Found user: mongomockuser (mongomock@example.com)
	//
	// === Repository Pattern Benefits ===
	// 1. Abstraction: Hides database implementation details
	// 2. Testability: Easy to mock for unit testing
	// 3. Flexibility: Can switch database types without changing business logic
	// 4. Separation of Concerns: Keeps data access logic separate from business logic
	// 5. Consistency: Provides a consistent interface for data access
}
