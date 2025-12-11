// Package renderer handles HTML template rendering for site generation.
//
// This package uses Go's [html/template] with embedded templates to produce
// HTML pages. Templates follow a BearBlog-inspired design with responsive
// layout and dark mode support.
//
// # Templates
//
// The package provides templates for:
//
//   - Homepage (index.html)
//   - Static pages (/about/, /contact/, etc.)
//   - Blog listing (/blog/)
//   - Individual blog posts (/blog/slug/)
//   - 404 error page
//
// All templates include navigation, consistent styling, and link to /style.css.
//
// # Usage
//
// Create a renderer and call the appropriate Render method:
//
//	r, err := renderer.New()
//	if err != nil {
//		return err
//	}
//	html, err := r.RenderPage(site, page)
//
// Templates are embedded at compile time and require no external files.
//
// # Methods
//
// SetVersion sets the version string to be included in rendered templates.
//
// RenderBase renders the base template with site data and content.
// It accepts raw HTML content and wraps it with the site's base template.
//
// RenderFeed generates an RSS 2.0 feed from blog posts. It includes up to
// 20 most recent posts with full content wrapped in CDATA sections.
package renderer
