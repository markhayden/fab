# Records

Read this before writing an ADR, a status line, or a glossary entry in fab:knowledge. A repo convention always wins.

## ADR

Default location `docs/adr/`, files `0001-short-title.md`, numbers continuing the sequence.

```markdown
# 0004. Short title of the decision

- Status: accepted (2026-11-21)
- Supersedes: 0002-short-title (only when it replaces one)

## Context
What forced a choice, and the constraints that mattered.

## Decision
What was chosen, in one or two sentences.

## Alternatives
Each real option, and why it lost.

## Consequences
What this makes easier, what it makes harder, and what would make us revisit it.
```

Never delete or rewrite an accepted ADR. When a decision changes, write a new one that names the old one.

## Status line

A dated record (ADR, dated spec or plan) only ever has its status changed, in place (add one under the title if
it has none); nothing else in it is edited. One of:

```markdown
- Status: superseded by 0007-event-log (2026-11-21)
- Status: outdated as of 4f2a9c1: exports moved to the jobs queue
```

Use "superseded" only when the replacing record lands in the same commit; otherwise "outdated as of".

## Glossary

One `GLOSSARY.md` at the root; a repo with several distinct areas uses `GLOSSARY-MAP.md` pointing to one per area.

```markdown
# <Project or area>

One or two sentences on what this area is.

## Language

**Order**: a customer's request to buy one or more items, from checkout until delivery.
_Avoid_: purchase, transaction
```

- Pick one word per concept and list the others under _Avoid_.
- Define what the thing is, in one or two sentences, with no implementation detail.
- Only terms specific to this project; general programming words don't belong.

ADR shape and lifecycle adapted from Addy Osmani's `documentation-and-adrs`
([addyosmani/agent-skills](https://github.com/addyosmani/agent-skills), MIT); the three-part ADR test and glossary
format from Matt Pocock's `domain-modeling` ([mattpocock/skills](https://github.com/mattpocock/skills), MIT).
