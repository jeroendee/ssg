## Why

The "Jump to date" navigation list grows indefinitely as new daily entries are added. With multiple months of content, the list becomes unwieldy. Users need a way to navigate to older dates without scrolling through a massive list, while keeping the current month's dates readily accessible.

## What Changes

- Split date navigation into two sections: current month (left) and archived months (right)
- Current month shows individual date links as before
- Previous months collapse into expandable `<details>` sections, nested under collapsible year headers
- "Current month" is defined as the month containing the most recent date entry
- Archived months are ordered newest-first
- When no archives exist, show "No archives yet" placeholder
- Mobile layout stacks vertically (current month on top)

## Capabilities

### New Capabilities

- `date-grouping`: Groups date anchors by year/month, nests months under years, and identifies which month is "current" (most recent with entries)

### Modified Capabilities

- `html-rendering`: The "Render Date Navigation" requirement changes from a single flat list to a two-column layout with current month and collapsible archive sections

## Impact

- **Parser**: New `GroupDatesByMonth()` and `GroupMonthsByYear()` functions to categorize dates
- **Model**: New `MonthGroup` and `YearGroup` types; Page struct gains `CurrentMonthDates` and `ArchivedYears` fields
- **Renderer**: Template data struct expanded; `base.html` template updated
- **CSS**: New flexbox layout styles for `.date-nav-container`, responsive breakpoints
- **No breaking changes**: Existing `DateAnchors` field retained for backwards compatibility
