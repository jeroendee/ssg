## 1. Model Changes [SEQUENTIAL - must complete first]

- [x] 1.1 Add `FooterContent string` field to `Site` struct in `internal/model/model.go`
  - **Skills:** `go-writer`, `go-tester`

## 2. Builder Changes [PARALLEL - after group 1]

- [x] 2.1 Add footer loading logic to `ScanContent()` in `internal/builder/builder.go`
  - **Skills:** `go-writer`, `go-tester`, `go-lsp`
- [x] 2.2 Write tests for footer loading (exists, not exists, parse error)
  - **Skills:** `go-tester`

## 3. Template Changes [PARALLEL - after group 1]

- [x] 3.1 Update `_footer.html` to render optional footer content with `<hr>`
  - **Skills:** `html-css-writer`
- [x] 3.2 Verify footer content HTML is not escaped (may need `template.HTML`)
  - **Skills:** `go-writer`, `go-tester`

## 4. Content Changes [PARALLEL - after group 1]

- [x] 4.1 Create `dev/content/_footer.md` with contact information
  - **Skills:** none
- [x] 4.2 Delete `dev/content/contact.md`
  - **Skills:** none
- [x] 4.3 Remove Contact entry from navigation in `dev/ssg.yaml`
  - **Skills:** none

## 5. Verification [SEQUENTIAL - must be last]

- [x] 5.1 Run `make build` and verify footer renders correctly
  - **Skills:** `go-tester`, `go-analyzer`
- [x] 5.2 Run `go test ./...` and ensure all tests pass
  - **Skills:** `go-tester`
- [x] 5.3 Verify `/contact/` returns 404
  - **Skills:** none
