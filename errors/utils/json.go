// Copyright (c) 2024 A Bit of Help, Inc.

package utils

import (
	"encoding/json"
	"fmt"
	"github.com/abitofhelp/servicelib/errors/core"
)

// ToJSON converts an error to a JSON string.
func ToJSON(err error) string {
	if err == nil {
		return ""
	}

	if ce, ok := err.(*core.ContextualError); ok {
		bytes, _ := json.Marshal(ce)
		return string(bytes)
	}

	// Create a simple JSON object for non-contextual errors
	return fmt.Sprintf(`{"message": "%s"}`, err.Error())
}

// FromJSON converts a JSON string to an error.
func FromJSON(jsonStr string) (error, error) {
	var ce core.ContextualError
	if err := json.Unmarshal([]byte(jsonStr), &ce); err != nil {
		return nil, fmt.Errorf("failed to parse error JSON: %w", err)
	}
	return &ce, nil
}
