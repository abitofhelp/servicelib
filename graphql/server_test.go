// Copyright (c) 2025 A Bit of Help, Inc.

package graphql

import (
	"context"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/abitofhelp/servicelib/graphql/mocks"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"
)

// TestNewDefaultServerConfig tests the NewDefaultServerConfig function
func TestNewDefaultServerConfig(t *testing.T) {
	// Call the function
	config := NewDefaultServerConfig()

	// Verify the default values
	assert.Equal(t, 25, config.MaxQueryDepth, "Default max query depth should be 25")
	assert.Equal(t, 100, config.MaxQueryComplexity, "Default max query complexity should be 100")
	assert.Equal(t, 30*time.Second, config.RequestTimeout, "Default request timeout should be 30 seconds")
}

// TestNewServer tests the NewServer function
func TestNewServer(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock schema
	mockSchema := mocks.NewMockExecutableSchema(ctrl)

	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Call the function with default config
	server := NewServer(mockSchema, logger, NewDefaultServerConfig())

	// Verify that the server is not nil
	assert.NotNil(t, server)
	assert.IsType(t, &handler.Server{}, server)
}

// TestNewServerComplexity tests the complexity limit configuration of the server
func TestNewServerComplexity(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock schema
	mockSchema := mocks.NewMockExecutableSchema(ctrl)

	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a custom config with specific complexity limits
	config := ServerConfig{
		MaxQueryDepth:      10,
		MaxQueryComplexity: 50,
		RequestTimeout:     30 * time.Second,
	}

	// Call the function with custom config
	server := NewServer(mockSchema, logger, config)

	// We can't directly test the complexity limits, but we can verify that the server was created
	assert.NotNil(t, server)
}

// TestNewServerWithCustomTimeout tests the server with a custom timeout
func TestNewServerWithCustomTimeout(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock schema
	mockSchema := mocks.NewMockExecutableSchema(ctrl)

	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a custom config with a specific timeout
	config := ServerConfig{
		MaxQueryDepth:      25,
		MaxQueryComplexity: 100,
		RequestTimeout:     5 * time.Second,
	}

	// Call the function with custom config
	server := NewServer(mockSchema, logger, config)

	// We can't directly test the timeout, but we can verify that the server was created
	assert.NotNil(t, server)
}

// TestNewServerWithEmptyConfig tests the server with an empty config
func TestNewServerWithEmptyConfig(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock schema
	mockSchema := mocks.NewMockExecutableSchema(ctrl)

	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create an empty config
	config := ServerConfig{}

	// Call the function with empty config
	server := NewServer(mockSchema, logger, config)

	// We can't directly test the config values, but we can verify that the server was created
	assert.NotNil(t, server)
}

// TestNewServerWithExtremeValues tests the server with extreme config values
func TestNewServerWithExtremeValues(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock schema
	mockSchema := mocks.NewMockExecutableSchema(ctrl)

	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a config with extreme values
	config := ServerConfig{
		MaxQueryDepth:      1000000,
		MaxQueryComplexity: 1000000,
		RequestTimeout:     24 * time.Hour,
	}

	// Call the function with extreme config
	server := NewServer(mockSchema, logger, config)

	// We can't directly test the config values, but we can verify that the server was created
	assert.NotNil(t, server)
}

// TestNewServerWithZeroTimeout tests the server with a zero timeout
func TestNewServerWithZeroTimeout(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock schema
	mockSchema := mocks.NewMockExecutableSchema(ctrl)

	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a custom config with a zero timeout
	config := ServerConfig{
		MaxQueryDepth:      25,
		MaxQueryComplexity: 100,
		RequestTimeout:     0,
	}

	// Call the function with custom config
	server := NewServer(mockSchema, logger, config)

	// We can't directly test the timeout, but we can verify that the server was created
	assert.NotNil(t, server)
}

// TestNewServerWithNegativeTimeout tests the server with a negative timeout
func TestNewServerWithNegativeTimeout(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock schema
	mockSchema := mocks.NewMockExecutableSchema(ctrl)

	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a custom config with a negative timeout
	config := ServerConfig{
		MaxQueryDepth:      25,
		MaxQueryComplexity: 100,
		RequestTimeout:     -1 * time.Second,
	}

	// Call the function with custom config
	server := NewServer(mockSchema, logger, config)

	// We can't directly test the timeout, but we can verify that the server was created
	assert.NotNil(t, server)
}

