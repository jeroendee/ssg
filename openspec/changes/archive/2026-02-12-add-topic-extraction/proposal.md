## Why

Date-anchored pages like `/moments/` accumulate content over time, making it hard for readers to discover recurring themes. Auto-extracted topic words — displayed as a frequency-ranked list — give readers a quick overview of what subjects appear most, similar to an old-school tag cloud but rendered as plain comma-separated text.

## What Changes

- New `topics` config section with a `pages` list, following the same pattern as `feed.pages`
- New `internal/topics/` package that extracts subject words from markdown, filters stop words, counts frequencies, and returns the top N results
- `model.Page` gains a `Topics` field holding extracted topic/count pairs
- Builder wires topic extraction for configured pages during `ScanContent`
- `base.html` template renders a topics bar (between date-nav and content) for pages that have topics
- Stop-word list filters English common words; only words ≥ 3 characters appearing ≥ 2 times qualify
- Top 20 topics displayed by frequency descending, with counts: `agent (27), claude (21), ...`

## Capabilities

### New Capabilities
- `topic-extraction`: Extracts recurring subject words from page markdown content and renders them as a frequency-ranked discovery bar on configured pages

### Modified Capabilities
- `site-configuration`: Adds `topics.pages` config section for opting pages into topic extraction
- `html-rendering`: Adds topics bar rendering in the page template for pages with extracted topics

## Impact

- **New package**: `internal/topics/` (word tokenizer, stop words, frequency counter)
- **Model**: `model.Page` gets new `Topics` field
- **Config**: `model.Config` gets new `TopicPages` field; `internal/config/` parses `topics.pages`
- **Builder**: `ScanContent` calls topic extraction for configured pages
- **Templates**: `base.html` gains a conditional topics section
- **CSS**: Minimal styling for the topics bar in default stylesheet
