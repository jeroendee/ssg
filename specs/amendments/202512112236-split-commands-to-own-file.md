# 202512112236 - Technical amendment specification

## Objective

The different subcommands should live in their own file in the cmd/ssg directory. 

## Technology & Design

- fetch and read: https://cobra.dev/docs/how-to-guides/working-with-commands/ 
```
## How to Organize Commands in Packages
Small apps are fine with all commands in one package (cmd/). Large apps benefit from modular packages where each feature exposes a constructor (e.g., NewCommand()). This keeps imports narrow, improves compile times, and lets subtrees evolve independently—an approach used by bigger projects like Hugo.

Two layouts you can choose from:

- Simple (default Cobra): one cmd package, one file per command. => use this.
- Modular (recommended at scale): each feature gets its own package that returns a *cobra.Command.
```

### Skills

Use these Skills

- ...

## Examples
