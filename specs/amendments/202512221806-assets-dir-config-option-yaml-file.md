# 202512221806 - Technical amendment specification

## Objective

In the current ssg.yaml config file, the section: `build` has the config options:
- `content`
- `output`

This section should also get the option:
- `assets` => directory where the assets are stored. 

```yaml
build:
  assets: assets
  content: content
  output: public
```

This way the build options are all available in the yaml config file. 

## Technology & Design

### Skills

Use these Skills

- ...

## Examples

EXAMPLE project site directory structure:

```tree
.
├── assets
│   ├── favicon.ico
│   ├── quality-shepherd-logo.png
│   ├── quality-shepherd-logo.svg
│   └── style.css
├── content
│   ├── about.md
│   ├── blog
│   │   ├── 2021-03-26-hello-world.md
│   │   ├── 2023-06-15-testing-in-go.md
│   │   └── 2025-12-08-test-post.md
│   └── contact.md
├── public
│   ├── 404.html
│   ├── about
│   │   └── index.html
│   ├── blog
│   │   ├── hello-world
│   │   │   └── index.html
│   │   ├── index.html
│   │   ├── test-post
│   │   │   └── index.html
│   │   └── testing-in-go
│   │       └── index.html
│   ├── contact
│   │   └── index.html
│   ├── favicon.ico
│   ├── feed
│   │   └── index.xml
│   ├── index.html
│   ├── quality-shepherd-logo.png
│   ├── quality-shepherd-logo.svg
│   ├── robots.txt
│   ├── sitemap.xml
│   └── style.css
└── ssg.yaml
```
