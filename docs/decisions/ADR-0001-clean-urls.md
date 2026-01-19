# ADR-0001: Clean URL Pattern

> Status: Accepted
> Date: 2026-01-19

## Context

Static site generators must decide how to map content files to output URLs. The main options are:

1. **Direct mapping**: `about.md` → `/about.html`
2. **Clean URLs**: `about.md` → `/about/index.html` (accessed as `/about/`)
3. **No extension, no slash**: `about.md` → `/about` (requires server config)

Considerations:
- URL aesthetics and memorability
- Server compatibility (works without special configuration)
- Trailing slash consistency
- Link portability across hosting platforms

## Decision

Use the **clean URL pattern**: Each page generates a directory with an `index.html` file.

**Mapping**:
| Input | Output | URL |
|-------|--------|-----|
| `about.md` | `about/index.html` | `/about/` |
| `contact.md` | `contact/index.html` | `/contact/` |
| `blog/2024-01-15-post.md` | `blog/post/index.html` | `/blog/post/` |

**Special case**: `home.md` → `index.html` (root `/`)

## Consequences

### Positive

- **Cleaner URLs**: No `.html` extension visible to users
- **Universal compatibility**: Works on any static file server without configuration
- **Consistent trailing slash**: All paths end with `/`
- **Future-proof**: URLs don't reveal implementation details

### Negative

- **More directories**: Each page creates a directory
- **Larger output**: Slightly more filesystem overhead
- **Trailing slash required**: Must ensure consistent linking with trailing slashes

### Neutral

- Server configuration optional (but can redirect non-slash to slash version)
- Some hosting platforms handle this automatically

## Implementation Notes

- Generate `{slug}/index.html` for all content except homepage
- Homepage generates `index.html` at root
- Internal links should always include trailing slash
- Sitemap URLs should include trailing slash

## References

- [Google on URL structure](https://developers.google.com/search/docs/advanced/guidelines/url-structure)
- [Netlify's pretty URLs](https://docs.netlify.com/routing/redirects/redirect-options/#trailing-slash)
