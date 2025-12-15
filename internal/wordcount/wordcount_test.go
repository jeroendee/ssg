package wordcount_test

import (
	"testing"

	"github.com/jeroendee/ssg/internal/wordcount"
)

func TestCount(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{
			name: "empty string returns 0",
			text: "",
			want: 0,
		},
		{
			name: "whitespace only returns 0",
			text: "   \t\n  ",
			want: 0,
		},
		{
			name: "single word returns 1",
			text: "Hello",
			want: 1,
		},
		{
			name: "spec example returns 5",
			text: "Hello! My name is Jeroen.",
			want: 5,
		},
		{
			name: "multiple whitespace returns 2",
			text: "Hello    World",
			want: 2,
		},
		{
			name: "heading notation excluded",
			text: "## Heading",
			want: 1,
		},
		{
			name: "unordered list notation excluded",
			text: "- list item",
			want: 2,
		},
		{
			name: "asterisk list notation excluded",
			text: "* another item",
			want: 2,
		},
		{
			name: "ordered list notation excluded",
			text: "1. numbered item",
			want: 2,
		},
		{
			name: "bold notation excluded",
			text: "**bold text**",
			want: 2,
		},
		{
			name: "italic notation excluded",
			text: "_italic text_",
			want: 2,
		},
		{
			name: "link notation excluded, only link text counted",
			text: "[link text](url)",
			want: 2,
		},
		{
			name: "inline code notation excluded",
			text: "`inline code`",
			want: 2,
		},
		{
			name: "fenced code block excluded entirely",
			text: "```\ncode here\n```",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := wordcount.Count(tt.text)

			if got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}
