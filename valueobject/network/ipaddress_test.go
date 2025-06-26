// Copyright (c) 2025 A Bit of Help, Inc.

package network

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIPAddress(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{"Valid IPv4", "192.168.1.1", "192.168.1.1", false},
		{"Valid IPv6", "2001:db8::1", "2001:db8::1", false},
		{"Valid IP with Spaces", " 192.168.1.1 ", "192.168.1.1", false},
		{"Empty Value", "", "", false}, // Empty is allowed
		{"Invalid IP", "invalid", "", true},
		{"Invalid Format", "256.256.256.256", "", true},
		{"Incomplete IPv4", "192.168.1", "", true},
		{"Invalid IPv6", "2001:db8::1::1", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := NewIPAddress(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, value.String())
			}
		})
	}
}

func TestIPAddress_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"IPv4 Address", "192.168.1.1", "192.168.1.1"},
		{"IPv6 Address", "2001:db8::1", "2001:db8::1"},
		{"Empty Value", "", ""},
		{"IPv4 with Leading Zeros", "192.168.001.001", "192.168.1.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewIPAddress(tt.input)
			assert.Equal(t, tt.expected, value.String())
		})
	}
}

func TestIPAddress_Equals(t *testing.T) {
	tests := []struct {
		name     string
		value1   string
		value2   string
		expected bool
	}{
		{"Same IPv4", "192.168.1.1", "192.168.1.1", true},
		{"Same IPv6", "2001:db8::1", "2001:db8::1", true},
		{"IPv4 with Leading Zeros", "192.168.1.1", "192.168.001.001", true},
		{"Different IPv6 Format", "2001:db8::1", "2001:0db8:0000:0000:0000:0000:0000:0001", true},
		{"Different IPs", "192.168.1.1", "192.168.1.2", false},
		{"IPv4 vs IPv6", "192.168.1.1", "2001:db8::1", false},
		{"Empty Values", "", "", true},
		{"One Empty Value", "192.168.1.1", "", false},
		{"Invalid IPs", "invalid1", "invalid2", false},
		{"Invalid vs Valid", "invalid", "192.168.1.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value1, _ := NewIPAddress(tt.value1)
			value2, _ := NewIPAddress(tt.value2)
			assert.Equal(t, tt.expected, value1.Equals(value2))
		})
	}
}

func TestIPAddress_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"Empty Value", "", true},
		{"IPv4 Address", "192.168.1.1", false},
		{"IPv6 Address", "2001:db8::1", false},
		{"Whitespace", "  ", true}, // NewIPAddress trims whitespace
		{"Invalid IP", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewIPAddress(tt.value)
			assert.Equal(t, tt.expected, value.IsEmpty())
		})
	}
}

func TestIPAddress_Validate(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"Valid IPv4", "192.168.1.1", false},
		{"Valid IPv6", "2001:db8::1", false},
		{"Empty Value", "", false},
		{"Invalid IP", "invalid", true},
		{"Incomplete IP", "192.168.1", true},
		{"Invalid IPv4", "256.256.256.256", true},
		{"Invalid IPv6", "2001:db8::1::1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var value IPAddress
			if tt.expectError {
				value = IPAddress(tt.input)
			} else {
				var err error
				value, err = NewIPAddress(tt.input)
				require.NoError(t, err)
			}

			err := value.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIPAddress_MarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedJSON  string
		expectError   bool
		unmarshalTest bool
	}{
		{"IPv4 Address", "192.168.1.1", `"192.168.1.1"`, false, true},
		{"IPv6 Address", "2001:db8::1", `"2001:db8::1"`, false, true},
		{"Empty Value", "", `""`, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewIPAddress(tt.input)
			data, err := value.MarshalJSON()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedJSON, string(data))

				if tt.unmarshalTest {
					var str string
					err = json.Unmarshal(data, &str)
					assert.NoError(t, err)
					assert.Equal(t, value.String(), str)
				}
			}
		})
	}
}

func TestIPAddress_IsIPv4(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid IPv4", "192.168.1.1", true},
		{"IPv6 Address", "2001:db8::1", false},
		{"Empty Value", "", false},
		{"Invalid IP", "invalid", false},
		{"Incomplete IPv4", "192.168.1", false},
		{"IPv4-Mapped IPv6", "::ffff:192.168.1.1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewIPAddress(tt.input)
			assert.Equal(t, tt.expected, value.IsIPv4())
		})
	}
}

func TestIPAddress_IsIPv6(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid IPv6", "2001:db8::1", true},
		{"IPv6 Full Format", "2001:0db8:0000:0000:0000:0000:0000:0001", true},
		{"IPv4 Address", "192.168.1.1", false},
		{"Empty Value", "", false},
		{"Invalid IP", "invalid", false},
		{"Invalid IPv6", "2001:db8::1::1", false},
		{"IPv4-Mapped IPv6", "::ffff:192.168.1.1", false}, // This should be treated as IPv4
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewIPAddress(tt.input)
			assert.Equal(t, tt.expected, value.IsIPv6())
		})
	}
}

func TestIPAddress_IsLoopback(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"IPv4 Loopback", "127.0.0.1", true},
		{"IPv4 Loopback Range", "127.1.2.3", true},
		{"IPv6 Loopback", "::1", true},
		{"Non-Loopback IPv4", "192.168.1.1", false},
		{"Non-Loopback IPv6", "2001:db8::1", false},
		{"Empty Value", "", false},
		{"Invalid IP", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewIPAddress(tt.input)
			assert.Equal(t, tt.expected, value.IsLoopback())
		})
	}
}

func TestIPAddress_IsPrivate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"IPv4 Private 10.x.x.x", "10.0.0.1", true},
		{"IPv4 Private 172.16.x.x", "172.16.0.1", true},
		{"IPv4 Private 172.31.x.x", "172.31.255.255", true},
		{"IPv4 Private 192.168.x.x", "192.168.1.1", true},
		{"IPv4 Public", "8.8.8.8", false},
		{"IPv6 Private fc00::/7", "fc00::1", true},
		{"IPv6 Private fd00::/7", "fd00::1", true},
		{"IPv6 Public", "2001:db8::1", false},
		{"Empty Value", "", false},
		{"Invalid IP", "invalid", false},
		{"Edge Case - Not Private", "172.32.0.1", false},
		{"Edge Case - Not Private", "172.15.0.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewIPAddress(tt.input)
			assert.Equal(t, tt.expected, value.IsPrivate())
		})
	}
}
