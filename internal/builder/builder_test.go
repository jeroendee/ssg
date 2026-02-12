package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jeroendee/ssg/internal/assets"
	"github.com/jeroendee/ssg/internal/model"
)

func TestScanContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(t *testing.T) string
		wantPages int
		wantPosts int
		wantErr   bool
	}{
		{
			name: "scans pages and posts",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, filepath.Join(dir, "home.md"), "---\ntitle: Home\n---\nWelcome")
				// Create pages
				writeFile(t, filepath.Join(dir, "about.md"), "---\ntitle: About\n---\nAbout content")
				writeFile(t, filepath.Join(dir, "contact.md"), "---\ntitle: Contact\n---\nContact content")
				// Create blog posts
				os.MkdirAll(filepath.Join(dir, "blog"), 0755)
				writeFile(t, filepath.Join(dir, "blog", "2021-03-26-first-post.md"), "---\ntitle: First Post\n---\nFirst content")
				writeFile(t, filepath.Join(dir, "blog", "2021-04-15-second-post.md"), "---\ntitle: Second Post\n---\nSecond content")
				return dir
			},
			wantPages: 3, // about + contact + home (home.md now parsed as page with empty slug)
			wantPosts: 2,
			wantErr:   false,
		},
		{
			name: "handles content directory with only home.md",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, filepath.Join(dir, "home.md"), "---\ntitle: Home\n---\nWelcome")
				return dir
			},
			wantPages: 1, // home.md parsed as page with empty slug
			wantPosts: 0,
			wantErr:   false,
		},
		{
			name: "handles pages only",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, filepath.Join(dir, "home.md"), "---\ntitle: Home\n---\nWelcome")
				writeFile(t, filepath.Join(dir, "about.md"), "---\ntitle: About\n---\nAbout content")
				return dir
			},
			wantPages: 2, // about + home
			wantPosts: 0,
			wantErr:   false,
		},
		{
			name: "handles posts only",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, filepath.Join(dir, "home.md"), "---\ntitle: Home\n---\nWelcome")
				os.MkdirAll(filepath.Join(dir, "blog"), 0755)
				writeFile(t, filepath.Join(dir, "blog", "2021-03-26-post.md"), "---\ntitle: Post\n---\nContent")
				return dir
			},
			wantPages: 1, // home.md parsed as page with empty slug
			wantPosts: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			contentDir := tt.setup(t)
			cfg := &model.Config{
				Title:      "Test Site",
				BaseURL:    "https://example.com",
				ContentDir: contentDir,
			}

			b := New(cfg)
			site, err := b.ScanContent()

			if (err != nil) != tt.wantErr {
				t.Errorf("ScanContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(site.Pages) != tt.wantPages {
					t.Errorf("ScanContent() pages = %d, want %d", len(site.Pages), tt.wantPages)
				}
				if len(site.Posts) != tt.wantPosts {
					t.Errorf("ScanContent() posts = %d, want %d", len(site.Posts), tt.wantPosts)
				}
			}
		})
	}
}

func TestScanContent_NonExistentDir(t *testing.T) {
	t.Parallel()

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: "/nonexistent/path",
	}

	b := New(cfg)
	_, err := b.ScanContent()

	if err == nil {
		t.Error("ScanContent() expected error for non-existent directory")
	}
}

func TestScanContent_PostsSortedByDate(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeHomeMd(t, dir)
	os.MkdirAll(filepath.Join(dir, "blog"), 0755)
	// Create posts in non-chronological order
	writeFile(t, filepath.Join(dir, "blog", "2021-01-01-old.md"), "---\ntitle: Old Post\n---\nOld")
	writeFile(t, filepath.Join(dir, "blog", "2023-06-15-newest.md"), "---\ntitle: Newest Post\n---\nNewest")
	writeFile(t, filepath.Join(dir, "blog", "2022-03-10-middle.md"), "---\ntitle: Middle Post\n---\nMiddle")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: dir,
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	if len(site.Posts) != 3 {
		t.Fatalf("ScanContent() posts = %d, want 3", len(site.Posts))
	}

	// Verify posts are sorted newest first
	if site.Posts[0].Title != "Newest Post" {
		t.Errorf("First post = %q, want 'Newest Post'", site.Posts[0].Title)
	}
	if site.Posts[1].Title != "Middle Post" {
		t.Errorf("Second post = %q, want 'Middle Post'", site.Posts[1].Title)
	}
	if site.Posts[2].Title != "Old Post" {
		t.Errorf("Third post = %q, want 'Old Post'", site.Posts[2].Title)
	}
}

func TestScanContent_SiteMetadata(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeHomeMd(t, dir)
	cfg := &model.Config{
		Title:   "Quality Shepherd",
		BaseURL: "https://www.qualityshepherd.nl",
		Author:  "Jeroen",
		Navigation: []model.NavItem{
			{Title: "Home", URL: "/"},
			{Title: "Blog", URL: "/blog/"},
		},
		ContentDir: dir,
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	if site.Title != cfg.Title {
		t.Errorf("Site.Title = %q, want %q", site.Title, cfg.Title)
	}
	if site.BaseURL != cfg.BaseURL {
		t.Errorf("Site.BaseURL = %q, want %q", site.BaseURL, cfg.BaseURL)
	}
	if site.Author != cfg.Author {
		t.Errorf("Site.Author = %q, want %q", site.Author, cfg.Author)
	}
	if len(site.Navigation) != len(cfg.Navigation) {
		t.Errorf("Site.Navigation len = %d, want %d", len(site.Navigation), len(cfg.Navigation))
	}
}

func TestScanContent_IgnoresNonMarkdownFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeHomeMd(t, dir)
	writeFile(t, filepath.Join(dir, "about.md"), "---\ntitle: About\n---\nAbout")
	writeFile(t, filepath.Join(dir, "readme.txt"), "not markdown")
	writeFile(t, filepath.Join(dir, "image.png"), "binary data")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: dir,
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	// Expect 2: about.md + home.md (home.md is now parsed as a page with empty slug)
	if len(site.Pages) != 2 {
		t.Errorf("ScanContent() pages = %d, want 2 (should ignore non-md files)", len(site.Pages))
	}
}

