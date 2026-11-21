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
| `grill` | Interviews you until intent and every design decision are settled, then writes `grill.md` for the spec. |
| `plan` | Breaks the approved spec (and design) into `plan.md`: vertical-slice tasks with checks that fail before and pass after, one commit each, and rollback checkpoints. |
| `spec` | Turns the settled decisions and user stories into `spec.md`: testable success criteria, scope, implementation and testing decisions, boundaries. |
| `stories` | Writes the user stories for settled work into `stories.md`: one per behavior, each with a concrete example. Skips (with your OK) when nothing observable changes. |
| `build` | Carries out the approved plan task by task: checks first, one commit per task, a stop at each checkpoint, a fresh-eyes review before the last, and `build.md` with SHAs and deviations. |
| `design` | Designs the spec's visual parts in a connected design tool (such as Paper) or from mockups you provide, and indexes them in `design.md`. Skips when there's nothing visual or nothing to design with. |
| `test` | Verifies a build against its spec and design: traces every criterion to a passing test, fills gaps, rewrites weak tests, checks UI with a browser tool, and writes `tests.md`. Proves bugs with a failing test first. |
| `voice` | Writes or rewrites content in your voice, from your private voice profile. Offers to build the profile if you don't have one. |

## Recipes

None shipped yet. `bigdev` (the full grill → stories → spec → design → plan → build → test flow) and `tinydev` (a lighter flow for small changes) are in
progress.

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
