## 1. Template Change

- [x] 1.1 Remove `({{$t.Count}})` from the topic rendering line in `internal/renderer/templates/base.html`

## 2. Test Updates

- [x] 2.1 Update `TestRenderPage_WithTopics` in `internal/renderer/renderer_test.go` to assert words without counts

## 3. Verify

- [x] 3.1 Run `make test` — all tests pass
- [x] 3.2 Run `make dev` and visually confirm topics bar shows words only
