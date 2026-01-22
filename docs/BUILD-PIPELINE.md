# Build Pipeline

> Version: 1.0.0
> Last Updated: 2026-01-19

## Overview

The build pipeline transforms source content into a deployable static website through six sequential phases.

```
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│  Config  │-->│  Scan    │-->│  Parse   │-->│  Render  │-->│  Assets  │-->│   SEO    │
│  Load    │   │  Content │   │  Content │   │   HTML   │   │  Copy    │   │  Generate│
└──────────┘   └──────────┘   └──────────┘   └──────────┘   └──────────┘   └──────────┘
```

## Phase 1: Configuration Loading

**Purpose**: Load and validate site configuration

### Steps

1. Locate configuration file (`ssg.yaml`)
2. Parse YAML content
3. Apply default values
4. Merge CLI overrides
5. Validate required fields

### Validation Rules

**Required**:
- `site.title` must be non-empty
- `site.baseURL` must be non-empty

**Defaults Applied**:
| Field | Default Value |
|-------|---------------|
| `build.content` | `"content"` |
| `build.output` | `"public"` |
| `build.assets` | `"assets"` |

### CLI Override Precedence

```
CLI flag > Config file > Default value
```

### Failure Conditions

- Config file not found → Error
- Invalid YAML syntax → Error
- Missing required field → Error

## Phase 2: Content Scanning

**Purpose**: Discover all content files

### Steps

1. Validate output directory is safe
2. Verify `home.md` exists
3. Scan content root for page files (`*.md`)
4. Scan `blog/` subdirectory for post files (`*.md`)
5. Filter out index files (`_*.md`)

### Output Directory Validation

**Rejected Paths**:
- Root directory (`/`)
- Home directory (`~` or `$HOME`)
- Current project directory
- Paths containing `..`
- Paths outside project directory

### Discovery Rules

**Pages**:
- Location: `{content}/`
- Pattern: `*.md`
- Exclude: `_*.md` (index files)

**Posts**:
- Location: `{content}/blog/`
- Pattern: `*.md`
- Exclude: `_*.md` (index files)
- Validate: Filename matches `YYYY-MM-DD-*.md`

### Failure Conditions

- `home.md` not found → Error
- Output directory unsafe → Error
- Invalid post filename format → Error

## Phase 3: Content Parsing

**Purpose**: Extract metadata and convert content

### Steps

For each content file:

1. Read file content
2. Extract YAML frontmatter
3. Parse frontmatter fields
4. Extract date (for posts)
5. Convert Markdown to HTML
6. Calculate word count (for posts)
7. Extract asset references (for posts)

### Frontmatter Extraction

**Format**:
```
---
YAML content
---
Markdown content
```

**Algorithm**:
1. Check file starts with `---\n`
2. Find closing `---\n`
3. Parse content between delimiters as YAML
4. Rest of file is Markdown body

### Date Extraction (Posts)

**Primary**: Filename pattern
```
YYYY-MM-DD-slug.md
└─┬──┴─┴─┘
  date portion
```

**Fallback**: Frontmatter `date` field (also used as override)

### Markdown Conversion

**Input**: Markdown text
**Output**: HTML string

**Features Required**:
- Headings (`# h1` through `###### h6`)
- Paragraphs
- Bold, italic, strikethrough
- Links (external and internal)
- Images with alt text
- Code blocks with language hint
- Inline code
- Blockquotes
- Ordered and unordered lists
- Horizontal rules
- Tables (GFM extension)

### Word Count Calculation

1. Convert Markdown to HTML
2. Remove `<pre>...</pre>` blocks
3. Strip all HTML tags
4. Split on whitespace
5. Count non-empty tokens

### Asset Reference Extraction

**Pattern**: `!\[.*?\]\((assets/[^)]+)\)`

**Extracted**: Relative path after `assets/`

### Failure Conditions

- Invalid frontmatter YAML → Error with file path
- Missing required `title` field → Error with file path

## Phase 4: HTML Rendering

**Purpose**: Generate HTML from templates and data

### Steps

1. Prepare output directory (clean if exists)
2. Render static pages
3. Render blog posts
4. Render blog listing (if posts exist)
5. Render 404 page

### Page Rendering

For each page:

1. Prepare template data:
   - Page title
   - Page content (HTML)
   - Site configuration
   - Canonical URL
2. Render with page template
3. Write to `{output}/{slug}/index.html`
4. Special: Homepage writes to `{output}/index.html`

