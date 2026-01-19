# ADR-0003: YAML Frontmatter

> Status: Accepted
> Date: 2026-01-19

## Context

Markdown files need metadata (title, date, summary, etc.). Common approaches:

1. **YAML frontmatter**: Metadata in YAML block at file start
2. **TOML frontmatter**: Similar, using TOML syntax
3. **JSON frontmatter**: Metadata as JSON object
4. **Filename only**: Extract all metadata from filename
5. **Separate files**: `.json` or `.yaml` companion files

Considerations:
- Human readability and writability
- Tooling support (editors, linters)
- Flexibility for future fields
- Industry standard practices

## Decision

Use **YAML frontmatter** delimited by `---` markers.

**Format**:
```markdown
---
title: "My Post Title"
summary: "A brief description"
date: "2024-01-15"
---

Content starts here...
```

## Consequences

### Positive

- **Industry standard**: Most SSGs use this format
- **Human readable**: YAML is easy to read and write
- **Editor support**: Many editors highlight frontmatter
- **Flexible**: Supports strings, lists, nested objects
- **Familiar**: Matches Hugo, Jekyll, Eleventy conventions

### Negative

- **YAML quirks**: Requires quoting some strings, indentation matters
- **No schema validation**: Invalid YAML fails at build time
- **Mixed format**: File contains both YAML and Markdown

### Neutral

- Requires YAML parser in addition to Markdown parser
- Standard `---` delimiters are unambiguous

## Implementation Notes

- Frontmatter must start on line 1 with `---\n`
- Closing `---\n` marks end of frontmatter
- Everything after is Markdown content
- Parse YAML into structured object
- Validate required fields (e.g., `title`)
- Handle missing optional fields gracefully

### Required Fields

| Field | Content Type | Description |
|-------|--------------|-------------|
| `title` | All | Page/post title |

### Optional Fields

| Field | Content Type | Description |
|-------|--------------|-------------|
| `summary` | Posts | Brief description |
| `date` | Posts | Date override (YYYY-MM-DD) |

## References

- [YAML 1.2 Specification](https://yaml.org/spec/1.2/spec.html)
- ADR-0004-date-from-filename.md (primary date source)
