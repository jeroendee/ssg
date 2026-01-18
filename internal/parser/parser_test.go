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
	if !strings.Contains(html, "<h1>Hello</h1>") {
		t.Errorf("expected <h1>Hello</h1>, got %s", html)
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
	if !strings.Contains(html, "<h1>Hello</h1>") {
		t.Errorf("expected <h1>Hello</h1>, got %s", html)
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