// TestCreateAroundOperationsFunc tests the createAroundOperationsFunc function
func TestCreateAroundOperationsFunc(t *testing.T) {
	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create the middleware function
	middleware := createAroundOperationsFunc(logger, 30*time.Second)

	// Verify that the middleware function is not nil
	assert.NotNil(t, middleware)
}

// TestCreateAroundOperationsFuncValidation tests validation in the middleware
func TestCreateAroundOperationsFuncValidation(t *testing.T) {
	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a context
	ctx := context.Background()

	// Create the middleware function
	middleware := createAroundOperationsFunc(logger, 30*time.Second)

	// Create a mock next handler that returns a response handler
	nextCalled := false
	next := func(ctx context.Context) graphql.ResponseHandler {
		nextCalled = true
		return func(ctx context.Context) *graphql.Response {
			return &graphql.Response{}
		}
	}

	// Test with missing operation name
	// We need to set up a context with operation context that has no name
	opCtx := &graphql.OperationContext{
		Operation: &ast.OperationDefinition{
			Name: "",
		},
	}
	ctx = graphql.WithOperationContext(ctx, opCtx)

	// Call the middleware
	responseHandler := middleware(ctx, next)

	// The next handler should not be called because validation failed
	assert.False(t, nextCalled)

	// Call the response handler
	response := responseHandler(ctx)

	// Verify that we got an error response
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Errors)
	assert.Equal(t, "Operation must have a name", response.Errors[0].Message)
}

// TestCreateAroundOperationsFuncWithValidOperation tests the middleware with a valid operation
func TestCreateAroundOperationsFuncWithValidOperation(t *testing.T) {
	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a context
	ctx := context.Background()

	// Create the middleware function
	middleware := createAroundOperationsFunc(logger, 30*time.Second)

	// Create a mock next handler that returns a response handler
	nextCalled := false
	next := func(ctx context.Context) graphql.ResponseHandler {
		nextCalled = true
		return func(ctx context.Context) *graphql.Response {
			return &graphql.Response{
				Data: []byte(`{"data":{"test":"value"}}`),
			}
		}
	}

	// Set up a context with operation context that has a name
	opCtx := &graphql.OperationContext{
		Operation: &ast.OperationDefinition{
			Name: "TestQuery",
		},
	}
	ctx = graphql.WithOperationContext(ctx, opCtx)

	// Call the middleware
	responseHandler := middleware(ctx, next)

	// The next handler should be called because the operation is valid
	assert.True(t, nextCalled)

	// Call the response handler
	response := responseHandler(ctx)

	// Verify that we got a successful response
	assert.NotNil(t, response)
	assert.Empty(t, response.Errors)
	assert.NotEmpty(t, response.Data)
}

// TestCreateAroundOperationsFuncWithCancelledContext tests the middleware with a cancelled context
func TestCreateAroundOperationsFuncWithCancelledContext(t *testing.T) {
	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	// Create the middleware function
	middleware := createAroundOperationsFunc(logger, 30*time.Second)

	// Create a mock next handler that returns a response handler
	nextCalled := false
	next := func(ctx context.Context) graphql.ResponseHandler {
		nextCalled = true
		return func(ctx context.Context) *graphql.Response {
			return &graphql.Response{}
		}
	}

	// Set up a context with operation context that has a name
	opCtx := &graphql.OperationContext{
		Operation: &ast.OperationDefinition{
			Name: "TestQuery",
		},
	}
	ctx = graphql.WithOperationContext(ctx, opCtx)

	// Call the middleware
	responseHandler := middleware(ctx, next)

	// The next handler should still be called
	assert.True(t, nextCalled)

	// Call the response handler
	response := responseHandler(ctx)

	// Verify that we got an error response about the cancelled context
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Errors)
	assert.Contains(t, response.Errors[0].Message, "interrupted")
	assert.Equal(t, "REQUEST_INTERRUPTED", response.Errors[0].Extensions["code"])
}

