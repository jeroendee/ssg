package parser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jeroendee/ssg/internal/model"
	"github.com/jeroendee/ssg/internal/wordcount"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

// Frontmatter holds metadata extracted from markdown files.
type Frontmatter struct {
	Title    string `yaml:"title"`
	Template string `yaml:"template"`
	Summary  string `yaml:"summary"`
	Date     string `yaml:"date"`
}

// MarkdownToHTML converts markdown content to HTML.
// Deprecated: Use MarkdownToHTMLWithError for proper error handling.
func MarkdownToHTML(md string) string {
	html, _ := MarkdownToHTMLWithError(md)
	return html
}

// MarkdownToHTMLWithError converts markdown content to HTML and returns any conversion error.
func MarkdownToHTMLWithError(md string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return "", fmt.Errorf("markdown conversion failed: %w", err)
	}
	return buf.String(), nil
}

// ExtractFrontmatter separates YAML frontmatter from markdown content.
func ExtractFrontmatter(content string) (Frontmatter, string, error) {
	var fm Frontmatter
	if !strings.HasPrefix(content, "---") {
		return fm, content, nil
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return fm, content, nil
	}

	if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
		return fm, content, err
	}

	return fm, strings.TrimSpace(parts[2]), nil
}

// ParsePage reads a markdown file and returns a Page.
func ParsePage(path string) (*model.Page, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fm, body, err := ExtractFrontmatter(string(data))
	if err != nil {
		return nil, err
	}

	slug := strings.TrimSuffix(filepath.Base(path), ".md")
	if slug == "home" {
		slug = ""
	}
	template := fm.Template
	if template == "" {
		template = "page"
	}

	html, err := MarkdownToHTMLWithError(body)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	pagePath := "/" + slug + "/"
	if slug == "" {
		pagePath = "/"
	}

	return &model.Page{
		Title:    fm.Title,
		Slug:     slug,
		Content:  html,
		Path:     pagePath,
		Template: template,
	}, nil
}

var dateFilenameRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})-(.+)\.md$`)

// ParsePost reads a markdown file and returns a Post.
func ParsePost(path string) (*model.Post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fm, body, err := ExtractFrontmatter(string(data))
	if err != nil {
		return nil, err
	}

	filename := filepath.Base(path)
	var postDate time.Time
	slug := strings.TrimSuffix(filename, ".md")

	// Try to parse date from filename first
	if matches := dateFilenameRegex.FindStringSubmatch(filename); matches != nil {
		postDate, err = time.Parse("2006-01-02", matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid date in filename %s: %w", filename, err)
		}
		slug = matches[2]
	}

	// Frontmatter date overrides filename date
	if fm.Date != "" {
		postDate, err = time.Parse("2006-01-02", fm.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date in frontmatter: %w", err)
		}
	}

	template := fm.Template
	if template == "" {
		template = "post"
	}

	html, err := MarkdownToHTMLWithError(body)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	return &model.Post{
		Page: model.Page{
			Title:    fm.Title,
			Slug:     slug,
			Content:  html,
			Path:     "/blog/" + slug + "/",
			Template: template,
		},
		Date:      postDate,
		Summary:   fm.Summary,
		WordCount: wordcount.Count(body),
	}, nil
}
