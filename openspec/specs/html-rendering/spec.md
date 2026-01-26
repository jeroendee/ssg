# html-rendering Specification

## Purpose

The HTML rendering capability transforms site data and content into complete HTML documents with proper structure, metadata, and navigation for static site generation.
## Requirements
### Requirement: Render HTML Document

The system SHALL produce valid HTML documents with proper structure.

#### Scenario: Valid document structure

- **WHEN** rendering any page type
- **THEN** the rendered output SHALL include a DOCTYPE declaration
- **AND** the rendered output SHALL include an html element
- **AND** the rendered output SHALL include header and footer elements

### Requirement: Render Static Page

The system SHALL render static pages with title and content.

#### Scenario: Page with title and content

- **WHEN** rendering a page with title "About" and content containing text
- **THEN** the rendered output SHALL include the page title
- **AND** the rendered output SHALL include the page content

#### Scenario: Page with special characters

- **WHEN** rendering a page with content containing special characters
- **THEN** the rendered output SHALL preserve the special characters in the content

### Requirement: Render Blog List

The system SHALL render a blog listing page showing all posts.

#### Scenario: Blog list with posts

- **WHEN** rendering a blog list with multiple posts
- **THEN** the rendered output SHALL include each post title as a link
- **AND** the rendered output SHALL include each post date in YYYY-MM-DD format
- **AND** the rendered output SHALL include each post word count

#### Scenario: Empty blog list

- **WHEN** rendering a blog list with no posts
- **THEN** the rendered output SHALL be a valid HTML document

### Requirement: Render Blog Post

The system SHALL render individual blog posts with title, date, and content.

#### Scenario: Blog post with metadata

- **WHEN** rendering a blog post with title, date, and content
- **THEN** the rendered output SHALL include the post title
- **AND** the rendered output SHALL include the date in YYYY-MM-DD format
- **AND** the rendered output SHALL include the post content

#### Scenario: Blog post with HTML content

- **WHEN** rendering a blog post containing preformatted code blocks
- **THEN** the rendered output SHALL preserve the HTML markup in the content

#### Scenario: Blog post with word count

- **WHEN** rendering a blog post with a word count value
- **THEN** the rendered output SHALL display the word count

#### Scenario: Blog post with zero word count

- **WHEN** rendering a blog post with zero word count
- **THEN** the rendered output SHALL display "0 words"

### Requirement: Render Error Page

The system SHALL render a 404 error page.

#### Scenario: Not found page

- **WHEN** rendering the 404 error page
- **THEN** the rendered output SHALL include the site title
- **AND** the rendered output SHALL include a 404 indicator

### Requirement: Include Navigation

The system SHALL include site navigation in rendered pages.

#### Scenario: Site with navigation items

- **WHEN** rendering a page for a site with navigation items
- **THEN** the rendered output SHALL include each navigation item

#### Scenario: Site with empty navigation

- **WHEN** rendering a page for a site with no navigation items
- **THEN** the rendered output SHALL be a valid HTML document

### Requirement: Include Site Logo

The system SHALL conditionally include a logo image.

#### Scenario: Logo configured

- **WHEN** site configuration includes a logo path
- **THEN** the rendered output SHALL include an image element with the logo path
- **AND** the image element SHALL have an alt attribute containing the site title

#### Scenario: Logo not configured

- **WHEN** site configuration omits the logo
- **THEN** the rendered output SHALL NOT include a logo image element

### Requirement: Include Favicon

The system SHALL conditionally include a favicon link.

#### Scenario: Favicon configured

- **WHEN** site configuration includes a favicon path
- **THEN** the rendered output SHALL include a link element for the favicon

#### Scenario: Favicon not configured

- **WHEN** site configuration omits the favicon
- **THEN** the rendered output SHALL NOT include a favicon link element

### Requirement: Include Canonical URL

The system SHALL include canonical URLs in rendered pages.

#### Scenario: Static page canonical URL

