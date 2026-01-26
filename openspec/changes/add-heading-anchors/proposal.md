# Proposal: Add Heading Anchor Links

## Summary

Enable automatic ID generation for all heading levels (h1-h6) so users can link directly to specific sections of pages and posts.

## Motivation

Currently, headings in rendered content have no ID attributes. Users cannot share direct links to specific sections of a page (e.g., `/about/#features`). This is a standard feature in most static site generators and documentation tools.

## Scope

- **Capability affected**: content-parsing (Transform-Markup requirement)
- **Files affected**: internal/parser/parser.go
- **Risk**: Low - additive change, backward compatible

## Approach

Configure Goldmark (already in use) with `parser.WithAutoHeadingID()` to automatically generate slug-based IDs for all heading elements.

**ID Generation Rules** (Goldmark defaults):
- Lowercase the heading text
- Replace spaces with hyphens
- Remove special characters
- Example: `## Hello World!` becomes `<h2 id="hello-world">Hello World!</h2>`

## User Impact

- **Before**: `<h2>Features</h2>`
- **After**: `<h2 id="features">Features</h2>`
- **Benefit**: Users can share links like `/about/#features` to jump directly to that section

## Out of Scope

- Table of contents generation
- Custom ID specification in frontmatter
- Anchor link icons (visible link symbols next to headings)
