// Copyright (c) 2025 A Bit of Help, Inc.

package graphql

// DefaultServerConfig returns a default configuration for the GraphQL server.
// This is an alias for NewDefaultServerConfig for compatibility with the configscan tool.
func DefaultServerConfig() ServerConfig {
	return NewDefaultServerConfig()
}
