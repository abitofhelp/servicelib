# Repository Package

The `repository` package provides generic repository interfaces for entity persistence operations in Go applications. It implements the Repository Pattern, which is a key component of Domain-Driven Design (DDD) and Hexagonal Architecture.

## Features

- **Generic Interfaces**: Type-safe repository interfaces using Go generics
- **Hexagonal Architecture**: Supports ports and adapters pattern
- **Domain-Driven Design**: Facilitates separation of domain and infrastructure concerns
- **Repository Pattern**: Standardized approach to data access
- **Repository Factory**: Interface for creating repositories

## Installation

```bash
go get github.com/abitofhelp/servicelib/repository
```

## Usage

### Basic Repository Implementation

```go
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
```

### Using Repository Factory

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/abitofhelp/servicelib/repository"
)

// User is a domain entity
type User struct {
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

// RepositoryFactoryImpl implements the RepositoryFactory interface
type RepositoryFactoryImpl struct {
    userRepo  repository.Repository[User]
    orderRepo repository.Repository[Order]
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(
    userRepo repository.Repository[User],
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
func (f *RepositoryFactoryImpl) GetUserRepository() repository.Repository[User] {
    return f.userRepo
}

// GetOrderRepository returns the order repository
func (f *RepositoryFactoryImpl) GetOrderRepository() repository.Repository[Order] {
    return f.orderRepo
}

func main() {
    // Create repositories
    userRepo := NewInMemoryUserRepository()
    orderRepo := NewInMemoryOrderRepository()
    
    // Create a repository factory
    factory := NewRepositoryFactory(userRepo, orderRepo)
    
    // Get repositories from the factory
    userRepository := factory.GetUserRepository()
    orderRepository := factory.GetOrderRepository()
    
    // Use the repositories
    ctx := context.Background()
    
    user := User{ID: "user123", Username: "johndoe", Email: "john@example.com"}
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
}

// InMemoryUserRepository and InMemoryOrderRepository implementations would be similar to the previous example
```

### Integration with Dependency Injection

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/abitofhelp/servicelib/di"
    "github.com/abitofhelp/servicelib/repository"
)

// User is a domain entity
type User struct {
    ID       string
    Username string
    Email    string
}

// UserService is a domain service that uses a repository
type UserService struct {
    userRepo repository.Repository[User]
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.Repository[User]) *UserService {
    return &UserService{
        userRepo: userRepo,
    }
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (User, error) {
    return s.userRepo.GetByID(ctx, id)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, username, email string) (User, error) {
    // In a real application, you would generate a unique ID
    user := User{
        ID:       "user123",
        Username: username,
        Email:    email,
    }
    
    err := s.userRepo.Save(ctx, user)
    if err != nil {
        return User{}, err
    }
    
    return user, nil
}

func main() {
    // Create a DI container
    container := di.NewContainer()
    
    // Register the user repository
    container.Register("userRepository", func(c di.Container) (interface{}, error) {
        return NewInMemoryUserRepository(), nil
    })
    
    // Register the user service with a dependency on the user repository
    container.Register("userService", func(c di.Container) (interface{}, error) {
        repo, err := c.Get("userRepository")
        if err != nil {
            return nil, err
        }
        
        return NewUserService(repo.(repository.Repository[User])), nil
    })
    
    // Resolve the user service
    service, err := container.Get("userService")
    if err != nil {
        fmt.Printf("Error resolving user service: %v\n", err)
        return
    }
    
    // Use the user service
    userService := service.(*UserService)
    ctx := context.Background()
    
    user, err := userService.CreateUser(ctx, "johndoe", "john.doe@example.com")
    if err != nil {
        fmt.Printf("Error creating user: %v\n", err)
        return
    }
    
    fmt.Printf("Created user: %+v\n", user)
}
```

## Best Practices

1. **Interface Segregation**: Keep repository interfaces focused on specific entity types.

2. **Dependency Inversion**: Depend on repository interfaces, not concrete implementations.

3. **Testability**: Use in-memory repository implementations for testing.

4. **Transaction Management**: Consider adding transaction support for operations that span multiple repositories.

5. **Error Handling**: Use domain-specific errors for repository operations.

6. **Context Usage**: Always pass a context to repository methods for cancellation and timeout support.

7. **Repository Factory**: Use a factory to create and manage repositories when dealing with multiple entity types.

8. **Concurrency**: Ensure thread safety in repository implementations.

## License

This project is licensed under the MIT License - see the LICENSE file for details.