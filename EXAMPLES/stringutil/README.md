# StringUtil Package Examples

This directory contains examples demonstrating how to use the `stringutil` package, which provides utilities for string manipulation and processing in Go applications. The package offers various helper functions for common string operations that extend beyond what's available in the standard library.

## Examples

### 1. Case Insensitive Operations Example

[case_insensitive_operations_example.go](case_insensitive_operations_example.go)

Demonstrates how to perform case-insensitive string operations.

Key concepts:
- Case-insensitive string comparison
- Case-insensitive string searching
- Case-insensitive string replacement
- Handling Unicode characters in case-insensitive operations
- Performance considerations for case-insensitive operations

### 2. Path Normalization Example

[path_normalization_example.go](path_normalization_example.go)

Shows how to normalize file and URL paths.

Key concepts:
- Normalizing file paths across different operating systems
- Handling relative paths
- Resolving path references
- Cleaning and simplifying paths
- URL path normalization

### 3. Prefix Checking Example

[prefix_checking_example.go](prefix_checking_example.go)

Demonstrates how to check if strings have specific prefixes.

Key concepts:
- Checking if a string starts with a prefix
- Case-sensitive and case-insensitive prefix checking
- Handling multiple possible prefixes
- Removing prefixes from strings
- Performance optimizations for prefix operations

### 4. String Joining Example

[string_joining_example.go](string_joining_example.go)

Shows how to join strings with various delimiters.

Key concepts:
- Joining string slices with delimiters
- Handling empty strings in join operations
- Conditional joining based on string content
- Performance considerations for string joining
- Building complex strings from multiple parts

### 5. String Truncation Example

[string_truncation_example.go](string_truncation_example.go)

Demonstrates how to truncate strings to a specified length.

Key concepts:
- Truncating strings to a maximum length
- Adding ellipsis to truncated strings
- Handling edge cases (negative or zero length)
- Preserving word boundaries during truncation
- Real-world applications of string truncation

### 6. Whitespace Handling Example

[whitespace_handling_example.go](whitespace_handling_example.go)

Shows how to handle whitespace in strings.

Key concepts:
- Trimming whitespace from strings
- Normalizing whitespace within strings
- Collapsing multiple whitespace characters
- Handling different types of whitespace (spaces, tabs, newlines)
- Unicode whitespace considerations

## Running the Examples

To run any of the examples, use the `go run` command:

```bash
go run examples/stringutil/string_truncation_example.go
```

## Additional Resources

For more information about the stringutil package, see the [stringutil package documentation](../../stringutil/README.md).