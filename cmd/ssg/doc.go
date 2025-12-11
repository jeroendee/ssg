// Ssg generates static websites from markdown content.
//
// Ssg is a command-line tool that converts markdown files with frontmatter
// into a complete static website with HTML pages, RSS feeds, and sitemaps.
//
// # Commands
//
// Build the site from markdown content:
//
//	ssg build [flags]
//
// Start a development server:
//
//	ssg serve [flags]
//
// Display version information:
//
//	ssg version
//
// # Version Information
//
// The version and build date are set at build time via -ldflags:
//
//	go build -ldflags="-X main.Version=v1.0.0 -X main.BuildDate=2024-01-15"
//
// If not set, Version defaults to "dev" and BuildDate to "unknown".
//
// # Configuration
//
// Commands read configuration from ssg.yaml by default. Use --config to
// specify a different configuration file. See the config package for
// configuration structure and required fields.
package main
