## Why

Date-anchored pages (e.g., /moments/) display an archive section where year groups sit at the top level alongside "Jump to date". Adding a "History" wrapper creates a clearer visual hierarchy — the two top-level items ("Jump to date" and "History") are symmetrical, and the year/month/date tree is nested one level deeper where it belongs.

## What Changes

- Wrap the archive year groups in a collapsible "History" `<details>` element in the page template
- Hide the "History" section entirely when there are no archived years (instead of showing "No archives yet")

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `html-rendering`: The archive section of date navigation gains a "History" wrapper `<details>` element. The "No archives yet" placeholder is removed — when empty, the entire archive section is hidden.

## Impact

- `internal/renderer/templates/base.html` — template change to add the "History" wrapper and adjust the empty-state conditional
- No Go code changes (model, parser, builder untouched)
- No CSS changes (existing `.archive details details` indentation cascades naturally to the new depth)
