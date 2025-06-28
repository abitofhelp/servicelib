// Copyright (c) 2025 A Bit of Help, Inc.

// Example of using timeouts with the context package
package main

import (
	"context" // Standard context package
	"fmt"
	"time"

	svcctx "github.com/abitofhelp/servicelib/context" // Aliased to avoid name collision
)

func main() {
	// Demonstrate different timeout functions
	fmt.Println("Context Timeout Examples")
	fmt.Println("=======================")

	// Create a base context
	baseCtx := svcctx.Background()

	// 1. Default timeout (30 seconds)
	defaultCtx, defaultCancel := svcctx.WithDefaultTimeout(baseCtx)
	defer defaultCancel()
	fmt.Println("\n1. Default Timeout Context:")
	printContextDeadline(defaultCtx)

	// 2. Short timeout (5 seconds)
	shortCtx, shortCancel := svcctx.WithShortTimeout(baseCtx)
	defer shortCancel()
	fmt.Println("\n2. Short Timeout Context:")
	printContextDeadline(shortCtx)

	// 3. Long timeout (60 seconds)
	longCtx, longCancel := svcctx.WithLongTimeout(baseCtx)
	defer longCancel()
	fmt.Println("\n3. Long Timeout Context:")
	printContextDeadline(longCtx)

	// 4. Database timeout (10 seconds)
	dbCtx, dbCancel := svcctx.WithDatabaseTimeout(baseCtx)
	defer dbCancel()
	fmt.Println("\n4. Database Timeout Context:")
	printContextDeadline(dbCtx)
	fmt.Println("Operation:", svcctx.GetOperation(dbCtx))

	// 5. Network timeout (15 seconds)
	networkCtx, networkCancel := svcctx.WithNetworkTimeout(baseCtx)
	defer networkCancel()
	fmt.Println("\n5. Network Timeout Context:")
	printContextDeadline(networkCtx)
	fmt.Println("Operation:", svcctx.GetOperation(networkCtx))

	// 6. External service timeout (20 seconds)
	externalCtx, externalCancel := svcctx.WithExternalServiceTimeout(baseCtx, "payment-service")
	defer externalCancel()
	fmt.Println("\n6. External Service Timeout Context:")
	printContextDeadline(externalCtx)
	fmt.Println("Operation:", svcctx.GetOperation(externalCtx))
	fmt.Println("Service Name:", svcctx.GetServiceName(externalCtx))

	// 7. Custom timeout
	customCtx, customCancel := svcctx.WithTimeout(baseCtx, 25*time.Second, svcctx.ContextOptions{
		Operation: "custom-operation",
	})
	defer customCancel()
	fmt.Println("\n7. Custom Timeout Context:")
	printContextDeadline(customCtx)
	fmt.Println("Operation:", svcctx.GetOperation(customCtx))

	// Demonstrate a timeout expiring
	fmt.Println("\nDemonstrating a timeout expiring:")
	quickTimeoutCtx, quickTimeoutCancel := svcctx.WithTimeout(baseCtx, 1*time.Second, svcctx.ContextOptions{
		Operation: "quick-operation",
	})
	defer quickTimeoutCancel()

	// Wait for the timeout to expire
	fmt.Println("Waiting for timeout to expire...")
	time.Sleep(1100 * time.Millisecond)

	// Check if the context is done
	if err := svcctx.CheckContext(quickTimeoutCtx); err != nil {
		fmt.Printf("Context error: %v\n", err)
	} else {
		fmt.Println("Context is still valid (unexpected)")
	}
}

// Helper function to print context deadline
func printContextDeadline(ctx context.Context) {
	deadline, ok := ctx.Deadline()
	if ok {
		fmt.Printf("Deadline: %s (in %v)\n",
			deadline.Format(time.RFC3339),
			time.Until(deadline).Round(time.Millisecond))
	} else {
		fmt.Println("No deadline set")
	}
}