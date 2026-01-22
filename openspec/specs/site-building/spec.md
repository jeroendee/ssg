# site-building Specification

## Purpose

The system orchestrates the complete build process for a static site, transforming content files into a structured output directory containing rendered pages, blog posts, feeds, and supporting files.

## Requirements

### Requirement: Build-Orchestration

The system SHALL execute a complete site build by scanning content, rendering output, and copying assets.

#### Scenario: Successful build

- **WHEN** the build process is initiated with valid configuration
- **THEN** the system SHALL create the output directory
- **AND** the system SHALL render all pages and posts to the output directory
- **AND** the system SHALL copy static assets to the output directory

#### Scenario: Content directory does not exist

- **WHEN** the build process is initiated
- **AND** the content directory does not exist
- **THEN** the system SHALL return an error

### Requirement: Output-Directory-Management

The system SHALL manage the output directory to ensure clean builds.

#### Scenario: Output directory does not exist

- **WHEN** the build process is initiated
- **AND** the output directory does not exist
- **THEN** the system SHALL create the output directory

#### Scenario: Output directory contains stale files

- **WHEN** the build process is initiated
- **AND** the output directory contains files from previous builds
- **THEN** the system SHALL remove all existing files before generating new output

### Requirement: Output-Directory-Validation

The system SHALL validate the output directory path to prevent catastrophic data loss.

#### Scenario: Empty output directory path

- **WHEN** the output directory path is empty
- **THEN** the system SHALL return an error indicating the path cannot be empty

#### Scenario: Output directory is root

- **WHEN** the output directory resolves to the filesystem root
- **THEN** the system SHALL return an error

#### Scenario: Output directory is home directory

- **WHEN** the output directory resolves to the user home directory
- **THEN** the system SHALL return an error

#### Scenario: Output directory is current directory

- **WHEN** the output directory is specified as the current directory
- **THEN** the system SHALL return an error

#### Scenario: Output directory is parent directory

- **WHEN** the output directory traverses to a parent directory
- **THEN** the system SHALL return an error

#### Scenario: Output directory is project root

- **WHEN** the output directory equals the project root
- **THEN** the system SHALL return an error

#### Scenario: Output directory is outside project root

- **WHEN** the output directory is located outside the project root
- **THEN** the system SHALL return an error

#### Scenario: Valid output directory

- **WHEN** the output directory is a subdirectory within the project root
- **THEN** the system SHALL accept the path

### Requirement: Homepage-Generation

The system SHALL generate a homepage from designated content.

#### Scenario: Homepage content exists

- **WHEN** the content directory contains the homepage file
- **THEN** the system SHALL render the homepage to the root index file

#### Scenario: Homepage content missing

- **WHEN** the content directory does not contain the homepage file
- **THEN** the system SHALL return an error indicating the homepage is required

#### Scenario: Homepage output location

- **WHEN** the homepage is rendered
- **THEN** the system SHALL write it directly to the output root
- **AND** the system SHALL NOT create a subdirectory for the homepage

### Requirement: Page-Rendering

The system SHALL render content pages with clean URLs.

#### Scenario: Page output structure

- **WHEN** a page is rendered
- **THEN** the system SHALL create a directory named after the page slug
- **AND** the system SHALL write the rendered content as an index file within that directory

#### Scenario: Pages only

- **WHEN** the content directory contains only pages
- **THEN** the system SHALL render all pages

### Requirement: Post-Rendering

The system SHALL render blog posts with clean URLs.

#### Scenario: Post output structure

- **WHEN** a post is rendered
- **THEN** the system SHALL create a directory under the blog path named after the post slug
- **AND** the system SHALL write the rendered content as an index file within that directory

#### Scenario: Posts only

- **WHEN** the content directory contains only posts
- **THEN** the system SHALL render all posts

### Requirement: Content-Scanning

The system SHALL discover and organize content files.

#### Scenario: Mixed content types

- **WHEN** the content directory contains pages and posts
- **THEN** the system SHALL scan and categorize all content

#### Scenario: Non-content files ignored

- **WHEN** the content directory contains non-markdown files
- **THEN** the system SHALL ignore those files

#### Scenario: Index files ignored

- **WHEN** the content directory contains index files with underscore prefix
- **THEN** the system SHALL ignore those files

#### Scenario: Posts sorted by date

- **WHEN** posts are scanned
- **THEN** the system SHALL sort posts by date in descending order (newest first)

### Requirement: Site-Metadata-Transfer

The system SHALL transfer configuration metadata to the site representation.

#### Scenario: Description from configuration

- **WHEN** the configuration includes a description
- **THEN** the site representation SHALL use that description

