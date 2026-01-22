# feed-generation Specification

## Purpose
Enable content syndication by generating standardized RSS feed documents that allow users to subscribe to site updates through feed readers and aggregators.

## Requirements

### Requirement: Generate RSS Feed
The system SHALL generate a valid RSS 2.0 feed document containing recent blog posts.

#### Scenario: Valid RSS structure
- **WHEN** generating a feed with one or more posts
- **THEN** the output SHALL include an XML declaration with UTF-8 encoding
- **AND** the output SHALL include an RSS root element with version "2.0"
- **AND** the output SHALL include a channel element

#### Scenario: Channel metadata
- **WHEN** generating a feed
- **THEN** the channel SHALL include the site title
- **AND** the channel SHALL include the site base URL as the link
- **AND** the channel SHALL include the site description
- **AND** the channel SHALL include the last build date

### Requirement: Generate Feed Entries
The system SHALL generate an entry for each blog post in the feed.

#### Scenario: Entry contains required fields
- **WHEN** generating a feed entry for a post
- **THEN** the entry SHALL include the post title
- **AND** the entry SHALL include an absolute URL to the post
- **AND** the entry SHALL include the post content
- **AND** the entry SHALL include the publication date

#### Scenario: Absolute URLs in entries
- **WHEN** generating entry links
- **THEN** the link SHALL be an absolute URL constructed from the base URL and post slug
- **AND** the URL format SHALL be "{baseURL}/blog/{slug}/"

### Requirement: Format Publication Dates
The system SHALL format publication dates according to the RFC 822 standard.

#### Scenario: RFC 822 date format
- **WHEN** a post has a publication date
- **THEN** the pubDate element SHALL use RFC 822 format (e.g., "Mon, 15 Jan 2024 10:30:00 +0000")

### Requirement: Handle HTML Content
The system SHALL preserve HTML content in feed entries without escaping.

#### Scenario: CDATA wrapping for HTML content
- **WHEN** post content contains HTML markup
- **THEN** the description element SHALL wrap the content in a CDATA section
- **AND** the HTML tags SHALL be preserved unescaped

### Requirement: Limit Feed Size
The system SHALL limit the number of entries in the feed to maintain reasonable file sizes.

#### Scenario: Maximum 20 entries
- **WHEN** the site has more than 20 posts
- **THEN** the feed SHALL include only the first 20 posts
- **AND** older posts SHALL be excluded from the feed

### Requirement: Handle Empty Post List
The system SHALL gracefully handle sites with no posts.

#### Scenario: No posts available
- **WHEN** there are no posts to include
- **THEN** the system SHALL return an empty result
- **AND** no feed file SHALL be written

### Requirement: Write Feed File
The system SHALL write the generated feed to a standard location.

#### Scenario: Feed file location
- **WHEN** a feed is generated successfully
- **THEN** the feed SHALL be written to "feed.xml" in the output directory

#### Scenario: Skip writing when no posts
- **WHEN** the feed content is empty (no posts)
- **THEN** no feed file SHALL be written to the output directory
