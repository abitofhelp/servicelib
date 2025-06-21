//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example of basic usage of the env package
package main

import (
	"fmt"
	"os"

	"github.com/abitofhelp/servicelib/env"
)

func main() {
	fmt.Println("Environment Variables Basic Usage Example")
	fmt.Println("=========================================")

	// Get an environment variable with a fallback value
	port := env.GetEnv("PORT", "8080")
	fmt.Printf("Server will run on port: %s\n", port)

	// Get a database URL with a fallback
	dbURL := env.GetEnv("DATABASE_URL", "postgres://localhost:5432/mydb")
	fmt.Printf("Database URL: %s\n", dbURL)

	// Get API keys (sensitive information)
	apiKey := env.GetEnv("API_KEY", "")
	if apiKey == "" {
		fmt.Println("Warning: API_KEY environment variable is not set")
	} else {
		// In a real application, you wouldn't print the API key
		// This is just for demonstration purposes
		fmt.Printf("API Key: %s\n", apiKey)
	}

	// Show all environment variables
	fmt.Println("\nCurrent Environment Variables:")
	fmt.Println("-----------------------------")
	for _, e := range os.Environ() {
		fmt.Println(e)
	}

	fmt.Println("\nTry running this example with different environment variables set:")
	fmt.Println("export PORT=9090")
	fmt.Println("export DATABASE_URL=\"postgres://user:password@localhost:5432/mydb\"")
	fmt.Println("export API_KEY=\"your-api-key\"")
	fmt.Println("go run examples/env/basic_usage_example.go")
}
