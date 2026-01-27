## Context

The monthly-date-archive feature added `<details open>` to the current-month date navigation in `base.html`. The `open` attribute causes the dropdown to be expanded on page load. The original date-navigation feature used `<details>` without `open`, which kept it collapsed.

## Goals / Non-Goals

**Goals:**
- Restore collapsed-by-default behavior for the "Jump to date" dropdown

**Non-Goals:**
- Changing the archived years/months behavior (already collapsed correctly)
- Adding JavaScript toggle functionality
- Changing the visual styling

## Decisions

**Remove `open` attribute from `<details>` element**

The HTML `<details>` element is collapsed by default. Adding the `open` attribute expands it. Simply removing `open` restores the expected behavior.

Location: `internal/renderer/templates/base.html`, line 11

```html
<!-- Current (bug) -->
<details open>

<!-- Fixed -->
<details>
```

No alternatives needed - this is the standard HTML approach.

## Risks / Trade-offs

None. This is a one-attribute fix with no side effects.
