## Why

The site has a dedicated `/contact/` page that adds navigation clutter while containing minimal content. Contact information belongs in the footer where it's accessible from every page. This consolidation simplifies site structure and improves discoverability.

## What Changes

- Add support for rendering markdown content in the site footer
- Footer loads content from `content/_footer.md` (markdown partial)
- Add horizontal rule (`<hr>`) at top of footer for visual separation
- Keep existing "Made with ❤️ in Amsterdam" tagline with version
- Remove dedicated `/contact/` page
- Remove Contact from navigation

## Capabilities

### New Capabilities
- `footer-content`: Support for optional markdown content in the site footer, loaded from `_footer.md` partial file

### Modified Capabilities
- `html-rendering`: Footer template needs to render optional markdown content with `<hr>` separator
- `content-scanning`: Builder must detect and parse `_footer.md` partial (already skips `_` prefixed files from pages)

## Impact

- **Model**: `Site` struct gains `FooterContent string` field
- **Builder**: `ScanContent()` loads `_footer.md` if present
- **Templates**: `_footer.html` renders content with `<hr>`
- **Config**: Navigation loses Contact entry
- **Content**: `contact.md` deleted, `_footer.md` created
