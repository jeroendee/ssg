## 1. Configuration

> **Skill**: `go-writer` (1.1-1.3), `go-tester` (1.4) — TDD workflow for config parsing
> **Sub-agent**: `go-engineer` — Coordinates go-writer and go-tester for the full TDD cycle

- [x] 1.1 Add Feed struct to yamlConfig in config/config.go with Pages field
- [x] 1.2 Add FeedPages []string to model.Config
- [x] 1.3 Map yaml feed.pages to model.Config.FeedPages in Load function
- [x] 1.4 Write tests for feed config parsing (specified, empty, omitted)

## 2. Feed Item Abstraction

> **Skill**: `go-writer` — Interface definition and implementation following Google Go Style Guide
> **Sub-agent**: `go-engineer` — Validates interface design with go-analyzer

- [x] 2.1 Define FeedItem interface in model package (FeedTitle, FeedLink, FeedContent, FeedDate, FeedGUID)
- [x] 2.2 Implement FeedItem interface on Post type
- [x] 2.3 Create DateSection type for page date sections
- [x] 2.4 Implement FeedItem interface on DateSection type

## 3. Date Section Parsing

> **Skill**: `go-tester` (3.5), `go-writer` (3.1-3.4) — TDD workflow: write tests first, then implementation
> **Sub-agent**: `go-engineer` — Coordinates go-writer and go-tester for the full TDD cycle

- [x] 3.1 Create function to extract date sections from page content (regex for `{#YYYY-MM-DD}` anchors)
- [x] 3.2 Parse date from anchor ID to time.Time
- [x] 3.3 Extract content between date headers
- [x] 3.4 Generate title in format "{PageTitle} - {Month} {Day}, {Year}"
- [x] 3.5 Write tests for date section extraction (single date, multiple dates, no dates)

## 4. Builder Integration

> **Skill**: `go-writer` — Builder modifications following existing patterns
> **Sub-agent**: `go-engineer` — Validates with go-diagnostics after changes

- [x] 4.1 Load feed pages from config in builder
- [x] 4.2 Find matching Page objects for configured feed page paths
- [x] 4.3 Extract DateSection items from each feed page
- [x] 4.4 Merge post FeedItems with DateSection FeedItems
- [x] 4.5 Sort combined items by date descending
- [x] 4.6 Apply 20-item limit to combined feed

## 5. Renderer Updates

> **Skill**: `go-writer` (5.1, 5.2), `go-tester` (5.3) — Update renderer to use FeedItem interface
> **Sub-agent**: `go-engineer` — Coordinates refactoring with test updates

- [x] 5.1 Update RenderFeed to accept []FeedItem instead of []Post
- [x] 5.2 Generate entries using FeedItem interface methods
- [x] 5.3 Update tests for mixed feed items

## 6. Integration Testing

> **Skill**: `go-tester` — Integration tests for end-to-end feed generation
> **Sub-agent**: `go-engineer` — Runs diagnostics and validates build

- [x] 6.1 Test feed generation with posts only (existing behavior)
- [x] 6.2 Test feed generation with page sections only
- [x] 6.3 Test feed generation with mixed posts and page sections
- [x] 6.4 Test feed with interleaved dates (verify sort order)
- [x] 6.5 Test 20-item limit with mixed sources
