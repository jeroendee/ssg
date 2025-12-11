# 202512070900 - Header/Logo Responsive Behavior Fix

## Objective

Fix the header logo clipping issue at certain viewport widths. The logo must be fully visible at ALL viewport widths.

## Problem

The header logo gets clipped at viewport widths between ~480-720px. The current CSS has only one breakpoint at 480px which is too narrow to catch the "middle range" where logo and title compete for horizontal space.

## Reference Behavior (from qualityshepherd.nl)

| Viewport | Layout |
|----------|--------|
| Narrow (~320px) | Logo ABOVE title, left-aligned. Title and nav below. |
| Wider (~375px+) | Logo to the RIGHT of title (horizontal layout) |

## Technology

### Skills

Use these Skills:

- `html-css-writer` - For writing CSS styling changes

## Implementation Plan

### Analysis

**Root Cause:** Desktop-first CSS with single breakpoint at 480px doesn't handle mid-range viewports properly.

| Element | Current | Target |
|---------|---------|--------|
| Default layout | Row (horizontal) | Column (vertical, stacked) |
| Breakpoint | 480px max-width | 500px min-width |
| Logo protection | None | `flex-shrink: 0`, `max-width: 100%` |

### Files to Modify

- `/Users/jeroendee/ACID/40-49/44/44.S/ssg/assets/style.css`
- `/Users/jeroendee/ACID/40-49/44/44.S/ssg/dev/assets/style.css` (keep in sync)

### Task 1: Update header CSS to mobile-first approach

**Skill:** Invoke `html-css-writer`

Replace the header CSS block (lines 34-74) and the mobile breakpoint (lines 160-175) with:

```css
/* Header styling - mobile-first approach */
header {
  padding: 20px 20px;
  display: flex;
  flex-direction: column;  /* Stack by default (narrow screens) */
  align-items: flex-start;
  gap: 1rem;
  margin: -20px -20px 20px -20px;
  width: calc(100% + 40px);
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.site-title {
  font-size: 1.8rem;
  font-weight: bold;
  text-decoration: none;
}

.site-title h2 {
  font-size: inherit;
  margin: 0;
}

header nav a {
  text-decoration: underline;
  margin-right: 15px;
}

.header-right {
  display: flex;
  align-items: center;
  order: -1;  /* Logo appears first (above title) on narrow screens */
  align-self: flex-start;
}

header .logo {
  width: auto;
  height: 80px;
  max-width: 100%;  /* Never exceed container */
  flex-shrink: 0;   /* Never shrink the logo */
}

/* Wider viewports: side-by-side layout */
@media (min-width: 500px) {
  header {
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
  }

  .header-right {
    order: 0;  /* Reset order - logo on right */
    align-self: auto;
  }
}
```

### Task 2: Remove old 480px mobile breakpoint

Delete the old `@media (max-width: 480px)` block for header (it's replaced by mobile-first default).

### Task 3: Rebuild SSG binary

```bash
go build -o /tmp/ssg ./cmd/ssg
```

### Task 4: Visual verification with Playwright

Test at these viewport widths:
- 320px (narrow mobile)
- 375px (iPhone SE)
- 430px (iPhone 14 Pro Max)
- 500px (breakpoint boundary)
- 600px (small tablet)
- 720px (max-width boundary)
- 1024px (desktop)

Verify:
- [ ] Logo is 100% visible at ALL widths (no clipping)
- [ ] Narrow viewports (<500px): Logo appears above title
- [ ] Wider viewports (>=500px): Logo appears to right of title
