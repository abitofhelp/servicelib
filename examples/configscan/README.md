# ConfigScan Examples

This directory contains examples demonstrating how to use the configscan package to scan packages for configuration requirements and check if they have default values set.

## configscan_example.go

This example demonstrates how to use the configscan package to:

1. Scan packages for configuration requirements
2. Generate a report of packages with configuration requirements that don't have default values set
3. Check if a specific package has default values for its configuration requirements

### Running the Example

```bash
go run configscan_example.go
```

### Expected Output

The example will output:

1. A list of all configuration requirements found in the project
2. A report of packages with configuration requirements that don't have default values set
3. An example of how to check if a specific package has default values for its configuration requirements

### Code Explanation

The example:

1. Gets the root directory of the project
2. Creates a scanner with the root directory
3. Scans packages for configuration requirements
4. Prints all configuration requirements
5. Generates a report of packages with configuration requirements that don't have default values set
6. Prints the report
7. Demonstrates how to check if a specific package has default values for its configuration requirements

### Best Practices

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