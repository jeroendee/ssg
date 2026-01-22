# content-parsing Specification

## Purpose

The system parses content files to extract structured metadata and transform markup into rendered output for site generation.

## Requirements

### Requirement: Parse-Page

The system SHALL parse page content files and produce page content objects.

#### Scenario: Page with frontmatter

- **WHEN** a page content file contains frontmatter with a title field
- **THEN** the parsed page object SHALL contain the extracted title
- **AND** the content body SHALL be transformed to rendered markup

#### Scenario: Page slug derivation

- **WHEN** a page content file is named "about.md"
- **THEN** the parsed page object SHALL have slug "about"
- **AND** the parsed page object SHALL have path "/about/"

#### Scenario: Home page special handling

- **WHEN** a page content file is named "home.md"
- **THEN** the parsed page object SHALL have an empty slug
- **AND** the parsed page object SHALL have path "/"

#### Scenario: Non-existent page file

- **WHEN** the specified page file does not exist
- **THEN** the system SHALL return an error

### Requirement: Parse-Post

The system SHALL parse post content files and produce post content objects with additional metadata.

#### Scenario: Post with frontmatter date

- **WHEN** a post content file contains frontmatter with date and summary fields
- **THEN** the parsed post object SHALL contain the extracted date
- **AND** the parsed post object SHALL contain the extracted summary

#### Scenario: Post date from filename

- **WHEN** a post content file is named with date prefix (YYYY-MM-DD-slug.md format)
- **THEN** the parsed post object SHALL extract the date from the filename
- **AND** the parsed post object SHALL derive the slug from the portion after the date

#### Scenario: Frontmatter date overrides filename date

- **WHEN** a post content file has both a date-prefixed filename and a frontmatter date
- **THEN** the frontmatter date SHALL take precedence

#### Scenario: Invalid frontmatter date

- **WHEN** a post content file contains an invalid date in frontmatter
- **THEN** the system SHALL return an error

#### Scenario: Invalid filename date

- **WHEN** a post content file has an invalid date in the filename prefix
- **THEN** the system SHALL return an error

#### Scenario: Non-existent post file

- **WHEN** the specified post file does not exist
- **THEN** the system SHALL return an error

#### Scenario: Post path generation

- **WHEN** a post is parsed with slug "my-post"
- **THEN** the parsed post object SHALL have path "/blog/my-post/"

### Requirement: Calculate-Word-Count

The system SHALL calculate word count from content body.

#### Scenario: Word count calculation

- **WHEN** a post content body contains words
- **THEN** the parsed post object SHALL include the word count

### Requirement: Transform-Markup

The system SHALL transform markup content to rendered output format.

#### Scenario: Heading transformation

- **WHEN** content contains a level-one heading
- **THEN** the rendered output SHALL contain the appropriate heading element

#### Scenario: Text formatting transformation

- **WHEN** content contains bold text markers
- **THEN** the rendered output SHALL contain emphasized text elements

### Requirement: Extract-Frontmatter

The system SHALL extract structured metadata from content file frontmatter.

#### Scenario: Frontmatter present

- **WHEN** content begins with frontmatter delimiters
- **THEN** the system SHALL parse the metadata section
- **AND** the system SHALL return the content body without frontmatter

#### Scenario: No frontmatter

- **WHEN** content does not begin with frontmatter delimiters
- **THEN** the system SHALL return the full content as body
- **AND** metadata fields SHALL be empty

### Requirement: Extract-Asset-References

The system SHALL identify local asset references within content.

#### Scenario: Single asset reference

- **WHEN** content contains an image reference to a local asset path
- **THEN** the system SHALL return the asset path

#### Scenario: Multiple asset references

- **WHEN** content contains multiple local asset references
- **THEN** the system SHALL return all asset paths in order

#### Scenario: No asset references

- **WHEN** content contains no local asset references
- **THEN** the system SHALL return an empty collection

#### Scenario: External URLs excluded

- **WHEN** content contains image references to external URLs
- **THEN** the system SHALL NOT include external URLs in asset references

#### Scenario: Mixed internal and external references

- **WHEN** content contains both local asset references and external URLs
- **THEN** the system SHALL return only the local asset paths

#### Scenario: Nested asset paths

- **WHEN** content references assets in subdirectories
- **THEN** the system SHALL preserve the full relative path

#### Scenario: Empty content

- **WHEN** content is empty
- **THEN** the system SHALL return an empty collection

#### Scenario: Post includes extracted assets

- **WHEN** a post is parsed
- **THEN** the parsed post object SHALL include all extracted asset references
