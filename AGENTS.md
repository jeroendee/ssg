# Project Context

Project-specific instructions and guidelines:

## Overview

**ssg** is a minimal static site generator for blogs, written in Go. It converts Markdown files with YAML frontmatter into a complete static website with clean URLs, RSS feed, sitemap, SEO metadata, and automatic light/dark theme support.

## Tech Stack

- **Language**: Go 1.25
- **Markdown**: Goldmark (CommonMark compliant, with GFM table extension)
- **CLI**: Cobra
- **Config**: YAML via `gopkg.in/yaml.v3`
- **Build**: Make with version injection via ldflags
- **Change management**: OpenSpec (artifact-driven workflow in `openspec/`)

## Architecture

```
cmd/ssg/           CLI entry point (Cobra commands: build, serve, version)
internal/
  model/           Data types: Site, Page, Post, Config, NavItem, FeedItem
  config/          YAML config loading with CLI override support
  parser/          Markdown parsing, frontmatter extraction, date anchors
  builder/         Site orchestration: scanning, rendering, asset copying
  renderer/        HTML template rendering (embedded templates)
  assets/          Embedded default stylesheet
  server/          Development HTTP server
  wordcount/       Word count for blog posts
```

**Dependency flow**: `cmd/ssg` → `builder` → `parser`, `renderer`, `model`, `assets`. No circular dependencies.

## Key Conventions

### Code Style
- Standard `go fmt` formatting
- Package names: lowercase, single-word
- Exported functions have doc comments starting with function name
- Error messages: lowercase, no trailing punctuation
- Constructor pattern: `New(cfg)` returns initialized structs
- Setter injection for optional config: `SetVersion()`, `SetAssetsDir()`

### Testing
- Test files live alongside source (`*_test.go`)
- Table-driven tests with named cases
- **All tests use `t.Parallel()`** — no exceptions
- File system isolation via `t.TempDir()`
- Helper pattern: `writeFile(t, path, content)`, `writeHomeMd(t, dir)`
- Integration tests in `integration_test.go` at package root
- Run all tests: `make test` or `go test ./...`

### Build & Development
- `make all` — CI pipeline: fmt → vet → test → build
- `make dev` — quick restart: kill server, copy assets, serve
- `make serve` — full workflow: kill, assets, build, serve on `:8080`
- Binary output: `bin/ssg`
- Version injected from git SHA via ldflags

## Content Model

- **`content/home.md`** — Required. Build fails without it.
- **`content/*.md`** — Static pages. Slug derived from filename.
- **`content/_footer.md`** — Optional. Rendered as site-wide footer.
- **`content/blog/YYYY-MM-DD-slug.md`** — Blog posts. Date extracted from filename, overridable via frontmatter `date` field.
- **`content/blog/assets/`** — Co-located images for posts. Paths rewritten at build time.
- **Frontmatter fields**: `title` (required), `summary`, `date` (optional override)
- **Date-anchored pages**: Headings like `# *2026-01-15*` create navigable date sections with archive grouping.

## Output

Clean URLs: `content/about.md` → `public/about/index.html` → `/about/`

Generated files: `index.html`, `404.html`, `robots.txt`, `sitemap.xml`, `feed.xml`, `build.json`, `style.css`

## Important Constraints

- **Output directory validation**: Builder rejects `/`, `~`, `.`, `..`, or paths outside the project root before `rm -rf`.
- **Feed pages**: Configurable in `ssg.yaml` under `feed.pages` — date-anchored page sections appear in RSS alongside blog posts.
- **Default stylesheet**: Embedded in `internal/assets/`. A custom `assets/style.css` overrides it.
- **No runtime dependencies**: Output is self-contained static HTML.

## Configuration

Config file: `ssg.yaml` (see `ssg.yaml.example` for full reference).

Required fields: `site.title`, `site.baseURL`. Everything else has defaults.

## Development Site

The `dev/` directory contains a working site for local development. It is gitignored and used by `make serve` / `make dev`.
