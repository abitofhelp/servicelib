// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"testing"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		name        string
		major       int
		minor       int
		patch       int
		preRelease  string
		build       string
		expectError bool
	}{
		{"Valid Version", 1, 2, 3, "", "", false},
		{"Valid Version with PreRelease", 1, 2, 3, "alpha", "", false},
		{"Valid Version with Build", 1, 2, 3, "", "build123", false},
		{"Valid Version with PreRelease and Build", 1, 2, 3, "beta", "build456", false},
		{"Valid Version with Complex PreRelease", 1, 2, 3, "alpha.1.beta", "", false},
		{"Valid Version with Complex Build", 1, 2, 3, "", "build.123.456", false},
		{"Invalid Negative Major", -1, 2, 3, "", "", true},
		{"Invalid Negative Minor", 1, -2, 3, "", "", true},
		{"Invalid Negative Patch", 1, 2, -3, "", "", true},
		{"Invalid PreRelease Format", 1, 2, 3, "alpha!", "", true},
		{"Invalid Build Format", 1, 2, 3, "", "build!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := NewVersion(tt.major, tt.minor, tt.patch, tt.preRelease, tt.build)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if version.major != tt.major {
					t.Errorf("Expected major %d, got %d", tt.major, version.major)
				}
				if version.minor != tt.minor {
					t.Errorf("Expected minor %d, got %d", tt.minor, version.minor)
				}
				if version.patch != tt.patch {
					t.Errorf("Expected patch %d, got %d", tt.patch, version.patch)
				}
				if version.preRelease != tt.preRelease {
					t.Errorf("Expected preRelease %s, got %s", tt.preRelease, version.preRelease)
				}
				if version.build != tt.build {
					t.Errorf("Expected build %s, got %s", tt.build, version.build)
				}
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		major       int
		minor       int
		patch       int
		preRelease  string
		build       string
		expectError bool
	}{
		{"Valid Version", "1.2.3", 1, 2, 3, "", "", false},
		{"Valid Version with v prefix", "v1.2.3", 1, 2, 3, "", "", false},
		{"Valid Version with PreRelease", "1.2.3-alpha", 1, 2, 3, "alpha", "", false},
		{"Valid Version with Build", "1.2.3+build123", 1, 2, 3, "", "build123", false},
		{"Valid Version with PreRelease and Build", "1.2.3-beta+build456", 1, 2, 3, "beta", "build456", false},
		{"Valid Version with Complex PreRelease", "1.2.3-alpha.1.beta", 1, 2, 3, "alpha.1.beta", "", false},
		{"Valid Version with Complex Build", "1.2.3+build.123.456", 1, 2, 3, "", "build.123.456", false},
		{"Valid Version with Spaces", " 1.2.3 ", 1, 2, 3, "", "", false},
		{"Empty String", "", 0, 0, 0, "", "", true},
		{"Invalid Format", "1.2", 0, 0, 0, "", "", true},
		{"Invalid Format with Letters", "a.b.c", 0, 0, 0, "", "", true},
		{"Invalid PreRelease Format", "1.2.3-alpha!", 0, 0, 0, "", "", true},
		{"Invalid Build Format", "1.2.3+build!", 0, 0, 0, "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := ParseVersion(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if version.major != tt.major {
					t.Errorf("Expected major %d, got %d", tt.major, version.major)
				}
				if version.minor != tt.minor {
					t.Errorf("Expected minor %d, got %d", tt.minor, version.minor)
				}
				if version.patch != tt.patch {
					t.Errorf("Expected patch %d, got %d", tt.patch, version.patch)
				}
				if version.preRelease != tt.preRelease {
					t.Errorf("Expected preRelease %s, got %s", tt.preRelease, version.preRelease)
				}
				if version.build != tt.build {
					t.Errorf("Expected build %s, got %s", tt.build, version.build)
				}
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name       string
		major      int
		minor      int
		patch      int
		preRelease string
		build      string
		expected   string
	}{
		{"Simple Version", 1, 2, 3, "", "", "1.2.3"},
		{"Version with PreRelease", 1, 2, 3, "alpha", "", "1.2.3-alpha"},
		{"Version with Build", 1, 2, 3, "", "build123", "1.2.3+build123"},
		{"Version with PreRelease and Build", 1, 2, 3, "beta", "build456", "1.2.3-beta+build456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, _ := NewVersion(tt.major, tt.minor, tt.patch, tt.preRelease, tt.build)
			if version.String() != tt.expected {
				t.Errorf("Expected string %s, got %s", tt.expected, version.String())
			}
		})
	}
}

