package listeners

import (
	"htruong/kt-crawler/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChapterScanListener_parseInnerBlock(t *testing.T) {
	l := &ChapterScanListener{}

	tests := []struct {
		name     string
		input    string
		expected internal.InnerBlock
	}{
		{
			name:  "no a tags",
			input: "This is a simple text.",
			expected: internal.InnerBlock{
				Content:    "This is a simple text.",
				References: []string{},
			},
		},
		{
			name:  "one a tag",
			input: `This is a text with <a data-toggle="tooltip" data-placement="bottom" title="Reference 1">⚓</a>.`,
			expected: internal.InnerBlock{
				Content:    `This is a text with [⚓](tooltip:{"title":"Reference 1"}).`,
				References: []string{"Reference 1"},
			},
		},
		{
			name:  "multiple a tags",
			input: `Text with <a data-toggle="tooltip" data-placement="bottom" title="Ref 1">⚓</a> and <a data-toggle="tooltip" data-placement="bottom" title="Ref 2; Ref 3">⚓</a>.`,
			expected: internal.InnerBlock{
				Content:    `Text with [⚓](tooltip:{"title":"Ref 1"}) and [⚓](tooltip:{"title":"Ref 2; Ref 3"}).`,
				References: []string{"Ref 1", "Ref 2", "Ref 3"},
			},
		},
		{
			name:  "a tag with special characters in title",
			input: `Text with <a data-toggle="tooltip" data-placement="bottom" title="title with &quot;quotes&quot;">⚓</a>.`,
			expected: internal.InnerBlock{
				Content:    `Text with [⚓](tooltip:{"title":"title with \"quotes\""}).`,
				References: []string{`title with "quotes"`},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := l.parseInnerBlock(tt.input)
			assert.Equal(t, tt.expected.Content, actual.Content)
			assert.ElementsMatch(t, tt.expected.References, actual.References)
		})
	}
}
