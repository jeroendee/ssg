package wordcount

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
)

// Count returns the number of words in the given markdown text.
// It excludes code blocks from the word count.
func Count(text string) int {
	if text == "" || strings.TrimSpace(text) == "" {
		return 0
	}

	// Convert markdown to HTML using goldmark
	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert([]byte(text), &buf); err != nil {
		return 0
	}

	html := buf.String()

	// Remove <pre>...</pre> and <code>...</code> blocks entirely
	html = removeCodeBlocks(html)

	// Strip remaining HTML tags
	html = stripHTMLTags(html)

	// Count words using strings.Fields on plain text
	return len(strings.Fields(html))
}

// removeCodeBlocks removes <pre>...</pre> blocks entirely (fenced code blocks).
// Inline <code> tags are preserved (will be stripped later with other HTML tags).
func removeCodeBlocks(html string) string {
	preRegex := regexp.MustCompile(`(?s)<pre>.*?</pre>`)
	return preRegex.ReplaceAllString(html, "")
}

// stripHTMLTags removes all HTML tags from the input.
func stripHTMLTags(html string) string {
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	return tagRegex.ReplaceAllString(html, " ")
}
