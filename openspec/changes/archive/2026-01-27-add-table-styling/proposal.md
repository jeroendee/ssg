## Why

Tables rendered from markdown have no visual styling, making them difficult to read. Only `width: 100%` is applied, with no borders, cell padding, or header differentiation.

## What Changes

- Add bordered table styling to `default_style.css`
- Include cell padding for readability
- Differentiate header row with background color
- Use existing Solarized CSS variables for consistency

## Capabilities

### New Capabilities
- `table-styling`: CSS styling rules for HTML tables including borders, padding, and header differentiation

### Modified Capabilities
None - this is purely additive CSS, no existing spec behavior changes.

## Impact

- `internal/assets/default_style.css`: Add table styling rules
- All pages with markdown tables will display with improved readability
