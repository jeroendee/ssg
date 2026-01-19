# ADR-0004: Date from Filename

> Status: Accepted
> Date: 2026-01-19

## Context

Blog posts need publication dates for sorting and display. Options for specifying dates:

1. **Filename**: `2024-01-15-my-post.md`
2. **Frontmatter only**: `date: 2024-01-15`
3. **File modification time**: Use filesystem mtime
4. **Directory structure**: `2024/01/15/my-post.md`

Considerations:
- Visibility (date apparent in file listing)
- Chronological sorting in filesystem
- Rename-ability (can date be changed?)
- Git-friendliness (mtime is unreliable with git)

## Decision

Use **filename pattern** `YYYY-MM-DD-slug.md` as the primary date source, with frontmatter `date` field as override.

**Format**: `YYYY-MM-DD-slug.md`
- `YYYY`: 4-digit year
- `MM`: 2-digit month (01-12)
- `DD`: 2-digit day (01-31)
- `-slug`: URL slug (rest of filename)

**Examples**:
```
2024-01-15-first-post.md
2024-12-25-christmas-update.md
2025-03-08-march-thoughts.md
```

**Override**:
```yaml
---
title: "My Post"
date: "2024-01-20"  # Overrides filename date
---
```

## Consequences

### Positive

- **Visible in file listing**: Date apparent without opening file
- **Natural sorting**: Files sort chronologically in filesystem
- **Git-friendly**: Date doesn't depend on filesystem metadata
- **Explicit**: No ambiguity about intended publication date
- **Override available**: Frontmatter allows corrections

### Negative

- **Filename coupling**: Changing date requires renaming file
- **Strict format**: Invalid patterns cause build failure
- **URL leakage**: Date visible in file, not in URL (which is fine)

### Neutral

- Jekyll convention (familiar to many users)
- Slug extracted separately from date

## Implementation Notes

### Parsing Algorithm

1. Check if filename matches `^(\d{4})-(\d{2})-(\d{2})-(.+)\.md$`
2. Extract year, month, day, slug
3. Parse into date object
4. Check frontmatter for `date` field
5. If frontmatter `date` exists, use it instead
6. Validate date is real (not Feb 30, etc.)

### Error Handling

- Non-matching filename in blog directory → Build error
- Invalid date (Feb 30) → Build error
- Future dates → Allowed (scheduled posts)

### Date Format Output

| Context | Format |
|---------|--------|
| Template display | `YYYY-MM-DD` |
| RSS pubDate | RFC 1123Z |
| Sitemap lastmod | `YYYY-MM-DD` |
| JSON-LD datePublished | `YYYY-MM-DD` |

## References

- [Jekyll post naming](https://jekyllrb.com/docs/posts/)
- ADR-0003-frontmatter-yaml.md (override mechanism)
