// Copyright (c) 2025 A Bit of Help, Inc.

// Example demonstrating how to use the configscan package to scan packages for configuration requirements
// and check if they have default values set.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abitofhelp/servicelib/configscan"
)

func main() {
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get the root directory of the project
	rootDir := filepath.Join(currentDir, "../..")

	// Create a scanner with the root directory
	scanner := configscan.NewScanner(rootDir)

	// Scan packages for configuration requirements
	requirements, err := scanner.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print all configuration requirements
	fmt.Println("All configuration requirements:")
	for _, req := range requirements {
		fmt.Printf("- Package: %s, Config Type: %s, Fields: %v, Has Default: %v\n",
			req.PackageName, req.ConfigType, req.Fields, req.HasDefault)
	}

	// Generate a report of packages with configuration requirements that don't have default values set
	missingDefaults := configscan.Report(requirements)

	// Print the report
	fmt.Println("\nPackages with configuration requirements that don't have default values set:")
	if len(missingDefaults) == 0 {
		fmt.Println("All packages with configuration requirements have default values set.")
	} else {
		for _, req := range missingDefaults {
			fmt.Printf("- Package: %s, Config Type: %s, Fields: %v\n",
				req.PackageName, req.ConfigType, req.Fields)
		}
	}

	// Example of how to use the configscan package programmatically
	fmt.Println("\nExample of how to check if a specific package has default values for its configuration requirements:")
	for _, req := range requirements {
		if req.PackageName == "db" && req.ConfigType == "PostgresConfig" {
			if req.HasDefault {
				fmt.Println("The PostgresConfig in the db package has default values set.")
			} else {
				fmt.Println("The PostgresConfig in the db package does not have default values set.")
				fmt.Println("Consider adding a DefaultPostgresConfig function to provide default values.")
			}
		}
	}
}