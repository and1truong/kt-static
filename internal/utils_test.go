package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanMarkdownContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "header with no blank lines after",
			input:    "# Header\nSome content",
			expected: "# Header\n\nSome content",
		},
		{
			name:     "header with one blank line after",
			input:    "# Header\n\nSome content",
			expected: "# Header\n\nSome content",
		},
		{
			name:     "header with multiple blank lines after",
			input:    "# Header\n\n\nSome content",
			expected: "# Header\n\nSome content",
		},
		{
			name:     "multiple blank lines in content",
			input:    "Line 1\n\n\nLine 2",
			expected: "Line 1\n\nLine 2",
		},
		{
			name:     "header and multiple blank lines",
			input:    "# Header\n\n\nLine 1\n\n\n\nLine 2",
			expected: "# Header\n\nLine 1\n\nLine 2",
		},
		{
			name:     "no headers or multiple blank lines",
			input:    "Line 1\nLine 2",
			expected: "Line 1\nLine 2",
		},
		{
			name:     "content with multiple headers",
			input:    "# Header 1\nLine 1\n## Header 2\n\n\nLine 2",
			expected: "# Header 1\n\nLine 1\n## Header 2\n\nLine 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CleanMarkdownContent(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
