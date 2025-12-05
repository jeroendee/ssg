package builder

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jeroendee/ssg/internal/model"
	"github.com/jeroendee/ssg/internal/parser"
	"github.com/jeroendee/ssg/internal/renderer"
)

// Builder orchestrates site building from content files.
type Builder struct {
	cfg       *model.Config
	assetsDir string
}

// New creates a new Builder with the given configuration.
func New(cfg *model.Config) *Builder {
	return &Builder{cfg: cfg}
}

// SetAssetsDir sets the assets directory for static files.
func (b *Builder) SetAssetsDir(dir string) {
	b.assetsDir = dir
}

// ScanContent scans the content directory and returns a Site with all pages and posts.
func (b *Builder) ScanContent() (*model.Site, error) {
	site := &model.Site{
		Title:      b.cfg.Title,
		BaseURL:    b.cfg.BaseURL,
		Author:     b.cfg.Author,
		Logo:       b.cfg.Logo,
		Favicon:    b.cfg.Favicon,
		Navigation: b.cfg.Navigation,
	}

	// Check if content directory exists
	if _, err := os.Stat(b.cfg.ContentDir); os.IsNotExist(err) {
		return nil, err
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

	return nil
}

// cleanOutputDir removes and recreates the output directory.
func (b *Builder) cleanOutputDir() error {
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
	if b.assetsDir == "" {
		return nil
	}

	entries, err := os.ReadDir(b.assetsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

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
