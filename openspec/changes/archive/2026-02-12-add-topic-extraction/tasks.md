## 1. Model & Config

- [x] 1.1 Add `Topic` struct (Word, Count) to `internal/model/model.go` and `Topics []Topic` field to `Page`
- [x] 1.2 Add `TopicPages []string` field to `model.Config`
- [x] 1.3 Add `topics.pages` parsing to `internal/config/config.go` (mirror `feed.pages` pattern), initialize to empty slice when omitted
- [x] 1.4 Add config tests: topic pages specified, not specified, empty array

## 2. Topic Extraction Package

- [x] 2.1 Create `internal/topics/topics.go` with `Extract(markdown string) []Topic` function
- [x] 2.2 Implement stop-word list (~150 common English words)
- [x] 2.3 Implement markdown stripping: links (keep text, discard URL), image refs, HTML entities
- [x] 2.4 Implement tokenization: split on non-alphanumeric (preserve hyphens within words), lowercase
- [x] 2.5 Implement filtering: min length 3, min frequency 2, stop-word exclusion
- [x] 2.6 Implement sorting: frequency descending, alphabetical tiebreaker, cap at 20
- [x] 2.7 Add comprehensive tests: basic extraction, stop words, short words, frequency threshold, top-20 cap, sort order, empty input, markdown stripping, hyphenated words

## 3. Builder Integration

- [x] 3.1 Wire topic extraction in `builder.ScanContent`: check page path against `cfg.TopicPages`, call `topics.Extract` on raw markdown body, set `page.Topics`
- [x] 3.2 Pass raw markdown body to topic extraction (requires storing body before HTML conversion in ParsePage or re-reading file in builder)
- [x] 3.3 Add builder tests: page with topics config gets topics populated, page without config gets empty topics

## 4. Template & Styling

- [x] 4.1 Add topics bar to `base.html`: conditional `<div class="topics">` between date-nav and content, render `word (count), word (count), ...`
- [x] 4.2 Add `.topics` CSS class to default stylesheet: small font, muted color
- [x] 4.3 Add renderer tests: page with topics renders topics bar, page without topics omits bar

## 5. Documentation & Config Example

- [x] 5.1 Add `topics.pages` section to `ssg.yaml.example` with comments
- [x] 5.2 Run full test suite (`make test`) and verify all tests pass
