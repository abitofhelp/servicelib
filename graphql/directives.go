// Copyright (c) 2025 A Bit of Help, Inc.

// Package graphql provides utilities for working with GraphQL.
package graphql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/abitofhelp/servicelib/auth/middleware"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/telemetry"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

var (
	// AuthorizationCheckDuration measures the duration of authorization checks
	AuthorizationCheckDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "authorization_check_duration_seconds",
		Help:    "Duration of authorization checks in seconds",
		Buckets: prometheus.DefBuckets,
	})

	// AuthorizationFailures counts the number of failed authorization attempts
	AuthorizationFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "authorization_failures_total",
		Help: "Total number of authorization failures",
	})
)

func init() {
	// Register metrics
	prometheus.MustRegister(AuthorizationCheckDuration)
	prometheus.MustRegister(AuthorizationFailures)
}

// IsAuthorizedDirective implements the @isAuthorized directive for GraphQL
// It checks if the user has any of the allowed roles and all of the required scopes for the specified resource
//
// Parameters:
//   - ctx: The context of the request
//   - obj: The object being resolved
//   - next: The next resolver in the chain
//   - allowedRoles: The roles that are allowed to access this field
//   - requiredScopes: The scopes that are required to access this field (optional)
//   - resource: The resource being accessed (optional)
//   - logger: A context logger for logging (optional, can be nil)
//
// Returns:
//   - The result of the next resolver if authorized
//   - An error if not authorized
func IsAuthorizedDirective(ctx context.Context, obj interface{}, next graphql.Resolver, allowedRoles []string, requiredScopes []string, resource string, logger *logging.ContextLogger) (interface{}, error) {
	// Start timing the authorization check
	start := time.Now()
	defer func() {
		// Record metrics about authorization checks
		AuthorizationCheckDuration.Observe(time.Since(start).Seconds())
	}()

	// Create a span for tracing
	ctx, span := telemetry.StartSpan(ctx, "graphql.IsAuthorizedDirective")
	defer span.End()

	// Add span attributes
	telemetry.AddSpanAttributes(ctx,
		attribute.StringSlice("auth.allowed_roles", allowedRoles),
		attribute.StringSlice("auth.required_scopes", requiredScopes),
		attribute.String("auth.resource", resource),
	)

	// Get the operation name from the GraphQL context
	operationName := "graphql_operation"
	if graphqlContext := graphql.GetOperationContext(ctx); graphqlContext != nil {
		operationName = graphqlContext.OperationName
		if operationName == "" {
			// If no operation name is provided, use the raw query
			operationName = graphqlContext.RawQuery
		}
	}

	// Check if the user has any of the allowed roles and all of the required scopes for the specified resource
	if !middleware.IsAuthorizedWithScopes(ctx, allowedRoles, requiredScopes, resource) {
		// Record the authorization failure
		AuthorizationFailures.Inc()
		telemetry.RecordErrorMetric(ctx, "graphql", "authorization_failure")

		// Log the unauthorized access attempt
		if logger != nil {
			logger.Warn(ctx, "Unauthorized access attempt",
				zap.String("operation", operationName),
				zap.Strings("allowed_roles", allowedRoles),
				zap.Strings("required_scopes", requiredScopes),
				zap.String("resource", resource),
			)
		}

		// Return an error
		err := errors.New("you are not authorized to perform this operation")
		telemetry.RecordErrorSpan(ctx, err)
		return nil, err
	}

	// If authorized, proceed with the next resolver
	return next(ctx)
}

// CheckAuthorization is a helper function to check if the user is authorized to perform an operation
// It can be used in resolvers to perform custom authorization checks
//
// Parameters:
//   - ctx: The context of the request
//   - allowedRoles: The roles that are allowed to perform the operation
//   - requiredScopes: The scopes that are required to perform the operation (optional)
//   - resource: The resource being accessed (optional)
//   - operation: The name of the operation being performed
//   - logger: A context logger for logging (optional, can be nil)
//
// Returns:
//   - An error if not authorized, nil if authorized
func CheckAuthorization(ctx context.Context, allowedRoles []string, requiredScopes []string, resource string, operation string, logger *logging.ContextLogger) error {
	// Start timing the authorization check
	start := time.Now()
	defer func() {
		// Record metrics about authorization checks
		AuthorizationCheckDuration.Observe(time.Since(start).Seconds())
	}()

	// Create a span for tracing
	ctx, span := telemetry.StartSpan(ctx, "graphql.CheckAuthorization")
	defer span.End()

	// Add span attributes
	telemetry.AddSpanAttributes(ctx,
		attribute.StringSlice("auth.allowed_roles", allowedRoles),
		attribute.StringSlice("auth.required_scopes", requiredScopes),
		attribute.String("auth.resource", resource),
		attribute.String("operation", operation),
	)

	// Check if the user has any of the allowed roles and all of the required scopes for the specified resource
	if !middleware.IsAuthorizedWithScopes(ctx, allowedRoles, requiredScopes, resource) {
		// Record the authorization failure
		AuthorizationFailures.Inc()
		telemetry.RecordErrorMetric(ctx, "graphql", "authorization_failure")

		// Log the unauthorized access attempt
		if logger != nil {
			logger.Warn(ctx, "Unauthorized access attempt",
				zap.String("operation", operation),
				zap.Strings("allowed_roles", allowedRoles),
				zap.Strings("required_scopes", requiredScopes),
				zap.String("resource", resource),
			)
		}

		// Return an error
		err := errors.New("you are not authorized to perform this operation")
		telemetry.RecordErrorSpan(ctx, err)
		return err
	}

	return nil
}

// ConvertRolesToStrings converts a slice of enum roles to a slice of strings
// This is useful when working with generated enum types for roles
//
// Parameters:
//   - roles: A slice of role enums that have a String() method
//
// Returns:
//   - A slice of strings representing the roles
func ConvertRolesToStrings[T fmt.Stringer](roles []T) []string {
	strRoles := make([]string, len(roles))
	for i, role := range roles {
		strRoles[i] = role.String()
	}
	return strRoles
}
