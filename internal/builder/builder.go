package builder

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jeroendee/ssg/internal/assets"
	"github.com/jeroendee/ssg/internal/model"
	"github.com/jeroendee/ssg/internal/parser"
	"github.com/jeroendee/ssg/internal/renderer"
)

// Builder orchestrates site building from content files.
type Builder struct {
	cfg       *model.Config
	assetsDir string
	version   string
}

// New creates a new Builder with the given configuration.
func New(cfg *model.Config) *Builder {
	return &Builder{cfg: cfg}
}

// SetAssetsDir sets the assets directory for static files.
func (b *Builder) SetAssetsDir(dir string) {
	b.assetsDir = dir
}

// SetVersion sets the version string for the build.
func (b *Builder) SetVersion(version string) {
	b.version = version
}

// ScanContent scans the content directory and returns a Site with all pages and posts.
func (b *Builder) ScanContent() (*model.Site, error) {
	description := b.cfg.Description
	if description == "" {
		description = b.cfg.Title
	}

	site := &model.Site{
		Title:       b.cfg.Title,
		Description: description,
		BaseURL:     b.cfg.BaseURL,
		Author:      b.cfg.Author,
		Logo:        b.cfg.Logo,
		Favicon:     b.cfg.Favicon,
		Navigation:  b.cfg.Navigation,
		Analytics:   b.cfg.Analytics,
	}

	// Check if content directory exists
	if _, err := os.Stat(b.cfg.ContentDir); os.IsNotExist(err) {
		return nil, err
	}

	// Verify home.md exists - homepage is required
	homePath := filepath.Join(b.cfg.ContentDir, "home.md")
	if _, err := os.Stat(homePath); os.IsNotExist(err) {
		return nil, errors.New("home.md not found in content directory - homepage is required")
	}

	// Scan for pages (markdown files in root content directory)
	rootEntries, err := os.ReadDir(b.cfg.ContentDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range rootEntries {
		if entry.IsDir() {
			continue
		}
		if !isMarkdownFile(entry.Name()) {
			continue
		}
		if isIndexFile(entry.Name()) {
			continue
		}
		// Skip home.md - it's handled separately as the homepage
		if entry.Name() == "home.md" {
			continue
		}

		path := filepath.Join(b.cfg.ContentDir, entry.Name())
		page, err := parser.ParsePage(path)
		if err != nil {
			return nil, err
		}
		site.Pages = append(site.Pages, *page)
	}

	// Scan for posts (markdown files in blog/ subdirectory)
	blogDir := filepath.Join(b.cfg.ContentDir, "blog")
	if _, err := os.Stat(blogDir); err == nil {
		blogEntries, err := os.ReadDir(blogDir)
		if err != nil {
			return nil, err
		}

		for _, entry := range blogEntries {
			if entry.IsDir() {
				continue
			}
			if !isMarkdownFile(entry.Name()) {
				continue
			}
			if isIndexFile(entry.Name()) {
				continue
			}

			path := filepath.Join(blogDir, entry.Name())
			post, err := parser.ParsePost(path)
			if err != nil {
				return nil, err
			}
			site.Posts = append(site.Posts, *post)
		}
	}

	// Sort posts by date (newest first)
	sort.Slice(site.Posts, func(i, j int) bool {
		return site.Posts[i].Date.After(site.Posts[j].Date)
	})

	return site, nil
}

// isMarkdownFile returns true if the filename has a .md extension.
func isMarkdownFile(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".md")
}

// isIndexFile returns true if the filename is an index file (_index.md).
func isIndexFile(name string) bool {
	return strings.HasPrefix(name, "_")
}

// Build performs a complete site build: scans content, renders HTML, copies assets.
func (b *Builder) Build() error {
	// Clean and create output directory
	if err := b.cleanOutputDir(); err != nil {
		return err
	}

	// Scan content
	site, err := b.ScanContent()
	if err != nil {
		return err
	}

	// Create renderer
	r, err := renderer.New()
	if err != nil {
		return err
	}
	r.SetVersion(b.version)

	// Render pages with clean URLs
	for _, page := range site.Pages {
		if err := b.writePage(r, *site, page); err != nil {
			return err
		}
	}

	// Render blog posts with clean URLs
	for _, post := range site.Posts {
		if err := b.writePost(r, *site, post); err != nil {
			return err
		}
	}

	// Generate blog listing if there are posts
	if len(site.Posts) > 0 {
		if err := b.writeBlogListing(r, *site); err != nil {
			return err
		}
	}

	// Generate RSS feed
	if err := b.writeFeed(r, *site); err != nil {
		return err
	}

	// Generate homepage
	if err := b.writeHomepage(r, *site); err != nil {
		return err
	}

	// Generate 404 page
	if err := b.write404(r, *site); err != nil {
		return err
	}

	// Copy static assets
	if err := b.copyAssets(); err != nil {
		return err
	}

	// Generate sitemap.xml
	if err := b.generateSitemap(*site); err != nil {
		return err
	}

	// Generate robots.txt
	if err := b.generateRobotsTxt(); err != nil {
		return err
	}

	return nil
}

