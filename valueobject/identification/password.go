// Copyright (c) 2025 A Bit of Help, Inc.

// Package identification provides value objects related to identification.
package identification

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/argon2"
	"strings"
)

// Password represents a password value object
// It stores the hashed password, not the plain text
type Password struct {
	hashedPassword string
	salt           []byte
}

// NewPassword creates a new Password with validation and hashing
func NewPassword(plainPassword string) (Password, error) {
	// Trim whitespace
	trimmedPassword := strings.TrimSpace(plainPassword)

	// Validate password
	if err := validatePassword(trimmedPassword); err != nil {
		return Password{}, err
	}

	// Generate salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return Password{}, errors.New("failed to generate salt")
	}

	// Hash password
	hashedPassword := hashPassword(trimmedPassword, salt)

	return Password{
		hashedPassword: hashedPassword,
		salt:           salt,
	}, nil
}

// validatePassword checks if the password meets security requirements
func validatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check for at least one uppercase letter
	hasUpper := false
	for _, char := range password {
		if strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZ", char) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check for at least one digit
	hasDigit := false
	for _, char := range password {
		if strings.ContainsRune("0123456789", char) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}

	// Check for at least one special character
	hasSpecial := false
	specialChars := "!@#$%^&*()-_=+[]{}|;:,.<>?/~`"
	for _, char := range password {
		if strings.ContainsRune(specialChars, char) {
			hasSpecial = true
			break
		}
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// hashPassword hashes the password using Argon2id
func hashPassword(password string, salt []byte) string {
	// Argon2id parameters
	time := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)

	// Hash the password
	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// Encode as base64
	return base64.StdEncoding.EncodeToString(hash)
}

// Verify checks if the provided plain password matches the stored hashed password
func (p Password) Verify(plainPassword string) bool {
	// Trim whitespace, just like in NewPassword
	trimmedPassword := strings.TrimSpace(plainPassword)

	// Hash the provided password with the same salt
	hashedInput := hashPassword(trimmedPassword, p.salt)

	// Compare hashes in constant time to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(p.hashedPassword), []byte(hashedInput)) == 1
}

// String returns a masked representation of the password
// Never returns the actual password or hash
func (p Password) String() string {
	return "********"
}

// Equals checks if two Password values are equal
// This compares the underlying hashed passwords
func (p Password) Equals(other Password) bool {
	return p.hashedPassword == other.hashedPassword
}

// IsEmpty checks if the Password is empty
func (p Password) IsEmpty() bool {
	return p.hashedPassword == ""
}

// Salt returns a copy of the salt used for hashing
func (p Password) Salt() []byte {
	saltCopy := make([]byte, len(p.salt))
	copy(saltCopy, p.salt)
	return saltCopy
}

// HashedPassword returns the base64-encoded hashed password
func (p Password) HashedPassword() string {
	return p.hashedPassword
}