## Context

Date-anchored pages render a two-column date navigation: "Jump to date" (current month) on the left, and an archive section on the right. The archive currently has year `<details>` at the top level. We want to wrap all year groups in a single "History" `<details>` and hide the section entirely when empty.

## Goals / Non-Goals

**Goals:**
- Add a "History" `<details><summary>` wrapper around the existing year groups in `base.html`
- Remove the "No archives yet" placeholder — hide the archive `<div>` entirely when empty

**Non-Goals:**
- Changing the Go model, parser, or builder
- Modifying CSS (existing `.archive details details` rule cascades to the new nesting depth)
- Changing the "Jump to date" (current month) section

## Decisions

**Template-only change in `base.html`**: The archive `<div>` content is replaced with a single `{{if .Page.ArchivedYears}}` block containing a `<details><summary>History</summary>` that wraps the existing `{{range .Page.ArchivedYears}}` loop. When `ArchivedYears` is empty, no archive markup is rendered at all.

_Alternative considered_: Keeping the archive `<div>` visible with a "History" summary that opens to "No archives yet". Rejected — an empty disclosure adds no value and breaks the visual symmetry with "Jump to date".

## Risks / Trade-offs

**[Minor layout shift when archive is empty]** → On pages with date anchors but only one month of data, the right column disappears entirely. Acceptable — flexbox collapses gracefully and the "Jump to date" column fills the space.
