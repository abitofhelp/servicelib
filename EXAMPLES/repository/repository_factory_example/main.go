// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using repository factory
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/abitofhelp/servicelib/repository"
)

// FactoryUser is a domain entity
type FactoryUser struct {
	ID       string
	Username string
	Email    string
}

// Order is another domain entity
type Order struct {
	ID     string
	UserID string
	Amount float64
}

// FactoryInMemoryUserRepository implements the Repository interface for FactoryUser entities
type FactoryInMemoryUserRepository struct {
	users map[string]FactoryUser
	mutex sync.RWMutex
}

// NewFactoryInMemoryUserRepository creates a new in-memory user repository
func NewFactoryInMemoryUserRepository() *FactoryInMemoryUserRepository {
	return &FactoryInMemoryUserRepository{
		users: make(map[string]FactoryUser),
	}
}

// GetByID retrieves a user by ID
func (r *FactoryInMemoryUserRepository) GetByID(ctx context.Context, id string) (FactoryUser, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return FactoryUser{}, errors.New("user not found")
	}

	return user, nil
}

// GetAll retrieves all users
func (r *FactoryInMemoryUserRepository) GetAll(ctx context.Context) ([]FactoryUser, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]FactoryUser, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

// Save persists a user
func (r *FactoryInMemoryUserRepository) Save(ctx context.Context, user FactoryUser) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.users[user.ID] = user
	return nil
}

// InMemoryOrderRepository implements the Repository interface for Order entities
type InMemoryOrderRepository struct {
	orders map[string]Order
	mutex  sync.RWMutex
}

// NewInMemoryOrderRepository creates a new in-memory order repository
func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]Order),
	}
}

// GetByID retrieves an order by ID
func (r *InMemoryOrderRepository) GetByID(ctx context.Context, id string) (Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return Order{}, errors.New("order not found")
	}

	return order, nil
}

// GetAll retrieves all orders
func (r *InMemoryOrderRepository) GetAll(ctx context.Context) ([]Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	orders := make([]Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}

	return orders, nil
}

// Save persists an order
func (r *InMemoryOrderRepository) Save(ctx context.Context, order Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.orders[order.ID] = order
	return nil
}

// RepositoryFactoryImpl implements the RepositoryFactory interface
type RepositoryFactoryImpl struct {
	userRepo  repository.Repository[FactoryUser]
	orderRepo repository.Repository[Order]
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(
	userRepo repository.Repository[FactoryUser],
	orderRepo repository.Repository[Order],
) *RepositoryFactoryImpl {
	return &RepositoryFactoryImpl{
		userRepo:  userRepo,
		orderRepo: orderRepo,
	}
}

// GetRepository returns a repository for the given entity type
func (f *RepositoryFactoryImpl) GetRepository() any {
	// This method is typically used with type assertions in the calling code
	return f
}

// GetUserRepository returns the user repository
func (f *RepositoryFactoryImpl) GetUserRepository() repository.Repository[FactoryUser] {
	return f.userRepo
}

// GetOrderRepository returns the order repository
func (f *RepositoryFactoryImpl) GetOrderRepository() repository.Repository[Order] {
	return f.orderRepo
}

func main() {
	// Create repositories
	userRepo := NewFactoryInMemoryUserRepository()
	orderRepo := NewInMemoryOrderRepository()

	// Create a repository factory
	factory := NewRepositoryFactory(userRepo, orderRepo)

	// Get repositories from the factory
	userRepository := factory.GetUserRepository()
	orderRepository := factory.GetOrderRepository()

	// Use the repositories
	ctx := context.Background()

	user := FactoryUser{ID: "user123", Username: "johndoe", Email: "john@example.com"}
	err := userRepository.Save(ctx, user)
	if err != nil {
		fmt.Printf("Error saving user: %v\n", err)
		return
	}

	order := Order{ID: "order123", UserID: "user123", Amount: 99.99}
	err = orderRepository.Save(ctx, order)
	if err != nil {
		fmt.Printf("Error saving order: %v\n", err)
		return
	}

	fmt.Println("User and order saved successfully")

	// Retrieve the user and order
	retrievedUser, err := userRepository.GetByID(ctx, "user123")
	if err != nil {
		fmt.Printf("Error retrieving user: %v\n", err)
		return
	}

	retrievedOrder, err := orderRepository.GetByID(ctx, "order123")
	if err != nil {
		fmt.Printf("Error retrieving order: %v\n", err)
		return
	}

	fmt.Printf("Retrieved user: %+v\n", retrievedUser)
	fmt.Printf("Retrieved order: %+v\n", retrievedOrder)
}
