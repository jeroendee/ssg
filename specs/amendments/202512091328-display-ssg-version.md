# 202512091328 - Technical amendment specification

## Objective
- using `ssg version` => output version information
- the version information shouldn't be semver. But a sha from a commit. 
- using ssg to build a site output the shortened version sha in the footer

```
# example footer
<span>Made with ❤️ in Amsterdam ({sha})</span>
```

## Technology & Design

- fetch and read: https://pkg.go.dev/github.com/spf13/cobra
- fetch and read: https://cobra.dev/docs/how-to-guides/working-with-commands/
- fetch and read: https://cobra.dev/docs/how-to-guides/working-with-flags/
- have a look in `/Users/jeroendee/ACID/40-49/44/44.W/watcher` how the version is implemented. 

### Skills

Use these Skills

- ...

## Examples
