//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating basic usage of the cache package
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/cache"
)

// User is a simple struct for demonstration purposes
type User struct {
	ID   string
	Name string
	Age  int
}

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a cache with default configuration
	cfg := cache.DefaultConfig()
	options := cache.DefaultOptions().WithName("example")
	userCache := cache.NewCache[User](cfg, options)

	// Define a function that simulates fetching a user from a database
	fetchUser := func(ctx context.Context) (User, error) {
		// In a real application, this would be a database query
		fmt.Println("Fetching user from database...")
		// Simulate some work
		time.Sleep(100 * time.Millisecond)
		return User{ID: "123", Name: "John Doe", Age: 30}, nil
	}

	// Use the cache to get a user
	fmt.Println("Getting user for the first time (should fetch from database)...")
	user, err := cache.WithCache(ctx, userCache, "user:123", fetchUser)
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
	} else {
		fmt.Printf("Got user: %+v\n", user)
	}

	// Get the same user again (should be cached)
	fmt.Println("\nGetting user for the second time (should be cached)...")
	user, err = cache.WithCache(ctx, userCache, "user:123", fetchUser)
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
	} else {
		fmt.Printf("Got user: %+v\n", user)
	}

	// Manually set a value in the cache
	fmt.Println("\nManually setting a value in the cache...")
	userCache.Set(ctx, "user:456", User{ID: "456", Name: "Jane Smith", Age: 25})

	// Get the manually set value
	fmt.Println("\nGetting manually set value...")
	user, found := userCache.Get(ctx, "user:456")
	if found {
		fmt.Printf("Got user from cache: %+v\n", user)
	} else {
		fmt.Println("User not found in cache")
	}

	// Delete a value from the cache
	fmt.Println("\nDeleting a value from the cache...")
	userCache.Delete(ctx, "user:123")

	// Try to get the deleted value
	fmt.Println("\nTrying to get deleted value...")
	user, found = userCache.Get(ctx, "user:123")
	if found {
		fmt.Printf("Got user from cache: %+v\n", user)
	} else {
		fmt.Println("User not found in cache (as expected)")
	}

	// Get the user again (should fetch from database again)
	fmt.Println("\nGetting user after deletion (should fetch from database again)...")
	user, err = cache.WithCache(ctx, userCache, "user:123", fetchUser)
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
	} else {
		fmt.Printf("Got user: %+v\n", user)
	}

	// Shutdown the cache
	fmt.Println("\nShutting down the cache...")
	userCache.Shutdown()
}