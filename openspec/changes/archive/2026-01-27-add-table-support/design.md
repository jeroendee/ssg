## Context

The `ssg` static site generator uses goldmark v1.7.13 for markdown-to-HTML conversion. The parser is configured as a package-level singleton in `internal/parser/parser.go` with:
- `gmparser.WithAutoHeadingID()` - auto-generates heading IDs
- Custom `anchorHeadingRenderer` - wraps headings in anchor links

Goldmark has a built-in `extension.Table` that parses GFM-style table syntax. It's not currently enabled.

## Goals / Non-Goals

**Goals:**
- Enable markdown table rendering in generated HTML
- Maintain existing parser behavior for all other content

**Non-Goals:**
- Enhanced table CSS styling (existing `table { width: 100% }` is sufficient)
- Custom table rendering behavior
- Table of contents generation from tables

## Decisions

### Decision 1: Use goldmark's built-in extension.Table

**Choice**: Add `goldmark.WithExtensions(extension.Table)` to the parser initialization.

**Rationale**:
- Built-in, well-tested, no additional dependencies
- GFM-compatible syntax (pipe-based tables)
- Renders standard HTML table elements that work with existing CSS

**Alternatives considered**:
- Custom table parser: Unnecessary complexity for standard table support
- Different markdown library: goldmark is already in use and well-suited

### Decision 2: No changes to wordcount package

The `internal/wordcount/wordcount.go` creates its own fresh goldmark instance without extensions. Table content will still be counted correctly since it strips HTML tags after conversion.

## Risks / Trade-offs

**[Minimal risk]** Table extension slightly increases parsing time → Negligible for typical content sizes.

**[Non-issue]** Existing content with pipe characters → GFM table syntax requires specific formatting (header row, separator row), so existing pipe usage won't accidentally become tables.
