package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchPatterns(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		patterns []string
		expected bool
	}{
		{
			name:     "Single matching pattern",
			str:      "file.txt",
			patterns: []string{"*.txt"},
			expected: true,
		},
		{
			name:     "Multiple patterns with one match",
			str:      "document.pdf",
			patterns: []string{"*.doc", "*.pdf", "*.txt"},
			expected: true,
		},
		{
			name:     "No matching patterns",
			str:      "image.png",
			patterns: []string{"*.jpg", "*.gif"},
			expected: false,
		},
		{
			name:     "Empty pattern list",
			str:      "file.txt",
			patterns: []string{},
			expected: false,
		},
		{
			name:     "Pattern with special characters",
			str:      "file[1].txt",
			patterns: []string{"file[[]*.txt"},
			expected: true,
		},
		{
			name:     "Pattern with layer 1 directory",
			str:      "prod/file.txt",
			patterns: []string{"prod/*"},
			expected: true,
		},
		{
			name:     "Pattern with layer 2 directory",
			str:      "prod/abc/file.txt",
			patterns: []string{"prod/**/*"},
			expected: true,
		},
		{
			name:     "Pattern with layer 3+ directory",
			str:      "prod/abc/ecd/file.txt",
			patterns: []string{"**/*"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchPatterns(tt.str, false, tt.patterns)
			assert.Equal(t, tt.expected, result, "matchPatterns(%q, %v) = %v, want %v", tt.str, tt.patterns, result, tt.expected)
		})
	}
}

func TestFilterStringsEdgeCases(t *testing.T) {
	tests := []struct {
		name            string
		input           []string
		includePatterns []string
		excludePatterns []string
		expected        []string
	}{
		{
			name:            "Empty input",
			input:           []string{},
			includePatterns: []string{"*"},
			excludePatterns: []string{},
			expected:        []string(nil),
		},
		{
			name:            "Empty include patterns",
			input:           []string{"file1.txt", "file2.txt"},
			includePatterns: []string{},
			excludePatterns: []string{},
			expected:        []string{"file1.txt", "file2.txt"},
		},
		{
			name:            "Conflicting include and exclude patterns",
			input:           []string{"file1.txt", "file2.txt", "file3.log"},
			includePatterns: []string{"*.txt", "*.log"},
			excludePatterns: []string{"*"},
			expected:        []string(nil),
		},
		{
			name:            "Case sensitivity",
			input:           []string{"File1.TXT", "file2.txt", "FILE3.LOG"},
			includePatterns: []string{"*.txt"},
			excludePatterns: []string{},
			expected:        []string{"file2.txt"},
		},
		{
			name:            "Complex patterns",
			input:           []string{"file1.txt", "a/b/c/d/file2.txt", "dir/file3.txt", "dir/subdir/file4.txt"},
			includePatterns: []string{"**/*.txt"},
			excludePatterns: []string{"*/subdir/*"},
			expected:        []string{"file1.txt", "a/b/c/d/file2.txt", "dir/file3.txt"},
		},
		{
			name:            "Overlapping patterns",
			input:           []string{"file1.txt", "file2.log", "file3.tmp"},
			includePatterns: []string{"*.txt", "*.log", "*"},
			excludePatterns: []string{"*.log", "*.tmp"},
			expected:        []string{"file1.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterStrings(tt.input, tt.includePatterns, tt.excludePatterns)
			assert.Equal(t, tt.expected, result, "FilterStrings(%v, %v, %v) = %v, want %v", tt.input, tt.includePatterns, tt.excludePatterns, result, tt.expected)
		})
	}
}
