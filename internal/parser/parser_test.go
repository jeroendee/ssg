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
	html := parser.MarkdownToHTML(md)
	if !strings.Contains(html, "<h1>Hello</h1>") {
		t.Errorf("expected <h1>Hello</h1>, got %s", html)
	}
	if !strings.Contains(html, "<strong>bold</strong>") {
		t.Errorf("expected <strong>bold</strong>, got %s", html)
	}
}

func TestExtractFrontmatter(t *testing.T) {
	t.Parallel()
	content := `---
title: "My Page"
template: "custom"
summary: "A brief summary"
---
# Content here`

	fm, body, err := parser.ExtractFrontmatter(content)
	if err != nil {
		t.Fatalf("ExtractFrontmatter() error = %v", err)
	}
	if fm.Title != "My Page" {
		t.Errorf("Title = %q, want %q", fm.Title, "My Page")
	}
	if fm.Template != "custom" {
		t.Errorf("Template = %q, want %q", fm.Template, "custom")
	}
	if fm.Summary != "A brief summary" {
		t.Errorf("Summary = %q, want %q", fm.Summary, "A brief summary")
	}
	if !strings.Contains(body, "# Content here") {
		t.Errorf("body should contain content, got %q", body)
	}
}

func TestExtractFrontmatter_NoFrontmatter(t *testing.T) {
	t.Parallel()
	content := "# Just content\n\nNo frontmatter here."
	fm, body, err := parser.ExtractFrontmatter(content)
	if err != nil {
		t.Fatalf("ExtractFrontmatter() error = %v", err)
	}
	if fm.Title != "" {
		t.Errorf("Title = %q, want empty", fm.Title)
	}
	if !strings.Contains(body, "# Just content") {
		t.Errorf("body = %q, want original content", body)
	}
}

func TestParsePage(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "about.md")
	content := `---
title: "About Me"
template: "page"
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
	if page.Template != "page" {
		t.Errorf("Template = %q, want %q", page.Template, "page")
	}
	if !strings.Contains(page.Content, "<h1>About</h1>") {
		t.Errorf("Content should have HTML, got %q", page.Content)
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
