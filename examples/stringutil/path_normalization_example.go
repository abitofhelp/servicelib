// Copyright (c) 2025 A Bit of Help, Inc.

// Example of path normalization operations
package main

import (
    "fmt"
    "path/filepath"
    "runtime"
    
    "github.com/abitofhelp/servicelib/stringutil"
)

func main() {
    // Convert Windows-style paths to Unix-style
    windowsPath := "C:\\Users\\username\\Documents\\file.txt"
    unixPath := stringutil.ForwardSlashPath(windowsPath)
    fmt.Printf("Windows path: %s\n", windowsPath)
    fmt.Printf("Normalized to Unix: %s\n", unixPath) // Output: C:/Users/username/Documents/file.txt
    
    // Already Unix-style paths remain unchanged
    alreadyUnix := "/home/username/documents/file.txt"
    normalized := stringutil.ForwardSlashPath(alreadyUnix)
    fmt.Printf("\nUnix path: %s\n", alreadyUnix)
    fmt.Printf("After normalization: %s\n", normalized) // Output: /home/username/documents/file.txt
    
    // Mixed path separators
    mixedPath := "C:/Users\\username/Documents\\file.txt"
    normalizedMixed := stringutil.ForwardSlashPath(mixedPath)
    fmt.Printf("\nMixed separators: %s\n", mixedPath)
    fmt.Printf("Normalized: %s\n", normalizedMixed)
    
    // Real-world examples
    fmt.Println("\nReal-world examples:")
    
    // Example 1: Normalizing paths for cross-platform configuration
    configPaths := []string{
        "C:\\Program Files\\App\\config.json",
        "/etc/app/config.json",
        "..\\..\\config\\settings.yaml",
        "./config/local.json",
    }
    
    fmt.Println("Normalized configuration paths:")
    for _, path := range configPaths {
        fmt.Printf("  %s -> %s\n", path, stringutil.ForwardSlashPath(path))
    }
    
    // Example 2: Working with current OS path and normalizing
    if runtime.GOOS == "windows" {
        fmt.Println("\nRunning on Windows, normalizing current OS paths:")
    } else {
        fmt.Println("\nRunning on Unix-like OS, normalizing current OS paths:")
    }
    
    // Get absolute path to current directory
    absPath, err := filepath.Abs(".")
    if err != nil {
        fmt.Printf("Error getting absolute path: %v\n", err)
    } else {
        fmt.Printf("Current directory: %s\n", absPath)
        fmt.Printf("Normalized: %s\n", stringutil.ForwardSlashPath(absPath))
    }
    
    // Example 3: URL-like paths
    urlPaths := []string{
        "http:\\\\example.com\\path\\to\\resource",
        "file:\\\\server\\share\\document.pdf",
    }
    
    fmt.Println("\nNormalizing URL-like paths:")
    for _, path := range urlPaths {
        fmt.Printf("  %s -> %s\n", path, stringutil.ForwardSlashPath(path))
    }
}