// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic usage of the context package
package example_context

import (
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/context"
)

func main() {
	// Create a new context with default options
	ctx, cancel := context.NewContext(context.ContextOptions{})
	defer cancel() // Always call cancel to release resources

	// The context automatically generates a request ID, trace ID, and correlation ID
	fmt.Println("Basic context information:")
	fmt.Println("Request ID:", context.GetRequestID(ctx))
	fmt.Println("Trace ID:", context.GetTraceID(ctx))
	fmt.Println("Correlation ID:", context.GetCorrelationID(ctx))

	// Create a context with specific options
	customCtx, customCancel := context.NewContext(context.ContextOptions{
		Timeout:     10 * time.Second,
		Operation:   "example-operation",
		UserID:      "user-123",
		TenantID:    "tenant-456",
		ServiceName: "example-service",
	})
	defer customCancel()

	// Print all context information
	fmt.Println("\nCustom context information:")
	fmt.Println(context.ContextInfo(customCtx))

	// Check if the context is still valid
	if err := context.CheckContext(customCtx); err != nil {
		fmt.Printf("Context error: %v\n", err)
	} else {
		fmt.Println("\nContext is valid")
	}

	// Demonstrate using the Background and TODO contexts
	bgCtx := context.Background()
	todoCtx := context.TODO()

	fmt.Println("\nBackground context:", context.ContextInfo(bgCtx))
	fmt.Println("TODO context:", context.ContextInfo(todoCtx))

	// Expected output:
	// Basic context information:
	// Request ID: <some-uuid>
	// Trace ID: <some-uuid>
	// Correlation ID: <some-uuid>
	//
	// Custom context information:
	// RequestID: <some-uuid>, TraceID: <some-uuid>, UserID: user-123, TenantID: tenant-456, Operation: example-operation, CorrelationID: <some-uuid>, ServiceName: example-service, Deadline: <timestamp>
	//
	// Context is valid
	//
	// Background context: RequestID: , TraceID:
	// TODO context: RequestID: , TraceID:
}
