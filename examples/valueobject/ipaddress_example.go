// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the IPAddress value object
package main

import (
	"fmt"

	"github.com/abitofhelp/servicelib/valueobject"
)

func main() {
	// Create a new IPv4 address
	ipv4, err := valueobject.NewIPAddress("192.168.1.1")
	if err != nil {
		// Handle error
		fmt.Println("Error creating IPv4 address:", err)
		return
	}

	// Create a new IPv6 address
	ipv6, err := valueobject.NewIPAddress("2001:db8::1")
	if err != nil {
		// Handle error
		fmt.Println("Error creating IPv6 address:", err)
		return
	}

	// Check IP types
	fmt.Printf("Is %s an IPv4 address? %v\n", ipv4, ipv4.IsIPv4()) // true
	fmt.Printf("Is %s an IPv6 address? %v\n", ipv6, ipv6.IsIPv6()) // true

	// Check if IP is loopback
	loopback, _ := valueobject.NewIPAddress("127.0.0.1")
	fmt.Printf("Is %s a loopback address? %v\n", loopback, loopback.IsLoopback()) // true

	// Check if IP is private
	private, _ := valueobject.NewIPAddress("10.0.0.1")
	fmt.Printf("Is %s a private address? %v\n", private, private.IsPrivate()) // true

	// Compare IP addresses
	sameIP, _ := valueobject.NewIPAddress("192.168.1.1")
	areEqual := ipv4.Equals(sameIP)
	fmt.Printf("Are IPs equal? %v\n", areEqual) // true
}
