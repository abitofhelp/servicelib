// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import "github.com/abitofhelp/servicelib/errors/core"

// AddDetail adds a single detail to an error's context
func AddDetail(err error, key string, value interface{}) error {
	if err == nil {
		return nil
	}

	if ce, ok := err.(*core.ContextualError); ok {
		if ce.Context.Details == nil {
			ce.Context.Details = make(map[string]interface{})
		}
		ce.Context.Details[key] = value
		return ce
	}

	// Create a new ContextualError if the input error is not already one
	return &core.ContextualError{
		Original: err,
		Context: core.ErrorContext{
			Details: map[string]interface{}{
				key: value,
			},
		},
	}
}

// AddDetails adds multiple details to an error's context
func AddDetails(err error, details map[string]interface{}) error {
	if err == nil {
		return nil
	}

	if ce, ok := err.(*core.ContextualError); ok {
		if ce.Context.Details == nil {
			ce.Context.Details = make(map[string]interface{})
		}
		for k, v := range details {
			ce.Context.Details[k] = v
		}
		return ce
	}

	// Create a new ContextualError if the input error is not already one
	return &core.ContextualError{
		Original: err,
		Context: core.ErrorContext{
			Details: details,
		},
	}
}

// GetDetail retrieves a detail from an error's context
func GetDetail(err error, key string) (interface{}, bool) {
	if err == nil {
		return nil, false
	}

	if ce, ok := err.(*core.ContextualError); ok {
		if ce.Context.Details != nil {
			value, exists := ce.Context.Details[key]
			return value, exists
		}
	}

	return nil, false
}

// GetDetails retrieves all details from an error's context
func GetDetails(err error) map[string]interface{} {
	if err == nil {
		return nil
	}

	if ce, ok := err.(*core.ContextualError); ok {
		if ce.Context.Details == nil {
			return make(map[string]interface{})
		}
		return ce.Context.Details
	}

	return nil
}
