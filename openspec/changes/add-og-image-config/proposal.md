# Change: Add Dedicated OG Image Configuration

## Why
LinkedIn and other social platforms don't support SVG format for Open Graph images. Currently, the site logo is used as og:image fallback, but when the logo is SVG (common for scalable web logos), social sharing previews fail to display an image.

## What Changes
- Add `ogImage` field to Site and Config models
- Add `site.ogImage` YAML config option
- Update `ogImageURL()` to prefer `OGImage` with `Logo` fallback
- Allow PNG/JPG images optimized for social sharing (1200x627 recommended)

## Impact
- Affected specs: social-sharing (new capability)
- Affected code:
  - `internal/model/model.go` - Add OGImage field
  - `internal/config/config.go` - Parse from YAML
  - `internal/renderer/renderer.go` - Update ogImageURL()
