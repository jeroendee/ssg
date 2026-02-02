## ADDED Requirements

### Requirement: Render footer content

The system SHALL render optional markdown content in the site footer.

#### Scenario: Footer with content

- **WHEN** rendering a page and site has footer content
- **THEN** the rendered footer SHALL include a horizontal rule (`<hr>`) at the top
- **AND** the rendered footer SHALL include the footer content HTML after the horizontal rule
- **AND** the rendered footer SHALL include the "Made with ❤️" tagline after the footer content

#### Scenario: Footer without content

- **WHEN** rendering a page and site has no footer content
- **THEN** the rendered footer SHALL NOT include a horizontal rule
- **AND** the rendered footer SHALL only include the "Made with ❤️" tagline

#### Scenario: Footer content HTML rendering

- **WHEN** rendering footer content containing HTML (links, formatting)
- **THEN** the rendered output SHALL preserve the HTML markup
- **AND** the HTML SHALL NOT be escaped

### Requirement: Footer visual structure

The system SHALL render the footer with clear visual hierarchy.

#### Scenario: Footer structure order

- **WHEN** rendering a footer with content
- **THEN** the elements SHALL appear in order: horizontal rule, footer content, tagline
- **AND** the footer content SHALL be wrapped in a container element
