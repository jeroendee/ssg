# SSG Documentation

> Version: 1.0.0
> Last Updated: 2026-01-19

## Overview

SSG is a minimal, opinionated static site generator designed for personal blogs and simple websites. It transforms Markdown content into a fully-featured static website with:

- **Clean URLs** — No `.html` extensions in URLs
- **Blog Support** — Built-in blog with RSS feeds and listings
- **SEO Ready** — Sitemaps, meta tags, Open Graph, JSON-LD structured data
- **Fast** — Single binary, no runtime dependencies
- **Themeable** — Customizable templates and styles

## Quick Reference

| Concept | Description |
|---------|-------------|
| Content Directory | Where Markdown files live (default: `content/`) |
| Output Directory | Where generated site goes (default: `public/`) |
| Pages | Markdown files in content root → static pages |
| Posts | Markdown files in `content/blog/` → blog posts |
| Clean URLs | `/about/` instead of `/about.html` |
| Frontmatter | YAML metadata at top of Markdown files |

## Documentation Index

### Core Specifications

| Document | Description |
|----------|-------------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | System components and data flow |
| [CONTENT-MODEL.md](CONTENT-MODEL.md) | Page and Post structures, frontmatter |
| [BUILD-PIPELINE.md](BUILD-PIPELINE.md) | Build phases and processing order |
| [CONFIGURATION.md](CONFIGURATION.md) | Configuration schema and validation |
| [TEMPLATING.md](TEMPLATING.md) | Template system and theming |
| [SEO.md](SEO.md) | SEO features and metadata |
| [CLI.md](CLI.md) | Command-line interface |
| [DEV-SERVER.md](DEV-SERVER.md) | Development server behavior |

### Architecture Decision Records

| ADR | Title |
|-----|-------|
| [ADR-0001](decisions/ADR-0001-clean-urls.md) | Clean URL Pattern |
| [ADR-0002](decisions/ADR-0002-embedded-templates.md) | Embedded Templates |
| [ADR-0003](decisions/ADR-0003-frontmatter-yaml.md) | YAML Frontmatter |
| [ADR-0004](decisions/ADR-0004-date-from-filename.md) | Date from Filename |
| [ADR-0005](decisions/ADR-0005-no-javascript.md) | No JavaScript |
| [ADR-0006](decisions/ADR-0006-post-asset-colocation.md) | Post Asset Colocation |
| [ADR-TEMPLATE](decisions/ADR-TEMPLATE.md) | Template for new ADRs |

### Machine-Readable Schemas

| Schema | Description |
|--------|-------------|
| [config.schema.json](schemas/config.schema.json) | Configuration file validation |
| [frontmatter.schema.json](schemas/frontmatter.schema.json) | Content frontmatter validation |

## Project Structure

```
project/
├── ssg.yaml              # Configuration file
├── content/
│   ├── home.md           # Homepage (required)
│   ├── about.md          # Static page → /about/
│   └── blog/
│       ├── 2024-01-15-first-post.md
│       └── assets/       # Post assets
├── assets/
│   ├── style.css         # Custom stylesheet (optional)
│   ├── favicon.svg       # Site favicon
│   └── logo.svg          # Site logo
└── public/               # Generated output
```

## Generated Output

```
public/
├── index.html            # Homepage
├── 404.html              # Error page
├── robots.txt            # SEO file
├── sitemap.xml           # Site map
├── style.css             # Stylesheet
├── feed/
│   └── index.xml         # RSS feed
├── about/
│   └── index.html        # Static page
└── blog/
    ├── index.html        # Blog listing
    └── first-post/
        └── index.html    # Blog post
```

## Design Philosophy

1. **Simplicity** — Convention over configuration
2. **Speed** — Fast builds, no runtime dependencies
3. **SEO** — Built-in best practices for search engines
4. **Portability** — Pure HTML/CSS output, works anywhere
5. **Minimal** — No JavaScript by default

## Implementation Notes

This documentation describes WHAT the SSG does, not HOW it's implemented. The specifications are language-agnostic and can be used as a reference for implementing the SSG in any programming language.
