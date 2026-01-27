## Context

The SSG currently extracts date anchors from markdown content using `ExtractDateAnchors()` in `parser.go`. These dates are passed through to the template as a flat `[]string` slice and rendered as a single collapsible list in `base.html`. As content grows (daily journal entries), this list becomes unwieldy.

Current data flow:
```
markdown → ExtractDateAnchors() → Page.DateAnchors → template → flat <ul>
```

Proposed data flow:
```
markdown → ExtractDateAnchors() → GroupDatesByMonth() → GroupMonthsByYear() → Page.CurrentMonthDates
                                                                                        → Page.ArchivedYears
                                                                                                ↓
                                                                              template → two-column nav
                                                                                         (year > month > date)
```

## Goals / Non-Goals

**Goals:**
- Limit the "Jump to date" section to current month only
- Group previous months into collapsible archive sections
- Maintain pure HTML/CSS (no JavaScript)
- Responsive layout (side-by-side on desktop, stacked on mobile)
- Preserve backwards compatibility (keep `DateAnchors` field)

**Non-Goals:**
- Pagination or "load more" functionality
- Customizable archive display format
- Search/filter functionality

## Decisions

### Decision 1: "Current month" definition

**Choice**: The month containing the most recent date entry (content-driven).

**Alternatives considered**:
- Build-time current (use `time.Now().Month()`) - Rejected because if user doesn't post for a week into February, January would disappear from current view
- Explicit configuration - Rejected as unnecessary complexity

**Rationale**: Content-driven is intuitive—whatever month you last posted in stays "current" until you post in a new month. Since the site builds every minute, this aligns with user expectations.

### Decision 2: Data structure for grouped dates

**Choice**: Add `MonthGroup` struct and two new Page fields.

```go
type MonthGroup struct {
    Year   int
    Month  string   // "January", "February"
    Dates  []string
}

type YearGroup struct {
    Year   int
    Months []MonthGroup
}

// Page additions:
CurrentMonthDates []string
ArchivedYears     []YearGroup
```

**Alternatives considered**:
- Single `map[string][]string` keyed by "YYYY-MM" - Rejected because maps don't preserve order in Go templates
- Compute grouping in template - Rejected because Go templates lack this capability

**Rationale**: Explicit struct with year/month fields gives clear template access. Separate slices for current vs. archived simplifies template logic. YearGroup allows collapsible year-level hierarchy in the archive.

### Decision 3: Collapsible mechanism

**Choice**: Native HTML5 `<details>/<summary>` elements.

**Alternatives considered**:
- CSS-only checkbox hack - More complex, less semantic
- JavaScript accordion - Violates no-JS constraint

**Rationale**: Already using `<details>` for current date nav. Consistent and semantic.

### Decision 4: Layout approach

**Choice**: CSS Flexbox with media query breakpoint at 768px.

```css
.date-nav-container { display: flex; }
@media (max-width: 768px) { flex-direction: column; }
```

**Alternatives considered**:
- CSS Grid - Overkill for two-column layout
- Float-based layout - Outdated, harder to make responsive

**Rationale**: Flexbox is simple, well-supported, and handles the responsive switch elegantly.

### Decision 5: Year-level collapsible grouping

**Choice**: Nest months under collapsible year headers in the archive section.

**Structure**:
```
> 2025 (collapsed)
  > December
    - 2025-12-31
    - 2025-12-20
  > November
    - 2025-11-15
```

**Alternatives considered**:
- Flat month list (original design) - Rejected because the list grows unwieldy with many months
- Year headers without collapse - Rejected because it still presents a long list

**Rationale**: As the archive grows, year-level collapse provides a scalable hierarchy. Users can expand a specific year to see its months, then expand a month to see dates.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Date parsing failures | Use strict `time.Parse()` with "2006-01-02" format; malformed dates excluded from grouping |
| Empty archive looks awkward | Show "No archives yet" placeholder text |
| Month name localization | Use English names only (consistent with existing ISO date format) |
| Many archived months | Months nested under year headers; users expand year first, then month |
| Nested collapse depth | Two levels (year > month) is manageable; avoid deeper nesting |
