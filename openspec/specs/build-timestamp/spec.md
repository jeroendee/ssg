# build-timestamp Specification

## Purpose
TBD - created by archiving change add-build-timestamp. Update Purpose after archive.
## Requirements
### Requirement: Build-Timestamp-Generation

The system SHALL generate a build timestamp file for sync verification.

#### Scenario: Build timestamp created

- **WHEN** the build completes successfully
- **THEN** the system SHALL generate a build timestamp file in the output root
- **AND** the file SHALL be named `build.json`

#### Scenario: Timestamp format

- **WHEN** the build timestamp file is generated
- **THEN** the file SHALL contain valid JSON
- **AND** the JSON SHALL have a `buildTime` field
- **AND** the value SHALL be in RFC 3339 format (e.g., `2026-01-25T14:30:00Z`)
- **AND** the timestamp SHALL be in UTC timezone

#### Scenario: Timestamp reflects build time

- **WHEN** the build timestamp file is generated
- **THEN** the `buildTime` value SHALL reflect the actual time the build was executed
- **AND** the timestamp SHALL be generated fresh on each build

#### Scenario: File encoding

- **WHEN** the build timestamp file is generated
- **THEN** the file SHALL be UTF-8 encoded

