# content-parsing Specification Delta

## MODIFIED Requirements

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
