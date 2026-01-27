package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jeroendee/ssg/internal/model"
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

func TestParsePage_ExtractsDateAnchors(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "now.md")
	content := `---
title: "Now"
---
# What I'm doing now

#### *2026-01-26*

Working on the SSG project.

#### *2026-01-20*

Started a new feature.`
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	page, err := parser.ParsePage(file)
	if err != nil {
		t.Fatalf("ParsePage() error = %v", err)
	}
	wantAnchors := []string{"2026-01-26", "2026-01-20"}
	if len(page.DateAnchors) != len(wantAnchors) {
		t.Errorf("DateAnchors = %v, want %v", page.DateAnchors, wantAnchors)
		return
	}
	for i := range page.DateAnchors {
		if page.DateAnchors[i] != wantAnchors[i] {
			t.Errorf("DateAnchors[%d] = %q, want %q", i, page.DateAnchors[i], wantAnchors[i])
		}
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

func TestExtractDateAnchors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		markdown string
		want     []string
	}{
		{
			name:     "single date anchor",
			markdown: "#### *2026-01-26*\n\nSome content here.",
			want:     []string{"2026-01-26"},
		},
		{
			name:     "multiple date anchors",
			markdown: "#### *2026-01-26*\n\nFirst entry\n\n#### *2026-01-20*\n\nSecond entry",
			want:     []string{"2026-01-26", "2026-01-20"},
		},
		{
			name:     "no date anchors",
			markdown: "## Regular heading\n\nJust plain content.",
			want:     []string{},
		},
		{
			name:     "accept dates at all heading levels",
			markdown: "# *2026-01-26*\n## *2026-01-25*\n### *2026-01-24*\n#### *2026-01-23*\n##### *2026-01-22*\n###### *2026-01-21*",
			want:     []string{"2026-01-26", "2026-01-25", "2026-01-24", "2026-01-23", "2026-01-22", "2026-01-21"},
		},
		{
			name:     "ignore dates without italic",
			markdown: "#### 2026-01-26\n\nNot in italics.",
			want:     []string{},
		},
		{
			name:     "empty markdown",
			markdown: "",
			want:     []string{},
		},
		{
			name:     "date with surrounding content on heading",
			markdown: "#### *2026-01-26* - Weekly update\n\nContent here.",
			want:     []string{"2026-01-26"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.ExtractDateAnchors(tt.markdown)
			if len(got) != len(tt.want) {
				t.Errorf("ExtractDateAnchors() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ExtractDateAnchors()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGroupDatesByMonth(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                  string
		dates                 []string
		wantCurrentMonth      []string
		wantArchivedMonths    int
		wantFirstArchiveYear  int
		wantFirstArchiveMonth string
	}{
		{
			name:                  "multiple months of dates",
			dates:                 []string{"2026-02-03", "2026-02-01", "2026-01-31", "2026-01-15", "2025-12-20"},
			wantCurrentMonth:      []string{"2026-02-03", "2026-02-01"},
			wantArchivedMonths:    2,
			wantFirstArchiveYear:  2026,
			wantFirstArchiveMonth: "January",
		},
		{
			name:               "single month of dates",
			dates:              []string{"2026-01-27", "2026-01-26", "2026-01-25"},
			wantCurrentMonth:   []string{"2026-01-27", "2026-01-26", "2026-01-25"},
			wantArchivedMonths: 0,
		},
		{
			name:               "empty input",
			dates:              []string{},
			wantCurrentMonth:   nil,
			wantArchivedMonths: 0,
		},
		{
			name:               "malformed dates excluded",
			dates:              []string{"2026-01-27", "not-a-date", "2026-01-25"},
			wantCurrentMonth:   []string{"2026-01-27", "2026-01-25"},
			wantArchivedMonths: 0,
		},
		{
			name:                  "archives ordered newest first",
			dates:                 []string{"2026-03-01", "2026-01-15", "2025-11-20"},
			wantCurrentMonth:      []string{"2026-03-01"},
			wantArchivedMonths:    2,
			wantFirstArchiveYear:  2026,
			wantFirstArchiveMonth: "January",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			current, archived := parser.GroupDatesByMonth(tt.dates)

			// Check current month dates
			if len(current) != len(tt.wantCurrentMonth) {
				t.Errorf("CurrentMonthDates = %v, want %v", current, tt.wantCurrentMonth)
			} else {
				for i := range current {
					if current[i] != tt.wantCurrentMonth[i] {
						t.Errorf("CurrentMonthDates[%d] = %q, want %q", i, current[i], tt.wantCurrentMonth[i])
					}
				}
			}

			// Check archived months count
			if len(archived) != tt.wantArchivedMonths {
				t.Errorf("ArchivedMonths count = %d, want %d", len(archived), tt.wantArchivedMonths)
			}

			// Check first archive details if applicable
			if tt.wantArchivedMonths > 0 && len(archived) > 0 {
				if archived[0].Year != tt.wantFirstArchiveYear {
					t.Errorf("First archive Year = %d, want %d", archived[0].Year, tt.wantFirstArchiveYear)
				}
				if archived[0].Month != tt.wantFirstArchiveMonth {
					t.Errorf("First archive Month = %q, want %q", archived[0].Month, tt.wantFirstArchiveMonth)
				}
			}
		})
	}
}

func TestGroupMonthsByYear(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		months                  []model.MonthGroup
		wantYearCount           int
		wantFirstYear           int
		wantFirstYearMonthCount int
	}{
		{
			name: "multiple years",
			months: []model.MonthGroup{
				{Year: 2026, Month: "January", Dates: []string{"2026-01-15"}},
				{Year: 2025, Month: "December", Dates: []string{"2025-12-20"}},
				{Year: 2025, Month: "November", Dates: []string{"2025-11-15"}},
			},
			wantYearCount:           2,
			wantFirstYear:           2026,
			wantFirstYearMonthCount: 1,
		},
		{
			name: "single year",
			months: []model.MonthGroup{
				{Year: 2025, Month: "December", Dates: []string{"2025-12-20"}},
				{Year: 2025, Month: "November", Dates: []string{"2025-11-15"}},
				{Year: 2025, Month: "October", Dates: []string{"2025-10-10"}},
			},
			wantYearCount:           1,
			wantFirstYear:           2025,
			wantFirstYearMonthCount: 3,
		},
		{
			name:          "empty input",
			months:        []model.MonthGroup{},
			wantYearCount: 0,
		},
		{
			name: "years ordered newest first",
			months: []model.MonthGroup{
				{Year: 2024, Month: "March", Dates: []string{"2024-03-01"}},
				{Year: 2026, Month: "January", Dates: []string{"2026-01-15"}},
				{Year: 2025, Month: "June", Dates: []string{"2025-06-20"}},
			},
			wantYearCount:           3,
			wantFirstYear:           2026,
			wantFirstYearMonthCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.GroupMonthsByYear(tt.months)

			if len(got) != tt.wantYearCount {
				t.Errorf("GroupMonthsByYear() year count = %d, want %d", len(got), tt.wantYearCount)
			}

			if tt.wantYearCount > 0 && len(got) > 0 {
				if got[0].Year != tt.wantFirstYear {
					t.Errorf("GroupMonthsByYear() first year = %d, want %d", got[0].Year, tt.wantFirstYear)
				}
				if len(got[0].Months) != tt.wantFirstYearMonthCount {
					t.Errorf("GroupMonthsByYear() first year month count = %d, want %d", len(got[0].Months), tt.wantFirstYearMonthCount)
				}
			}
		})
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
