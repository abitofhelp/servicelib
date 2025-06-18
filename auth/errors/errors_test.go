// Copyright (c) 2025 A Bit of Help, Inc.

package errors_test

import (
	"errors"
	"testing"

	autherrors "github.com/abitofhelp/servicelib/auth/errors"
	"github.com/stretchr/testify/assert"
)

func TestStandardErrors(t *testing.T) {
	// Test that all standard errors are defined
	assert.NotNil(t, autherrors.ErrInvalidToken)
	assert.NotNil(t, autherrors.ErrExpiredToken)
	assert.NotNil(t, autherrors.ErrMissingToken)
	assert.NotNil(t, autherrors.ErrInvalidSignature)
	assert.NotNil(t, autherrors.ErrInvalidClaims)
	assert.NotNil(t, autherrors.ErrUnauthorized)
	assert.NotNil(t, autherrors.ErrForbidden)
	assert.NotNil(t, autherrors.ErrInvalidConfig)
	assert.NotNil(t, autherrors.ErrInternal)
}

func TestAuthError_Error(t *testing.T) {
	// Test with message
	err := autherrors.NewAuthError(autherrors.ErrInvalidToken, "op", "custom message", nil)
	assert.Equal(t, "custom message", err.Error())

	// Test without message
	err = autherrors.NewAuthError(autherrors.ErrInvalidToken, "op", "", nil)
	assert.Equal(t, autherrors.ErrInvalidToken.Error(), err.Error())
}

func TestAuthError_Unwrap(t *testing.T) {
	baseErr := autherrors.ErrInvalidToken
	err := autherrors.NewAuthError(baseErr, "op", "message", nil)

	// Test unwrap
	unwrapped := err.Unwrap()
	assert.Equal(t, baseErr, unwrapped)
}

func TestAuthError_Is(t *testing.T) {
	// Test Is with matching error
	err := autherrors.NewAuthError(autherrors.ErrInvalidToken, "op", "message", nil)
	assert.True(t, errors.Is(err, autherrors.ErrInvalidToken))

	// Test Is with non-matching error
	assert.False(t, errors.Is(err, autherrors.ErrExpiredToken))
}

func TestNewAuthError(t *testing.T) {
	// Test creating a new AuthError
	baseErr := autherrors.ErrInvalidToken
	op := "test.op"
	message := "test message"
	context := map[string]interface{}{"key": "value"}

	err := autherrors.NewAuthError(baseErr, op, message, context)

	assert.Equal(t, baseErr, err.Err)
	assert.Equal(t, op, err.Op)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, context, err.Context)
}

func TestWithContext(t *testing.T) {
	// Test adding context to an existing AuthError
	baseErr := autherrors.NewAuthError(autherrors.ErrInvalidToken, "op", "message", nil)
	err := autherrors.WithContext(baseErr, "key", "value")

	authErr, ok := err.(*autherrors.AuthError)
	assert.True(t, ok)
	assert.Equal(t, "value", authErr.Context["key"])

	// Test adding context to a non-AuthError
	stdErr := errors.New("standard error")
	err = autherrors.WithContext(stdErr, "key", "value")

	authErr, ok = err.(*autherrors.AuthError)
	assert.True(t, ok)
	assert.Equal(t, stdErr, authErr.Err)
	assert.Equal(t, "value", authErr.Context["key"])
}

func TestWithOp(t *testing.T) {
	// Test adding op to an existing AuthError
	baseErr := autherrors.NewAuthError(autherrors.ErrInvalidToken, "original.op", "message", nil)
	err := autherrors.WithOp(baseErr, "new.op")

	authErr, ok := err.(*autherrors.AuthError)
	assert.True(t, ok)
	assert.Equal(t, "new.op", authErr.Op)

	// Test adding op to a non-AuthError
	stdErr := errors.New("standard error")
	err = autherrors.WithOp(stdErr, "new.op")

	authErr, ok = err.(*autherrors.AuthError)
	assert.True(t, ok)
	assert.Equal(t, stdErr, authErr.Err)
	assert.Equal(t, "new.op", authErr.Op)
}

func TestWithMessage(t *testing.T) {
	// Test adding message to an existing AuthError
	baseErr := autherrors.NewAuthError(autherrors.ErrInvalidToken, "op", "original message", nil)
	err := autherrors.WithMessage(baseErr, "new message")

	authErr, ok := err.(*autherrors.AuthError)
	assert.True(t, ok)
	assert.Equal(t, "new message", authErr.Message)

	// Test adding message to a non-AuthError
	stdErr := errors.New("standard error")
	err = autherrors.WithMessage(stdErr, "new message")

	authErr, ok = err.(*autherrors.AuthError)
	assert.True(t, ok)
	assert.Equal(t, stdErr, authErr.Err)
	assert.Equal(t, "new message", authErr.Message)
}

func TestWrap(t *testing.T) {
	// Test wrapping an error with a message
	baseErr := errors.New("base error")
	err := autherrors.Wrap(baseErr, "wrapped")
	assert.Contains(t, err.Error(), "wrapped")
	assert.Contains(t, err.Error(), "base error")

	// Test wrapping nil error
	err = autherrors.Wrap(nil, "wrapped")
	assert.Nil(t, err)
}

func TestGetContext(t *testing.T) {
	// Test getting context from an AuthError
	context := map[string]interface{}{"key": "value"}
	err := autherrors.NewAuthError(autherrors.ErrInvalidToken, "op", "message", context)

	value, ok := autherrors.GetContext(err, "key")
	assert.True(t, ok)
	assert.Equal(t, "value", value)

	// Test getting non-existent context key
	value, ok = autherrors.GetContext(err, "nonexistent")
	assert.False(t, ok)
	assert.Nil(t, value)

	// Test getting context from a non-AuthError
	stdErr := errors.New("standard error")
	value, ok = autherrors.GetContext(stdErr, "key")
	assert.False(t, ok)
	assert.Nil(t, value)
}

func TestGetOp(t *testing.T) {
	// Test getting op from an AuthError
	err := autherrors.NewAuthError(autherrors.ErrInvalidToken, "test.op", "message", nil)

	op, ok := autherrors.GetOp(err)
	assert.True(t, ok)
	assert.Equal(t, "test.op", op)

	// Test getting op from an AuthError with empty op
	err = autherrors.NewAuthError(autherrors.ErrInvalidToken, "", "message", nil)
	op, ok = autherrors.GetOp(err)
	assert.False(t, ok)
	assert.Equal(t, "", op)

	// Test getting op from a non-AuthError
	stdErr := errors.New("standard error")
	op, ok = autherrors.GetOp(stdErr)
	assert.False(t, ok)
	assert.Equal(t, "", op)
}

func TestGetMessage(t *testing.T) {
	// Test getting message from an AuthError
	err := autherrors.NewAuthError(autherrors.ErrInvalidToken, "op", "test message", nil)

	message, ok := autherrors.GetMessage(err)
	assert.True(t, ok)
	assert.Equal(t, "test message", message)

	// Test getting message from an AuthError with empty message
	err = autherrors.NewAuthError(autherrors.ErrInvalidToken, "op", "", nil)
	message, ok = autherrors.GetMessage(err)
	assert.False(t, ok)
	assert.Equal(t, "", message)

	// Test getting message from a non-AuthError
	stdErr := errors.New("standard error")
	message, ok = autherrors.GetMessage(stdErr)
	assert.False(t, ok)
	assert.Equal(t, "", message)
}
