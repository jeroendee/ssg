## Why

The "Jump to date" navigation list on pages with date anchors is expanded by default on page load. This was introduced during the monthly-date-archive feature. The expected behavior is for it to be collapsed by default.

## What Changes

- Remove the `open` attribute from the current-month `<details>` element in `base.html`
- Restores collapsed-by-default behavior for the date navigation dropdown

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `html-rendering`: The date navigation should be collapsed by default, not expanded

## Impact

- `internal/renderer/templates/base.html`: Single attribute removal
- No API changes
- No breaking changes
