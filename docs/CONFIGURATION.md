# Configuration

> Version: 1.0.0
> Last Updated: 2026-01-19

## Overview

Configuration is defined in a YAML file that specifies site metadata, build paths, and optional features.

## Configuration File

**Default Location**: `ssg.yaml` in project root

**Format**: YAML

## Schema

### Full Configuration Example

```yaml
site:
  title: "My Blog"
  baseURL: "https://example.com"
  description: "A personal blog about technology"
  author: "John Doe"
  logo: "/logo.svg"
  favicon: "/favicon.svg"

build:
  content: "content"
  output: "public"
  assets: "assets"

navigation:
  - title: "Home"
    url: "/"
  - title: "About"
    url: "/about/"
  - title: "Blog"
    url: "/blog/"

analytics:
  goatcounter: "mysite"
```

## Fields

### Site Section

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `title` | String | Yes | - | Site name displayed in header, titles, metadata |
| `baseURL` | String | Yes | - | Root URL where site is hosted (no trailing slash) |
| `description` | String | No | - | Site description for RSS and meta tags |
| `author` | String | No | - | Author name for JSON-LD and RSS |
| `logo` | String | No | - | Path to logo image (relative to site root) |
| `favicon` | String | No | - | Path to favicon (relative to site root) |

### Build Section

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `content` | String | No | `"content"` | Directory containing Markdown files |
| `output` | String | No | `"public"` | Directory for generated site |
| `assets` | String | No | `"assets"` | Directory containing static assets |

### Navigation Section

List of navigation menu items.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | String | Yes | Display text for link |
| `url` | String | Yes | URL path (relative or absolute) |

### Analytics Section

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `goatcounter` | String | No | GoatCounter site code |

## Required Fields

The following fields must be present and non-empty:

1. `site.title` — Site name
2. `site.baseURL` — Root URL

**Build will fail** if these fields are missing or empty.

## Default Values

| Field | Default |
|-------|---------|
| `build.content` | `"content"` |
| `build.output` | `"public"` |
| `build.assets` | `"assets"` |

## CLI Overrides

Command-line flags override configuration file values.

### Build Command Flags

| Flag | Short | Overrides |
|------|-------|-----------|
| `--config` | `-c` | Configuration file path |
| `--output` | `-o` | `build.output` |
| `--content` | - | `build.content` |
| `--assets` | - | `build.assets` |

### Serve Command Flags

| Flag | Short | Overrides |
|------|-------|-----------|
| `--config` | `-c` | Configuration file path |
| `--dir` | `-d` | Directory to serve (not `build.output`) |
| `--port` | `-p` | Server port (not in config) |

### Override Precedence

```
CLI flag > Config file value > Default value
```

**Example**:
```bash
# Config has output: "dist"
# CLI overrides to "build"
ssg build --output build
```

## Validation Rules

### Site Title
- Must be non-empty string
- No length limit
- Whitespace-only values are invalid

### Base URL
- Must be non-empty string
- Should be valid URL format
- No trailing slash recommended
- Used for absolute URLs in sitemap, RSS, meta tags

**Valid**:
```yaml
baseURL: "https://example.com"
baseURL: "https://blog.example.com"
```

**Invalid**:
```yaml
baseURL: ""
baseURL: "https://example.com/"  # Trailing slash not recommended
```

### Content Directory
- Must be relative path
- Must exist at build time
- Must contain `home.md`

### Output Directory
- Must be relative path
- Cannot be root (`/`)
- Cannot be home directory (`~`)
- Cannot be project root
- Cannot contain `..`

### Assets Directory
- Must be relative path
- Optional (build succeeds without it)

### Navigation Items
- Each item must have `title` and `url`
- Empty navigation list is valid
- URLs can be relative or absolute

## Favicon MIME Type

The favicon MIME type is derived from file extension:

| Extension | MIME Type |
|-----------|-----------|
| `.svg` | `image/svg+xml` |
| `.png` | `image/png` |
| `.ico` | `image/x-icon` |
| `.gif` | `image/gif` |
| `.jpg`, `.jpeg` | `image/jpeg` |

## Minimal Configuration

```yaml
site:
  title: "My Site"
  baseURL: "https://example.com"
```

## Environment-Specific Configuration

For different environments, use separate config files:

```bash
ssg build --config ssg.production.yaml
ssg build --config ssg.staging.yaml
```

Or override specific values:

```bash
ssg build --output dist
```

## Configuration Loading Process

1. Look for config file at specified path (default: `ssg.yaml`)
2. Parse YAML content
3. Validate syntax
4. Apply default values for missing optional fields
5. Merge CLI overrides
6. Validate required fields
7. Return configuration object or error

## Error Messages

| Condition | Error |
|-----------|-------|
| File not found | "config file not found: {path}" |
| Invalid YAML | "invalid config syntax: {details}" |
| Missing title | "site.title is required" |
| Missing baseURL | "site.baseURL is required" |
| Invalid output path | "invalid output directory: {reason}" |
