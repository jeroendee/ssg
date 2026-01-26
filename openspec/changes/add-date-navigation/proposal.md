# Proposal: Add Date Navigation

## Summary

Add a build-time generated navigation dropdown to pages with date-based headings, allowing users to jump between date sections without JavaScript.

## Motivation

The `/now` page contains date-based entries (e.g., `#### *2026-01-26*`) that users need to navigate. Currently, users must scroll or use browser search. A dropdown navigation provides quick access to any date section.

## Scope

- **In scope**: Extract date anchors from markdown, render navigation dropdown
- **Out of scope**: JavaScript-based navigation, custom date formats, filtering

## Approach

1. **Parsing**: Extract h4 date headings (`#### *YYYY-MM-DD*`) from markdown source during page parsing
2. **Model**: Add `DateAnchors []string` field to Page struct
3. **Rendering**: Render `<details>/<summary>` dropdown with anchor links when dates exist
4. **Styling**: Add CSS for dropdown using existing Solarized color variables

## Constraints

- No JavaScript - uses native HTML5 `<details>/<summary>` element
- Build-time extraction - dates extracted from markdown source, not HTML
- Conditional rendering - dropdown only appears when page has date anchors
- Progressive enhancement - page works without dropdown support

## Affected Specs

- `content-parsing`: New requirement for extracting date anchors
- `html-rendering`: New requirement for rendering date navigation

## Risks

- **Low**: Feature is additive, no breaking changes
- **Low**: Uses well-supported HTML5 elements