- **WHEN** rendering a static page with slug "about" and base URL "https://example.com"
- **THEN** the rendered output SHALL include canonical URL "https://example.com/about/"

#### Scenario: Blog list canonical URL

- **WHEN** rendering the blog list with base URL "https://example.com"
- **THEN** the rendered output SHALL include canonical URL "https://example.com/blog/"

#### Scenario: Blog post canonical URL

- **WHEN** rendering a blog post with slug "my-post" and base URL "https://example.com"
- **THEN** the rendered output SHALL include canonical URL "https://example.com/blog/my-post/"

#### Scenario: Error page canonical URL

- **WHEN** rendering the 404 page with base URL "https://example.com"
- **THEN** the rendered output SHALL include canonical URL "https://example.com/404/"

### Requirement: Include Meta Description

The system SHALL include meta description tags.

#### Scenario: Blog post with summary

- **WHEN** rendering a blog post that has a summary
- **THEN** the rendered output SHALL include a meta description with the post summary

#### Scenario: Blog post without summary

- **WHEN** rendering a blog post that has no summary
- **THEN** the rendered output SHALL include a meta description with the site description

#### Scenario: Static page meta description

- **WHEN** rendering a static page
- **THEN** the rendered output SHALL include a meta description with the site description

### Requirement: Include Open Graph Tags

The system SHALL include Open Graph metadata for social sharing.

#### Scenario: Blog post Open Graph tags

- **WHEN** rendering a blog post
- **THEN** the rendered output SHALL include og:title with the post title
- **AND** the rendered output SHALL include og:description with the summary
- **AND** the rendered output SHALL include og:url with the canonical URL
- **AND** the rendered output SHALL include og:type with value "article"
- **AND** the rendered output SHALL include og:site_name with the site title

#### Scenario: Static page Open Graph tags

- **WHEN** rendering a static page
- **THEN** the rendered output SHALL include og:title with the page title
- **AND** the rendered output SHALL include og:type with value "website"
- **AND** the rendered output SHALL include og:site_name with the site title

### Requirement: Include Twitter Card Tags

The system SHALL include Twitter Card metadata.

#### Scenario: Blog post Twitter Card tags

- **WHEN** rendering a blog post
- **THEN** the rendered output SHALL include twitter:card with value "summary"
- **AND** the rendered output SHALL include twitter:title with the post title
- **AND** the rendered output SHALL include twitter:description with the summary

#### Scenario: Static page Twitter Card tags

- **WHEN** rendering a static page
- **THEN** the rendered output SHALL include twitter:card with value "summary"
- **AND** the rendered output SHALL include twitter:title with the page title

### Requirement: Include Social Image

The system SHALL include social sharing image tags when configured.

#### Scenario: OG image configured

- **WHEN** site configuration includes a dedicated social image path
- **THEN** the rendered output SHALL include og:image with the absolute URL to the social image
- **AND** the rendered output SHALL include twitter:image with the absolute URL to the social image

#### Scenario: Logo fallback for social image

- **WHEN** site configuration omits social image but includes a logo
- **THEN** the rendered output SHALL include og:image with the absolute URL to the logo
- **AND** the rendered output SHALL include twitter:image with the absolute URL to the logo

#### Scenario: No social image available

- **WHEN** site configuration omits both social image and logo
- **THEN** the rendered output SHALL NOT include og:image tag
- **AND** the rendered output SHALL NOT include twitter:image tag

### Requirement: Include Structured Data

The system SHALL include JSON-LD structured data.

#### Scenario: Blog post structured data

- **WHEN** rendering a blog post
- **THEN** the rendered output SHALL include a JSON-LD script element
- **AND** the structured data SHALL have type "Article"
- **AND** the structured data SHALL include the headline
- **AND** the structured data SHALL include the publication date
- **AND** the structured data SHALL include the author
- **AND** the structured data SHALL include the description

#### Scenario: Static page structured data

- **WHEN** rendering a static page
- **THEN** the rendered output SHALL include structured data with type "WebSite"
- **AND** the rendered output SHALL NOT include structured data with type "Article"

