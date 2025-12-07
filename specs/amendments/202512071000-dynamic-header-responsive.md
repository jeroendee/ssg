# 202512071000 - Dynamic Header Responsive with Intrinsic Flexbox

## Objective

Replace fixed media query breakpoints with pure flexbox intrinsic sizing for header layout. The logo must be fully visible at ALL viewport widths without relying on fixed pixel breakpoints.

## Problem

The current implementation uses a fixed 768px media query breakpoint. At 769px viewport width, the header's effective content width is only 760px (due to body max-width: 720px + padding: 20px × 2 = 760px). This "off-by-8px" mismatch causes the logo to clip when the viewport is just above the breakpoint.

**Current approach:** Fixed breakpoint at 768px determines layout switch.
**Issue:** Breakpoint doesn't account for actual available header width.

## Solution: Intrinsic Flexbox Sizing

Use flexbox's natural wrapping behavior to let the content itself determine when to stack vs. display horizontally. No fixed breakpoint needed.

### Key Technique: flex-wrap + flex-basis

```css
header {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.header-left {
  flex: 1 1 300px;  /* Grow, shrink, min 300px before wrapping */
}

.header-right {
  flex: 0 0 auto;   /* Don't grow/shrink, use natural size */
}
```

**How it works:**
- `flex-wrap: wrap` allows items to wrap to next line when space is insufficient
- `flex-basis: 300px` on header-left means it needs at least 300px
- When header-left + header-right + gap exceeds container width, flexbox automatically wraps
- The breakpoint is now *intrinsic* to the content, not a fixed pixel value

## Technology

### Skills

Use these Skills:

- `html-css-writer` - For writing CSS styling changes

## Implementation Plan

### Files to Modify

- `/Users/jeroendee/ACID/40-49/44/44.S/ssg/assets/style.css`

### Task 1: Update header CSS with intrinsic flexbox

**Skill:** Invoke `html-css-writer`

Replace the header CSS (lines 34-83) and remove the 768px media query (lines 165-177) with:

```css
/* Header styling - intrinsic responsive using flexbox */
header {
  padding: 20px;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 1rem;
  margin: -20px -20px 20px -20px;
  width: calc(100% + 40px);
}

.header-left {
  flex: 1 1 300px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
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
  flex: 0 0 auto;
  display: flex;
  align-items: center;
}

header .logo {
  height: 80px;
  width: auto;
  max-width: 100%;
  flex-shrink: 0;
}
```

### Task 2: Remove fixed 768px media query

Delete the `@media (min-width: 768px)` block for header (lines 165-177). The intrinsic flexbox approach eliminates the need for this breakpoint.

### Task 3: Visual verification with Playwright

Test at these viewport widths to verify logo is never clipped:
- 320px (narrow mobile)
- 480px (wide mobile)
- 600px (small tablet)
- 720px (body max-width boundary)
- 768px (old breakpoint - critical test)
- 769px (the previously broken width)
- 800px (comfortable desktop)
- 1024px (desktop)

**Acceptance criteria:**
- [ ] Logo is 100% visible at ALL widths (no clipping)
- [ ] Layout wraps naturally based on available space
- [ ] No fixed pixel breakpoint in header CSS
