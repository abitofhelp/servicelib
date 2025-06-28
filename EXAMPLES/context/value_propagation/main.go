//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of value propagation with the context package
package main

import (
	"context" // Standard context package
	"fmt"

	svcctx "github.com/abitofhelp/servicelib/context" // Aliased to avoid name collision
)

func main() {
	fmt.Println("Context Value Propagation Examples")
	fmt.Println("=================================")

	// Create a base context
	baseCtx, cancel := svcctx.NewContext(svcctx.ContextOptions{})
	defer cancel()

	// 1. Adding individual values to context
	fmt.Println("\n1. Adding Individual Values:")

	// Add user ID
	userCtx := svcctx.WithUserID(baseCtx, "user-123")
	fmt.Println("User ID:", svcctx.GetUserID(userCtx))

	// Add tenant ID
	tenantCtx := svcctx.WithTenantID(userCtx, "tenant-456")
	fmt.Println("Tenant ID:", svcctx.GetTenantID(tenantCtx))

	// Add operation
	opCtx := svcctx.WithOperation(tenantCtx, "example-operation")
	fmt.Println("Operation:", svcctx.GetOperation(opCtx))

	// Add service name
	serviceCtx := svcctx.WithServiceName(opCtx, "example-service")
	fmt.Println("Service Name:", svcctx.GetServiceName(serviceCtx))

	// Add correlation ID
	corrCtx := svcctx.WithCorrelationID(serviceCtx, "corr-789")
	fmt.Println("Correlation ID:", svcctx.GetCorrelationID(corrCtx))

	// 2. Adding multiple values at once
	fmt.Println("\n2. Adding Multiple Values at Once:")
	multiCtx, multiCancel := svcctx.NewContext(svcctx.ContextOptions{
		UserID:        "multi-user-123",
		TenantID:      "multi-tenant-456",
		Operation:     "multi-operation",
		ServiceName:   "multi-service",
		CorrelationID: "multi-corr-789",
	})
	defer multiCancel()
	fmt.Println(svcctx.ContextInfo(multiCtx))

	// 3. Propagating values through function calls
	fmt.Println("\n3. Propagating Values Through Function Calls:")
	processRequest(corrCtx)

	// 4. Using WithValues for arbitrary key-value pairs
	fmt.Println("\n4. Using WithValues for Arbitrary Key-Value Pairs:")
	type customKey string
	const (
		appVersionKey  customKey = "app_version"
		featureFlagKey customKey = "feature_flag"
	)

	customCtx := svcctx.WithValues(baseCtx,
		appVersionKey, "1.2.3",
		featureFlagKey, true,
	)

	// Retrieve the values
	if appVersion, ok := customCtx.Value(appVersionKey).(string); ok {
		fmt.Println("App Version:", appVersion)
	}

	if featureFlag, ok := customCtx.Value(featureFlagKey).(bool); ok {
		fmt.Println("Feature Flag:", featureFlag)
	}

	// 5. Automatic generation of IDs
	fmt.Println("\n5. Automatic Generation of IDs:")
	autoCtx, autoCancel := svcctx.NewContext(svcctx.ContextOptions{})
	defer autoCancel()

	fmt.Println("Auto-generated Request ID:", svcctx.GetRequestID(autoCtx))
	fmt.Println("Auto-generated Trace ID:", svcctx.GetTraceID(autoCtx))
	fmt.Println("Auto-generated Correlation ID:", svcctx.GetCorrelationID(autoCtx))
}

// processRequest demonstrates how context values are propagated through function calls
func processRequest(ctx context.Context) {
	fmt.Println("In processRequest function:")
	fmt.Println("- User ID:", svcctx.GetUserID(ctx))
	fmt.Println("- Tenant ID:", svcctx.GetTenantID(ctx))
	fmt.Println("- Operation:", svcctx.GetOperation(ctx))
	fmt.Println("- Service Name:", svcctx.GetServiceName(ctx))
	fmt.Println("- Correlation ID:", svcctx.GetCorrelationID(ctx))

	// Call another function to demonstrate further propagation
	validateRequest(ctx)
}

// validateRequest demonstrates further propagation of context values
func validateRequest(ctx context.Context) {
	fmt.Println("In validateRequest function:")
	fmt.Println("- User ID:", svcctx.GetUserID(ctx))
	fmt.Println("- Operation:", svcctx.GetOperation(ctx))
}