// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"bytes"
	"testing"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{"Valid Password", "Password1!", false},
		{"Valid Complex Password", "C0mpl3x!P@ssw0rd", false},
		{"Empty Password", "", true},
		{"Too Short Password", "Pass1!", true},
		{"No Uppercase", "password1!", true},
		{"No Digit", "Password!", true},
		{"No Special Character", "Password1", true},
		{"Whitespace Trimmed", " Password1! ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := NewPassword(tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Verify the password was hashed
				if password.hashedPassword == "" {
					t.Errorf("Expected hashed password, got empty string")
				}

				// Verify the salt was generated
				if len(password.salt) == 0 {
					t.Errorf("Expected salt to be generated")
				}

				// Verify the password can be verified
				if !password.Verify(tt.password) {
					t.Errorf("Expected password to verify")
				}
			}
		})
	}
}

func TestPassword_Verify(t *testing.T) {
	plainPassword := "Password1!"
	password, _ := NewPassword(plainPassword)

	tests := []struct {
		name        string
		input       string
		shouldMatch bool
	}{
		{"Correct Password", plainPassword, true},
		{"Incorrect Password", "WrongPassword1!", false},
		{"Case Sensitive", "password1!", false},
		{"With Extra Character", plainPassword + "a", false},
		{"Missing Character", plainPassword[:len(plainPassword)-1], false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := password.Verify(tt.input)
			if result != tt.shouldMatch {
				t.Errorf("Expected Verify to return %v for %s", tt.shouldMatch, tt.input)
			}
		})
	}
}

func TestPassword_String(t *testing.T) {
	password, _ := NewPassword("Password1!")
	expected := "********"

	if password.String() != expected {
		t.Errorf("Expected string %s, got %s", expected, password.String())
	}
}

func TestPassword_Equals(t *testing.T) {
	password1, _ := NewPassword("Password1!")
	
	// Create a copy with the same hash and salt
	password2 := Password{
		hashedPassword: password1.hashedPassword,
		salt:           password1.salt,
	}
	
	// Different password
	password3, _ := NewPassword("DifferentPassword1!")

	if !password1.Equals(password2) {
		t.Errorf("Expected identical passwords to be equal")
	}

	if password1.Equals(password3) {
		t.Errorf("Expected different passwords to not be equal")
	}
}

func TestPassword_IsEmpty(t *testing.T) {
	password, _ := NewPassword("Password1!")
	emptyPassword := Password{}

	if password.IsEmpty() {
		t.Errorf("Expected non-empty password")
	}

	if !emptyPassword.IsEmpty() {
		t.Errorf("Expected empty password")
	}
}

func TestPassword_Salt(t *testing.T) {
	password, _ := NewPassword("Password1!")
	
	// Get the salt
	salt := password.Salt()
	
	// Verify it's a copy, not the original
	if &salt[0] == &password.salt[0] {
		t.Errorf("Expected Salt() to return a copy, not the original")
	}
	
	// Verify the content is the same
	if !bytes.Equal(salt, password.salt) {
		t.Errorf("Expected salt to be equal to the original")
	}
	
	// Modify the returned salt and verify the original is unchanged
	salt[0] = salt[0] + 1
	if bytes.Equal(salt, password.salt) {
		t.Errorf("Expected modifying the returned salt to not affect the original")
	}
}

func TestPassword_HashedPassword(t *testing.T) {
	password, _ := NewPassword("Password1!")
	
	// Get the hashed password
	hashedPassword := password.HashedPassword()
	
	// Verify it's the same as the internal value
	if hashedPassword != password.hashedPassword {
		t.Errorf("Expected HashedPassword() to return the internal hashed password")
	}
}

func Test_validatePassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{"Valid Password", "Password1!", false},
		{"Empty Password", "", true},
		{"Too Short Password", "Pass1!", true},
		{"No Uppercase", "password1!", true},
		{"No Digit", "Password!", true},
		{"No Special Character", "Password1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func Test_hashPassword(t *testing.T) {
	// Test that the same password and salt always produces the same hash
	password := "Password1!"
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	
	hash1 := hashPassword(password, salt)
	hash2 := hashPassword(password, salt)
	
	if hash1 != hash2 {
		t.Errorf("Expected the same password and salt to produce the same hash")
	}
	
	// Test that different passwords produce different hashes
	differentPassword := "DifferentPassword1!"
	hash3 := hashPassword(differentPassword, salt)
	
	if hash1 == hash3 {
		t.Errorf("Expected different passwords to produce different hashes")
	}
	
	// Test that different salts produce different hashes
	differentSalt := []byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	hash4 := hashPassword(password, differentSalt)
	
	if hash1 == hash4 {
		t.Errorf("Expected different salts to produce different hashes")
	}
}