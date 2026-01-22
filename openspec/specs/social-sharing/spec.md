# social-sharing Specification

## Purpose
TBD - created by archiving change add-og-image-config. Update Purpose after archive.
## Requirements
### Requirement: OG Image Configuration
The system SHALL support a dedicated `ogImage` configuration field for social media sharing images.

#### Scenario: OG image specified
- **WHEN** site configuration includes `site.ogImage: /social-image.png`
- **THEN** the rendered HTML SHALL include `<meta property="og:image" content="{baseURL}/social-image.png">`
- **AND** the rendered HTML SHALL include `<meta name="twitter:image" content="{baseURL}/social-image.png">`

#### Scenario: OG image not specified with logo fallback
- **WHEN** site configuration omits `ogImage` but includes `logo`
- **THEN** the rendered HTML SHALL use the logo path for og:image

#### Scenario: Neither OG image nor logo specified
- **WHEN** site configuration omits both `ogImage` and `logo`
- **THEN** the rendered HTML SHALL NOT include og:image or twitter:image meta tags