### Post Rendering

For each post (sorted by date descending):

1. Prepare template data:
   - Post title, date, summary
   - Post content (HTML) with rewritten asset paths
   - Word count
   - Site configuration
   - Canonical URL
2. Render with blog post template
3. Write to `{output}/blog/{slug}/index.html`
4. Copy referenced assets to post directory

### Blog Listing Rendering

If any posts exist:

1. Prepare listing data:
   - All posts (title, slug, date, word count)
   - Site configuration
2. Render with blog list template
3. Write to `{output}/blog/index.html`

### 404 Page Rendering

1. Prepare data:
   - Site configuration
   - Standard 404 message
2. Render with 404 template
3. Write to `{output}/404.html`

### Clean URL Structure

| Content | Output Path | URL |
|---------|-------------|-----|
| Homepage | `/index.html` | `/` |
| Page | `/{slug}/index.html` | `/{slug}/` |
| Post | `/blog/{slug}/index.html` | `/blog/{slug}/` |
| Blog list | `/blog/index.html` | `/blog/` |

## Phase 5: Asset Management

**Purpose**: Copy static assets to output

### Steps

1. Copy global assets from assets directory
2. Apply default stylesheet if none provided
3. Copy post-specific assets

### Global Assets

**Source**: `{assets}/`
**Destination**: `{output}/`

All files copied maintaining directory structure.

### Default Stylesheet

If `{assets}/style.css` does not exist:
- Embed default stylesheet
- Write to `{output}/style.css`

If `{assets}/style.css` exists:
- Copy user's stylesheet (overrides default)

### Post Assets

For each post with referenced assets:

**Source**: `{content}/blog/assets/{filename}`
**Destination**: `{output}/blog/{post-slug}/{filename}`

### Asset Path Rewriting

In post HTML content:
- Find: `src="assets/filename.ext"`
- Replace: `src="filename.ext"`

### Failure Conditions

- Referenced asset not found → Error with details
- Asset copy failure → Error with path

## Phase 6: SEO Artifact Generation

**Purpose**: Generate SEO-related files

### Sitemap Generation

**Output**: `{output}/sitemap.xml`

**Format**: XML Sitemap Protocol

**Contents**:
- All pages with `<loc>` (absolute URL)
- All posts with `<loc>` and `<lastmod>` (post date)

**Structure**:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
  </url>
  <url>
    <loc>https://example.com/blog/post-slug/</loc>
    <lastmod>2024-01-15</lastmod>
  </url>
</urlset>
```

### Robots.txt Generation

**Output**: `{output}/robots.txt`

**Contents**:
```
User-agent: *
Allow: /
Sitemap: {baseURL}/sitemap.xml
```

### RSS Feed Generation

**Output**: `{output}/feed.xml`

**Format**: RSS 2.0

**Contents**:
- Channel info (title, link, description)
- Up to 20 most recent posts
- Each item: title, link, description (CDATA HTML), pubDate

**Structure**:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Site Title</title>
    <link>https://example.com/</link>
    <description>Site description</description>
    <lastBuildDate>Mon, 15 Jan 2024 00:00:00 +0000</lastBuildDate>
    <item>
      <title>Post Title</title>
      <link>https://example.com/blog/post-slug/</link>
      <description><![CDATA[<p>Post content...</p>]]></description>
      <pubDate>Mon, 15 Jan 2024 00:00:00 +0000</pubDate>
    </item>
  </channel>
</rss>
```

**Date Format**: RFC 1123Z (`Mon, 02 Jan 2006 15:04:05 -0700`)

## Build Output Summary

```
{output}/
├── index.html           # Homepage
├── 404.html             # Error page
├── robots.txt           # Robots directive
├── sitemap.xml          # Sitemap
├── style.css            # Stylesheet
├── [global assets]      # Copied from assets/
├── feed.xml             # RSS feed
├── {page-slug}/
│   └── index.html       # Static page
└── blog/
    ├── index.html       # Blog listing
    └── {post-slug}/
        ├── index.html   # Blog post
        └── [assets]     # Post-specific assets
```

## Error Recovery

The build pipeline does not support partial builds. Any error causes complete build failure.

**Recommended**: Always build to a new output directory, then swap directories atomically for deployment.

## Performance Optimization

### Parallelization Opportunities

- Page parsing (independent files)
- Post parsing (independent files)
- Asset copying (independent files)

### Caching Opportunities

- Template compilation (once per build)
- Unchanged file detection (future enhancement)
