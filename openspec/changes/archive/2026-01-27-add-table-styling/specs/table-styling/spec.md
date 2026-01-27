## ADDED Requirements

### Requirement: Table border styling

The system SHALL render tables with visible borders.

#### Scenario: Table has collapsed borders

- **WHEN** rendering a page containing a table
- **THEN** the table SHALL have collapsed borders (no double lines)

#### Scenario: All cells have borders

- **WHEN** rendering a page containing a table
- **THEN** all `<th>` and `<td>` elements SHALL have a 1px solid border
- **AND** the border color SHALL use the `--border` CSS variable

### Requirement: Table cell padding

The system SHALL render table cells with adequate padding.

#### Scenario: Cells have readable padding

- **WHEN** rendering a page containing a table
- **THEN** all `<th>` and `<td>` elements SHALL have padding
- **AND** the padding SHALL be at least 0.5rem vertical and 0.75rem horizontal

### Requirement: Table header differentiation

The system SHALL visually differentiate table headers from body cells.

#### Scenario: Header background color

- **WHEN** rendering a page containing a table with `<th>` elements
- **THEN** the header cells SHALL have a background color using `--bg-secondary`

#### Scenario: Header text emphasis

- **WHEN** rendering a page containing a table with `<th>` elements
- **THEN** the header text SHALL use the `--text-emphasis` color

### Requirement: Table spacing

The system SHALL render tables with appropriate vertical spacing.

#### Scenario: Table has vertical margin

- **WHEN** rendering a page containing a table
- **THEN** the table SHALL have vertical margin of at least 1rem
