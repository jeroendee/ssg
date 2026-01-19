# ADR-0005: No JavaScript

> Status: Accepted
> Date: 2026-01-19

## Context

Modern web development often includes JavaScript for interactivity. Options for a static site generator:

1. **No JavaScript**: Pure HTML/CSS output
2. **Optional JavaScript**: User can add if needed
3. **Built-in JavaScript**: Include for features like search, dark mode toggle
4. **Framework-based**: React, Vue, etc. for hydration

Considerations:
- Page load performance
- Accessibility
- Privacy (no tracking by default)
- Maintenance burden
- Target use case (personal blogs)

## Decision

Generate **pure HTML/CSS output with no JavaScript** by default.

**Behavior**:
- Templates produce semantic HTML
- Styling via CSS only (including dark mode)
- No JavaScript in default templates
- Analytics only via external services (opt-in)
- Users can add JavaScript via custom templates/assets

## Consequences

### Positive

- **Fast**: No JavaScript parsing or execution
- **Accessible**: Works without JavaScript enabled
- **Private**: No client-side tracking or fingerprinting
- **Secure**: Smaller attack surface
- **Simple**: Easier to reason about, maintain, test
- **Universal**: Works on any browser, any device

### Negative

- **No client interactivity**: No search, no dynamic features
- **CSS limitations**: Some features harder without JS
- **No SPA benefits**: Full page loads for navigation
- **Dark mode**: CSS-only (uses prefers-color-scheme, no toggle)

### Neutral

- Users who need JavaScript can add it manually
- External analytics (GoatCounter) is JS-based but opt-in
- Future: Could add JavaScript features as opt-in

## Implementation Notes

### Dark Mode

Use CSS media query instead of JavaScript toggle:

```css
@media (prefers-color-scheme: dark) {
  :root {
    --bg: #002b36;
    --text: #839496;
  }
}
```

Users get dark mode based on system preference automatically.

### Features to Implement in CSS

- Responsive layout (media queries)
- Dark mode (prefers-color-scheme)
- Hover states (`:hover`, `:focus`)
- Print styles (@media print)
- Smooth scrolling (`scroll-behavior: smooth`)

### Optional JavaScript

Users can add JavaScript by:
1. Creating `assets/script.js`
2. Adding custom templates that include it
3. Using external services with `<script>` tags

### Analytics Exception

GoatCounter or similar privacy-respecting analytics are allowed as opt-in:
- Configured via `analytics.goatcounter` in config
- Only included if explicitly configured
- Uses external script, not bundled

## References

- [No-JS Club](https://no-js.club/)
- [The Web Without JavaScript](https://www.wired.com/2015/11/i-turned-off-javascript-for-a-whole-week-and-it-was-glorious/)
- ADR-0002-embedded-templates.md (templates don't include JS)
