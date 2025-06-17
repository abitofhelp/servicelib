// Copyright (c) 2025 A Bit of Help, Inc.

//go:generate go run github.com/golang/mock/mockgen -package=mocks -destination=mock_graphql.go github.com/99designs/gqlgen/graphql ExecutableSchema
package mocks