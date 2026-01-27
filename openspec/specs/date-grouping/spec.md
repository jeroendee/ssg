# date-grouping Specification

## Purpose

The date grouping capability organizes date anchors into hierarchical structures (year > month > date) for navigation and archive display in static site generation.
## Requirements
### Requirement: Group dates by month

The system SHALL group date anchors by year and month.

#### Scenario: Multiple months of dates

- **WHEN** grouping dates ["2026-02-03", "2026-02-01", "2026-01-31", "2026-01-15", "2025-12-20"]
- **THEN** the system SHALL return groups for 2026-February, 2026-January, and 2025-December
- **AND** each group SHALL contain only dates from that month
- **AND** dates within each group SHALL preserve their original order

#### Scenario: Single month of dates

- **WHEN** grouping dates ["2026-01-27", "2026-01-26", "2026-01-25"]
- **THEN** the system SHALL return one group for 2026-January
- **AND** that group SHALL contain all three dates

#### Scenario: Empty input

- **WHEN** grouping an empty date list
- **THEN** the system SHALL return empty results for both current and archived

### Requirement: Identify current month

The system SHALL identify the "current month" as the month containing the most recent date entry.

#### Scenario: Current month identification

- **WHEN** grouping dates where the most recent is "2026-01-27"
- **THEN** the system SHALL identify January 2026 as the current month
- **AND** the current month dates SHALL be returned separately from archived months

#### Scenario: Current month with future month entries

- **WHEN** grouping dates where February has entries but January also has entries
- **THEN** the system SHALL identify February as the current month (most recent date wins)

### Requirement: Order archived months

The system SHALL order archived months from newest to oldest.

#### Scenario: Archive ordering

- **WHEN** archiving months January 2026, December 2025, and November 2025
- **THEN** the archived list SHALL be ordered [January 2026, December 2025, November 2025]

### Requirement: Provide month display name

The system SHALL provide human-readable month names for each group.

#### Scenario: Month name format

- **WHEN** grouping dates from January 2026
- **THEN** the month group SHALL have Year=2026 and Month="January"

### Requirement: Handle malformed dates

The system SHALL gracefully handle malformed date strings.

#### Scenario: Invalid date format

- **WHEN** grouping dates containing "not-a-date"
- **THEN** the system SHALL exclude the malformed entry from all groups
- **AND** the system SHALL continue processing valid dates

### Requirement: Group months by year

The system SHALL group archived months under year headers.

#### Scenario: Multiple years of archived months

- **WHEN** archiving months January 2026, December 2025, and November 2025
- **THEN** the system SHALL group into two year groups: 2026 and 2025
- **AND** year 2026 SHALL contain [January]
- **AND** year 2025 SHALL contain [December, November] in that order

#### Scenario: Single year of archived months

- **WHEN** archiving months December 2025, November 2025, October 2025
- **THEN** the system SHALL return one year group for 2025
- **AND** that year group SHALL contain all three months in newest-first order

### Requirement: Order year groups

The system SHALL order year groups from newest to oldest.

#### Scenario: Year ordering

- **WHEN** archiving months from years 2026, 2025, and 2024
- **THEN** the year groups SHALL be ordered [2026, 2025, 2024]

### Requirement: Provide YearGroup structure

The system SHALL provide year groups with numeric year and ordered months.

#### Scenario: YearGroup fields

- **WHEN** grouping months from 2025
- **THEN** the year group SHALL have Year=2025
- **AND** the year group SHALL have Months as a slice of MonthGroup

