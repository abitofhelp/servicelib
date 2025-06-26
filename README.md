# ServiceLib

## Overview

ServiceLib is a comprehensive Go library designed to provide a robust foundation for building scalable, maintainable, and production-ready microservices. It offers a collection of packages that address common challenges in service development, from authentication and configuration to error handling and telemetry.

## Features

- **Authentication & Authorization**: Secure your services with JWT, OIDC, and role-based access control
- **Configuration Management**: Flexible configuration with support for multiple sources and formats
- **Error Handling**: Structured error types with context, stack traces, and categorization
- **Database Access**: Connection management, health checks, and transaction support
- **Dependency Injection**: Simple yet powerful DI container for managing service dependencies
- **Telemetry**: Integrated logging, metrics, and distributed tracing
- **Health Checks**: Standardized health check endpoints and status reporting
- **Middleware**: Common HTTP middleware for logging, error handling, and more
- **Retry & Circuit Breaking**: Resilience patterns for handling transient failures
- **Validation**: Comprehensive input validation utilities

## Installation

```bash
go get github.com/abitofhelp/servicelib
```

## Quick Start

See the [Quick Start examples](./EXAMPLES/quickstart/README.md) for complete, runnable examples of how to use ServiceLib to build a basic microservice.

## Packages

ServiceLib is organized into the following packages:

- [auth](./auth/README.md) - Authentication and authorization
- [cache](./cache/README.md) - Caching utilities
- [circuit](./circuit/README.md) - Circuit breaker implementation
- [config](./config/README.md) - Configuration management
- [context](./context/README.md) - Context utilities
- [date](./date/README.md) - Date and time utilities
- [db](./db/README.md) - Database access and management
- [di](./di/README.md) - Dependency injection
- [env](./env/README.md) - Environment variable utilities
- [errors](./errors/README.md) - Error handling, management, and recovery patterns
- [graphql](./graphql/README.md) - GraphQL utilities
- [health](./health/README.md) - Health check utilities
- [logging](./logging/README.md) - Structured logging
- [middleware](./middleware/README.md) - HTTP middleware
- [model](./model/README.md) - Model utilities
- [rate](./rate/README.md) - Rate limiting
- [repository](./repository/README.md) - Repository pattern implementation
- [retry](./retry/README.md) - Retry utilities
- [shutdown](./shutdown/README.md) - Graceful shutdown utilities
- [signal](./signal/README.md) - Signal handling
- [stringutil](./stringutil/README.md) - String utilities
- [telemetry](./telemetry/README.md) - Telemetry (metrics, tracing)
- [transaction](./transaction/README.md) - Transaction management
- [validation](./validation/README.md) - Input validation
- [valueobject](./valueobject/README.md) - Value object implementations

## Examples

For complete, runnable examples of each component, see the [EXAMPLES](./EXAMPLES/README.md) directory.

## Contributing

Contributions to ServiceLib are welcome! Please see the [Contributing Guide](./CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
