// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating custom configuration of the cache package
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/cache"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// Product is a simple struct for demonstration purposes
type Product struct {
	ID    string
	Name  string
	Price float64
}

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a logger
	logger, _ := zap.NewDevelopment()
	contextLogger := logging.NewContextLogger(logger)

	// Create a cache with custom configuration
	cfg := cache.DefaultConfig().
		WithEnabled(true).                // Explicitly enable the cache
		WithTTL(10 * time.Second).        // Set TTL to 10 seconds
		WithMaxSize(100).                 // Set max size to 100 items
		WithPurgeInterval(30 * time.Second) // Purge expired items every 30 seconds

	options := cache.DefaultOptions().
		WithName("product-cache").        // Set a custom name
		WithLogger(contextLogger)         // Use a custom logger

	productCache := cache.NewCache[Product](cfg, options)

	// Define a function that simulates fetching a product from a database
	fetchProduct := func(ctx context.Context) (Product, error) {
		// In a real application, this would be a database query
		fmt.Println("Fetching product from database...")
		// Simulate some work
		time.Sleep(100 * time.Millisecond)
		return Product{ID: "123", Name: "Laptop", Price: 999.99}, nil
	}

	// Use the cache to get a product
	fmt.Println("Getting product for the first time (should fetch from database)...")
	product, err := cache.WithCache(ctx, productCache, "product:123", fetchProduct)
	if err != nil {
		fmt.Printf("Error getting product: %v\n", err)
	} else {
		fmt.Printf("Got product: %+v\n", product)
	}

	// Get the same product again (should be cached)
	fmt.Println("\nGetting product for the second time (should be cached)...")
	product, err = cache.WithCache(ctx, productCache, "product:123", fetchProduct)
	if err != nil {
		fmt.Printf("Error getting product: %v\n", err)
	} else {
		fmt.Printf("Got product: %+v\n", product)
	}

	// Demonstrate custom TTL with WithCacheTTL
	fmt.Println("\nDemonstrating custom TTL with WithCacheTTL...")
	
	// Define a function that simulates fetching a product with a short TTL
	fetchProductWithShortTTL := func(ctx context.Context) (Product, error) {
		fmt.Println("Fetching product with short TTL from database...")
		time.Sleep(100 * time.Millisecond)
		return Product{ID: "456", Name: "Phone", Price: 499.99}, nil
	}
	
	// Use WithCacheTTL to cache with a custom TTL
	fmt.Println("Getting product with short TTL (2 seconds)...")
	product, err = cache.WithCacheTTL(ctx, productCache, "product:456", 2*time.Second, fetchProductWithShortTTL)
	if err != nil {
		fmt.Printf("Error getting product: %v\n", err)
	} else {
		fmt.Printf("Got product: %+v\n", product)
	}
	
	// Get the same product again (should be cached)
	fmt.Println("\nGetting product with short TTL again (should be cached)...")
	product, err = cache.WithCacheTTL(ctx, productCache, "product:456", 2*time.Second, fetchProductWithShortTTL)
	if err != nil {
		fmt.Printf("Error getting product: %v\n", err)
	} else {
		fmt.Printf("Got product: %+v\n", product)
	}
	
	// Wait for the short TTL to expire
	fmt.Println("\nWaiting for short TTL to expire (3 seconds)...")
	time.Sleep(3 * time.Second)
	
	// Get the product again (should fetch from database again)
	fmt.Println("\nGetting product after TTL expiration (should fetch from database again)...")
	product, err = cache.WithCacheTTL(ctx, productCache, "product:456", 2*time.Second, fetchProductWithShortTTL)
	if err != nil {
		fmt.Printf("Error getting product: %v\n", err)
	} else {
		fmt.Printf("Got product: %+v\n", product)
	}

	// Demonstrate individual parameter configuration
	fmt.Println("\nDemonstrating individual parameter configuration:")
	
	// Start with default config
	cfg = cache.DefaultConfig()
	
	// Configure TTL only
	cfg = cfg.WithTTL(20 * time.Second)
	fmt.Printf("TTL: %v\n", cfg.TTL)
	
	// Configure max size only
	cfg = cfg.WithMaxSize(200)
	fmt.Printf("Max size: %d\n", cfg.MaxSize)
	
	// Configure purge interval only
	cfg = cfg.WithPurgeInterval(60 * time.Second)
	fmt.Printf("Purge interval: %v\n", cfg.PurgeInterval)
	
	// Configure enabled only
	cfg = cfg.WithEnabled(false)
	fmt.Printf("Enabled: %t\n", cfg.Enabled)

	// Shutdown the cache
	fmt.Println("\nShutting down the cache...")
	productCache.Shutdown()
}