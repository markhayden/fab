# PR body

Read this when build drafts a pull request. Skip preambles; use the repo's glossary terms.

```markdown
## Summary
<the smallest view that shows the change>

## Evidence
- **Before:** <failing check output, old behavior, or screenshot>
- **After:** <passing check output, new behavior, or screenshot>

## Merge risk
**Door:** one-way or two-way. <what makes it hard to undo, if one-way>
**Blast radius:** <one word: a module, a service, all users> <what could break>
```

## Summary: pick one view

- **Pseudocode** for logic or an algorithm.
- **A call tree** for control flow (`submitForm` → `createSession` → `persistPrompt`).
- **A component tree** for UI structure, with the state and boundaries that matter.
- **A shallow file tree** for what each file now owns, or a broad refactor.
- **A diagram** (Mermaid) for interaction between parts.
- **A `diff` of that view** when the point is what changed in a shape that already exists.

## Evidence

Take it from what the run already recorded: `build.md`'s Checks (failed before, passed after), `tests.md`'s
results, and `debug-*.md` loops. Don't paste whole logs; quote the lines that show the difference.

## Merge risk

A **one-way door** can't be cleanly reverted after merge: a data migration that drops or rewrites data, a published
API or package version, a message already sent. Everything else is a two-way door.

Adapted from Matt Pocock's `pr` skill ([mattpocock/skills](https://github.com/mattpocock/skills), MIT), which
credits Dex Horthy's `show-me` skill ([humanlayer/skills](https://github.com/humanlayer/skills)).
