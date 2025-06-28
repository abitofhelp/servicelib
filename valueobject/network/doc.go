// Copyright (c) 2025 A Bit of Help, Inc.

// Package network provides value objects related to network information.
//
// This package contains value objects that represent network-related concepts
// such as IP addresses and URLs. These value objects are immutable and follow
// the Value Object pattern from Domain-Driven Design.
//
// Key value objects in this package:
//   - IPAddress: Represents an IP address (IPv4 or IPv6)
//   - URL: Represents a URL with validation and parsing capabilities
//
// The IPAddress value object provides methods for:
//   - Creating and validating IP addresses
//   - String representation
//   - Equality comparison
//   - Checking IP address properties (IPv4, IPv6, loopback, private)
//   - JSON marshaling
//
// The URL value object provides methods for:
//   - Creating and validating URLs
//   - String representation
//   - Equality comparison
//   - Extracting URL components (domain, path, query)
//   - Validation of scheme and host
//
// Example usage:
//
//	// Create a new IP address
//	ip, err := network.NewIPAddress("192.168.1.1")
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Check IP address properties
//	if ip.IsIPv4() {
//	    // Handle IPv4 address
//	}
//
//	if ip.IsPrivate() {
//	    // Handle private IP address
//	}
//
//	// Create a new URL
//	url, err := network.NewURL("https://example.com/path?query=value")
//	if err != nil {
//	    // Handle validation error
//	}
//
//	// Extract URL components
//	domain, err := url.Domain()
//	if err != nil {
//	    // Handle error
//	}
//
//	path, err := url.Path()
//	if err != nil {
//	    // Handle error
//	}
//
// All value objects in this package are designed to be immutable, so they cannot
// be changed after creation. To modify a value object, create a new instance with
// the desired values.
package network
