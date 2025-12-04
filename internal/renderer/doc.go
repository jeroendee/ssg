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
package renderer
