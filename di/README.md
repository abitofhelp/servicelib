# Dependency Injection Package

The `di` package provides a container-based dependency injection system for Go applications. It helps manage dependencies between components, making your code more modular, testable, and maintainable.

## Features

- **Container Types**:
  - Base container
  - Service container
  - Repository container
  - Generic container

- **Features**:
  - Constructor injection
  - Singleton instances
  - Lazy initialization
  - Scoped instances
  - Circular dependency detection

## Installation

```bash
go get github.com/abitofhelp/servicelib/di
```

## Usage

### Basic Container Usage

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/di"
)

func main() {
    // Create a new DI container
    container := di.NewContainer()

    // Register a simple value
    container.Register("greeting", "Hello, World!")

    // Register a function that returns a value
    container.Register("counter", func() int {
        return 42
    })

    // Retrieve values from the container
    greeting, err := container.Get("greeting")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    counter, err := container.Get("counter")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Println(greeting.(string))  // Output: Hello, World!
    fmt.Println(counter.(int))      // Output: 42
}
```

### Dependency Resolution

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/di"
)

// Define some interfaces and implementations
type Logger interface {
    Log(message string)
}

type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {
    fmt.Println("LOG:", message)
}

type Service interface {
    DoSomething()
}

type MyService struct {
    logger Logger
}

func NewMyService(logger Logger) Service {
    return &MyService{logger: logger}
}

func (s *MyService) DoSomething() {
    s.logger.Log("Doing something...")
}

func main() {
    // Create a new DI container
    container := di.NewContainer()

    // Register the logger
    container.Register("logger", func(c di.Container) (interface{}, error) {
        return &ConsoleLogger{}, nil
    })

    // Register the service with a dependency on the logger
    container.Register("service", func(c di.Container) (interface{}, error) {
        // Resolve the logger dependency
        logger, err := c.Get("logger")
        if err != nil {
            return nil, err
        }

        // Create and return the service
        return NewMyService(logger.(Logger)), nil
    })

    // Resolve and use the service
    service, err := container.Get("service")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Use the service
    service.(Service).DoSomething()  // Output: LOG: Doing something...
}
```

### Singleton vs. Transient Instances

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/di"
)

type Counter struct {
    value int
}

func (c *Counter) Increment() {
    c.value++
}

func (c *Counter) Value() int {
    return c.value
}

func main() {
    // Create a new DI container
    container := di.NewContainer()

    // Register a singleton counter
    container.RegisterSingleton("singletonCounter", func(c di.Container) (interface{}, error) {
        return &Counter{}, nil
    })

    // Register a transient counter
    container.Register("transientCounter", func(c di.Container) (interface{}, error) {
        return &Counter{}, nil
    })

    // Get the singleton counter twice and increment it
    counter1, _ := container.Get("singletonCounter")
    counter1.(*Counter).Increment()

    counter2, _ := container.Get("singletonCounter")
    fmt.Printf("Singleton counter value: %d\n", counter2.(*Counter).Value())  // Output: 1 (shared instance)

    // Get the transient counter twice and increment it
    counter3, _ := container.Get("transientCounter")
    counter3.(*Counter).Increment()

    counter4, _ := container.Get("transientCounter")
    fmt.Printf("Transient counter value: %d\n", counter4.(*Counter).Value())  // Output: 0 (new instance)
}
```

### Scoped Containers

```go
package main

import (
    "context"
    "fmt"
    "github.com/abitofhelp/servicelib/di"
)

type RequestContext struct {
    UserID string
}

type UserService struct {
    requestContext *RequestContext
}

func (s *UserService) GetUserID() string {
    return s.requestContext.UserID
}

func main() {
    // Create a parent container
    parentContainer := di.NewContainer()

    // Handle a request for user1
    handleRequest(parentContainer, "user1")

    // Handle a request for user2
    handleRequest(parentContainer, "user2")
}

