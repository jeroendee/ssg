# Content Model

> Version: 1.0.0
> Last Updated: 2026-01-19

## Overview

The SSG uses a hierarchical content model with two primary types: **Pages** and **Posts**. Posts extend Pages with additional blog-specific fields.

## Content Types

### Page

A static page represents non-blog content like homepage, about, contact pages.

**Fields**:

| Field | Type | Source | Description |
|-------|------|--------|-------------|
| `title` | String | Frontmatter | Page title (required) |
| `slug` | String | Filename | URL segment derived from filename |
| `content` | HTML | Markdown body | Rendered HTML content |
| `path` | String | Derived | Full URL path (e.g., `/about/`) |

**Slug Derivation**:
- Filename without extension becomes slug
- `about.md` → slug: `about` → path: `/about/`
- Special case: `home.md` → slug: `""` → path: `/`

### Post

A blog post extends Page with date, summary, and asset support.

**Fields** (inherits all Page fields):

| Field | Type | Source | Description |
|-------|------|--------|-------------|
| `date` | Date | Filename or frontmatter | Publication date |
| `summary` | String | Frontmatter | Brief description |
| `wordCount` | Integer | Calculated | Word count of content |
| `assets` | List | Extracted | Referenced asset paths |

**Date Derivation**:
- Primary: Extracted from filename (`YYYY-MM-DD-slug.md`)
- Override: Frontmatter `date` field takes precedence
- Format: ISO 8601 date (`2024-01-15`)

### Site

Aggregate object containing all site data.

**Fields**:

| Field | Type | Source | Description |
|-------|------|--------|-------------|
| `title` | String | Config | Site name |
| `baseURL` | String | Config | Root URL |
| `description` | String | Config | Site description |
| `author` | String | Config | Author name |
| `logo` | String | Config | Logo path |
| `favicon` | String | Config | Favicon path |
| `pages` | List[Page] | Content | All static pages |
| `posts` | List[Post] | Content | All blog posts (sorted by date desc) |
| `navigation` | List | Config | Navigation menu items |

## Frontmatter

Frontmatter is YAML metadata at the beginning of Markdown files, delimited by `---`.

### Page Frontmatter

```yaml
---
title: "About Me"
---
```

**Required Fields**:
- `title` — Page title displayed in heading and browser tab

### Post Frontmatter

```yaml
---
title: "My First Post"
summary: "An introduction to my blog"
date: "2024-01-15"
---
```

**Required Fields**:
- `title` — Post title

**Optional Fields**:
- `summary` — Brief description for listings and meta tags
- `date` — Override filename date (format: `YYYY-MM-DD`)

### Validation Rules

1. Frontmatter must be valid YAML
2. Frontmatter must be at the very beginning of the file
3. Opening `---` must be on line 1
4. Closing `---` must appear before content
5. `title` field is required for all content types

## Directory Structure

### Content Directory Layout

```
content/
├── home.md              # Homepage (required)
├── about.md             # Static page
├── contact.md           # Static page
└── blog/
    ├── 2024-01-15-first-post.md
    ├── 2024-01-20-second-post.md
    └── assets/
        ├── image1.jpg
        └── diagram.png
```

### Rules

1. **Homepage Required**: `home.md` must exist in content root
2. **Pages**: Markdown files (`.md`) in content root become pages
3. **Posts**: Markdown files in `blog/` subdirectory become posts
4. **Index Files**: Files prefixed with `_` are ignored
5. **Blog Directory**: Optional; site builds without it

### File Discovery

**Scanner Behavior**:
1. Scan content root for `*.md` files (pages)
2. Verify `home.md` exists
3. Scan `content/blog/` for `*.md` files (posts)
4. Skip files starting with `_`
5. Return separate lists for pages and posts

## Filename Conventions

### Pages

```
{slug}.md
```

Examples:
- `home.md` → `/` (special case)
- `about.md` → `/about/`
- `contact.md` → `/contact/`
- `privacy-policy.md` → `/privacy-policy/`

### Posts

```
YYYY-MM-DD-{slug}.md
```

Examples:
- `2024-01-15-first-post.md` → `/blog/first-post/`
- `2024-12-25-christmas-update.md` → `/blog/christmas-update/`

**Date Extraction Rules**:
1. Filename must start with `YYYY-MM-DD-`
2. Year: 4 digits
3. Month: 2 digits (01-12)
4. Day: 2 digits (01-31)
5. Separator: single hyphen after date
6. Rest of filename becomes slug

## Asset References

### Syntax

In Markdown content, reference assets using relative paths:

```markdown
![Alt text](assets/image.jpg)
```

### Asset Extraction

The parser extracts asset references using pattern matching:

**Pattern**: `![...](assets/...)`

**Examples**:
```markdown
![Photo](assets/photo.jpg)           → assets/photo.jpg
![Diagram](assets/diagrams/arch.png) → assets/diagrams/arch.png
```

### Asset Processing

1. Parser extracts all asset references from content
2. Builder validates all referenced assets exist
3. Assets copied to post output directory
4. Paths rewritten in HTML (remove `assets/` prefix)

**Before** (Markdown):
```markdown
![Photo](assets/photo.jpg)
```

**After** (HTML in output):
```html
<img src="photo.jpg" alt="Photo">
```

## Word Count Calculation

### Algorithm

1. Convert Markdown to HTML
2. Remove `<pre>` blocks (code blocks excluded)
3. Strip all HTML tags
4. Split on whitespace
5. Count non-empty tokens

### Rules

- Code blocks (`<pre>` tags) are excluded
- Inline code is included
- HTML entities count as single words
- Only visible text is counted

## Content Sorting

### Posts

Posts are always sorted by date in descending order (newest first).

**Sort Key**: Publication date (from filename or frontmatter)

### Pages

Pages maintain discovery order (filesystem order).

## Output Mapping

### Pages

| Input | Output |
|-------|--------|
| `content/home.md` | `public/index.html` |
| `content/about.md` | `public/about/index.html` |
| `content/contact.md` | `public/contact/index.html` |

### Posts

| Input | Output |
|-------|--------|
| `content/blog/2024-01-15-first-post.md` | `public/blog/first-post/index.html` |
| `content/blog/assets/photo.jpg` | `public/blog/first-post/photo.jpg` |

### Special Files

| Generated | Output |
|-----------|--------|
| Blog listing | `public/blog/index.html` |
| RSS feed | `public/feed/index.xml` |
| Sitemap | `public/sitemap.xml` |
| Robots | `public/robots.txt` |
| 404 page | `public/404.html` |
