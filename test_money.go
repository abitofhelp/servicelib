// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Change to the valueobject directory
	err := os.Chdir("valueobject")
	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		os.Exit(1)
	}

	// Run the specific test
	cmd := exec.Command("go", "test", "-v", "-run", "TestNewMoneyFromUint64/Large_Amount")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Test failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Test passed successfully!")
}