func TestScanContent_IgnoresIndexFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeHomeMd(t, dir)
	writeFile(t, filepath.Join(dir, "about.md"), "---\ntitle: About\n---\nAbout")
	writeFile(t, filepath.Join(dir, "_index.md"), "---\ntitle: Home\n---\nHome")
	os.MkdirAll(filepath.Join(dir, "blog"), 0755)
	writeFile(t, filepath.Join(dir, "blog", "_index.md"), "---\ntitle: Blog\n---\nBlog listing")
	writeFile(t, filepath.Join(dir, "blog", "2021-03-26-post.md"), "---\ntitle: Post\n---\nContent")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: dir,
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	// Should have 2 pages (about.md + home.md) and 1 post, ignoring _index.md files
	if len(site.Pages) != 2 {
		t.Errorf("ScanContent() pages = %d, want 2", len(site.Pages))
	}
	if len(site.Posts) != 1 {
		t.Errorf("ScanContent() posts = %d, want 1", len(site.Posts))
	}
}

// Helper to write test files
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
}

// writeHomeMd creates a minimal home.md in the given directory
func writeHomeMd(t *testing.T, dir string) {
	t.Helper()
	writeFile(t, filepath.Join(dir, "home.md"), "---\ntitle: Home\n---\nWelcome")
}

func TestBuild_CreatesOutputDirectory(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := filepath.Join(t.TempDir(), "public")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("Build() did not create output directory")
	}
}

func TestBuild_CleansOutputDirectory(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	// Create a file that should be cleaned
	oldFile := filepath.Join(outputDir, "old-file.html")
	writeFile(t, oldFile, "old content")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("Build() did not clean old files from output directory")
	}
}

func TestBuild_WritesPageWithCleanURL(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	writeFile(t, filepath.Join(contentDir, "about.md"), "---\ntitle: About\n---\nAbout content")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Check clean URL structure: /about/index.html
	pagePath := filepath.Join(outputDir, "about", "index.html")
	if _, err := os.Stat(pagePath); os.IsNotExist(err) {
		t.Errorf("Build() did not create page at %s", pagePath)
	}

	content, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatalf("failed to read page: %v", err)
	}
	if !strings.Contains(string(content), "About") {
		t.Error("Build() page does not contain expected content")
	}
}

func TestBuild_WritesPostWithCleanURL(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	os.MkdirAll(filepath.Join(contentDir, "blog"), 0755)
	writeFile(t, filepath.Join(contentDir, "blog", "2021-03-26-first-post.md"),
		"---\ntitle: First Post\n---\nFirst content")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Check clean URL structure: /blog/first-post/index.html
	postPath := filepath.Join(outputDir, "blog", "first-post", "index.html")
	if _, err := os.Stat(postPath); os.IsNotExist(err) {
		t.Errorf("Build() did not create post at %s", postPath)
	}

	content, err := os.ReadFile(postPath)
	if err != nil {
		t.Fatalf("failed to read post: %v", err)
	}
	if !strings.Contains(string(content), "First Post") {
		t.Error("Build() post does not contain expected content")
	}
}

func TestBuild_CopiesStaticAssets(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()
	assetsDir := t.TempDir()

	writeFile(t, filepath.Join(assetsDir, "style.css"), "body { color: black; }")
	writeFile(t, filepath.Join(assetsDir, "favicon.ico"), "fake favicon")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	b.SetAssetsDir(assetsDir)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Check assets are copied to output root
	cssPath := filepath.Join(outputDir, "style.css")
	if _, err := os.Stat(cssPath); os.IsNotExist(err) {
		t.Errorf("Build() did not copy style.css to %s", cssPath)
	}

	faviconPath := filepath.Join(outputDir, "favicon.ico")
	if _, err := os.Stat(faviconPath); os.IsNotExist(err) {
		t.Errorf("Build() did not copy favicon.ico to %s", faviconPath)
	}
}

func TestBuild_GeneratesHomepage(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	indexPath := filepath.Join(outputDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Errorf("Build() did not generate index.html at %s", indexPath)
	}

	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index.html: %v", err)
	}
	if !strings.Contains(string(content), "Test Site") {
		t.Error("Build() index.html does not contain site title")
	}
}

func TestBuild_Generates404Page(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	notFoundPath := filepath.Join(outputDir, "404.html")
	if _, err := os.Stat(notFoundPath); os.IsNotExist(err) {
		t.Errorf("Build() did not generate 404.html at %s", notFoundPath)
	}

	content, err := os.ReadFile(notFoundPath)
	if err != nil {
		t.Fatalf("failed to read 404.html: %v", err)
	}
	if !strings.Contains(string(content), "Not Found") && !strings.Contains(string(content), "404") {
		t.Error("Build() 404.html does not contain expected content")
	}
}

func TestBuild_GeneratesBlogListing(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	os.MkdirAll(filepath.Join(contentDir, "blog"), 0755)
	writeFile(t, filepath.Join(contentDir, "blog", "2021-03-26-first-post.md"),
		"---\ntitle: First Post\n---\nFirst content")
	writeFile(t, filepath.Join(contentDir, "blog", "2021-04-15-second-post.md"),
		"---\ntitle: Second Post\n---\nSecond content")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	blogIndexPath := filepath.Join(outputDir, "blog", "index.html")
	if _, err := os.Stat(blogIndexPath); os.IsNotExist(err) {
		t.Errorf("Build() did not generate blog/index.html at %s", blogIndexPath)
	}

	content, err := os.ReadFile(blogIndexPath)
	if err != nil {
		t.Fatalf("failed to read blog/index.html: %v", err)
	}
	if !strings.Contains(string(content), "First Post") {
		t.Error("Build() blog listing does not contain first post")
	}
	if !strings.Contains(string(content), "Second Post") {
		t.Error("Build() blog listing does not contain second post")
	}
}

