## 1. Template Change

- [x] 1.1 In `internal/renderer/templates/base.html`, wrap the archive year loop in a `<details><summary>History</summary>...</details>` element, conditioned on `{{if .Page.ArchivedYears}}`
- [x] 1.2 Remove the `{{else}}` branch that renders the "No archives yet" placeholder

## 2. Verification

- [x] 2.1 Update the archive-related test in `internal/renderer/renderer_test.go` to assert the "History" summary text and verify no "No archives yet" output for empty archives
- [x] 2.2 Run `make test` — all tests pass
- [x] 2.3 Run `make serve`, verify /moments/-style page shows "History" wrapper with correct nesting
