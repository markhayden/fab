<p align="center"><img src="docs/assets/fab.png" alt="fab: a red toolbox holding a blueprint, a code window, a gear and a wrench" width="384"></p>

# fab

Dev skills and recipes for Claude Code and Codex, from one `skills/` directory.

**Skills** are building blocks: one job each. **Recipes** chain skills into a workflow with approval gates, and only run
when you call them. Every name is identical in both hosts, only the sigil differs: `/fab:<name>` in Claude Code,
`$fab:<name>` in Codex.

## Install

Claude Code:

```sh
claude plugin marketplace add markhayden/fab
claude plugin install fab@fab
```

Codex:

```sh
codex plugin marketplace add markhayden/fab
codex plugin add fab@fab
```

## Update

Claude Code (then restart the session):

```sh
claude plugin marketplace update fab
claude plugin update fab@fab
```

Codex:

```sh
codex plugin marketplace upgrade fab
```

## Skills

| Skill | What it does |
|---|---|
| `docs` | Keeps public docs true after a change (README, guides, API reference, changelog, examples), adds what users now need, and runs the commands and examples it documents wherever that's safe. |
| `grill` | Interviews you until intent and every design decision are settled, then writes `grill.md` for the spec. |
| `knowledge` | Keeps internal docs true (knowledge base, AGENTS.md, ADRs, glossary, specs and plans) without rewriting dated records, and carries a run's lasting decisions out of `.fab/`. |
| `pitch` | Pitches your idea back to you as its honest advocate; you challenge it until you call it sold, killed, or parked. A sold pitch seeds `grill`. |
| `plan` | Breaks the approved spec (and design) into `plan.md`: vertical-slice tasks with checks that fail before and pass after, one commit each, and rollback checkpoints. |
| `retro` | Looks back over a finished run and proposes environment fixes (checks for mechanical mistakes, standards for judgment calls, pointers, preferences), then applies the ones you pick and writes `retro.md`. |
| `review` | Adversarial review by the other host's agent (Codex from Claude Code, Claude Code from Codex): a spec reviewer and a standards reviewer each round, every finding fixed or rebutted, at least 2 rounds and at most 5. `build` runs it before its last checkpoint. |
| `simplify` | Cleans up a diff or an area without changing behavior (duplication, dead weight, needless layers), one commit per change (tests untouched for safe ones), and deepens shallow modules you pick. Writes `simplify-<target>.md`. |
| `spec` | Turns the settled decisions and user stories into `spec.md`: testable success criteria, scope, implementation and testing decisions, boundaries. |
| `stories` | Writes the user stories for settled work into `stories.md`: one per behavior, each with a concrete example. Skips (with your OK) when nothing observable changes. |
| `build` | Carries out the approved plan task by task: checks first, one commit per task, a stop at each checkpoint, a fresh-eyes review before the last, and `build.md` with SHAs and deviations. |
| `debug` | Finds and fixes a bug's root cause: first one fast command that fails on the exact symptom, then a shrunk case, ranked hypotheses, a fix with a regression test, and `debug-<bug>.md`. |
| `design` | Designs the spec's visual parts in a connected design tool (such as Paper) or from mockups you provide, and indexes them in `design.md`. Skips when there's nothing visual or nothing to design with. |
| `test` | Verifies a build against its spec and design: traces every criterion to a passing test, fills gaps, rewrites weak tests, checks UI with a browser tool, and writes `tests.md`. Bugs it finds go to `debug`. |
| `ui-kit` | Keeps UI built from the project's shared components and tokens: tells plan what to reuse, offers one-off UI to `simplify` to become shared components, and syncs the style guide (Storybook, a Paper page) after a run. |
| `voice` | Writes or rewrites content in your voice, from your private voice profile. Offers to build the profile if you don't have one. |

## Recipes

Recipes are call-only: they run only when you name them (`/fab:bigdev` in Claude Code, `$fab:bigdev` in Codex).

| Recipe | What it does |
|---|---|
| `bigdev` | The full flow, from interview to an open pull request: grill, stories, spec (reviewed), UI reuse, design, plan (reviewed), build (reviewed), simplify, test, then internal docs, public docs and style guide upkeep at one closing gate, a retro offered, then the PR with CI watched. |
| `tinydev` | The light flow for small, clear changes: questions only for real gaps, a short spec you approve, then plan, build (reviewed), simplify, test, and the same docs and style guide upkeep, all shown at an end gate, then the PR with CI watched. Bug reports go through `debug`. Offers to hand over to `bigdev` when the work turns out bigger. |

## Make it yours

fab reads optional private files from `~/.config/fab/`. They live outside the plugin, so updates never touch them.

| File | Purpose |
|---|---|
| `voice.md` | Your voice profile, used by `voice`. Start from `skills/voice/references/voice-template.md`. |
| `preferences.md` | Rules for every skill, e.g. "no compatibility shims on personal projects". |
| `skills/<name>.md` | Extra instructions for one skill or recipe. |

A per-skill file wins over `preferences.md`, which wins over fab's defaults. An override can change how a skill works,
but if it conflicts with an approval gate, the agent asks instead of skipping the gate. Write overrides as outcomes
("always include a rollback section"), not step numbers, so they keep working as fab changes. Keeping these files in
sync across machines is up to you (a private dotfiles repo works).

`~/.config/fab/` is outside your project, so Claude Code asks before reading it. To skip the prompt, allow it once in
`~/.claude/settings.json`: `"permissions": { "allow": ["Read(~/.config/fab/**)"] }`. In headless runs (`claude -p`),
pass `--add-dir ~/.config/fab`; without access the skill warns and falls back to defaults.

## Docs

- [Principles](docs/principles.md): why fab is built the way it is.
- [Conventions](docs/conventions.md): the rules every skill follows.
- [Development](docs/development.md): working on fab itself.
- [Skill audit](docs/skill-audit.md): where each skill came from, and what's next.

## Credits

Many ideas are adapted from [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) and
[mattpocock/skills](https://github.com/mattpocock/skills) (both MIT). Every fab skill
is rewritten through its own review; nothing is vendored.

## License

MIT, Copyright (c) 2026 WaffleHammock.
