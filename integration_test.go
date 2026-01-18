package ssg_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeroendee/ssg/internal/builder"
	"github.com/jeroendee/ssg/internal/model"
)

func TestIntegration_FullBuild(t *testing.T) {
	t.Parallel()

	// Setup test directories
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")
	assetsDir := filepath.Join(tmpDir, "assets")

	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(filepath.Join(contentDir, "blog"), 0755)
	os.MkdirAll(assetsDir, 0755)

	// Create required home.md
	writeTestFile(t, filepath.Join(contentDir, "home.md"), `---
title: Home
---
Welcome to Quality Shepherd`)

	// Create sample pages matching qualityshepherd.nl structure
	writeTestFile(t, filepath.Join(contentDir, "about.md"), `---
title: About
---
Hello! My name is Jeroen.

I'm a software engineer passionate about quality.`)

	writeTestFile(t, filepath.Join(contentDir, "contact.md"), `---
title: Contact
---
You can reach me at example@example.com`)

	// Create sample blog posts
	writeTestFile(t, filepath.Join(contentDir, "blog", "2021-03-26-hello-world.md"), `---
title: Hello World
---
This is my first blog post!`)

	writeTestFile(t, filepath.Join(contentDir, "blog", "2023-06-15-testing-tips.md"), `---
title: Testing Tips
---
Here are some tips for testing Go applications.`)

	// Create assets
	writeTestFile(t, filepath.Join(assetsDir, "style.css"), `body { color: black; }
@media (prefers-color-scheme: dark) { body { color: white; } }`)

	// Create config
	cfg := &model.Config{
		Title:   "Quality Shepherd",
		BaseURL: "https://www.qualityshepherd.nl",
		Author:  "Jeroen",
		Navigation: []model.NavItem{
			{Title: "Home", URL: "/"},
			{Title: "About", URL: "/about/"},
			{Title: "Blog", URL: "/blog/"},
		},
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	// Build site
	b := builder.New(cfg)
	b.SetAssetsDir(assetsDir)

	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Verify output structure
	verifyFileExists(t, outputDir, "index.html")
	verifyFileExists(t, outputDir, "404.html")
	verifyFileExists(t, outputDir, "style.css")
	verifyFileExists(t, outputDir, "about/index.html")
	verifyFileExists(t, outputDir, "contact/index.html")
	verifyFileExists(t, outputDir, "blog/index.html")
	verifyFileExists(t, outputDir, "blog/hello-world/index.html")
	verifyFileExists(t, outputDir, "blog/testing-tips/index.html")
}

func TestIntegration_HTMLOutputCorrectness(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")

	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(filepath.Join(contentDir, "blog"), 0755)

	// Create required home.md
	writeTestFile(t, filepath.Join(contentDir, "home.md"), `---
title: Home
---
Welcome`)

	writeTestFile(t, filepath.Join(contentDir, "about.md"), `---
title: About Me
---
This is the about page content.`)

	writeTestFile(t, filepath.Join(contentDir, "blog", "2023-01-15-test-post.md"), `---
title: Test Post
---
This is the post content.`)

	cfg := &model.Config{
		Title:   "Test Site",
		BaseURL: "https://example.com",
		Navigation: []model.NavItem{
			{Title: "Home", URL: "/"},
			{Title: "Blog", URL: "/blog/"},
		},
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := builder.New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Verify homepage HTML
	homepage := readTestFile(t, filepath.Join(outputDir, "index.html"))
	assertContains(t, homepage, "<!DOCTYPE html>", "homepage should have DOCTYPE")
	assertContains(t, homepage, "Test Site", "homepage should have site title")
	assertContains(t, homepage, "<nav", "homepage should have navigation")

	// Verify about page HTML
	aboutPage := readTestFile(t, filepath.Join(outputDir, "about", "index.html"))
	assertContains(t, aboutPage, "About Me", "about page should have title")
	assertContains(t, aboutPage, "This is the about page content", "about page should have content")
	assertContains(t, aboutPage, "<h1>", "about page should have h1")

	// Verify blog listing HTML
	blogList := readTestFile(t, filepath.Join(outputDir, "blog", "index.html"))
	assertContains(t, blogList, "Test Post", "blog list should show post title")
	assertContains(t, blogList, "2023-01-15", "blog list should show formatted date")
	assertContains(t, blogList, `href="/blog/test-post/"`, "blog list should link to post")

	// Verify blog post HTML
	blogPost := readTestFile(t, filepath.Join(outputDir, "blog", "test-post", "index.html"))
	assertContains(t, blogPost, "Test Post", "blog post should have title")
	assertContains(t, blogPost, "2023-01-15", "blog post should show date")
	assertContains(t, blogPost, "This is the post content", "blog post should have content")

	// Verify 404 page HTML
	notFound := readTestFile(t, filepath.Join(outputDir, "404.html"))
	assertContains(t, notFound, "404", "404 page should indicate error")
	assertContains(t, notFound, "Not Found", "404 page should say Not Found")
}

func TestIntegration_CSSInclusion(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")
	assetsDir := filepath.Join(tmpDir, "assets")

	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(assetsDir, 0755)

	// Create required home.md
	writeTestFile(t, filepath.Join(contentDir, "home.md"), `---
title: Home
---
Welcome`)

	// CSS with dark mode
	cssContent := `body { color: black; }
@media (prefers-color-scheme: dark) {
  body { color: white; background-color: #333; }
}`
	writeTestFile(t, filepath.Join(assetsDir, "style.css"), cssContent)

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := builder.New(cfg)
	b.SetAssetsDir(assetsDir)

	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Verify CSS file is copied
	verifyFileExists(t, outputDir, "style.css")

	// Verify CSS content is preserved
	outputCSS := readTestFile(t, filepath.Join(outputDir, "style.css"))
	assertContains(t, outputCSS, "prefers-color-scheme: dark", "CSS should have dark mode media query")

	// Verify HTML links to CSS
	homepage := readTestFile(t, filepath.Join(outputDir, "index.html"))
	assertContains(t, homepage, `href="/style.css"`, "HTML should link to CSS")
}

func TestIntegration_NavigationLinks(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")

	os.MkdirAll(contentDir, 0755)

	// Create required home.md
	writeTestFile(t, filepath.Join(contentDir, "home.md"), `---
title: Home
---
Welcome`)

	cfg := &model.Config{
		Title:   "Test Site",
		BaseURL: "https://example.com",
		Navigation: []model.NavItem{
			{Title: "Home", URL: "/"},
			{Title: "About", URL: "/about/"},
			{Title: "Blog", URL: "/blog/"},
			{Title: "Contact", URL: "/contact/"},
		},
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := builder.New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	homepage := readTestFile(t, filepath.Join(outputDir, "index.html"))

	// Verify all navigation links are present
	assertContains(t, homepage, `href="/"`, "should have home link")
	assertContains(t, homepage, `href="/about/"`, "should have about link")
	assertContains(t, homepage, `href="/blog/"`, "should have blog link")
	assertContains(t, homepage, `href="/contact/"`, "should have contact link")
	assertContains(t, homepage, ">Home<", "should show Home text")
	assertContains(t, homepage, ">About<", "should show About text")
	assertContains(t, homepage, ">Blog<", "should show Blog text")
	assertContains(t, homepage, ">Contact<", "should show Contact text")
}

func TestIntegration_PostsSortedByDate(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")

	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(filepath.Join(contentDir, "blog"), 0755)

	// Create required home.md
	writeTestFile(t, filepath.Join(contentDir, "home.md"), `---
title: Home
---
Welcome`)

	// Create posts in non-chronological order
	writeTestFile(t, filepath.Join(contentDir, "blog", "2021-01-01-oldest.md"), `---
title: Oldest Post
---
Content`)

	writeTestFile(t, filepath.Join(contentDir, "blog", "2023-12-25-newest.md"), `---
title: Newest Post
---
Content`)

	writeTestFile(t, filepath.Join(contentDir, "blog", "2022-06-15-middle.md"), `---
title: Middle Post
---
Content`)

	cfg := &model.Config{
		Title:      "Test Site",
		BaseURL:    "https://example.com",
		ContentDir: contentDir,
		OutputDir:  outputDir,
	}

	b := builder.New(cfg)
	err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	blogList := readTestFile(t, filepath.Join(outputDir, "blog", "index.html"))

	// Verify posts appear in correct order (newest first)
	newestIdx := strings.Index(blogList, "Newest Post")
	middleIdx := strings.Index(blogList, "Middle Post")
	oldestIdx := strings.Index(blogList, "Oldest Post")

	if newestIdx == -1 || middleIdx == -1 || oldestIdx == -1 {
		t.Fatal("blog list missing posts")
	}

	if newestIdx > middleIdx {
		t.Error("Newest post should appear before middle post")
	}
	if middleIdx > oldestIdx {
		t.Error("Middle post should appear before oldest post")
	}
}

// Helper functions

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(content)
}

func verifyFileExists(t *testing.T, dir, filename string) {
	t.Helper()
	path := filepath.Join(dir, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", path)
	}
}

func assertContains(t *testing.T, content, substr, msg string) {
	t.Helper()
	if !strings.Contains(content, substr) {
		t.Errorf("%s: expected to contain %q", msg, substr)
	}
}
