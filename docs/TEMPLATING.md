# Templating

> Version: 1.0.0
> Last Updated: 2026-01-19

## Overview

The SSG uses a template system to render HTML pages. Templates define the structure and layout of pages, while content and configuration provide the data.

## Template Types

### Core Templates

| Template | Purpose | Output |
|----------|---------|--------|
| `base.html` | Static pages layout | `/{slug}/index.html` |
| `blog_post.html` | Individual blog post | `/blog/{slug}/index.html` |
| `blog_list.html` | Blog listing page | `/blog/index.html` |
| `404.html` | Error page | `/404.html` |

### Partial Templates

| Partial | Purpose |
|---------|---------|
| `_head.html` | HTML head section (meta tags, styles) |
| `_header.html` | Site header and navigation |
| `_footer.html` | Site footer |

## Template Composition

Templates include partials to form complete pages:

```
┌─────────────────────────────────────┐
│           _head.html                │  ← Meta tags, styles
├─────────────────────────────────────┤
│          _header.html               │  ← Navigation
├─────────────────────────────────────┤
│                                     │
│      Main Template Content          │  ← Page/post specific
│   (base.html, blog_post.html, etc)  │
│                                     │
├─────────────────────────────────────┤
│          _footer.html               │  ← Footer, analytics
└─────────────────────────────────────┘
```

## Template Variables

### Global Variables

Available in all templates:

| Variable | Type | Description |
|----------|------|-------------|
| `Site.Title` | String | Site name |
| `Site.BaseURL` | String | Root URL |
| `Site.Description` | String | Site description |
| `Site.Author` | String | Author name |
| `Site.Logo` | String | Logo path |
| `Site.Favicon` | String | Favicon path |
| `Site.Navigation` | List | Navigation items |
| `Site.Analytics.GoatCounter` | String | Analytics code |
| `PageTitle` | String | Current page/post title |
| `CanonicalURL` | String | Absolute URL of current page |
| `Summary` | String | Page/post summary |
| `IsPost` | Boolean | True if rendering a blog post |
| `OGImage` | String | Open Graph image URL |
| `Version` | String | Build version (git SHA) |
| `PageType` | String | Type identifier |

### Navigation Item

| Field | Type | Description |
|-------|------|-------------|
| `Title` | String | Display text |
| `URL` | String | Link URL |

### Page-Specific Variables

**For `base.html` (pages)**:

| Variable | Type | Description |
|----------|------|-------------|
| `Page.Title` | String | Page title |
| `Page.Content` | HTML | Rendered page content |

**For `blog_post.html`**:

| Variable | Type | Description |
|----------|------|-------------|
| `Post.Title` | String | Post title |
| `Post.DateFormatted` | String | Date in `YYYY-MM-DD` format |
| `Post.Content` | HTML | Rendered post content |
| `Post.WordCount` | Integer | Word count |

**For `blog_list.html`**:

| Variable | Type | Description |
|----------|------|-------------|
| `Posts` | List | All posts |
| `Posts[].Title` | String | Post title |
| `Posts[].Slug` | String | URL slug |
| `Posts[].DateFormatted` | String | Date string |
| `Posts[].WordCount` | Integer | Word count |

## PageType Values

| Value | Context |
|-------|---------|
| `"page"` | Static pages |
| `"blog_post"` | Blog post |
| `"blog_list"` | Blog listing |
| `"404"` | Error page |

## Template Syntax

Templates use a standard templating syntax with the following constructs:

### Variable Output

```html
{{ .Variable }}
{{ .Site.Title }}
{{ .Post.Title }}
```

### Conditionals

```html
{{ if .Site.Logo }}
  <img src="{{ .Site.Logo }}" alt="Logo">
{{ end }}

{{ if .IsPost }}
  <time>{{ .Post.DateFormatted }}</time>
{{ end }}
```

### Loops

```html
{{ range .Site.Navigation }}
  <a href="{{ .URL }}">{{ .Title }}</a>
{{ end }}

{{ range .Posts }}
  <article>
    <h2>{{ .Title }}</h2>
    <time>{{ .DateFormatted }}</time>
  </article>
{{ end }}
```

### Raw HTML Output

Content fields containing HTML should be output without escaping:

```html
{{ .Page.Content }}
{{ .Post.Content }}
```

## Default Stylesheet

The SSG includes a default stylesheet with:

- **Color Scheme**: Solarized light and dark modes
- **Typography**: Clean, readable fonts
- **Layout**: Centered content (max-width: 720px)
- **Responsive**: Mobile-friendly design
- **Dark Mode**: Automatic via `prefers-color-scheme`

### Color Variables

