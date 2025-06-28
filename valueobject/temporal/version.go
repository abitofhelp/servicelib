// Copyright (c) 2025 A Bit of Help, Inc.

// Package temporal provides value objects related to time and versioning information.
package temporal

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a semantic version (major.minor.patch)
type Version struct {
	major      int
	minor      int
	patch      int
	preRelease string
	build      string
}

// Regular expression for validating semantic version format
var versionRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// NewVersion creates a new Version with validation
func NewVersion(major, minor, patch int, preRelease, build string) (Version, error) {
	// Validate version components
	if major < 0 {
		return Version{}, errors.New("major version cannot be negative")
	}
	if minor < 0 {
		return Version{}, errors.New("minor version cannot be negative")
	}
	if patch < 0 {
		return Version{}, errors.New("patch version cannot be negative")
	}

	// Validate pre-release format if provided
	if preRelease != "" {
		if !regexp.MustCompile(`^[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*$`).MatchString(preRelease) {
			return Version{}, errors.New("invalid pre-release format")
		}
	}

	// Validate build metadata format if provided
	if build != "" {
		if !regexp.MustCompile(`^[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*$`).MatchString(build) {
			return Version{}, errors.New("invalid build metadata format")
		}
	}

	return Version{
		major:      major,
		minor:      minor,
		patch:      patch,
		preRelease: preRelease,
		build:      build,
	}, nil
}

// ParseVersion creates a new Version from a string in format "major.minor.patch[-prerelease][+build]"
func ParseVersion(s string) (Version, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(s)

	// Empty string is not allowed
	if trimmed == "" {
		return Version{}, errors.New("version string cannot be empty")
	}

	// Match against regex
	matches := versionRegex.FindStringSubmatch(trimmed)
	if matches == nil {
		return Version{}, errors.New("invalid version format, expected 'major.minor.patch[-prerelease][+build]'")
	}

	// Parse major version
	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return Version{}, errors.New("invalid major version")
	}

	// Parse minor version
	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return Version{}, errors.New("invalid minor version")
	}

	// Parse patch version
	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return Version{}, errors.New("invalid patch version")
	}

	// Get pre-release and build metadata
	preRelease := ""
	if len(matches) > 4 && matches[4] != "" {
		preRelease = matches[4]
	}

	build := ""
	if len(matches) > 5 && matches[5] != "" {
		build = matches[5]
	}

	return NewVersion(major, minor, patch, preRelease, build)
}

// String returns the string representation of the Version
func (v Version) String() string {
	result := fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	if v.preRelease != "" {
		result += "-" + v.preRelease
	}
	if v.build != "" {
		result += "+" + v.build
	}
	return result
}

// Equals checks if two Versions are equal
func (v Version) Equals(other Version) bool {
	return v.major == other.major &&
		v.minor == other.minor &&
		v.patch == other.patch &&
		v.preRelease == other.preRelease &&
		v.build == other.build
}

// IsEmpty checks if the Version is empty (zero value)
func (v Version) IsEmpty() bool {
	return v.major == 0 && v.minor == 0 && v.patch == 0 && v.preRelease == "" && v.build == ""
}

// Validate checks if the Version is valid
func (v Version) Validate() error {
	// Validate version components
	if v.major < 0 {
		return errors.New("major version cannot be negative")
	}
	if v.minor < 0 {
		return errors.New("minor version cannot be negative")
	}
	if v.patch < 0 {
		return errors.New("patch version cannot be negative")
	}

	// Validate pre-release format if provided
	if v.preRelease != "" {
		if !regexp.MustCompile(`^[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*$`).MatchString(v.preRelease) {
			return errors.New("invalid pre-release format")
		}
	}

	// Validate build metadata format if provided
	if v.build != "" {
		if !regexp.MustCompile(`^[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*$`).MatchString(v.build) {
			return errors.New("invalid build metadata format")
		}
	}

	return nil
}

// Major returns the major version number
func (v Version) Major() int {
	return v.major
}

// Minor returns the minor version number
func (v Version) Minor() int {
	return v.minor
}

// Patch returns the patch version number
func (v Version) Patch() int {
	return v.patch
}

// PreRelease returns the pre-release identifier
func (v Version) PreRelease() string {
	return v.preRelease
}

// Build returns the build metadata
func (v Version) Build() string {
	return v.build
}

// IsPreRelease checks if this is a pre-release version
func (v Version) IsPreRelease() bool {
	return v.preRelease != ""
}

// CompareTo compares this version to another version
// Returns:
//
//	-1 if this version is less than the other
//	 0 if this version is equal to the other
//	 1 if this version is greater than the other
func (v Version) CompareTo(other Version) int {
	// Compare major version
	if v.major < other.major {
		return -1
	}
	if v.major > other.major {
		return 1
	}

	// Compare minor version
	if v.minor < other.minor {
		return -1
	}
	if v.minor > other.minor {
		return 1
	}

	// Compare patch version
	if v.patch < other.patch {
		return -1
	}
	if v.patch > other.patch {
		return 1
	}

	// Compare pre-release (pre-release versions are less than the associated normal version)
	if v.preRelease == "" && other.preRelease != "" {
		return 1
	}
	if v.preRelease != "" && other.preRelease == "" {
		return -1
	}
	if v.preRelease != other.preRelease {
		// Compare pre-release identifiers
		return comparePreRelease(v.preRelease, other.preRelease)
	}

	// Versions are equal (build metadata does not affect precedence)
	return 0
}

// comparePreRelease compares two pre-release strings according to SemVer rules
func comparePreRelease(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	// Compare each identifier
	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		aIsNum := true
		aNum, aErr := strconv.Atoi(aParts[i])
		if aErr != nil {
			aIsNum = false
		}

		bIsNum := true
		bNum, bErr := strconv.Atoi(bParts[i])
		if bErr != nil {
			bIsNum = false
		}

		// Numeric identifiers always have lower precedence than non-numeric identifiers
		if aIsNum && !bIsNum {
			return -1
		}
		if !aIsNum && bIsNum {
			return 1
		}

		// If both are numeric, compare numerically
		if aIsNum && bIsNum {
			if aNum < bNum {
				return -1
			}
			if aNum > bNum {
				return 1
			}
			continue
		}

		// If both are non-numeric, compare lexically
		if aParts[i] < bParts[i] {
			return -1
		}
		if aParts[i] > bParts[i] {
			return 1
		}
	}

	// If all identifiers so far are equal, the one with more identifiers has higher precedence
	if len(aParts) < len(bParts) {
		return -1
	}
	if len(aParts) > len(bParts) {
		return 1
	}

	return 0
}

// NextMajor returns the next major version
func (v Version) NextMajor() Version {
	result, _ := NewVersion(v.major+1, 0, 0, "", "")
	return result
}

// NextMinor returns the next minor version
func (v Version) NextMinor() Version {
	result, _ := NewVersion(v.major, v.minor+1, 0, "", "")
	return result
}

// NextPatch returns the next patch version
func (v Version) NextPatch() Version {
	result, _ := NewVersion(v.major, v.minor, v.patch+1, "", "")
	return result
}

// ToMap converts the Version to a map[string]interface{}
func (v Version) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"major":      v.major,
		"minor":      v.minor,
		"patch":      v.patch,
		"preRelease": v.preRelease,
		"build":      v.build,
	}
}