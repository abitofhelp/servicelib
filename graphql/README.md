# GraphQL

## Overview

The GraphQL component provides utilities for working with GraphQL in Go applications. It offers a robust server implementation with built-in error handling, request validation, timeout management, and context cancellation handling. This component is designed to work seamlessly with the gqlgen library while adding additional features for production-ready GraphQL APIs.

## Features

- **Server Configuration**: Configurable query depth, complexity limits, and request timeouts
- **Error Handling**: Comprehensive error handling with appropriate error codes and messages
- **Context Awareness**: Proper handling of context cancellation and timeouts
- **Request Validation**: Validation of GraphQL operations with complexity limits
- **Logging Integration**: Seamless integration with the logging component
- **Security Features**: Protection against malicious queries with depth and complexity limits
- **Client Disconnect Handling**: Graceful handling of client disconnections

## Installation

```bash
go get github.com/abitofhelp/servicelib/graphql
```

## Quick Start

```go
package main

import (
    "net/http"
    
    "github.com/99designs/gqlgen/graphql/playground"
    "github.com/abitofhelp/servicelib/graphql"
    "github.com/abitofhelp/servicelib/logging"
    "go.uber.org/zap"
)

func main() {
    // Create a logger
    logger, _ := zap.NewProduction()
    contextLogger := logging.NewContextLogger(logger)
    
    // Create a GraphQL server with default configuration
    config := graphql.NewDefaultServerConfig()
    server := graphql.NewServer(
        GeneratedExecutableSchema(),
        contextLogger,
        config,
    )
    
    // Set up routes
    http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
    http.Handle("/query", server)
    
    // Start the server
    http.ListenAndServe(":8080", nil)
}
```

## API Documentation

### Core Types

#### ServerConfig

Configuration for the GraphQL server.

```go
type ServerConfig struct {
    MaxQueryDepth      int
    MaxQueryComplexity int
    RequestTimeout     time.Duration
}
```

### Key Methods

#### NewDefaultServerConfig

Creates a new server configuration with default values.

```go
func NewDefaultServerConfig() ServerConfig
```

#### NewServer

Creates a new GraphQL server with the given schema and configuration.

```go
func NewServer(schema graphql.ExecutableSchema, logger *logging.ContextLogger, cfg ServerConfig) *handler.Server
```

#### HandleError

Processes an error and returns an appropriate GraphQL error.

```go
func HandleError(ctx context.Context, err error, operation string, logger *logging.ContextLogger) error
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Basic Server](../EXAMPLES/graphql/basic_server/README.md) - Shows how to set up a basic GraphQL server
- [Error Handling](../EXAMPLES/graphql/error_handling/README.md) - Shows how to handle errors in GraphQL resolvers
- [Authentication](../EXAMPLES/graphql/authentication/README.md) - Shows how to implement authentication in a GraphQL API
- [File Upload](../EXAMPLES/graphql/file_upload/README.md) - Shows how to handle file uploads in a GraphQL API

## Best Practices

1. **Name Operations**: Always name your GraphQL operations for better error reporting and tracing
2. **Set Appropriate Limits**: Configure query depth and complexity limits based on your API's needs
3. **Handle Errors Properly**: Use the HandleError function to convert domain errors to GraphQL errors
4. **Use Context**: Pass context through resolvers to enable proper cancellation and timeouts
5. **Implement Validation**: Validate input data before processing to return meaningful errors

## Troubleshooting

### Common Issues

#### Query Complexity Exceeded

If you're seeing "query complexity limit exceeded" errors:
- The client is sending queries that are too complex for your server configuration
- Consider increasing the MaxQueryComplexity in your ServerConfig
- Advise clients to simplify their queries or use pagination

#### Request Timeout

If you're seeing timeout errors:
- The GraphQL operation is taking longer than the configured RequestTimeout
- Consider increasing the timeout for complex operations
- Optimize your resolvers and database queries
- Implement pagination for large result sets

## Related Components

- [Errors](../errors/README.md) - Error types used by the GraphQL error handler
- [Logging](../logging/README.md) - Logging integration for GraphQL operations
- [Middleware](../middleware/README.md) - Middleware for adding request ID and other context values
- [Context](../context/README.md) - Context utilities for timeout and cancellation handling

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.