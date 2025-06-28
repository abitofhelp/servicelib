// Copyright (c) 2025 A Bit of Help, Inc.

// Package transaction provides utilities for managing transactions in distributed systems.
//
// This package offers tools and patterns for ensuring data consistency across multiple
// services and databases in a distributed environment. It focuses on maintaining the
// ACID properties (Atomicity, Consistency, Isolation, Durability) as much as possible
// in distributed scenarios where traditional database transactions are not feasible.
//
// The package is organized into several subpackages:
//   - saga: Implementation of the Saga pattern for distributed transactions
//
// The Saga pattern, implemented in the saga subpackage, is particularly useful for
// maintaining data consistency across multiple services without using two-phase commit.
// It works by defining a sequence of local transactions, each with a corresponding
// compensating transaction that can undo its effects if a failure occurs.
//
// Example usage of the saga pattern:
//
//	err := saga.WithTransaction(ctx, logger, func(tx *saga.Transaction) error {
//	    // Add operations and their rollbacks to the transaction
//	    tx.AddOperation(
//	        // Operation to create a user
//	        func(ctx context.Context) error {
//	            return userRepo.Create(ctx, user)
//	        },
//	        // Rollback operation to delete the user if a later operation fails
//	        func(ctx context.Context) error {
//	            return userRepo.Delete(ctx, user.ID)
//	        },
//	    )
//
//	    // Add more operations as needed
//	    return nil
//	})
//
// Future additions to this package may include:
//   - Two-phase commit implementation
//   - Outbox pattern for reliable message publishing
//   - Distributed locking mechanisms
//   - Transaction coordination services
//
// When working with distributed transactions, it's important to consider:
//   - The CAP theorem (Consistency, Availability, Partition tolerance)
//   - The trade-offs between strong consistency and eventual consistency
//   - The impact of network failures and service unavailability
//   - The performance implications of distributed transaction patterns
//
// This package aims to provide tools that help developers make informed decisions
// about these trade-offs and implement reliable transaction management in their
// distributed systems.
package transaction
