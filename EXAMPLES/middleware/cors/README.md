# middleware cors Example

## Overview

This example demonstrates how to use the CORS (Cross-Origin Resource Sharing) middleware from the ServiceLib middleware package to enable cross-origin requests to your API endpoints.

## Features

- **CORS Middleware**: Apply CORS headers to API responses
- **Request Context**: Add request context to handlers
- **API Endpoint**: Create a simple JSON API endpoint
- **Preflight Requests**: Handle OPTIONS preflight requests automatically

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### API Handler

Create a simple API handler that returns a JSON response:

```go
apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Set content type
    w.Header().Set("Content-Type", "application/json")

    // Return a simple JSON response
    w.Write([]byte(`{"message":"This is a CORS-enabled API endpoint"}`))
})
```

### CORS Middleware

Apply the CORS middleware to the API handler:

```go
// Apply CORS middleware to the handler
corsHandler := middleware.WithCORS(apiHandler)
```

### Request Context

Register the handler with request context:

```go
// Register the handler with request context
http.Handle("/api", middleware.WithRequestContext(corsHandler))
```

## Expected Output

When you run the example, you'll see instructions for testing the CORS functionality:

```
Server starting on :8080
API endpoint: http://localhost:8080/api

To test CORS, you can use curl:
  curl -H "Origin: http://example.com" -v http://localhost:8080/api

Or test a preflight request:
  curl -H "Origin: http://example.com" -H "Access-Control-Request-Method: POST" -H "Access-Control-Request-Headers: Content-Type" -X OPTIONS -v http://localhost:8080/api
```

## Related Examples


- [basic_usage](../basic_usage/README.md) - Related example for basic_usage
- [error_handling](../error_handling/README.md) - Related example for error_handling
- [logging](../logging/README.md) - Related example for logging

## Related Components

- [middleware Package](../../../middleware/README.md) - The middleware package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
