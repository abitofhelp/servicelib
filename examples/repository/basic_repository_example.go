// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic repository implementation
package main

import (
    "context"
    "errors"
    "fmt"
    "sync"

    "github.com/abitofhelp/servicelib/repository"
)

// User is a domain entity
type User struct {
    ID       string
    Username string
    Email    string
    Active   bool
}

// InMemoryUserRepository implements the Repository interface for User entities
type InMemoryUserRepository struct {
    users  map[string]User
    mutex  sync.RWMutex
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{
        users: make(map[string]User),
    }
}

// GetByID retrieves a user by ID
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    user, exists := r.users[id]
    if !exists {
        return User{}, errors.New("user not found")
    }

    return user, nil
}

// GetAll retrieves all users
func (r *InMemoryUserRepository) GetAll(ctx context.Context) ([]User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    users := make([]User, 0, len(r.users))
    for _, user := range r.users {
        users = append(users, user)
    }

    return users, nil
}

// Save persists a user
func (r *InMemoryUserRepository) Save(ctx context.Context, user User) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    r.users[user.ID] = user
    return nil
}

func main() {
    // Create a repository
    var userRepo repository.Repository[User] = NewInMemoryUserRepository()

    // Create a context
    ctx := context.Background()

    // Create and save a user
    user := User{
        ID:       "user123",
        Username: "johndoe",
        Email:    "john.doe@example.com",
        Active:   true,
    }

    err := userRepo.Save(ctx, user)
    if err != nil {
        fmt.Printf("Error saving user: %v\n", err)
        return
    }

    // Retrieve the user
    retrievedUser, err := userRepo.GetByID(ctx, "user123")
    if err != nil {
        fmt.Printf("Error retrieving user: %v\n", err)
        return
    }

    fmt.Printf("Retrieved user: %+v\n", retrievedUser)

    // Get all users
    allUsers, err := userRepo.GetAll(ctx)
    if err != nil {
        fmt.Printf("Error retrieving all users: %v\n", err)
        return
    }

    fmt.Printf("All users: %+v\n", allUsers)
}
