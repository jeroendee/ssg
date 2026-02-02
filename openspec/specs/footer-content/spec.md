# footer-content Specification

## Purpose

The footer content capability enables site-wide footer markdown content loaded from a partial file.

## Requirements

### Requirement: Load footer partial

The system SHALL load and render markdown content from `_footer.md` when present in the content directory.

#### Scenario: Footer file exists

- **WHEN** the content directory contains `_footer.md`
- **THEN** the system SHALL parse the file as markdown
- **AND** the system SHALL convert the markdown to HTML
- **AND** the system SHALL store the rendered HTML in the site's footer content

#### Scenario: Footer file does not exist

- **WHEN** the content directory does not contain `_footer.md`
- **THEN** the system SHALL set footer content to empty string
- **AND** the system SHALL NOT return an error

#### Scenario: Footer file parse error

- **WHEN** `_footer.md` exists but cannot be parsed
- **THEN** the system SHALL return an error indicating the parse failure

### Requirement: Footer content is site-wide

The system SHALL make footer content available to all rendered pages.

#### Scenario: Footer content in page templates

- **WHEN** rendering any page type (static page, blog post, blog list, 404)
- **THEN** the template SHALL have access to the site's footer content

### Requirement: Footer partial excluded from pages

The system SHALL NOT generate a page for `_footer.md`.

#### Scenario: Footer file not in site pages

- **WHEN** the content directory contains `_footer.md`
- **THEN** the discovered site pages SHALL NOT include `_footer.md`
- **AND** no HTML page SHALL be generated for `_footer.md`
