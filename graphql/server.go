// Copyright (c) 2025 A Bit of Help, Inc.

// Package graphql provides utilities for working with GraphQL.
package graphql

import (
	"context"
	"errors"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/middleware"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

// ServerConfig contains configuration for the GraphQL server
type ServerConfig struct {
	MaxQueryDepth      int
	MaxQueryComplexity int
	RequestTimeout     time.Duration
}

// NewDefaultServerConfig creates a new server configuration with default values
func NewDefaultServerConfig() ServerConfig {
	return ServerConfig{
		MaxQueryDepth:      25,
		MaxQueryComplexity: 100,
		RequestTimeout:     30 * time.Second,
	}
}

// NewServer creates a new GraphQL server with the given schema and configuration
func NewServer(schema graphql.ExecutableSchema, logger *logging.ContextLogger, cfg ServerConfig) *handler.Server {
	// Create server with configurations
	server := handler.NewDefaultServer(schema)

	// 1. Request validation
	server.Use(extension.FixedComplexityLimit(cfg.MaxQueryComplexity))
	server.Use(&extension.ComplexityLimit{
		Func: func(ctx context.Context, rc *graphql.OperationContext) int {
			return cfg.MaxQueryDepth
		},
	})

	// 2. Error handling
	server.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		logger.Error(ctx, "GraphQL panic recovered",
			zap.Any("error", err),
			zap.String("request_id", middleware.RequestID(ctx)),
		)
		return &gqlerror.Error{
			Message: "Internal server error",
			Extensions: map[string]interface{}{
				"code": "INTERNAL_ERROR",
				"time": time.Now().Format(time.RFC3339),
			},
		}
	})

	// 3. Response formatting with improved context cancellation handling
	server.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		// Get request ID for logging and error responses
		requestID := middleware.RequestID(ctx)
		timestamp := time.Now().UTC().Format(time.RFC3339)

		// Check if it's already a GraphQL error
		var gqlErr *gqlerror.Error
		if errors.As(err, &gqlErr) {
			return gqlErr
		}

		// Handle context cancellation
		if errors.Is(err, context.Canceled) {
			// Log at debug level
			logger.Debug(ctx, "Request context was canceled",
				zap.String("request_id", requestID),
				zap.Error(err),
			)

			// Return a more informative error message
			return &gqlerror.Error{
				Message: "The request was interrupted. This could be due to a client disconnect.",
				Extensions: map[string]interface{}{
					"code":       "CLIENT_DISCONNECTED",
					"timestamp":  timestamp,
					"request_id": requestID,
				},
			}
		}

		// Handle timeout
		if errors.Is(err, context.DeadlineExceeded) {
			// Log at debug level
			logger.Debug(ctx, "Request timed out",
				zap.String("request_id", requestID),
				zap.Error(err),
			)

			// Return a timeout-specific error message
			return &gqlerror.Error{
				Message: "The request timed out. Please try again with a simpler query or contact support if the issue persists.",
				Extensions: map[string]interface{}{
					"code":       "TIMEOUT",
					"timestamp":  timestamp,
					"request_id": requestID,
				},
			}
		}

		// Log detailed error internally for other types of errors
		logger.Error(ctx, "GraphQL error",
			zap.Error(err),
			zap.String("request_id", requestID),
		)

		// Return formatted error to client
		return &gqlerror.Error{
			Message: "An error occurred while processing your request",
			Extensions: map[string]interface{}{
				"code":       "INTERNAL_ERROR",
				"timestamp":  timestamp,
				"request_id": requestID,
			},
		}
	})

	// 4. Query complexity limits
	server.Use(extension.FixedComplexityLimit(cfg.MaxQueryComplexity))

	// Add request validation middleware with proper context cancellation handling
	server.AroundOperations(createAroundOperationsFunc(logger, cfg.RequestTimeout))

	return server
}

// createAroundOperationsFunc creates a function that wraps GraphQL operations with timeout and cancellation handling
func createAroundOperationsFunc(logger *logging.ContextLogger, requestTimeout time.Duration) graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		op := graphql.GetOperationContext(ctx)

		// Validate operation name
		if op.Operation.Name == "" {
			return func(ctx context.Context) *graphql.Response {
				return &graphql.Response{
					Errors: gqlerror.List{{
						Message: "Operation must have a name",
						Extensions: map[string]interface{}{
							"code": "VALIDATION_ERROR",
						},
					}},
				}
			}
		}

		// Create a new context with timeout
		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, requestTimeout)

		// Create a done channel to signal when processing is complete
		done := make(chan struct{})

		// Get the response handler using the timeout context
		responseHandler := next(timeoutCtx)

		// Return a wrapped response handler that handles context cancellation
		return func(ctx context.Context) *graphql.Response {
			// Check if the parent context was already cancelled
			if ctx.Err() != nil {
				timeoutCancel() // Cancel our timeout context
				logger.Debug(ctx, "Parent context was cancelled before processing",
					zap.String("request_id", middleware.RequestID(ctx)),
					zap.Error(ctx.Err()),
				)
				return &graphql.Response{
					Errors: gqlerror.List{{
						Message: "The request was interrupted. This could be due to a client disconnect or a server-side timeout.",
						Extensions: map[string]interface{}{
							"code":       "REQUEST_INTERRUPTED",
							"timestamp":  time.Now().UTC().Format(time.RFC3339),
							"request_id": middleware.RequestID(ctx),
						},
					}},
				}
			}

			// Process the response in a goroutine
			var resp *graphql.Response
			go func() {
				resp = responseHandler(timeoutCtx)
				close(done)
			}()

			// Wait for either completion, parent context cancellation, or timeout
			select {
			case <-done:
				// Request completed normally
				timeoutCancel() // Clean up the timeout context
				return resp

			case <-ctx.Done():
				// Parent context was cancelled (client disconnected)
				timeoutCancel() // Cancel our timeout context
				logger.Debug(ctx, "Parent context was cancelled during processing",
					zap.String("request_id", middleware.RequestID(ctx)),
					zap.Error(ctx.Err()),
				)
				return &graphql.Response{
					Errors: gqlerror.List{{
						Message: "The request was interrupted. This could be due to a client disconnect.",
						Extensions: map[string]interface{}{
							"code":       "CLIENT_DISCONNECTED",
							"timestamp":  time.Now().UTC().Format(time.RFC3339),
							"request_id": middleware.RequestID(ctx),
						},
					}},
				}

			case <-timeoutCtx.Done():
				// Our timeout context expired
				logger.Debug(ctx, "Request timed out",
					zap.String("request_id", middleware.RequestID(ctx)),
					zap.Duration("timeout", requestTimeout),
				)
				return &graphql.Response{
					Errors: gqlerror.List{{
						Message: "The request timed out. Please try again with a simpler query or contact support if the issue persists.",
						Extensions: map[string]interface{}{
							"code":       "TIMEOUT",
							"timestamp":  time.Now().UTC().Format(time.RFC3339),
							"request_id": middleware.RequestID(ctx),
						},
					}},
				}
			}
		}
	}
}