func TestVersion_Equals(t *testing.T) {
	tests := []struct {
		name        string
		v1Major     int
		v1Minor     int
		v1Patch     int
		v1Pre       string
		v1Build     string
		v2Major     int
		v2Minor     int
		v2Patch     int
		v2Pre       string
		v2Build     string
		shouldEqual bool
	}{
		{"Same Version", 1, 2, 3, "", "", 1, 2, 3, "", "", true},
		{"Different Major", 2, 2, 3, "", "", 1, 2, 3, "", "", false},
		{"Different Minor", 1, 3, 3, "", "", 1, 2, 3, "", "", false},
		{"Different Patch", 1, 2, 4, "", "", 1, 2, 3, "", "", false},
		{"Different PreRelease", 1, 2, 3, "alpha", "", 1, 2, 3, "beta", "", false},
		{"Different Build", 1, 2, 3, "", "build1", 1, 2, 3, "", "build2", false},
		{"One with PreRelease", 1, 2, 3, "alpha", "", 1, 2, 3, "", "", false},
		{"One with Build", 1, 2, 3, "", "build", 1, 2, 3, "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, _ := NewVersion(tt.v1Major, tt.v1Minor, tt.v1Patch, tt.v1Pre, tt.v1Build)
			v2, _ := NewVersion(tt.v2Major, tt.v2Minor, tt.v2Patch, tt.v2Pre, tt.v2Build)

			if v1.Equals(v2) != tt.shouldEqual {
				t.Errorf("Expected Equals to return %v for %s and %s", tt.shouldEqual, v1.String(), v2.String())
			}
		})
	}
}

func TestVersion_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		major    int
		minor    int
		patch    int
		pre      string
		build    string
		expected bool
	}{
		{"Empty Version", 0, 0, 0, "", "", true},
		{"Non-Empty Major", 1, 0, 0, "", "", false},
		{"Non-Empty Minor", 0, 1, 0, "", "", false},
		{"Non-Empty Patch", 0, 0, 1, "", "", false},
		{"Non-Empty PreRelease", 0, 0, 0, "alpha", "", false},
		{"Non-Empty Build", 0, 0, 0, "", "build", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var version Version
			if tt.major == 0 && tt.minor == 0 && tt.patch == 0 && tt.pre == "" && tt.build == "" {
				version = Version{} // Empty version
			} else {
				version, _ = NewVersion(tt.major, tt.minor, tt.patch, tt.pre, tt.build)
			}

			if version.IsEmpty() != tt.expected {
				t.Errorf("Expected IsEmpty to return %v for %s", tt.expected, version.String())
			}
		})
	}
}

func TestVersion_Getters(t *testing.T) {
	major, minor, patch := 1, 2, 3
	preRelease, build := "alpha", "build123"

	version, _ := NewVersion(major, minor, patch, preRelease, build)

	if version.Major() != major {
		t.Errorf("Expected Major() to return %d, got %d", major, version.Major())
	}

	if version.Minor() != minor {
		t.Errorf("Expected Minor() to return %d, got %d", minor, version.Minor())
	}

	if version.Patch() != patch {
		t.Errorf("Expected Patch() to return %d, got %d", patch, version.Patch())
	}

	if version.PreRelease() != preRelease {
		t.Errorf("Expected PreRelease() to return %s, got %s", preRelease, version.PreRelease())
	}

	if version.Build() != build {
		t.Errorf("Expected Build() to return %s, got %s", build, version.Build())
	}
}

func TestVersion_IsPreRelease(t *testing.T) {
	tests := []struct {
		name     string
		pre      string
		expected bool
	}{
		{"With PreRelease", "alpha", true},
		{"Without PreRelease", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, _ := NewVersion(1, 2, 3, tt.pre, "")
			if version.IsPreRelease() != tt.expected {
				t.Errorf("Expected IsPreRelease to return %v for %s", tt.expected, version.String())
			}
		})
	}
}

