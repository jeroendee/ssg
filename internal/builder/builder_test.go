package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
				// Create pages
				writeFile(t, filepath.Join(dir, "about.md"), "---\ntitle: About\n---\nAbout content")
				writeFile(t, filepath.Join(dir, "contact.md"), "---\ntitle: Contact\n---\nContact content")
				// Create blog posts
				os.MkdirAll(filepath.Join(dir, "blog"), 0755)
				writeFile(t, filepath.Join(dir, "blog", "2021-03-26-first-post.md"), "---\ntitle: First Post\n---\nFirst content")
				writeFile(t, filepath.Join(dir, "blog", "2021-04-15-second-post.md"), "---\ntitle: Second Post\n---\nSecond content")
				return dir
			},
			wantPages: 2,
			wantPosts: 2,
			wantErr:   false,
		},
		{
			name: "handles empty content directory",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantPages: 0,
			wantPosts: 0,
			wantErr:   false,
		},
		{
			name: "handles pages only",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, filepath.Join(dir, "about.md"), "---\ntitle: About\n---\nAbout content")
				return dir
			},
			wantPages: 1,
			wantPosts: 0,
			wantErr:   false,
		},
		{
			name: "handles posts only",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				os.MkdirAll(filepath.Join(dir, "blog"), 0755)
				writeFile(t, filepath.Join(dir, "blog", "2021-03-26-post.md"), "---\ntitle: Post\n---\nContent")
				return dir
			},
			wantPages: 0,
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

	if len(site.Pages) != 1 {
		t.Errorf("ScanContent() pages = %d, want 1 (should ignore non-md files)", len(site.Pages))
	}
}

func TestScanContent_IgnoresIndexFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
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

	// Should have 1 page (about.md) and 1 post, ignoring _index.md files
	if len(site.Pages) != 1 {
		t.Errorf("ScanContent() pages = %d, want 1", len(site.Pages))
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

// Helper to compare times
func timeEqual(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func TestBuild_CreatesOutputDirectory(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
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

	feedPath := filepath.Join(outputDir, "feed", "index.xml")
	if _, err := os.Stat(feedPath); os.IsNotExist(err) {
		t.Errorf("Build() did not generate feed/index.xml at %s", feedPath)
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

	feedPath := filepath.Join(outputDir, "feed", "index.xml")
	if _, err := os.Stat(feedPath); !os.IsNotExist(err) {
		t.Errorf("Build() should not create feed when no posts exist")
	}
}

func TestBuild_GeneratesRobotsTxt(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
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

func TestBuild_SitemapWithNoPosts(t *testing.T) {
	t.Parallel()

	contentDir := t.TempDir()
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
