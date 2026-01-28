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
	"github.com/yuin/goldmark/extension"
	gmparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"gopkg.in/yaml.v3"
)

// frontmatter holds metadata extracted from markdown files.
type frontmatter struct {
	Title   string `yaml:"title"`
	Summary string `yaml:"summary"`
	Date    string `yaml:"date"`
}

// md is the configured Goldmark instance with auto heading IDs and anchor links.
var md = goldmark.New(
	goldmark.WithExtensions(extension.Table),
	goldmark.WithParserOptions(
		gmparser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(newAnchorHeadingRenderer(), 100),
		),
	),
)

// MarkdownToHTMLWithError converts markdown content to HTML and returns any conversion error.
func MarkdownToHTMLWithError(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", fmt.Errorf("markdown conversion failed: %w", err)
	}
	return buf.String(), nil
}

// extractFrontmatter separates YAML frontmatter from markdown content.
func extractFrontmatter(content string) (frontmatter, string, error) {
	var fm frontmatter
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

	fm, body, err := extractFrontmatter(string(data))
	if err != nil {
		return nil, err
	}

	slug := strings.TrimSuffix(filepath.Base(path), ".md")
	if slug == "home" {
		slug = ""
	}

	html, err := MarkdownToHTMLWithError(body)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	pagePath := "/" + slug + "/"
	if slug == "" {
		pagePath = "/"
	}

	dateAnchors := ExtractDateAnchors(body)
	currentMonth, archivedMonths := GroupDatesByMonth(dateAnchors)
	archivedYears := GroupMonthsByYear(archivedMonths)

	return &model.Page{
		Title:             fm.Title,
		Slug:              slug,
		Content:           html,
		Path:              pagePath,
		DateAnchors:       dateAnchors,
		CurrentMonthDates: currentMonth,
		ArchivedYears:     archivedYears,
	}, nil
}

var dateFilenameRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})-(.+)\.md$`)
var dateAnchorRegex = regexp.MustCompile(`(?m)^#{1,6} \*(\d{4}-\d{2}-\d{2})\*`)

// ExtractDateAnchors finds date headings (h1-h6) in italic format from markdown.
func ExtractDateAnchors(markdown string) []string {
	matches := dateAnchorRegex.FindAllStringSubmatch(markdown, -1)
	if matches == nil {
		return []string{}
	}
	dates := make([]string, len(matches))
	for i, m := range matches {
		dates[i] = m[1]
	}
	return dates
}

// GroupDatesByMonth separates date anchors into current month and archived months.
// Current month is the month containing the most recent valid date.
// Archived months are returned newest-first. Malformed dates are excluded.
func GroupDatesByMonth(dates []string) (currentMonth []string, archived []model.MonthGroup) {
	if len(dates) == 0 {
		return nil, nil
	}

	// Parse and group dates by year-month
	type yearMonth struct {
		year  int
		month time.Month
	}
	groups := make(map[yearMonth][]string)
	var validDates []yearMonth

	for _, d := range dates {
		t, err := time.Parse("2006-01-02", d)
		if err != nil {
			continue
		}
		ym := yearMonth{year: t.Year(), month: t.Month()}
		if _, exists := groups[ym]; !exists {
			validDates = append(validDates, ym)
		}
		groups[ym] = append(groups[ym], d)
	}

	if len(validDates) == 0 {
		return nil, nil
	}

	// Find the most recent month (first valid date determines current month)
	currentYM := validDates[0]
	for _, ym := range validDates {
		if ym.year > currentYM.year || (ym.year == currentYM.year && ym.month > currentYM.month) {
			currentYM = ym
		}
	}

	currentMonth = groups[currentYM]

	// Collect archived months (all except current), sorted newest-first
	var archiveYMs []yearMonth
	for ym := range groups {
		if ym != currentYM {
			archiveYMs = append(archiveYMs, ym)
		}
	}

	// Sort newest-first
	for i := 0; i < len(archiveYMs)-1; i++ {
		for j := i + 1; j < len(archiveYMs); j++ {
			if archiveYMs[j].year > archiveYMs[i].year ||
				(archiveYMs[j].year == archiveYMs[i].year && archiveYMs[j].month > archiveYMs[i].month) {
				archiveYMs[i], archiveYMs[j] = archiveYMs[j], archiveYMs[i]
			}
		}
	}

	for _, ym := range archiveYMs {
		archived = append(archived, model.MonthGroup{
			Year:  ym.year,
			Month: ym.month.String(),
			Dates: groups[ym],
		})
	}

	return currentMonth, archived
}

