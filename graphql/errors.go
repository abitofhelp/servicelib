// Copyright (c) 2025 A Bit of Help, Inc.

// Package graphql provides utilities for working with GraphQL.
package graphql

import (
	"context"
	"errors"

	myerrors "github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

// HandleError processes an error and returns an appropriate GraphQL error.
// It logs the error and converts it to a GraphQL error with appropriate extensions.
// Parameters:
//   - ctx: The context containing trace information
//   - err: The error to handle
//   - operation: The name of the operation that caused the error
//   - logger: The logger to use for logging the error
//
// Returns:
//   - error: A GraphQL error with appropriate extensions
func HandleError(ctx context.Context, err error, operation string, logger *logging.ContextLogger) error {
	// Check for context cancellation
	// Client cancellations are normal behavior for GraphQL clients like Apollo Client,
	// which may cancel in-flight requests when new requests are made for the same operation.
	// We log these at DEBUG level since they're not actual errors.
	if errors.Is(err, context.Canceled) {
		logger.Debug(ctx, "Operation canceled by client",
			zap.String("operation", operation))
		return &gqlerror.Error{
			Message: "Operation canceled",
			Extensions: map[string]interface{}{
				"code": "CANCELED",
			},
		}
	}

	if errors.Is(err, context.DeadlineExceeded) {
		logger.Error(ctx, "Operation timed out",
			zap.String("operation", operation),
			zap.Error(err))
		return &gqlerror.Error{
			Message: "Operation timed out",
			Extensions: map[string]interface{}{
				"code": "TIMEOUT",
			},
		}
	}

	// Log the error with context
	logger.Error(ctx, "GraphQL error",
		zap.String("operation", operation),
		zap.Error(err))

	// Convert to appropriate GraphQL error based on type
	switch e := err.(type) {
	case *myerrors.ValidationError, *myerrors.ValidationErrors:
		return &gqlerror.Error{
			Message: e.Error(),
			Extensions: map[string]interface{}{
				"code": "VALIDATION_ERROR",
			},
		}
	case *myerrors.NotFoundError:
		return &gqlerror.Error{
			Message: e.Error(),
			Extensions: map[string]interface{}{
				"code": "NOT_FOUND",
			},
		}
	case *myerrors.DomainError:
		return &gqlerror.Error{
			Message: e.Error(),
			Extensions: map[string]interface{}{
				"code": e.Code,
			},
		}
	case *myerrors.ApplicationError:
		return &gqlerror.Error{
			Message: e.Error(),
			Extensions: map[string]interface{}{
				"code": e.Code,
			},
		}
	case *myerrors.RepositoryError:
		// Don't expose internal database errors to clients
		return &gqlerror.Error{
			Message: "An internal error occurred",
			Extensions: map[string]interface{}{
				"code": "INTERNAL_ERROR",
			},
		}
	default:
		// Generic error handling
		return &gqlerror.Error{
			Message: "An unexpected error occurred",
			Extensions: map[string]interface{}{
				"code": "INTERNAL_ERROR",
			},
		}
	}
}
