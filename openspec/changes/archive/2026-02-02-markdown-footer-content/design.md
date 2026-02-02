## Context

The site footer currently displays a hardcoded "Made with ❤️ in Amsterdam" message with an optional version badge. Contact information lives on a separate `/contact/` page, requiring navigation to access.

**Current state:**
- Footer defined in `_footer.html` template with static text
- Contact page at `content/contact.md` with its own nav entry
- Builder already skips `_` prefixed files from page generation (`isIndexFile()`)
- Renderer passes `Site` struct to all templates

**Constraint:** Footer content must remain optional - sites without `_footer.md` should work unchanged.

## Goals / Non-Goals

**Goals:**
- Enable markdown content in footer via `content/_footer.md`
- Consolidate contact info into footer, removing dedicated page
- Add visual separator (`<hr>`) between main content and footer
- Preserve existing "Made with ❤️" tagline and version badge

**Non-Goals:**
- Multiple footer partials or footer sections
- Dynamic footer content (different per page)
- Footer configuration in ssg.yaml

## Decisions

### Decision 1: FooterContent as Site-level field

Add `FooterContent string` to `model.Site` struct rather than template-level data.

**Rationale:** Footer is site-wide, not page-specific. Site struct already flows to all templates. Keeps template data structures simple.

**Alternative considered:** Add to each `templateData` struct → More code duplication, no benefit.

### Decision 2: Load in ScanContent, not Build

Parse `_footer.md` during `ScanContent()` alongside pages/posts.

**Rationale:** Content discovery belongs together. `ScanContent()` already handles markdown parsing for all other content types.

**Alternative considered:** Load during `Build()` → Would fragment content loading logic.

### Decision 3: Raw HTML in template (no escaping)

Store rendered HTML in `FooterContent` and output with `{{.Site.FooterContent}}` (Go templates auto-escape by default, so we'll need `template.HTML` type or explicit raw output).

**Rationale:** Markdown is already converted to HTML. Double-escaping would break links.

**Implementation note:** Check if `Site.FooterContent` needs to be `template.HTML` or if template changes are needed.

### Decision 4: Horizontal rule styling

Use semantic `<hr>` element with CSS styling rather than border tricks.

**Rationale:** Semantic HTML is accessible. CSS can style as needed.

## Risks / Trade-offs

**Risk:** Empty `_footer.md` file produces empty div in HTML.
→ Mitigation: Template checks `{{if .Site.FooterContent}}` before rendering.

**Risk:** Malformed markdown in `_footer.md` could break build.
→ Mitigation: Use existing `MarkdownToHTMLWithError` which returns error on parse failure.

**Trade-off:** Footer content is not configurable via ssg.yaml.
→ Acceptable: Markdown file gives more flexibility for formatting.

## Implementation Guidance

### Go Skills Required

| Task | Skills | Rationale |
|------|--------|-----------|
| 1.1 Model | `go-writer`, `go-tester` | Struct field addition with test |
| 2.1 Builder | `go-writer`, `go-tester`, `go-lsp` | New loading logic, LSP for symbol refs |
| 2.2 Tests | `go-tester` | Table-driven tests for 3 scenarios |
| 3.2 template.HTML | `go-writer`, `go-tester` | Type change if template escapes HTML |
| 5.x Verify | `go-tester`, `go-analyzer` | Integration test + static analysis |

**Non-Go tasks:** 3.1 uses `html-css-writer`; 4.x are file operations (no skill needed).

### Sub-Agent Parallelization

```
           ┌────────────┐
           │ 1.1 Model  │ ◄─── SEQUENTIAL: Must complete first
           │ FooterContent
           └─────┬──────┘
                 │
     ┌───────────┼───────────┐
     ▼           ▼           ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│ Builder │ │Template │ │ Content │  ◄─── PARALLEL WINDOW
│ 2.1-2.2 │ │ 3.1-3.2 │ │ 4.1-4.3 │
└────┬────┘ └────┬────┘ └────┬────┘
     │           │           │
     └───────────┼───────────┘
                 ▼
           ┌──────────┐
           │ 5.1-5.3  │ ◄─── SEQUENTIAL: Integration verification
           │ Verify   │
           └──────────┘
```

**Parallel agents (after 1.1):**

1. **Builder Agent** — Tasks 2.1, 2.2
   - Skills: `go-writer`, `go-tester`, `go-lsp`
   - Files: `internal/builder/builder.go`, `internal/builder/builder_test.go`

2. **Template Agent** — Tasks 3.1, 3.2
   - Skills: `html-css-writer`, `go-writer` (if type change needed)
   - Files: `templates/_footer.html`, possibly `internal/model/model.go`

3. **Content Agent** — Tasks 4.1, 4.2, 4.3
   - Skills: none (file operations)
   - Files: `dev/content/_footer.md`, `dev/content/contact.md`, `dev/ssg.yaml`
