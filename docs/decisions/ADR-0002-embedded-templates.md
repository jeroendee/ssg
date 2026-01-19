# ADR-0002: Embedded Templates

> Status: Accepted
> Date: 2026-01-19

## Context

Static site generators need templates to render HTML. These templates can be:

1. **External files**: Shipped separately, loaded at runtime
2. **Embedded in binary**: Compiled into the executable
3. **User-provided only**: No defaults, user must supply all templates

Considerations:
- Ease of use (zero-config experience)
- Customization flexibility
- Distribution simplicity
- Maintenance of templates vs binary

## Decision

**Embed default templates in the binary** with optional user override capability.

**Behavior**:
1. Templates are compiled into the executable
2. SSG works out-of-the-box without any template files
3. Users can override templates (future enhancement)
4. Default templates provide a complete, styled experience

## Consequences

### Positive

- **Zero configuration**: Works immediately after installation
- **Single binary**: No external dependencies or files needed
- **Consistent behavior**: Same templates across all installations
- **Simpler distribution**: Just ship one executable
- **Version consistency**: Templates match SSG version

### Negative

- **Larger binary**: Templates increase executable size
- **Rebuild required**: Template changes require recompilation
- **Less flexibility**: Users can't easily tweak defaults
- **Opinionated design**: Default aesthetic may not suit everyone

### Neutral

- Template overrides can be added later without breaking compatibility
- Embedded templates act as documentation for expected variables

## Implementation Notes

- Embed templates at compile time using language-appropriate mechanism
- Include complete set: base, blog post, blog list, 404, partials
- Embed default stylesheet alongside templates
- User-provided `style.css` replaces default entirely
- Future: Look for templates in `templates/` directory before using embedded

## References

- ADR-0005-no-javascript.md (templates output pure HTML/CSS)
