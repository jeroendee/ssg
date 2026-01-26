# html-rendering Spec Delta

## ADDED Requirements

### Requirement: Render Date Navigation

The system SHALL render a date navigation dropdown for pages with date anchors.

#### Scenario: Page with date anchors

- **WHEN** rendering a page with DateAnchors containing dates
- **THEN** the rendered output SHALL include a date navigation element
- **AND** the navigation SHALL appear after the page title
- **AND** the navigation SHALL appear before the page content

#### Scenario: Page without date anchors

- **WHEN** rendering a page with empty DateAnchors
- **THEN** the rendered output SHALL NOT include a date navigation element

#### Scenario: Date navigation structure

- **WHEN** rendering a page with date anchors
- **THEN** the rendered output SHALL include a details element with class "date-nav"
- **AND** the details element SHALL contain a summary element with text "Jump to date"
- **AND** the details element SHALL contain an unordered list of anchor links

#### Scenario: Date anchor links

- **WHEN** rendering date navigation with date "2026-01-26"
- **THEN** the rendered output SHALL include an anchor element with href "#2026-01-26"
- **AND** the anchor text SHALL be "2026-01-26"

#### Scenario: Multiple date links

- **WHEN** rendering date navigation with multiple dates
- **THEN** each date SHALL be rendered as a list item with an anchor link
- **AND** dates SHALL appear in the same order as in DateAnchors

#### Scenario: Date navigation styling

- **WHEN** rendering date navigation
- **THEN** the navigation element SHALL use CSS class "date-nav"
- **AND** the navigation SHALL be styled using existing Solarized color variables

#### Scenario: No JavaScript required

- **WHEN** rendering date navigation
- **THEN** the navigation SHALL function without JavaScript
- **AND** the navigation SHALL use native HTML5 details/summary elements
