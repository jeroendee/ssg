# feed-pages Specification

## Purpose
Enable pages with date-anchored content sections to generate RSS feed entries, allowing subscribers to receive notifications when new daily content is added.

## Requirements

### Requirement: Parse Date Sections from Page Content
The system SHALL extract date sections from configured feed pages based on date-anchored headers.

#### Scenario: Page with date-anchored headers
- **WHEN** a feed page contains headers with date anchors in format `{#YYYY-MM-DD}`
- **THEN** the system SHALL extract each date section as a separate feed item
- **AND** each feed item SHALL contain the content between the header and the next date header (or end of content)

#### Scenario: Page without date anchors
- **WHEN** a feed page contains no date-anchored headers
- **THEN** the system SHALL produce no feed items for that page

#### Scenario: Multiple dates on same page
- **WHEN** a feed page contains multiple date-anchored headers
- **THEN** the system SHALL produce one feed item per date
- **AND** each feed item SHALL be associated with its respective date

### Requirement: Generate Feed Item Metadata
The system SHALL generate appropriate metadata for each date section feed item.

#### Scenario: Feed item title
- **WHEN** generating a feed item for a date section
- **THEN** the title SHALL be the page title followed by the formatted date
- **AND** the format SHALL be "{PageTitle} - {Month} {Day}, {Year}"

#### Scenario: Feed item link
- **WHEN** generating a feed item for a date section
- **THEN** the link SHALL be the absolute URL to the page with the date anchor
- **AND** the format SHALL be "{baseURL}{pagePath}#{date}"

#### Scenario: Feed item GUID
- **WHEN** generating a feed item for a date section
- **THEN** the GUID SHALL be the same as the link (unique per page and date)

#### Scenario: Feed item date
- **WHEN** generating a feed item for a date section
- **THEN** the publication date SHALL be derived from the date anchor
- **AND** the time component SHALL be set to midnight UTC

### Requirement: Include Date Section Content
The system SHALL include the full content of each date section in the feed entry.

#### Scenario: Content extraction
- **WHEN** generating a feed item for a date section
- **THEN** the content SHALL include all text between the date header and the next date header
- **AND** the content SHALL exclude the header itself
- **AND** HTML markup SHALL be preserved
