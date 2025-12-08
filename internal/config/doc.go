// Package config handles loading and validating site configuration.
//
// Configuration is loaded from a YAML file (typically ssg.yaml) with the
// following structure:
//
//	site:
//	  title: My Site
//	  description: A brief description of my site
//	  baseURL: https://example.com
//	  author: Author Name
//	  logo: /images/logo.png
//	  favicon: /images/favicon.ico
//	build:
//	  content: content
//	  output: public
//	navigation:
//	  - title: Home
//	    url: /
//
// Use [Load] to read configuration from a file, or [LoadWithOptions] to
// apply CLI flag overrides for content and output directories.
//
// Required fields are site.title and site.baseURL. Default values are
// applied for content ("content") and output ("public") directories.
package config
