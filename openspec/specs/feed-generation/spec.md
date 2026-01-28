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

### Requirement: Write Feed File
The system SHALL write the generated feed to a standard location.

#### Scenario: Feed file location
- **WHEN** a feed is generated successfully
- **THEN** the feed SHALL be written to "feed.xml" in the output directory

#### Scenario: Skip writing when no posts
- **WHEN** the feed content is empty (no posts)
- **THEN** no feed file SHALL be written to the output directory
