# 202512050830 - Reference Assets: Fix serve --build and Add Logo Support

## Objective

Make the example static site visually match https://www.qualityshepherd.nl by:

1. **Fix bug**: `ssg serve --build` doesn't copy assets (CSS missing)
2. **Add feature**: Logo support in site header
3. **Update example**: Configure logo in example site

### Root Cause Analysis

**Styling issue**: The serve command (`cmd/ssg/main.go:161`) creates a builder but never calls `SetAssetsDir()`, so assets aren't copied to `public/`.

**No images**: The template (`base.html`) has no logo support, and no logo file exists in `example/assets/`.

## Technology

- Go 1.23+
- HTML templates (Go `html/template`)
- YAML configuration

### Skills

Use these Skills:

- `tdd-workflow` - Write failing tests first for each task
- `go-writer` - Implement Go code changes
- `go-tester` - Verify test coverage

## Examples

**Live site header** (qualityshepherd.nl):
```html
<header>
  <img src="/logo.svg" alt="Quality Shepherd logo">
  <a href="/">Quality Shepherd</a>
  <nav>...</nav>
</header>
```

**Current example site header** (missing logo, no styling):
```html
<header>
  <a href="/" class="title">Quality Shepherd</a>
  <nav>...</nav>
</header>
```

## Implementation Plan

### Task 1: Fix serve command to copy assets

**Skill**: `tdd-workflow`, `go-tester`, `go-writer`

**Files**:
- `cmd/ssg/main.go:161` - Add `b.SetAssetsDir("assets")` call

**Steps**:
1. RED: Write test verifying assets are copied when using `serve --build`
2. GREEN: Add `b.SetAssetsDir("assets")` before `b.Build()` in serve command
3. REFACTOR: Ensure consistency with build command's asset handling

### Task 2: Add logo support to templates

**Skill**: `tdd-workflow`, `go-tester`, `go-writer`

**Files**:
- `internal/model/config.go` - Add `Logo string` field to Site struct
- `internal/config/config.go` - Parse `site.logo` from YAML
- `internal/renderer/templates/base.html` - Render logo conditionally

**Steps**:
1. RED: Write test for logo field in config parsing
2. GREEN: Add Logo field to model and config loader
3. RED: Write test for logo rendering in template
4. GREEN: Update base.html template:
   ```html
   <header>
       {{if .Site.Logo}}<img src="{{.Site.Logo}}" alt="{{.Site.Title}} logo">{{end}}
       <a href="/" class="title">{{.Site.Title}}</a>
       <nav>...</nav>
   </header>
   ```
5. REFACTOR: Ensure template handles missing logo gracefully

### Task 3: Update example site configuration

**Files**:
- `example/ssg.yaml` - Add `logo: /logo.svg`
- `example/assets/logo.svg` - Add logo file

**Steps**:
1. Add logo.svg file to `example/assets/`
2. Update `example/ssg.yaml` with logo configuration:
   ```yaml
   site:
     title: Quality Shepherd
     logo: /logo.svg
   ```

## Critical Files

| Task | File | Change |
|------|------|--------|
| 1 | `cmd/ssg/main.go` | Add `b.SetAssetsDir("assets")` at line ~161 |
| 2 | `internal/model/config.go` | Add `Logo string` to Site struct |
| 2 | `internal/config/config.go` | Parse logo from YAML |
| 2 | `internal/renderer/templates/base.html` | Add conditional logo img |
| 3 | `example/ssg.yaml` | Add logo config |
| 3 | `example/assets/logo.svg` | Add logo file |
