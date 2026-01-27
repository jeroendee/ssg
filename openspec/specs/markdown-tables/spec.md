## ADDED Requirements

### Requirement: Markdown tables render as HTML tables

The markdown parser SHALL convert GFM-style table syntax into HTML `<table>` elements.

#### Scenario: Simple table renders correctly
- **WHEN** markdown contains a table with headers and rows:
  ```
  | Header 1 | Header 2 |
  |----------|----------|
  | Cell 1   | Cell 2   |
  ```
- **THEN** output contains `<table>`, `<thead>`, `<tbody>`, `<tr>`, `<th>`, and `<td>` elements

#### Scenario: Table with alignment renders correctly
- **WHEN** markdown contains a table with column alignment:
  ```
  | Left | Center | Right |
  |:-----|:------:|------:|
  | L    | C      | R     |
  ```
- **THEN** output contains appropriate alignment attributes or styles

#### Scenario: Table content preserves inline formatting
- **WHEN** table cells contain inline markdown (bold, links, code)
- **THEN** inline formatting is rendered within the table cells
