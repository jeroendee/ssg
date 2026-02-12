## Why

Topic counts (e.g., "agent (27)") add visual noise without providing value to readers. The words alone, ordered by frequency, already communicate what a page is about. Stripping the counts gives a cleaner, calmer topic bar.

## What Changes

- Remove the `(count)` display from topic words in the HTML template
- Topics remain sorted by frequency descending — only the visible count is removed
- The `Topic.Count` field and extraction logic stay unchanged (count drives sort order internally)

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `html-rendering`: The "Render topics bar" requirement changes from `"word (count)"` format to displaying only the word.

## Impact

- **Template**: `internal/renderer/templates/base.html` — remove `({{$t.Count}})` from topic rendering
- **Tests**: `internal/renderer/renderer_test.go` — update assertions that check for count in displayed output
- **Spec**: `openspec/specs/html-rendering/spec.md` — update the topics bar scenario
