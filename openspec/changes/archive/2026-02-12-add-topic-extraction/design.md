## Context

ssg is a static site generator that builds date-anchored pages like `/moments/`. These pages grow over time with daily entries. Currently readers rely on `ctrl+F` or scrolling to discover content. The `feed.pages` config pattern already exists for opting pages into RSS feed generation — topic extraction follows the same approach.

## Goals / Non-Goals

**Goals:**
- Auto-extract recurring subject words from page markdown at build time
- Display top 20 topics with frequency counts, ordered by frequency descending
- Config-driven opt-in via `topics.pages` (mirrors `feed.pages` pattern)
- Clean separation: new `internal/topics/` package owns all extraction logic

**Non-Goals:**
- Clickable/interactive topics (plain text only)
- NLP or stemming (simple word tokenization is sufficient)
- Topic extraction for blog posts (pages only)
- Custom stop-word lists via config

## Decisions

### 1. New `internal/topics/` package for extraction logic

All word tokenization, stop-word filtering, and frequency counting lives in a dedicated package. This keeps the parser and builder clean.

**Alternative considered**: Adding extraction to `internal/parser/`. Rejected because topic extraction is a distinct concern from markdown parsing — it operates on raw text, not on parsed AST nodes.

### 2. Topic type as a simple struct

```go
type Topic struct {
    Word  string
    Count int
}
```

`Extract(markdown string) []Topic` returns up to 20 topics sorted by count descending. The function accepts raw markdown (after frontmatter removal) so it can operate on the text the reader actually sees.

### 3. Config via `topics.pages` (mirrors `feed.pages`)

```yaml
topics:
  pages:
    - /moments/
```

Parsed into `model.Config.TopicPages []string`. The builder checks each page's path against this list during `ScanContent` and runs extraction for matches.

**Alternative considered**: Frontmatter `topics: true`. Rejected because Dee prefers config-level control, consistent with how `feed.pages` works.

### 4. Stop-word filtering strategy

A hardcoded list of ~150 English stop words (articles, prepositions, pronouns, common verbs). Words must be ≥ 3 characters and appear ≥ 2 times. Words are lowercased before counting. Markdown syntax artifacts (URLs, image refs, code fences) are stripped before tokenization.

### 5. Tokenization approach

1. Strip markdown links `[text](url)` — keep text, discard URL
2. Strip image references `![alt](path)`
3. Strip inline code backticks
4. Strip HTML entities (`&amp;`, `&quot;`, etc.)
5. Split on non-alphanumeric boundaries (keep hyphens within words for terms like `pre-push`)
6. Lowercase everything
7. Filter stop words and short words (< 3 chars)
8. Count and sort

### 6. Template rendering

Add `Topics []Topic` to `model.Page`. In `base.html`, render a `<div class="topics">` between the date-nav and content div, only when `Topics` is non-empty. Format: `word (count), word (count), ...`

### 7. Minimal CSS

A single `.topics` class in the default stylesheet — small font, muted color, consistent with the existing date-nav styling. No special layout needed; the comma-separated text wraps naturally to ~2 lines.

## Risks / Trade-offs

- **[Stop-word list completeness]** → Some irrelevant words may slip through. Mitigation: start with a well-known list, easy to extend later.
- **[Markdown stripping imperfection]** → Regex-based stripping may miss edge cases. Mitigation: good-enough extraction — a stray URL fragment in the topic list is low-impact and easily caught in review.
- **[Performance on large pages]** → Topic extraction adds string processing per configured page. Mitigation: negligible for realistic page sizes (even 1000 entries would be milliseconds).
