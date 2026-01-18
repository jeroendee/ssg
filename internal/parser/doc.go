// Package parser handles markdown content parsing with frontmatter extraction.
//
// This package converts markdown files into [model.Page] and [model.Post]
// structures. It supports YAML frontmatter for metadata and converts
// markdown content to HTML using [github.com/yuin/goldmark].
//
// # Frontmatter Format
//
// Frontmatter is delimited by "---" and contains YAML. The [Frontmatter]
// struct supports these fields:
//
//	---
//	title: My Page Title
//	summary: Brief description for RSS feeds
//	date: 2024-01-15
//	---
//
//	Markdown content here...
//
// Use [ExtractFrontmatter] to separate frontmatter from content, and
// [MarkdownToHTML] to convert the markdown body to HTML.
//
// # Blog Post Dates
//
// Blog post dates are extracted from filenames in the format:
//
//	YYYY-MM-DD-slug.md
//
// For example, "2024-01-15-hello-world.md" produces a post with
// date 2024-01-15 and slug "hello-world". A date in frontmatter
// overrides the filename date.
//
// Use [ParsePage] for static pages and [ParsePost] for blog posts.
package parser
