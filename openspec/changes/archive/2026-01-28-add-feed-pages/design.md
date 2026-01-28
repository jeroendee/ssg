## Context

Currently, the RSS feed only contains blog posts. Pages with frequently updated content (like `/moments/` or `/now/`) cannot notify subscribers of updates. Users who curate daily content want one notification per active day when they add new date sections.

The page content uses date-anchored headers:
```markdown
## January 27, 2026 {#2026-01-27}
Content for that day...

## January 26, 2026 {#2026-01-26}
Earlier content...
```

## Goals / Non-Goals

**Goals:**
- Allow configured pages to contribute date-based entries to the RSS feed
- One feed entry per unique date section on the page
- Entries link directly to the date anchor on the page
- Mixed feed (posts + page sections) sorted by date
- Preserve existing blog post feed behavior

**Non-Goals:**
- Real-time notifications (one per update) - daily granularity is sufficient
- Automatic detection of "moments-like" pages - explicit config required
- Changes to the date anchor format - use existing `{#YYYY-MM-DD}` convention

## Decisions

### 1. Configuration in `feed.pages` array

Add a new `feed` section to `ssg.yaml`:
```yaml
feed:
  pages:
    - /now/
    - /moments/
```

**Rationale**: Config describes site behavior, frontmatter describes content. Feed generation is behavior.

**Alternative considered**: Frontmatter `feed: daily` flag. Rejected because it conflates content metadata with delivery mechanism.

### 2. Introduce `FeedItem` interface

Create a common type that both posts and page sections can satisfy:
```go
type FeedItem interface {
    FeedTitle() string
    FeedLink() string
    FeedContent() string
    FeedDate() time.Time
    FeedGUID() string
}
```

**Rationale**: Clean abstraction allows the renderer to handle mixed item types without special-casing. Posts and page sections implement the same interface.

**Alternative considered**: Create synthetic `Post` objects from page sections. Rejected because it overloads the Post type with page-specific concerns.

### 3. Date section parsing via regex

Parse page content for date headers matching:
```
## <Month> <Day>, <Year> {#YYYY-MM-DD}
```

Extract:
- The anchor ID as the date identifier
- Content between this header and the next (or end of file)

**Rationale**: Consistent with existing date anchor convention already used for navigation.

### 4. GUID uses page URL with date anchor

Format: `{baseURL}{pagePath}#{date}`
Example: `https://aishepherd.nl/moments/#2026-01-27`

**Rationale**:
- Unique per date (avoids re-notification for same day updates)
- Deep-links directly to the relevant section
- Follows RSS guid best practices (unique, permanent)

### 5. Feed entry content is full day's content

Include the entire content block for that date section in the feed entry.

**Rationale**: Allows feed readers to show the full content inline, matching blog post behavior.

## Risks / Trade-offs

**Risk**: Date parsing regex may not handle all date formats
→ **Mitigation**: Only support the existing `{#YYYY-MM-DD}` anchor format. Document this requirement.

**Risk**: Large pages with many dates could bloat the feed
→ **Mitigation**: Apply the existing 20-item limit to the combined feed (posts + page sections). Most recent items win.

**Risk**: Pages without date anchors silently produce no entries
→ **Mitigation**: This is acceptable - config specifies intent, page content determines output. Could add warning in verbose mode.
