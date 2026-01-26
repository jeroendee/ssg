# Tasks: Add Date Navigation

## Implementation Order

### 1. Add DateAnchors field to Page model
- [x] **File**: `internal/model/model.go`
- [x] **Change**: Add `DateAnchors []string` field to Page struct
- [x] **Validation**: Existing tests pass, field accessible

### 2. Implement date anchor extraction in parser
- [x] **File**: `internal/parser/parser.go`
- [x] **Change**: Add `ExtractDateAnchors(markdown string) []string` function
- [x] **Pattern**: Regex `^#### \*(\d{4}-\d{2}-\d{2})\*` for h4 date headings
- [x] **Integration**: Call in `ParsePage()` before markdown conversion
- [x] **Validation**: Unit tests for extraction with various inputs

### 3. Pass DateAnchors to template in renderer
- [x] **File**: `internal/renderer/renderer.go`
- [x] **Change**: Add DateAnchors to templateData.Page struct, populate in `RenderPage()`
- [x] **Validation**: Template receives date anchors

### 4. Add date navigation template markup
- [x] **File**: `internal/renderer/templates/base.html`
- [x] **Change**: Add `<details>/<summary>` dropdown when `.Page.DateAnchors` non-empty
- [x] **Validation**: HTML renders correctly, anchor links work

### 5. Add CSS styling for date navigation
- [x] **File**: `dev/public/style.css`
- [x] **Change**: Add `.date-nav` styles using existing Solarized variables
- [x] **Validation**: Dropdown styled correctly in light/dark mode

### 6. End-to-end verification
- [x] **Action**: Build site, test dropdown on `/now` page
- [x] **Validation**: Dropdown appears, dates listed, anchor navigation works

## Dependencies

- Task 2 depends on Task 1 (needs DateAnchors field)
- Task 3 depends on Task 2 (needs extraction function)
- Task 4 depends on Task 3 (needs template data)
- Task 5 can run in parallel with Tasks 2-4
- Task 6 depends on all previous tasks

## Parallelizable Work

- Tasks 1, 2, 3, 4 are sequential (data flow dependency)
- Task 5 (CSS) can be done in parallel with Tasks 2-4
