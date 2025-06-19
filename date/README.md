# Date Package

The `date` package provides utilities for working with dates and times in Go applications. It offers a consistent way to parse, format, and manipulate dates using a standardized format.

## Features

- **Standardized Date Format**: Uses RFC3339 format consistently throughout the application
- **Date Parsing**: Parse date strings into time.Time objects
- **Date Formatting**: Format time.Time objects into standardized strings
- **Optional Date Handling**: Special handling for nil date pointers
- **Error Handling**: Proper validation and error handling for invalid date formats

## Installation

```bash
go get github.com/abitofhelp/servicelib/date
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/abitofhelp/servicelib/date"
)

func main() {
    // Parse a date string
    dateStr := "2023-01-15T14:30:00Z"
    parsedDate, err := date.ParseDate(dateStr)
    if err != nil {
        fmt.Printf("Error parsing date: %v\n", err)
        return
    }
    fmt.Printf("Parsed date: %v\n", parsedDate)
    
    // Format a date
    now := time.Now()
    formatted := date.FormatDate(now)
    fmt.Printf("Formatted current time: %s\n", formatted)
}
```

### Working with Optional Dates

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/abitofhelp/servicelib/date"
)

func main() {
    // Parse an optional date string
    var dateStr *string
    
    // Handle nil case
    parsedNil, err := date.ParseOptionalDate(dateStr)
    if err != nil {
        fmt.Printf("Error parsing date: %v\n", err)
        return
    }
    fmt.Printf("Parsed nil date: %v\n", parsedNil) // Will be nil
    
    // Handle non-nil case
    validDateStr := "2023-01-15T14:30:00Z"
    dateStr = &validDateStr
    
    parsed, err := date.ParseOptionalDate(dateStr)
    if err != nil {
        fmt.Printf("Error parsing date: %v\n", err)
        return
    }
    fmt.Printf("Parsed optional date: %v\n", *parsed)
    
    // Format an optional date
    var optionalDate *time.Time
    
    // Handle nil case
    formattedNil := date.FormatOptionalDate(optionalDate)
    fmt.Printf("Formatted nil date: %v\n", formattedNil) // Will be nil
    
    // Handle non-nil case
    now := time.Now()
    optionalDate = &now
    
    formatted := date.FormatOptionalDate(optionalDate)
    fmt.Printf("Formatted optional date: %s\n", *formatted)
}
```

### Error Handling

```go
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/date"
)

func main() {
    // Try to parse an invalid date string
    invalidDateStr := "not-a-date"
    _, err := date.ParseDate(invalidDateStr)
    if err != nil {
        fmt.Printf("Error parsing invalid date: %v\n", err)
    }
    
    // Try to parse an invalid optional date string
    invalidOptStr := "also-not-a-date"
    _, err = date.ParseOptionalDate(&invalidOptStr)
    if err != nil {
        fmt.Printf("Error parsing invalid optional date: %v\n", err)
    }
}
```

## Best Practices

1. **Consistent Date Format**: Always use the package's standard date format for consistency across your application.

2. **Error Handling**: Always check for errors when parsing dates.

3. **Nil Handling**: Use the optional date functions when dealing with dates that might be nil.

4. **Time Zones**: Be aware of time zone implications when parsing and formatting dates.

5. **Date Comparisons**: When comparing dates, ensure they are in the same time zone to avoid unexpected results.

## License

This project is licensed under the MIT License - see the LICENSE file for details.