func TestScanContent_DescriptionFromConfig(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeHomeMd(t, dir)
	cfg := &model.Config{
		Title:       "Test Site",
		BaseURL:     "https://example.com",
		Description: "A blog about testing",
		ContentDir:  dir,
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	if site.Description != "A blog about testing" {
		t.Errorf("Site.Description = %q, want %q", site.Description, "A blog about testing")
	}
}

func TestScanContent_DescriptionFallsBackToTitle(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeHomeMd(t, dir)
	cfg := &model.Config{
		Title:       "Test Site",
		BaseURL:     "https://example.com",
		Description: "", // empty description
		ContentDir:  dir,
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	if site.Description != "Test Site" {
		t.Errorf("Site.Description = %q, want fallback to Title %q", site.Description, "Test Site")
	}
}

func TestBuild_GeneratesRSSFeed(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	os.MkdirAll(filepath.Join(contentDir, "blog"), 0755)
	writeFile(t, filepath.Join(contentDir, "blog", "2021-03-26-first-post.md"),
		"---\ntitle: First Post\n---\nFirst content")

	cfg := &model.Config{
		Title:       "Test Site",
		BaseURL:     "https://example.com",
		Description: "A test blog",
		ContentDir:  contentDir,
		OutputDir:   outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	feedPath := filepath.Join(outputDir, "feed.xml")
	if _, err := os.Stat(feedPath); os.IsNotExist(err) {
		t.Errorf("Build() did not generate feed.xml at %s", feedPath)
	}

	content, err := os.ReadFile(feedPath)
	if err != nil {
		t.Fatalf("failed to read feed: %v", err)
	}
	if !strings.Contains(string(content), "<rss") {
		t.Error("Build() feed does not contain RSS element")
	}
	if !strings.Contains(string(content), "First Post") {
		t.Error("Build() feed does not contain post title")
	}
}

func TestBuild_SkipsFeedWhenNoPosts(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	feedPath := filepath.Join(outputDir, "feed.xml")
	if _, err := os.Stat(feedPath); !os.IsNotExist(err) {
		t.Errorf("Build() should not create feed.xml when no posts exist")
	}
}

func TestBuild_GeneratesRobotsTxt(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	robotsPath := filepath.Join(outputDir, "robots.txt")
	if _, err := os.Stat(robotsPath); os.IsNotExist(err) {
		t.Errorf("Build() did not generate robots.txt at %s", robotsPath)
	}

	content, err := os.ReadFile(robotsPath)
	if err != nil {
		t.Fatalf("failed to read robots.txt: %v", err)
	}

	expected := "User-agent: *\nAllow: /\nSitemap: https://example.com/sitemap.xml\n"
	if string(content) != expected {
		t.Errorf("robots.txt content = %q, want %q", string(content), expected)
	}
}

func TestBuild_GeneratesSitemap(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	// Create pages and posts
	writeFile(t, filepath.Join(contentDir, "about.md"), "---\ntitle: About\n---\nAbout content")
	writeFile(t, filepath.Join(contentDir, "contact.md"), "---\ntitle: Contact\n---\nContact content")

	os.MkdirAll(filepath.Join(contentDir, "blog"), 0755)
	writeFile(t, filepath.Join(contentDir, "blog", "2021-03-26-first-post.md"),
		"---\ntitle: First Post\n---\nFirst content")
	writeFile(t, filepath.Join(contentDir, "blog", "2021-04-15-second-post.md"),
		"---\ntitle: Second Post\n---\nSecond content")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	sitemapPath := filepath.Join(outputDir, "sitemap.xml")
	if _, err := os.Stat(sitemapPath); os.IsNotExist(err) {
		t.Fatalf("Build() did not generate sitemap.xml at %s", sitemapPath)
	}

	content, err := os.ReadFile(sitemapPath)
	if err != nil {
		t.Fatalf("failed to read sitemap.xml: %v", err)
	}

	sitemapStr := string(content)

	// Verify sitemap XML namespace
	if !strings.Contains(sitemapStr, `xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`) {
		t.Error("sitemap.xml missing correct namespace")
	}

	// Verify pages are included with full URLs
	if !strings.Contains(sitemapStr, "https://example.com/about/") {
		t.Error("sitemap.xml missing about page URL")
	}
	if !strings.Contains(sitemapStr, "https://example.com/contact/") {
		t.Error("sitemap.xml missing contact page URL")
	}

	// Verify posts are included with full URLs and dates
	if !strings.Contains(sitemapStr, "https://example.com/blog/first-post/") {
		t.Error("sitemap.xml missing first post URL")
	}
	if !strings.Contains(sitemapStr, "https://example.com/blog/second-post/") {
		t.Error("sitemap.xml missing second post URL")
	}
	if !strings.Contains(sitemapStr, "2021-03-26") {
		t.Error("sitemap.xml missing first post date")
	}
	if !strings.Contains(sitemapStr, "2021-04-15") {
		t.Error("sitemap.xml missing second post date")
	}

	// Verify lastmod elements exist for posts
	if !strings.Contains(sitemapStr, "<lastmod>") {
		t.Error("sitemap.xml missing lastmod elements")
	}
}

func TestBuild_SitemapHomepageURL(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	outputDir := t.TempDir()

	// Create home.md (homepage with empty slug)
	writeFile(t, filepath.Join(contentDir, "home.md"), "---\ntitle: Home\n---\nWelcome")
	// Create a regular page for comparison
	writeFile(t, filepath.Join(contentDir, "about.md"), "---\ntitle: About\n---\nAbout content")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	sitemapPath := filepath.Join(outputDir, "sitemap.xml")
	content, err := os.ReadFile(sitemapPath)
	if err != nil {
		t.Fatalf("failed to read sitemap.xml: %v", err)
	}

	sitemapStr := string(content)

	// Homepage must appear as baseURL/ without double slashes
	if !strings.Contains(sitemapStr, "<loc>https://example.com/</loc>") {
		t.Errorf("sitemap.xml homepage URL incorrect. Want <loc>https://example.com/</loc>\ngot:\n%s", sitemapStr)
	}

	// Must NOT contain double slashes (except in http://)
	// Check for the pattern that would indicate double slash in path: .com//
	if strings.Contains(sitemapStr, "example.com//") {
		t.Error("sitemap.xml contains double slashes in homepage URL")
	}

	// Regular page should still work correctly
	if !strings.Contains(sitemapStr, "https://example.com/about/") {
		t.Error("sitemap.xml missing about page URL")
	}
}

func TestBuild_SitemapWithNoPosts(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	// Create only pages
	writeFile(t, filepath.Join(contentDir, "about.md"), "---\ntitle: About\n---\nAbout content")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	sitemapPath := filepath.Join(outputDir, "sitemap.xml")
	if _, err := os.Stat(sitemapPath); os.IsNotExist(err) {
		t.Fatalf("Build() did not generate sitemap.xml at %s", sitemapPath)
	}

	content, err := os.ReadFile(sitemapPath)
	if err != nil {
		t.Fatalf("failed to read sitemap.xml: %v", err)
	}

	sitemapStr := string(content)

	// Verify page is included
	if !strings.Contains(sitemapStr, "https://example.com/about/") {
		t.Error("sitemap.xml missing about page URL")
	}
}

func TestSetVersion(t *testing.T) {
	t.Parallel()

	cfg := &model.Config{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}

	b := New(cfg)
	version := "1.2.3"
	b.SetVersion(version)

	if b.version != version {
		t.Errorf("SetVersion() version = %q, want %q", b.version, version)
	}
}

func TestBuilder_VersionStoredAndRetrievable(t *testing.T) {
	t.Parallel()

	cfg := &model.Config{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}

	b := New(cfg)
	testVersion := "v1.0.0-test"
	b.SetVersion(testVersion)

	if b.version != testVersion {
		t.Errorf("Builder.version = %q, want %q", b.version, testVersion)
	}
}

func TestBuild_PassesVersionToRenderer(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	writeFile(t, filepath.Join(contentDir, "about.md"), "---\ntitle: About\n---\nAbout content")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	testVersion := "test-version-123"
	b.SetVersion(testVersion)

	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Read generated HTML
	pagePath := filepath.Join(outputDir, "about", "index.html")
	content, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatalf("failed to read page: %v", err)
	}

	// Version should appear in footer
	if !strings.Contains(string(content), testVersion) {
		t.Errorf("Build() page should contain version %q in footer", testVersion)
	}
}

func TestBuild_EmbeddedStyleCSS_NoAssetsDir(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Embedded style.css should be written to output
	stylePath := filepath.Join(outputDir, "style.css")
	if _, err := os.Stat(stylePath); os.IsNotExist(err) {
		t.Fatalf("Build() did not write embedded style.css to %s", stylePath)
	}

	content, err := os.ReadFile(stylePath)
	if err != nil {
		t.Fatalf("failed to read style.css: %v", err)
	}

	if !bytes.Equal(content, assets.DefaultStyleCSS()) {
		t.Errorf("Build() style.css content does not match embedded default\ngot %d bytes, want %d bytes", len(content), len(assets.DefaultStyleCSS()))
	}
}

func TestBuild_CustomStyleCSS_Override(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()
	assetsDir := t.TempDir()

	customCSS := "body { color: red; }"
	writeFile(t, filepath.Join(assetsDir, "style.css"), customCSS)

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	b.SetAssetsDir(assetsDir)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Custom style.css should be used
	stylePath := filepath.Join(outputDir, "style.css")
	content, err := os.ReadFile(stylePath)
	if err != nil {
		t.Fatalf("failed to read style.css: %v", err)
	}

	if !bytes.Equal(content, []byte(customCSS)) {
		t.Errorf("Build() style.css content does not match custom CSS\ngot:  %q\nwant: %q", string(content), customCSS)
	}
}

func TestBuild_EmbeddedStyleCSS_WithOtherAssets(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()
	assetsDir := t.TempDir()

	writeFile(t, filepath.Join(assetsDir, "favicon.ico"), "fake favicon")
	writeFile(t, filepath.Join(assetsDir, "logo.png"), "fake logo")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	b.SetAssetsDir(assetsDir)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Embedded style.css should be written
	stylePath := filepath.Join(outputDir, "style.css")
	if _, err := os.Stat(stylePath); os.IsNotExist(err) {
		t.Fatalf("Build() did not write embedded style.css to %s", stylePath)
	}

	// Other assets should also be copied
	faviconPath := filepath.Join(outputDir, "favicon.ico")
	if _, err := os.Stat(faviconPath); os.IsNotExist(err) {
		t.Fatalf("Build() did not copy favicon.ico to %s", faviconPath)
	}

	logoPath := filepath.Join(outputDir, "logo.png")
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		t.Fatalf("Build() did not copy logo.png to %s", logoPath)
	}

	content, err := os.ReadFile(stylePath)
	if err != nil {
		t.Fatalf("failed to read style.css: %v", err)
	}

	if !bytes.Equal(content, assets.DefaultStyleCSS()) {
		t.Errorf("Build() style.css content does not match embedded default\ngot %d bytes, want %d bytes", len(content), len(assets.DefaultStyleCSS()))
	}
}

func TestValidateOutputDir(t *testing.T) {
	t.Parallel()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	contentDir := filepath.Join(t.TempDir(), "project", "content")
	projectRoot := filepath.Dir(contentDir)

	tests := []struct {
		name       string
		outputDir  string
		contentDir string
		wantErr    bool
		errContain string
	}{
		{
			name:       "reject empty output directory",
			outputDir:  "",
			contentDir: contentDir,
			wantErr:    true,
			errContain: "empty",
		},
		{
			name:       "reject root directory",
			outputDir:  "/",
			contentDir: contentDir,
			wantErr:    true,
			errContain: "root",
		},
		{
			name:       "reject home directory",
			outputDir:  homeDir,
			contentDir: contentDir,
			wantErr:    true,
			errContain: "home",
		},
		{
			name:       "reject tilde home path",
			outputDir:  "~",
			contentDir: contentDir,
			wantErr:    true,
			errContain: "home",
		},
		{
			name:       "reject current directory dot",
			outputDir:  ".",
			contentDir: contentDir,
			wantErr:    true,
			errContain: "current",
		},
		{
			name:       "reject parent directory dotdot",
			outputDir:  "..",
			contentDir: contentDir,
			wantErr:    true,
			errContain: "parent",
		},
		{
			name:       "reject path resolving to root",
			outputDir:  "/foo/..",
			contentDir: contentDir,
			wantErr:    true,
			errContain: "root",
		},
		{
			name:       "reject path resolving to home",
			outputDir:  filepath.Join(homeDir, "subdir", ".."),
			contentDir: contentDir,
			wantErr:    true,
			errContain: "home",
		},
		{
			name:       "reject path equal to project root",
			outputDir:  projectRoot,
			contentDir: contentDir,
			wantErr:    true,
			errContain: "project",
		},
		{
			name:       "reject path that is parent of content dir",
			outputDir:  filepath.Dir(projectRoot),
			contentDir: contentDir,
			wantErr:    true,
			errContain: "outside",
		},
		{
			name:       "accept valid output directory",
			outputDir:  filepath.Join(projectRoot, "public"),
			contentDir: contentDir,
			wantErr:    false,
		},
		{
			name:       "accept subdirectory of project",
			outputDir:  filepath.Join(projectRoot, "dist", "output"),
			contentDir: contentDir,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateOutputDir(tt.outputDir, tt.contentDir)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateOutputDir() expected error, got nil")
					return
				}
				if tt.errContain != "" && !strings.Contains(strings.ToLower(err.Error()), tt.errContain) {
					t.Errorf("ValidateOutputDir() error = %q, should contain %q", err.Error(), tt.errContain)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateOutputDir() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestBuild_HomeMdMissing(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	outputDir := t.TempDir()

	// No home.md exists in content directory
	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()

	if err == nil {
		t.Fatal("Build() should fail when home.md is missing")
	}

	wantMsg := "home.md not found in content directory"
	if !strings.Contains(err.Error(), wantMsg) {
		t.Errorf("Build() error = %q, want to contain %q", err.Error(), wantMsg)
	}
}

func TestBuild_HomepageOutput_WritesToRootIndex(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	outputDir := t.TempDir()

	// Create home.md with distinctive content
	writeFile(t, filepath.Join(contentDir, "home.md"),
		"---\ntitle: Welcome Home\n---\nThis is my unique homepage content from home.md")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Verify homepage is written to /index.html
	indexPath := filepath.Join(outputDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Fatal("Build() did not create /index.html")
	}

	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index.html: %v", err)
	}

	// Verify content comes from home.md, not hardcoded
	if !strings.Contains(string(content), "unique homepage content from home.md") {
		t.Error("Build() index.html does not contain content from home.md")
	}

	// Verify /home/index.html is NOT created
	homeDirPath := filepath.Join(outputDir, "home", "index.html")
	if _, err := os.Stat(homeDirPath); !os.IsNotExist(err) {
		t.Error("Build() should NOT create /home/index.html - homepage should go to /index.html")
	}
}

func TestCleanOutputDir_RejectsUnsafePaths(t *testing.T) {
	t.Parallel()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	contentDir := t.TempDir()

	unsafePaths := []string{
		"",
		"/",
		homeDir,
		".",
		"..",
	}

	for _, unsafePath := range unsafePaths {
		t.Run("reject_"+unsafePath, func(t *testing.T) {
			t.Parallel()

			cfg := &model.Config{
				Title:      "Test Site",
				BaseURL:    "https://example.com",
				ContentDir: contentDir,
				OutputDir:  unsafePath,
			}

			b := New(cfg)
			err := b.Build()

			if err == nil {
				t.Errorf("Build() with OutputDir=%q should return error", unsafePath)
			}
		})
	}
}

func TestCopyPostAssets_Success(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	outputDir := t.TempDir()

	// Create asset source directory and file
	assetsDir := filepath.Join(contentDir, "blog", "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		t.Fatalf("failed to create assets dir: %v", err)
	}
	assetContent := []byte("test image content")
	writeFile(t, filepath.Join(assetsDir, "test.png"), string(assetContent))

	// Create output directory for post
	postOutputDir := filepath.Join(outputDir, "blog", "my-post")
	if err := os.MkdirAll(postOutputDir, 0755); err != nil {
		t.Fatalf("failed to create post output dir: %v", err)
	}

	cfg := &model.Config{
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	post := model.Post{
		Page: model.Page{
			Slug: "my-post",
		},
		Assets: []string{"assets/test.png"},
	}

	err := b.copyPostAssets(post, postOutputDir)
	if err != nil {
		t.Fatalf("copyPostAssets() error = %v", err)
	}

	// Verify asset was copied
	copiedPath := filepath.Join(postOutputDir, "test.png")
	copiedContent, err := os.ReadFile(copiedPath)
	if err != nil {
		t.Fatalf("failed to read copied asset: %v", err)
	}

	if !bytes.Equal(copiedContent, assetContent) {
		t.Errorf("copied content = %q, want %q", copiedContent, assetContent)
	}
}

func TestCopyPostAssets_MissingAsset_ReturnsError(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	outputDir := t.TempDir()

	// Create blog/assets dir but NO actual asset file
	assetsDir := filepath.Join(contentDir, "blog", "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		t.Fatalf("failed to create assets dir: %v", err)
	}

	postOutputDir := filepath.Join(outputDir, "blog", "my-post")
	if err := os.MkdirAll(postOutputDir, 0755); err != nil {
		t.Fatalf("failed to create post output dir: %v", err)
	}

	cfg := &model.Config{
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	post := model.Post{
		Page: model.Page{
			Slug: "my-post",
		},
		Assets: []string{"assets/missing.png"},
	}

	err := b.copyPostAssets(post, postOutputDir)
	if err == nil {
		t.Fatal("copyPostAssets() expected error for missing asset")
	}

	// Verify error message is descriptive
	if !strings.Contains(err.Error(), "missing.png") {
		t.Errorf("error should mention asset name, got: %v", err)
	}
	if !strings.Contains(err.Error(), "my-post") {
		t.Errorf("error should mention post slug, got: %v", err)
	}
}

func TestCopyPostAssets_NoAssets(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	outputDir := t.TempDir()

	postOutputDir := filepath.Join(outputDir, "blog", "my-post")
	if err := os.MkdirAll(postOutputDir, 0755); err != nil {
		t.Fatalf("failed to create post output dir: %v", err)
	}

	cfg := &model.Config{
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	post := model.Post{
		Page: model.Page{
			Slug: "my-post",
		},
		Assets: []string{}, // Empty assets
	}

	err := b.copyPostAssets(post, postOutputDir)
	if err != nil {
		t.Errorf("copyPostAssets() with no assets should return nil, got: %v", err)
	}
}

func TestCopyPostAssets_MultipleAssets(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	outputDir := t.TempDir()

	// Create asset source directory and multiple files
	assetsDir := filepath.Join(contentDir, "blog", "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		t.Fatalf("failed to create assets dir: %v", err)
	}
	writeFile(t, filepath.Join(assetsDir, "image1.png"), "image 1 content")
	writeFile(t, filepath.Join(assetsDir, "image2.jpg"), "image 2 content")
	writeFile(t, filepath.Join(assetsDir, "diagram.svg"), "svg content")

	postOutputDir := filepath.Join(outputDir, "blog", "my-post")
	if err := os.MkdirAll(postOutputDir, 0755); err != nil {
		t.Fatalf("failed to create post output dir: %v", err)
	}

	cfg := &model.Config{
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := New(cfg)
	post := model.Post{
		Page: model.Page{
			Slug: "my-post",
		},
		Assets: []string{
			"assets/image1.png",
			"assets/image2.jpg",
			"assets/diagram.svg",
		},
	}

	err := b.copyPostAssets(post, postOutputDir)
	if err != nil {
		t.Fatalf("copyPostAssets() error = %v", err)
	}

	// Verify all assets were copied
	expectedFiles := map[string]string{
		"image1.png":  "image 1 content",
		"image2.jpg":  "image 2 content",
		"diagram.svg": "svg content",
	}

	for name, expectedContent := range expectedFiles {
		copiedPath := filepath.Join(postOutputDir, name)
		content, err := os.ReadFile(copiedPath)
		if err != nil {
			t.Errorf("failed to read copied asset %s: %v", name, err)
			continue
		}
		if string(content) != expectedContent {
			t.Errorf("asset %s content = %q, want %q", name, content, expectedContent)
		}
	}
}

func TestRewriteAssetPaths_DoubleQuotes(t *testing.T) {
	t.Parallel()

	input := `<img src="assets/pic.jpg" alt="test">`
	want := `<img src="pic.jpg" alt="test">`

	got := rewriteAssetPaths(input)

	if got != want {
		t.Errorf("rewriteAssetPaths() = %q, want %q", got, want)
	}
}

func TestRewriteAssetPaths_SingleQuotes(t *testing.T) {
	t.Parallel()

	input := `<img src='assets/pic.jpg' alt='test'>`
	want := `<img src='pic.jpg' alt='test'>`

	got := rewriteAssetPaths(input)

	if got != want {
		t.Errorf("rewriteAssetPaths() = %q, want %q", got, want)
	}
}

func TestRewriteAssetPaths_MultipleImages(t *testing.T) {
	t.Parallel()

	input := `<p><img src="assets/one.png"><img src="assets/two.jpg"></p>`
	want := `<p><img src="one.png"><img src="two.jpg"></p>`

	got := rewriteAssetPaths(input)

	if got != want {
		t.Errorf("rewriteAssetPaths() = %q, want %q", got, want)
	}
}

func TestRewriteAssetPaths_NoAssets(t *testing.T) {
	t.Parallel()

	input := `<img src="photo.jpg" alt="test">`
	want := `<img src="photo.jpg" alt="test">`

	got := rewriteAssetPaths(input)

	if got != want {
		t.Errorf("rewriteAssetPaths() = %q, want %q", got, want)
	}
}

func TestRewriteAssetPaths_ExternalUrls(t *testing.T) {
	t.Parallel()

	input := `<img src="https://example.com/assets/image.jpg">`
	want := `<img src="https://example.com/assets/image.jpg">`

	got := rewriteAssetPaths(input)

	if got != want {
		t.Errorf("rewriteAssetPaths() should not modify external URLs, got %q, want %q", got, want)
	}
}

func TestBlogAssetIntegration(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	outputDir := t.TempDir()

	// Setup: Create home.md (required by builder)
	writeFile(t, filepath.Join(contentDir, "home.md"), "---\ntitle: Home\n---\nWelcome")

	// Setup: Create blog/assets directory with test image
	assetsDir := filepath.Join(contentDir, "blog", "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		t.Fatalf("failed to create assets dir: %v", err)
	}
	writeFile(t, filepath.Join(assetsDir, "test-image.png"), "fake png content")

	// Setup: Create blog post that references the asset
	blogDir := filepath.Join(contentDir, "blog")
	postContent := `---
title: Test Post
---
Here is an image:

![](assets/test-image.png)
`
	writeFile(t, filepath.Join(blogDir, "2024-01-01-test.md"), postContent)

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	// Execute: Run full build
	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Assert 1: HTML output exists at expected path
	htmlPath := filepath.Join(outputDir, "blog", "test", "index.html")
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		t.Errorf("Build() did not create HTML at %s", htmlPath)
	}

	// Assert 2: Asset was copied to post directory
	assetPath := filepath.Join(outputDir, "blog", "test", "test-image.png")
	if _, err := os.Stat(assetPath); os.IsNotExist(err) {
		t.Errorf("Build() did not copy asset to %s", assetPath)
	}

	// Assert 3: HTML contains rewritten path (not assets/ prefix)
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		t.Fatalf("failed to read HTML: %v", err)
	}

	// Should contain: src="test-image.png"
	if !strings.Contains(string(htmlContent), `src="test-image.png"`) {
		t.Errorf("HTML should contain src=\"test-image.png\", got:\n%s", htmlContent)
	}

	// Should NOT contain: src="assets/test-image.png"
	if strings.Contains(string(htmlContent), `src="assets/test-image.png"`) {
		t.Errorf("HTML should NOT contain assets/ prefix, got:\n%s", htmlContent)
	}
}

func TestCollectFeedItems_PostsOnly(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)

	// Create blog posts
	blogDir := filepath.Join(contentDir, "blog")
	os.MkdirAll(blogDir, 0755)
	writeFile(t, filepath.Join(blogDir, "2026-01-27-post1.md"), "---\ntitle: Post 1\n---\nContent 1")
	writeFile(t, filepath.Join(blogDir, "2026-01-26-post2.md"), "---\ntitle: Post 2\n---\nContent 2")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		FeedPages:  []string{}, // No feed pages configured
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	items := b.collectFeedItems(*site)
	if len(items) != 2 {
		t.Errorf("collectFeedItems() returned %d items, want 2", len(items))
	}

	// Verify sorted by date (newest first)
	if items[0].FeedTitle() != "Post 1" {
		t.Errorf("First item title = %q, want %q", items[0].FeedTitle(), "Post 1")
	}
}

func TestCollectFeedItems_WithFeedPages(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)

	// Create a feed page with date sections using markdown headers that generate IDs
	// The parser uses the date format YYYY-MM-DD as heading IDs
	feedPageContent := `---
title: Moments
---
# My Moments

## 2026-01-28

Today's moment.

## 2026-01-25

Earlier moment.
`
	writeFile(t, filepath.Join(contentDir, "moments.md"), feedPageContent)

	// Create blog posts
	blogDir := filepath.Join(contentDir, "blog")
	os.MkdirAll(blogDir, 0755)
	writeFile(t, filepath.Join(blogDir, "2026-01-27-post.md"), "---\ntitle: Blog Post\n---\nPost content")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		FeedPages:  []string{"/moments/"},
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	items := b.collectFeedItems(*site)
	if len(items) != 3 {
		t.Errorf("collectFeedItems() returned %d items, want 3 (1 post + 2 date sections)", len(items))
		return
	}

	// Verify sorted by date (newest first): 2026-01-28, 2026-01-27, 2026-01-25
	if !strings.Contains(items[0].FeedTitle(), "January 28, 2026") {
		t.Errorf("First item title = %q, expected date section from Jan 28", items[0].FeedTitle())
	}
	if items[1].FeedTitle() != "Blog Post" {
		t.Errorf("Second item title = %q, want %q", items[1].FeedTitle(), "Blog Post")
	}
}

func TestCollectFeedItems_NonExistentFeedPage(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		FeedPages:  []string{"/nonexistent/"},
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	// Should not panic, just skip the missing page
	items := b.collectFeedItems(*site)
	if len(items) != 0 {
		t.Errorf("collectFeedItems() returned %d items, want 0 (no posts, missing feed page)", len(items))
	}
}

func TestCollectFeedItems_LimitsTo20(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)

	// Create 25 blog posts
	blogDir := filepath.Join(contentDir, "blog")
	os.MkdirAll(blogDir, 0755)
	for i := 1; i <= 25; i++ {
		filename := fmt.Sprintf("2026-01-%02d-post%d.md", i, i)
		content := fmt.Sprintf("---\ntitle: Post %d\n---\nContent", i)
		writeFile(t, filepath.Join(blogDir, filename), content)
	}

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		FeedPages:  []string{},
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	items := b.collectFeedItems(*site)
	if len(items) != 20 {
		t.Errorf("collectFeedItems() returned %d items, want 20 (limit)", len(items))
	}
}

func TestFindPageByPath(t *testing.T) {
	t.Parallel()

	pages := []Page{
		{Title: "About", Slug: "about", Path: "/about/"},
		{Title: "Now", Slug: "now", Path: "/now/"},
		{Title: "Home", Slug: "", Path: "/"},
	}

	tests := []struct {
		path     string
		wantName string
		wantNil  bool
	}{
		{"/now/", "Now", false},
		{"now/", "Now", false}, // Missing leading slash
		{"/now", "Now", false}, // Missing trailing slash
		{"/about/", "About", false},
		{"/missing/", "", true},
		{"/", "Home", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			got := findPageByPath(pages, tt.path)
			if tt.wantNil {
				if got != nil {
					t.Errorf("findPageByPath(%q) = %v, want nil", tt.path, got)
				}
				return
			}
			if got == nil {
				t.Errorf("findPageByPath(%q) = nil, want page with title %q", tt.path, tt.wantName)
				return
			}
			if got.Title != tt.wantName {
				t.Errorf("findPageByPath(%q) = %q, want %q", tt.path, got.Title, tt.wantName)
			}
		})
	}
}

func TestBuild_FeedIncludesPages(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	outputDir := t.TempDir()

	writeHomeMd(t, contentDir)

	// Create a feed page with date sections
	feedPageContent := `---
title: Moments
---
# My Moments

## 2026-01-28

Today's moment.
`
	writeFile(t, filepath.Join(contentDir, "moments.md"), feedPageContent)

	// Create blog post
	blogDir := filepath.Join(contentDir, "blog")
	os.MkdirAll(blogDir, 0755)
	writeFile(t, filepath.Join(blogDir, "2026-01-27-post.md"), "---\ntitle: Blog Post\n---\nPost content")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
		FeedPages:  []string{"/moments/"},
	}

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Read feed.xml
	feedPath := filepath.Join(outputDir, "feed.xml")
	content, err := os.ReadFile(feedPath)
	if err != nil {
		t.Fatalf("failed to read feed.xml: %v", err)
	}

	feedStr := string(content)

	// Check blog post is in feed
	if !strings.Contains(feedStr, "<title>Blog Post</title>") {
		t.Error("feed.xml missing blog post")
	}

	// Check page date section is in feed
	if !strings.Contains(feedStr, "Moments - January 28, 2026") {
		t.Error("feed.xml missing page date section")
	}

	// Check link format for page section
	if !strings.Contains(feedStr, "https://example.com/moments/#2026-01-28") {
		t.Error("feed.xml missing page date section link with anchor")
	}
}

