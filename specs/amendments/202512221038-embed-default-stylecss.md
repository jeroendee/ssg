# 202512221038 - Technical amendment specification

## Objective

Embed a style.css file as the default style for the static generated site.

- `ssg build` should place the embedded style.css in the build output directory (public/)
- if a style.css is present in the configured `assets` build directory, use this style.css. This overrides the embedded default style.css file. 

## Technology & Design

Go embed directive is used for the default style.css

### Skills

Use these Skills

- ...

## Examples

Scenario: only embedded style.css file exists
Given the build output is set to a directory named: `public`
And no sibling assets directory with a style.css file exists
When the site is built
Then the `public` directory will contain the embedded style.css file

Scenario: a specific style.css file available in `assets` directory overrides the embedded style.css
Given a site project directory contains an `assets` directory with a style.css
And the build output is configured as `public`
When the site is built
Then the `public` directory contains the style.css from the `assets` directory
