# Project Context

## Purpose
ssg is a minimal static site generator for blogs, written in Go. It converts markdown files with YAML frontmatter into a complete, clean-URL website with RSS feed, sitemap, SEO metadata, and automatic light/dark theme support.

## Tech Stack
- **Language**: Go 1.25
- **Markdown**: Goldmark (github.com/yuin/goldmark) - CommonMark compliant
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **Configuration**: YAML via gopkg.in/yaml.v3
- **Build**: Make with version injection via ldflags
- **Issue Tracking**: bd (beads) - stored in .beads/

## Project Conventions

### Code Style
- Standard Go formatting via `go fmt`
- Package names: lowercase, single-word (model, parser, builder, renderer)
- File names: lowercase with underscores for tests (*_test.go)
- Exported functions have doc comments starting with function name
- Unexported helpers are concise with minimal comments
- Error messages are lowercase, no punctuation at end

### Architecture Patterns
- **Package layout**: cmd/ssg (CLI entry), internal/* (core packages)
- **Core packages**:
  - `model` - Data types (Site, Page, Post, Config, NavItem)
  - `config` - YAML config loading
  - `parser` - Markdown parsing with frontmatter extraction
  - `builder` - Site orchestration (scanning, building, asset copying)
  - `renderer` - HTML template rendering
  - `assets` - Embedded default stylesheet
  - `server` - Development HTTP server
  - `wordcount` - Content analysis
- **Dependency flow**: builder depends on parser, renderer, model; no circular deps
- **Constructor pattern**: `New(cfg)` returns initialized structs
- **Setter injection**: `SetVersion()`, `SetAssetsDir()` for optional config

### Testing Strategy
- **Test location**: *_test.go files alongside source
- **Table-driven tests**: Standard pattern with named test cases
- **Parallel execution**: All tests use `t.Parallel()` for speed
- **Isolation**: `t.TempDir()` for file system tests, no shared state
- **Helper functions**: `writeFile(t, path, content)`, `writeHomeMd(t, dir)`
- **Naming**: TestFunctionName_Scenario for clarity
- **Coverage focus**: Business logic, edge cases, error conditions
- **Integration test**: integration_test.go at package root

### Git Workflow
- **Main branch**: `main`
- **Commit convention**: type(scope): message (e.g., feat(ssg-xxx): Add feature)
- **CI pipeline**: `make all` runs fmt → vet → test → build
- **Version injection**: Git SHA via `-ldflags "-X 'main.Version=$(VERSION)'"`
- **Issue tracking**: bd (beads) for all tasks, synced via `bd sync`

## Domain Context
- **Content structure**:
  - `content/home.md` - Required homepage
  - `content/*.md` - Static pages
  - `content/blog/YYYY-MM-DD-slug.md` - Blog posts with date in filename
  - `content/blog/assets/` - Co-located images for posts
- **Output structure**: Clean URLs (`/about/index.html` for `/about/`)
- **Generated files**: index.html, 404.html, robots.txt, sitemap.xml, feed.xml
- **Frontmatter fields**: title (required), summary, date (optional override)

## Important Constraints
- **home.md is required** - Build fails without it
- **Output directory validation** - Cannot be /, ~, ., .., or outside project
- **Date format**: Filename-based YYYY-MM-DD, frontmatter can override
- **Asset paths**: `assets/` prefix rewritten to co-locate with posts
- **No external dependencies at runtime** - Self-contained HTML output

## External Dependencies
- **Build time only**: Go toolchain, Make
- **Optional**: staticcheck for linting
- **Runtime**: None (generates static HTML)
- **Dev server**: Built-in, no external server needed
