# String Utilities

## Overview

The String Utilities component provides additional string manipulation functions beyond what's available in the standard library. It offers case-insensitive operations, natural language formatting, whitespace handling, and other common string manipulation tasks that are frequently needed in applications.

## Features

- **Case-Insensitive Operations**: Perform prefix and substring checks without case sensitivity
- **Multiple Prefix Checking**: Check if a string starts with any of multiple prefixes
- **Natural Language Formatting**: Join strings with commas and "and" for human-readable output
- **Whitespace Handling**: Check for empty strings and remove whitespace
- **String Truncation**: Safely truncate strings with ellipsis
- **Path Normalization**: Convert backslashes to forward slashes for cross-platform compatibility
- **UTF-8 Safe**: All operations are safe for UTF-8 encoded strings
- **Zero Dependencies**: No external dependencies beyond the Go standard library

## Installation

```bash
go get github.com/abitofhelp/servicelib/stringutil
```

## Quick Start

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/stringutil"
)

func main() {
    // Case-insensitive operations
    fmt.Println(stringutil.HasPrefixIgnoreCase("Hello World", "hello"))  // true
    fmt.Println(stringutil.ContainsIgnoreCase("Hello World", "WORLD"))   // true
    
    // Join strings with natural language formatting
    items := []string{"apples", "bananas", "oranges"}
    fmt.Println(stringutil.JoinWithAnd(items, true))  // "apples, bananas, and oranges"
    
    // Check for empty strings
    fmt.Println(stringutil.IsEmpty("   "))  // true
    fmt.Println(stringutil.IsEmpty("Hello"))  // false
    
    // Truncate long strings
    longText := "This is a very long string that needs to be truncated"
    fmt.Println(stringutil.Truncate(longText, 20))  // "This is a very long..."
    
    // Normalize file paths
    path := "C:\\Users\\Username\\Documents"
    fmt.Println(stringutil.ForwardSlashPath(path))  // "C:/Users/Username/Documents"
}
```

## API Documentation

### Key Functions

#### HasPrefixIgnoreCase

Checks if a string begins with a prefix, ignoring case.

```go
func HasPrefixIgnoreCase(s, prefix string) bool
```

#### ContainsIgnoreCase

Checks if a substring is within a string, ignoring case.

```go
func ContainsIgnoreCase(s, substr string) bool
```

#### HasAnyPrefix

Checks if a string begins with any of the specified prefixes.

```go
func HasAnyPrefix(s string, prefixes ...string) bool
```

#### ToLowerCase

Converts a string to lowercase.

```go
func ToLowerCase(s string) string
```

#### JoinWithAnd

Joins a slice of strings with commas and "and".

```go
func JoinWithAnd(items []string, useOxfordComma bool) string
```

#### IsEmpty

Checks if a string is empty or contains only whitespace.

```go
func IsEmpty(s string) bool
```

#### IsNotEmpty

Checks if a string is not empty and contains non-whitespace characters.

```go
func IsNotEmpty(s string) bool
```

#### Truncate

Truncates a string to a specified maximum length.

```go
func Truncate(s string, maxLength int) string
```

#### RemoveWhitespace

Removes all whitespace characters from a string.

```go
func RemoveWhitespace(s string) string
```

#### ForwardSlashPath

Converts backslashes in a path to forward slashes.

```go
func ForwardSlashPath(path string) string
```

## Examples

Currently, there are no dedicated examples for the stringutil package in the EXAMPLES directory. The following code snippets demonstrate common usage patterns:

### Case-Insensitive String Operations

```go
// Check if a string starts with a prefix, ignoring case
if stringutil.HasPrefixIgnoreCase(userInput, "cmd:") {
    // Process as a command
}

// Check if a string contains a substring, ignoring case
if stringutil.ContainsIgnoreCase(email, "@example.com") {
    // Handle example.com email addresses
}
```

### Natural Language Formatting

```go
// Join items with commas and "and"
selectedOptions := []string{"Option A", "Option B", "Option C"}
message := "You selected " + stringutil.JoinWithAnd(selectedOptions, true)
// "You selected Option A, Option B, and Option C"

// Without Oxford comma
message = "You selected " + stringutil.JoinWithAnd(selectedOptions, false)
// "You selected Option A, Option B and Option C"
```

### String Validation and Manipulation

```go
// Check if a string is empty or whitespace-only
if stringutil.IsEmpty(userInput) {
    return errors.New("input cannot be empty")
}

// Truncate a string for display
title := stringutil.Truncate(longArticleTitle, 50)

// Remove all whitespace
compactString := stringutil.RemoveWhitespace("  This   has   extra   spaces  ")
// "Thishasextraspaces"
```

## Best Practices

1. **Use Case-Insensitive Functions When Appropriate**: For user input and non-exact matching, use the case-insensitive functions
2. **Consider Oxford Comma Preference**: When using JoinWithAnd, consider your audience's preference for Oxford commas
3. **Check for Empty Strings**: Use IsEmpty instead of checking for "" to also catch whitespace-only strings
4. **Set Appropriate Truncation Lengths**: When truncating strings, ensure the maxLength is appropriate for your display context
5. **Normalize Paths**: Use ForwardSlashPath when working with file paths across different operating systems

## Troubleshooting

### Common Issues

#### Performance Considerations

If you're processing large volumes of strings:
- Consider caching results of expensive operations
- Be aware that regular expressions (used in RemoveWhitespace) can be slower for very large strings
- For high-performance needs, consider using string builders for concatenation operations

#### Unicode Handling

All functions are UTF-8 safe, but be aware of potential issues:
- Some operations may behave differently with non-ASCII characters
- Length calculations are based on bytes, not characters, which can affect truncation of multi-byte characters

## Related Components

- [Validation](../validation/README.md) - String validation utilities
- [Errors](../errors/README.md) - Error handling that may use string formatting

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
