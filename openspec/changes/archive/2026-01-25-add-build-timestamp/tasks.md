# Tasks: Add Build Timestamp File

## Implementation Tasks

1. **Add build timestamp generation to builder**
   - Create `writeBuildTimestamp()` method in `builder.go`
   - Generate RFC 3339 UTC timestamp
   - Write `build.json` to output root
   - Dependencies: None
   - Verification: Unit test

2. **Call timestamp generation in Build() pipeline**
   - Add call after sitemap/robots generation
   - Always generate (not conditional on posts existing)
   - Dependencies: Task 1
   - Verification: Integration test

3. **Add tests for build timestamp**
   - Table-driven test for `writeBuildTimestamp()`
   - Integration test verifying file exists and format
   - Dependencies: Task 1, 2
   - Verification: `make test` passes

4. **Update spec delta**
   - Apply spec delta to `site-building` spec
   - Dependencies: Task 3
   - Verification: `openspec validate`

## Validation Checklist

- [x] `build.json` created in output directory
- [x] JSON contains valid `buildTime` field
- [x] Timestamp is RFC 3339 format (e.g., `2026-01-25T14:30:00Z`)
- [x] File is UTF-8 encoded
- [x] All tests pass
- [x] `make all` succeeds
