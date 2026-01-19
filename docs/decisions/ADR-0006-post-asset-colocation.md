# ADR-0006: Post Asset Colocation

> Status: Accepted
> Date: 2026-01-19

## Context

Blog posts often include images and other assets. Options for organizing these:

1. **Global assets folder**: All images in `assets/images/`
2. **Post-specific folders**: Each post has its own directory with assets
3. **Shared blog assets**: `content/blog/assets/` with references from posts
4. **External hosting**: CDN or image hosting service

Considerations:
- Content portability (moving posts between sites)
- Reference simplicity (short paths in Markdown)
- Output organization (where assets end up)
- Namespace conflicts (duplicate filenames)

## Decision

Use a **shared blog assets folder** with path rewriting for output colocation.

**Source Structure**:
```
content/
└── blog/
    ├── 2024-01-15-my-post.md
    └── assets/
        ├── photo1.jpg      # Used by my-post
        └── diagram.png     # Used by my-post
```

**Markdown Reference**:
```markdown
![Photo](assets/photo1.jpg)
```

**Output Structure**:
```
public/
└── blog/
    └── my-post/
        ├── index.html
        ├── photo1.jpg      # Copied from assets/
        └── diagram.png     # Copied from assets/
```

**HTML Output** (path rewritten):
```html
<img src="photo1.jpg" alt="Photo">
```

## Consequences

### Positive

- **Simple references**: Just `assets/filename` in Markdown
- **Output colocation**: Assets live with their posts in output
- **Portable posts**: Copy post directory and it works
- **Clean URLs**: Assets at `/blog/post/image.jpg`
- **No conflicts**: Each post gets its own output directory

### Negative

- **Shared source folder**: All post assets in one directory
- **Manual tracking**: Must know which assets belong to which post
- **Path rewriting**: Adds build complexity
- **Source/output mismatch**: Structure differs between input and output

### Neutral

- Assets only copied if referenced in post
- Unused assets in `content/blog/assets/` are ignored
- Missing referenced assets cause build failure

## Implementation Notes

### Asset Reference Detection

Extract asset references from Markdown using pattern:
```
!\[.*?\]\((assets/[^)]+)\)
```

This matches:
- `![alt](assets/image.jpg)`
- `![](assets/nested/path/file.png)`

### Path Rewriting

During rendering, rewrite asset paths in HTML:
- Find: `src="assets/filename.ext"`
- Replace: `src="filename.ext"`

### Asset Copying

For each post:
1. Parse Markdown for asset references
2. Validate all referenced assets exist
3. Copy referenced assets to post output directory
4. Assets not referenced are not copied

### Validation

- Reference to non-existent asset → Build error
- Include asset path in error message
- Check before any file writing begins

### Supported Asset Types

Any file type is supported:
- Images: `.jpg`, `.png`, `.gif`, `.svg`, `.webp`
- Documents: `.pdf`
- Other: Any file referenced in Markdown

## References

- ADR-0001-clean-urls.md (output directory structure)
- CONTENT-MODEL.md (Post.Assets field)
