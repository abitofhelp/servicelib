// Package mocks contains mock implementations for telemetry package testing
package mocks

//go:generate go run github.com/golang/mock/mockgen -package=mocks -destination=mock_meter.go go.opentelemetry.io/otel/metric Meter,MeterProvider
//go:generate go run github.com/golang/mock/mockgen -package=mocks -destination=mock_tracer.go go.opentelemetry.io/otel/trace Tracer,TracerProvider
//go:generate go run github.com/golang/mock/mockgen -package=mocks -destination=mock_span.go go.opentelemetry.io/otel/trace Span