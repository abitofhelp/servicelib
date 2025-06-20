// Copyright (c) 2025 A Bit of Help, Inc.

// Example of error handling with model package
package main

import (
    "fmt"
    
    "github.com/abitofhelp/servicelib/model"
)

type Source struct {
    Name string
    Age  int
}

type Destination struct {
    Name    string
    Address string // Different field type
}

func main() {
    // Create source and destination
    src := &Source{Name: "John", Age: 30}
    dst := &Destination{}
    
    // Copy fields
    err := model.CopyFields(dst, src)
    if err != nil {
        fmt.Printf("Error copying fields: %v\n", err)
        return
    }
    
    // Print result
    fmt.Printf("Destination: %+v\n", dst)
    // Output: Destination: {Name:John Address:}
    
    // Try to copy between incompatible types
    type NotAStruct string
    srcNotStruct := NotAStruct("test")
    dstNotStruct := NotAStruct("")
    
    err = model.CopyFields(&dstNotStruct, &srcNotStruct)
    if err != nil {
        fmt.Printf("Expected error: %v\n", err)
        // Output: Expected error: both source and destination must be pointers to structs
    }
    
    // Try to deep copy between different types
    differentSrc := &Source{Name: "Alice", Age: 25}
    differentDst := &Destination{}
    
    err = model.DeepCopy(differentDst, differentSrc)
    if err != nil {
        fmt.Printf("Expected error from DeepCopy: %v\n", err)
        // Output: Expected error from DeepCopy: source and destination must be of the same type
    }
    
    // Try to copy with non-pointer values
    err = model.CopyFields(Source{}, Destination{})
    if err != nil {
        fmt.Printf("Expected error with non-pointers: %v\n", err)
        // Output: Expected error with non-pointers: both source and destination must be pointers
    }
}