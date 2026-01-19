# Development Server

> Version: 1.0.0
> Last Updated: 2026-01-19

## Overview

The development server provides a local HTTP server for previewing the generated site during development.

## Purpose

- Preview site locally before deployment
- Test clean URLs without deployment
- Rapid iteration during development

## Command

```bash
ssg serve [flags]
```

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--config` | `-c` | string | `ssg.yaml` | Configuration file path |
| `--port` | `-p` | int | `8080` | Port to listen on |
| `--dir` | `-d` | string | (from config) | Directory to serve |
| `--build` | `-b` | bool | `false` | Build before serving |

## Behavior

### Static File Serving

The server serves static files from the specified directory (default: config's `build.output`).

**Request → Response Mapping**:

| Request | Served File |
|---------|-------------|
| `/` | `/index.html` |
| `/about/` | `/about/index.html` |
| `/blog/my-post/` | `/blog/my-post/index.html` |
| `/style.css` | `/style.css` |
| `/feed/` | `/feed/index.xml` |

### Clean URL Handling

When a request is made for a directory path:

1. Check if directory contains `index.html`
2. If yes, serve `index.html`
3. If no, check for `index.xml`
4. If yes, serve `index.xml` with XML content type
5. Otherwise, return 404

### Content Types

File extensions are mapped to MIME types:

| Extension | Content-Type |
|-----------|--------------|
| `.html` | `text/html` |
| `.css` | `text/css` |
| `.js` | `application/javascript` |
| `.xml` | `application/xml` |
| `.json` | `application/json` |
| `.svg` | `image/svg+xml` |
| `.png` | `image/png` |
| `.jpg`, `.jpeg` | `image/jpeg` |
| `.gif` | `image/gif` |
| `.ico` | `image/x-icon` |
| `.woff` | `font/woff` |
| `.woff2` | `font/woff2` |

### Port Configuration

Default port: `8080`

If the port is in use, the server fails with an error.

```bash
# Custom port
ssg serve --port 3000
```

### Pre-Build Option

The `--build` flag triggers a site build before starting the server:

```bash
ssg serve --build
```

**Process**:
1. Execute build pipeline
2. If build fails, exit with error
3. If build succeeds, start server

## Security

### Path Traversal Prevention

The server rejects requests containing:
- `..` (parent directory traversal)
- Absolute paths
- Paths outside the served directory

**Rejected**:
```
/../etc/passwd
/../../secret.txt
```

### Request Timeouts

| Timeout | Value | Purpose |
|---------|-------|---------|
| Read header | 10 seconds | Prevent slowloris attacks |
| Idle | 120 seconds | Clean up idle connections |

### No Directory Listing

Directory requests without index files return 404, not directory listings.

## Graceful Shutdown

When receiving shutdown signal (SIGINT/SIGTERM):

1. Stop accepting new connections
2. Wait for active connections to complete (up to 120s)
3. Force close remaining connections
4. Exit cleanly

### Signals

| Signal | Action |
|--------|--------|
| `SIGINT` (Ctrl+C) | Graceful shutdown |
| `SIGTERM` | Graceful shutdown |

## Usage Examples

### Basic Usage

```bash
# Build site first
ssg build

# Start server
ssg serve
# Output: Serving site at http://localhost:8080
```

### Combined Build and Serve

```bash
ssg serve --build
```

### Custom Port

```bash
ssg serve --port 3000
```

### Serve Different Directory

```bash
ssg serve --dir dist
```

### With Custom Config

```bash
ssg serve --config ssg.dev.yaml --build
```

## Output

### Startup

```
Serving site at http://localhost:8080
Press Ctrl+C to stop
```

### Shutdown

```
Server stopped
```

### Errors

```
Error: port 8080 is already in use
Error: directory not found: public
```

## Request Logging

The server does not log individual requests in the current specification. Implementations may optionally add request logging for debugging.

**Optional logging format**:
```
GET /about/ 200 15ms
GET /style.css 200 2ms
GET /nonexistent/ 404 1ms
```

## Limitations

### No Live Reload

The server does not include live reload functionality. Changes require:
1. Rebuild site (`ssg build`)
2. Refresh browser manually

### No Watch Mode

The server does not watch for file changes. Future implementations may add:
- File watching
- Automatic rebuilds
- Browser refresh

### Single Site Only

The server serves one directory at a time. It does not support:
- Virtual hosts
- Multiple sites
- Proxying

### HTTP Only

The server uses HTTP only, not HTTPS. For local development, this is acceptable. For testing HTTPS-specific features, use a reverse proxy.

## Implementation Notes

### HTTP/1.1

The server should support HTTP/1.1 with:
- Keep-alive connections
- Proper Content-Length headers
- Correct status codes

### Error Responses

| Status | Condition |
|--------|-----------|
| 200 | File found and served |
| 404 | File not found |
| 405 | Method not allowed (non-GET) |
| 500 | Server error |

### Caching Headers

The development server should NOT set aggressive caching headers. This ensures changes are visible immediately on refresh.

**Recommended**:
```
Cache-Control: no-cache
```

## Testing the Server

### Manual Testing

```bash
# Start server
ssg serve --build

# Test in browser
open http://localhost:8080

# Test with curl
curl -I http://localhost:8080/
curl http://localhost:8080/about/
curl http://localhost:8080/feed/
```

### Verification Checklist

- [ ] Homepage loads at `/`
- [ ] Pages load at `/{slug}/`
- [ ] Posts load at `/blog/{slug}/`
- [ ] Blog listing loads at `/blog/`
- [ ] RSS feed loads at `/feed/`
- [ ] CSS loads correctly
- [ ] 404 page shows for missing paths
- [ ] Server shuts down cleanly with Ctrl+C
