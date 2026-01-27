## 1. Model Layer

> **Skill**: `go-writer` — Implement Go struct additions following Google Go Style Guide

- [x] 1.1 Add `MonthGroup` struct to `internal/model/model.go`
- [x] 1.2 Add `CurrentMonthDates []string` field to Page struct
- [x] 1.3 Add `ArchivedMonths []MonthGroup` field to Page struct (to be replaced)
- [x] 1.4 Add `YearGroup` struct to `internal/model/model.go`
- [x] 1.5 Replace `ArchivedMonths` with `ArchivedYears []YearGroup` field in Page struct

## 2. Parser Layer

> **Skill**: `go-tester` (2.2, 2.5), `go-writer` (2.1, 2.3, 2.4, 2.6) — TDD workflow: write tests first, then implementation
> **Sub-agent**: `go-engineer` — Coordinates go-writer and go-tester for the full TDD cycle

- [x] 2.1 Create `GroupDatesByMonth()` function in `internal/parser/parser.go`
- [x] 2.2 Write unit tests for `GroupDatesByMonth()`: single month, multiple months, empty input, malformed dates
- [x] 2.3 Update `ParsePage()` to call `GroupDatesByMonth()` and populate new Page fields
- [x] 2.4 Create `GroupMonthsByYear()` function in `internal/parser/parser.go`
- [x] 2.5 Write unit tests for `GroupMonthsByYear()`: multiple years, single year, empty input, year ordering
- [x] 2.6 Update `ParsePage()` to call `GroupMonthsByYear()` and populate `ArchivedYears`

## 3. Renderer Layer

> **Skill**: `go-writer` (3.1, 3.2, 3.4), `html-css-writer` (3.3, 3.5) — Go struct updates and HTML template changes

- [x] 3.1 Add `CurrentMonthDates` and `ArchivedMonths` to `templateData.Page` struct in `internal/renderer/renderer.go`
- [x] 3.2 Pass new fields through in `RenderPage()`
- [x] 3.3 Update `base.html` template with two-column date navigation layout
- [x] 3.4 Replace `ArchivedMonths` with `ArchivedYears` in `templateData.Page` struct
- [x] 3.5 Update `base.html` template with nested year > month > date structure

## 4. Styling

> **Skill**: `html-css-writer` — Semantic CSS with flexbox responsive patterns

- [x] 4.1 Add `.date-nav-container` flexbox styles to `internal/assets/default_style.css`
- [x] 4.2 Add `.current-month` and `.archive` section styles
- [x] 4.3 Add `.archive-empty` placeholder styles
- [x] 4.4 Add mobile responsive styles with `@media (max-width: 768px)` breakpoint
- [x] 4.5 Add nested details indentation style (`.archive details details { margin-left: 1rem; }`)

## 5. Integration Testing

> **Skill**: `go-tester` — Manual verification with build commands
> **Sub-agent**: `go-engineer` — Runs diagnostics and validates build

- [x] 5.1 Build site with `go run ./cmd/ssg build`
- [x] 5.2 Verify desktop layout (side-by-side sections)
- [x] 5.3 Verify mobile layout (stacked sections)
- [x] 5.4 Verify archive links navigate to correct date headings
- [x] 5.5 Verify "No archives yet" appears when only one month of content exists
- [x] 5.6 Verify year headers appear collapsed by default
- [x] 5.7 Verify clicking year expands to show months
- [x] 5.8 Verify clicking month expands to show dates
- [x] 5.9 Verify nested indentation displays correctly
