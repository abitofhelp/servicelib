# String Utilities Package

The `stringutil` package provides additional string manipulation utilities beyond what's available in the standard library. It offers a collection of helper functions for common string operations that are frequently needed in Go applications.

## Features

- **Case-Insensitive Operations**: Functions for case-insensitive string comparisons
- **Prefix Checking**: Utilities for checking string prefixes with various options
- **String Joining**: Advanced string joining with natural language formatting
- **Whitespace Handling**: Functions for checking and manipulating whitespace
- **String Truncation**: Safely truncate strings with ellipsis
- **Path Normalization**: Convert file paths to use consistent separators
- **Regular Expression Utilities**: String manipulation using regular expressions

## Installation

```bash
go get github.com/abitofhelp/servicelib/stringutil
```

## Usage

### Case-Insensitive Operations

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/stringutil"
)

func main() {
    // Check if a string starts with a prefix (case-insensitive)
    hasPrefix := stringutil.HasPrefixIgnoreCase("Hello, World!", "hello")
    fmt.Printf("Has prefix 'hello': %v\n", hasPrefix) // Output: true
    
    // Check if a string contains a substring (case-insensitive)
    contains := stringutil.ContainsIgnoreCase("Hello, World!", "WORLD")
    fmt.Printf("Contains 'WORLD': %v\n", contains) // Output: true
    
    // Convert to lowercase (wrapper around strings.ToLower)
    lower := stringutil.ToLowerCase("Hello, World!")
    fmt.Printf("Lowercase: %s\n", lower) // Output: hello, world!
}
```

### Prefix Checking

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/stringutil"
)

func main() {
    // Check if a string starts with any of the given prefixes
    hasAnyPrefix := stringutil.HasAnyPrefix("https://example.com", "http://", "https://")
    fmt.Printf("Has any prefix: %v\n", hasAnyPrefix) // Output: true
    
    // Check multiple strings for prefixes
    urls := []string{
        "https://example.com",
        "http://example.org",
        "ftp://example.net",
        "example.io",
    }
    
    for _, url := range urls {
        if stringutil.HasAnyPrefix(url, "http://", "https://") {
            fmt.Printf("%s is a web URL\n", url)
        } else {
            fmt.Printf("%s is not a web URL\n", url)
        }
    }
}
```

### String Joining with Natural Language Formatting

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/stringutil"
)

func main() {
    // Join strings with "and" (no Oxford comma)
    fruits := []string{"apples", "bananas", "oranges"}
    joined := stringutil.JoinWithAnd(fruits, false)
    fmt.Printf("Without Oxford comma: %s\n", joined) // Output: apples, bananas and oranges
    
    // Join strings with "and" (with Oxford comma)
    joinedOxford := stringutil.JoinWithAnd(fruits, true)
    fmt.Printf("With Oxford comma: %s\n", joinedOxford) // Output: apples, bananas, and oranges
    
    // Handle different list lengths
    fmt.Println(stringutil.JoinWithAnd([]string{}, false))                // Output: ""
    fmt.Println(stringutil.JoinWithAnd([]string{"apples"}, false))        // Output: "apples"
    fmt.Println(stringutil.JoinWithAnd([]string{"apples", "bananas"}, false)) // Output: "apples and bananas"
}
```

### Whitespace Handling

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/stringutil"
)

func main() {
    // Check if a string is empty or contains only whitespace
    isEmpty := stringutil.IsEmpty("   \t\n   ")
    fmt.Printf("Is empty: %v\n", isEmpty) // Output: true
    
    // Check if a string contains non-whitespace characters
    isNotEmpty := stringutil.IsNotEmpty("  Hello  ")
    fmt.Printf("Is not empty: %v\n", isNotEmpty) // Output: true
    
    // Remove all whitespace from a string
    noWhitespace := stringutil.RemoveWhitespace("Hello, \t World! \n")
    fmt.Printf("No whitespace: %s\n", noWhitespace) // Output: Hello,World!
}
```

### String Truncation

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/stringutil"
)

func main() {
    // Truncate a long string
    longText := "This is a very long string that needs to be truncated to fit in a limited space."
    truncated := stringutil.Truncate(longText, 20)
    fmt.Printf("Truncated: %s\n", truncated) // Output: This is a very long...
    
    // No truncation needed for short strings
    shortText := "Short text"
    truncatedShort := stringutil.Truncate(shortText, 20)
    fmt.Printf("Truncated short: %s\n", truncatedShort) // Output: Short text
    
    // Handle negative max length
    negativeMax := stringutil.Truncate(longText, -5)
    fmt.Printf("Negative max: %s\n", negativeMax) // Output: ...
}
```

### Path Normalization

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/stringutil"
)

func main() {
    // Convert Windows-style paths to Unix-style
    windowsPath := "C:\\Users\\username\\Documents\\file.txt"
    unixPath := stringutil.ForwardSlashPath(windowsPath)
    fmt.Printf("Unix path: %s\n", unixPath) // Output: C:/Users/username/Documents/file.txt
    
    // Already Unix-style paths remain unchanged
    alreadyUnix := "/home/username/documents/file.txt"
    normalized := stringutil.ForwardSlashPath(alreadyUnix)
    fmt.Printf("Normalized: %s\n", normalized) // Output: /home/username/documents/file.txt
}
```

## Best Practices

1. **Performance Considerations**: For performance-critical code, be aware that some functions like `ContainsIgnoreCase` create temporary strings.

2. **String Immutability**: Remember that strings in Go are immutable, so functions like `Truncate` and `RemoveWhitespace` return new strings.

3. **Unicode Support**: The functions in this package work with UTF-8 encoded strings, which is Go's default string encoding.

4. **Error Handling**: These utility functions don't return errors, so validate input when necessary before calling them.

5. **Path Handling**: For more comprehensive path manipulation, consider using the standard library's `path/filepath` package alongside `ForwardSlashPath`.

6. **Regular Expressions**: `RemoveWhitespace` uses regular expressions, which can be expensive for very large strings or frequent calls.

7. **String Builders**: For complex string manipulations, consider using `strings.Builder` for better performance.

## License

This project is licensed under the MIT License - see the LICENSE file for details.