// Package model defines domain types for the static site generator.
//
// This package contains the core data structures used throughout ssg:
//
//   - [Config] holds site configuration loaded from ssg.yaml
//   - [Site] represents the complete site with pages and posts
//   - [Site.FaviconMIMEType] returns MIME type based on favicon extension
//   - [Page] represents a static page with title and content
//   - [Post] extends Page with date and summary for blog posts
//   - [NavItem] represents a navigation menu entry
//
// These types form the foundation for content processing and rendering.
package model