// ValidateOutputDir checks if the output directory is safe for deletion.
// It rejects paths that could cause catastrophic data loss.
func ValidateOutputDir(outputDir, contentDir string) error {
	if outputDir == "" {
		return errors.New("output directory cannot be empty")
	}

	// Expand tilde to home directory
	if outputDir == "~" || strings.HasPrefix(outputDir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		if outputDir == "~" {
			outputDir = home
		} else {
			outputDir = filepath.Join(home, outputDir[2:])
		}
	}

	// Clean the path to resolve . and .. components
	cleanPath := filepath.Clean(outputDir)

	// Reject current directory
	if cleanPath == "." {
		return errors.New("output directory cannot be current directory")
	}

	// Reject parent directory traversal
	if cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) {
		return errors.New("output directory cannot be parent directory")
	}

	// Get absolute path for comparison
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return err
	}

	// Reject root directory
	if absPath == "/" {
		return errors.New("output directory cannot be root directory")
	}

	// Reject home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	if absPath == home {
		return errors.New("output directory cannot be home directory")
	}

	// Get project root (parent of content directory)
	absContentDir, err := filepath.Abs(contentDir)
	if err != nil {
		return err
	}
	projectRoot := filepath.Dir(absContentDir)

	// Reject if output equals project root
	if absPath == projectRoot {
		return errors.New("output directory cannot be project root")
	}

	// Reject if output is outside or parent of project root
	if !strings.HasPrefix(absPath, projectRoot+string(filepath.Separator)) {
		return errors.New("output directory is outside project root")
	}

	return nil
}

// cleanOutputDir removes and recreates the output directory.
func (b *Builder) cleanOutputDir() error {
	// Validate output directory before deletion
	if err := ValidateOutputDir(b.cfg.OutputDir, b.cfg.ContentDir); err != nil {
		return err
	}

	// Remove existing output directory
	if _, err := os.Stat(b.cfg.OutputDir); err == nil {
		if err := os.RemoveAll(b.cfg.OutputDir); err != nil {
			return err
		}
	}

	// Create fresh output directory
	return os.MkdirAll(b.cfg.OutputDir, 0755)
}

// writePage writes a page to its clean URL path.
func (b *Builder) writePage(r *renderer.Renderer, site model.Site, page model.Page) error {
	html, err := r.RenderPage(site, page)
	if err != nil {
		return err
	}

	// Create clean URL directory: /slug/index.html
	dir := filepath.Join(b.cfg.OutputDir, page.Slug)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "index.html"), []byte(html), 0644)
}

// writePost writes a post to its clean URL path.
func (b *Builder) writePost(r *renderer.Renderer, site model.Site, post model.Post) error {
	html, err := r.RenderBlogPost(site, post)
	if err != nil {
		return err
	}

	// Create clean URL directory: /blog/slug/index.html
	dir := filepath.Join(b.cfg.OutputDir, "blog", post.Slug)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "index.html"), []byte(html), 0644)
}

// writeBlogListing writes the blog listing page.
func (b *Builder) writeBlogListing(r *renderer.Renderer, site model.Site) error {
	html, err := r.RenderBlogList(site, site.Posts)
	if err != nil {
		return err
	}

	dir := filepath.Join(b.cfg.OutputDir, "blog")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "index.html"), []byte(html), 0644)
}

// writeFeed writes the RSS feed.
func (b *Builder) writeFeed(r *renderer.Renderer, site model.Site) error {
	xml, err := r.RenderFeed(site, site.Posts)
	if err != nil {
		return err
	}

	// Skip if no posts (RenderFeed returns empty string)
	if xml == "" {
		return nil
	}

	dir := filepath.Join(b.cfg.OutputDir, "feed")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "index.xml"), []byte(xml), 0644)
}

// writeHomepage writes the homepage.
func (b *Builder) writeHomepage(r *renderer.Renderer, site model.Site) error {
	html, err := r.RenderHome(site)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(b.cfg.OutputDir, "index.html"), []byte(html), 0644)
}

// write404 writes the 404 error page.
func (b *Builder) write404(r *renderer.Renderer, site model.Site) error {
	html, err := r.Render404(site)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(b.cfg.OutputDir, "404.html"), []byte(html), 0644)
}

// copyAssets copies static assets to the output directory.
func (b *Builder) copyAssets() error {
	if b.assetsDir != "" {
		entries, err := os.ReadDir(b.assetsDir)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				src := filepath.Join(b.assetsDir, entry.Name())
				dst := filepath.Join(b.cfg.OutputDir, entry.Name())

				if err := copyFile(src, dst); err != nil {
					return err
				}
			}
		}
	}

	// Write embedded default style.css if no custom one exists
	stylePath := filepath.Join(b.cfg.OutputDir, "style.css")
	if _, err := os.Stat(stylePath); os.IsNotExist(err) {
		if err := os.WriteFile(stylePath, assets.DefaultStyleCSS(), 0644); err != nil {
			return err
		}
	}

	return nil
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// generateRobotsTxt creates robots.txt with sitemap reference.
func (b *Builder) generateRobotsTxt() error {
	content := "User-agent: *\nAllow: /\nSitemap: " + b.cfg.BaseURL + "/sitemap.xml\n"
	robotsPath := filepath.Join(b.cfg.OutputDir, "robots.txt")
	return os.WriteFile(robotsPath, []byte(content), 0644)
}

// sitemapURL represents a URL entry in the sitemap.
type sitemapURL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

// sitemapURLSet represents the root urlset element.
type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

// generateSitemap creates sitemap.xml with all pages and posts.
func (b *Builder) generateSitemap(site model.Site) error {
	var urls []sitemapURL

	// Add pages (without lastmod)
	for _, page := range site.Pages {
		urls = append(urls, sitemapURL{
			Loc: site.BaseURL + "/" + page.Slug + "/",
		})
	}

	// Add posts (with lastmod using post date)
	for _, post := range site.Posts {
		urls = append(urls, sitemapURL{
			Loc:     site.BaseURL + "/blog/" + post.Slug + "/",
			LastMod: post.Date.Format("2006-01-02"),
		})
	}

	urlset := sitemapURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	var buf strings.Builder
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(urlset); err != nil {
		return err
	}

	sitemapPath := filepath.Join(b.cfg.OutputDir, "sitemap.xml")
	return os.WriteFile(sitemapPath, []byte(buf.String()), 0644)
}