func TestVersion_CompareTo(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		// Major version comparison
		{"Major Greater", "2.0.0", "1.0.0", 1},
		{"Major Less", "1.0.0", "2.0.0", -1},

		// Minor version comparison
		{"Minor Greater", "1.2.0", "1.1.0", 1},
		{"Minor Less", "1.1.0", "1.2.0", -1},

		// Patch version comparison
		{"Patch Greater", "1.1.2", "1.1.1", 1},
		{"Patch Less", "1.1.1", "1.1.2", -1},

		// Pre-release comparison
		{"Normal vs Pre-release", "1.0.0", "1.0.0-alpha", 1},
		{"Pre-release vs Normal", "1.0.0-alpha", "1.0.0", -1},
		{"Different Pre-release - Alphabetical", "1.0.0-alpha", "1.0.0-beta", -1},
		{"Different Pre-release - Reverse Alphabetical", "1.0.0-beta", "1.0.0-alpha", 1},

		// Numeric identifiers in pre-release
		{"Numeric Pre-release Comparison", "1.0.0-alpha.1", "1.0.0-alpha.2", -1},
		{"Numeric vs Non-numeric Pre-release", "1.0.0-1", "1.0.0-alpha", -1},
		{"Non-numeric vs Numeric Pre-release", "1.0.0-alpha", "1.0.0-1", 1},

		// Pre-release with different number of identifiers
		{"More Pre-release Identifiers", "1.0.0-alpha.beta.1", "1.0.0-alpha.beta", 1},
		{"Fewer Pre-release Identifiers", "1.0.0-alpha.beta", "1.0.0-alpha.beta.1", -1},

		// Equal versions
		{"Equal Versions", "1.0.0", "1.0.0", 0},
		{"Equal Versions with Pre-release", "1.0.0-alpha", "1.0.0-alpha", 0},

		// Build metadata (should not affect comparison)
		{"Equal with Different Build", "1.0.0+build1", "1.0.0+build2", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, _ := ParseVersion(tt.v1)
			v2, _ := ParseVersion(tt.v2)

			result := v1.CompareTo(v2)
			if result != tt.expected {
				t.Errorf("Expected CompareTo to return %d for %s and %s, got %d",
					tt.expected, tt.v1, tt.v2, result)
			}
		})
	}
}

func TestVersion_NextMajor(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{"From 1.2.3", "1.2.3", "2.0.0"},
		{"From 1.2.3-alpha", "1.2.3-alpha", "2.0.0"},
		{"From 1.2.3+build", "1.2.3+build", "2.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, _ := ParseVersion(tt.version)
			next := version.NextMajor()

			if next.String() != tt.expected {
				t.Errorf("Expected NextMajor to return %s for %s, got %s",
					tt.expected, tt.version, next.String())
			}
		})
	}
}

func TestVersion_NextMinor(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{"From 1.2.3", "1.2.3", "1.3.0"},
		{"From 1.2.3-alpha", "1.2.3-alpha", "1.3.0"},
		{"From 1.2.3+build", "1.2.3+build", "1.3.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, _ := ParseVersion(tt.version)
			next := version.NextMinor()

			if next.String() != tt.expected {
				t.Errorf("Expected NextMinor to return %s for %s, got %s",
					tt.expected, tt.version, next.String())
			}
		})
	}
}

func TestVersion_NextPatch(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{"From 1.2.3", "1.2.3", "1.2.4"},
		{"From 1.2.3-alpha", "1.2.3-alpha", "1.2.4"},
		{"From 1.2.3+build", "1.2.3+build", "1.2.4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, _ := ParseVersion(tt.version)
			next := version.NextPatch()

			if next.String() != tt.expected {
				t.Errorf("Expected NextPatch to return %s for %s, got %s",
					tt.expected, tt.version, next.String())
			}
		})
	}
}

func Test_comparePreRelease(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected int
	}{
		{"Alphabetical Order", "alpha", "beta", -1},
		{"Reverse Alphabetical Order", "beta", "alpha", 1},
		{"Same String", "alpha", "alpha", 0},
		{"Numeric Order", "1", "2", -1},
		{"Numeric vs Non-numeric", "1", "alpha", -1},
		{"Non-numeric vs Numeric", "alpha", "1", 1},
		{"More Identifiers", "alpha.beta", "alpha", 1},
		{"Fewer Identifiers", "alpha", "alpha.beta", -1},
		{"Mixed Comparison", "alpha.1", "alpha.beta", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := comparePreRelease(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Expected comparePreRelease to return %d for %s and %s, got %d",
					tt.expected, tt.a, tt.b, result)
			}
		})
	}
}
