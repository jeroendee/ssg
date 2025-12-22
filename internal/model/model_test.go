package model_test

import (
	"testing"
	"time"

	"github.com/jeroendee/ssg/internal/model"
)

func TestNavItem(t *testing.T) {
	t.Parallel()
	nav := model.NavItem{
		Title: "Home",
		URL:   "/",
	}
	if nav.Title != "Home" {
		t.Errorf("NavItem.Title = %q, want %q", nav.Title, "Home")
	}
	if nav.URL != "/" {
		t.Errorf("NavItem.URL = %q, want %q", nav.URL, "/")
	}
}

func TestPage(t *testing.T) {
	t.Parallel()
	page := model.Page{
		Title:    "About",
		Slug:     "about",
		Content:  "<p>About content</p>",
		Path:     "/about/",
		Template: "page",
	}
	if page.Title != "About" {
		t.Errorf("Page.Title = %q, want %q", page.Title, "About")
	}
	if page.Slug != "about" {
		t.Errorf("Page.Slug = %q, want %q", page.Slug, "about")
	}
	if page.Content != "<p>About content</p>" {
		t.Errorf("Page.Content = %q, want %q", page.Content, "<p>About content</p>")
	}
	if page.Path != "/about/" {
		t.Errorf("Page.Path = %q, want %q", page.Path, "/about/")
	}
	if page.Template != "page" {
		t.Errorf("Page.Template = %q, want %q", page.Template, "page")
	}
}

func TestPost(t *testing.T) {
	t.Parallel()
	date := time.Date(2021, 3, 26, 0, 0, 0, 0, time.UTC)
	post := model.Post{
		Page: model.Page{
			Title:    "My First Post",
			Slug:     "my-first-post",
			Content:  "<p>Post content</p>",
			Path:     "/blog/my-first-post/",
			Template: "post",
		},
		Date:    date,
		Summary: "A summary of my first post",
	}
	if post.Title != "My First Post" {
		t.Errorf("Post.Title = %q, want %q", post.Title, "My First Post")
	}
	if !post.Date.Equal(date) {
		t.Errorf("Post.Date = %v, want %v", post.Date, date)
	}
	if post.Summary != "A summary of my first post" {
		t.Errorf("Post.Summary = %q, want %q", post.Summary, "A summary of my first post")
	}
}

func TestSite(t *testing.T) {
	t.Parallel()
	site := model.Site{
		Title:   "Quality Shepherd",
		BaseURL: "https://www.qualityshepherd.nl",
		Author:  "Jeroen",
		Navigation: []model.NavItem{
			{Title: "Home", URL: "/"},
			{Title: "About", URL: "/about/"},
		},
		Pages: []model.Page{
			{Title: "About", Slug: "about", Path: "/about/"},
		},
		Posts: []model.Post{
			{Page: model.Page{Title: "Post 1", Slug: "post-1"}},
		},
	}
	if site.Title != "Quality Shepherd" {
		t.Errorf("Site.Title = %q, want %q", site.Title, "Quality Shepherd")
	}
	if site.BaseURL != "https://www.qualityshepherd.nl" {
		t.Errorf("Site.BaseURL = %q, want %q", site.BaseURL, "https://www.qualityshepherd.nl")
	}
	if site.Author != "Jeroen" {
		t.Errorf("Site.Author = %q, want %q", site.Author, "Jeroen")
	}
	if len(site.Navigation) != 2 {
		t.Errorf("len(Site.Navigation) = %d, want %d", len(site.Navigation), 2)
	}
	if len(site.Pages) != 1 {
		t.Errorf("len(Site.Pages) = %d, want %d", len(site.Pages), 1)
	}
	if len(site.Posts) != 1 {
		t.Errorf("len(Site.Posts) = %d, want %d", len(site.Posts), 1)
	}
}

func TestConfig(t *testing.T) {
	t.Parallel()
	cfg := model.Config{
		Title:      "Quality Shepherd",
		BaseURL:    "https://www.qualityshepherd.nl",
		Author:     "Jeroen",
		ContentDir: "content",
		OutputDir:  "public",
		AssetsDir:  "assets",
		Navigation: []model.NavItem{
			{Title: "Home", URL: "/"},
		},
	}
	if cfg.Title != "Quality Shepherd" {
		t.Errorf("Config.Title = %q, want %q", cfg.Title, "Quality Shepherd")
	}
	if cfg.ContentDir != "content" {
		t.Errorf("Config.ContentDir = %q, want %q", cfg.ContentDir, "content")
	}
	if cfg.OutputDir != "public" {
		t.Errorf("Config.OutputDir = %q, want %q", cfg.OutputDir, "public")
	}
	if cfg.AssetsDir != "assets" {
		t.Errorf("Config.AssetsDir = %q, want %q", cfg.AssetsDir, "assets")
	}
}

func TestSite_FaviconMIMEType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		favicon string
		want    string
	}{
		{
			name:    "svg extension returns svg+xml",
			favicon: "/favicon.svg",
			want:    "image/svg+xml",
		},
		{
			name:    "ico extension returns x-icon",
			favicon: "/favicon.ico",
			want:    "image/x-icon",
		},
		{
			name:    "png extension returns png",
			favicon: "/favicon.png",
			want:    "image/png",
		},
		{
			name:    "gif extension returns gif",
			favicon: "/favicon.gif",
			want:    "image/gif",
		},
		{
			name:    "unknown extension returns default x-icon",
			favicon: "/favicon.webp",
			want:    "image/x-icon",
		},
		{
			name:    "empty favicon returns default x-icon",
			favicon: "",
			want:    "image/x-icon",
		},
		{
			name:    "uppercase SVG extension returns svg+xml",
			favicon: "/favicon.SVG",
			want:    "image/svg+xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			site := model.Site{Favicon: tt.favicon}
			got := site.FaviconMIMEType()

			if got != tt.want {
				t.Errorf("Site.FaviconMIMEType() = %q, want %q", got, tt.want)
			}
		})
	}
}
