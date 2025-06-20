// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the FileSize value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new file size
	fileSize, err := valueobject.NewFileSize(1024, valueobject.Kilobytes)
	if err != nil {
		// Handle error
		fmt.Println("Error creating file size:", err)
		return
	}

	// Access values in different units
	bytes := fileSize.Bytes()
	kb := fileSize.Kilobytes()
	mb := fileSize.Megabytes()
	fmt.Printf("File size: %d bytes, %.2f KB, %.2f MB\n", bytes, kb, mb)

	// Format in different ways
	fmt.Println(fileSize.String())          // Auto format (1.00 MB)
	fmt.Println(fileSize.Format("B"))       // 1048576 B
	fmt.Println(fileSize.Format("MB"))      // 1.00 MB
	fmt.Println(fileSize.Format("binary"))  // 1.00 MiB
	fmt.Println(fileSize.Format("decimal")) // 1.05 MB

	// Parse from string
	parsedSize, err := valueobject.ParseFileSize("2.5GB")
	if err != nil {
		// Handle error
		fmt.Println("Error parsing file size:", err)
		return
	}
	fmt.Println(parsedSize.String()) // 2.50 GB

	// Perform calculations
	otherSize, _ := valueobject.NewFileSize(500, valueobject.Megabytes)
	sum := fileSize.Add(otherSize)
	fmt.Println(sum.String()) // 1.49 GB

	// Compare file sizes
	isLarger := fileSize.IsLargerThan(otherSize)
	fmt.Printf("Is 1024 KB larger than 500 MB? %v\n", isLarger) // false
}
