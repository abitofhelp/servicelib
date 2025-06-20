// Copyright (c) 2025 A Bit of Help, Inc.

// Example of whitespace handling operations
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/stringutil"
)

func main() {
    // Check if a string is empty or contains only whitespace
    emptyStrings := []string{
        "",
        "   ",
        "\t\n\r",
        "  \t  \n  ",
    }
    
    fmt.Println("IsEmpty examples:")
    for _, s := range emptyStrings {
        isEmpty := stringutil.IsEmpty(s)
        fmt.Printf("IsEmpty(%q) = %v\n", s, isEmpty) // All should be true
    }
    
    // Check if a string contains non-whitespace characters
    nonEmptyStrings := []string{
        "Hello",
        "  Hello  ",
        "\tHello\n",
        " . ",  // Just a period with spaces
    }
    
    fmt.Println("\nIsNotEmpty examples:")
    for _, s := range nonEmptyStrings {
        isNotEmpty := stringutil.IsNotEmpty(s)
        fmt.Printf("IsNotEmpty(%q) = %v\n", s, isNotEmpty) // All should be true
    }
    
    // Remove all whitespace from a string
    stringsWithWhitespace := []string{
        "Hello, World!",
        "  Spaces  at  ends  ",
        "Tabs\tand\tnewlines\n",
        "Mixed   whitespace \t\n characters",
    }
    
    fmt.Println("\nRemoveWhitespace examples:")
    for _, s := range stringsWithWhitespace {
        noWhitespace := stringutil.RemoveWhitespace(s)
        fmt.Printf("RemoveWhitespace(%q) = %q\n", s, noWhitespace)
    }
    
    // Real-world examples
    fmt.Println("\nReal-world examples:")
    
    // Example 1: Validating user input
    userInput := "   "
    if stringutil.IsEmpty(userInput) {
        fmt.Println("Please enter a non-empty value")
    } else {
        fmt.Println("Input is valid")
    }
    
    // Example 2: Normalizing phone numbers
    phoneNumber := "(123) 456-7890"
    normalized := stringutil.RemoveWhitespace(phoneNumber)
    fmt.Printf("Original phone: %s\nNormalized: %s\n", phoneNumber, normalized)
    
    // Example 3: Checking if a configuration value is set
    configValue := "\n"
    if stringutil.IsEmpty(configValue) {
        fmt.Println("Configuration value is not set, using default")
    } else {
        fmt.Println("Using configured value")
    }
}