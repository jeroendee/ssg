## Context

The SSG uses a BearBlog-inspired CSS stylesheet with Solarized color scheme. Tables are rendered from markdown via goldmark's table extension, but currently only have `width: 100%` styling, making them difficult to read.

CSS variables already exist for consistent theming:
- `--bg-secondary`: Secondary background color
- `--border`: Border color
- `--text-emphasis`: Emphasized text color

## Goals / Non-Goals

**Goals:**
- Make tables visually distinct and readable
- Use existing CSS variables for dark/light mode compatibility
- Keep styling minimal and consistent with site aesthetic

**Non-Goals:**
- Zebra striping (alternating row colors)
- Responsive table behavior (horizontal scroll)
- Custom table classes or variants

## Decisions

### Decision 1: Bordered table style

**Choice**: Full borders on all cells with collapsed border style.

**Rationale**: Provides clear cell delineation without visual clutter. Matches the clean aesthetic of the existing stylesheet.

**Alternatives considered**:
- Lines-only (horizontal rules only): More minimal but harder to scan wide tables
- Borderless with padding: Insufficient visual separation

### Decision 2: Header differentiation

**Choice**: Background color (`--bg-secondary`) and emphasized text color on `<th>` elements.

**Rationale**: Clear visual hierarchy without adding new CSS variables. Leverages existing color scheme.

## Risks / Trade-offs

**[Minimal risk]** Tables with many columns may feel cramped → Existing `width: 100%` allows natural flow; padding provides breathing room.
