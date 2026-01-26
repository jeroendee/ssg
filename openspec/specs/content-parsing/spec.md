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

#### Scenario: Heading anchor ID generation

- **WHEN** content contains a heading at any level (h1-h6)
- **THEN** the rendered output SHALL include an id attribute on the heading element
- **AND** the id value SHALL be derived from the heading text
- **AND** the id value SHALL be lowercase
- **AND** spaces in the heading text SHALL be replaced with hyphens in the id
- **AND** special characters SHALL be removed from the id

#### Scenario: Heading level one with ID

- **WHEN** content contains `# Hello World`
- **THEN** the rendered output SHALL contain `<h1 id="hello-world">Hello World</h1>`

#### Scenario: Heading level two with ID

- **WHEN** content contains `## Features`
- **THEN** the rendered output SHALL contain `<h2 id="features">Features</h2>`

#### Scenario: Heading with special characters

- **WHEN** content contains `### What's New?`
- **THEN** the rendered output SHALL contain an id without special characters

#### Scenario: All heading levels supported

- **WHEN** content contains headings at levels 1 through 6
- **THEN** each heading element SHALL have an id attribute

#### Scenario: Heading content wrapped in anchor link

- **WHEN** content contains a heading at any level (h1-h6)
- **THEN** the heading content SHALL be wrapped in an anchor element
- **AND** the anchor href SHALL reference the heading id with a hash prefix

#### Scenario: Clickable heading level one

- **WHEN** content contains `# Hello World`
- **THEN** the rendered output SHALL contain `<h1 id="hello-world"><a href="#hello-world">Hello World</a></h1>`

#### Scenario: Clickable heading level two

- **WHEN** content contains `## Features`
- **THEN** the rendered output SHALL contain `<h2 id="features"><a href="#features">Features</a></h2>`

#### Scenario: URL updates on heading click

- **WHEN** a user clicks on a heading
- **THEN** the browser URL SHALL update to include the heading fragment
- **AND** the page position SHALL scroll to the heading if not already visible

#### Scenario: Copy link functionality

- **WHEN** a user right-clicks a heading and selects "Copy link"
- **THEN** the copied URL SHALL include the heading fragment

#### Scenario: Anchor link default styling

- **WHEN** a heading anchor link is rendered
- **THEN** the link color SHALL inherit from the heading text color
- **AND** the link SHALL NOT display an underline by default

#### Scenario: Anchor link hover styling

- **WHEN** a user hovers over a heading anchor link
- **THEN** the link SHALL display an underline

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

