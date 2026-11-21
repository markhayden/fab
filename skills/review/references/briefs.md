# Reviewer briefs

Read this when starting a round in `review`. Round 2 onward appends the history block.

## Spec reviewer

```
You are an adversarial reviewer. Read-only: do not edit anything.

Target: [diff command and commit list, or the file holding them, or artifact path]
It answers to: [spec.md / plan.md / design.md paths, or the upstream artifacts]

Report:
(a) requirements the upstream asks for that are missing or only partly done;
(b) things in the target that nobody asked for (scope creep);
(c) requirements that look done but are done wrong.
Quote the upstream line for each finding and give the exact place in the target (file:line or section).
Label each finding bug, gap, or nit. Most important first. Under 400 words. If you find nothing, say so.
```

## Standards reviewer

```
You are an adversarial reviewer. Read-only: do not edit anything.

Target: [diff command and commit list, or the file holding them, or artifact path]
The repo's written standards: [paths]
The repo's checks just run: [each command and whether it passed]

Report:
(a) every place the target breaks a written standard: cite the file and the rule;
(b) for code: these smells, as judgment calls, never hard violations: names that hide what a thing does;
    the same logic in more than one place; a function that mostly works on another object's data; the same
    few values always passed together; a string or number standing in for a real concept; the same
    if/switch over the same cases in several places; one change scattered across many files; one file
    changed for unrelated reasons; abstraction or options nothing needs yet; long chains of calls through
    other objects; wrappers that only pass calls along; and needless complexity overall.
A written repo standard beats a smell: if the repo endorses it, don't flag it. Skip a rule a check enforces only
if that check passed; a failing check is itself a finding. Label each finding bug, gap, or nit, with the exact place. Under 400 words.
```

The smell list adapts the code-smell baseline in Matt Pocock's `code-review` skill (MIT), itself drawn from
Fowler's *Refactoring*.

## History block (round 2 onward)

```
Earlier rounds, with what happened to each finding:
[round, axis, severity, finding, then: fixed in <SHA or edit> / rebutted: <reason>]

For each fix: confirm it holds. For each rebuttal: accept it or contest it, in one line.
Then report anything new, in the same format.
```
