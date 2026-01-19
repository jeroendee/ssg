# CLI

> Version: 1.0.0
> Last Updated: 2026-01-19

## Overview

The SSG provides a command-line interface with three main commands: `build`, `serve`, and `version`.

## Installation

The SSG is distributed as a single binary. After building or downloading:

```bash
# Add to PATH or use directly
./ssg --help
```

## Commands

### build

Generate the static website.

```
ssg build [flags]
```

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--config` | `-c` | string | `ssg.yaml` | Configuration file path |
| `--output` | `-o` | string | (from config) | Output directory |
| `--content` | | string | (from config) | Content directory |
| `--assets` | | string | (from config) | Assets directory |

#### Examples

```bash
# Build with defaults
ssg build

# Build with custom config
ssg build --config production.yaml

# Override output directory
ssg build --output dist

# Multiple overrides
ssg build --content src/content --output build
```

#### Process

1. Load configuration
2. Validate paths and content
3. Parse all content files
4. Render HTML pages
5. Copy assets
6. Generate SEO files
7. Report success or failure

#### Output

On success:
```
Site built successfully to public/
```

On failure:
```
Error: {error message}
```

### serve

Run a development server.

```
ssg serve [flags]
```

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--config` | `-c` | string | `ssg.yaml` | Configuration file path |
| `--port` | `-p` | int | `8080` | Port to listen on |
| `--dir` | `-d` | string | (from config) | Directory to serve |
| `--build` | `-b` | bool | `false` | Build site before serving |

#### Examples

```bash
# Serve with defaults (port 8080)
ssg serve

# Serve on custom port
ssg serve --port 3000

# Build then serve
ssg serve --build

# Serve specific directory
ssg serve --dir dist
```

#### Behavior

1. Optionally build site first (`--build`)
2. Start HTTP server on specified port
3. Serve static files from directory
4. Handle clean URLs (serve index.html for directories)
5. Continue until interrupted (Ctrl+C)

#### Output

On start:
```
Serving site at http://localhost:8080
Press Ctrl+C to stop
```

On shutdown:
```
Server stopped
```

### version

Display version information.

```
ssg version
```

#### Output

```
ssg version {VERSION}
```

Where `{VERSION}` is:
- Git commit SHA (when built with version injection via build system)
- `dev` (when built without version injection)

#### Examples

```
ssg version f4a9750
ssg version dev
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--help` | Show help for command |

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Error (build failure, invalid config, etc.) |

## Error Messages

### Configuration Errors

```
Error: config file not found: ssg.yaml
Error: invalid config syntax: {yaml error}
Error: site.title is required
Error: site.baseURL is required
```

### Content Errors

```
Error: home.md not found in content directory
Error: invalid post filename format: {filename}
Error: invalid frontmatter in {file}: {error}
Error: missing required field 'title' in {file}
```

### Path Errors

```
Error: invalid output directory: cannot be root
Error: invalid output directory: cannot be home directory
Error: referenced asset not found: {path}
```

### Server Errors

```
Error: port {port} is already in use
Error: directory not found: {path}
```

## Usage Patterns

### Development Workflow

```bash
# Initial build
ssg build

# Start dev server
ssg serve

# Or combined
ssg serve --build
```

### Production Build

```bash
# Build with production config
ssg build --config ssg.production.yaml

# Build to specific directory
ssg build --output dist
```

### CI/CD Integration

```bash
#!/bin/bash
set -e

# Build site
ssg build --output public

# Deploy (example)
rsync -avz public/ server:/var/www/html/
```

## Configuration File Discovery

The CLI looks for configuration in this order:

1. Path specified by `--config` flag
2. `ssg.yaml` in current directory

If not found, build fails with error.

## Signal Handling

### serve Command

| Signal | Action |
|--------|--------|
| `SIGINT` (Ctrl+C) | Graceful shutdown |
| `SIGTERM` | Graceful shutdown |

Graceful shutdown:
1. Stop accepting new connections
2. Wait for active requests (timeout: 120s)
3. Exit cleanly

## Help Output

```
$ ssg --help
SSG - A minimal static site generator

Usage:
  ssg [command]

Available Commands:
  build       Build the static site
  serve       Start development server
  version     Print version information
  help        Help about any command

Flags:
  -h, --help   help for ssg

Use "ssg [command] --help" for more information about a command.
```

```
$ ssg build --help
Build the static site

Usage:
  ssg build [flags]

Flags:
  -c, --config string    Config file (default "ssg.yaml")
  -o, --output string    Output directory
      --content string   Content directory
      --assets string    Assets directory
  -h, --help            help for build
```

```
$ ssg serve --help
Start development server

Usage:
  ssg serve [flags]

Flags:
  -c, --config string   Config file (default "ssg.yaml")
  -p, --port int        Port to listen on (default 8080)
  -d, --dir string      Directory to serve
  -b, --build           Build site before serving
  -h, --help           help for serve
```