func TestScanContent_FooterFileExists(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeHomeMd(t, dir)
	writeFile(t, filepath.Join(dir, "_footer.md"), "Contact: [email](mailto:test@example.com)")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: dir,
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	if site.FooterContent == "" {
		t.Error("ScanContent() FooterContent should not be empty when _footer.md exists")
	}

	// Verify markdown was converted to HTML
	if !strings.Contains(site.FooterContent, "<a href=") {
		t.Errorf("ScanContent() FooterContent should contain HTML link, got: %s", site.FooterContent)
	}
}

func TestScanContent_FooterFileNotExists(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeHomeMd(t, dir)
	// No _footer.md created

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: dir,
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	if site.FooterContent != "" {
		t.Errorf("ScanContent() FooterContent should be empty when _footer.md doesn't exist, got: %s", site.FooterContent)
	}
}

func TestScanContent_FooterExcludedFromPages(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeHomeMd(t, dir)
	writeFile(t, filepath.Join(dir, "_footer.md"), "Footer content")
	writeFile(t, filepath.Join(dir, "about.md"), "---\ntitle: About\n---\nAbout")

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: dir,
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	// Should have 2 pages (home + about), not 3
	if len(site.Pages) != 2 {
		t.Errorf("ScanContent() pages = %d, want 2 (footer should be excluded)", len(site.Pages))
	}

	// Verify footer is not in pages
	for _, page := range site.Pages {
		if page.Slug == "_footer" || strings.Contains(page.Title, "footer") {
			t.Errorf("ScanContent() _footer.md should not be in pages, found: %s", page.Slug)
		}
	}
}

