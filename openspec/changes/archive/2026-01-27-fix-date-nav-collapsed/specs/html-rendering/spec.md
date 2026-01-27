## MODIFIED Requirements

### Requirement: Render Date Navigation

The system SHALL render date navigation with current month and archived months sections.

#### Scenario: Page with date anchors spanning multiple months

- **WHEN** rendering a page with dates from multiple months
- **THEN** the rendered output SHALL include a date navigation container
- **AND** the container SHALL have two sections: current month and archive

#### Scenario: Current month section

- **WHEN** rendering date navigation with current month dates
- **THEN** the current month section SHALL appear on the left (desktop) or top (mobile)
- **AND** the current month section SHALL use a details element with summary "Jump to date"
- **AND** the details element SHALL be collapsed by default
- **AND** each date SHALL be rendered as an anchor link with href matching the date

#### Scenario: Archive section with previous months

- **WHEN** rendering date navigation with archived months
- **THEN** the archive section SHALL appear on the right (desktop) or below current month (mobile)
- **AND** archived months SHALL be nested under collapsible year headers
- **AND** each year SHALL be rendered as a collapsible details element with summary showing the year number
- **AND** within each year, each month SHALL be rendered as a nested collapsible details element
- **AND** the month summary text SHALL show the month name (e.g., "December")
- **AND** years SHALL be ordered newest to oldest
- **AND** months within each year SHALL be ordered newest to oldest
- **AND** each date within a month SHALL be rendered as an anchor link

#### Scenario: Year-level collapse

- **WHEN** rendering archived years
- **THEN** each year details element SHALL be collapsed by default
- **AND** clicking a year summary SHALL expand to show months
- **AND** clicking a month summary SHALL expand to show dates

#### Scenario: Archive section with no previous months

- **WHEN** rendering date navigation with no archived years
- **THEN** the archive section SHALL display "No archives yet" placeholder text

#### Scenario: Page without date anchors

- **WHEN** rendering a page with empty DateAnchors
- **THEN** the rendered output SHALL NOT include a date navigation element

#### Scenario: Date anchor links

- **WHEN** rendering date navigation with date "2026-01-26"
- **THEN** the rendered output SHALL include an anchor element with href "#2026-01-26"
- **AND** the anchor text SHALL be "2026-01-26"

#### Scenario: Responsive layout

- **WHEN** rendering date navigation on desktop viewport (>768px)
- **THEN** current month and archive sections SHALL display side-by-side

#### Scenario: Mobile layout

- **WHEN** rendering date navigation on mobile viewport (≤768px)
- **THEN** current month section SHALL appear above archive section
- **AND** both sections SHALL be full-width

#### Scenario: No JavaScript required

- **WHEN** rendering date navigation
- **THEN** the navigation SHALL function without JavaScript
- **AND** the navigation SHALL use native HTML5 details/summary elements
- **AND** nested year/month collapse SHALL work without JavaScript

#### Scenario: Nested details styling

- **WHEN** rendering nested month details within year details
- **THEN** the nested details SHALL have left indentation for visual hierarchy
