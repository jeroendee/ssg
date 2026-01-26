package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jeroendee/ssg/internal/parser"
)

func TestParseMarkdown(t *testing.T) {
	t.Parallel()
	md := "# Hello\n\nThis is **bold** text."
	html, err := parser.MarkdownToHTMLWithError(md)
	if err != nil {
		t.Fatalf("MarkdownToHTMLWithError() error = %v", err)
	}
	if !strings.Contains(html, `<h1 id="hello"><a href="#hello">Hello</a></h1>`) {
		t.Errorf("expected <h1 id=\"hello\"><a href=\"#hello\">Hello</a></h1>, got %s", html)
	}
	if !strings.Contains(html, "<strong>bold</strong>") {
		t.Errorf("expected <strong>bold</strong>, got %s", html)
	}
}

func TestParsePage(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "about.md")
	content := `---
title: "About Me"
---
# About

This is the about page.`
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	page, err := parser.ParsePage(file)
	if err != nil {
		t.Fatalf("ParsePage() error = %v", err)
	}
	if page.Title != "About Me" {
		t.Errorf("Title = %q, want %q", page.Title, "About Me")
	}
	if page.Slug != "about" {
		t.Errorf("Slug = %q, want %q", page.Slug, "about")
	}
	if !strings.Contains(page.Content, `<h1 id="about"><a href="#about">About</a></h1>`) {
		t.Errorf("Content should have HTML with ID and anchor, got %q", page.Content)
	}
}

func TestParsePost_DateFromFrontmatter(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "my-post.md")
	content := `---
title: "My Post"
date: "2021-03-26"
summary: "Post summary"
---
# Post content`
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	post, err := parser.ParsePost(file)
	if err != nil {
		t.Fatalf("ParsePost() error = %v", err)
	}
	if post.Title != "My Post" {
		t.Errorf("Title = %q, want %q", post.Title, "My Post")
	}
	expectedDate := time.Date(2021, 3, 26, 0, 0, 0, 0, time.UTC)
	if !post.Date.Equal(expectedDate) {
		t.Errorf("Date = %v, want %v", post.Date, expectedDate)
	}
	if post.Summary != "Post summary" {
		t.Errorf("Summary = %q, want %q", post.Summary, "Post summary")
	}
}

func TestParsePost_DateFromFilename(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "2021-03-26-my-post.md")
	content := `---
title: "My Post"
---
# Post content`
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	post, err := parser.ParsePost(file)
	if err != nil {
		t.Fatalf("ParsePost() error = %v", err)
	}
	expectedDate := time.Date(2021, 3, 26, 0, 0, 0, 0, time.UTC)
	if !post.Date.Equal(expectedDate) {
		t.Errorf("Date = %v, want %v", post.Date, expectedDate)
	}
	if post.Slug != "my-post" {
		t.Errorf("Slug = %q, want %q", post.Slug, "my-post")
	}
}

func TestParsePage_NonExistentFile(t *testing.T) {
	t.Parallel()
	_, err := parser.ParsePage("/nonexistent/file.md")
	if err == nil {
		t.Error("ParsePage() expected error for non-existent file")
	}
}

func TestParsePost_NonExistentFile(t *testing.T) {
	t.Parallel()
	_, err := parser.ParsePost("/nonexistent/file.md")
	if err == nil {
		t.Error("ParsePost() expected error for non-existent file")
	}
}

func TestParsePost_WordCount(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "test-post.md")
	content := `---
title: "Test Post"
---
Hello! My name is Jeroen.`
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	post, err := parser.ParsePost(file)
	if err != nil {
		t.Fatalf("ParsePost() error = %v", err)
	}
	if post.WordCount != 5 {
		t.Errorf("WordCount = %d, want 5", post.WordCount)
	}
}

func TestMarkdownToHTML_ReturnsError(t *testing.T) {
	t.Parallel()
	// Valid markdown should not return error
	md := "# Hello"
	html, err := parser.MarkdownToHTMLWithError(md)
	if err != nil {
		t.Errorf("MarkdownToHTMLWithError() unexpected error = %v", err)
	}
	if !strings.Contains(html, `<h1 id="hello"><a href="#hello">Hello</a></h1>`) {
		t.Errorf("expected <h1 id=\"hello\"><a href=\"#hello\">Hello</a></h1>, got %s", html)
	}
}

func TestParsePost_InvalidFrontmatterDate(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "my-post.md")
	content := `---
title: "My Post"
date: "not-a-date"
---
# Content`
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := parser.ParsePost(file)
	if err == nil {
		t.Error("ParsePost() expected error for invalid frontmatter date")
	}
}

func TestParsePost_InvalidFilenameDate(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	// Invalid date: month 13 doesn't exist
	file := filepath.Join(dir, "2021-13-45-my-post.md")
	content := `---
title: "My Post"
---
# Content`
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := parser.ParsePost(file)
	if err == nil {
		t.Error("ParsePost() expected error for invalid filename date")
	}
}

