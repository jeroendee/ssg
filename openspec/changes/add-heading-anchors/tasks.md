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

## Dependencies

None - single package change with no external dependencies.

## Parallelizable Work

All tasks are sequential (TDD workflow).
