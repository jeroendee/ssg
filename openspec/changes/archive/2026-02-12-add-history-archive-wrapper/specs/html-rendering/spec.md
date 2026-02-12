## MODIFIED Requirements

### Requirement: Render Date Navigation

The system SHALL render date navigation with current month and archived months sections.

#### Scenario: Archive section with previous months

- **WHEN** rendering date navigation with archived months
- **THEN** the archive section SHALL contain a single collapsible `<details>` element with summary text "History"
- **AND** within the "History" element, archived months SHALL be nested under collapsible year headers
- **AND** each year SHALL be rendered as a collapsible details element with summary showing the year number
- **AND** within each year, each month SHALL be rendered as a nested collapsible details element
- **AND** the month summary text SHALL show the month name (e.g., "December")
- **AND** years SHALL be ordered newest to oldest
- **AND** months within each year SHALL be ordered newest to oldest
- **AND** each date within a month SHALL be rendered as an anchor link
- **AND** the "History" details element SHALL be collapsed by default

#### Scenario: Archive section with no previous months

- **WHEN** rendering date navigation with no archived years
- **THEN** the archive section SHALL NOT render any markup
- **AND** the archive section SHALL NOT display "No archives yet" placeholder text

#### Scenario: Nested details styling

- **WHEN** rendering nested month details within year details within the History wrapper
- **THEN** the nested details SHALL have left indentation for visual hierarchy
