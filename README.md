# ssg

A minimal static site generator for blogs, written in Go.

## Features

- Markdown to HTML conversion with [Goldmark](https://github.com/yuin/goldmark)
- YAML frontmatter support
- Blog posts with automatic date extraction from filenames
- Static pages
- Clean URLs (`/about/` instead of `/about.html`)
- Configurable navigation menu
- Asset copying (CSS, images, etc.)
- Auto-generated homepage, blog listing, and 404 page
- RSS 2.0 feed generation for blog posts
- Development server for local preview

## Installation

```bash
go install github.com/jeroendee/ssg/cmd/ssg@latest
```

Or build from source:

```bash
git clone https://github.com/jeroendee/ssg.git
cd ssg
go build -o ssg ./cmd/ssg
```

## Quick Start

1. Create a configuration file `ssg.yaml`:

```yaml
site:
  title: My Blog
  baseURL: https://example.com
  author: Your Name
  logo: /logo.svg
  favicon: /favicon.svg

build:
  content: content
  output: public

navigation:
  - title: Home
    url: /
  - title: About
    url: /about/
  - title: Blog
    url: /blog/
```

2. Create your content directory structure:

```
project/
├── ssg.yaml
├── assets/
│   └── style.css
└── content/
    ├── about.md
    └── blog/
        └── 2024-01-15-my-first-post.md
```

3. Build the site:

```bash
ssg build
```

Your generated site will be in the `public/` directory.

## Usage

### Build

```
ssg build [flags]

Flags:
  -c, --config string    path to config file (default "ssg.yaml")
  -o, --output string    output directory (overrides config)
      --content string   content directory (overrides config)
      --assets string    assets directory (default "assets")
```

### Serve

Start a local development server to preview your site:

```
ssg serve [flags]

Flags:
  -c, --config string    path to config file (default "ssg.yaml")
  -p, --port int         port to serve on (default 8080)
  -d, --dir string       directory to serve (overrides config output_dir)
  -b, --build            build site before serving
```

Examples:

```bash
ssg serve                     # Serve on http://localhost:8080
ssg serve -p 3000             # Serve on port 3000
ssg serve --build             # Build first, then serve
ssg serve --dir ./other-dir   # Serve a specific directory
```

Press `Ctrl+C` to stop the server.

## Content Format

### Pages

Static pages are markdown files in the root of your content directory:

```markdown
---
title: About
---

Your page content here in **markdown**.
```

Pages become clean URLs: `content/about.md` → `/about/`

### Blog Posts

Blog posts go in the `content/blog/` directory. Filenames must follow the format `YYYY-MM-DD-slug.md`:

```markdown
---
title: My First Post
---

Your post content here.
```

The date is extracted from the filename: `2024-01-15-my-first-post.md` → published January 15, 2024, accessible at `/blog/my-first-post/`

## Configuration

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `site.title` | Yes | - | Site title for header and page titles |
| `site.baseURL` | Yes | - | Root URL where the site will be hosted |
| `site.author` | No | - | Author name for metadata |
| `site.logo` | No | - | Path to site logo (e.g., `/logo.svg`) |
| `site.favicon` | No | - | Path to favicon (e.g., `/favicon.svg`) |
| `site.description` | No | - | Site description for RSS feed |
| `build.content` | No | `content` | Directory containing markdown files |
| `build.output` | No | `public` | Directory for generated HTML |
| `navigation` | No | - | List of navigation menu items |

## RSS Feed

An RSS 2.0 feed is automatically generated at `/feed/index.xml` containing the 20 most recent blog posts. The feed uses `site.description` if provided.

## Output Structure

```
public/
├── index.html          # Homepage
├── 404.html            # Error page
├── style.css           # Copied from assets/
├── feed/
│   └── index.xml       # RSS feed
├── about/
│   └── index.html      # About page
└── blog/
    ├── index.html      # Blog listing
    └── my-first-post/
        └── index.html  # Blog post
```

## Example

See the `example/` directory for a working site. To build it:

```bash
cd example
../ssg build
```

## License

MIT
