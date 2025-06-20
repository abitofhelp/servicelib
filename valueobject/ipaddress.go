// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"net"
	"strings"
)

// IPAddress represents an IP address value object that supports both IPv4 and IPv6 formats
type IPAddress string

// NewIPAddress creates a new IPAddress with validation for both IPv4 and IPv6 formats
// It uses Go's net.ParseIP function which supports both IPv4 (e.g., "192.168.1.1")
// and IPv6 (e.g., "2001:db8::1") address formats
func NewIPAddress(ip string) (IPAddress, error) {
	// Trim whitespace
	trimmedIP := strings.TrimSpace(ip)

	// Empty IP is allowed (optional field)
	if trimmedIP == "" {
		return "", nil
	}

	// Validate IP format (both IPv4 and IPv6)
	parsedIP := net.ParseIP(trimmedIP)
	if parsedIP == nil {
		return "", errors.New("invalid IP address format")
	}

	return IPAddress(trimmedIP), nil
}

// String returns the string representation of the IPAddress
func (ip IPAddress) String() string {
	return string(ip)
}

// Equals checks if two IPAddresses are equal
// Properly handles comparison of both IPv4 and IPv6 addresses in different representations
// (e.g., IPv4 mapped to IPv6, or different IPv6 notations for the same address)
func (ip IPAddress) Equals(other IPAddress) bool {
	// Parse both IPs for comparison to handle different representations
	parsedThis := net.ParseIP(string(ip))
	parsedOther := net.ParseIP(string(other))

	// If either IP can't be parsed, fall back to string comparison
	if parsedThis == nil || parsedOther == nil {
		return string(ip) == string(other)
	}

	// Compare normalized IPs
	return parsedThis.Equal(parsedOther)
}

// IsEmpty checks if the IPAddress is empty
func (ip IPAddress) IsEmpty() bool {
	return ip == ""
}

// IsIPv4 checks if the IP address is in IPv4 format (e.g., "192.168.1.1")
// Returns true only for valid IPv4 addresses, false for IPv6 addresses or invalid formats
func (ip IPAddress) IsIPv4() bool {
	if ip.IsEmpty() {
		return false
	}

	parsedIP := net.ParseIP(string(ip))
	if parsedIP == nil {
		return false
	}

	return parsedIP.To4() != nil
}

// IsIPv6 checks if the IP address is in IPv6 format (e.g., "2001:db8::1")
// Returns true only for valid IPv6 addresses, false for IPv4 addresses or invalid formats
func (ip IPAddress) IsIPv6() bool {
	if ip.IsEmpty() {
		return false
	}

	parsedIP := net.ParseIP(string(ip))
	if parsedIP == nil {
		return false
	}

	return parsedIP.To4() == nil
}

// IsLoopback checks if the IP address is a loopback address
// Works with both IPv4 loopback (127.0.0.0/8) and IPv6 loopback (::1) addresses
func (ip IPAddress) IsLoopback() bool {
	if ip.IsEmpty() {
		return false
	}

	parsedIP := net.ParseIP(string(ip))
	if parsedIP == nil {
		return false
	}

	return parsedIP.IsLoopback()
}

// IsPrivate checks if the IP address is in a private range
// Supports both IPv4 private ranges (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16)
// and IPv6 private ranges (fc00::/7)
func (ip IPAddress) IsPrivate() bool {
	if ip.IsEmpty() {
		return false
	}

	parsedIP := net.ParseIP(string(ip))
	if parsedIP == nil {
		return false
	}

	// Check for private IPv4 ranges
	if ipv4 := parsedIP.To4(); ipv4 != nil {
		// 10.0.0.0/8
		if ipv4[0] == 10 {
			return true
		}
		// 172.16.0.0/12
		if ipv4[0] == 172 && ipv4[1] >= 16 && ipv4[1] <= 31 {
			return true
		}
		// 192.168.0.0/16
		if ipv4[0] == 192 && ipv4[1] == 168 {
			return true
		}
	}

	// Check for private IPv6 ranges (fc00::/7)
	if ipv6 := parsedIP.To16(); ipv6 != nil && ipv6[0] == 0xfc || ipv6[0] == 0xfd {
		return true
	}

	return false
}
