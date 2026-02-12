## MODIFIED Requirements

### Requirement: Render topics bar

The system SHALL render a topics discovery bar on pages that have extracted topics.

#### Scenario: Page with topics

- **WHEN** rendering a page that has topics extracted
- **THEN** the rendered output SHALL include a topics container element
- **AND** the topics container SHALL appear between the date navigation and the content
- **AND** each topic SHALL be displayed as the word only, without a count
- **AND** topics SHALL be separated by commas and spaces
- **AND** topics SHALL be ordered by frequency descending (matching extraction order)

#### Scenario: Page without topics

- **WHEN** rendering a page that has no topics (empty Topics slice)
- **THEN** the rendered output SHALL NOT include a topics container element

#### Scenario: Topics bar styling

- **WHEN** rendering the topics bar
- **THEN** the topics container SHALL have a CSS class "topics"
- **AND** the text SHALL use a smaller font size than the main content
- **AND** the text color SHALL be muted relative to the main content
