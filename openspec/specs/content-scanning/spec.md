# content-scanning Specification

## Purpose
The system discovers and organizes content files from the content directory, distinguishing between pages and posts while ensuring required content exists and proper ordering of results.

## Requirements

### Requirement: Verify Content Directory
The system SHALL verify the content directory exists before scanning.

#### Scenario: Content directory does not exist
- **WHEN** the content directory path does not exist
- **THEN** the system SHALL return an error

### Requirement: Require Homepage
The system SHALL require a homepage file in the content directory.

#### Scenario: Homepage file missing
- **WHEN** the content directory does not contain a homepage file named "home.md"
- **THEN** the system SHALL return an error indicating the homepage is required

#### Scenario: Homepage file present
- **WHEN** the content directory contains "home.md"
- **THEN** the system SHALL proceed with content scanning

### Requirement: Discover Pages
The system SHALL discover page content files from the root content directory.

#### Scenario: Markdown files in root directory
- **WHEN** the content directory contains markdown files
- **THEN** the system SHALL include each markdown file as a page
- **AND** the homepage file SHALL be included as a page with empty identifier

#### Scenario: Non-markdown files present
- **WHEN** the content directory contains non-markdown files (text files, images, etc.)
- **THEN** the system SHALL ignore non-markdown files
- **AND** the system SHALL only return markdown files as pages

#### Scenario: Subdirectories in content root
- **WHEN** the content directory contains subdirectories
- **THEN** the system SHALL ignore subdirectories when scanning for pages

#### Scenario: Index files present
- **WHEN** the content directory contains index files (files prefixed with underscore)
- **THEN** the system SHALL ignore index files
- **AND** the system SHALL not return index files as pages

### Requirement: Discover Posts
The system SHALL discover post content files from the blog subdirectory.

#### Scenario: Blog directory exists with posts
- **WHEN** the content directory contains a "blog" subdirectory with markdown files
- **THEN** the system SHALL include each markdown file as a post

#### Scenario: Blog directory does not exist
- **WHEN** the content directory does not contain a "blog" subdirectory
- **THEN** the system SHALL return zero posts
- **AND** the system SHALL not return an error

#### Scenario: Blog directory with index files
- **WHEN** the blog directory contains index files (prefixed with underscore)
- **THEN** the system SHALL ignore index files in the blog directory

#### Scenario: Blog directory with non-markdown files
- **WHEN** the blog directory contains non-markdown files
- **THEN** the system SHALL ignore non-markdown files in the blog directory

### Requirement: Order Posts Chronologically
The system SHALL sort discovered posts by publication date.

#### Scenario: Multiple posts with different dates
- **WHEN** the blog directory contains posts with different publication dates
- **THEN** the system SHALL return posts sorted by date in descending order (newest first)

### Requirement: Propagate Site Metadata
The system SHALL propagate site metadata from configuration to the discovered site.

#### Scenario: Configuration with all metadata
- **WHEN** the configuration includes title, base URL, author, and navigation
- **THEN** the discovered site SHALL include all configured metadata

#### Scenario: Description configured
- **WHEN** the configuration includes a description
- **THEN** the discovered site description SHALL match the configured description

#### Scenario: Description not configured
- **WHEN** the configuration omits description
- **THEN** the discovered site description SHALL default to the site title

### Requirement: Handle Content Parsing Errors
The system SHALL propagate errors encountered during content parsing.

#### Scenario: Invalid page content
- **WHEN** a page file cannot be parsed
- **THEN** the system SHALL return an error

#### Scenario: Invalid post content
- **WHEN** a post file cannot be parsed
- **THEN** the system SHALL return an error