// TestCreateAroundOperationsFuncWithCancelledParentContext tests the middleware with a parent context that's cancelled before processing
func TestCreateAroundOperationsFuncWithCancelledParentContext(t *testing.T) {
	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a context that's not cancelled yet
	ctx := context.Background()

	// Create the middleware function
	middleware := createAroundOperationsFunc(logger, 30*time.Second)

	// Create a mock next handler that returns a response handler
	nextCalled := false
	next := func(ctx context.Context) graphql.ResponseHandler {
		nextCalled = true
		return func(ctx context.Context) *graphql.Response {
			return &graphql.Response{}
		}
	}

	// Set up a context with operation context that has a name
	opCtx := &graphql.OperationContext{
		Operation: &ast.OperationDefinition{
			Name: "TestQuery",
		},
	}
	ctx = graphql.WithOperationContext(ctx, opCtx)

	// Call the middleware to get the response handler
	responseHandler := middleware(ctx, next)

	// The next handler should be called
	assert.True(t, nextCalled)

	// Now create a new context that's cancelled for the response handler
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	// Call the response handler with the cancelled context
	response := responseHandler(cancelledCtx)

	// Verify that we got an error response about the cancelled context
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Errors)
	assert.Contains(t, response.Errors[0].Message, "interrupted")
	assert.Equal(t, "REQUEST_INTERRUPTED", response.Errors[0].Extensions["code"])
}

// TestCreateAroundOperationsFuncWithCancellationDuringProcessing tests the middleware with a context that's cancelled during processing
func TestCreateAroundOperationsFuncWithCancellationDuringProcessing(t *testing.T) {
	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a context with a cancel function
	ctx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	// Create the middleware function
	middleware := createAroundOperationsFunc(logger, 30*time.Second)

	// Create a channel to signal when the response handler is called
	handlerCalled := make(chan struct{})

	// Create a mock next handler that returns a response handler that blocks until cancelled
	next := func(ctx context.Context) graphql.ResponseHandler {
		return func(ctx context.Context) *graphql.Response {
			// Signal that the handler was called
			close(handlerCalled)

			// Block until the context is cancelled
			<-ctx.Done()

			return &graphql.Response{}
		}
	}

	// Set up a context with operation context that has a name
	opCtx := &graphql.OperationContext{
		Operation: &ast.OperationDefinition{
			Name: "TestQuery",
		},
	}
	ctx = graphql.WithOperationContext(ctx, opCtx)

	// Call the middleware to get the response handler
	responseHandler := middleware(ctx, next)

	// Start a goroutine to call the response handler
	responseCh := make(chan *graphql.Response)
	go func() {
		responseCh <- responseHandler(ctx)
	}()

	// Wait for the handler to be called
	<-handlerCalled

	// Cancel the parent context
	parentCancel()

	// Get the response
	response := <-responseCh

	// Verify that we got an error response about the client disconnect
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Errors)
	assert.Contains(t, response.Errors[0].Message, "client disconnect")
	assert.Equal(t, "CLIENT_DISCONNECTED", response.Errors[0].Extensions["code"])
}

// TestCreateAroundOperationsFuncWithTimeout tests the middleware with a timeout
func TestCreateAroundOperationsFuncWithTimeout(t *testing.T) {
	// Create a logger for testing
	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewContextLogger(zapLogger)

	// Create a context
	ctx := context.Background()

	// Create the middleware function with a very short timeout
	middleware := createAroundOperationsFunc(logger, 10*time.Millisecond)

	// Create a mock next handler that returns a response handler that sleeps
	nextCalled := false
	next := func(ctx context.Context) graphql.ResponseHandler {
		nextCalled = true
		return func(ctx context.Context) *graphql.Response {
			// Sleep longer than the timeout
			time.Sleep(20 * time.Millisecond)
			return &graphql.Response{}
		}
	}

	// Set up a context with operation context that has a name
	opCtx := &graphql.OperationContext{
		Operation: &ast.OperationDefinition{
			Name: "TestQuery",
		},
	}
	ctx = graphql.WithOperationContext(ctx, opCtx)

	// Call the middleware
	responseHandler := middleware(ctx, next)

	// The next handler should be called
	assert.True(t, nextCalled)

	// Call the response handler
	response := responseHandler(ctx)

	// Verify that we got a timeout error
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Errors)
	assert.Contains(t, response.Errors[0].Message, "timed out")
	assert.Equal(t, "TIMEOUT", response.Errors[0].Extensions["code"])
}
