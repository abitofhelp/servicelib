// Copyright (c) 2025 A Bit of Help, Inc.

// Example of integration with dependency injection
package example_repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/abitofhelp/servicelib/repository"
)

// DIUser is a domain entity
type DIUser struct {
	ID       string
	Username string
	Email    string
}

// DIInMemoryUserRepository implements the Repository interface for DIUser entities
type DIInMemoryUserRepository struct {
	users map[string]DIUser
	mutex sync.RWMutex
}

// NewDIInMemoryUserRepository creates a new in-memory user repository
func NewDIInMemoryUserRepository() *DIInMemoryUserRepository {
	return &DIInMemoryUserRepository{
		users: make(map[string]DIUser),
	}
}

// GetByID retrieves a user by ID
func (r *DIInMemoryUserRepository) GetByID(ctx context.Context, id string) (DIUser, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return DIUser{}, errors.New("user not found")
	}

	return user, nil
}

// GetAll retrieves all users
func (r *DIInMemoryUserRepository) GetAll(ctx context.Context) ([]DIUser, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]DIUser, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

// Save persists a user
func (r *DIInMemoryUserRepository) Save(ctx context.Context, user DIUser) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.users[user.ID] = user
	return nil
}

// UserService is a domain service that uses a repository
type UserService struct {
	userRepo repository.Repository[DIUser]
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.Repository[DIUser]) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (DIUser, error) {
	return s.userRepo.GetByID(ctx, id)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, username, email string) (DIUser, error) {
	// In a real application, you would generate a unique ID
	user := DIUser{
		ID:       "user123",
		Username: username,
		Email:    email,
	}

	err := s.userRepo.Save(ctx, user)
	if err != nil {
		return DIUser{}, err
	}

	return user, nil
}

func main() {
	// Create a simple dependency injection setup

	// Create a repository
	userRepo := NewDIInMemoryUserRepository()

	// Create a service with the repository dependency
	userService := NewUserService(userRepo)
	ctx := context.Background()

	user, err := userService.CreateUser(ctx, "johndoe", "john.doe@example.com")
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return
	}

	fmt.Printf("Created user: %+v\n", user)

	// Retrieve the user
	retrievedUser, err := userService.GetUser(ctx, "user123")
	if err != nil {
		fmt.Printf("Error retrieving user: %v\n", err)
		return
	}

	fmt.Printf("Retrieved user: %+v\n", retrievedUser)
}
