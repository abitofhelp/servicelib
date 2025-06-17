// Copyright (c) 2025 A Bit of Help, Inc.

package stringutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPrefixIgnoreCase(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		prefix   string
		expected bool
	}{
		{
			name:     "empty string and prefix",
			s:        "",
			prefix:   "",
			expected: true,
		},
		{
			name:     "empty string with non-empty prefix",
			s:        "",
			prefix:   "prefix",
			expected: false,
		},
		{
			name:     "non-empty string with empty prefix",
			s:        "string",
			prefix:   "",
			expected: true,
		},
		{
			name:     "matching prefix same case",
			s:        "Hello World",
			prefix:   "Hello",
			expected: true,
		},
		{
			name:     "matching prefix different case",
			s:        "Hello World",
			prefix:   "hello",
			expected: true,
		},
		{
			name:     "non-matching prefix",
			s:        "Hello World",
			prefix:   "World",
			expected: false,
		},
		{
			name:     "prefix longer than string",
			s:        "Hi",
			prefix:   "Hello",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasPrefixIgnoreCase(tt.s, tt.prefix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "empty string and substring",
			s:        "",
			substr:   "",
			expected: true,
		},
		{
			name:     "empty string with non-empty substring",
			s:        "",
			substr:   "substr",
			expected: false,
		},
		{
			name:     "non-empty string with empty substring",
			s:        "string",
			substr:   "",
			expected: true,
		},
		{
			name:     "matching substring same case",
			s:        "Hello World",
			substr:   "World",
			expected: true,
		},
		{
			name:     "matching substring different case",
			s:        "Hello World",
			substr:   "world",
			expected: true,
		},
		{
			name:     "non-matching substring",
			s:        "Hello World",
			substr:   "Universe",
			expected: false,
		},
		{
			name:     "substring in middle",
			s:        "Hello World",
			substr:   "lo Wo",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsIgnoreCase(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasAnyPrefix(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		prefixes []string
		expected bool
	}{
		{
			name:     "empty string and no prefixes",
			s:        "",
			prefixes: []string{},
			expected: false,
		},
		{
			name:     "empty string with prefixes",
			s:        "",
			prefixes: []string{"prefix1", "prefix2"},
			expected: false,
		},
		{
			name:     "empty string with empty prefix",
			s:        "",
			prefixes: []string{""},
			expected: true,
		},
		{
			name:     "string with matching first prefix",
			s:        "Hello World",
			prefixes: []string{"Hello", "Hi"},
			expected: true,
		},
		{
			name:     "string with matching second prefix",
			s:        "Hello World",
			prefixes: []string{"Hi", "Hello"},
			expected: true,
		},
		{
			name:     "string with no matching prefixes",
			s:        "Hello World",
			prefixes: []string{"Hi", "Hey"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasAnyPrefix(tt.s, tt.prefixes...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToLowerCase(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{
			name:     "empty string",
			s:        "",
			expected: "",
		},
		{
			name:     "lowercase string",
			s:        "hello world",
			expected: "hello world",
		},
		{
			name:     "uppercase string",
			s:        "HELLO WORLD",
			expected: "hello world",
		},
		{
			name:     "mixed case string",
			s:        "Hello World",
			expected: "hello world",
		},
		{
			name:     "string with numbers and symbols",
			s:        "Hello123!@#",
			expected: "hello123!@#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToLowerCase(tt.s)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJoinWithAnd(t *testing.T) {
	tests := []struct {
		name           string
		items          []string
		useOxfordComma bool
		expected       string
	}{
		{
			name:           "empty slice",
			items:          []string{},
			useOxfordComma: false,
			expected:       "",
		},
		{
			name:           "single item",
			items:          []string{"apple"},
			useOxfordComma: false,
			expected:       "apple",
		},
		{
			name:           "two items without Oxford comma",
			items:          []string{"apple", "banana"},
			useOxfordComma: false,
			expected:       "apple and banana",
		},
		{
			name:           "two items with Oxford comma (should be the same)",
			items:          []string{"apple", "banana"},
			useOxfordComma: true,
			expected:       "apple and banana",
		},
		{
			name:           "three items without Oxford comma",
			items:          []string{"apple", "banana", "cherry"},
			useOxfordComma: false,
			expected:       "apple, banana and cherry",
		},
		{
			name:           "three items with Oxford comma",
			items:          []string{"apple", "banana", "cherry"},
			useOxfordComma: true,
			expected:       "apple, banana, and cherry",
		},
		{
			name:           "four items without Oxford comma",
			items:          []string{"apple", "banana", "cherry", "date"},
			useOxfordComma: false,
			expected:       "apple, banana, cherry and date",
		},
		{
			name:           "four items with Oxford comma",
			items:          []string{"apple", "banana", "cherry", "date"},
			useOxfordComma: true,
			expected:       "apple, banana, cherry, and date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinWithAnd(tt.items, tt.useOxfordComma)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{
			name:     "empty string",
			s:        "",
			expected: true,
		},
		{
			name:     "whitespace only",
			s:        "   \t\n",
			expected: true,
		},
		{
			name:     "non-empty string",
			s:        "hello",
			expected: false,
		},
		{
			name:     "string with whitespace",
			s:        "  hello  ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmpty(tt.s)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{
			name:     "empty string",
			s:        "",
			expected: false,
		},
		{
			name:     "whitespace only",
			s:        "   \t\n",
			expected: false,
		},
		{
			name:     "non-empty string",
			s:        "hello",
			expected: true,
		},
		{
			name:     "string with whitespace",
			s:        "  hello  ",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotEmpty(tt.s)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		maxLength int
		expected  string
	}{
		{
			name:      "empty string",
			s:         "",
			maxLength: 10,
			expected:  "",
		},
		{
			name:      "string shorter than max length",
			s:         "hello",
			maxLength: 10,
			expected:  "hello",
		},
		{
			name:      "string equal to max length",
			s:         "hello",
			maxLength: 5,
			expected:  "hello",
		},
		{
			name:      "string longer than max length",
			s:         "hello world",
			maxLength: 5,
			expected:  "hello...",
		},
		{
			name:      "max length of 0",
			s:         "hello",
			maxLength: 0,
			expected:  "...",
		},
		{
			name:      "negative max length",
			s:         "hello",
			maxLength: -1,
			expected:  "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Truncate(tt.s, tt.maxLength)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRemoveWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{
			name:     "empty string",
			s:        "",
			expected: "",
		},
		{
			name:     "string with no whitespace",
			s:        "hello",
			expected: "hello",
		},
		{
			name:     "string with spaces",
			s:        "hello world",
			expected: "helloworld",
		},
		{
			name:     "string with tabs and newlines",
			s:        "hello\tworld\n",
			expected: "helloworld",
		},
		{
			name:     "string with mixed whitespace",
			s:        "  hello  \t world \n",
			expected: "helloworld",
		},
		{
			name:     "only whitespace",
			s:        "   \t\n",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveWhitespace(tt.s)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestForwardSlashPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "empty string",
			path:     "",
			expected: "",
		},
		{
			name:     "path with no backslashes",
			path:     "/usr/local/bin",
			expected: "/usr/local/bin",
		},
		{
			name:     "Windows path",
			path:     "C:\\Users\\user\\Documents",
			expected: "C:/Users/user/Documents",
		},
		{
			name:     "Mixed path",
			path:     "C:\\Users/user\\Documents/file.txt",
			expected: "C:/Users/user/Documents/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ForwardSlashPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}