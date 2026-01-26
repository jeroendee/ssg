# Tasks: Add Heading Anchor Links

## Implementation Order

- [x] **Write failing tests for heading ID generation**
   - Test all heading levels (h1-h6) produce IDs
   - Test ID format (lowercase, hyphenated)
   - Test special character handling
   - Verification: `go test ./internal/parser/... -run TestMarkdown` fails

- [x] **Configure Goldmark with auto heading IDs**
   - Create configured Goldmark instance with `parser.WithAutoHeadingID()`
   - Update `MarkdownToHTMLWithError` to use configured instance
   - Verification: Tests pass

- [x] **Manual verification**
   - Build dev site: `cd dev && ../bin/ssg build`
   - Verify heading IDs in output: `grep -r 'id="' public/`
   - Test browser navigation to `#anchor` links

---

## Phase 2: Clickable Anchor Links

- [x] **Write failing tests for anchor link wrapping**
   - Test heading content is wrapped in `<a href="#id">`
   - Test all heading levels (h1-h6) have anchor links
   - Verification: `go test ./internal/parser/... -run TestMarkdown` fails

- [x] **Create custom Goldmark heading renderer**
   - Create `internal/parser/anchor_renderer.go`
   - Implement custom renderer that wraps heading content in anchor tags
   - Register renderer with Goldmark instance
   - Verification: Tests pass

- [x] **Add CSS styling for anchor links**
   - Heading links inherit text color
   - Remove default underline
   - Add hover underline
   - Verification: Visual inspection in browser

- [x] **Manual verification**
   - Build dev site: `make dev`
   - Click heading → URL updates with fragment
   - Right-click heading → "Copy link" includes fragment
   - Refresh with fragment → scrolls to heading

## Dependencies

None - single package change with no external dependencies.

## Parallelizable Work

All tasks are sequential (TDD workflow).
