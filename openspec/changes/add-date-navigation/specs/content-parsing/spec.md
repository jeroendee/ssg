# content-parsing Spec Delta

## ADDED Requirements

### Requirement: Extract-Date-Anchors

The system SHALL extract date-based heading anchors from markdown content.

#### Scenario: Page with date headings

- **WHEN** a page content contains headings (h1-h6) with dates in format `# *YYYY-MM-DD*` through `###### *YYYY-MM-DD*`
- **THEN** the parsed page object SHALL contain the extracted dates in DateAnchors field
- **AND** the dates SHALL be in the order they appear in the content

#### Scenario: Multiple date headings at various levels

- **WHEN** a page content contains multiple date headings at any level (h1-h6)
- **THEN** the DateAnchors field SHALL contain all dates in document order

#### Scenario: No date headings

- **WHEN** a page content contains no date headings at any level
- **THEN** the DateAnchors field SHALL be empty

#### Scenario: Date heading format validation

- **WHEN** a page content contains `## *2026-01-26*`
- **THEN** the extracted date SHALL be "2026-01-26"

#### Scenario: Non-date headings excluded

- **WHEN** a page content contains headings without date format
- **THEN** those headings SHALL NOT be included in DateAnchors

#### Scenario: Extraction from markdown source

- **WHEN** extracting date anchors
- **THEN** the system SHALL extract from the raw markdown before HTML conversion

#### Scenario: All heading levels supported

- **WHEN** a page contains date headings at h1, h2, h3, h4, h5, and h6 levels
- **THEN** all dates SHALL be extracted regardless of heading level
