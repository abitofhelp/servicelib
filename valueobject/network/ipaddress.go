// Copyright (c) 2025 A Bit of Help, Inc.

// Package network provides value objects related to network information.
package network

import (
	"encoding/json"
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

	// Use the normalized string format from net.ParseIP
	return IPAddress(parsedIP.String()), nil
}

// String returns the string representation of the IPAddress
func (ip IPAddress) String() string {
	return string(ip)
}

// normalizeIPv4 removes leading zeros from each segment of an IPv4 address
// e.g., "192.168.001.001" -> "192.168.1.1"
func normalizeIPv4(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return ip // Not an IPv4 address in dotted decimal notation
	}

	for i, part := range parts {
		// Remove leading zeros
		parts[i] = strings.TrimLeft(part, "0")
		if parts[i] == "" {
			parts[i] = "0" // If the part was all zeros, keep one zero
		}
	}

	return strings.Join(parts, ".")
}

// Equals checks if two IPAddresses are equal
// Properly handles comparison of both IPv4 and IPv6 addresses in different representations
// (e.g., IPv4 mapped to IPv6, or different IPv6 notations for the same address)
func (ip IPAddress) Equals(other IPAddress) bool {
	// Handle empty cases
	if ip.IsEmpty() && other.IsEmpty() {
		return true
	}
	if ip.IsEmpty() || other.IsEmpty() {
		return false
	}

	// Special handling for IPv4 addresses with leading zeros
	ipStr := string(ip)
	otherStr := string(other)

	// Check if both might be IPv4 addresses (contains dots)
	if strings.Contains(ipStr, ".") && strings.Contains(otherStr, ".") {
		// Normalize both addresses by removing leading zeros
		normalizedIP := normalizeIPv4(ipStr)
		normalizedOther := normalizeIPv4(otherStr)

		// If normalization worked (both are valid IPv4), compare the normalized strings
		if normalizedIP != ipStr || normalizedOther != otherStr {
			return normalizedIP == normalizedOther
		}
	}

	// Parse both IPs for comparison to handle different representations
	parsedThis := net.ParseIP(ipStr)
	parsedOther := net.ParseIP(otherStr)

	// If either IP can't be parsed, fall back to string comparison
	if parsedThis == nil || parsedOther == nil {
		return ipStr == otherStr
	}

	// For IPv4 addresses, convert to 4-byte representation for comparison
	thisIPv4 := parsedThis.To4()
	otherIPv4 := parsedOther.To4()

	if thisIPv4 != nil && otherIPv4 != nil {
		// Both are IPv4, compare the 4-byte representation
		return thisIPv4.Equal(otherIPv4)
	}

	// Compare normalized IPs
	return parsedThis.Equal(parsedOther)
}

// IsEmpty checks if the IPAddress is empty
func (ip IPAddress) IsEmpty() bool {
	return ip == ""
}

// Validate checks if the IPAddress is valid
func (ip IPAddress) Validate() error {
	// Empty value is allowed (optional field)
	if ip.IsEmpty() {
		return nil
	}

	// Validate IP format (both IPv4 and IPv6)
	parsedIP := net.ParseIP(string(ip))
	if parsedIP == nil {
		return errors.New("invalid IP address format")
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (ip IPAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(ip))
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
	if ipv6 := parsedIP.To16(); ipv6 != nil && (ipv6[0] == 0xfc || ipv6[0] == 0xfd) {
		return true
	}

	return false
}
