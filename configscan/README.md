# ConfigScan

ConfigScan is a tool for scanning packages to determine whether any have configuration requirements and don't have default values set.

## Overview

Many packages in a Go project define configuration structs to allow users to configure their behavior. It's a best practice to provide default values for these configuration structs so that users don't have to specify every configuration parameter.

ConfigScan scans packages in a Go project and identifies configuration structs that don't have default values set. This helps ensure that all packages with configuration requirements provide sensible defaults.

## Usage

### As a Command-Line Tool

```bash
# Scan the current directory and its subdirectories
go run cmd/main.go

# Scan a specific directory
go run cmd/main.go -dir=/path/to/directory
```

### As a Library

```go
package main

import (
	"fmt"
	"os"

	"github.com/abitofhelp/servicelib/configscan"
)

func main() {
	// Create a scanner with the directory to scan
	scanner := configscan.NewScanner(".")

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
	}
}
```

## How It Works

ConfigScan works by:

1. Walking through all directories in the specified root directory
2. Parsing Go files
3. Finding structs with names containing "Config"
4. Checking if there's a default function for each config struct
5. Reporting config structs that don't have default functions

A default function is a function with a name matching one of these patterns:
- `Default` + ConfigType (e.g., `DefaultMyConfig` for a struct named `MyConfig`)
- `Default` + TrimSuffix(ConfigType, "Config") + "Config" (e.g., `DefaultConfig` for a struct named `Config`)

## Best Practices

When defining a configuration struct in your package, always provide a default function that returns a configuration struct with sensible default values. For example:

```go
// Config contains configuration parameters for the package
type Config struct {
	Timeout      time.Duration
	MaxRetries   int
	BackoffFactor float64
}

// DefaultConfig returns a default configuration with sensible values
func DefaultConfig() Config {
	return Config{
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		BackoffFactor: 2.0,
	}
}
```

This ensures that users of your package don't have to specify every configuration parameter and can start with sensible defaults.