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

// Analytics holds analytics configuration.
type Analytics struct {
	GoatCounter string
}

// Page represents a static page.
type Page struct {
	Title       string
	Slug        string
	Content     string
	Path        string
	DateAnchors []string // Date anchors for navigation (e.g., "2026-01-26")
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
	Title       string
	Description string
	BaseURL     string
	Author      string
	Logo        string
	OGImage     string
	Favicon     string
	Navigation  []NavItem
	Pages       []Page
	Posts       []Post
	Analytics   Analytics
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
