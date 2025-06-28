// Copyright (c) 2025 A Bit of Help, Inc.

// Package base provides common interfaces and utilities for value objects.
package base

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
)

var (
	// ErrEmptyValue is returned when a required value is empty.
	ErrEmptyValue = errors.New("value cannot be empty")

	// ErrInvalidFormat is returned when a value has an invalid format.
	ErrInvalidFormat = errors.New("invalid format")

	// ErrOutOfRange is returned when a value is out of the allowed range.
	ErrOutOfRange = errors.New("value out of allowed range")
)

// Common validation regular expressions
var (
	// EmailRegex is a regular expression for validating email addresses.
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// URLRegex is a regular expression for validating URLs.
	URLRegex = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)

	// PhoneRegex is a regular expression for validating phone numbers.
	PhoneRegex = regexp.MustCompile(`^\+?[0-9]{1,3}[-\s]?[0-9]{1,14}$`)

	// UsernameRegex is a regular expression for validating usernames.
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)

	// CurrencyCodeRegex is a regular expression for validating ISO currency codes.
	CurrencyCodeRegex = regexp.MustCompile(`^[A-Z]{3}$`)
)

// ValidateEmail validates an email address.
func ValidateEmail(email string) error {
	// Trim whitespace
	trimmedEmail := strings.TrimSpace(email)

	// Empty email is allowed (optional field)
	if trimmedEmail == "" {
		return nil
	}

	// Check if the email contains @ symbol
	if !strings.Contains(trimmedEmail, "@") {
		return errors.New("invalid email format: missing @ symbol")
	}

	// Validate email format
	_, err := mail.ParseAddress(trimmedEmail)
	if err != nil {
		return errors.New("invalid email format")
	}

	// Additional validation with regex to restrict special characters
	// Skip regex validation for emails with names in the format "Name <email@example.com>"
	if !strings.Contains(trimmedEmail, "<") && !strings.Contains(trimmedEmail, ">") {
		if !EmailRegex.MatchString(trimmedEmail) {
			return errors.New("invalid email format: contains invalid characters")
		}
	}

	return nil
}

// ValidateURL validates a URL.
func ValidateURL(url string) error {
	// Trim whitespace
	trimmedURL := strings.TrimSpace(url)

	// Empty URL is allowed (optional field)
	if trimmedURL == "" {
		return nil
	}

	// Validate URL format
	if !URLRegex.MatchString(trimmedURL) {
		return errors.New("invalid URL format")
	}

	return nil
}

// ValidatePhone validates a phone number.
func ValidatePhone(phone string) error {
	// Trim whitespace
	trimmedPhone := strings.TrimSpace(phone)

	// Empty phone is allowed (optional field)
	if trimmedPhone == "" {
		return nil
	}

	// Check if the phone contains any letters
	if regexp.MustCompile(`[a-zA-Z]`).MatchString(trimmedPhone) {
		return errors.New("invalid phone number format: contains letters")
	}

	// Use the same regex pattern as the original implementation
	phoneRegex := regexp.MustCompile(`^[\+]?[\s]?[0-9]{0,3}[\s]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$`)

	// Validate phone format with a simple regex
	// This is a basic validation and might need to be enhanced for specific requirements
	if !phoneRegex.MatchString(trimmedPhone) {
		return errors.New("invalid phone number format")
	}

	return nil
}

// ValidateUsername validates a username.
func ValidateUsername(username string) error {
	// Trim whitespace
	trimmedUsername := strings.TrimSpace(username)

	// Empty username is not allowed
	if trimmedUsername == "" {
		return errors.New("username cannot be empty")
	}

	// Validate username format
	if !UsernameRegex.MatchString(trimmedUsername) {
		return errors.New("invalid username format: must be 3-20 characters and contain only letters, numbers, underscores, and hyphens")
	}

	return nil
}

// ValidateCurrencyCode validates an ISO currency code.
func ValidateCurrencyCode(code string) error {
	// Trim whitespace
	trimmedCode := strings.TrimSpace(code)

	// Empty code is not allowed
	if trimmedCode == "" {
		return errors.New("currency code cannot be empty")
	}

	// Convert to uppercase
	trimmedCode = strings.ToUpper(trimmedCode)

	// Validate currency code format
	if !CurrencyCodeRegex.MatchString(trimmedCode) {
		return errors.New("invalid currency code format: must be a 3-letter ISO code")
	}

	return nil
}