| Variable | Light | Dark |
|----------|-------|------|
| Background | `#fdf6e3` | `#002b36` |
| Text | `#657b83` | `#839496` |
| Headings | `#073642` | `#93a1a1` |
| Links | `#268bd2` | `#268bd2` |
| Code background | `#eee8d5` | `#073642` |

## Custom Stylesheet Override

To override the default stylesheet:

1. Create `assets/style.css`
2. Build site
3. Your stylesheet replaces the default entirely

**Note**: Custom stylesheets completely replace the default. Copy and modify the default if you want to extend it.

## Template Files Reference

### _head.html

Generates the `<head>` section:

```html
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{ .PageTitle }} - {{ .Site.Title }}</title>

  <!-- SEO Meta Tags -->
  <meta name="description" content="{{ .Summary }}">
  <link rel="canonical" href="{{ .CanonicalURL }}">

  <!-- Open Graph -->
  <meta property="og:title" content="{{ .PageTitle }}">
  <meta property="og:description" content="{{ .Summary }}">
  <meta property="og:url" content="{{ .CanonicalURL }}">
  <meta property="og:type" content="...">
  <meta property="og:site_name" content="{{ .Site.Title }}">

  <!-- Twitter Card -->
  <meta name="twitter:card" content="summary">
  <meta name="twitter:title" content="{{ .PageTitle }}">
  <meta name="twitter:description" content="{{ .Summary }}">

  <!-- Favicon -->
  {{ if .Site.Favicon }}
  <link rel="icon" href="{{ .Site.Favicon }}" type="...">
  {{ end }}

  <!-- Stylesheet -->
  <link rel="stylesheet" href="/style.css">

  <!-- JSON-LD Structured Data -->
  <script type="application/ld+json">...</script>
</head>
```

### _header.html

Generates site header and navigation:

```html
<header>
  <nav>
    <a href="/" class="site-title">{{ .Site.Title }}</a>
    <ul>
      {{ range .Site.Navigation }}
      <li><a href="{{ .URL }}">{{ .Title }}</a></li>
      {{ end }}
    </ul>
  </nav>
</header>
```

### _footer.html

Generates site footer:

```html
<footer>
  <p>Made with ❤️ ({{ .Version }})</p>

  {{ if .Site.Analytics.GoatCounter }}
  <script data-goatcounter="https://{{ .Site.Analytics.GoatCounter }}.goatcounter.com/count"
          async src="//gc.zgo.at/count.js"></script>
  {{ end }}
</footer>
```

### base.html

Layout for static pages:

```html
<!DOCTYPE html>
<html lang="en">
{{ template "_head.html" . }}
<body>
  {{ template "_header.html" . }}
  <main>
    <h1>{{ .Page.Title }}</h1>
    {{ .Page.Content }}
  </main>
  {{ template "_footer.html" . }}
</body>
</html>
```

### blog_post.html

Layout for blog posts:

```html
<!DOCTYPE html>
<html lang="en">
{{ template "_head.html" . }}
<body>
  {{ template "_header.html" . }}
  <main>
    <article>
      <header>
        <h1>{{ .Post.Title }}</h1>
        <time datetime="{{ .Post.DateFormatted }}">{{ .Post.DateFormatted }}</time>
        <span>{{ .Post.WordCount }} words</span>
      </header>
      {{ .Post.Content }}
    </article>
  </main>
  {{ template "_footer.html" . }}
</body>
</html>
```

### blog_list.html

Layout for blog listing:

```html
<!DOCTYPE html>
<html lang="en">
{{ template "_head.html" . }}
<body>
  {{ template "_header.html" . }}
  <main>
    <h1>Blog</h1>
    <ul>
      {{ range .Posts }}
      <li>
        <a href="/blog/{{ .Slug }}/">{{ .Title }}</a>
        <time>{{ .DateFormatted }}</time>
        <span>{{ .WordCount }} words</span>
      </li>
      {{ end }}
    </ul>
  </main>
  {{ template "_footer.html" . }}
</body>
</html>
```

### 404.html

Error page layout:

```html
<!DOCTYPE html>
<html lang="en">
{{ template "_head.html" . }}
<body>
  {{ template "_header.html" . }}
  <main>
    <h1>Page Not Found</h1>
    <p>The page you're looking for doesn't exist.</p>
    <p><a href="/">Go home</a></p>
  </main>
  {{ template "_footer.html" . }}
</body>
</html>
```

## Custom Templates (Future)

To support custom templates in future implementations:

1. Check for user-provided templates in a templates directory
2. Fall back to embedded defaults if not found
3. Allow partial overrides (e.g., only override `_footer.html`)
