// Copyright (c) 2025 A Bit of Help, Inc.

package identification

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewID(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		expectError bool
	}{
		{"Valid UUID", "f47ac10b-58cc-0372-8567-0e02b2c3d479", false},
		{"Empty String", "", true},
		{"Invalid UUID", "not-a-uuid", true},
		{"Invalid UUID Format", "f47ac10b-58cc-0372-8567-0e02b2c3d479-extra", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewID(tt.id)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if id.String() != tt.id {
					t.Errorf("Expected ID %s, got %s", tt.id, id.String())
				}
			}
		})
	}
}

func TestGenerateID(t *testing.T) {
	id := GenerateID()

	// Verify it's not empty
	if id.IsEmpty() {
		t.Errorf("Generated ID should not be empty")
	}

	// Verify it's a valid UUID
	_, err := uuid.Parse(id.String())
	if err != nil {
		t.Errorf("Generated ID is not a valid UUID: %v", err)
	}

	// Generate another ID and verify it's different
	id2 := GenerateID()
	if id.Equals(id2) {
		t.Errorf("Two generated IDs should not be equal")
	}
}

func TestID_String(t *testing.T) {
	uuidStr := "f47ac10b-58cc-0372-8567-0e02b2c3d479"
	id, _ := NewID(uuidStr)

	if id.String() != uuidStr {
		t.Errorf("Expected string %s, got %s", uuidStr, id.String())
	}
}

func TestID_Equals(t *testing.T) {
	id1, _ := NewID("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	id2, _ := NewID("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	id3, _ := NewID("550e8400-e29b-41d4-a716-446655440000")

	if !id1.Equals(id2) {
		t.Errorf("Expected id1 to equal id2")
	}

	if id1.Equals(id3) {
		t.Errorf("Expected id1 to not equal id3")
	}
}

func TestID_IsEmpty(t *testing.T) {
	// We can't create an empty ID with NewID, so we'll create it directly
	var emptyID ID = ""
	id, _ := NewID("f47ac10b-58cc-0372-8567-0e02b2c3d479")

	if !emptyID.IsEmpty() {
		t.Errorf("Expected empty ID to be empty")
	}

	if id.IsEmpty() {
		t.Errorf("Expected non-empty ID to not be empty")
	}
}
