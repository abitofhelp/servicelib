# auth quickstart Example

## Overview

This example demonstrates how to quickly set up authentication and authorization using the ServiceLib auth package. It shows the basic steps to create an auth instance, use auth middleware, and perform authorization checks.

## Features

- **Auth Configuration**: Set up auth configuration with JWT secret key
- **Auth Middleware**: Protect HTTP endpoints with auth middleware
- **Authorization Checks**: Check if users are authorized to perform operations
- **User Context**: Access user information from the request context

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Auth Configuration and Instance Creation

This example starts by creating an auth configuration with a JWT secret key and then creating an auth instance:

```go
// Create a configuration
config := auth.DefaultConfig()
config.JWT.SecretKey = "your-secret-key-that-is-at-least-32-characters-long"

// Create an auth instance
authInstance, err := auth.New(ctx, config, logger)
if err != nil {
    logger.Fatal("Failed to create auth instance", zap.Error(err))
}
```

### Auth Middleware

The example shows how to use the auth middleware to protect an HTTP endpoint:

```go
http.Handle("/", authInstance.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Handler code
})))
```

### Authorization Checks

The example demonstrates how to check if a user is authorized to perform an operation:

```go
// Check if the user is authorized to perform an operation
authorized, err := authInstance.IsAuthorized(r.Context(), "read:resource")
if err != nil {
    http.Error(w, "Authorization error", http.StatusInternalServerError)
    return
}

if !authorized {
    http.Error(w, "Forbidden", http.StatusForbidden)
    return
}
```

## Expected Output

When you run the example and access the HTTP endpoint with a valid JWT token that has the "read:resource" permission, you should see:

```
Hello, [user-id]
```

Where [user-id] is the ID of the authenticated user. If the token is invalid or the user doesn't have the required permission, you'll see an error message.

## Related Examples


- [auth_instance](../auth_instance/README.md) - Related example for auth_instance
- [authorization](../authorization/README.md) - Related example for authorization
- [configuration](../configuration/README.md) - Related example for configuration

## Related Components

- [auth Package](../../../auth/README.md) - The auth package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
