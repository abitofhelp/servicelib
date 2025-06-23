// Copyright (c) 2025 A Bit of Help, Inc.

// Command configscan scans packages for configuration requirements and checks if they have default values set.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/abitofhelp/servicelib/configscan"
)

func main() {
	// Parse command-line flags
	rootDir := flag.String("dir", ".", "Root directory to scan")
	flag.Parse()

	// Get absolute path of the root directory
	absRootDir, err := filepath.Abs(*rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create a new scanner
	scanner := configscan.NewScanner(absRootDir)

	// Scan packages for configuration requirements
	requirements, err := scanner.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Generate a report of packages with configuration requirements that don't have default values set
	missingDefaults := configscan.Report(requirements)

	// Print the report
	if len(missingDefaults) == 0 {
		fmt.Println("All packages with configuration requirements have default values set.")
	} else {
		fmt.Println("Packages with configuration requirements that don't have default values set:")
		for _, req := range missingDefaults {
			fmt.Printf("- Package: %s, Config Type: %s, Fields: %v\n", req.PackageName, req.ConfigType, req.Fields)
		}
		os.Exit(1)
	}
}