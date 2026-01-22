# SEO

> Version: 1.0.0
> Last Updated: 2026-01-19

## Overview

The SSG generates comprehensive SEO artifacts including meta tags, structured data, sitemaps, RSS feeds, and robots directives.

## Meta Tags

### Standard Meta Tags

Generated for all pages:

```html
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{PageTitle} - {SiteTitle}</title>
<meta name="description" content="{Summary}">
<link rel="canonical" href="{CanonicalURL}">
```

### Title Tag

**Format**: `{PageTitle} - {SiteTitle}`

**Examples**:
- Homepage: `Home - My Blog`
- About page: `About - My Blog`
- Blog post: `My First Post - My Blog`

### Description Meta Tag

**Source Priority**:
1. Page/post `summary` from frontmatter
2. Site `description` from config
3. Empty if neither available

**Recommendations**:
- Keep under 160 characters
- Include relevant keywords naturally
- Make it compelling for click-through

### Canonical URL

**Format**: Absolute URL to the current page

**Examples**:
- `https://example.com/`
- `https://example.com/about/`
- `https://example.com/blog/my-post/`

## Open Graph Tags

Open Graph tags enable rich previews when sharing on social platforms.

### All Pages

```html
<meta property="og:title" content="{PageTitle}">
<meta property="og:description" content="{Summary}">
<meta property="og:url" content="{CanonicalURL}">
<meta property="og:site_name" content="{SiteTitle}">
```

### Page Type

| Content Type | og:type |
|--------------|---------|
| Static pages | `website` |
| Blog posts | `article` |

```html
<meta property="og:type" content="website">
<!-- or -->
<meta property="og:type" content="article">
```

### Image Tag

If `site.logo` is configured:

```html
<meta property="og:image" content="{BaseURL}{Logo}">
```

## Twitter Cards

Twitter Card tags for Twitter/X previews:

```html
<meta name="twitter:card" content="summary">
<meta name="twitter:title" content="{PageTitle}">
<meta name="twitter:description" content="{Summary}">
```

If `site.logo` is configured:

```html
<meta name="twitter:image" content="{BaseURL}{Logo}">
```

### Card Types

| Type | Description |
|------|-------------|
| `summary` | Small square image with title and description |

**Note**: `summary_large_image` can be used for larger images if desired.

## JSON-LD Structured Data

Structured data helps search engines understand page content.

### WebSite Schema

Generated on all pages:

```json
{
  "@context": "https://schema.org",
  "@type": "WebSite",
  "name": "{SiteTitle}",
  "url": "{BaseURL}",
  "description": "{SiteDescription}"
}
```

### Article Schema

Generated on blog posts only:

```json
{
  "@context": "https://schema.org",
  "@type": "Article",
  "headline": "{PostTitle}",
  "datePublished": "{PostDate}",
  "description": "{PostSummary}",
  "author": {
    "@type": "Person",
    "name": "{SiteAuthor}"
  }
}
```

**Date Format**: ISO 8601 (`2024-01-15`)

**Note**: Article schema only includes author if `site.author` is configured.

## Sitemap

### Location

`{output}/sitemap.xml`

### Format

XML Sitemap Protocol (https://www.sitemaps.org/protocol.html)

### Structure

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <!-- Static pages -->
  <url>
    <loc>https://example.com/</loc>
  </url>
  <url>
    <loc>https://example.com/about/</loc>
  </url>

  <!-- Blog posts with dates -->
  <url>
    <loc>https://example.com/blog/my-post/</loc>
    <lastmod>2024-01-15</lastmod>
  </url>
</urlset>
```

### Contents

| Content | Included | lastmod |
|---------|----------|---------|
| Homepage | Yes | No |
| Static pages | Yes | No |
| Blog posts | Yes | Post date |
| Blog listing | No | - |
| 404 page | No | - |
| RSS feed | No | - |

### URL Format

- All URLs are absolute
- URLs use clean format with trailing slash
- No query parameters or fragments

### Date Format

ISO 8601 date: `YYYY-MM-DD`

## RSS Feed

### Location

`{output}/feed.xml`

**Note**: Located at `/feed.xml`

### Format

RSS 2.0 specification

### Structure

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>{SiteTitle}</title>
    <link>{BaseURL}</link>
    <description>{SiteDescription}</description>
    <lastBuildDate>{BuildDate}</lastBuildDate>

    <item>
      <title>{PostTitle}</title>
      <link>{PostURL}</link>
      <description><![CDATA[{PostHTML}]]></description>
      <pubDate>{PostDate}</pubDate>
    </item>
    <!-- ... more items ... -->
  </channel>
</rss>
```

### Channel Elements

| Element | Source | Required |
|---------|--------|----------|
| `title` | `site.title` | Yes |
| `link` | `site.baseURL` | Yes |
| `description` | `site.description` or `site.title` | Yes |
| `lastBuildDate` | Build timestamp | Yes |

### Item Elements

| Element | Source | Format |
|---------|--------|--------|
| `title` | Post title | Plain text |
| `link` | Post URL | Absolute URL |
| `description` | Post content | CDATA-wrapped HTML |
| `pubDate` | Post date | RFC 1123Z |

### Date Formats

**lastBuildDate**: RFC 1123Z
```
Mon, 15 Jan 2024 00:00:00 +0000
```

**pubDate**: RFC 1123Z
```
Mon, 15 Jan 2024 00:00:00 +0000
```

### Post Limit

Maximum 20 most recent posts included in feed.

### Description Content

Full HTML content of post wrapped in CDATA:

```xml
<description><![CDATA[<p>Post content here...</p>]]></description>
```

## robots.txt

### Location

`{output}/robots.txt`

### Content

```
User-agent: *
Allow: /
Sitemap: {BaseURL}/sitemap.xml
```

### Directives

| Directive | Value | Meaning |
|-----------|-------|---------|
| User-agent | `*` | Applies to all crawlers |
| Allow | `/` | Allow crawling all paths |
| Sitemap | URL | Location of sitemap |

## Favicon

### Configuration

```yaml
site:
  favicon: "/favicon.svg"
```

### HTML Output

```html
<link rel="icon" href="/favicon.svg" type="image/svg+xml">
```

### MIME Type Detection

| Extension | MIME Type |
|-----------|-----------|
| `.svg` | `image/svg+xml` |
| `.png` | `image/png` |
| `.ico` | `image/x-icon` |
| `.gif` | `image/gif` |
| `.jpg`, `.jpeg` | `image/jpeg` |

## SEO Best Practices

### URL Structure

- Clean URLs without extensions
- Lowercase, hyphen-separated slugs
- Trailing slashes for consistency
- Descriptive but concise

**Good**: `/blog/my-first-post/`
**Avoid**: `/blog/my_first_post.html`

### Content Guidelines

- Use descriptive `title` in frontmatter
- Provide meaningful `summary` for important pages
- Keep titles under 60 characters
- Keep descriptions under 160 characters

### Images

- Always include `alt` text in Markdown
- Use descriptive filenames
- Consider image optimization (external to SSG)

### Performance

- Single stylesheet
- No JavaScript by default
- Clean HTML output
- Small page sizes

## Testing SEO

### Tools

- Google Search Console
- Google Rich Results Test
- Facebook Sharing Debugger
- Twitter Card Validator
- Lighthouse SEO audit

### Validation

1. **Sitemap**: Submit to search consoles
2. **RSS**: Validate with feed validators
3. **Structured Data**: Test with Google's Rich Results Test
4. **Social**: Preview with platform debuggers