// GroupMonthsByYear groups archived months under year headers.
// Years are ordered newest to oldest. Months within each year preserve their input order.
func GroupMonthsByYear(months []model.MonthGroup) []model.YearGroup {
	if len(months) == 0 {
		return nil
	}

	// Group by year, preserving month order
	yearMap := make(map[int][]model.MonthGroup)
	var years []int

	for _, m := range months {
		if _, exists := yearMap[m.Year]; !exists {
			years = append(years, m.Year)
		}
		yearMap[m.Year] = append(yearMap[m.Year], m)
	}

	// Sort years newest first
	for i := 0; i < len(years)-1; i++ {
		for j := i + 1; j < len(years); j++ {
			if years[j] > years[i] {
				years[i], years[j] = years[j], years[i]
			}
		}
	}

	// Build result
	result := make([]model.YearGroup, len(years))
	for i, y := range years {
		result[i] = model.YearGroup{
			Year:   y,
			Months: yearMap[y],
		}
	}

	return result
}

var assetRefRegex = regexp.MustCompile(`!\[.*?\]\((assets/[^)]+)\)`)

// dateHeaderRegex matches HTML headings with date anchor IDs (YYYY-MM-DD format).
// Handles headings like: <h2 id="2026-01-27"><a href="#2026-01-27">January 27, 2026</a></h2>
var dateHeaderRegex = regexp.MustCompile(`<h[1-6] id="(\d{4}-\d{2}-\d{2})">.*?</h[1-6]>`)

// FeedDateSection represents a date-anchored section extracted from HTML for feed generation.
type FeedDateSection struct {
	Anchor  string
	Content string
}

// ExtractFeedDateSections finds date-anchored sections in HTML content.
// Returns sections with their date anchors and the content between that header and the next.
func ExtractFeedDateSections(html string) []FeedDateSection {
	// Find all date headers
	matches := dateHeaderRegex.FindAllStringSubmatchIndex(html, -1)
	if matches == nil {
		return []FeedDateSection{}
	}

	sections := make([]FeedDateSection, 0, len(matches))

	for i, match := range matches {
		// match[0]:match[1] is the full header match
		// match[2]:match[3] is the date anchor capture group
		anchor := html[match[2]:match[3]]

		// Content starts after this header
		contentStart := match[1]

		// Content ends at the next heading (any level) or end of document
		var contentEnd int
		if i+1 < len(matches) {
			// Find the start of the next date header
			contentEnd = matches[i+1][0]
		} else {
			contentEnd = len(html)
		}

		// Extract and trim content
		content := strings.TrimSpace(html[contentStart:contentEnd])

		// Skip if there's no content
		if content == "" {
			continue
		}

		sections = append(sections, FeedDateSection{
			Anchor:  anchor,
			Content: content,
		})
	}

	return sections
}

// ParseDateFromAnchor parses a date string in YYYY-MM-DD format to time.Time.
func ParseDateFromAnchor(anchor string) (time.Time, error) {
	if anchor == "" {
		return time.Time{}, fmt.Errorf("empty anchor")
	}
	return time.Parse("2006-01-02", anchor)
}

// ExtractAssetReferences finds all asset references in markdown content.
func ExtractAssetReferences(markdown string) []string {
	matches := assetRefRegex.FindAllStringSubmatch(markdown, -1)
	if matches == nil {
		return []string{}
	}
	assets := make([]string, len(matches))
	for i, m := range matches {
		assets[i] = m[1]
	}
	return assets
}

// ParsePost reads a markdown file and returns a Post.
func ParsePost(path string) (*model.Post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fm, body, err := extractFrontmatter(string(data))
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

	html, err := MarkdownToHTMLWithError(body)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	return &model.Post{
		Page: model.Page{
			Title:   fm.Title,
			Slug:    slug,
			Content: html,
			Path:    "/blog/" + slug + "/",
		},
		Date:      postDate,
		Summary:   fm.Summary,
		WordCount: wordcount.Count(body),
		Assets:    ExtractAssetReferences(body),
	}, nil
}
