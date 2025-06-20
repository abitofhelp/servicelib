// Copyright (c) 2025 A Bit of Help, Inc.

// Example of idempotent operations in transactions
package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/abitofhelp/servicelib/transaction/saga"
	"go.uber.org/zap"
)

// Simple in-memory store to simulate a database
type ResourceStore struct {
	resources map[string]string
	mutex     sync.RWMutex
}

func NewResourceStore() *ResourceStore {
	return &ResourceStore{
		resources: make(map[string]string),
	}
}

func (s *ResourceStore) Create(id, value string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if resource already exists (idempotent check)
	if _, exists := s.resources[id]; exists {
		// Resource already exists, return true without changing anything
		fmt.Printf("Resource %s already exists, skipping creation\n", id)
		return true
	}

	// Create the resource
	s.resources[id] = value
	fmt.Printf("Resource %s created with value %s\n", id, value)
	return true
}

func (s *ResourceStore) Delete(id string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if resource exists
	if _, exists := s.resources[id]; !exists {
		// Resource doesn't exist, return true without doing anything (idempotent)
		fmt.Printf("Resource %s doesn't exist, skipping deletion\n", id)
		return true
	}

	// Delete the resource
	delete(s.resources, id)
	fmt.Printf("Resource %s deleted\n", id)
	return true
}

func main() {
	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Sync()

	// Create a context
	ctx := context.Background()

	// Create a resource store
	store := NewResourceStore()

	// First transaction - creates resources
	fmt.Println("=== First Transaction ===")
	err = saga.WithTransaction(ctx, logger, func(tx *saga.Transaction) error {
		// Add idempotent operations
		tx.AddOperation(
			// Idempotent operation to create resource A
			func(ctx context.Context) error {
				success := store.Create("A", "Value A")
				if !success {
					return fmt.Errorf("failed to create resource A")
				}
				return nil
			},
			// Idempotent rollback to delete resource A
			func(ctx context.Context) error {
				success := store.Delete("A")
				if !success {
					return fmt.Errorf("failed to delete resource A")
				}
				return nil
			},
		)

		tx.AddOperation(
			// Idempotent operation to create resource B
			func(ctx context.Context) error {
				success := store.Create("B", "Value B")
				if !success {
					return fmt.Errorf("failed to create resource B")
				}
				return nil
			},
			// Idempotent rollback to delete resource B
			func(ctx context.Context) error {
				success := store.Delete("B")
				if !success {
					return fmt.Errorf("failed to delete resource B")
				}
				return nil
			},
		)

		return nil
	})

	if err != nil {
		fmt.Printf("First transaction failed: %v\n", err)
	} else {
		fmt.Println("First transaction completed successfully")
	}

	// Second transaction - tries to create the same resources (demonstrates idempotency)
	fmt.Println("\n=== Second Transaction ===")
	err = saga.WithTransaction(ctx, logger, func(tx *saga.Transaction) error {
		// Add the same operations again
		tx.AddOperation(
			// Idempotent operation to create resource A (already exists)
			func(ctx context.Context) error {
				success := store.Create("A", "Value A")
				if !success {
					return fmt.Errorf("failed to create resource A")
				}
				return nil
			},
			// Idempotent rollback to delete resource A
			func(ctx context.Context) error {
				success := store.Delete("A")
				if !success {
					return fmt.Errorf("failed to delete resource A")
				}
				return nil
			},
		)

		tx.AddOperation(
			// Idempotent operation to create resource B (already exists)
			func(ctx context.Context) error {
				success := store.Create("B", "Value B")
				if !success {
					return fmt.Errorf("failed to create resource B")
				}
				return nil
			},
			// Idempotent rollback to delete resource B
			func(ctx context.Context) error {
				success := store.Delete("B")
				if !success {
					return fmt.Errorf("failed to delete resource B")
				}
				return nil
			},
		)

		// This operation will fail
		tx.AddOperation(
			func(ctx context.Context) error {
				fmt.Println("Failing operation...")
				return fmt.Errorf("operation failed")
			},
			saga.NoopRollback(),
		)

		return nil
	})

	if err != nil {
		fmt.Printf("Second transaction failed: %v\n", err)
	} else {
		fmt.Println("Second transaction completed successfully")
	}

	// Expected output:
	// === First Transaction ===
	// Resource A created with value Value A
	// Resource B created with value Value B
	// First transaction completed successfully
	//
	// === Second Transaction ===
	// Resource A already exists, skipping creation
	// Resource B already exists, skipping creation
	// Failing operation...
	// Resource B deleted
	// Resource A deleted
	// Second transaction failed: transaction execution failed: transaction operation failed: operation failed
}
