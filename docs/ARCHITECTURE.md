# Architecture

> Version: 1.0.0
> Last Updated: 2026-01-19

## System Context

The SSG operates as a command-line tool that transforms a source directory into a deployable static website.

```
┌─────────────────┐     ┌─────────────┐     ┌──────────────────┐
│  Source Files   │ --> │     SSG     │ --> │  Static Website  │
│  - Markdown     │     │   (build)   │     │  - HTML          │
│  - Assets       │     │             │     │  - CSS           │
│  - Config       │     │             │     │  - Assets        │
└─────────────────┘     └─────────────┘     └──────────────────┘
```

## Component Overview

```
┌────────────────────────────────────────────────────────────────────┐
│                              CLI                                    │
│  Handles user commands: build, serve, version                       │
└────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌────────────────────────────────────────────────────────────────────┐
│                            Config                                   │
│  Loads, validates, and provides site configuration                  │
└────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌────────────────────────────────────────────────────────────────────┐
│                           Builder                                   │
│  Orchestrates the entire build process                              │
└────────────────────────────────────────────────────────────────────┘
          │                      │                      │
          ▼                      ▼                      ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│     Scanner      │  │     Parser       │  │    Renderer      │
│ Finds content    │  │ Extracts meta,   │  │ Generates HTML   │
│ files            │  │ converts MD      │  │ from templates   │
└──────────────────┘  └──────────────────┘  └──────────────────┘
                                                      │
                                                      ▼
                               ┌──────────────────────────────────┐
                               │             Writer               │
                               │  Writes files to output dir      │
                               └──────────────────────────────────┘
```

## Components

### CLI

**Responsibility**: User interaction and command routing

**Commands**:
- `build` — Generate static site
- `serve` — Run development server
- `version` — Display version information

**Behavior**:
- Parses command-line arguments and flags
- Validates input parameters
- Invokes appropriate subsystems
- Reports success/failure with exit codes

### Config

**Responsibility**: Configuration management

**Functions**:
- Load configuration from YAML file
- Apply default values for optional fields
- Validate required fields
- Merge CLI overrides with file configuration

**Outputs**: Validated configuration object

### Scanner

**Responsibility**: Content discovery

**Functions**:
- Traverse content directory
- Identify page files (root-level Markdown)
- Identify post files (`blog/` subdirectory Markdown)
- Filter out index files (prefixed with `_`)
- Verify required files exist (`home.md`)

**Outputs**: Lists of page and post file paths

### Parser

**Responsibility**: Content processing

**Functions**:
- Extract YAML frontmatter
- Convert Markdown to HTML
- Extract dates from filenames
- Calculate word counts
- Detect asset references in content

**Outputs**: Structured Page and Post objects with metadata

### Renderer

**Responsibility**: HTML generation

**Functions**:
- Load and manage templates
- Bind data to templates
- Generate HTML for pages, posts, listings
- Apply template inheritance (partials)

**Outputs**: Rendered HTML strings

### Writer

**Responsibility**: File system operations

**Functions**:
- Validate output directory safety
- Clean output directory
- Write HTML files with clean URL structure
- Copy static assets
- Copy post-specific assets
- Generate SEO files (robots.txt, sitemap.xml)

**Outputs**: Complete static website in output directory

### Server (Dev Mode Only)

**Responsibility**: Local development server

**Functions**:
- Serve static files
- Handle clean URLs
- Provide graceful shutdown

## Data Flow

```
1. INPUT
   ├── ssg.yaml (configuration)
   ├── content/
   │   ├── home.md
   │   ├── *.md (pages)
   │   └── blog/
   │       ├── YYYY-MM-DD-slug.md (posts)
   │       └── assets/ (post assets)
   └── assets/ (global assets)

2. PROCESSING
   Config.Load() → Site configuration
        │
        ▼
   Scanner.Scan() → File paths
        │
        ▼
   Parser.Parse() → Page/Post objects
        │
        ▼
   Renderer.Render() → HTML content
        │
        ▼
   Writer.Write() → Files on disk

3. OUTPUT
   public/
   ├── index.html
   ├── 404.html
   ├── robots.txt
   ├── sitemap.xml
   ├── style.css
   ├── feed/index.xml
   ├── {page}/index.html
   └── blog/
       ├── index.html
       └── {post}/index.html
```

## Extension Points

### Custom Templates

Templates can be overridden by providing custom template files. The system should check for user-provided templates before falling back to defaults.

### Custom Stylesheets

A custom `style.css` in the assets directory replaces the default stylesheet entirely.

### Additional Content Types

The architecture supports extension to new content types by:
1. Adding new scanner patterns
2. Creating new model structures
3. Adding new templates
4. Extending the build pipeline

### Plugin System (Future)

The component-based architecture allows for future plugin hooks at:
- Post-parse (content transformation)
- Pre-render (template customization)
- Post-write (deployment automation)

## Error Handling

### Validation Errors
- Missing required config fields → Build failure with message
- Missing `home.md` → Build failure with message
- Invalid frontmatter → Build failure with file path

### Safety Errors
- Invalid output directory (root, home, etc.) → Build failure
- Missing referenced assets → Build failure with details

### Runtime Errors
- File write failures → Build failure with path
- Template errors → Build failure with details

## Performance Considerations

### Build Performance
- Content files processed in parallel where possible
- Templates parsed once and reused
- Assets copied without modification

### Memory Usage
- Stream large files when copying
- Process files individually rather than loading all into memory
- Clean up intermediate objects after use

## Security Considerations

### Output Directory Validation
- Never delete root directory (`/`)
- Never delete home directory (`~`)
- Never delete project directory
- Reject paths with parent traversal (`..`)

### Path Traversal Prevention
- Sanitize all file paths
- Validate paths are within expected directories
- Reject absolute paths in content references
