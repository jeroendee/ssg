# content-parsing Spec Delta

## ADDED Requirements

### Requirement: Extract-Date-Anchors

The system SHALL extract date-based heading anchors from markdown content.

#### Scenario: Page with date headings

- **WHEN** a page content contains h4 headings with dates in format `#### *YYYY-MM-DD*`
- **THEN** the parsed page object SHALL contain the extracted dates in DateAnchors field
- **AND** the dates SHALL be in the order they appear in the content

#### Scenario: Multiple date headings

- **WHEN** a page content contains multiple h4 date headings
- **THEN** the DateAnchors field SHALL contain all dates in document order

#### Scenario: No date headings

- **WHEN** a page content contains no h4 date headings
- **THEN** the DateAnchors field SHALL be empty

#### Scenario: Date heading format validation

- **WHEN** a page content contains `#### *2026-01-26*`
- **THEN** the extracted date SHALL be "2026-01-26"

#### Scenario: Non-date h4 headings excluded

- **WHEN** a page content contains h4 headings without date format
- **THEN** those headings SHALL NOT be included in DateAnchors

#### Scenario: Extraction from markdown source

- **WHEN** extracting date anchors
- **THEN** the system SHALL extract from the raw markdown before HTML conversion