func TestParsePage_HomeMdGetsEmptySlug(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "home.md")
	content := `---
title: "Home"
---
# Welcome

This is the homepage content.`
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	page, err := parser.ParsePage(file)
	if err != nil {
		t.Fatalf("ParsePage() error = %v", err)
	}
	if page.Slug != "" {
		t.Errorf("Slug = %q, want empty string for home.md", page.Slug)
	}
	if page.Path != "/" {
		t.Errorf("Path = %q, want %q for home.md", page.Path, "/")
	}
}

func TestExtractAssetReferences(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		markdown string
		want     []string
	}{
		{
			name:     "single asset reference",
			markdown: "Here is an image: ![diagram](assets/diagram.png)",
			want:     []string{"assets/diagram.png"},
		},
		{
			name:     "multiple asset references",
			markdown: "![first](assets/one.png)\n\nSome text\n\n![second](assets/two.jpg)",
			want:     []string{"assets/one.png", "assets/two.jpg"},
		},
		{
			name:     "no asset references",
			markdown: "Just plain text without any images.",
			want:     []string{},
		},
		{
			name:     "external URLs not matched",
			markdown: "![external](https://example.com/image.png)",
			want:     []string{},
		},
		{
			name:     "mixed internal and external",
			markdown: "![local](assets/local.png)\n![external](https://example.com/img.png)\n![another](assets/another.gif)",
			want:     []string{"assets/local.png", "assets/another.gif"},
		},
		{
			name:     "empty markdown",
			markdown: "",
			want:     []string{},
		},
		{
			name:     "asset in subdirectory",
			markdown: "![nested](assets/images/photo.jpg)",
			want:     []string{"assets/images/photo.jpg"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.ExtractAssetReferences(tt.markdown)
			if len(got) != len(tt.want) {
				t.Errorf("ExtractAssetReferences() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ExtractAssetReferences()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParsePost_PopulatesAssets(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "2024-01-15-with-assets.md")
	content := `---
title: "Post with Assets"
---
# My Post

Here's a diagram:

![diagram](assets/diagram.png)

And another image:

![photo](assets/photos/vacation.jpg)`
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	post, err := parser.ParsePost(file)
	if err != nil {
		t.Fatalf("ParsePost() error = %v", err)
	}
	wantAssets := []string{"assets/diagram.png", "assets/photos/vacation.jpg"}
	if len(post.Assets) != len(wantAssets) {
		t.Errorf("Assets = %v, want %v", post.Assets, wantAssets)
		return
	}
	for i := range post.Assets {
		if post.Assets[i] != wantAssets[i] {
			t.Errorf("Assets[%d] = %q, want %q", i, post.Assets[i], wantAssets[i])
		}
	}
}

func TestMarkdownToHTML_HeadingIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		markdown     string
		wantContains string
	}{
		{
			name:         "h1 gets ID and anchor",
			markdown:     "# Hello",
			wantContains: `<h1 id="hello"><a href="#hello">Hello</a></h1>`,
		},
		{
			name:         "h2 gets ID and anchor",
			markdown:     "## Features",
			wantContains: `<h2 id="features"><a href="#features">Features</a></h2>`,
		},
		{
			name:         "h3 gets ID and anchor",
			markdown:     "### Getting Started",
			wantContains: `<h3 id="getting-started"><a href="#getting-started">Getting Started</a></h3>`,
		},
		{
			name:         "h4 gets ID and anchor",
			markdown:     "#### Installation Steps",
			wantContains: `<h4 id="installation-steps"><a href="#installation-steps">Installation Steps</a></h4>`,
		},
		{
			name:         "h5 gets ID and anchor",
			markdown:     "##### Advanced Options",
			wantContains: `<h5 id="advanced-options"><a href="#advanced-options">Advanced Options</a></h5>`,
		},
		{
			name:         "h6 gets ID and anchor",
			markdown:     "###### Footnotes",
			wantContains: `<h6 id="footnotes"><a href="#footnotes">Footnotes</a></h6>`,
		},
		{
			name:         "spaces become hyphens in anchor",
			markdown:     "## Hello World",
			wantContains: `<h2 id="hello-world"><a href="#hello-world">Hello World</a></h2>`,
		},
		{
			name:         "special characters removed from ID but preserved in text",
			markdown:     "## Hello World!",
			wantContains: `<h2 id="hello-world"><a href="#hello-world">Hello World!</a></h2>`,
		},
		{
			name:         "mixed case becomes lowercase in ID",
			markdown:     "## Hello WORLD",
			wantContains: `<h2 id="hello-world"><a href="#hello-world">Hello WORLD</a></h2>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			html, err := parser.MarkdownToHTMLWithError(tt.markdown)
			if err != nil {
				t.Fatalf("MarkdownToHTMLWithError() error = %v", err)
			}
			if !strings.Contains(html, tt.wantContains) {
				t.Errorf("MarkdownToHTMLWithError() = %q, want to contain %q", html, tt.wantContains)
			}
		})
	}
}
