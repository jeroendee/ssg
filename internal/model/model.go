package model

import "time"

// NavItem represents a navigation menu entry.
type NavItem struct {
	Title string
	URL   string
}

// Page represents a static page.
type Page struct {
	Title    string
	Slug     string
	Content  string
	Path     string
	Template string
}

// Post represents a blog post with date and summary.
type Post struct {
	Page
	Date    time.Time
	Summary string
}

// Site represents the complete site with all pages and posts.
type Site struct {
	Title      string
	BaseURL    string
	Author     string
	Navigation []NavItem
	Pages      []Page
	Posts      []Post
}

// Config holds site configuration loaded from ssg.yaml.
type Config struct {
	Title      string
	BaseURL    string
	Author     string
	ContentDir string
	OutputDir  string
	Navigation []NavItem
}
