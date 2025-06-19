// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"net/url"
	"strings"
)

// URL represents a URL value object
type URL string

// NewURL creates a new URL with validation
func NewURL(rawURL string) (URL, error) {
	// Trim whitespace
	trimmedURL := strings.TrimSpace(rawURL)

	// Empty URL is allowed (optional field)
	if trimmedURL == "" {
		return "", nil
	}

	// Parse and validate URL
	parsedURL, err := url.Parse(trimmedURL)
	if err != nil {
		return "", errors.New("invalid URL format")
	}

	// Ensure URL has a scheme and host
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errors.New("URL must have a scheme and host")
	}

	// Ensure scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", errors.New("URL scheme must be http or https")
	}

	return URL(trimmedURL), nil
}

// String returns the string representation of the URL
func (u URL) String() string {
	return string(u)
}

// Equals checks if two URLs are equal
func (u URL) Equals(other URL) bool {
	// Parse both URLs for comparison
	parsedThis, err1 := url.Parse(string(u))
	parsedOther, err2 := url.Parse(string(other))
	
	// If either URL can't be parsed, fall back to string comparison
	if err1 != nil || err2 != nil {
		return string(u) == string(other)
	}
	
	// Compare normalized URLs
	return parsedThis.String() == parsedOther.String()
}

// IsEmpty checks if the URL is empty
func (u URL) IsEmpty() bool {
	return u == ""
}

// Domain returns the domain part of the URL
func (u URL) Domain() (string, error) {
	if u.IsEmpty() {
		return "", errors.New("URL is empty")
	}
	
	parsedURL, err := url.Parse(string(u))
	if err != nil {
		return "", errors.New("invalid URL format")
	}
	
	return parsedURL.Host, nil
}

// Path returns the path part of the URL
func (u URL) Path() (string, error) {
	if u.IsEmpty() {
		return "", errors.New("URL is empty")
	}
	
	parsedURL, err := url.Parse(string(u))
	if err != nil {
		return "", errors.New("invalid URL format")
	}
	
	return parsedURL.Path, nil
}

// Query returns the query part of the URL
func (u URL) Query() (url.Values, error) {
	if u.IsEmpty() {
		return nil, errors.New("URL is empty")
	}
	
	parsedURL, err := url.Parse(string(u))
	if err != nil {
		return nil, errors.New("invalid URL format")
	}
	
	return parsedURL.Query(), nil
}