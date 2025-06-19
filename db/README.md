# Database Package

The `db` package provides utilities for database connection management and operations in Go applications. It supports multiple database types and provides features for connection pooling, transaction management, and query execution.

## Features

- **Connection Management**:
  - Connection pooling
  - Automatic reconnection
  - Health checks

- **Supported Databases**:
  - PostgreSQL (via pgx)
  - SQLite
  - MongoDB

- **Features**:
  - Transaction management
  - Query execution with retries
  - Result mapping
  - Migrations

## Installation

```bash
go get github.com/abitofhelp/servicelib/db
```

## Usage

### Connecting to a Database

```go
package main

import (
    "context"
    "log"
    "github.com/abitofhelp/servicelib/db"
    "github.com/abitofhelp/servicelib/config"
)

func main() {
    // Load configuration
    cfg, err := config.New("config.yaml")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Create a database connection
    ctx := context.Background()
    database, err := db.New(ctx, cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.Close()

    // Configure connection pool
    database.SetMaxOpenConns(25)
    database.SetMaxIdleConns(10)
    database.SetConnMaxLifetime(5 * time.Minute)

    // Check connection
    if err := database.Ping(ctx); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    log.Println("Successfully connected to database")
}
```

### Executing Queries

```go
package main

import (
    "context"
    "log"
    "github.com/abitofhelp/servicelib/db"
)

type User struct {
    ID       string
    Username string
    Email    string
    Active   bool
}

func main() {
    // Assume database is already connected
    ctx := context.Background()
    
    // Execute a query
    rows, err := database.Query(ctx, "SELECT id, username, email, active FROM users WHERE active = $1", true)
    if err != nil {
        log.Fatalf("Failed to execute query: %v", err)
    }
    defer rows.Close()

    // Process results
    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Active); err != nil {
            log.Fatalf("Failed to scan row: %v", err)
        }
        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        log.Fatalf("Error iterating rows: %v", err)
    }

    log.Printf("Found %d active users", len(users))
}
```

### Transaction Management

```go
package main

import (
    "context"
    "log"
    "github.com/abitofhelp/servicelib/db"
)

func main() {
    // Assume database is already connected
    ctx := context.Background()
    
    // Execute a transaction
    err := database.Transaction(ctx, func(tx *sql.Tx) error {
        // Execute first query
        _, err := tx.Exec("INSERT INTO users (id, username, email, active) VALUES ($1, $2, $3, $4)",
            "user123", "johndoe", "john.doe@example.com", true)
        if err != nil {
            return err
        }

        // Execute second query
        _, err = tx.Exec("INSERT INTO user_roles (user_id, role) VALUES ($1, $2)",
            "user123", "admin")
        if err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        log.Fatalf("Transaction failed: %v", err)
    }

    log.Println("Transaction completed successfully")
}
```

### Query with Retries

```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/abitofhelp/servicelib/db"
)

func main() {
    // Assume database is already connected
    ctx := context.Background()
    
    // Configure retry options
    options := db.RetryOptions{
        MaxRetries:  3,
        InitialDelay: 100 * time.Millisecond,
        MaxDelay:     1 * time.Second,
        Multiplier:   2.0,
    }

    // Execute query with retries
    var count int
    err := db.QueryWithRetries(ctx, database, options, func(ctx context.Context, db *sql.DB) error {
        return db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
    })

    if err != nil {
        log.Fatalf("Query failed after retries: %v", err)
    }

    log.Printf("User count: %d", count)
}
```

### Using the Repository Pattern

```go
package main

import (
    "context"
    "database/sql"
    "errors"
    "log"
    "github.com/abitofhelp/servicelib/db"
)

// User model
type User struct {
    ID       string
    Username string
    Email    string
    Active   bool
}

// UserRepository interface
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}

// SQLUserRepository implementation
type SQLUserRepository struct {
    db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) UserRepository {
    return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.db.QueryRowContext(ctx, 
        "SELECT id, username, email, active FROM users WHERE id = $1", 
        id,
    ).Scan(&user.ID, &user.Username, &user.Email, &user.Active)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // Not found
        }
        return nil, err
    }

    return &user, nil
}

func (r *SQLUserRepository) Create(ctx context.Context, user *User) error {
    _, err := r.db.ExecContext(ctx,
        "INSERT INTO users (id, username, email, active) VALUES ($1, $2, $3, $4)",
        user.ID, user.Username, user.Email, user.Active,
    )
    return err
}

func (r *SQLUserRepository) Update(ctx context.Context, user *User) error {
    _, err := r.db.ExecContext(ctx,
        "UPDATE users SET username = $1, email = $2, active = $3 WHERE id = $4",
        user.Username, user.Email, user.Active, user.ID,
    )
    return err
}

func (r *SQLUserRepository) Delete(ctx context.Context, id string) error {
    _, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
    return err
}

func main() {
    // Assume database is already connected
    ctx := context.Background()
    
    // Create repository
    userRepo := NewSQLUserRepository(database)
    
    // Use repository
    user, err := userRepo.GetByID(ctx, "user123")
    if err != nil {
        log.Fatalf("Failed to get user: %v", err)
    }
    
    if user != nil {
        log.Printf("Found user: %s (%s)", user.Username, user.Email)
    } else {
        log.Println("User not found")
    }
}
```

## Configuration

Example database configuration in YAML:

```yaml
database:
  driver: "postgres"
  url: "postgres://user:password@localhost:5432/mydb?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime_seconds: 300
  retry:
    max_retries: 3
    initial_delay_ms: 100
    max_delay_ms: 1000
    multiplier: 2.0
```

## Best Practices

1. **Connection Pooling**: Configure connection pools based on your application's needs and the database's capacity.

2. **Transactions**: Use transactions for operations that need to be atomic.

3. **Prepared Statements**: Use prepared statements to prevent SQL injection and improve performance.

4. **Context**: Always pass a context to database operations to enable cancellation and timeouts.

5. **Error Handling**: Handle database errors appropriately, distinguishing between different types of errors.

6. **Repository Pattern**: Use the repository pattern to abstract database access and make testing easier.

7. **Migrations**: Use database migrations to manage schema changes.

## License

This project is licensed under the MIT License - see the LICENSE file for details.