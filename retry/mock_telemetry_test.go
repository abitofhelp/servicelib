// Copyright (c) 2025 A Bit of Help, Inc.

package retry

import (
	"context"

	"github.com/abitofhelp/servicelib/telemetry"
	"github.com/golang/mock/gomock"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// MockSpan is a mock implementation of telemetry.Span
type MockSpan struct {
	ctrl     *gomock.Controller
	recorder *MockSpanRecorder
}

// MockSpanRecorder is the mock recorder for MockSpan
type MockSpanRecorder struct {
	mock *MockSpan
}

// NewMockSpan creates a new mock instance
func NewMockSpan(ctrl *gomock.Controller) *MockSpan {
	mock := &MockSpan{ctrl: ctrl}
	mock.recorder = &MockSpanRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSpan) EXPECT() *MockSpanRecorder {
	return m.recorder
}

// End mocks base method
func (m *MockSpan) End() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "End")
}

// End indicates an expected call of End
func (mr *MockSpanRecorder) End() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "End", nil)
}

// SetAttributes mocks base method
func (m *MockSpan) SetAttributes(attributes ...attribute.KeyValue) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range attributes {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "SetAttributes", varargs...)
}

// SetAttributes indicates an expected call of SetAttributes
func (mr *MockSpanRecorder) SetAttributes(attributes ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAttributes", nil, attributes...)
}

// RecordError mocks base method
func (m *MockSpan) RecordError(err error, opts ...trace.EventOption) {
	m.ctrl.T.Helper()
	varargs := []interface{}{err}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "RecordError", varargs...)
}

// RecordError indicates an expected call of RecordError
func (mr *MockSpanRecorder) RecordError(err interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := []interface{}{err}
	varargs = append(varargs, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordError", nil, varargs...)
}

// MockTracer is a mock implementation of telemetry.Tracer
type MockTracer struct {
	ctrl     *gomock.Controller
	recorder *MockTracerRecorder
}

// MockTracerRecorder is the mock recorder for MockTracer
type MockTracerRecorder struct {
	mock *MockTracer
}

// NewMockTracer creates a new mock instance
func NewMockTracer(ctrl *gomock.Controller) *MockTracer {
	mock := &MockTracer{ctrl: ctrl}
	mock.recorder = &MockTracerRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTracer) EXPECT() *MockTracerRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockTracer) Start(ctx context.Context, name string) (context.Context, telemetry.Span) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx, name)
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(telemetry.Span)
	return ret0, ret1
}

// Start indicates an expected call of Start
func (mr *MockTracerRecorder) Start(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", nil, ctx, name)
}