func TestBuild_GeneratesBuildTimestamp(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
	writeHomeMd(t, contentDir)
	outputDir := t.TempDir()

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	beforeBuild := time.Now().UTC()

	b := New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	afterBuild := time.Now().UTC()

	// Verify build.json exists
	buildPath := filepath.Join(outputDir, "build.json")
	if _, err := os.Stat(buildPath); os.IsNotExist(err) {
		t.Fatalf("Build() did not generate build.json at %s", buildPath)
	}

	// Verify JSON structure and format
	content, err := os.ReadFile(buildPath)
	if err != nil {
		t.Fatalf("failed to read build.json: %v", err)
	}

	// Parse as JSON
	var ts struct {
		BuildTime string `json:"buildTime"`
	}
	if err := json.Unmarshal(content, &ts); err != nil {
		t.Fatalf("build.json is not valid JSON: %v", err)
	}

	// Verify buildTime field exists
	if ts.BuildTime == "" {
		t.Error("build.json buildTime field is empty")
	}

	// Verify RFC 3339 format by parsing
	buildTime, err := time.Parse(time.RFC3339, ts.BuildTime)
	if err != nil {
		t.Errorf("build.json buildTime is not RFC 3339 format: %v", err)
	}

	// Verify timestamp is within expected range
	if buildTime.Before(beforeBuild.Add(-time.Second)) || buildTime.After(afterBuild.Add(time.Second)) {
		t.Errorf("build.json buildTime %v is outside expected range [%v, %v]", buildTime, beforeBuild, afterBuild)
	}
}

