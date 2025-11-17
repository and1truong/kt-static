package internal

import (
	"regexp"
)

// CleanMarkdownContent applies several cleanup rules to markdown content:
// 1. Ensures there's exactly one blank line after header lines (e.g., # Header).
// 2. Replaces two or more consecutive blank lines with a single blank line.
func CleanMarkdownContent(content string) string {
	// Rule 1: Ensure one blank line after header lines.
	// This regex looks for a header line (starting with #) followed by non-blank lines,
	// and ensures there's exactly one newline after the header.
	// It avoids adding extra newlines if there's already one or more.
	reHeaderBlankLine := regexp.MustCompile(`(?m)^(\s*#+\s*.*?)(\n{2,}|$)`)
	content = reHeaderBlankLine.ReplaceAllString(content, "$1\n\n")

	// Rule 2: Replace 2 (or 3 or more) blank lines to only one.
	// This regex looks for two or more consecutive newline characters and replaces them with two newlines (one blank line).
	reMultipleBlankLines := regexp.MustCompile(`\n\n\n+`)
	content = reMultipleBlankLines.ReplaceAllString(content, "\n\n")

	return content
}
