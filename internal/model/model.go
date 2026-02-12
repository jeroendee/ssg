package model

import (
	"path"
	"strings"
	"time"
)

// NavItem represents a navigation menu entry.
type NavItem struct {
	Title string
	URL   string
}

// FeedItem is the interface for items that can appear in the RSS feed.
type FeedItem interface {
	FeedTitle() string
	FeedLink() string
	FeedContent() string
	FeedDate() time.Time
	FeedGUID() string
}

// MonthGroup groups date anchors by year and month for archive navigation.
type MonthGroup struct {
	Year  int
	Month string
	Dates []string
}

// YearGroup groups months by year for hierarchical archive navigation.
type YearGroup struct {
	Year   int
	Months []MonthGroup
}

// Topic represents a word with its frequency count.
type Topic struct {
	Word  string
	Count int
}

// Analytics holds analytics configuration.
type Analytics struct {
	GoatCounter string
}

// Page represents a static page.
type Page struct {
	Title             string
	Slug              string
	Content           string
	Path              string
	DateAnchors       []string    // Date anchors for navigation (e.g., "2026-01-26")
	CurrentMonthDates []string    // Dates from the most recent month
	ArchivedYears     []YearGroup // Previous months grouped by year for archive navigation
	Topics            []Topic     // Extracted topic words with frequency counts
}

// Post represents a blog post with date and summary.
type Post struct {
	Page
	Date      time.Time
	Summary   string
	WordCount int
	Assets    []string // Referenced asset paths from markdown
}

// Site represents the complete site with all pages and posts.
type Site struct {
	Title         string
	Description   string
	BaseURL       string
	Author        string
	Logo          string
	OGImage       string
	Favicon       string
	Navigation    []NavItem
	Pages         []Page
	Posts         []Post
	Analytics     Analytics
	FooterContent string
}

// Config holds site configuration loaded from ssg.yaml.
type Config struct {
	Title       string
	Description string
	BaseURL     string
	Author      string
	Logo        string
	OGImage     string
	Favicon     string
	ContentDir  string
	OutputDir   string
	AssetsDir   string
	Navigation  []NavItem
	Analytics   Analytics
	FeedPages   []string
	TopicPages  []string
}

// FaviconMIMEType returns the MIME type for the favicon based on file extension.
func (s Site) FaviconMIMEType() string {
	ext := strings.ToLower(path.Ext(s.Favicon))
	switch ext {
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "image/x-icon"
	}
}

// PostFeedAdapter wraps a Post to implement FeedItem with site context.
type PostFeedAdapter struct {
	Post    *Post
	BaseURL string
}

// FeedTitle returns the post title.
func (p PostFeedAdapter) FeedTitle() string {
	return p.Post.Title
}

// FeedLink returns the absolute URL to the post.
func (p PostFeedAdapter) FeedLink() string {
	return p.BaseURL + "/blog/" + p.Post.Slug + "/"
}

// FeedContent returns the post's HTML content.
func (p PostFeedAdapter) FeedContent() string {
	return p.Post.Content
}

// FeedDate returns the post's publication date.
func (p PostFeedAdapter) FeedDate() time.Time {
	return p.Post.Date
}

// FeedGUID returns the unique identifier (same as link for posts).
func (p PostFeedAdapter) FeedGUID() string {
	return p.FeedLink()
}

// DateSection represents a date-anchored section from a page for the RSS feed.
type DateSection struct {
	PageTitle string
	PagePath  string
	Date      time.Time
	Anchor    string
	Content   string
	BaseURL   string
}

// FeedTitle returns the title in format "{PageTitle} - {Month} {Day}, {Year}".
func (d DateSection) FeedTitle() string {
	return d.PageTitle + " - " + d.Date.Format("January 2, 2006")
}

// FeedLink returns the absolute URL to the page with the date anchor.
func (d DateSection) FeedLink() string {
	return d.BaseURL + d.PagePath + "#" + d.Anchor
}

// FeedContent returns the section's HTML content.
func (d DateSection) FeedContent() string {
	return d.Content
}

// FeedDate returns the date of this section.
func (d DateSection) FeedDate() time.Time {
	return d.Date
}

// FeedGUID returns the unique identifier (same as link).
func (d DateSection) FeedGUID() string {
	return d.FeedLink()
}
