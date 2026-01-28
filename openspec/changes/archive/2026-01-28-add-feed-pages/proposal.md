## Why

Pages with frequently updated content (like `/moments/` or `/now/` pages) cannot notify subscribers of updates via RSS. Currently, only blog posts appear in the feed. Users who curate daily content want followers to be notified when new daily entries are added, without creating separate blog posts for each update.

## What Changes

- Add `feed.pages` configuration option to `ssg.yaml` to specify pages that should generate feed entries
- Parse configured pages for date-anchored sections (e.g., `## January 27, 2026 {#2026-01-27}`)
- Generate one feed entry per date section, with guid pointing to the date anchor
- Merge page-based entries with blog post entries, sorted by date
- Existing blog post feed behavior remains unchanged

## Capabilities

### New Capabilities
- `feed-pages`: Configuration and processing of pages that contribute date-based entries to the RSS feed

### Modified Capabilities
- `site-configuration`: Add `feed.pages` configuration option
- `feed-generation`: Accept mixed feed items (posts and page sections), generate entries for both

## Impact

- `internal/config/config.go`: Add Feed section to yaml struct
- `internal/model/model.go`: Add FeedPages to Config, possibly new FeedItem type
- `internal/builder/builder.go`: Parse feed pages for date sections
- `internal/renderer/renderer.go`: Handle mixed feed item types
- Config file schema: New `feed.pages` array under `feed` section
