// Package builder orchestrates site building from content files.
//
// This package coordinates the complete build pipeline:
//
//  1. Scan content directory for markdown files
//  2. Parse pages and blog posts with frontmatter
//  3. Render HTML using templates
//  4. Write output with clean URLs
//  5. Copy static assets
//
// # Clean URLs
//
// Pages are output with clean URL structure:
//
//	content/about.md     → public/about/index.html
//	content/blog/post.md → public/blog/post/index.html
//
// # Usage
//
//	b := builder.New(cfg)
//	b.SetAssetsDir("assets")
//	if err := b.Build(); err != nil {
//		return err
//	}
//
// Use [Builder.ScanContent] to scan without building, or [Builder.Build]
// for the complete build pipeline.
package builder
