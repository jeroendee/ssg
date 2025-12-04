# 202512031308 - original Static Site Generator specification

## Objective

Static Site Generator (ssg) should build a static website based on markdown files. The site cannot use any javascript. Some styling is allowed. 

- There should be a `blog` section, that contains blog entries per date. 
- There should be the possibility to use image assets. Like the use of a favicon. 
- There should be an automatic dark / light theme based on time.
- The site should be reponsive. 

The final product should be: `https://www.qualityshepherd.nl` scrape that website for it's contents. 
- https://www.qualityshepherd.nl
- https://www.qualityshepherd.nl/about/
- https://www.qualityshepherd.nl/contact/
- https://www.qualityshepherd.nl/values/
- https://www.qualityshepherd.nl/blog/

`ssg` should be a command line tool that will parse a set of directories / files and will output a static site. 

## Technology

ssg should be built with the Go programming language. 

- use Cobra (https://cobra.dev) for the command line interface. Fetch and read all the docs: https://cobra.dev/docs/ 
- use goldmark for the markdown parsing. Fetch and read: https://github.com/yuin/goldmark 
- fetch and read: https://www.qualityshepherd.nl => site to generate by `ssg`
- fetch and read: https://github.com/janraasch/hugo-bearblog/ => hugo theme used
- fetch and read: https://pkg.go.dev/html/template => Go html template engine, if this is needed at all. 
- fetch and read: https://www.fastmail.help/hc/en-us/articles/1500000280141-How-to-set-up-a-website => static website hosting at Fastmail

### Skills

Use these Skills (invoke at appropriate phases):

| Phase | Skill Invocation | Purpose |
|-------|------------------|---------|
| Setup | `skill: "go-scaffolder"` | Scaffold Mixed Repository project structure |
| All coding | `skill: "tdd-workflow"` | Red-Green-Refactor for every implementation |
| All coding | `skill: "go-writer"` | Idiomatic Go code following Google Style Guide |
| All testing | `skill: "go-tester"` | Table-driven tests, parallel execution, coverage |
| Documentation | `skill: "go-documenter"` | Package docs, doc.go files, exported symbols |
| Code review | `skill: "go-reviewer"` | Manual code review per Go wiki standards |
| Code review | `skill: "go-analyzer"` | Static analysis with gopls, staticcheck |

## Examples

## Implementation Plan

### Configuration

- **Module path**: `github.com/jeroendee/ssg`
- **Project type**: Mixed Repository (library + CLI)
- **Content directory**: `content/` (Hugo-style)
- **Configuration**: `ssg.yaml` with CLI flag overrides
- **Output directory**: `public/` (default)

### Architecture

```
ssg/
├── cmd/ssg/main.go           # CLI entry point (Cobra)
├── internal/
│   ├── config/               # Configuration loading (YAML + flags)
│   ├── parser/               # Markdown parsing (Goldmark)
│   ├── builder/              # Site building orchestration
│   ├── renderer/             # HTML rendering (html/template)
│   └── model/                # Domain models (Site, Page, Post)
├── templates/                # HTML templates (embedded)
├── assets/                   # Static assets (CSS, favicon)
├── go.mod
└── ssg.yaml.example
```

### Tasks

#### Phase 1: Project Setup

**Task 1**: Scaffold Go project
- Invoke: `skill: "go-scaffolder"` (Mixed Repository pattern)
- Create directory structure per architecture above
- Initialize `go.mod` with `github.com/jeroendee/ssg`

**Task 2**: Add external dependencies
- Add `github.com/spf13/cobra` (CLI framework)
- Add `github.com/yuin/goldmark` (Markdown parser)
- Add `gopkg.in/yaml.v3` (YAML config parsing)
- Run `go mod tidy`

#### Phase 2: Domain Models

**Task 3**: Implement `internal/model` package
- Invoke: `skill: "tdd-workflow"`, `skill: "go-writer"`, `skill: "go-tester"`
- Define `Site` struct (title, baseURL, author, navigation, pages, posts)
- Define `Page` struct (title, slug, content, path, template)
- Define `Post` struct (extends Page with date, summary)
- Define `Config` struct (site metadata, paths, build options)

#### Phase 3: Configuration

**Task 4**: Implement `internal/config` package
- Invoke: `skill: "tdd-workflow"`, `skill: "go-writer"`, `skill: "go-tester"`
- Load configuration from `ssg.yaml`
- Support CLI flag overrides
- Validate required fields
- Define default values

#### Phase 4: Markdown Parsing

**Task 5**: Implement `internal/parser` package
- Invoke: `skill: "tdd-workflow"`, `skill: "go-writer"`, `skill: "go-tester"`
- Parse markdown files using Goldmark
- Extract frontmatter (YAML) for metadata
- Convert markdown to HTML
- Handle blog posts (date from filename `YYYY-MM-DD-title.md` or frontmatter)

#### Phase 5: HTML Rendering

**Task 6**: Implement `internal/renderer` package
- Invoke: `skill: "tdd-workflow"`, `skill: "go-writer"`, `skill: "go-tester"`
- Create HTML templates matching qualityshepherd.nl design
- Base template with navigation, header, footer
- Page template for static pages
- Blog list template (date + title format)
- Blog post template
- Embed templates using `embed.FS`

**Task 7**: Implement CSS styling
- Port BearBlog-style CSS from qualityshepherd.nl
- Responsive design (max-width: 720px)
- Dark mode via `prefers-color-scheme` media query
- Blog post list styling (flexbox with fixed date width)

#### Phase 6: Site Building

**Task 8**: Implement `internal/builder` package
- Invoke: `skill: "tdd-workflow"`, `skill: "go-writer"`, `skill: "go-tester"`
- Scan `content/` directory for markdown files
- Identify pages vs blog posts (by directory)
- Parse all content files
- Generate site navigation from pages
- Sort blog posts by date (newest first)

**Task 9**: Implement output generation
- Invoke: `skill: "tdd-workflow"`, `skill: "go-writer"`, `skill: "go-tester"`
- Create output directory structure
- Write HTML files (clean URLs: `/about/index.html`)
- Copy static assets (CSS, favicon, images)
- Generate `index.html` for homepage
- Generate `404.html` for Fastmail compatibility

#### Phase 7: CLI Implementation

**Task 10**: Implement Cobra CLI in `cmd/ssg`
- Invoke: `skill: "tdd-workflow"`, `skill: "go-writer"`, `skill: "go-tester"`
- Root command with version info
- `build` subcommand (default action)
- CLI flags: `--config`, `--output`, `--content`
- Help text and usage documentation

#### Phase 8: Integration & Polish

**Task 11**: Integration testing
- Invoke: `skill: "go-tester"`
- Create sample content matching qualityshepherd.nl structure
- Test full build pipeline
- Verify HTML output correctness
- Test dark/light theme CSS

**Task 12**: Documentation and review
- Invoke: `skill: "go-documenter"`, `skill: "go-reviewer"`, `skill: "go-analyzer"`
- Add package documentation (doc.go files)
- Create example `ssg.yaml.example`
- Run static analysis (gopls, staticcheck)
- Code review per Go wiki standards

### Content Structure (Expected Input)

```
content/
├── _index.md           # Homepage
├── about.md            # About page
├── contact.md          # Contact page
├── values.md           # Values page
└── blog/
    ├── _index.md       # Blog listing config
    └── 2021-03-26-post-title.md
```

### Configuration File (ssg.yaml)

```yaml
site:
  title: "Quality Shepherd"
  baseURL: "https://www.qualityshepherd.nl"
  author: "Jeroen"

build:
  content: "content"
  output: "public"

navigation:
  - title: "Home"
    url: "/"
  - title: "About"
    url: "/about/"
  - title: "Contact"
    url: "/contact/"
  - title: "Values"
    url: "/values/"
  - title: "Blog"
    url: "/blog/"
```

### Key Design Decisions

1. **No JavaScript**: Pure HTML/CSS as per spec
2. **CSS dark mode**: Uses `prefers-color-scheme` (automatic, browser-based)
3. **Embedded templates**: Templates compiled into binary via `embed.FS`
4. **Frontmatter**: YAML frontmatter in markdown for page metadata
5. **Clean URLs**: Output as `/about/index.html` accessible as `/about/`
6. **Fastmail compatible**: Includes `index.html` and `404.html`

### External Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| github.com/spf13/cobra | latest | CLI framework |
| github.com/yuin/goldmark | latest | Markdown to HTML |
| gopkg.in/yaml.v3 | latest | YAML config parsing |

### Success Criteria

1. `ssg build` generates complete static site in `public/`
2. Output matches qualityshepherd.nl in structure and styling
3. All tests pass with high coverage of business logic
4. Code follows Google Go Style Guide
5. Documentation complete per go-documenter guidelines
