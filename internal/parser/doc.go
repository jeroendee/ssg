// Package parser handles markdown content parsing with frontmatter extraction.
//
// This package converts markdown files into [model.Page] and [model.Post]
// structures. It supports YAML frontmatter for metadata and converts
// markdown content to HTML using [github.com/gomarkdown/markdown].
//
// # Frontmatter Format
//
// Frontmatter is delimited by "---" and contains YAML:
//
//	---
//	title: My Page Title
//	---
//
//	Markdown content here...
//
// # Blog Post Dates
//
// Blog post dates are extracted from filenames in the format:
//
//	YYYY-MM-DD-slug.md
//
// For example, "2024-01-15-hello-world.md" produces a post with
// date 2024-01-15 and slug "hello-world".
//
// Use [ParsePage] for static pages and [ParsePost] for blog posts.
package parser