#### Scenario: Description fallback to title

- **WHEN** the configuration omits the description
- **THEN** the site representation SHALL use the title as the description

### Requirement: Blog-Listing-Generation

The system SHALL generate a blog listing page when posts exist.

#### Scenario: Posts exist

- **WHEN** the content includes blog posts
- **THEN** the system SHALL generate a blog listing page
- **AND** the listing SHALL include all posts

#### Scenario: No posts exist

- **WHEN** the content contains no blog posts
- **THEN** the system SHALL NOT generate a blog listing page

### Requirement: RSS-Feed-Generation

The system SHALL generate an RSS feed for blog content.

#### Scenario: Posts exist

- **WHEN** the content includes blog posts
- **THEN** the system SHALL generate an RSS feed file
- **AND** the feed SHALL include post content

#### Scenario: No posts exist

- **WHEN** the content contains no blog posts
- **THEN** the system SHALL NOT generate an RSS feed file

### Requirement: Error-Page-Generation

The system SHALL generate a 404 error page.

#### Scenario: Build completion

- **WHEN** the build completes successfully
- **THEN** the system SHALL generate a 404 error page

### Requirement: Sitemap-Generation

The system SHALL generate an XML sitemap.

#### Scenario: Pages included

- **WHEN** the build completes
- **THEN** the sitemap SHALL include all page URLs

#### Scenario: Posts included with modification dates

- **WHEN** the build completes with blog posts
- **THEN** the sitemap SHALL include all post URLs
- **AND** the sitemap SHALL include the last modification date for each post

#### Scenario: Homepage URL format

- **WHEN** the homepage is included in the sitemap
- **THEN** the URL SHALL be the base URL with a trailing slash
- **AND** the URL SHALL NOT contain double slashes in the path

#### Scenario: Sitemap with no posts

- **WHEN** the content contains only pages
- **THEN** the sitemap SHALL include only page URLs

### Requirement: Robots-File-Generation

The system SHALL generate a robots.txt file.

#### Scenario: Build completion

- **WHEN** the build completes
- **THEN** the system SHALL generate a robots.txt file
- **AND** the file SHALL allow all user agents
- **AND** the file SHALL reference the sitemap URL

### Requirement: Static-Asset-Handling

The system SHALL copy static assets to the output directory.

#### Scenario: Assets directory specified

- **WHEN** a static assets directory is configured
- **THEN** the system SHALL copy all files from that directory to the output root

#### Scenario: Assets directory not specified

- **WHEN** no static assets directory is configured
- **THEN** the system SHALL write the default stylesheet to the output

#### Scenario: Custom stylesheet provided

- **WHEN** the assets directory contains a stylesheet
- **THEN** the system SHALL use the custom stylesheet

#### Scenario: Other assets with default stylesheet

- **WHEN** the assets directory contains assets but no stylesheet
- **THEN** the system SHALL copy all provided assets
- **AND** the system SHALL write the default stylesheet

### Requirement: Post-Asset-Handling

The system SHALL copy assets referenced in blog posts.

#### Scenario: Post references assets

- **WHEN** a post references asset files
- **THEN** the system SHALL copy those assets to the post output directory

#### Scenario: Post references missing asset

- **WHEN** a post references an asset that does not exist
- **THEN** the system SHALL return an error identifying the missing asset and the post

#### Scenario: Post references no assets

- **WHEN** a post does not reference any assets
- **THEN** the system SHALL proceed without copying assets

#### Scenario: Post references multiple assets

- **WHEN** a post references multiple assets
- **THEN** the system SHALL copy all referenced assets

### Requirement: Asset-Path-Rewriting

The system SHALL rewrite asset paths in rendered post content.

#### Scenario: Double-quoted asset path

- **WHEN** rendered content contains a double-quoted asset path with prefix
- **THEN** the system SHALL remove the prefix from the path

#### Scenario: Single-quoted asset path

- **WHEN** rendered content contains a single-quoted asset path with prefix
- **THEN** the system SHALL remove the prefix from the path

#### Scenario: Multiple asset paths

- **WHEN** rendered content contains multiple asset paths
- **THEN** the system SHALL rewrite all paths

#### Scenario: No asset prefix

- **WHEN** rendered content contains paths without the asset prefix
- **THEN** the system SHALL leave those paths unchanged

#### Scenario: External URLs

- **WHEN** rendered content contains external URLs with asset-like paths
- **THEN** the system SHALL NOT modify those URLs

### Requirement: Version-Tracking

The system SHALL support build version tracking.

#### Scenario: Version specified

- **WHEN** a version is set before building
- **THEN** the rendered output SHALL include the version string
