# feed-generation Delta Specification

## ADDED Requirements

### Requirement: Accept Mixed Feed Items
The system SHALL generate feed entries for both blog posts and page date sections.

#### Scenario: Feed with posts and page sections
- **WHEN** generating a feed with blog posts and configured feed pages with date sections
- **THEN** the feed SHALL include entries for both posts and page date sections

#### Scenario: Feed with only posts
- **WHEN** generating a feed with blog posts and no configured feed pages
- **THEN** the feed SHALL include only post entries (existing behavior)

#### Scenario: Feed with only page sections
- **WHEN** generating a feed with no posts but configured feed pages with date sections
- **THEN** the feed SHALL include only page date section entries

### Requirement: Sort Mixed Feed by Date
The system SHALL sort all feed entries by date regardless of source.

#### Scenario: Interleaved dates
- **WHEN** generating a feed with posts and page sections having overlapping dates
- **THEN** all entries SHALL be sorted by date in descending order (newest first)
- **AND** entries from different sources SHALL be interleaved based on their dates

## MODIFIED Requirements

### Requirement: Limit Feed Size
The system SHALL limit the number of entries in the feed to maintain reasonable file sizes.

#### Scenario: Maximum 20 entries
- **WHEN** the combined total of posts and page sections exceeds 20
- **THEN** the feed SHALL include only the 20 most recent entries
- **AND** older entries SHALL be excluded regardless of source

### Requirement: Handle Empty Post List
The system SHALL gracefully handle sites with no feed content.

#### Scenario: No posts or page sections available
- **WHEN** there are no posts and no page sections to include
- **THEN** the system SHALL return an empty result
- **AND** no feed file SHALL be written

#### Scenario: No posts but page sections available
- **WHEN** there are no posts but page sections exist
- **THEN** the system SHALL generate a feed with the page section entries
