## Context

The topics bar currently renders each topic as `word (count)`, e.g. `agent (27), claude (21), docker (5)`. The count is extracted in `internal/topics/topics.go` and stored in `model.Topic.Count`. The template at `base.html:47` formats it. The count is also used for sort order.

## Goals / Non-Goals

**Goals:**
- Remove the visible count from the topics bar display
- Keep sort order by frequency intact

**Non-Goals:**
- Changing the `Topic` struct or extraction logic
- Changing topic sort order or filtering rules

## Decisions

**Display-only change in the template.**
Remove `({{$t.Count}})` from `base.html`. The `Count` field remains in the model — it drives sorting. This is a one-line template edit.

**No model changes needed.**
Alternative considered: removing `Count` from `Topic` entirely. Rejected because `Count` is needed internally for sorting and could be useful for future features (e.g., font sizing by frequency).

## Risks / Trade-offs

No meaningful risks. The change is purely cosmetic and confined to one template line plus test assertions.