func handleRequest(parentContainer di.Container, userID string) {
    // Create a scoped container for this request
    scopedContainer := parentContainer.CreateScope()

    // Register request-specific dependencies
    scopedContainer.Register("requestContext", &RequestContext{UserID: userID})

    // Register a service that depends on the request context
    scopedContainer.Register("userService", func(c di.Container) (interface{}, error) {
        ctx, err := c.Get("requestContext")
        if err != nil {
            return nil, err
        }
        return &UserService{requestContext: ctx.(*RequestContext)}, nil
    })

    // Resolve and use the service
    service, _ := scopedContainer.Get("userService")
    userService := service.(*UserService)

    fmt.Printf("Request for user: %s\n", userService.GetUserID())
}
```

### Using with Context

```go
package main

import (
    "context"
    "fmt"
    "github.com/abitofhelp/servicelib/di"
)

func main() {
    // Create a base context
    ctx := context.Background()

    // Create a container with context
    container, err := di.NewContainerWithContext(ctx)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Register a service that uses the context
    container.Register("contextAwareService", func(c di.Container) (interface{}, error) {
        // Get the context from the container
        ctx := c.GetContext()
        
        // Use the context
        return &ContextAwareService{ctx: ctx}, nil
    })

    // Use the service
    service, _ := container.Get("contextAwareService")
    service.(*ContextAwareService).DoSomething()
}

type ContextAwareService struct {
    ctx context.Context
}

func (s *ContextAwareService) DoSomething() {
    // Use the context for cancellation, timeouts, etc.
    select {
    case <-s.ctx.Done():
        fmt.Println("Context cancelled")
    default:
        fmt.Println("Context is still valid")
    }
}
```

## Advanced Usage

### Circular Dependency Detection

The DI container automatically detects circular dependencies:

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/di"
)

func main() {
    container := di.NewContainer()

    // Register service A which depends on service B
    container.Register("serviceA", func(c di.Container) (interface{}, error) {
        b, err := c.Get("serviceB")
        if err != nil {
            return nil, err
        }
        return &ServiceA{b: b.(*ServiceB)}, nil
    })

    // Register service B which depends on service A
    container.Register("serviceB", func(c di.Container) (interface{}, error) {
        a, err := c.Get("serviceA")
        if err != nil {
            return nil, err
        }
        return &ServiceB{a: a.(*ServiceA)}, nil
    })

    // This will result in a circular dependency error
    _, err := container.Get("serviceA")
    fmt.Printf("Error: %v\n", err)
    // Output: Error: circular dependency detected: serviceA -> serviceB -> serviceA
}

type ServiceA struct {
    b *ServiceB
}

type ServiceB struct {
    a *ServiceA
}
```

### Lazy Initialization

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/servicelib/di"
)

func main() {
    container := di.NewContainer()

    // Register a service that will be initialized lazily
    container.RegisterLazy("expensiveService", func(c di.Container) (interface{}, error) {
        fmt.Println("Initializing expensive service...")
        // Simulate expensive initialization
        return &ExpensiveService{}, nil
    })

    fmt.Println("Container created, but expensive service not yet initialized")

    // Service is only initialized when requested
    service, _ := container.Get("expensiveService")
    fmt.Println("Service retrieved:", service != nil)
}

type ExpensiveService struct{}
```

## Best Practices

1. **Use Interfaces**: Register interfaces rather than concrete types to make your code more flexible and testable.

2. **Singleton vs. Transient**: Use singletons for stateless services and transient instances for stateful ones.

3. **Scoped Containers**: Use scoped containers for request-specific dependencies.

4. **Avoid Service Locator Pattern**: Inject dependencies directly rather than passing the container around.

5. **Constructor Injection**: Prefer constructor injection over property or method injection.

6. **Circular Dependencies**: Avoid circular dependencies by redesigning your components.

7. **Testing**: Use the DI container to easily mock dependencies in tests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.