func TestScanContent_TopicsPopulated(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeHomeMd(t, dir)

	// Create a page with repeated words that should produce topics
	momentsContent := "---\ntitle: Moments\n---\n" +
		"claude claude claude agent agent agent docker docker docker\n"
	writeFile(t, filepath.Join(dir, "moments.md"), momentsContent)

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: dir,
		TopicPages: []string{"/moments/"},
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	// Find the moments page
	var momentsPage *model.Page
	for i := range site.Pages {
		if site.Pages[i].Slug == "moments" {
			momentsPage = &site.Pages[i]
			break
		}
	}
	if momentsPage == nil {
		t.Fatal("moments page not found")
	}

	if len(momentsPage.Topics) == 0 {
		t.Error("moments page should have topics populated")
	}

	// Verify expected topics
	topicMap := make(map[string]int)
	for _, topic := range momentsPage.Topics {
		topicMap[topic.Word] = topic.Count
	}
	if topicMap["claude"] != 3 {
		t.Errorf("expected 'claude' with count 3, got %d", topicMap["claude"])
	}
	if topicMap["agent"] != 3 {
		t.Errorf("expected 'agent' with count 3, got %d", topicMap["agent"])
	}
}

func TestScanContent_TopicsNotPopulatedWithoutConfig(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeHomeMd(t, dir)

	momentsContent := "---\ntitle: Moments\n---\n" +
		"claude claude claude agent agent agent docker docker docker\n"
	writeFile(t, filepath.Join(dir, "moments.md"), momentsContent)

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: dir,
		TopicPages: []string{}, // No topic pages configured
	}

	b := New(cfg)
	site, err := b.ScanContent()
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	// Find the moments page
	for _, page := range site.Pages {
		if page.Slug == "moments" {
			if len(page.Topics) != 0 {
				t.Errorf("moments page should NOT have topics when not configured, got %v", page.Topics)
			}
			return
		}
	}
	t.Fatal("moments page not found")
}
