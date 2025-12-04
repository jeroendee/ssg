---
title: Testing in Go
---

Go has excellent built-in testing support. Here's what makes it great:

## Table-Driven Tests

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        a, b, want int
    }{
        {1, 2, 3},
        {0, 0, 0},
        {-1, 1, 0},
    }

    for _, tt := range tests {
        got := Add(tt.a, tt.b)
        if got != tt.want {
            t.Errorf("Add(%d, %d) = %d, want %d",
                tt.a, tt.b, got, tt.want)
        }
    }
}
```

## Parallel Tests

Use `t.Parallel()` to run tests concurrently and make your test suite faster.
