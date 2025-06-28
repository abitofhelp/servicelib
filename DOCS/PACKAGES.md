# Service Library Package Organization

This document provides an overview of how our packages are organized logically, making it easier to understand their purposes and relationships. While our directory structure is flat (which follows Go best practices), packages can be categorized into the following logical groups:

## 1. Core Infrastructure
These packages form the foundation of our service architecture:

- `config`: Configuration management and loading
- `env`: Environment variable handling and validation
- `di`: Dependency injection container and utilities
- `context`: Context management and propagation
- `errors`: Error types, handling, and custom error definitions
- `logging`: Structured logging and log level management
- `shutdown`: Graceful shutdown coordination
- `signal`: OS signal handling and management

**When to use**: These packages are fundamental building blocks. You'll likely need them when setting up a new service or managing service lifecycle.

## 2. Security & Access Control
Packages handling authentication, authorization, and validation:

- `auth`: Authentication and authorization mechanisms
- `validation`: Input validation and sanitization

**When to use**: Use these when implementing security features or validating input data.

## 3. Data Layer
Packages managing data operations and persistence:

- `db`: Database connections and core operations
- `repository`: Data access patterns and implementations
- `transaction`: Transaction management
- `cache`: Caching mechanisms
- `model`: Core data models
- `valueobject`: Value object implementations

**When to use**: These packages are essential when working with data persistence, retrieval, or data modeling.

## 4. Communication & API
Packages handling external communication and API concerns:

- `graphql`: GraphQL schema and resolvers
- `middleware`: HTTP middleware components
- `telemetry`: Metrics, tracing, and monitoring

**When to use**: Use these when building APIs or implementing observability.

## 5. Resilience & Performance
Packages implementing reliability patterns:

- `circuit`: Circuit breaker pattern implementation
- `rate`: Rate limiting functionality
- `retry`: Retry mechanisms for failed operations

**When to use**: These packages are crucial when building resilient services that need to handle failures gracefully.

## 6. Utilities
General-purpose utility packages:

- `date`: Date/time handling utilities
- `stringutil`: String manipulation utilities

**When to use**: These provide common helper functions used across the codebase.

## Best Practices

1. **Package Dependencies**
   - Prefer depending on packages in the same group or lower-level groups
   - Avoid circular dependencies between packages
   - Core Infrastructure packages should have minimal dependencies on other groups

2. **When Creating New Code**
   - First, check if your functionality fits into an existing package
   - If creating a new package, consider which logical group it belongs to
   - Document the package's purpose in its `doc.go` file

3. **Import Guidelines**
   - Import packages directly (e.g., `github.com/abitofhelp/servicelib/config`)
   - Don't try to import by logical groups as they're just for documentation

## Common Use Cases

### Setting Up a New Service
```go
import (
    "github.com/abitofhelp/servicelib/config"
    "github.com/abitofhelp/servicelib/di"
    "github.com/abitofhelp/servicelib/logging"
)
```

### Implementing Data Access
```go
import (
    "github.com/abitofhelp/servicelib/db"
    "github.com/abitofhelp/servicelib/repository"
    "github.com/abitofhelp/servicelib/transaction"
)
```

### Adding API Endpoints
```go
import (
    "github.com/abitofhelp/servicelib/auth"
    "github.com/abitofhelp/servicelib/middleware"
    "github.com/abitofhelp/servicelib/validation"
)
```

## Contributing

When contributing new packages:
1. Review existing packages to avoid duplication
2. Consider which logical group the package belongs to
3. Update this document to include the new package
4. Provide comprehensive documentation in the package's `doc.go` file

## Package Stability

- 🟢 Stable: Production-ready, API unlikely to change
- 🟡 Beta: API might change but suitable for production use
- 🔴 Alpha: API likely to change, use with caution

[Package stability status table omitted for brevity - can be added if needed]

---

Remember: This logical grouping is for documentation purposes only. The actual package imports should always reference the direct package path.
