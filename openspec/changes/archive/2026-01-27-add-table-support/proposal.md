## Why

Markdown tables are a common formatting need for documentation and blog posts. Currently, `ssg` uses goldmark for markdown parsing but doesn't enable the built-in table extension, so table syntax is rendered as plain text.

## What Changes

- Enable goldmark's `extension.Table` in the markdown parser configuration
- Tables in markdown content will render as HTML `<table>` elements

## Capabilities

### New Capabilities

- `markdown-tables`: Support for GFM-style markdown table syntax rendering to HTML

### Modified Capabilities

None - existing table CSS (`table { width: 100% }`) is sufficient.

## Impact

- `internal/parser/parser.go` - Add table extension to goldmark initialization
- Tests - Add test case verifying table markdown converts to HTML