### Requirement: Include Version

The system SHALL include a version identifier when configured.

#### Scenario: Version configured

- **WHEN** a version string is configured
- **THEN** the rendered output for static pages SHALL include the version
- **AND** the rendered output for blog lists SHALL include the version
- **AND** the rendered output for blog posts SHALL include the version
- **AND** the rendered output for error pages SHALL include the version
- **AND** the rendered output for base template SHALL include the version

### Requirement: Include Analytics

The system SHALL conditionally include analytics scripts.

#### Scenario: Analytics configured

- **WHEN** site configuration includes an analytics code
- **THEN** the rendered output SHALL include the analytics script with the site code

#### Scenario: Analytics not configured

- **WHEN** site configuration omits analytics code
- **THEN** the rendered output SHALL NOT include analytics scripts

#### Scenario: Analytics configured with empty value

- **WHEN** site configuration includes analytics with empty value
- **THEN** the rendered output SHALL NOT include analytics scripts

### Requirement: Render RSS Feed

The system SHALL render an RSS 2.0 feed for blog posts.

#### Scenario: RSS feed structure

- **WHEN** rendering an RSS feed with posts
- **THEN** the rendered output SHALL include an XML declaration
- **AND** the rendered output SHALL include an RSS root element with version 2.0
- **AND** the rendered output SHALL include a channel element
- **AND** the channel SHALL include the site title
- **AND** the channel SHALL include the site description
- **AND** the channel SHALL include the site link

#### Scenario: RSS feed item URLs

- **WHEN** rendering an RSS feed with a post
- **THEN** each item link SHALL be an absolute URL

#### Scenario: RSS feed date format

- **WHEN** rendering an RSS feed with posts
- **THEN** each item publication date SHALL use RFC 822 format

#### Scenario: RSS feed HTML content

- **WHEN** rendering an RSS feed with posts containing HTML content
- **THEN** the item description SHALL wrap HTML content in CDATA

#### Scenario: RSS feed maximum items

- **WHEN** rendering an RSS feed with more than 20 posts
- **THEN** the rendered output SHALL include only 20 items

#### Scenario: Empty RSS feed

- **WHEN** rendering an RSS feed with no posts
- **THEN** the rendered output SHALL be empty

### Requirement: Render Date Navigation

The system SHALL render a date navigation dropdown for pages with date anchors.

#### Scenario: Page with date anchors

- **WHEN** rendering a page with DateAnchors containing dates
- **THEN** the rendered output SHALL include a date navigation element
- **AND** the navigation SHALL appear after the page title
- **AND** the navigation SHALL appear before the page content

#### Scenario: Page without date anchors

- **WHEN** rendering a page with empty DateAnchors
- **THEN** the rendered output SHALL NOT include a date navigation element

#### Scenario: Date navigation structure

- **WHEN** rendering a page with date anchors
- **THEN** the rendered output SHALL include a details element with class "date-nav"
- **AND** the details element SHALL contain a summary element with text "Jump to date"
- **AND** the details element SHALL contain an unordered list of anchor links

#### Scenario: Date anchor links

- **WHEN** rendering date navigation with date "2026-01-26"
- **THEN** the rendered output SHALL include an anchor element with href "#2026-01-26"
- **AND** the anchor text SHALL be "2026-01-26"

#### Scenario: Multiple date links

- **WHEN** rendering date navigation with multiple dates
- **THEN** each date SHALL be rendered as a list item with an anchor link
- **AND** dates SHALL appear in the same order as in DateAnchors

#### Scenario: Date navigation styling

- **WHEN** rendering date navigation
- **THEN** the navigation element SHALL use CSS class "date-nav"
- **AND** the navigation SHALL be styled using existing Solarized color variables

#### Scenario: No JavaScript required

- **WHEN** rendering date navigation
- **THEN** the navigation SHALL function without JavaScript
- **AND** the navigation SHALL use native HTML5 details/summary elements

